/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package fprintd

import (
	"dbus/net/reactivated/fprint"
	"path"

	oldDBusLib "pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

type Device struct {
	service *dbusutil.Service
	core    *fprint.Device

	methods *struct {
		Claim                 func() `in:"username"`
		EnrollStart           func() `in:"finger"`
		VerifyStart           func() `in:"finger"`
		DeleteEnrolledFingers func() `in:"username"`
		ListEnrolledFingers   func() `in:"username" out:"fingers"`
	}

	// TODO: enroll image
	signals *struct {
		EnrollStatus struct {
			status string
			ok     bool
		}
		VerifyStatus struct {
			status string
			ok     bool
		}
		VerifyFingerSelected struct {
			finger string
		}
	}
}
type Devices []*Device

func newDevice(objPath dbus.ObjectPath, service *dbusutil.Service) *Device {
	var dev Device
	dev.service = service
	dev.core, _ = fprint.NewDevice(fprintDBusServiceName, oldDBusLib.ObjectPath(objPath))

	dev.core.ConnectEnrollStatus(func(status string, ok bool) {
		dev.service.Emit(&dev, "EnrollStatus", status, ok)
	})
	dev.core.ConnectVerifyStatus(func(status string, ok bool) {
		dev.service.Emit(&dev, "VerifyStatus", status, ok)
	})
	dev.core.ConnectVerifyFingerSelected(func(finger string) {
		dev.service.Emit(&dev, "VerifyFingerSelected", finger)
	})

	return &dev
}

func destroyDevice(dev *Device) {
	fprint.DestroyDevice(dev.core)
	dev = nil
}

func (dev *Device) Claim(username string) *dbus.Error {
	err := dev.core.Claim(username)
	return dbusutil.ToError(err)
}

func (dev *Device) Release() *dbus.Error {
	err := dev.core.Release()
	return dbusutil.ToError(err)
}

func (dev *Device) EnrollStart(finger string) *dbus.Error {
	err := dev.core.EnrollStart(finger)
	return dbusutil.ToError(err)
}

func (dev *Device) EnrollStop() *dbus.Error {
	err := dev.core.EnrollStop()
	return dbusutil.ToError(err)
}

func (dev *Device) VerifyStart(finger string) *dbus.Error {
	err := dev.core.VerifyStart(finger)
	return dbusutil.ToError(err)
}

func (dev *Device) VerifyStop() *dbus.Error {
	err := dev.core.VerifyStop()
	return dbusutil.ToError(err)
}

func (dev *Device) DeleteEnrolledFingers(username string) *dbus.Error {
	err := dev.core.DeleteEnrolledFingers(username)
	return dbusutil.ToError(err)
}

func (dev *Device) ListEnrolledFingers(username string) ([]string, *dbus.Error) {
	fingers, err := dev.core.ListEnrolledFingers(username)
	if err != nil {
		return nil, dbusutil.ToError(err)
	}
	return fingers, nil
}

func (*Device) GetInterfaceName() string {
	return dbusDeviceInterface
}

func (dev *Device) getPath() dbus.ObjectPath {
	return convertFPrintPath(dbus.ObjectPath(dev.core.Path))
}

func destroyDevices(list Devices) {
	for _, dev := range list {
		destroyDevice(dev)
	}
	list = nil
}

func (devList Devices) Add(objPath dbus.ObjectPath, service *dbusutil.Service) Devices {
	dev := devList.Get(objPath)
	if dev != nil {
		return devList
	}

	var v = newDevice(objPath, service)
	err := service.Export(v.getPath(), v)
	if err != nil {
		logger.Warning("Failed to export:", objPath)
		return devList
	}

	devList = append(devList, v)
	return devList
}

func (devList Devices) Get(objPath dbus.ObjectPath) *Device {
	for _, dev := range devList {
		if dbus.ObjectPath(dev.core.Path) == objPath {
			return dev
		}
	}
	return nil
}

func (devList Devices) Delete(objPath dbus.ObjectPath) Devices {
	var (
		list Devices
		v    *Device
	)
	for _, dev := range devList {
		if dbus.ObjectPath(dev.core.Path) == objPath {
			v = dev
			continue
		}
		list = append(list, dev)
	}
	if v != nil {
		destroyDevice(v)
	}
	return list
}

func convertFPrintPath(objPath dbus.ObjectPath) dbus.ObjectPath {
	return dbus.ObjectPath(dbusPath + "/Device/" + path.Base(string(objPath)))
}
