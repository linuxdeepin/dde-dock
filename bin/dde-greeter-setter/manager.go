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
	"fmt"
	"os"
	"strconv"
	"strings"
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

	kf := m.getKeyFile()
	value, err := kf.GetValue(kfGroupGeneral, kfKeyScaleFactors)
	if err != nil {
		logger.Warning(err)
		return nil, nil
	}
	value, err = strconv.Unquote(value)
	if err != nil {
		logger.Warning(err)
		return nil, nil
	}
	factors := parseIndividualScaling(value)
	return factors, nil
}

func (m *Manager) SetScreenScaleFactors(factors map[string]float64) *dbus.Error {
	m.service.DelayAutoQuit()
	m.mu.Lock()
	defer m.mu.Unlock()

	kf := m.getKeyFile()
	value := joinIndividualScaling(factors)
	value = strconv.Quote(value)
	kf.SetValue(kfGroupGeneral, kfKeyScaleFactors, value)
	err := kf.SaveToFile(greeterConfigFile)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}
	return nil
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func parseIndividualScaling(str string) map[string]float64 {
	pairs := strings.Split(str, ";")
	result := make(map[string]float64)
	for _, value := range pairs {
		kv := strings.SplitN(value, "=", 2)
		if len(kv) != 2 {
			continue
		}

		value, err := strconv.ParseFloat(kv[1], 64)
		if err != nil {
			logger.Warning(err)
			continue
		}

		result[kv[0]] = value
	}

	return result
}

func joinIndividualScaling(v map[string]float64) string {
	pairs := make([]string, len(v))
	idx := 0
	for key, value := range v {
		pairs[idx] = fmt.Sprintf("%s=%.2f", key, value)
		idx++
	}
	return strings.Join(pairs, ";")
}
