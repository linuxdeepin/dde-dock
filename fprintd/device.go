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
	"path"

	"github.com/linuxdeepin/go-dbus-factory/net.reactivated.fprint"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
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

func newDevice(objPath dbus.ObjectPath, service *dbusutil.Service,
	systemSigLoop *dbusutil.SignalLoop) *Device {
	var dev Device
	dev.service = service
	dev.core, _ = fprint.NewDevice(systemSigLoop.Conn(), objPath)

	dev.core.InitSignalExt(systemSigLoop, true)
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

func (dev *Device) destroy() {
	dev.core.RemoveHandler(proxy.RemoveAllHandlers)
	dev.service.StopExport(dev)
}

func (dev *Device) Claim(username string) *dbus.Error {
	err := dev.core.Claim(0, username)
	return dbusutil.ToError(err)
}

func (dev *Device) Release() *dbus.Error {
	err := dev.core.Release(0)
	return dbusutil.ToError(err)
}

func (dev *Device) EnrollStart(finger string) *dbus.Error {
	err := dev.core.EnrollStart(0, finger)
	return dbusutil.ToError(err)
}

func (dev *Device) EnrollStop() *dbus.Error {
	err := dev.core.EnrollStop(0)
	return dbusutil.ToError(err)
}

func (dev *Device) VerifyStart(finger string) *dbus.Error {
	err := dev.core.VerifyStart(0, finger)
	return dbusutil.ToError(err)
}

func (dev *Device) VerifyStop() *dbus.Error {
	err := dev.core.VerifyStop(0)
	return dbusutil.ToError(err)
}

func (dev *Device) DeleteEnrolledFingers(username string) *dbus.Error {
	err := dev.core.DeleteEnrolledFingers(0, username)
	return dbusutil.ToError(err)
}

func (dev *Device) ListEnrolledFingers(username string) ([]string, *dbus.Error) {
	fingers, err := dev.core.ListEnrolledFingers(0, username)
	if err != nil {
		return nil, dbusutil.ToError(err)
	}
	return fingers, nil
}

func (*Device) GetInterfaceName() string {
	return dbusDeviceInterface
}

func (dev *Device) getPath() dbus.ObjectPath {
	return convertFPrintPath(dev.core.Path_())
}

func destroyDevices(list Devices) {
	for _, dev := range list {
		dev.destroy()
	}
}

func (devList Devices) Add(objPath dbus.ObjectPath, service *dbusutil.Service,
	systemSigLoop *dbusutil.SignalLoop) Devices {
	dev := devList.Get(objPath)
	if dev != nil {
		return devList
	}

	var v = newDevice(objPath, service, systemSigLoop)
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
		if dev.core.Path_() == objPath {
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
		if dev.core.Path_() == objPath {
			v = dev
			continue
		}
		list = append(list, dev)
	}
	if v != nil {
		v.destroy()
	}
	return list
}

func convertFPrintPath(objPath dbus.ObjectPath) dbus.ObjectPath {
	return dbus.ObjectPath(dbusPath + "/Device/" + path.Base(string(objPath)))
}
