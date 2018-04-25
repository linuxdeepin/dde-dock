/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"time"

	"github.com/linuxdeepin/go-dbus-factory/org.bluez"
	"pkg.deepin.io/lib/dbus1"
)

func bluezNewObjectManager() (*bluez.ObjectManager, error) {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	return bluez.NewObjectManager(systemConn), nil
}

func bluezNewAdapter(apath dbus.ObjectPath) (bluezAdapter *bluez.HCI, err error) {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	return bluez.NewHCI(systemConn, apath)
}

func bluezNewDevice(dpath dbus.ObjectPath) (bluezDevice *bluez.Device, err error) {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	return bluez.NewDevice(systemConn, dpath)
}

func bluezGetAdapters() (apathes []dbus.ObjectPath) {
	objectManager, err := bluezNewObjectManager()
	if err != nil {
		return
	}
	objects, err := objectManager.GetManagedObjects(0)
	if err != nil {
		logger.Error(err)
		return
	}
	for path, data := range objects {
		if _, ok := data[bluezAdapterDBusInterface]; ok {
			apathes = append(apathes, dbus.ObjectPath(path))
		}
	}
	return
}

func bluezStartDiscovery(apath dbus.ObjectPath) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	err = bluezAdapter.StartDiscovery(0)
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
	err = bluezAdapter.StopDiscovery(0)
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
	address, err = bluezAdapter.Address().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	return
}

func bluezGetAdapterAlias(apath dbus.ObjectPath) (alias string) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	alias, err = bluezAdapter.Alias().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	return
}

func bluezSetAdapterAlias(apath dbus.ObjectPath, alias string) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	return bluezAdapter.Alias().Set(0, alias)
}

func bluezGetAdapterDiscoverable(apath dbus.ObjectPath) (discoverable bool) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	discoverable, err = bluezAdapter.Discoverable().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	return
}

func bluezSetAdapterDiscoverable(apath dbus.ObjectPath, discoverable bool) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	return bluezAdapter.Discoverable().Set(0, discoverable)
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
	discovering, err = bluezAdapter.Discovering().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	return
}

func bluezGetAdapterDiscoverableTimeout(apath dbus.ObjectPath) (discoverableTimeout uint32) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	discoverableTimeout, err = bluezAdapter.DiscoverableTimeout().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	return
}

func bluezSetAdapterDiscoverableTimeout(apath dbus.ObjectPath, discoverableTimeout uint32) (err error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	return bluezAdapter.DiscoverableTimeout().Set(0, discoverableTimeout)
}

func bluezGetAdapterPowered(apath dbus.ObjectPath) (powered bool) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	powered, err = bluezAdapter.Powered().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	return
}

func bluezSetAdapterPowered(apath dbus.ObjectPath, powered bool) (er error) {
	bluezAdapter, err := bluezNewAdapter(apath)
	if err != nil {
		return
	}
	return bluezAdapter.Powered().Set(0, powered)
}

func bluezPairDevice(dpath dbus.ObjectPath) (err error) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	err = bluezDevice.Pair(0)
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
	err = bluezDevice.Connect(0)
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
	err = bluezDevice.Disconnect(0)
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
	err = bluezAdapter.RemoveDevice(0, dpath)
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
	return bluezDevice.Alias().Set(0, alias)
}

func bluezSetDeviceTrusted(dpath dbus.ObjectPath, trusted bool) (err error) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	return bluezDevice.Trusted().Set(0, trusted)
}

func bluezGetDeviceAddress(dpath dbus.ObjectPath) (address string) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	address, err = bluezDevice.Address().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	return
}

func bluezGetDeviceTrusted(dpath dbus.ObjectPath) (trusted bool) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	trusted, err = bluezDevice.Trusted().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	return
}

func bluezGetDevicePaired(dpath dbus.ObjectPath) (paired bool) {
	bluezDevice, err := bluezNewDevice(dpath)
	if err != nil {
		return
	}
	paired, err = bluezDevice.Paired().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	return
}
