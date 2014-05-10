package main

import (
	"dbus/org/bluez"
	"dlib/dbus"
)

func bluezNewAdapter(apath dbus.ObjectPath) (bluezAdapter *bluez.Adapter1, err error) {
	bluezAdapter, err = bluez.NewAdapter1(dbusBluezDest, apath)
	if err != nil {
		logger.Error(err)
	}
	return
}

func bluezNewDevice(dpath dbus.ObjectPath) (bluezDevice *bluez.Device1, err error) {
	bluezDevice, err = bluez.NewDevice1(dbusBluezDest, dpath)
	if err != nil {
		logger.Error(err)
	}
	return
}

func bluezGetAdapterPowered(apath dbus.ObjectPath) (powered bool, err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	powered = bluezAdapter.Powered.Get()
	return
}

func bluezSetAdapterPowered(apath dbus.ObjectPath, powered bool) (er error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	bluezAdapter.Powered.Set(powered)
	return
}
