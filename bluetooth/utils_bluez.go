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
