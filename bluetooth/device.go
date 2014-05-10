package main

import (
	"dbus/org/bluez"
	"dlib/dbus"
)

const deviceRssiNotInRange = -1000 // -1000db means device not in range

type device struct {
	bluezDevice *bluez.Device1

	Path      dbus.ObjectPath
	Adapter   dbus.ObjectPath
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
	d.Adapter = d.bluezDevice.Adapter.Get()
	d.Alias = d.bluezDevice.Alias.Get()
	d.Trusted = d.bluezDevice.Trusted.Get()
	d.Paired = d.bluezDevice.Paired.Get()
	d.Connected = d.bluezDevice.Connected.Get()

	// optional properties, read from dbus object data in order to check if property is exists
	d.Icon = getDBusObjectValueString(data, "Icon")
	if isDBusObjectKeyExists(data, "RSSI") {
		d.RSSI = getDBusObjectValueInt16(data, "RSSI")
	} else {
		d.RSSI = deviceRssiNotInRange
	}

	// connect properties
	d.bluezDevice.Alias.ConnectChanged(func() {
		d.Alias = d.bluezDevice.Alias.Get()
		b.updatePropDevices()
	})
	d.bluezDevice.Trusted.ConnectChanged(func() {
		d.Trusted = d.bluezDevice.Trusted.Get()
		b.updatePropDevices()
	})
	d.bluezDevice.Paired.ConnectChanged(func() {
		d.Paired = d.bluezDevice.Paired.Get()
		b.updatePropDevices()
	})
	d.bluezDevice.Connected.ConnectChanged(func() {
		d.Connected = d.bluezDevice.Connected.Get()
		b.updatePropDevices()
	})
	d.bluezDevice.Icon.ConnectChanged(func() {
		d.Icon = d.bluezDevice.Icon.Get()
		b.updatePropDevices()
	})
	d.bluezDevice.RSSI.ConnectChanged(func() {
		d.RSSI = d.bluezDevice.RSSI.Get()
		b.updatePropDevices()
	})

	return
}

func (b *Bluetooth) addDevice(dpath dbus.ObjectPath, data map[string]dbus.Variant) {
	if b.isDeviceExists(dpath) {
		logger.Warning("repeat add device:", dpath)
		return
	}
	d := b.newDevice(dpath, data)
	b.devices = append(b.devices, d)
	b.updatePropDevices()
}

func (b *Bluetooth) removeDevice(dpath dbus.ObjectPath) {
	i := b.getDeviceIndex(dpath)
	if i < 0 {
		logger.Warning("repeat remove device:", dpath)
		return
	}
	copy(b.devices[i:], b.devices[i+1:])
	b.devices[len(b.devices)-1] = nil
	b.devices = b.devices[:len(b.devices)-1]
	b.updatePropDevices()
}

func (b *Bluetooth) isDeviceExists(dpath dbus.ObjectPath) bool {
	if b.getDeviceIndex(dpath) >= 0 {
		return true
	}
	return false
}

func (b *Bluetooth) getDeviceIndex(dpath dbus.ObjectPath) int {
	for i, d := range b.devices {
		if d.Path == dpath {
			return i
		}
	}
	return -1
}

// TODO
func (b *Bluetooth) ConnectDeivce(dpath dbus.ObjectPath) (err error) {
	return
}

func (b *Bluetooth) SetDeivceAlias(dpath dbus.ObjectPath, alias string) (err error) {
	err = bluezSetDeviceAlias(dpath, alias)
	return
}

func (b *Bluetooth) SetDeivceTrusted(dpath dbus.ObjectPath, trusted bool) (err error) {
	err = bluezSetDeviceTrusted(dpath, trusted)
	return
}
