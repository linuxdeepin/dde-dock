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
	"pkg.deepin.io/lib/dbus"
	"sync"
)

const (
	dbusDest      = "com.deepin.daemon.Fprintd"
	dbusPath      = "/com/deepin/daemon/Fprintd"
	dbusIFC       = dbusDest
	dbusDeviceIFC = dbusDest + ".Device"

	fprintDest       = "net.reactivated.Fprint"
	fprintPath       = "/net/reactivated/Fprint/"
	fprintManagerIFC = "net.reactivated.Fprint.Manager"
	fprintDeviceIFC  = "net.reactivated.Fprint.Device"
)

type Manager struct {
	core *fprint.Manager

	devList   Devices
	devLocker sync.Mutex
}

func newManager() *Manager {
	var m Manager
	m.core, _ = fprint.NewManager(fprintDest, fprintPath+"Manager")
	return &m
}

func (m *Manager) GetDefaultDevice() (dbus.ObjectPath, error) {
	objPath, err := m.core.GetDefaultDevice()
	if err != nil {
		logger.Warning("Failed to get default device:", err)
		return "/", err
	}
	m.addDevice(objPath)
	return convertFprintPath(objPath), nil
}

func (m *Manager) GetDevices() ([]dbus.ObjectPath, error) {
	list, err := m.core.GetDevices()
	if err != nil {
		return nil, err
	}
	var ret []dbus.ObjectPath
	for _, v := range list {
		m.addDevice(v)
		ret = append(ret, convertFprintPath(v))
	}
	return ret, nil
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) init() {
	list, err := m.core.GetDevices()
	if err != nil {
		logger.Warning("Failed to get fprint devices:", err)
		return
	}

	if len(list) == 0 {
		logger.Info("Not found fprint device")
		return
	}
	m.addDevices(list)
}

func (m *Manager) addDevice(objPath dbus.ObjectPath) {
	logger.Debug("Will add device:", objPath)
	m.devLocker.Lock()
	m.devList = m.devList.Add(objPath)
	m.devLocker.Unlock()
}

func (m *Manager) addDevices(pathList []dbus.ObjectPath) {
	logger.Debug("Will add device list:", pathList)
	m.devLocker.Lock()
	for _, v := range pathList {
		m.devList = m.devList.Add(v)
	}
	m.devLocker.Unlock()
}
