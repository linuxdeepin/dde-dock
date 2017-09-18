/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package bluetooth

import (
	"dbus/org/bluez"
	sysdbus "dbus/org/freedesktop/dbus/system"
	"pkg.deepin.io/lib/dbus"
	"time"
)

var bluezDBusDaemon *sysdbus.DBusDaemon

func bluezNewDBusDaemon() (*sysdbus.DBusDaemon, error) {
	return sysdbus.NewDBusDaemon("org.freedesktop.DBus", "/org/freedesktop/DBus")
}
func bluezDestroyDbusDaemon(d *sysdbus.DBusDaemon) {
	sysdbus.DestroyDBusDaemon(d)
}
func bluezNewObjectManager() (objectManager *sysdbus.ObjectManager, err error) {
	objectManager, err = sysdbus.NewObjectManager(dbusBluezDest, "/")
	if err != nil {
		logger.Error(err)
	}
	return
}
func bluezDestroyObjectManager(objectManager *sysdbus.ObjectManager) {
	sysdbus.DestroyObjectManager(objectManager)
}

func bluezNewAdapter(apath dbus.ObjectPath) (bluezAdapter *bluez.Adapter1, err error) {
	bluezAdapter, err = bluez.NewAdapter1(dbusBluezDest, apath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func bluezDestroyAdapter(bluezAdapter *bluez.Adapter1) {
	bluez.DestroyAdapter1(bluezAdapter)
}

func bluezNewDevice(dpath dbus.ObjectPath) (bluezDevice *bluez.Device1, err error) {
	bluezDevice, err = bluez.NewDevice1(dbusBluezDest, dpath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func bluezDestroyDevice(bluezDevice *bluez.Device1) {
	bluez.DestroyDevice1(bluezDevice)
}

func bluezGetAdapters() (apathes []dbus.ObjectPath) {
	objectManager, err := bluezNewObjectManager()
	if err != nil {
		return
	}
	defer bluezDestroyObjectManager(objectManager)

	objects, err := objectManager.GetManagedObjects()
	if err != nil {
		logger.Error(err)
		return
	}
	for path, data := range objects {
		if _, ok := data[dbusBluezIfsAdapter]; ok {
			apathes = append(apathes, path)
		}
	}
	return
}

func bluezStartDiscovery(apath dbus.ObjectPath) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	err = bluezAdapter.StartDiscovery()
	if err != nil {
		logger.Error(err)
	}
	return
}

func bluezStopDiscovery(apath dbus.ObjectPath) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	err = bluezAdapter.StopDiscovery()
	if err != nil {
		logger.Error(err)
	}
	return
}

func bluezGetAdapterAddress(apath dbus.ObjectPath) (address string) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	address = bluezAdapter.Address.Get()
	return
}

func bluezGetAdapterAlias(apath dbus.ObjectPath) (alias string) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	alias = bluezAdapter.Alias.Get()
	return
}
func bluezSetAdapterAlias(apath dbus.ObjectPath, alias string) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	bluezAdapter.Alias.Set(alias)
	return
}

func bluezGetAdapterDiscoverable(apath dbus.ObjectPath) (discoverable bool) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	discoverable = bluezAdapter.Discoverable.Get()
	return
}
func bluezSetAdapterDiscoverable(apath dbus.ObjectPath, discoverable bool) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	bluezAdapter.Discoverable.Set(discoverable)
	return
}

func bluezSetAdapterDiscovering(apath dbus.ObjectPath, discovering bool) (err error) {
	if discovering {
		err = bluezStartDiscovery(apath)
		go func() {
			// adapter is not ready, retry again
			if err != nil {
				time.Sleep(3 * time.Second)
				bluezStartDiscovery(apath)
			}

		}()
	} else {
		err = bluezStopDiscovery(apath)
	}
	return
}

func bluezGetAdapterDiscovering(apath dbus.ObjectPath) (discovering bool) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	discovering = bluezAdapter.Discovering.Get()
	return
}

func bluezGetAdapterDiscoverableTimeout(apath dbus.ObjectPath) (discoverableTimeout uint32) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	discoverableTimeout = bluezAdapter.DiscoverableTimeout.Get()
	return
}
func bluezSetAdapterDiscoverableTimeout(apath dbus.ObjectPath, discoverableTimeout uint32) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	bluezAdapter.DiscoverableTimeout.Set(discoverableTimeout)
	return
}

func bluezGetAdapterPowered(apath dbus.ObjectPath) (powered bool) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	powered = bluezAdapter.Powered.Get()
	return
}
func bluezSetAdapterPowered(apath dbus.ObjectPath, powered bool) (er error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	bluezAdapter.Powered.Set(powered)
	return
}

func bluezPairDevice(dpath dbus.ObjectPath) (err error) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	defer bluezDestroyDevice(bluezDevice)

	err = bluezDevice.Pair()
	if err != nil {
		logger.Error(err)
	}
	return
}

func bluezConnectDevice(dpath dbus.ObjectPath) (err error) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	defer bluezDestroyDevice(bluezDevice)

	err = bluezDevice.Connect()
	if err != nil {
		logger.Error(err)
	}
	return
}

func bluezDisconnectDevice(dpath dbus.ObjectPath) (err error) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	defer bluezDestroyDevice(bluezDevice)

	err = bluezDevice.Disconnect()
	if err != nil {
		logger.Error(err)
	}
	return
}

func bluezRemoveDevice(apath, dpath dbus.ObjectPath) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	defer bluezDestroyAdapter(bluezAdapter)

	err = bluezAdapter.RemoveDevice(dpath)
	if err != nil {
		logger.Error(err)
	}
	return
}

func bluezSetDeviceAlias(dpath dbus.ObjectPath, alias string) (err error) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	defer bluezDestroyDevice(bluezDevice)

	bluezDevice.Alias.Set(alias)
	return
}

func bluezSetDeviceTrusted(dpath dbus.ObjectPath, trusted bool) (err error) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	defer bluezDestroyDevice(bluezDevice)

	bluezDevice.Trusted.Set(trusted)
	return
}

func bluezGetDeviceAddress(dpath dbus.ObjectPath) (address string) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	defer bluezDestroyDevice(bluezDevice)

	address = bluezDevice.Address.Get()
	return
}
func bluezGetDeviceTrusted(dpath dbus.ObjectPath) (trusted bool) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	defer bluezDestroyDevice(bluezDevice)

	trusted = bluezDevice.Trusted.Get()
	return
}

func bluezGetDevicePaired(dpath dbus.ObjectPath) (paired bool) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	defer bluezDestroyDevice(bluezDevice)

	paired = bluezDevice.Paired.Get()
	return
}

func bluezWatchRestart() {
	bluezDBusDaemon, _ = bluezNewDBusDaemon()
	logger.Info("bluezWatchRestart", bluezDBusDaemon)
	bluezDBusDaemon.ConnectNameOwnerChanged(func(name, oldOwner, newOwner string) {
		if name == dbusBluezDest {
			// if a new dbus session was installed, the name and newOwner
			// will be not empty, if a dbus session was uninstalled, the
			// name and oldOwner will be not empty
			if len(newOwner) != 0 {
				// network-manager is starting
				logger.Info("bluetooth is starting")
				time.Sleep(1 * time.Second)
				initBluetooth()
			} else {
				// network-manager stopped
				logger.Info("bluetooth stopped")
				destroyBluetooth()
			}
		}
	})
}
