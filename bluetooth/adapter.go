package main

import (
	"dbus/org/bluez"
	"dlib/dbus"
)

type adapter struct {
	bluezAdapter *bluez.Adapter1

	Path dbus.ObjectPath
}

func newAdapter(apath dbus.ObjectPath) (a *adapter) {
	a = &adapter{Path: apath}
	a.bluezAdapter, _ = bluezNewAdapter(apath)
	// TODO
	// a.bluezAdapter.
	return
}
