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
	"errors"
	"path"

	"github.com/linuxdeepin/go-dbus-factory/net.reactivated.fprint"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

const (
	actionIdEnroll = "com.deepin.daemon.fprintd.enroll"
	actionIdDelete = "com.deepin.daemon.fprintd.delete-enrolled-fingers"
)

type deviceMethods struct {
	Claim                 func() `in:"username"`
	ClaimForce            func() `in:"username"`
	GetCapabilities       func() `out:"caps"`
	EnrollStart           func() `in:"finger"`
	VerifyStart           func() `in:"finger"`
	DeleteEnrolledFingers func() `in:"username"`
	DeleteEnrolledFinger  func() `in:"username,finger"`
	ListEnrolledFingers   func() `in:"username" out:"fingers"`
}

type deviceSignals struct {
	EnrollStatus struct {
		status string
		done   bool
	}
	VerifyStatus struct {
		status string
		done   bool
	}
	VerifyFingerSelected struct {
		finger string
	}
}

type IDevice interface {
	IsDevice()
	destroy()
	getCorePath() dbus.ObjectPath
	getPath() dbus.ObjectPath
	dbusutil.Implementer
}

type Device struct {
	service *dbusutil.Service
	core    *fprint.Device

	ScanType string
	methods  *deviceMethods

	// TODO: enroll image
	signals *deviceSignals
}

func (d *Device) IsDevice() {}

type Devices []IDevice

func newDevice(objPath dbus.ObjectPath, service *dbusutil.Service,
	systemSigLoop *dbusutil.SignalLoop) *Device {
	var dev Device
	dev.service = service
	dev.core, _ = fprint.NewDevice(systemSigLoop.Conn(), objPath)
	dev.ScanType, _ = dev.core.ScanType().Get(0)
	dev.listenDBusSignals(systemSigLoop)
	return &dev
}

func (dev *Device) listenDBusSignals(sigLoop *dbusutil.SignalLoop) {
	dev.core.InitSignalExt(sigLoop, true)
	_, err := dev.core.ConnectEnrollStatus(func(status string, ok bool) {
		err := dev.service.Emit(dev, "EnrollStatus", status, ok)
		if err != nil {
			logger.Warning(err)
		}
	})
	if err != nil {
		logger.Warning(err)
	}

	_, err = dev.core.ConnectVerifyStatus(func(status string, ok bool) {
		err := dev.service.Emit(dev, "VerifyStatus", status, ok)
		if err != nil {
			logger.Warning(err)
		}
	})
	if err != nil {
		logger.Warning(err)
	}

	_, err = dev.core.ConnectVerifyFingerSelected(func(finger string) {
		err := dev.service.Emit(dev, "VerifyFingerSelected", finger)
		if err != nil {
			logger.Warning(err)
		}
	})
	if err != nil {
		logger.Warning(err)
	}
}

func (dev *Device) destroy() {
	dev.core.RemoveHandler(proxy.RemoveAllHandlers)
	err := dev.service.StopExport(dev)
	if err != nil {
		logger.Warning(err)
	}
}

func (dev *Device) Claim(username string) *dbus.Error {
	err := dev.core.Claim(0, username)
	return dbusutil.ToError(err)
}

func (dev *Device) Release() *dbus.Error {
	err := dev.core.Release(0)
	return dbusutil.ToError(err)
}

func (dev *Device) EnrollStart(sender dbus.Sender, finger string) *dbus.Error {
	ok, err := checkAuth(actionIdEnroll, string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	if !ok {
		err = errors.New("authentication failed")
		return dbusutil.ToError(err)
	}

	err = dev.core.EnrollStart(0, finger)
	return dbusutil.ToError(err)
}

func (dev *Device) EnrollStop(sender dbus.Sender) *dbus.Error {
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

func (dev *Device) DeleteEnrolledFingers(sender dbus.Sender, username string) *dbus.Error {
	ok, err := checkAuth(actionIdDelete, string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	if !ok {
		err = errors.New("authentication failed")
		return dbusutil.ToError(err)
	}

	err = dev.core.DeleteEnrolledFingers(0, username)
	return dbusutil.ToError(err)
}

func (dev *Device) DeleteEnrolledFinger(sender dbus.Sender, username string, finger string) *dbus.Error {
	return dbusutil.ToError(errors.New("can not delete fprintd single finger"))
}

func (dev *Device) GetCapabilities() ([]string, *dbus.Error) {
	return nil, nil
}

func (dev *Device) ClaimForce(sender dbus.Sender, username string) *dbus.Error {
	return dbusutil.ToError(errors.New("can not claim force"))
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

func (dev *Device) getCorePath() dbus.ObjectPath {
	return dev.core.Path_()
}

func destroyDevices(list Devices) {
	for _, dev := range list {
		dev.destroy()
	}
}

func (devList Devices) Add(objPath dbus.ObjectPath, service *dbusutil.Service,
	systemSigLoop *dbusutil.SignalLoop) Devices {
	var v = newDevice(objPath, service, systemSigLoop)
	err := service.Export(v.getPath(), v)
	if err != nil {
		logger.Warning("failed to export:", objPath)
		return devList
	}

	devList = append(devList, v)
	return devList
}

func (devList Devices) Get(objPath dbus.ObjectPath) IDevice {
	for _, dev := range devList {
		if dev.getCorePath() == objPath {
			return dev
		}
	}
	return nil
}

func (devList Devices) Delete(objPath dbus.ObjectPath) Devices {
	var (
		list Devices
		v    IDevice
	)
	for _, dev := range devList {
		if dev.getCorePath() == objPath {
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
