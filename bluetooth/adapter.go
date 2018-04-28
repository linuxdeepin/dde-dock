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

	"github.com/linuxdeepin/go-dbus-factory/org.bluez"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

type adapter struct {
	bluezAdapter *bluez.HCI

	Path                dbus.ObjectPath
	adddress            string
	Name                string
	Alias               string
	Powered             bool
	Discovering         bool
	Discoverable        bool
	DiscoverableTimeout uint32
}

func newAdapter(systemSigLoop *dbusutil.SignalLoop, apath dbus.ObjectPath) (a *adapter) {
	a = &adapter{Path: apath}
	systemConn := systemSigLoop.Conn()
	a.bluezAdapter, _ = bluez.NewHCI(systemConn, apath)
	a.bluezAdapter.InitSignalExt(systemSigLoop, true)
	a.connectProperties()
	a.adddress, _ = a.bluezAdapter.Address().Get(0)
	a.Alias, _ = a.bluezAdapter.Alias().Get(0)
	a.Name, _ = a.bluezAdapter.Name().Get(0)
	a.Powered, _ = a.bluezAdapter.Powered().Get(0)
	a.Discovering, _ = a.bluezAdapter.Discovering().Get(0)
	a.Discoverable, _ = a.bluezAdapter.Discoverable().Get(0)
	a.DiscoverableTimeout, _ = a.bluezAdapter.DiscoverableTimeout().Get(0)
	return
}

func (a *adapter) destroy() {
	a.bluezAdapter.RemoveHandler(proxy.RemoveAllHandlers)
}

func (a *adapter) String() string {
	return fmt.Sprintf("adapter %s [%s]", a.Alias, a.adddress)
}

func (a *adapter) notifyAdapterAdded() {
	logger.Info("AdapterAdded", a)
	globalBluetooth.service.Emit(globalBluetooth, "AdapterAdded", marshalJSON(a))
	globalBluetooth.updateState()
}

func (a *adapter) notifyAdapterRemoved() {
	logger.Info("AdapterRemoved", marshalJSON(a))
	globalBluetooth.service.Emit(globalBluetooth, "AdapterRemoved", marshalJSON(a))
	globalBluetooth.updateState()
}

func (a *adapter) notifyPropertiesChanged() {
	logger.Debug("AdapterPropertiesChanged", marshalJSON(a))
	globalBluetooth.service.Emit(globalBluetooth, "AdapterPropertiesChanged", marshalJSON(a))
	globalBluetooth.updateState()
}

func (a *adapter) connectProperties() {
	a.bluezAdapter.Name().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		a.Name = value
		a.notifyPropertiesChanged()
	})
	a.bluezAdapter.Alias().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}
		a.Alias = value
		a.notifyPropertiesChanged()
	})
	a.bluezAdapter.Powered().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		a.Powered = value
		logger.Infof("adapter powered changed %#v", a)
		a.notifyPropertiesChanged()
	})
	a.bluezAdapter.Discovering().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		a.Discovering = value
		a.notifyPropertiesChanged()
	})
	a.bluezAdapter.Discoverable().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}
		a.Discoverable = value
		a.notifyPropertiesChanged()
	})
	a.bluezAdapter.DiscoverableTimeout().ConnectChanged(func(hasValue bool, value uint32) {
		if !hasValue {
			return
		}
		a.DiscoverableTimeout = value
		a.notifyPropertiesChanged()
	})
}
