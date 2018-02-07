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
	"dbus/org/bluez"
	"fmt"
	"pkg.deepin.io/lib/dbus"
)

type adapter struct {
	bluezAdapter *bluez.Adapter1

	Path                dbus.ObjectPath
	adddress            string
	Name                string
	Alias               string
	Powered             bool
	Discovering         bool
	Discoverable        bool
	DiscoverableTimeout uint32
}

func newAdapter(apath dbus.ObjectPath) (a *adapter) {
	a = &adapter{Path: apath}
	a.bluezAdapter, _ = bluezNewAdapter(apath)
	a.connectProperties()
	a.adddress = a.bluezAdapter.Address.Get()
	a.Alias = a.bluezAdapter.Alias.Get()
	a.Name = a.bluezAdapter.Name.Get()
	a.Powered = a.bluezAdapter.Powered.Get()
	a.Discovering = a.bluezAdapter.Discovering.Get()
	a.Discoverable = a.bluezAdapter.Discoverable.Get()
	a.DiscoverableTimeout = a.bluezAdapter.DiscoverableTimeout.Get()
	return
}
func destroyAdapter(a *adapter) {
	bluezDestroyAdapter(a.bluezAdapter)
}

func (a *adapter) notifyAdapterAdded() {
	logger.Info("AdapterAdded", marshalJSON(a))
	dbus.Emit(bluetooth, "AdapterAdded", marshalJSON(a))
	bluetooth.setPropState()
}
func (a *adapter) notifyAdapterRemoved() {
	logger.Info("AdapterRemoved", marshalJSON(a))
	dbus.Emit(bluetooth, "AdapterRemoved", marshalJSON(a))
	bluetooth.setPropState()
}
func (a *adapter) notifyPropertiesChanged() {
	logger.Debug("AdapterPropertiesChanged", marshalJSON(a))
	dbus.Emit(bluetooth, "AdapterPropertiesChanged", marshalJSON(a))
	bluetooth.setPropState()
}
func (a *adapter) connectProperties() {
	a.bluezAdapter.Name.ConnectChanged(func() {
		a.Name = a.bluezAdapter.Name.Get()
		a.notifyPropertiesChanged()
	})
	a.bluezAdapter.Alias.ConnectChanged(func() {
		a.Alias = a.bluezAdapter.Alias.Get()
		a.notifyPropertiesChanged()
	})
	a.bluezAdapter.Powered.ConnectChanged(func() {
		a.Powered = a.bluezAdapter.Powered.Get()
		logger.Infof("adapter powered changed %#v", a)
		a.notifyPropertiesChanged()
	})
	a.bluezAdapter.Discovering.ConnectChanged(func() {
		a.Discovering = a.bluezAdapter.Discovering.Get()
		a.notifyPropertiesChanged()
	})
	a.bluezAdapter.Discoverable.ConnectChanged(func() {
		a.Discoverable = a.bluezAdapter.Discoverable.Get()
		a.notifyPropertiesChanged()
	})
	a.bluezAdapter.DiscoverableTimeout.ConnectChanged(func() {
		a.DiscoverableTimeout = a.bluezAdapter.DiscoverableTimeout.Get()
		a.notifyPropertiesChanged()
	})
}

func (b *Bluetooth) addAdapter(apath dbus.ObjectPath) {
	if b.isAdapterExists(apath) {
		logger.Error("repeat add adapter", apath)
		return
	}

	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()

	a := newAdapter(apath)
	// initialize adapter power state
	b.config.addAdapterConfig(bluezGetAdapterAddress(apath))
	oldPowered := b.config.getAdapterConfigPowered(bluezGetAdapterAddress(apath))
	b.SetAdapterPowered(apath, oldPowered)
	b.SetAdapterDiscoverable(apath, false)
	if oldPowered {
		b.RequestDiscovery(apath)
	}

	b.adapters[apath] = a
	a.notifyAdapterAdded()
}

func (b *Bluetooth) removeAdapter(apath dbus.ObjectPath) {
	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()

	if b.adapters[apath] == nil {
		logger.Error("repeat remove adapter", apath)
		return
	}

	b.doRemoveAdapter(apath)
}

func (b *Bluetooth) doRemoveAdapter(apath dbus.ObjectPath) {
	removeAdapter := b.adapters[apath]
	delete(b.adapters, apath)

	removeAdapter.notifyAdapterRemoved()
	destroyAdapter(removeAdapter)
}

func (b *Bluetooth) getAdapter(apath dbus.ObjectPath) (a *adapter, err error) {
	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()

	a = b.adapters[apath]
	if a == nil {
		err = fmt.Errorf("adapter not exists %s", apath)
		logger.Error(err)
		return
	}
	return
}
func (b *Bluetooth) isAdapterExists(apath dbus.ObjectPath) bool {
	b.adaptersLock.Lock()
	defer b.adaptersLock.Unlock()
	if b.adapters[apath] != nil {
		return true
	}
	return false
}

// GetAdapters return all adapter objects that marshaled by json.
func (b *Bluetooth) GetAdapters() (adaptersJSON string, err error) {
	v := make([]*adapter, 0, len(b.adapters))
	for _, a := range b.adapters {
		v = append(v, a)
	}
	adaptersJSON = marshalJSON(v)
	return
}

func (b *Bluetooth) RequestDiscovery(apath dbus.ObjectPath) (err error) {
	// if adapter is discovering now, just ignore
	if bluezGetAdapterDiscovering(apath) {
		return
	}

	b.SetAdapterDiscovering(apath, true)
	return
}

func (b *Bluetooth) SetAdapterPowered(apath dbus.ObjectPath, powered bool) (err error) {
	b.ClearUnpairedDevice()
	err = bluezSetAdapterPowered(apath, powered)
	if err == nil {
		// save the powered state
		b.config.setAdapterConfigPowered(bluezGetAdapterAddress(apath), powered)
		if powered {
			b.SetAdapterDiscovering(apath, true)
		}
	}
	return
}
func (b *Bluetooth) SetAdapterAlias(apath dbus.ObjectPath, alias string) (err error) {
	err = bluezSetAdapterAlias(apath, alias)
	return
}
func (b *Bluetooth) SetAdapterDiscoverable(apath dbus.ObjectPath, discoverable bool) (err error) {
	err = bluezSetAdapterDiscoverable(apath, discoverable)
	return
}
func (b *Bluetooth) SetAdapterDiscovering(apath dbus.ObjectPath, discoverable bool) (err error) {
	err = bluezSetAdapterDiscovering(apath, discoverable)
	return
}
func (b *Bluetooth) SetAdapterDiscoverableTimeout(apath dbus.ObjectPath, discoverableTimeout uint32) (err error) {
	err = bluezSetAdapterDiscoverableTimeout(apath, discoverableTimeout)
	return
}
