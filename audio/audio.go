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
	"pkg.deepin.io/dde/api/soundutils"
	. "pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/pulse"
	dutils "pkg.deepin.io/lib/utils"
	"sync"
)

type Audio struct {
	init bool
	core *pulse.Context

	// 正常输出声音的程序列表
	SinkInputs []*SinkInput

	Cards string

	ActiveSinkPort   string
	ActiveSourcePort string

	DefaultSink   *Sink
	DefaultSource *Source

	// 最大音量
	MaxUIVolume float64

	cards CardInfos

	siEventChan  chan func()
	siPollerExit chan struct{}

	isSaving    bool
	saverLocker sync.Mutex

	sinkLocker sync.Mutex
	portLocker sync.Mutex
}

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

// SetPort activate the port for the special card.
// The available sinks and sources will also change with the profile changing.
func (a *Audio) SetPort(cardId uint32, portName string, direction int32) error {
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
	curCard, _ = a.cards.getByPort(oppositePort, oppositeDirection)

	var (
		destCard     *CardInfo
		destPortInfo pulse.CardPortInfo
		destProfile  string

		oppositePortInfo pulse.CardPortInfo
		commonProfiles   pulse.ProfileInfos2
	)
	if curCard != nil && curCard.Id == cardId {
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
		return a.trySetSinkByPort(portName)
	}
	return a.trySetSourceByPort(portName)
}

func NewSink(core *pulse.Sink) *Sink {
	s := &Sink{core: core}
	s.index = s.core.Index
	s.update()
	return s
}
func NewSource(core *pulse.Source) *Source {
	s := &Source{core: core}
	s.index = s.core.Index
	s.update()
	return s
}
func NewSinkInput(core *pulse.SinkInput) *SinkInput {
	s := &SinkInput{core: core}
	s.index = s.core.Index
	s.update()
	return s
}
func NewAudio(core *pulse.Context) *Audio {
	a := &Audio{core: core}
	a.MaxUIVolume = pulse.VolumeUIMax
	a.siEventChan = make(chan func(), 10)
	a.siPollerExit = make(chan struct{})
	a.applyConfig()
	a.update()
	a.initEventHandlers()

	go a.sinkInputPoller()

	return a
}

func (a *Audio) destroy() {
	close(a.siPollerExit)
	dbus.UnInstallObject(a)
}

func (a *Audio) SetDefaultSink(name string) {
	a.sinkLocker.Lock()
	defer a.sinkLocker.Unlock()

	if a.DefaultSink != nil && a.DefaultSink.Name == name {
		return
	}

	a.core.SetDefaultSink(name)
	a.update()
	a.saveConfig()

	var idxList []uint32
	for _, sinkInput := range a.SinkInputs {
		idxList = append(idxList, sinkInput.index)
	}
	if len(idxList) == 0 {
		return
	}
	a.core.MoveSinkInputsByName(idxList, name)
}
func (a *Audio) SetDefaultSource(name string) {
	if a.DefaultSource != nil && a.DefaultSource.Name == name {
		return
	}
	a.core.SetDefaultSource(name)
	a.update()
	a.saveConfig()
}

type Port struct {
	Name        string
	Description string
	Available   byte // Unknow:0, No:1, Yes:2
}
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
	s.core.SetMute(v)
	if !v {
		playFeedbackWithDevice(s.Name)
	}
}

// 设置此设备的当前使用端口
func (s *Sink) SetPort(name string) {
	s.core.SetPort(name)
}

type SinkInput struct {
	core  *pulse.SinkInput
	index uint32

	// process name
	Name string
	Icon string
	Mute bool

	Volume         float64
	Balance        float64
	SupportBalance bool
	Fade           float64
	SupportFade    bool
}

func (s *SinkInput) SetVolume(v float64, isPlay bool) error {
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
func (s *SinkInput) SetBalance(v float64, isPlay bool) error {
	if v < -1.00 || v > 1.00 {
		return fmt.Errorf("Invalid volume value: %v", v)
	}

	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
	if isPlay {
		playFeedback()
	}
	return nil
}
func (s *SinkInput) SetFade(v float64) error {
	if v < -1.00 || v > 1.00 {
		return fmt.Errorf("Invalid volume value: %v", v)
	}

	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
	playFeedback()
	return nil
}
func (s *SinkInput) SetMute(v bool) {
	s.core.SetMute(v)
	if !v {
		playFeedback()
	}
}

type Source struct {
	core  *pulse.Source
	index uint32

	Name        string
	Description string

	// 默认的输入音量
	BaseVolume float64

	Mute bool

	Volume         float64
	Balance        float64
	SupportBalance bool
	Fade           float64
	SupportFade    bool

	Ports      []Port
	ActivePort Port
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

type Daemon struct {
	*ModuleBase
}

func NewAudioDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = NewModuleBase("audio", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

var _audio *Audio

func finalize() {
	_audio.destroy()
	_audio = nil
	logger.EndTracing()
}

func (*Daemon) Start() error {
	if _audio != nil {
		return nil
	}

	logger.BeginTracing()

	ctx := pulse.GetContext()
	_audio = NewAudio(ctx)

	if err := dbus.InstallOnSession(_audio); err != nil {
		logger.Error("Failed InstallOnSession:", err)
		finalize()
		return err
	}

	initDefaultVolume(_audio)
	return nil
}

func (*Daemon) Stop() error {
	if _audio == nil {
		return nil
	}

	finalize()
	return nil
}

func playFeedback() {
	playFeedbackWithDevice("")
}

func playFeedbackWithDevice(device string) {
	soundutils.PlaySystemSound(soundutils.EventVolumeChanged, device, false)
}

func isVolumeValid(v float64) bool {
	if v < 0 || v > pulse.VolumeUIMax {
		return false
	}
	return true
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

func (a *Audio) trySetSinkByPort(portName string) error {
	for _, sink := range a.core.GetSinkList() {
		if !isPortExists(portName, sink.Ports) {
			continue
		}
		if sink.ActivePort.Name != portName {
			sink.SetPort(portName)
		}
		if a.DefaultSink == nil || a.DefaultSink.Name != sink.Name {
			a.SetDefaultSink(sink.Name)
		}
		return nil
	}
	return fmt.Errorf("Cann't find valid sink for port '%s'", portName)
}

func (a *Audio) trySetSourceByPort(portName string) error {
	for _, source := range a.core.GetSourceList() {
		if !isPortExists(portName, source.Ports) {
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

func (a *Audio) getDefaultSink(name string) *Sink {
	if a.DefaultSink != nil && a.DefaultSink.Name == name {
		return a.DefaultSink
	}
	for _, o := range a.core.GetSinkList() {
		if o.Name != name {
			continue
		}
		// TODO: Free old sink info
		a.DefaultSink = NewSink(o)
		break
	}
	return a.DefaultSink
}
func (a *Audio) getDefaultSource(name string) *Source {
	if a.DefaultSource != nil && a.DefaultSource.Name == name {
		return a.DefaultSource
	}
	for _, o := range a.core.GetSourceList() {
		if o.Name != name {
			continue
		}
		a.DefaultSource = NewSource(o)
		break
	}
	return a.DefaultSource
}

func isPortExists(name string, ports []pulse.PortInfo) bool {
	for _, port := range ports {
		if port.Name == name {
			return true
		}
	}
	return false
}
