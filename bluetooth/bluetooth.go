/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package bluetooth

import (
	idbus "dbus/org/freedesktop/dbus/system"
	"dlib/dbus"
)

const (
	dbusBluezDest       = "org.bluez"
	dbusBluezPath       = "/org/bluez"
	dbusBluezIfsAdapter = "org.bluez.Adapter1"
	dbusBluezIfsDevice  = "org.bluez.Device1"

	dbusBluetoothDest = "com.deepin.daemon.Bluetooth"
	dbusBluetoothPath = "/com/deepin/daemon/Bluetooth"
	dbusBluetoothIfs  = "com.deepin.daemon.Bluetooth"
)

var bluezObjectManager *idbus.ObjectManager

type dbusObjectData map[string]dbus.Variant
type dbusInterfaceData map[string]map[string]dbus.Variant
type dbusInterfacesData map[dbus.ObjectPath]map[string]map[string]dbus.Variant

type Bluetooth struct {
	// adapter
	PrimaryAdapter string `access:"readwrite"` // do not use dbus.ObjectPath here due to could not be reawrite
	adapters       []*adapter
	Adapters       string // array of adapters that marshaled by json

	// device
	devices map[dbus.ObjectPath][]*device
	Devices string // device objects that marshaled by json

	// alias properties for primary adapter
	Alias               dbus.Property `access:"readwrite"`
	Powered             dbus.Property `access:"readwrite"`
	Discoverable        dbus.Property `access:"readwrite"`
	DiscoverableTimeout dbus.Property `access:"readwrite"`

	// signals
	DeviceAdded      func(devJSON string)
	DeviceRemoved    func(devJSON string)
	RequestPinCode   func(devJSON string)
	AuthorizeService func(devJSON string, uuid string)
}

func NewBluetooth() (bluettoth *Bluetooth) {
	bluettoth = &Bluetooth{}
	return
}

func (b *Bluetooth) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		dbusBluetoothDest,
		dbusBluetoothPath,
		dbusBluetoothIfs,
	}
}

func DestroyBluetooth(b *Bluetooth) {
	dbus.UnInstallObject(bluetooth)
}

func (b *Bluetooth) initBluetooth() {
	b.devices = make(map[dbus.ObjectPath][]*device)

	// initialize dbus object manager
	var err error
	bluezObjectManager, err = idbus.NewObjectManager(dbusBluezDest, "/")
	if err != nil {
		panic(err)
	}
	objects, err := bluezObjectManager.GetManagedObjects()
	if err != nil {
		panic(err)
	}

	// add exists adapters and devices
	for path, data := range objects {
		b.handleInterfacesAdded(path, data)
	}

	// connect signals
	bluezObjectManager.ConnectInterfacesAdded(b.handleInterfacesAdded)
	bluezObjectManager.ConnectInterfacesRemoved(b.handleInterfacesRemoved)
}
func (b *Bluetooth) handleInterfacesAdded(path dbus.ObjectPath, data map[string]map[string]dbus.Variant) {
	if _, ok := data[dbusBluezIfsAdapter]; ok {
		b.addAdapter(path)
		if len(b.PrimaryAdapter) == 0 {
			b.updatePropPrimaryAdapter(path)
		}
	}
	if _, ok := data[dbusBluezIfsDevice]; ok {
		b.addDevice(path, data[dbusBluezIfsDevice])
	}
}
func (b *Bluetooth) handleInterfacesRemoved(path dbus.ObjectPath, interfaces []string) {
	if isStringInArray(dbusBluezIfsAdapter, interfaces) {
		b.removeAdapter(path)
		if dbus.ObjectPath(b.PrimaryAdapter) == path {
			if len(b.adapters) > 0 {
				b.updatePropPrimaryAdapter(b.adapters[0].Path)
			} else {
				b.updatePropPrimaryAdapter("")
			}
		}
	}
	if isStringInArray(dbusBluezIfsDevice, interfaces) {
		b.removeDevice(path)
	}
}
