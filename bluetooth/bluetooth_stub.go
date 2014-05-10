package main

import (
	"dlib/dbus"
)

func (b *Bluetooth) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	logger.Debug("OnPropertiesChanged()", name)
	switch name {
	// TODO
	}
}

// TODO
func (b *Bluetooth) updatePropPrimaryAdapter() {
	// b.PrimaryAdapter = "hh"
	dbus.NotifyChange(b, "PrimaryAdapter")
}
func (b *Bluetooth) doSetPrimaryAdapter(adapterName string) {
	// b.updatePropDevices()
	// b.updatePropAlias()
	// b.updatePropPowered()
	// b.updatePropDiscoverable()
	// b.updatePropDiscoverableTimeout()
}

func (b *Bluetooth) updatePropAdapters() {
	b.Adapters = marshalJSON(b.adapters)
	dbus.NotifyChange(b, "Adapters")
	logger.Debug(b.Adapters) // TODO test
	// TODO update alias properties for primary adapter
}

func (b *Bluetooth) updatePropDevices() {
	b.Devices = marshalJSON(b.devices)
	dbus.NotifyChange(b, "Devices")
	logger.Debug(b.Devices) // TODO test
}

// TODO
func (b *Bluetooth) updatePropAlias() {
	// dbus.NotifyChange(b, "Alias")
}

// TODO
func (b *Bluetooth) updatePropPowered() {
	// dbus.NotifyChange(b, "Powered")
}

// TODO
func (b *Bluetooth) updatePropDiscoverable() {
	// dbus.NotifyChange(b, "Discoverable")
}

// TODO
func (b *Bluetooth) updatePropDiscoverableTimeout() {
	// dbus.NotifyChange(b, "DiscoverableTimeout")
}

// TODO
