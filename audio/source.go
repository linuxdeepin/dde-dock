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

	"sync"

	"strconv"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/pulse"
)

type Source struct {
	service     *dbusutil.Service
	PropsMu     sync.RWMutex
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
	// dbusutil-gen: equal=portsEqual
	Ports      []Port
	ActivePort Port
	// 声卡的索引
	Card uint32

	methods *struct {
		SetVolume  func() `in:"value,isPlay"`
		SetBalance func() `in:"value,isPlay"`
		SetFade    func() `in:"value"`
		SetMute    func() `in:"value"`
		SetPort    func() `in:"name"`
		GetMeter   func() `out:"meter"`
	}
}

func NewSource(core *pulse.Source, service *dbusutil.Service) *Source {
	s := &Source{
		core:    core,
		service: service,
	}
	s.index = s.core.Index
	s.update()
	return s
}

// 如何反馈输入音量？
func (s *Source) SetVolume(v float64, isPlay bool) *dbus.Error {
	if !isVolumeValid(v) {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
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

func (s *Source) SetBalance(v float64, isPlay bool) *dbus.Error {
	if v < -1.00 || v > 1.00 {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
	if isPlay {
		playFeedback()
	}
	return nil
}

func (s *Source) SetFade(v float64) *dbus.Error {
	if v < -1.00 || v > 1.00 {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
	playFeedback()
	return nil
}

func (s *Source) SetMute(v bool) *dbus.Error {
	s.core.SetMute(v)
	if !v {
		playFeedback()
	}
	return nil
}

func (s *Source) SetPort(name string) *dbus.Error {
	s.core.SetPort(name)
	return nil
}

func (s *Source) GetMeter() (dbus.ObjectPath, *dbus.Error) {
	meterLocker.Lock()
	defer meterLocker.Unlock()
	id := fmt.Sprintf("source%d", s.core.Index)
	m, ok := meters[id]
	var meterPath dbus.ObjectPath
	if !ok {
		core := pulse.NewSourceMeter(pulse.GetContext(), s.core.Index)
		m = NewMeter(id, core, s.service)
		meterPath = m.getPath()
		s.service.Export(meterPath, m)

		meters[id] = m
		core.ConnectChanged(func(v float64) {
			m.PropsMu.Lock()
			m.setPropVolume(v)
			m.PropsMu.Unlock()
		})
	} else {
		meterPath = m.getPath()
	}
	return meterPath, nil
}

func (s *Source) getPath() dbus.ObjectPath {
	return dbus.ObjectPath(dbusPath + "/Source" + strconv.Itoa(int(s.index)))
}

func (*Source) GetInterfaceName() string {
	return dbusInterface + ".Source"
}

func (s *Source) update() {
	s.Name = s.core.Name
	s.Description = s.core.Description
	s.Card = s.core.Card
	s.BaseVolume = s.core.BaseVolume.ToPercent()

	s.PropsMu.Lock()
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

	s.PropsMu.Unlock()
}
