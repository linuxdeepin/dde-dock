/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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
	"dbus/org/freedesktop/miracle/wifi"
	"pkg.deepin.io/dde/daemon/iw"
	"pkg.deepin.io/lib/dbus"
	"strings"
	"sync"
	"time"
)

type LinkInfo struct {
	Name         string
	FriendlyName string
	MacAddress   string
	Managed      bool
	P2PScanning  bool
	Path         dbus.ObjectPath

	index  uint32
	core   *wifi.Link
	locker sync.Mutex
}
type LinkInfos []*LinkInfo

func newLinkInfo(dpath dbus.ObjectPath) (*LinkInfo, error) {
	link, err := wifi.NewLink(wifiDest, dpath)
	if err != nil {
		return nil, err
	}

	var info = &LinkInfo{
		index: link.InterfaceIndex.Get(),
		Path:  dpath,
		core:  link,
	}
	info.update()
	return info, nil
}

func destroyLinkInfo(link *LinkInfo) {
	if link.core == nil {
		return
	}
	link.locker.Lock()
	wifi.DestroyLink(link.core)
	link.core = nil
	link.locker.Unlock()
}

func (link *LinkInfo) update() {
	link.locker.Lock()
	defer link.locker.Unlock()
	if link.core == nil {
		return
	}
	link.Name = link.core.InterfaceName.Get()
	link.FriendlyName = link.core.FriendlyName.Get()
	link.MacAddress = link.core.MACAddress.Get()
	link.Managed = link.core.Managed.Get()
	link.P2PScanning = link.core.P2PScanning.Get()
}

func (link *LinkInfo) hasP2PSupported() bool {
	infos, err := iw.ListWirelessInfo()
	if err != nil {
		return false
	}
	return (infos.ListMiracastDevice().Get(link.MacAddress) != nil)
}

func (link *LinkInfo) EnableManaged(enabled bool) error {
	if enabled {
		return link.core.Manage()
	} else {
		return link.core.Unmanage()
	}
}

// SetName Must be set before scanning
func (link *LinkInfo) SetName(name string) {
	if link.core.FriendlyName.Get() == name {
		return
	}
	link.core.FriendlyName.Set(name)
	time.Sleep(time.Millisecond * 10)
	link.update()
}

func (link *LinkInfo) EnableP2PScanning(enabled bool) {
	if link.core.P2PScanning.Get() == enabled {
		return
	}
	link.core.P2PScanning.Set(enabled)
}

func (links LinkInfos) Get(dpath dbus.ObjectPath) *LinkInfo {
	if !isLinkObjectPath(dpath) {
		return nil
	}
	for _, link := range links {
		if link.Path == dpath {
			return link
		}
	}
	return nil
}

func (links LinkInfos) Add(dpath dbus.ObjectPath) (LinkInfos, error) {
	link := links.Get(dpath)
	if link != nil {
		link.update()
		return links, nil
	}
	tmp, err := newLinkInfo(dpath)
	if err != nil {
		return links, err
	}
	links = append(links, tmp)
	return links, nil
}

func (links LinkInfos) Remove(dpath dbus.ObjectPath) (LinkInfos, bool) {
	var (
		tmp    LinkInfos
		exists bool = false
	)
	for _, link := range links {
		if link.Path == dpath {
			exists = true
			destroyLinkInfo(link)
			continue
		}
		tmp = append(tmp, link)
	}
	return tmp, exists
}

func isLinkObjectPath(dpath dbus.ObjectPath) bool {
	return strings.Contains(string(dpath), linkPath)
}
