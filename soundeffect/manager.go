/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"sync"

	"gir/gio-2.0"
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
)

const (
	soundEffectSchema = "com.deepin.dde.sound-effect"
	settingKeyEnabled = "enabled"

	DBusServiceName = "com.deepin.daemon.SoundEffect"
	dbusPath        = "/com/deepin/daemon/SoundEffect"
	dbusInterface   = DBusServiceName
)

type Manager struct {
	service *dbusutil.Service
	setting *gio.Settings
	count   int
	countMu sync.Mutex

	Enabled gsprop.Bool `prop:"access:rw"`

	methods *struct {
		PlaySystemSound func() `in:"event"`
	}
}

func NewManager(service *dbusutil.Service) *Manager {
	var m = new(Manager)

	m.service = service
	m.setting = gio.NewSettings(soundEffectSchema)
	m.Enabled.Bind(m.setting, settingKeyEnabled)
	return m
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) PlaySystemSound(event string) *dbus.Error {
	m.service.DelayAutoQuit()

	if event == "" {
		return nil
	}

	go func() {
		m.countMu.Lock()
		m.count++
		logger.Debug("start", m.count)
		m.countMu.Unlock()

		err := soundutils.PlaySystemSound(event, "")
		if err != nil {
			logger.Error(err)
		}

		m.countMu.Lock()
		logger.Debug("end", m.count)
		m.count--
		m.countMu.Unlock()
	}()
	return nil
}

func (m *Manager) canQuit() bool {
	m.countMu.Lock()
	count := m.count
	m.countMu.Unlock()
	return count == 0
}
