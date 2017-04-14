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
	"fmt"
	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/pulse"
	dutils "pkg.deepin.io/lib/utils"
	"sync"
)

const (
	audioSchema       = "com.deepin.dde.audio"
	gsKeyFirstRun     = "first-run"
	gsKeyInputVolume  = "input-volume"
	gsKeyOutputVolume = "output-volume"

	soundEffectSchema     = "com.deepin.dde.sound-effect"
	soundEffectKeyEnabled = "enabled"
)

var (
	defaultInputVolume  float64
	defaultOutputVolume float64
)

type Audio struct {
	init bool
	core *pulse.Context
	// 正常输出声音的程序列表
	SinkInputs    []*SinkInput
	Cards         string
	DefaultSink   *Sink
	DefaultSource *Source

	// 最大音量
	MaxUIVolume float64
	cards       CardInfos

	siEventChan  chan func()
	siPollerExit chan struct{}

	isSaving    bool
	saverLocker sync.Mutex

	sinkLocker sync.Mutex
	portLocker sync.Mutex
}

func NewAudio(core *pulse.Context) *Audio {
	a := &Audio{core: core}
	a.MaxUIVolume = pulse.VolumeUIMax
	a.siEventChan = make(chan func(), 10)
	a.siPollerExit = make(chan struct{})
	go func() {
		a.update()
		a.applyConfig()
		a.initEventHandlers()
		a.sinkInputPoller()
	}()
	return a
}

func initDefaultVolume(audio *Audio) {
	setting := gio.NewSettings(audioSchema)
	defer setting.Unref()

	inVolumePer := float64(setting.GetInt(gsKeyInputVolume)) / 100.0
	outVolumePer := float64(setting.GetInt(gsKeyOutputVolume)) / 100.0
	defaultInputVolume = pulse.VolumeUIMax * inVolumePer
	defaultOutputVolume = pulse.VolumeUIMax * outVolumePer

	if !setting.GetBoolean(gsKeyFirstRun) {
		return
	}

	setting.SetBoolean(gsKeyFirstRun, false)
	audio.Reset()
}

func (a *Audio) SetDefaultSink(name string) {
	a.sinkLocker.Lock()
	defer a.sinkLocker.Unlock()

	if a.DefaultSink != nil && a.DefaultSink.Name == name {
		a.moveSinkInputsToSink(a.DefaultSink.index)
		return
	}

	logger.Debugf("audio.core.SetDefaultSink name: %q", name)
	a.core.SetDefaultSink(name)
	a.update()
	a.saveConfig()
	a.moveSinkInputsToSink(a.DefaultSink.index)
}

func (a *Audio) SetDefaultSource(name string) {
	if a.DefaultSource != nil && a.DefaultSource.Name == name {
		return
	}
	a.core.SetDefaultSource(name)
	a.update()
	a.saveConfig()
}

func (a *Audio) getActiveSinkPort() string {
	if a.DefaultSink == nil {
		return ""
	}

	for _, sink := range a.core.GetSinkList() {
		if sink.Name == a.DefaultSink.Name {
			return sink.ActivePort.Name
		}
	}
	return ""
}

// try set default sink and sink active port
func (a *Audio) trySetDefaultSink(cardId uint32, portName string) error {
	logger.Debugf("trySetDefaultSink card #%d port %q", cardId, portName)
	for _, sink := range a.core.GetSinkList() {
		if !isPortExists(portName, sink.Ports) || sink.Card != cardId {
			continue
		}
		if sink.ActivePort.Name != portName {
			sink.SetPort(portName)
		}
		a.SetDefaultSink(sink.Name)
		return nil
	}
	return fmt.Errorf("Cann't find valid sink for port '%s'", portName)
}

func (a *Audio) getActiveSourcePort() string {
	if a.DefaultSource == nil {
		return ""
	}

	for _, source := range a.core.GetSourceList() {
		if source.Name == a.DefaultSource.Name {
			return source.ActivePort.Name
		}
	}
	return ""
}

// try set default source and source active port
func (a *Audio) trySetDefaultSource(cardId uint32, portName string) error {
	logger.Debugf("trySetDefaultSource card #%d port %q", cardId, portName)
	for _, source := range a.core.GetSourceList() {
		if !isPortExists(portName, source.Ports) || source.Card != cardId {
			continue
		}
		if source.ActivePort.Name != portName {
			source.SetPort(portName)
		}
		if a.DefaultSource == nil || a.DefaultSource.Name != source.Name {
			a.SetDefaultSource(source.Name)
		}
		return nil
	}
	return fmt.Errorf("Cann't find valid source for port '%s'", portName)
}

// SetPort activate the port for the special card.
// The available sinks and sources will also change with the profile changing.
func (a *Audio) SetPort(cardId uint32, portName string, direction int32) error {
	logger.Debugf("Audio.SetPort card idx: %d, port name: %q, direction: %d", cardId, portName, direction)
	a.portLocker.Lock()
	defer a.portLocker.Unlock()
	var (
		curCard           *CardInfo
		oppositePort      string
		oppositeDirection int
	)
	if int(direction) == pulse.DirectionSink {
		oppositePort = a.getActiveSourcePort()
		oppositeDirection = pulse.DirectionSource
	} else if int(direction) == pulse.DirectionSource {
		oppositePort = a.getActiveSinkPort()
		oppositeDirection = pulse.DirectionSink
	} else {
		return fmt.Errorf("Invalid port direction: %d", direction)
	}
	curCard, _ = a.cards.get(cardId)
	if curCard != nil {
		_, err := curCard.Ports.Get(oppositePort, oppositeDirection)
		// curCard does not have the port
		if err != nil {
			curCard = nil
		}
	}

	var (
		destCard     *CardInfo
		destPortInfo pulse.CardPortInfo
		destProfile  string

		oppositePortInfo pulse.CardPortInfo
		commonProfiles   pulse.ProfileInfos2
	)
	if curCard != nil {
		portInfo, err := curCard.Ports.Get(portName, int(direction))
		if err != nil {
			return err
		}
		if portInfo.Profiles.Exists(curCard.ActiveProfile.Name) {
			goto setPort
		}
		destCard = curCard
		destPortInfo = portInfo
	} else {
		tmpCard, err := a.cards.get(cardId)
		if err != nil {
			return err
		}
		portInfo, err := tmpCard.Ports.Get(portName, int(direction))
		if err != nil {
			return err
		}
		destCard = tmpCard
		destPortInfo = portInfo
	}

	// match the common profile contain sinkPort and sourcePort
	oppositePortInfo, _ = destCard.Ports.Get(oppositePort, oppositeDirection)
	commonProfiles = getCommonProfiles(destPortInfo, oppositePortInfo)
	if len(commonProfiles) != 0 {
		destProfile = commonProfiles[0].Name
	} else {
		name, err := destCard.tryGetProfileByPort(portName)
		if err != nil {
			return err
		}
		destProfile = name
	}
	destCard.core.SetProfile(destProfile)

setPort:
	if int(direction) == pulse.DirectionSink {
		return a.trySetDefaultSink(cardId, portName)
	}
	return a.trySetDefaultSource(cardId, portName)
}

func (a *Audio) Reset() {
	for _, s := range a.core.GetSinkList() {
		s.SetMute(false)
		s.SetVolume(s.Volume.SetAvg(defaultOutputVolume))
		s.SetVolume(s.Volume.SetBalance(s.ChannelMap, 0))
		s.SetVolume(s.Volume.SetFade(s.ChannelMap, 0))
	}
	for _, s := range a.core.GetSourceList() {
		s.SetMute(false)
		s.SetVolume(s.Volume.SetAvg(defaultInputVolume))
		s.SetVolume(s.Volume.SetBalance(s.ChannelMap, 0))
		s.SetVolume(s.Volume.SetFade(s.ChannelMap, 0))
	}

	settings, err := dutils.CheckAndNewGSettings(soundEffectSchema)
	if err != nil {
		return
	}
	settings.SetBoolean(soundEffectKeyEnabled, true)
	settings.Unref()
}

func (a *Audio) destroy() {
	close(a.siPollerExit)
	dbus.UnInstallObject(a)
}

func (a *Audio) moveSinkInputsToSink(sinkId uint32) {
	// TODO: locker sinkinputs changed
	var list []uint32
	for _, sinkInput := range a.SinkInputs {
		if sinkInput.core.Sink == sinkId {
			continue
		}
		list = append(list, sinkInput.index)
	}
	if len(list) == 0 {
		return
	}
	a.core.MoveSinkInputsByIndex(list, sinkId)
}

func isPortExists(name string, ports []pulse.PortInfo) bool {
	for _, port := range ports {
		if port.Name == name {
			return true
		}
	}
	return false
}
