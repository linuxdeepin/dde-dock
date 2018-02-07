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
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/pulse"
)

type Source struct {
	core        *pulse.Source
	index       uint32
	Name        string
	Description string
	// 默认的输入音量
	BaseVolume     float64
	Mute           bool
	Volume         float64
	Balance        float64
	SupportBalance bool
	Fade           float64
	SupportFade    bool
	Ports          []Port
	ActivePort     Port
	// 声卡的索引
	Card uint32
}

func NewSource(core *pulse.Source) *Source {
	s := &Source{core: core}
	s.index = s.core.Index
	s.update()
	return s
}

// 如何反馈输入音量？
func (s *Source) SetVolume(v float64, isPlay bool) error {
	if !isVolumeValid(v) {
		return fmt.Errorf("Invalid volume value: %v", v)
	}

	if v == 0 {
		v = 0.001
	}
	s.core.SetVolume(s.core.Volume.SetAvg(v))
	if isPlay {
		playFeedback()
	}
	return nil
}

func (s *Source) SetBalance(v float64, isPlay bool) error {
	if v < -1.00 || v > 1.00 {
		return fmt.Errorf("Invalid volume value: %v", v)
	}

	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
	if isPlay {
		playFeedback()
	}
	return nil
}

func (s *Source) SetFade(v float64) error {
	if v < -1.00 || v > 1.00 {
		return fmt.Errorf("Invalid volume value: %v", v)
	}

	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
	playFeedback()
	return nil
}

func (s *Source) SetMute(v bool) {
	s.core.SetMute(v)
	if !v {
		playFeedback()
	}
}

func (s *Source) SetPort(name string) {
	s.core.SetPort(name)
}

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

func (s *Source) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       baseBusName,
		ObjectPath: fmt.Sprintf("%s/Source%d", baseBusPath, s.index),
		Interface:  baseBusIfc + ".Source",
	}
}

func (s *Source) update() {
	s.Name = s.core.Name
	s.Description = s.core.Description
	s.Card = s.core.Card
	s.BaseVolume = s.core.BaseVolume.ToPercent()

	s.setPropVolume(floatPrecision(s.core.Volume.Avg()))
	s.setPropMute(s.core.Mute)

	//TODO: handle this
	s.setPropSupportFade(false)
	s.setPropFade(s.core.Volume.Fade(s.core.ChannelMap))
	s.setPropSupportBalance(true)
	s.setPropBalance(s.core.Volume.Balance(s.core.ChannelMap))

	s.setPropActivePort(toPort(s.core.ActivePort))

	var ports []Port
	for _, p := range s.core.Ports {
		ports = append(ports, toPort(p))
	}
	s.setPropPorts(ports)
}

func (s *Source) setPropPorts(v []Port) {
	if !portsEqual(s.Ports, v) {
		s.Ports = v
		dbus.NotifyChange(s, "Ports")
	}
}

func (s *Source) setPropVolume(v float64) {
	if s.Volume != v {
		s.Volume = v
		dbus.NotifyChange(s, "Volume")
	}
}

func (s *Source) setPropSupportBalance(v bool) {
	if s.SupportBalance != v {
		s.SupportBalance = v
		dbus.NotifyChange(s, "SupportBalance")
	}
}

func (s *Source) setPropBalance(v float64) {
	if s.Balance != v {
		s.Balance = v
		dbus.NotifyChange(s, "Balance")
	}
}

func (s *Source) setPropSupportFade(v bool) {
	if s.SupportFade != v {
		s.SupportFade = v
		dbus.NotifyChange(s, "SupportFade")
	}
}

func (s *Source) setPropFade(v float64) {
	if s.Fade != v {
		s.Fade = v
		dbus.NotifyChange(s, "Fade")
	}
}

func (s *Source) setPropMute(v bool) {
	if s.Mute != v {
		s.Mute = v
		dbus.NotifyChange(s, "Mute")
	}
}

func (s *Source) setPropActivePort(v Port) {
	if s.ActivePort != v {
		s.ActivePort = v
		dbus.NotifyChange(s, "ActivePort")
	}
}
