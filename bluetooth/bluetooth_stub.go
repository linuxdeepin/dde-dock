package bluetooth

import (
	"dlib/dbus"
	"dlib/dbus/property"
)

func (b *Bluetooth) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	logger.Debug("OnPropertiesChanged()", name)
	switch name {
	case "PrimaryAdapter":
		oldPrimaryAdapter, _ := oldv.(string)
		if b.PrimaryAdapter != oldPrimaryAdapter {
			b.updatePropPrimaryAdapter(dbus.ObjectPath(b.PrimaryAdapter))
		}
	case "Powered":
		b.updateAdapterScanState(dbus.ObjectPath(b.PrimaryAdapter))
	}
}

func (b *Bluetooth) updatePropPrimaryAdapter(apath dbus.ObjectPath) {
	b.PrimaryAdapter = string(apath)

	// power on primary adapter and power off other adapters
	for _, a := range b.adapters {
		if a.Path == apath {
			bluezSetAdapterPowered(a.Path, true)
		} else {
			bluezSetAdapterPowered(a.Path, false)
		}
		b.updateAdapterScanState(a.Path)
	}

	// TODO update alias properties
	if b.isAdapterExists(apath) {
		bluezAdapter, _ := bluezNewAdapter(apath)
		b.Alias = property.NewWrapProperty(b, "Alias", bluezAdapter.Alias)
		b.Powered = property.NewWrapProperty(b, "Powered", bluezAdapter.Powered)
		b.Discoverable = property.NewWrapProperty(b, "Discoverable", bluezAdapter.Discoverable)
		b.DiscoverableTimeout = property.NewWrapProperty(b, "DiscoverableTimeout", bluezAdapter.DiscoverableTimeout)
		// TODO remove
		// b.updatePropAlias(bluezGetAdapterAlias(apath))
		// b.updatePropPowered(bluezGetAdapterPowered(apath))
		// b.updatePropDiscoverable(bluezGetAdapterDiscoverable(apath))
		// b.updatePropDiscoverableTimeout(bluezGetAdapterDiscoverableTimeout(apath))
	}

	dbus.NotifyChange(b, "PrimaryAdapter")
}

func (b *Bluetooth) updateAdapterScanState(apath dbus.ObjectPath) {
	// if adapter is power on, just start discovery
	powered := bluezGetAdapterPowered(apath)
	if powered {
		if !bluezGetAdapterDiscovering(apath) {
			bluezStartDiscovery(apath)
		}
	}
}

func (b *Bluetooth) updatePropAdapters() {
	b.Adapters = marshalJSON(b.adapters)
	dbus.NotifyChange(b, "Adapters")
	logger.Debug(b.Adapters) // TODO test

	// TODO update alias properties for primary adapter
}

func (b *Bluetooth) updatePropDevices() {
	devices := b.devices[dbus.ObjectPath(b.PrimaryAdapter)]
	b.Devices = marshalJSON(devices)
	dbus.NotifyChange(b, "Devices")
	logger.Debug(b.Devices) // TODO test
}

// TODO
func (b *Bluetooth) updatePropAlias(alias string) {
	// b.Alias = alias
	dbus.NotifyChange(b, "Alias")
}

// TODO
func (b *Bluetooth) updatePropPowered(powered bool) {
	// b.Powered = powered
	dbus.NotifyChange(b, "Powered")
}

// TODO
func (b *Bluetooth) updatePropDiscoverable(discoverable bool) {
	// b.Discoverable = discoverable
	dbus.NotifyChange(b, "Discoverable")
}

// TODO
func (b *Bluetooth) updatePropDiscoverableTimeout(discoverableTimeout uint32) {
	// b.DiscoverableTimeout = discoverableTimeout
	dbus.NotifyChange(b, "DiscoverableTimeout")
}
