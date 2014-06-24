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
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/dbus/property"
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
			b.syncConfigPowered()
		}
	case "Powered":
		if b.isPrimaryAdapterExists() {
			b.updateAdapterScanState(dbus.ObjectPath(b.PrimaryAdapter))
		}
		// else {
		// 	b.Powered.SetValue(false)
		// 	b.updatePropPowered()
		// }
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

	// update alias properties
	if b.isAdapterExists(apath) {
		bluezAdapter, _ := bluezNewAdapter(apath)
		b.Alias = property.NewWrapProperty(b, "Alias", bluezAdapter.Alias)
		b.Powered = property.NewWrapProperty(b, "Powered", bluezAdapter.Powered)
		b.Discoverable = property.NewWrapProperty(b, "Discoverable", bluezAdapter.Discoverable)
		b.DiscoverableTimeout = property.NewWrapProperty(b, "DiscoverableTimeout", bluezAdapter.DiscoverableTimeout)
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

func (b *Bluetooth) updatePropAlias() {
	dbus.NotifyChange(b, "Alias")
}

func (b *Bluetooth) updatePropPowered() {
	dbus.NotifyChange(b, "Powered")
}

func (b *Bluetooth) updatePropDiscoverable() {
	dbus.NotifyChange(b, "Discoverable")
}

func (b *Bluetooth) updatePropDiscoverableTimeout() {
	dbus.NotifyChange(b, "DiscoverableTimeout")
}
