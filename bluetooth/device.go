package main

import (
	"dbus/org/bluez"
	"dlib/dbus"
)

type device struct {
	bluezDevice *bluez.Device1

	Path    dbus.ObjectPath
	Adapter dbus.ObjectPath

	Icon      string
	Paired    bool
	Connected bool
	RSSI      int16

	Alias   string `access:"readwrite"`
	Trusted bool   `access:"readwrite"`
}

func newDevice(dpath dbus.ObjectPath) (d *device) {
	d = &device{Path: dpath}
	d.bluezDevice, _ = bluezNewDevice(dpath)
	// TODO
	return
}
