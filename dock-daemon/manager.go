package main

import "dlib/dbus"
import busdaemon "dbus/org/freedesktop/dbus"

type Manager struct {
	Entries []*EntryProxyer

	Added   func(dbus.ObjectPath)
	Removed func(string)
}

func (*Manager) watchEntries() {
	busdaemon, err := busdaemon.NewDBusDaemon("org.freedesktop.DBus", "/")
	if err != nil {
		panic(err)
	}
	busdaemon.ListNames()
	//TODO: monitor name lost/ name acquire
}
