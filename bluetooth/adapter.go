/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package bluetooth

import (
	"fmt"
	bluez "github.com/linuxdeepin/go-dbus-factory/org.bluez"
	"os"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"time"
)

type adapter struct {
	core    *bluez.HCI
	address string

	Path                dbus.ObjectPath
	Name                string
	Alias               string
	Powered             bool
	Discovering         bool
	Discoverable        bool
	DiscoverableTimeout uint32
	// discovering timer, when time is up, stop discovering until start button is clicked next time
	discoveringTimeout *time.Timer
}

var defaultDiscoveringTimeout = 1 * time.Minute

func newAdapter(systemSigLoop *dbusutil.SignalLoop, apath dbus.ObjectPath) (a *adapter) {
	a = &adapter{Path: apath}
	systemConn := systemSigLoop.Conn()
	a.core, _ = bluez.NewHCI(systemConn, apath)
	a.core.InitSignalExt(systemSigLoop, true)
	a.connectProperties()
	a.address, _ = a.core.Address().Get(0)
	// 用于定时停止扫描
	a.discoveringTimeout = time.AfterFunc(defaultDiscoveringTimeout, func() {
		logger.Debug("discovery time out, stop discovering")
		if err := a.core.StopDiscovery(0); err != nil {
			logger.Warningf("stop discovery failed, err:%v", err)
		}
	})
	// stop timer at first
	a.discoveringTimeout.Stop()
	// fix alias
	alias, _ := a.core.Alias().Get(0)
	if alias == "first-boot-hostname" {
		hostname, err := os.Hostname()
		if err == nil {
			if hostname != "first-boot-hostname" {
				// reset alias
				err = a.core.Alias().Set(0, "")
				if err != nil {
					logger.Warning(err)
				}
			}
		} else {
			logger.Warning("failed to get hostname:", err)
		}
	}

	a.Alias, _ = a.core.Alias().Get(0)
	a.Name, _ = a.core.Name().Get(0)
	a.Powered, _ = a.core.Powered().Get(0)
	a.Discovering, _ = a.core.Discovering().Get(0)
	a.Discoverable, _ = a.core.Discoverable().Get(0)
	a.DiscoverableTimeout, _ = a.core.DiscoverableTimeout().Get(0)
	return
}

func (a *adapter) destroy() {
	a.core.RemoveHandler(proxy.RemoveAllHandlers)
}

func (a *adapter) String() string {
	return fmt.Sprintf("adapter %s [%s]", a.Alias, a.address)
}

func (a *adapter) notifyAdapterAdded() {
	logger.Info("AdapterAdded", a)
	globalBluetooth.service.Emit(globalBluetooth, "AdapterAdded", marshalJSON(a))
	globalBluetooth.updateState()
}

func (a *adapter) notifyAdapterRemoved() {
	logger.Info("AdapterRemoved", a)
	globalBluetooth.service.Emit(globalBluetooth, "AdapterRemoved", marshalJSON(a))
	globalBluetooth.updateState()
}

func (a *adapter) notifyPropertiesChanged() {
	globalBluetooth.service.Emit(globalBluetooth, "AdapterPropertiesChanged", marshalJSON(a))
	globalBluetooth.updateState()
}

func (a *adapter) connectProperties() {
	a.core.Name().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		a.Name = value
		logger.Debugf("%s Name: %v", a, value)
		a.notifyPropertiesChanged()
	})

	a.core.Alias().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		a.Alias = value
		logger.Debugf("%s Alias: %v", a, value)
		a.notifyPropertiesChanged()
	})
	a.core.Powered().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		a.Powered = value
		logger.Debugf("%s Powered: %v", a, value)
		// set adapter powered config
		globalBluetooth.config.setAdapterConfigPowered(a.address, value)
		// check if powered state is true
		if a.Powered {
			// set discoverable time out
			err := a.core.DiscoverableTimeout().Set(0, 300)
			if err != nil {
				logger.Warningf("failed to set discoverable time out for %s: %v", a, err)
			}
			// set discoverable
			err = a.core.Discoverable().Set(0, true)
			if err != nil {
				logger.Warningf("failed to set discoverable for %s: %v", a, err)
			}
			// start discovery
			err = a.core.StartDiscovery(0)
			if err != nil {
				logger.Warningf("failed to start discovery for %s: %v", a, err)
			}
			// dont need to start discovering, according to blueZ, scan will be called, when power is set on
			a.discoveringTimeout.Reset(defaultDiscoveringTimeout)
			// sleep for 2 seconds in case user click paired device, but blueZ service is not ready
			time.Sleep(2 * time.Second)
			a.notifyPropertiesChanged()
			// in case auto connect to device failed, only when signal power on is received, try to auto connect device
			go globalBluetooth.tryConnectPairedDevices(a.Path)
		} else {
			// if power off, stop discovering time out
			a.discoveringTimeout.Stop()
			// sleep for 2 seconds in case
			time.Sleep(2 * time.Second)
			a.notifyPropertiesChanged()
		}
	})
	a.core.Discovering().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		a.Discovering = value
		logger.Debugf("%s Discovering: %v", a, value)
		a.notifyPropertiesChanged()
	})
	a.core.Discoverable().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		a.Discoverable = value
		logger.Debugf("%s Discoverable: %v", a, value)
		a.notifyPropertiesChanged()
	})
	a.core.DiscoverableTimeout().ConnectChanged(func(hasValue bool, value uint32) {
		if !hasValue {
			return
		}
		a.DiscoverableTimeout = value
		logger.Debugf("%s DiscoverableTimeout: %v", a, value)
		a.notifyPropertiesChanged()
	})
}
