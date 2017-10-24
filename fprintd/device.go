/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/lib/dbus"
)

type Device struct {
	core *fprint.Device

	// TODO: enroll image
	EnrollStatus         func(string, bool)
	VerifyStatus         func(string, bool)
	VerifyFingerSelected func(string)
}
type Devices []*Device

func newDevice(objPath dbus.ObjectPath) *Device {
	var dev Device
	dev.core, _ = fprint.NewDevice(fprintDest, objPath)

	dev.core.ConnectEnrollStatus(func(status string, ok bool) {
		dbus.Emit(&dev, "EnrollStatus", status, ok)
	})
	dev.core.ConnectVerifyStatus(func(status string, ok bool) {
		dbus.Emit(&dev, "VerifyStatus", status, ok)
	})
	dev.core.ConnectVerifyFingerSelected(func(finger string) {
		dbus.Emit(&dev, "VerifyFingerSelected", finger)
	})

	return &dev
}

func destroyDevice(dev *Device) {
	fprint.DestroyDevice(dev.core)
	dev = nil
}

func (dev *Device) Claim(username string) error {
	return dev.core.Claim(username)
}

func (dev *Device) Release() error {
	return dev.core.Release()
}

func (dev *Device) EnrollStart(finger string) error {
	return dev.core.EnrollStart(finger)
}

func (dev *Device) EnrollStop() error {
	return dev.core.EnrollStop()
}

func (dev *Device) VerifyStart(finger string) error {
	return dev.core.VerifyStart(finger)
}

func (dev *Device) VerifyStop() error {
	return dev.core.VerifyStop()
}

func (dev *Device) DeleteEnrolledFingers(username string) error {
	return dev.core.DeleteEnrolledFingers(username)
}

func (dev *Device) ListEnrolledFingers(username string) ([]string, error) {
	return dev.core.ListEnrolledFingers(username)
}

func (dev *Device) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: string(convertFprintPath(dev.core.Path)),
		Interface:  dbusDeviceIFC,
	}
}

func destroyDevices(list Devices) {
	for _, dev := range list {
		destroyDevice(dev)
	}
	list = nil
}

func (devList Devices) Add(objPath dbus.ObjectPath) Devices {
	dev := devList.Get(objPath)
	if dev != nil {
		return devList
	}

	var v = newDevice(objPath)
	err := dbus.InstallOnSession(v)
	if err != nil {
		logger.Warning("Failed to install dbus:", objPath)
		return devList
	}

	devList = append(devList, newDevice(objPath))
	return devList
}

func (devList Devices) Get(objPath dbus.ObjectPath) *Device {
	for _, dev := range devList {
		if dev.core.Path == objPath {
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
		if dev.core.Path == objPath {
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

func convertFprintPath(objPath dbus.ObjectPath) dbus.ObjectPath {
	return dbus.ObjectPath(dbusPath + "/Device/" + path.Base(string(objPath)))
}
