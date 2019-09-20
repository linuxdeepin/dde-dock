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
	"sync"

	"github.com/linuxdeepin/go-dbus-factory/net.reactivated.fprint"
	polkit "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName     = "com.deepin.daemon.Fprintd"
	dbusPath            = "/com/deepin/daemon/Fprintd"
	dbusInterface       = dbusServiceName
	dbusDeviceInterface = dbusServiceName + ".Device"

	fprintDBusServiceName  = "net.reactivated.Fprint"
	fprintDBusPath         = "/net/reactivated/Fprint/"
	fprintManagerInterface = "net.reactivated.Fprint.Manager"
	fprintDeviceInterface  = "net.reactivated.Fprint.Device"
)

type Manager struct {
	service       *dbusutil.Service
	systemSigLoop *dbusutil.SignalLoop
	core          *fprint.Manager
	devList       Devices
	devLocker     sync.Mutex

	methods *struct {
		GetDefaultDevice func() `out:"device"`
		GetDevices       func() `out:"devices"`
	}
}

func newManager(service *dbusutil.Service) (*Manager, error) {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	return &Manager{
		service:       service,
		core:          fprint.NewManager(systemConn),
		systemSigLoop: dbusutil.NewSignalLoop(systemConn, 10),
	}, nil
}

func (m *Manager) GetDefaultDevice() (dbus.ObjectPath, *dbus.Error) {
	objPath, err := m.core.GetDefaultDevice(0)
	if err != nil {
		logger.Warning("Failed to get default device:", err)
		return "/", dbusutil.ToError(err)
	}
	m.addDevice(objPath)
	return convertFPrintPath(objPath), nil
}

func (m *Manager) GetDevices() ([]dbus.ObjectPath, *dbus.Error) {
	devicePaths, err := m.core.GetDevices(0)
	if err != nil {
		return nil, dbusutil.ToError(err)
	}
	var ret []dbus.ObjectPath
	for _, devPath := range devicePaths {
		m.addDevice(devPath)
		ret = append(ret, convertFPrintPath(devPath))
	}
	return ret, nil
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) init() {
	m.systemSigLoop.Start()
	list, err := m.core.GetDevices(0)
	if err != nil {
		logger.Warning("Failed to get fprint devices:", err)
		return
	}

	if len(list) == 0 {
		logger.Info("Not found fprint device")
		return
	}
	list0 := make([]dbus.ObjectPath, len(list))
	for idx, devPath := range list {
		list0[idx] = devPath
	}
	m.addDevices(list0)
}

func (m *Manager) addDevice(objPath dbus.ObjectPath) {
	logger.Debug("Will add device:", objPath)
	m.devLocker.Lock()
	m.devList = m.devList.Add(objPath, m.service, m.systemSigLoop)
	m.devLocker.Unlock()
}

func (m *Manager) addDevices(pathList []dbus.ObjectPath) {
	logger.Debug("Will add device list:", pathList)
	m.devLocker.Lock()
	for _, v := range pathList {
		m.devList = m.devList.Add(v, m.service, m.systemSigLoop)
	}
	m.devLocker.Unlock()
}

func (m *Manager) destroy() {
	destroyDevices(m.devList)
	m.systemSigLoop.Stop()
}

func checkAuth(actionId string, busName string) (bool, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return false, err
	}
	authority := polkit.NewAuthority(systemBus)
	subject := polkit.MakeSubject(polkit.SubjectKindSystemBusName)
	subject.SetDetail("name", busName)

	ret, err := authority.CheckAuthorization(0, subject,
		actionId, nil,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}
	return ret.IsAuthorized, nil
}
