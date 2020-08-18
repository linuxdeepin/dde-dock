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

package audio

import (
	"sync"
	"time"

	dbus "github.com/godbus/dbus"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/pulse"
)

type Meter struct {
	audio   *Audio
	service *dbusutil.Service
	PropsMu sync.RWMutex
	Volume  float64
	id      string
	alive   bool
	core    *pulse.SourceMeter
}

//TODO: use pulse.Meter instead of remove pulse.SourceMeter
func newMeter(id string, core *pulse.SourceMeter, audio *Audio) *Meter {
	m := &Meter{
		id:      id,
		core:    core,
		audio:   audio,
		service: audio.service,
	}
	err := m.Tick()
	if err != nil {
		logger.Warning(err)
	}
	go m.tryQuit()
	return m
}

func (m *Meter) quit() {
	m.audio.mu.Lock()
	delete(m.audio.meters, m.id)
	m.audio.mu.Unlock()

	err := m.service.StopExport(m)
	if err != nil {
		logger.Warning(err)
	}
	m.core.Destroy()
}

func (m *Meter) tryQuit() {
	defer m.quit()

	for range time.After(time.Second * 10) {
		m.PropsMu.Lock()
		if !m.alive {
			m.PropsMu.Unlock()
			return
		}
		m.alive = false
		m.PropsMu.Unlock()
	}
}

func (m *Meter) Tick() *dbus.Error {
	m.PropsMu.Lock()
	m.alive = true
	m.PropsMu.Unlock()
	return nil
}

func (m *Meter) getPath() dbus.ObjectPath {
	return dbus.ObjectPath(dbusPath + "/Meter" + m.id)
}

func (*Meter) GetInterfaceName() string {
	return dbusInterface + ".Meter"
}
