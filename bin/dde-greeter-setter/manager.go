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
	"sync"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/keyfile"
	"pkg.deepin.io/lib/utils"
)

const (
	dbusDest = "com.deepin.daemon.Greeter"
	dbusPath = "/com/deepin/daemon/Greeter"
	dbusIFC  = dbusDest

	greeterConfigFile = "/etc/lightdm/lightdm-deepin-greeter.conf"
	kfGroupGeneral    = "General"
	kfKeyScaleFactor  = "ScreenScaleFactor"
)

type Manager struct {
	service *dbusutil.Service
	quit    bool
	locker  sync.Mutex

	methods *struct {
		SetScaleFactor func() `in:"factor"`
		GetScaleFactor func() `out:"factor"`
	}
}

func (m *Manager) SetScaleFactor(scale float64) *dbus.Error {
	m.setQuitFlag(false)
	defer m.setQuitFlag(true)
	m.service.DelayAutoQuit()

	kf, err := newKeyfile(greeterConfigFile)
	if err != nil {
		return dbusutil.ToError(err)
	}

	value, err := kf.GetFloat64(kfGroupGeneral, kfKeyScaleFactor)
	if err == nil && (value > scale-0.01 && value < scale+0.01) {
		return nil
	}
	kf.SetFloat64(kfGroupGeneral, kfKeyScaleFactor, scale)
	err = kf.SaveToFile(greeterConfigFile)
	return dbusutil.ToError(err)
}

func (m *Manager) GetScaleFactor() (float64, *dbus.Error) {
	m.setQuitFlag(false)
	defer m.setQuitFlag(true)
	m.service.DelayAutoQuit()

	kf, err := newKeyfile(greeterConfigFile)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}

	value, err := kf.GetFloat64(kfGroupGeneral, kfKeyScaleFactor)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}
	return value, nil
}

func (*Manager) GetDBusExportInfo() dbusutil.ExportInfo {
	return dbusutil.ExportInfo{
		Path:      dbusPath,
		Interface: dbusIFC,
	}
}

func (m *Manager) setQuitFlag(v bool) {
	m.locker.Lock()
	m.quit = v
	m.locker.Unlock()
}

func (m *Manager) canQuit() bool {
	m.locker.Lock()
	defer m.locker.Unlock()
	return m.quit
}

var _kf *keyfile.KeyFile

func newKeyfile(file string) (*keyfile.KeyFile, error) {
	if _kf != nil {
		return _kf, nil
	}

	if !utils.IsFileExist(file) {
		err := utils.CreateFile(file)
		if err != nil {
			return nil, err
		}
	}

	_kf = keyfile.NewKeyFile()
	err := _kf.LoadFromFile(file)
	if err != nil {
		_kf = nil
		return nil, err
	}
	return _kf, nil
}
