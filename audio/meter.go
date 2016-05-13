/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import "pkg.deepin.io/lib/dbus"
import "pkg.deepin.io/lib/pulse"
import "fmt"
import "time"
import "sync"

type Meter struct {
	Volume  float64
	id      string
	hasTick bool
	core    *pulse.SourceMeter
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

//TODO: use pulse.Meter instead of remove pulse.SourceMeter
func NewMeter(id string, core *pulse.SourceMeter) *Meter {
	m := &Meter{id: id, core: core}
	m.Tick()
	go m.tryQuit()
	return m
}

var (
	meterLocker sync.Mutex
	meters      = make(map[string]*Meter)
)

func (s *Source) GetMeter() *Meter {
	meterLocker.Lock()
	defer meterLocker.Unlock()
	id := fmt.Sprintf("source%d", s.core.Index)
	m, ok := meters[id]
	if !ok {
		core := pulse.NewSourceMeter(pulse.GetContext(), s.core.Index)
		m = NewMeter(id, core)
		dbus.InstallOnSession(m)
		meters[id] = m
		core.ConnectChanged(func(v float64) {
			m.setPropVolume(v)
		})
	}
	return m
}

func (s *Sink) GetMeter() *Meter {
	//TODO
	return nil
}
