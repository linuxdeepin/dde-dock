/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package soundeffect

import (
	"gir/gio-2.0"
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
	"sync"
)

const (
	soundEffectSchema = "com.deepin.dde.sound-effect"
	settingKeyEnabled = "enabled"

	DBusDest = "com.deepin.daemon.SoundEffect"
	dbusPath = "/com/deepin/daemon/SoundEffect"
	dbusIFC  = DBusDest
)

type Manager struct {
	Enabled *property.GSettingsBoolProperty `access:"readwrite"`
	setting *gio.Settings
	count   int
	mutex   sync.Mutex
}

func NewManager() *Manager {
	var m = new(Manager)

	m.setting = gio.NewSettings(soundEffectSchema)
	m.Enabled = property.NewGSettingsBoolProperty(m, "Enabled", m.setting, settingKeyEnabled)
	return m
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       DBusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) PlaySystemSound(event string) {
	if len(event) == 0 {
		return
	}

	go func() {
		m.mutex.Lock()
		m.count++
		logger.Debug("start", m.count)
		m.mutex.Unlock()

		err := soundutils.PlaySystemSound(event, "", true)
		if err != nil {
			logger.Error(err)
		}

		m.mutex.Lock()
		logger.Debug("end", m.count)
		m.count--
		m.mutex.Unlock()
	}()
}
