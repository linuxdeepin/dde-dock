package audio

import (
	libsound "dbus/com/deepin/api/sound"
	. "pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/pulse"
)

type Audio struct {
	init bool
	core *pulse.Context

	// 输出设备 ObjectPath 列表
	Sinks []*Sink
	// 输入设备ObjectPath 列表
	Sources []*Source
	// 正常输出声音的程序列表
	SinkInputs []*SinkInput

	// 默认的输出设备名称
	DefaultSink string
	// 默认的输入设备名称
	DefaultSource string

	// 最大音量
	MaxUIVolume float64

	siEventChan  chan func()
	siPollerExit chan struct{}
}

func (a *Audio) Reset() {
	for _, s := range a.Sinks {
		s.SetVolume(s.BaseVolume, false)
		s.SetBalance(0, false)
		s.SetFade(0)
	}
	for _, s := range a.Sources {
		s.SetVolume(s.BaseVolume, false)
		s.SetBalance(0, false)
		s.SetFade(0)
	}
}

func (s *Audio) GetDefaultSink() *Sink {
	for _, o := range s.Sinks {
		if o.Name == s.DefaultSink {
			return o
		}
	}
	return nil
}
func (s *Audio) GetDefaultSource() *Source {
	for _, o := range s.Sources {
		if o.Name == s.DefaultSource {
			return o
		}
	}
	return nil
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
	a.update()
	a.initEventHandlers()

	a.setupMediakeyMonitor()
	go a.sinkInputPoller()

	return a
}

func (a *Audio) destroy() {
	close(a.siPollerExit)
	dbus.UnInstallObject(a)
}

func (a *Audio) SetDefaultSink(name string) {
	a.core.SetDefaultSink(name)
	a.update()
}
func (a *Audio) SetDefaultSource(name string) {
	a.core.SetDefaultSource(name)
	a.update()
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
func (s *Sink) SetVolume(v float64, isPlay bool) {
	if v == 0 {
		v = 0.001
	}
	s.core.SetVolume(s.core.Volume.SetAvg(v))
	if isPlay {
		playFeedbackWithDevice(s.Name)
	}
}

// 设置左右声道平衡值
//
// v: 声道平衡值
//
// isPlay: 是否播放声音反馈
func (s *Sink) SetBalance(v float64, isPlay bool) {
	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
	if isPlay {
		playFeedbackWithDevice(s.Name)
	}
}

// 设置前后声道平衡值
//
// v: 声道平衡值
//
// isPlay: 是否播放声音反馈
func (s *Sink) SetFade(v float64) {
	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
	playFeedbackWithDevice(s.Name)
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

func (s *SinkInput) SetVolume(v float64, isPlay bool) {
	if v == 0 {
		v = 0.001
	}
	s.core.SetVolume(s.core.Volume.SetAvg(v))
	if isPlay {
		playFeedback()
	}
}
func (s *SinkInput) SetBalance(v float64, isPlay bool) {
	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
	if isPlay {
		playFeedback()
	}
}
func (s *SinkInput) SetFade(v float64) {
	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
	playFeedback()
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
func (s *Source) SetVolume(v float64, isPlay bool) {
	if v == 0 {
		v = 0.001
	}
	s.core.SetVolume(s.core.Volume.SetAvg(v))
	if isPlay {
		playFeedback()
	}
}
func (s *Source) SetBalance(v float64, isPlay bool) {
	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
	if isPlay {
		playFeedback()
	}
}
func (s *Source) SetFade(v float64) {
	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
	playFeedback()
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
	return nil
}

func (*Daemon) Stop() error {
	if _audio == nil {
		return nil
	}

	finalize()
	return nil
}

var playFeedback = func() func() {
	player, err := libsound.NewSound("com.deepin.api.Sound", "/com/deepin/api/Sound")
	if err != nil {
		logger.Error("Can't create com.deepin.api.Sound! Sound feedback support will be disabled", err)
		return nil
	}

	return func() {
		player.PlaySystemSound("audio-volume-change")
	}
}()

var playFeedbackWithDevice = func() func(string) {
	player, err := libsound.NewSound("com.deepin.api.Sound", "/com/deepin/api/Sound")
	if err != nil {
		logger.Error("Can't create com.deepin.api.Sound! Sound feedback support will be disabled", err)
		return nil
	}

	return func(device string) {
		player.PlaySystemSoundWithDevice("audio-volume-change", device)
	}
}()
