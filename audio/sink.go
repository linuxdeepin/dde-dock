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
	"strconv"
	"strings"
	"sync"

	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/pulse"
)

type Sink struct {
	audio   *Audio
	service *dbusutil.Service
	PropsMu sync.RWMutex
	index   uint32

	Name        string
	Description string

	// 默认音量值
	BaseVolume float64

	// 是否静音
	Mute bool

	// 当前音量
	Volume     float64
	cVolume    pulse.CVolume
	channelMap pulse.ChannelMap
	// 左右声道平衡值
	Balance float64
	// 是否支持左右声道调整
	SupportBalance bool
	// 前后声道平衡值
	Fade float64
	// 是否支持前后声道调整
	SupportFade bool

	// dbusutil-gen: equal=portsEqual
	// 支持的输出端口
	Ports []Port
	// 当前使用的输出端口
	ActivePort Port
	// 声卡的索引
	Card uint32

	props map[string]string

	methods *struct {
		SetVolume  func() `in:"value,isPlay"`
		SetBalance func() `in:"value,isPlay"`
		SetFade    func() `in:"value"`
		SetMute    func() `in:"value"`
		SetPort    func() `in:"name"`
		GetMeter   func() `out:"meter"`
	}
}

func newSink(sinkInfo *pulse.Sink, audio *Audio) *Sink {
	s := &Sink{
		audio:   audio,
		service: audio.service,
		index:   sinkInfo.Index,
		props:   sinkInfo.PropList,
	}
	s.update(sinkInfo)
	return s
}

// 设置音量大小
//
// v: 音量大小
//
// isPlay: 是否播放声音反馈
func (s *Sink) SetVolume(v float64, isPlay bool) *dbus.Error {
	if !isVolumeValid(v) {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	if v == 0 {
		v = 0.001
	}
	s.PropsMu.Lock()
	cv := s.cVolume.SetAvg(v)
	s.PropsMu.Unlock()
	s.audio.context().SetSinkVolumeByIndex(s.index, cv)

	configKeeper.SetVolume(s.audio.getCardNameById(s.Card), s.ActivePort.Name, v)
	err := configKeeper.Save(configKeeperFile)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	if isPlay {
		s.playFeedback()
	}
	return nil
}

// 设置左右声道平衡值
//
// v: 声道平衡值
//
// isPlay: 是否播放声音反馈
func (s *Sink) SetBalance(v float64, isPlay bool) *dbus.Error {
	if v < -1.00 || v > 1.00 {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	s.PropsMu.RLock()
	cv := s.cVolume.SetBalance(s.channelMap, v)
	s.PropsMu.RUnlock()
	s.audio.context().SetSinkVolumeByIndex(s.index, cv)

	configKeeper.SetBalance(s.audio.getCardNameById(s.Card), s.ActivePort.Name, v)
	err := configKeeper.Save(configKeeperFile)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	if isPlay {
		s.playFeedback()
	}
	return nil
}

// 设置前后声道平衡值
//
// v: 声道平衡值
//
// isPlay: 是否播放声音反馈
func (s *Sink) SetFade(v float64) *dbus.Error {
	if v < -1.00 || v > 1.00 {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	s.PropsMu.RLock()
	cv := s.cVolume.SetFade(s.channelMap, v)
	s.PropsMu.RUnlock()
	s.audio.context().SetSinkVolumeByIndex(s.index, cv)
	s.playFeedback()
	return nil
}

// 是否静音
func (s *Sink) SetMute(v bool) *dbus.Error {
	logger.Debugf("Sink #%d SetMute %v", s.index, v)
	s.audio.context().SetSinkMuteByIndex(s.index, v)

	configKeeper.SetMute(s.audio.getCardNameById(s.Card), s.ActivePort.Name, v)
	err := configKeeper.Save(configKeeperFile)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	if !v {
		s.playFeedback()
	}
	return nil
}

// 设置此设备的当前使用端口
func (s *Sink) SetPort(name string) *dbus.Error {
	s.audio.context().SetSinkPortByIndex(s.index, name)
	return nil
}

func (s *Sink) getPath() dbus.ObjectPath {
	return dbus.ObjectPath(dbusPath + "/Sink" + strconv.Itoa(int(s.index)))
}

func (*Sink) GetInterfaceName() string {
	return dbusInterface + ".Sink"
}

func (s *Sink) update(sinkInfo *pulse.Sink) {
	s.PropsMu.Lock()

	s.Name = sinkInfo.Name
	s.Description = sinkInfo.Description
	s.Card = sinkInfo.Card
	s.BaseVolume = sinkInfo.BaseVolume.ToPercent()
	s.cVolume = sinkInfo.Volume
	s.channelMap = sinkInfo.ChannelMap

	s.setPropMute(sinkInfo.Mute)
	s.setPropVolume(floatPrecision(sinkInfo.Volume.Avg()))

	s.setPropSupportFade(false)
	s.setPropFade(sinkInfo.Volume.Fade(sinkInfo.ChannelMap))
	s.setPropSupportBalance(true)
	s.setPropBalance(sinkInfo.Volume.Balance(sinkInfo.ChannelMap))

	oldActivePort := s.ActivePort
	newActivePort := toPort(sinkInfo.ActivePort)
	activePortChanged := s.setPropActivePort(newActivePort)

	s.setPropPorts(toPorts(sinkInfo.Ports))
	s.props = sinkInfo.PropList
	s.PropsMu.Unlock()

	// TODO(jouyouyun): Sometimes the default sink not in the same card, so the activePortChanged inaccurate.
	// The right way is saved the last default sink active port, then judge whether equal.
	if s.audio.headphoneUnplugAutoPause && activePortChanged {
		logger.Debugf("sink #%d active port changed, old %v, new %v",
			s.index, oldActivePort, newActivePort)
		// old port but has new available state
		oldPort, foundOldPort := getPortByName(s.Ports, oldActivePort.Name)
		var oldPortUnavailable bool
		if !foundOldPort {
			logger.Debug("Sink.update not found old port")
			oldPortUnavailable = true
		} else {
			oldPortUnavailable = int(oldPort.Available) == pulse.AvailableTypeNo
		}
		logger.Debugf("oldPortUnavailable: %v", oldPortUnavailable)

		handleUnplugedEvent(oldActivePort, newActivePort, oldPortUnavailable)
	}
}

func handleUnplugedEvent(oldActivePort, newActivePort Port, oldPortUnavailable bool) {
	logger.Debug("[handleUnplugedEvent] Old port:", oldActivePort.String(), oldPortUnavailable)
	logger.Debug("[handleUnplugedEvent] New port:", newActivePort.String())
	// old active port is headphone or bluetooth
	if isHeadphoneOrHeadsetPort(oldActivePort.Name) &&
		// old active port available is yes or unknown, not no
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

func (s *Sink) GetMeter() (dbus.ObjectPath, *dbus.Error) {
	//TODO
	return "/", nil
}

func (s *Sink) playFeedback() {
	s.PropsMu.RLock()
	name := s.Name
	s.PropsMu.RUnlock()
	playFeedbackWithDevice(name)
}
