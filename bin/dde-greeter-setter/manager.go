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

package main

import (
	"encoding/json"
	"os"
	"sync"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/keyfile"
)

const (
	dbusServiceName = "com.deepin.daemon.Greeter"
	dbusPath        = "/com/deepin/daemon/Greeter"
	dbusInterface   = dbusServiceName

	greeterConfigFile = "/etc/lightdm/lightdm-deepin-greeter.conf"
	kfGroupGeneral    = "General"
	kfKeyScaleFactor  = "ScreenScaleFactor"
	kfKeyScaleFactors = "ScreenScaleFactors"
)

type Manager struct {
	service *dbusutil.Service
	mu      sync.Mutex
	kf      *keyfile.KeyFile

	methods *struct {
		SetScaleFactor        func() `in:"factor"`
		GetScaleFactor        func() `out:"factor"`
		SetScreenScaleFactors func() `in:"factors"`
		GetScreenScaleFactors func() `out:"factors"`
	}
}

func (m *Manager) getKeyFile() *keyfile.KeyFile {
	if m.kf == nil {
		m.kf = keyfile.NewKeyFile()
		err := m.kf.LoadFromFile(greeterConfigFile)
		if err != nil && !os.IsNotExist(err) {
			logger.Warning(err)
		}
	}
	return m.kf
}

func (m *Manager) SetScaleFactor(scale float64) *dbus.Error {
	m.service.DelayAutoQuit()
	m.mu.Lock()
	defer m.mu.Unlock()

	kf := m.getKeyFile()
	value, err := kf.GetFloat64(kfGroupGeneral, kfKeyScaleFactor)
	if err == nil && (value > scale-0.01 && value < scale+0.01) {
		return nil
	}
	kf.SetFloat64(kfGroupGeneral, kfKeyScaleFactor, scale)
	err = kf.SaveToFile(greeterConfigFile)
	return dbusutil.ToError(err)
}

func (m *Manager) GetScaleFactor() (float64, *dbus.Error) {
	m.service.DelayAutoQuit()
	m.mu.Lock()
	defer m.mu.Unlock()

	kf := m.getKeyFile()
	value, err := kf.GetFloat64(kfGroupGeneral, kfKeyScaleFactor)
	if err != nil {
		return 1, nil
	}
	return value, nil
}

func (m *Manager) GetScreenScaleFactors() (map[string]float64, *dbus.Error) {
	m.service.DelayAutoQuit()
	m.mu.Lock()
	defer m.mu.Unlock()

	var factors map[string]float64
	kf := m.getKeyFile()
	value, err := kf.GetString(kfGroupGeneral, kfKeyScaleFactors)
	if err != nil {
		logger.Warning(err)
		return nil, nil
	}
	err = json.Unmarshal([]byte(value), &factors)
	if err != nil {
		logger.Warning(err)
		return nil, dbusutil.ToError(err)
	}

	return factors, nil
}

func (m *Manager) SetScreenScaleFactors(factors map[string]float64) *dbus.Error {
	m.service.DelayAutoQuit()
	m.mu.Lock()
	defer m.mu.Unlock()

	kf := m.getKeyFile()
	value, err := json.Marshal(factors)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}
	kf.SetString(kfGroupGeneral, kfKeyScaleFactors, string(value))
	err = kf.SaveToFile(greeterConfigFile)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}
	return nil
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}
