package main

import "dlib"
import "dlib/dbus"
import "dlib/logger"
import "dlib/pulse"
import "os"

var Logger = logger.NewLogger("com.deepin.daemon.Audio")

type Audio struct {
	init bool
	core *pulse.Context

	Sinks         []*Sink
	Sources       []*Source
	SinkInputs    []*SinkInput
	DefaultSink   string
	DefaultSource string

	MaxUIVolume float64
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
	s.update()
	return s
}
func NewSource(core *pulse.Source) *Source {
	s := &Source{core: core}
	s.update()
	return s
}
func NewSinkInput(core *pulse.SinkInput) *SinkInput {
	s := &SinkInput{core: core}
	s.update()
	return s
}
func NewAudio(core *pulse.Context) *Audio {
	a := &Audio{core: core}
	a.MaxUIVolume = pulse.VolumeUIMax
	a.update()
	a.initEventHandlers()
	return a
}

func (a *Audio) SetDefaultSink(name string) {
	a.core.SetDefaultSink(name)
}
func (a *Audio) SetDefaultSource(name string) {
	a.core.SetDefaultSource(name)
}

type Sink struct {
	core *pulse.Sink

	Name        string
	Description string

	BaseVolume float64

	Mute bool

	Volume         float64
	Balance        float64
	SupportBalance bool
	Fade           float64
	SupportFade    bool

	Ports      []string
	ActivePort string
}

func (s *Sink) SetVolume(v float64) {
	s.core.SetVolume(s.core.Volume.SetAvg(v))
}
func (s *Sink) SetBalance(v float64) {
	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
}
func (s *Sink) SetFade(v float64) {
	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
}
func (s *Sink) SetMute(v bool) {
	s.core.SetMute(v)
}
func (s *Sink) SetPort(name string) {
	s.core.SetPort(name)
}

type SinkInput struct {
	core *pulse.SinkInput
	Name string
	Icon string
	Mute bool

	Volume         float64
	Balance        float64
	SupportBalance bool
	Fade           float64
	SupportFade    bool
}

func (s *SinkInput) SetVolume(v float64) {
	s.core.SetVolume(s.core.Volume.SetAvg(v))
}
func (s *SinkInput) SetBalance(v float64) {
	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
}
func (s *SinkInput) SetFade(v float64) {
	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
}
func (s *SinkInput) SetMute(v bool) {
	s.core.SetMute(v)
}

type Source struct {
	core        *pulse.Source
	Name        string
	Description string

	BaseVolume float64

	Mute bool

	Volume         float64
	Balance        float64
	SupportBalance bool
	Fade           float64
	SupportFade    bool

	ActivePort string
	Ports      []string
}

func (s *Source) SetVolume(v float64) {
	s.core.SetVolume(s.core.Volume.SetAvg(v))
}
func (s *Source) SetBalance(v float64) {
	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
}
func (s *Source) SetFade(v float64) {
	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
}
func (s *Source) SetPort(name string) {
	s.core.SetPort(name)
}

func main() {
	defer Logger.EndTracing()
	if !dlib.UniqueOnSession("com.deepin.daemon.Audio") {
		Logger.Warning("There already has an Audio daemon running.")
		return
	}

	ctx := pulse.GetContext()
	audio := NewAudio(ctx)

	if err := dbus.InstallOnSession(audio); err != nil {
		Logger.Error("Failed InstallOnSession:", err)
		return
	}

	dbus.DealWithUnhandledMessage()
	audio.listenMediaKey()
	if err := dbus.Wait(); err != nil {
		Logger.Error("dbus.Wait recieve an error:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
	dbus.Wait()
}
