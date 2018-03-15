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

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/pulse"
)

var (
	meterLocker sync.Mutex
	meters      = make(map[string]*Meter)
)

type Meter struct {
	service *dbusutil.Service
	PropsMu sync.RWMutex
	Volume  float64
	id      string
	hasTick bool
	core    *pulse.SourceMeter
}

//TODO: use pulse.Meter instead of remove pulse.SourceMeter
func NewMeter(id string, core *pulse.SourceMeter, service *dbusutil.Service) *Meter {
	m := &Meter{
		id:      id,
		core:    core,
		service: service,
	}
	m.Tick()
	go m.tryQuit()
	return m
}

func (m *Meter) quit() {
	delete(meters, m.id)
	m.service.StopExport(m)
	m.core.Destroy()
}

func (m *Meter) tryQuit() {
	defer m.quit()

	for {
		select {
		case _, ok := <-time.After(time.Second * 10):
			if !ok {
				logger.Error("Invalid time event")
				return
			}

			if !m.hasTick {
				return
			}
			m.hasTick = false
		}
	}
}

func (m *Meter) Tick() *dbus.Error {
	m.hasTick = true
	return nil
}

func (m *Meter) getPath() dbus.ObjectPath {
	return dbus.ObjectPath(dbusPath + "/Meter" + m.id)
}

func (*Meter) GetInterfaceName() string {
	return dbusInterface + ".Meter"
}
