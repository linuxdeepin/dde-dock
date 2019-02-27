/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package miracast

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.miracle.wifi"
	"pkg.deepin.io/dde/daemon/iw"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

type LinkInfo struct {
	Name         string
	FriendlyName string
	MacAddress   string
	Managed      bool
	P2PScanning  bool
	Path         dbus.ObjectPath

	interfaceName string
	index         uint32
	core          *wifi.Link
	locker        sync.Mutex
	myManaged     bool
}

type LinkInfos []*LinkInfo

func newLinkInfo(objPath dbus.ObjectPath) (*LinkInfo, error) {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	link, err := wifi.NewLink(sysBus, objPath)
	if err != nil {
		return nil, err
	}

	index, err := link.InterfaceIndex().Get(0)
	if err != nil {
		return nil, err
	}

	var info = &LinkInfo{
		index: index,
		Path:  objPath,
		core:  link,
	}
	info.update()
	return info, nil
}

func (link *LinkInfo) connectSignal(m *Miracast) {
	link.core.InitSignalExt(m.sysSigLoop, true)
	err := link.core.Managed().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}

		link.locker.Lock()
		link.Managed = value
		link.locker.Unlock()

		logger.Debugf("link %s managed changed: %v", link.Path, value)
		if value {
			m.emitSignalEvent(EventLinkManaged, link.Path)
		} else {
			m.emitSignalEvent(EventLinkUnmanaged, link.Path)
		}

		if link.myManaged != value {
			logger.Warning("link.myManged != value")
			err := link.EnableManaged(link.myManaged)
			if err != nil {
				logger.Warning(err)
			}
		}
	})
	if err != nil {
		logger.Warning(err)
	}

	err = link.core.P2PScanning().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}

		link.locker.Lock()
		link.P2PScanning = value
		link.locker.Unlock()
	})
	if err != nil {
		logger.Warning(err)
	}
}

func destroyLinkInfo(link *LinkInfo) {
	link.core.RemoveHandler(proxy.RemoveAllHandlers)
}

func (link *LinkInfo) update() {
	link.locker.Lock()
	defer link.locker.Unlock()
	if link.core == nil {
		return
	}
	link.Name, _ = link.core.InterfaceName().Get(0)
	link.FriendlyName, _ = link.core.FriendlyName().Get(0)
	link.MacAddress, _ = link.core.MACAddress().Get(0)
	link.Managed, _ = link.core.Managed().Get(0)
	link.P2PScanning, _ = link.core.P2PScanning().Get(0)
	link.interfaceName, _ = link.core.InterfaceName().Get(0)
}

func (link *LinkInfo) hasP2PSupported() bool {
	infos, err := iw.ListWirelessInfo()
	if err != nil {
		return false
	}
	return infos.ListMiracastDevice().Get(link.MacAddress) != nil
}

func (link *LinkInfo) EnableManaged(enabled bool) error {
	managed, err := link.core.Managed().Get(0)
	if err != nil {
		return err
	}

	if managed == enabled {
		return nil
	}

	if enabled {
		return link.core.Manage(0)
	} else {
		return link.core.Unmanage(0)
	}
}

func (link *LinkInfo) waitManaged(wantManaged bool) error {
	name := fmt.Sprintf("link %s managed", link.Path)
	return waitChange(name, wantManaged, func() (b bool, err error) {
		return link.core.Managed().Get(0)
	})
}

func waitPeerConnected(peer *wifi.Peer, wantConnected bool) error {
	name := fmt.Sprintf("peer %s connected", peer.Path_())
	return waitChange(name, wantConnected, func() (b bool, err error) {
		return peer.Connected().Get(0)
	})
}

func waitChange(name string, wantVal bool, getValFn func() (bool, error)) error {
	max := int(defaultTimeout / defaultInterval)
	for i := 0; i < max; i++ {
		val, err := getValFn()
		if err != nil {
			return err
		}

		if val == wantVal {
			logger.Debugf("waitChange finish %s %v", name, wantVal)
			return nil
		}
		time.Sleep(defaultInterval)
		logger.Debug("wait tick", i, name, wantVal)
	}
	return fmt.Errorf("timeout wait %s %v", name, wantVal)
}

// SetName Must be set before scanning
func (link *LinkInfo) SetName(name string) {
	friendlyName, err := link.core.FriendlyName().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if friendlyName == name {
		return
	}
	err = link.core.FriendlyName().Set(0, name)
	if err != nil {
		logger.Warning(err)
		return
	}
	time.Sleep(time.Millisecond * 10)
	link.update()
}

func (link *LinkInfo) EnableP2PScanning(enabled bool) {
	scanning, err := link.core.P2PScanning().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if scanning == enabled {
		return
	}
	err = link.core.P2PScanning().Set(0, enabled)
	if err != nil {
		logger.Warning(err)
	}
}

func (link *LinkInfo) ConfigureForManaged() {
	err := link.core.FriendlyName().Set(0, os.Getenv("USER"))
	if err != nil {
		logger.Warning(err)
	}

	err = link.core.WfdSubelements().Set(0, "000600001c440036")
	if err != nil {
		logger.Warning(err)
	}
}

func (links LinkInfos) Get(objPath dbus.ObjectPath) *LinkInfo {
	if !isLinkObjectPath(objPath) {
		return nil
	}
	for _, link := range links {
		if link.Path == objPath {
			return link
		}
	}
	return nil
}

func (links LinkInfos) Remove(objPath dbus.ObjectPath) (LinkInfos, bool) {
	var (
		tmp    LinkInfos
		exists bool = false
	)
	for _, link := range links {
		if link.Path == objPath {
			exists = true
			destroyLinkInfo(link)
			continue
		}
		tmp = append(tmp, link)
	}
	return tmp, exists
}

func isLinkObjectPath(objPath dbus.ObjectPath) bool {
	return strings.HasPrefix(string(objPath), linkDBusPathPrefix)
}

func isPeerObjectPath(objPath dbus.ObjectPath) bool {
	return strings.HasPrefix(string(objPath), peerDBusPathPrefix)
}
