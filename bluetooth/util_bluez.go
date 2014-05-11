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

func bluezStartDiscovery(apath dbus.ObjectPath) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	err = bluezAdapter.StartDiscovery()
	return
}

func bluezStopDiscovery(apath dbus.ObjectPath) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	err = bluezAdapter.StopDiscovery()
	return
}

func bluezGetAdapterAlias(apath dbus.ObjectPath) (alias string) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	alias = bluezAdapter.Alias.Get()
	return
}
func bluezSetAdapterAlias(apath dbus.ObjectPath, alias string) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	bluezAdapter.Alias.Set(alias)
	return
}

func bluezGetAdapterDiscoverable(apath dbus.ObjectPath) (discoverable bool) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	discoverable = bluezAdapter.Discoverable.Get()
	return
}
func bluezSetAdapterDiscoverable(apath dbus.ObjectPath, discoverable bool) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	bluezAdapter.Discoverable.Set(discoverable)
	return
}

func bluezGetAdapterDiscovering(apath dbus.ObjectPath) (discovering bool) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	discovering = bluezAdapter.Discovering.Get()
	return
}

func bluezGetAdapterDiscoverableTimeout(apath dbus.ObjectPath) (discoverableTimeout uint32) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	discoverableTimeout = bluezAdapter.DiscoverableTimeout.Get()
	return
}
func bluezSetAdapterDiscoverableTimeout(apath dbus.ObjectPath, discoverableTimeout uint32) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	bluezAdapter.DiscoverableTimeout.Set(discoverableTimeout)
	return
}

func bluezGetAdapterPowered(apath dbus.ObjectPath) (powered bool) {
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

func bluezConnectDevice(dpath dbus.ObjectPath) (err error) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	err = bluezDevice.Connect()
	return
}

func bluezRemoveDevice(apath, dpath dbus.ObjectPath) (err error) {
	bluezAdapter, err := bluezNewAdapter(dpath)
	if err != nil {
		return
	}
	err = bluezAdapter.RemoveDevice(dpath)
	return
}

func bluezSetDeviceAlias(dpath dbus.ObjectPath, alias string) (err error) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	bluezDevice.Alias.Set(alias)
	return
}

func bluezSetDeviceTrusted(dpath dbus.ObjectPath, trusted bool) (err error) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	bluezDevice.Trusted.Set(trusted)
	return
}
