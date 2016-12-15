/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/pulse"
	"sync"
	"time"
)

var (
	meterLocker sync.Mutex
	meters      = make(map[string]*Meter)
)

type Meter struct {
	Volume  float64
	id      string
	hasTick bool
	core    *pulse.SourceMeter
}

//TODO: use pulse.Meter instead of remove pulse.SourceMeter
func NewMeter(id string, core *pulse.SourceMeter) *Meter {
	m := &Meter{id: id, core: core}
	m.Tick()
	go m.tryQuit()
	return m
}

func (m *Meter) quit() {
	delete(meters, m.id)
	dbus.UnInstallObject(m)
	m.core.Destroy()
}

func (m *Meter) tryQuit() {
	defer m.quit()

	for {
		select {
		case <-time.After(time.Second * 10):
			if !m.hasTick {
				return
			}
			m.hasTick = false
		}
	}
}

func (m *Meter) Tick() {
	m.hasTick = true
}

func (m *Meter) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       baseBusName,
		ObjectPath: baseBusPath + "/Meter" + m.id,
		Interface:  baseBusIfc + ".Meter",
	}
}

func (m *Meter) setPropVolume(v float64) {
	if m.Volume != v {
		m.Volume = v
		dbus.NotifyChange(m, "Volume")
	}
}
