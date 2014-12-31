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
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
	"time"
)

const deviceRssiNotInRange = -1000 // -1000db means device not in range

const (
	deviceStateDisconnected = 0
	deviceStateConnecting   = 1
	deviceStateConnected    = 2
)

type device struct {
	bluezDevice *bluez.Device1

	Path        dbus.ObjectPath
	AdapterPath dbus.ObjectPath

	Alias   string
	Trusted bool
	Paired  bool

	oldConnected bool
	connected    bool
	connecting   bool
	State        uint32

	// optinal
	Icon string
	RSSI int16
}

func newDevice(dpath dbus.ObjectPath, data map[string]dbus.Variant) (d *device) {
	d = &device{Path: dpath}
	d.bluezDevice, _ = bluezNewDevice(dpath)
	d.AdapterPath = d.bluezDevice.Adapter.Get()
	d.Alias = d.bluezDevice.Alias.Get()
	d.Trusted = d.bluezDevice.Trusted.Get()
	d.Paired = d.bluezDevice.Paired.Get()
	d.connected = d.bluezDevice.Connected.Get()
	d.oldConnected = d.connected
	d.notifyStateChanged()

	// optional properties, read from dbus object data in order to
	// check if property is exists
	d.Icon = getDBusObjectValueString(data, "Icon")
	if isDBusObjectKeyExists(data, "RSSI") {
		d.RSSI = getDBusObjectValueInt16(data, "RSSI")
	}
	d.fixRssi()

	d.connectProperties()
	return
}
func destroyDevice(d *device) {
	bluezDestroyDevice(d.bluezDevice)
}

func (d *device) notifyDeviceAdded() {
	logger.Debug("DeviceAdded", marshalJSON(d))
	dbus.Emit(bluetooth, "DeviceAdded", marshalJSON(d))
	bluetooth.setPropState()
}
func (d *device) notifyDeviceRemoved() {
	logger.Debug("DeviceRemoved", marshalJSON(d))
	dbus.Emit(bluetooth, "DeviceRemoved", marshalJSON(d))
	bluetooth.setPropState()
}
func (d *device) notifyDevicePropertiesChanged() {
	logger.Debug("DevicePropertiesChanged", marshalJSON(d))
	dbus.Emit(bluetooth, "DevicePropertiesChanged", marshalJSON(d))
	bluetooth.setPropState()
}

func (d *device) connectProperties() {
	d.bluezDevice.Connected.ConnectChanged(func() {
		d.connected = d.bluezDevice.Connected.Get()
		if d.oldConnected != d.connected {
			d.oldConnected = d.connected
			if d.connected {
				notifyBluetoothConnected(d.Alias)
			} else {
				notifyBluetoothDisconnected(d.Alias)
			}
		}
		d.notifyStateChanged()
	})
	d.bluezDevice.Alias.ConnectChanged(func() {
		d.Alias = d.bluezDevice.Alias.Get()
		d.notifyDevicePropertiesChanged()
		bluetooth.setPropDevices()
	})
	d.bluezDevice.Trusted.ConnectChanged(func() {
		d.Trusted = d.bluezDevice.Trusted.Get()
		d.notifyDevicePropertiesChanged()
		bluetooth.setPropDevices()
	})
	d.bluezDevice.Paired.ConnectChanged(func() {
		d.Paired = d.bluezDevice.Paired.Get()
		d.notifyDevicePropertiesChanged()
		bluetooth.setPropDevices()
	})
	d.bluezDevice.Icon.ConnectChanged(func() {
		d.Icon = d.bluezDevice.Icon.Get()
		d.notifyDevicePropertiesChanged()
		bluetooth.setPropDevices()
	})
	d.bluezDevice.RSSI.ConnectChanged(func() {
		d.RSSI = d.bluezDevice.RSSI.Get()
		d.fixRssi()
		d.notifyDevicePropertiesChanged()
		bluetooth.setPropDevices()
	})
}
func (d *device) notifyStateChanged() {
	if d.connected {
		d.connecting = false
		d.State = deviceStateConnected
	} else if d.connecting {
		d.State = deviceStateConnecting
	} else {
		d.State = deviceStateDisconnected
	}
	logger.Debugf("notifyStateChanged: %#v", d) // TODO test
	d.notifyDevicePropertiesChanged()
	bluetooth.setPropDevices()
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
	d := newDevice(dpath, data)
	if b.isDeviceExists(b.devices[d.AdapterPath], dpath) {
		logger.Warning("repeat add device:", dpath)
		return
	}
	d.notifyDeviceAdded()
	b.devices[d.AdapterPath] = append(b.devices[d.AdapterPath], d)
	b.setPropDevices()
}
func (b *Bluetooth) removeDevice(dpath dbus.ObjectPath) {
	// find adapter of the device
	for apath, devices := range b.devices {
		if b.isDeviceExists(devices, dpath) {
			d, _ := b.getDevice(dpath)
			d.notifyDeviceRemoved()

			b.devices[apath] = b.doRemoveDevice(devices, dpath)
			b.setPropDevices()
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
	destroyDevice(devices[i])
	copy(devices[i:], devices[i+1:])
	devices[len(devices)-1] = nil
	devices = devices[:len(devices)-1]
	return devices
}
func (b *Bluetooth) getDevice(dpath dbus.ObjectPath) (d *device, err error) {
	for _, devices := range b.devices {
		if i := b.getDeviceIndex(devices, dpath); i >= 0 {
			d = devices[i]
			return
		}
	}
	err = fmt.Errorf("device not exists %s", dpath)
	logger.Error(err)
	return
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

// GetDevices return all device objects that marshaled by json.
func (b *Bluetooth) GetDevices(apath dbus.ObjectPath) (devicesJSON string, err error) {
	devices := b.devices[apath]
	devicesJSON = marshalJSON(devices)
	return
}

func (b *Bluetooth) ConnectDevice(dpath dbus.ObjectPath) (err error) {
	// mark device is connecting
	d, err := b.getDevice(dpath)
	if err != nil {
		return
	}
	d.connecting = true
	d.notifyStateChanged()

	go func() {
		if !bluezGetDeviceTrusted(dpath) {
			bluezSetDeviceTrusted(dpath, true)
		}
		if !bluezGetDevicePaired(dpath) {
			bluezPairDevice(dpath)
			time.Sleep(200 * time.Millisecond)
		}
		err = bluezConnectDevice(dpath)

		if d.connecting {
			d.connecting = false
			d.notifyStateChanged()
		}
	}()
	return
}

func (b *Bluetooth) DisconnectDevice(dpath dbus.ObjectPath) (err error) {
	// mark disconnected in time
	d, err := b.getDevice(dpath)
	if err != nil {
		return
	}
	d.connected = false
	d.notifyStateChanged()

	go bluezDisconnectDevice(dpath)
	return
}

func (b *Bluetooth) RemoveDevice(apath, dpath dbus.ObjectPath) (err error) {
	// TODO
	go bluezRemoveDevice(apath, dpath)
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
