package main

import "dlib/dbus"
import "dlib/pulse"
import "fmt"
import "time"

type Meter struct {
	Volume float64
	id     string
}

func (m *Meter) Tick() {
	GetMeterManager().tick(m.id)
}

type MeterManager map[string]*Meter

var GetMeterManager = func() func() MeterManager {
	m := MeterManager(make(map[string]*Meter))
	return func() MeterManager {
		return m
	}
}()

func (mm MeterManager) clean() {
	select {
	case <-time.After(time.Second * 2):
	}
}
func (MeterManager) tick(id string) {
}
func (mm MeterManager) getSourceMeter(s *Source) *Meter {
	id := fmt.Sprintf("source%d", s.core.Index)
	m, ok := mm[id]
	if !ok {
		m = &Meter{id: id}
		dbus.InstallOnSession(m)
		mm[id] = m
		mm := pulse.NewSourceMeter(pulse.GetContext(), s.core.Index)
		fmt.Println("CreateMeter...", id)
		mm.ConnectChanged(func(v float64) {
			m.setPropVolume(v)
		})
	}
	return m
}

func (s *Source) GetMeter() *Meter {
	return GetMeterManager().getSourceMeter(s)
}

func (s *Sink) GetMeter() *Meter {
	//TODO
	return nil
}
