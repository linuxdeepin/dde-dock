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
	sysdbus "dbus/org/freedesktop/dbus/system"
	"pkg.deepin.io/lib/dbus"
	"sync"
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

const (
	StateUnavailable = 0
	StateAvailable   = 1
	StateConnected   = 2
)

type dbusObjectData map[string]dbus.Variant
type dbusInterfaceData map[string]map[string]dbus.Variant
type dbusInterfacesData map[dbus.ObjectPath]map[string]map[string]dbus.Variant

type Bluetooth struct {
	config        *config
	objectManager *sysdbus.ObjectManager
	agent         *agent

	// adapter
	adaptersLock sync.Mutex
	adapters     map[dbus.ObjectPath]*adapter

	// device
	devicesLock sync.Mutex
	devices     map[dbus.ObjectPath][]*device

	State uint32 // StateUnavailable/StateAvailable/StateConnected

	// Bluetooth adaper/device properties changed signals
	AdapterAdded             func(adapterJSON string)
	AdapterRemoved           func(adapterJSON string)
	AdapterPropertiesChanged func(adapterJSON string)
	DeviceAdded              func(devJSON string)
	DeviceRemoved            func(devJSON string)
	DevicePropertiesChanged  func(devJSON string)

	// Bluetooth pair request signals
	DisplayPinCode func(device dbus.ObjectPath, pincode string)
	DisplayPasskey func(device dbus.ObjectPath, passkey uint32, entered uint32)
	// RequestConfirmation you shoud call Confirm with accpet
	RequestConfirmation func(device dbus.ObjectPath, passkey string)
	// RequestAuthorization you shoud call Confirm with accpet
	RequestAuthorization func(device dbus.ObjectPath)
	// RequestPinCode you should call FeedPinCode with accpet and key
	RequestPinCode func(device dbus.ObjectPath)
	// RequestPasskey you should call FeedPasskey with accpet and key
	RequestPasskey func(device dbus.ObjectPath)
}

func newBluetooth() (b *Bluetooth) {
	b = &Bluetooth{}
	b.adapters = map[dbus.ObjectPath]*adapter{}
	return
}

func (b *Bluetooth) destroy() {
	bluezDestroyObjectManager(b.objectManager)
	dbus.UnInstallObject(b)
}

func (b *Bluetooth) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusBluetoothDest,
		ObjectPath: dbusBluetoothPath,
		Interface:  dbusBluetoothIfs,
	}
}

func (b *Bluetooth) init() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			b.destroy()
		}
	}()

	go func() {
		b.config = newConfig()
		b.config.save()
		b.devices = make(map[dbus.ObjectPath][]*device)

		// initialize dbus object manager
		var err error
		b.objectManager, err = bluezNewObjectManager()
		if err != nil {
			return
		}

		// connect signals
		b.objectManager.ConnectInterfacesAdded(b.handleInterfacesAdded)
		b.objectManager.ConnectInterfacesRemoved(b.handleInterfacesRemoved)

		// add exists adapters and devices
		objects, err := b.objectManager.GetManagedObjects()
		if err != nil {
			logger.Error(err)
			return
		}
		for path, data := range objects {
			b.handleInterfacesAdded(path, data)
		}
	}()
}
func (b *Bluetooth) handleInterfacesAdded(path dbus.ObjectPath, data map[string]map[string]dbus.Variant) {
	if _, ok := data[dbusBluezIfsAdapter]; ok {
		requestUnblockBluetoothDevice()
		b.addAdapter(path)
	}
	if _, ok := data[dbusBluezIfsDevice]; ok {
		b.addDevice(path, data[dbusBluezIfsDevice])
	}
}
func (b *Bluetooth) handleInterfacesRemoved(path dbus.ObjectPath, interfaces []string) {
	if isStringInArray(dbusBluezIfsAdapter, interfaces) {
		b.removeAdapter(path)
	}
	if isStringInArray(dbusBluezIfsDevice, interfaces) {
		b.removeDevice(path)
	}
}
