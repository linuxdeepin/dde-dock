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
	"sync"
	"time"

	"github.com/linuxdeepin/go-dbus-factory/net.reactivated.fprint"
	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	polkit "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName     = "com.deepin.daemon.Fprintd"
	dbusPath            = "/com/deepin/daemon/Fprintd"
	dbusInterface       = dbusServiceName
	dbusDeviceInterface = dbusServiceName + ".Device"

	systemdDBusServiceName = "org.freedesktop.systemd1"
	systemdDBusPath        = "/org/freedesktop/systemd1"
	systemdDBusInterface   = systemdDBusServiceName + ".Manager"
)

//go:generate dbusutil-gen -type Manager -import pkg.deepin.io/lib/dbus1 manager.go

type Manager struct {
	service       *dbusutil.Service
	sysSigLoop    *dbusutil.SignalLoop
	fprintManager *fprint.Manager
	dbusDaemon    *ofdbus.DBus
	devices       Devices
	devicesMu     sync.Mutex
	fprintCh      chan struct{}

	PropsMu sync.RWMutex
	// dbusutil-gen: equal=nil
	Devices []dbus.ObjectPath

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
		fprintManager: fprint.NewManager(systemConn),
		dbusDaemon:    ofdbus.NewDBus(systemConn),
		sysSigLoop:    dbusutil.NewSignalLoop(systemConn, 10),
	}, nil
}

func (m *Manager) GetDefaultDevice() (dbus.ObjectPath, *dbus.Error) {
	objPath, err := m.fprintManager.GetDefaultDevice(0)
	if err != nil {
		logger.Warning("Failed to get default device:", err)
		return "/", dbusutil.ToError(err)
	}
	m.addDevice(objPath)
	return convertFPrintPath(objPath), nil
}

func (m *Manager) GetDevices() ([]dbus.ObjectPath, *dbus.Error) {
	err := m.refreshDevices()
	if err != nil {
		return nil, dbusutil.ToError(err)
	}
	m.PropsMu.Lock()
	paths := m.Devices
	m.PropsMu.Unlock()
	return paths, nil
}

func (m *Manager) refreshDevices() error {
	devicePaths, err := m.fprintManager.GetDevices(0)
	if err != nil {
		return err
	}

	var needDelete []dbus.ObjectPath
	var needAdd []dbus.ObjectPath

	m.devicesMu.Lock()

	// 在 m.devList 但不在 devicePaths 中的记录在 needDelete
	for _, d := range m.devices {
		found := false
		for _, devPath := range devicePaths {
			if d.core.Path_() == devPath {
				found = true
				break
			}
		}
		if !found {
			needDelete = append(needDelete, d.core.Path_())
		}
	}

	// 在 devicePaths 但不在 m.devList 中的记录在 needAdd
	for _, devPath := range devicePaths {
		found := false
		for _, d := range m.devices {
			if d.core.Path_() == devPath {
				found = true
				break
			}
		}
		if !found {
			needAdd = append(needAdd, devPath)
		}
	}

	for _, devPath := range needDelete {
		m.devices = m.devices.Delete(devPath)
	}
	for _, devPath := range needAdd {
		m.devices = m.devices.Add(devPath, m.service, m.sysSigLoop)
	}
	m.devicesMu.Unlock()

	m.updatePropDevices()
	return nil
}

func (m *Manager) updatePropDevices() {
	m.devicesMu.Lock()
	paths := make([]dbus.ObjectPath, len(m.devices))
	for idx, d := range m.devices {
		paths[idx] = d.getPath()
	}
	m.devicesMu.Unlock()

	m.PropsMu.Lock()
	m.setPropDevices(paths)
	m.PropsMu.Unlock()
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) init() {
	m.sysSigLoop.Start()
	m.fprintCh = make(chan struct{}, 1)
	m.listenDBusSignals()

	paths, err := m.fprintManager.GetDevices(0)
	if err != nil {
		logger.Warning("Failed to get fprint devices:", err)
		return
	}
	for _, devPath := range paths {
		m.addDevice(devPath)
	}
	m.updatePropDevices()
}

func (m *Manager) listenDBusSignals() {
	m.dbusDaemon.InitSignalExt(m.sysSigLoop, true)
	_, err := m.dbusDaemon.ConnectNameOwnerChanged(func(name string, oldOwner string, newOwner string) {
		fprintDBusServiceName := m.fprintManager.ServiceName_()
		if name == fprintDBusServiceName && newOwner != "" {
			select {
			case m.fprintCh <- struct{}{}:
			default:
			}
		}
	})
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) addDevice(objPath dbus.ObjectPath) {
	logger.Debug("add device:", objPath)
	m.devicesMu.Lock()
	defer m.devicesMu.Unlock()

	d := m.devices.Get(objPath)
	if d != nil {
		return
	}
	m.devices = m.devices.Add(objPath, m.service, m.sysSigLoop)
}

func (m *Manager) destroy() {
	destroyDevices(m.devices)
	m.sysSigLoop.Stop()
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

func (m *Manager) TriggerUDevEvent(sender dbus.Sender) *dbus.Error {
	uid, err := m.service.GetConnUID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	if uid != 0 {
		err = errors.New("not root user")
		return dbusutil.ToError(err)
	}

	logger.Debug("udev event")

	select {
	case <-m.fprintCh:
	default:
	}

	err = restartSystemdService("fprintd.service", "replace")
	if err != nil {
		return dbusutil.ToError(err)
	}

	select {
	case <-m.fprintCh:
		logger.Debug("fprintd started")
	case <-time.After(5 * time.Second):
		logger.Warning("wait fprintd restart timed out!")
	}

	err = m.refreshDevices()
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}

func restartSystemdService(name, mode string) error {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	obj := sysBus.Object(systemdDBusServiceName, systemdDBusPath)
	var jobPath dbus.ObjectPath
	err = obj.Call(systemdDBusInterface+".RestartUnit", dbus.FlagNoAutoStart, name, mode).Store(&jobPath)
	return err
}
