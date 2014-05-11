package main

import (
	"dbus/org/bluez"
	"dlib/dbus"
)

type adapter struct {
	bluezAdapter *bluez.Adapter1

	Path                dbus.ObjectPath
	Alias               string
	Powered             bool
	Discoverable        bool
	DiscoverableTimeout uint32
}

func (b *Bluetooth) newAdapter(apath dbus.ObjectPath) (a *adapter) {
	a = &adapter{Path: apath}
	a.bluezAdapter, _ = bluezNewAdapter(apath)
	a.Alias = a.bluezAdapter.Alias.Get()
	a.Powered = a.bluezAdapter.Powered.Get()
	a.Discoverable = a.bluezAdapter.Discoverable.Get()
	a.DiscoverableTimeout = a.bluezAdapter.DiscoverableTimeout.Get()

	// TODO connect properties
	a.bluezAdapter.Alias.ConnectChanged(func() {
		a.Alias = a.bluezAdapter.Alias.Get()
		b.updatePropAdapters()
	})
	a.bluezAdapter.Powered.ConnectChanged(func() {
		a.Powered = a.bluezAdapter.Powered.Get()
		b.updatePropAdapters()
	})
	a.bluezAdapter.Discoverable.ConnectChanged(func() {
		a.Discoverable = a.bluezAdapter.Discoverable.Get()
		b.updatePropAdapters()
	})
	a.bluezAdapter.DiscoverableTimeout.ConnectChanged(func() {
		a.DiscoverableTimeout = a.bluezAdapter.DiscoverableTimeout.Get()
		b.updatePropAdapters()
	})

	return
}

func (b *Bluetooth) addAdapter(apath dbus.ObjectPath) {
	if b.isAdapterExists(apath) {
		logger.Warning("repeat add adapter:", apath)
		return
	}
	a := b.newAdapter(apath)
	b.adapters = append(b.adapters, a)
	b.updatePropAdapters()
}

func (b *Bluetooth) removeAdapter(apath dbus.ObjectPath) {
	i := b.getAdapterIndex(apath)
	if i < 0 {
		logger.Warning("repeat remove adapter:", apath)
		return
	}
	copy(b.adapters[i:], b.adapters[i+1:])
	b.adapters[len(b.adapters)-1] = nil
	b.adapters = b.adapters[:len(b.adapters)-1]
	b.updatePropAdapters()
}

func (b *Bluetooth) isAdapterExists(apath dbus.ObjectPath) bool {
	if b.getAdapterIndex(apath) >= 0 {
		return true
	}
	return false
}

func (b *Bluetooth) getAdapterIndex(apath dbus.ObjectPath) int {
	for i, a := range b.adapters {
		if a.Path == apath {
			return i
		}
	}
	return -1
}
