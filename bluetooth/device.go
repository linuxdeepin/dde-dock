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
	"dbus/org/bluez"
	"pkg.linuxdeepin.com/lib/dbus"
	"time"
)

const deviceRssiNotInRange = -1000 // -1000db means device not in range

type device struct {
	bluezDevice *bluez.Device1
	adapter     dbus.ObjectPath

	Path dbus.ObjectPath

	Alias     string
	Trusted   bool
	Paired    bool
	Connected bool

	// optinal
	Icon string
	RSSI int16
}

func (b *Bluetooth) newDevice(dpath dbus.ObjectPath, data map[string]dbus.Variant) (d *device) {
	d = &device{Path: dpath}
	d.bluezDevice, _ = bluezNewDevice(dpath)
	d.adapter = d.bluezDevice.Adapter.Get()
	d.Alias = d.bluezDevice.Alias.Get()
	d.Trusted = d.bluezDevice.Trusted.Get()
	d.Paired = d.bluezDevice.Paired.Get()
	d.Connected = d.bluezDevice.Connected.Get()

	// optional properties, read from dbus object data in order to check if property is exists
	d.Icon = getDBusObjectValueString(data, "Icon")
	if isDBusObjectKeyExists(data, "RSSI") {
		d.RSSI = getDBusObjectValueInt16(data, "RSSI")
	}
	d.fixRssi()

	d.connectProeprties()
	return
}
func (d *device) connectProeprties() {
	d.bluezDevice.Connected.ConnectChanged(func() {
		d.Connected = d.bluezDevice.Connected.Get()
		bluetooth.setPropDevices()
	})
	d.bluezDevice.Alias.ConnectChanged(func() {
		d.Alias = d.bluezDevice.Alias.Get()
		bluetooth.setPropDevices()
	})
	d.bluezDevice.Trusted.ConnectChanged(func() {
		d.Trusted = d.bluezDevice.Trusted.Get()
		bluetooth.setPropDevices()
	})
	d.bluezDevice.Paired.ConnectChanged(func() {
		d.Paired = d.bluezDevice.Paired.Get()
		bluetooth.setPropDevices()
	})
	d.bluezDevice.Icon.ConnectChanged(func() {
		d.Icon = d.bluezDevice.Icon.Get()
		bluetooth.setPropDevices()
	})
	d.bluezDevice.RSSI.ConnectChanged(func() {
		d.RSSI = d.bluezDevice.RSSI.Get()
		d.fixRssi()
		bluetooth.setPropDevices()
	})
}
func (d *device) resetConnect() {
	d.bluezDevice.Connected.Reset()
	d.bluezDevice.Alias.Reset()
	d.bluezDevice.Trusted.Reset()
	d.bluezDevice.Paired.Reset()
	d.bluezDevice.Icon.Reset()
	d.bluezDevice.RSSI.Reset()
}
func (d *device) fixRssi() {
	if d.RSSI == 0 {
		if d.Trusted {
			d.RSSI = deviceRssiNotInRange / 2
		} else {
			d.RSSI = deviceRssiNotInRange
		}
	}
}

func (b *Bluetooth) addDevice(dpath dbus.ObjectPath, data map[string]dbus.Variant) {
	d := b.newDevice(dpath, data)
	if b.isDeviceExists(b.devices[d.adapter], dpath) {
		logger.Warning("repeat add device:", dpath)
		return
	}
	b.devices[d.adapter] = append(b.devices[d.adapter], d)
	b.setPropDevices()

	// send signal DeviceAdded() if device is managed by primary adapter
	if dbus.ObjectPath(b.PrimaryAdapter) == d.adapter {
		dbus.Emit(b, "DeviceAdded", marshalJSON(d))
	}
}
func (b *Bluetooth) removeDevice(dpath dbus.ObjectPath) {
	// find adapter of the device
	for apath, devices := range b.devices {
		if b.isDeviceExists(devices, dpath) {
			b.devices[apath] = b.doRemoveDevice(devices, dpath)
			b.setPropDevices()

			// send signal DeviceRemoved() if device is managed by primary adapter
			if dbus.ObjectPath(b.PrimaryAdapter) == apath {
				d := &device{Path: dpath}
				dbus.Emit(b, "DeviceRemoved", marshalJSON(d))
			}
			return
		}
	}
}
func (b *Bluetooth) doRemoveDevice(devices []*device, dpath dbus.ObjectPath) []*device {
	i := b.getDeviceIndex(devices, dpath)
	if i < 0 {
		logger.Warning("repeat remove device:", dpath)
		return devices
	}
	devices[i].resetConnect()
	copy(devices[i:], devices[i+1:])
	devices[len(devices)-1] = nil
	devices = devices[:len(devices)-1]
	return devices
}
func (b *Bluetooth) isDeviceExists(devices []*device, dpath dbus.ObjectPath) bool {
	if b.getDeviceIndex(devices, dpath) >= 0 {
		return true
	}
	return false
}
func (b *Bluetooth) getDeviceIndex(devices []*device, dpath dbus.ObjectPath) int {
	for i, d := range devices {
		if d.Path == dpath {
			return i
		}
	}
	return -1
}

// GetDevices return all devices object that marshaled by json.
func (b *Bluetooth) GetDevices() (devicesJSON string) {
	devices := b.devices[dbus.ObjectPath(b.PrimaryAdapter)]
	devicesJSON = marshalJSON(devices)
	return
}

func (b *Bluetooth) ConnectDevice(dpath dbus.ObjectPath) (err error) {
	go func() {
		bluezSetDeviceTrusted(dpath, true)
		if !bluezGetDevicePaired(dpath) {
			bluezPairDevice(dpath)
			time.Sleep(200 * time.Millisecond)
		}
		bluezConnectDevice(dpath)
	}()
	return
}

func (b *Bluetooth) DisconnectDevice(dpath dbus.ObjectPath) (err error) {
	go bluezDisconnectDevice(dpath)
	return
}

func (b *Bluetooth) RemoveDevice(dpath dbus.ObjectPath) (err error) {
	// TODO
	go bluezRemoveDevice(dbus.ObjectPath(b.PrimaryAdapter), dpath)
	return
}

func (b *Bluetooth) SetDeviceAlias(dpath dbus.ObjectPath, alias string) (err error) {
	go bluezSetDeviceAlias(dpath, alias)
	return
}

func (b *Bluetooth) SetDeviceTrusted(dpath dbus.ObjectPath, trusted bool) (err error) {
	go bluezSetDeviceTrusted(dpath, trusted)
	return
}
