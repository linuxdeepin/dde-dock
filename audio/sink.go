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

package audio

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/pulse"
	"strings"
)

type Sink struct {
	core  *pulse.Sink
	index uint32

	Name        string
	Description string

	// 默认音量值
	BaseVolume float64

	// 是否静音
	Mute bool

	// 当前音量
	Volume float64
	// 左右声道平衡值
	Balance float64
	// 是否支持左右声道调整
	SupportBalance bool
	// 前后声道平衡值
	Fade float64
	// 是否支持前后声道调整
	SupportFade bool

	// 支持的输出端口
	Ports []Port
	// 当前使用的输出端口
	ActivePort Port
	// 声卡的索引
	Card uint32
}

func NewSink(core *pulse.Sink) *Sink {
	s := &Sink{core: core}
	s.index = s.core.Index
	s.update()
	return s
}

// 设置音量大小
//
// v: 音量大小
//
// isPlay: 是否播放声音反馈
func (s *Sink) SetVolume(v float64, isPlay bool) error {
	if !isVolumeValid(v) {
		return fmt.Errorf("Invalid volume value: %v", v)
	}

	if v == 0 {
		v = 0.001
	}
	s.core.SetVolume(s.core.Volume.SetAvg(v))
	if isPlay {
		playFeedbackWithDevice(s.Name)
	}
	return nil
}

// 设置左右声道平衡值
//
// v: 声道平衡值
//
// isPlay: 是否播放声音反馈
func (s *Sink) SetBalance(v float64, isPlay bool) error {
	if v < -1.00 || v > 1.00 {
		return fmt.Errorf("Invalid volume value: %v", v)
	}

	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
	if isPlay {
		playFeedbackWithDevice(s.Name)
	}
	return nil
}

// 设置前后声道平衡值
//
// v: 声道平衡值
//
// isPlay: 是否播放声音反馈
func (s *Sink) SetFade(v float64) error {
	if v < -1.00 || v > 1.00 {
		return fmt.Errorf("Invalid volume value: %v", v)
	}

	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
	playFeedbackWithDevice(s.Name)
	return nil
}

// 是否静音
func (s *Sink) SetMute(v bool) {
	logger.Debugf("Sink #%d SetMute %v", s.index, v)
	s.core.SetMute(v)
	if !v {
		playFeedbackWithDevice(s.Name)
	}
}

// 设置此设备的当前使用端口
func (s *Sink) SetPort(name string) {
	s.core.SetPort(name)
}

func (s *Sink) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       baseBusName,
		ObjectPath: fmt.Sprintf("%s/Sink%d", baseBusPath, s.index),
		Interface:  baseBusIfc + ".Sink",
	}
}

// if prev active port is 'headset-output',
// after unpluged, will lost the port info
var _prevSinkActivePort Port

func (s *Sink) update() {
	s.Name = s.core.Name
	s.Description = s.core.Description
	s.Card = s.core.Card
	s.BaseVolume = s.core.BaseVolume.ToPercent()

	s.setPropMute(s.core.Mute)
	s.setPropVolume(floatPrecision(s.core.Volume.Avg()))

	s.setPropSupportFade(false)
	s.setPropFade(s.core.Volume.Fade(s.core.ChannelMap))
	s.setPropSupportBalance(true)
	s.setPropBalance(s.core.Volume.Balance(s.core.ChannelMap))

	oldActivePort := s.ActivePort
	activePortChanged := s.setPropActivePort(toPort(s.core.ActivePort))

	var ports []Port
	for _, p := range s.core.Ports {
		ports = append(ports, toPort(p))
	}
	s.setPropPorts(ports)

	if activePortChanged {
		logger.Debugf("sink #%d active port changed, old %v, new %v", s.index, oldActivePort, s.ActivePort)
		// old port but has new available state
		oldPort, foundOldPort := getPortByName(s.Ports, oldActivePort.Name)
		var oldPortUnavailable bool
		if !foundOldPort {
			logger.Debug("Sink.update not found old port")
			oldPortUnavailable = true
		} else {
			oldPortUnavailable = (int(oldPort.Available) == pulse.AvailableTypeNo)
		}
		logger.Debugf("oldPortUnavailable: %v", oldPortUnavailable)
		handleUnplugedEvent(_prevSinkActivePort, s.ActivePort, oldPortUnavailable)
		_prevSinkActivePort = s.ActivePort
	}
}

func handleUnplugedEvent(oldActivePort, newActivePort Port, oldPortUnavailable bool) {
	logger.Debug("[handleUnplugedEvent] Old port:", oldActivePort.String(), oldPortUnavailable)
	logger.Debug("[handleUnplugedEvent] New port:", newActivePort.String())
	// old active port is headphone or bluetooth
	if isHeadphoneOrHeadsetPort(oldActivePort.Name) &&
		// old active port is not unavailable
		int(oldActivePort.Available) != pulse.AvailableTypeNo &&
		// new port is not headphone and bluetooth
		!isHeadphoneOrHeadsetPort(newActivePort.Name) && oldPortUnavailable {
		pauseAllPlayers()
	}
}

func isHeadphoneOrHeadsetPort(portName string) bool {
	name := strings.ToLower(portName)
	return strings.Contains(name, "headphone") || strings.Contains(name, "headset-output")
}

// return whether changed
func (s *Sink) setPropActivePort(v Port) bool {
	if s.ActivePort != v {
		s.ActivePort = v
		dbus.NotifyChange(s, "ActivePort")
		return true
	}
	return false
}

func (s *Sink) setPropPorts(v []Port) {
	if !portsEqual(s.Ports, v) {
		s.Ports = v
		dbus.NotifyChange(s, "Ports")
	}
}

func (s *Sink) setPropVolume(v float64) {
	if s.Volume != v {
		s.Volume = v
		dbus.NotifyChange(s, "Volume")
	}
}

func (s *Sink) setPropBalance(v float64) {
	if s.Balance != v {
		s.Balance = v
		dbus.NotifyChange(s, "Balance")
	}
}

func (s *Sink) setPropSupportBalance(v bool) {
	if s.SupportBalance != v {
		s.SupportBalance = v
		dbus.NotifyChange(s, "SupportBalance")
	}
}

func (s *Sink) setPropSupportFade(v bool) {
	if s.SupportFade != v {
		s.SupportFade = v
		dbus.NotifyChange(s, "SupportFade")
	}
}

func (s *Sink) setPropFade(v float64) {
	if s.Fade != v {
		s.Fade = v
		dbus.NotifyChange(s, "Fade")
	}
}

func (s *Sink) setPropMute(v bool) {
	if s.Mute != v {
		s.Mute = v
		dbus.NotifyChange(s, "Mute")
	}
}

func (s *Sink) GetMeter() *Meter {
	//TODO
	return nil
}
