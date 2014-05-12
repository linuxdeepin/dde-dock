package main

import "dlib/dbus"
import "fmt"

const (
	baseBusName = "com.deepin.daemon.Audio"
	baseBusPath = "/com/deepin/daemon/Audio"
	baseBusIfc  = "com.deepin.daemon.Audio"
)

const (
	PropAppIconName = "application.icon_name"
	PropAppName     = "application.name"
	PropAppPID      = "application.process.id"
)

func (*Audio) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		baseBusName,
		baseBusPath,
		baseBusIfc,
	}
}

func (m *Meter) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		baseBusName,
		baseBusPath + "/Meter" + m.id,
		baseBusIfc + ".Meter",
	}
}
func (m *Meter) setPropVolume(v float64) {
	if m.Volume != v {
		m.Volume = v
		dbus.NotifyChange(m, "Volume")
	}
}

func (s *Sink) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		baseBusName,
		fmt.Sprintf("%s/Sink%d", baseBusPath, s.core.Index),
		baseBusIfc + ".Sink",
	}
}

func (s *Source) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		baseBusName,
		fmt.Sprintf("%s/Source%d", baseBusPath, s.core.Index),
		baseBusIfc + ".Source",
	}
}

func (s *SinkInput) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		baseBusName,
		fmt.Sprintf("%s/SinkInput%d", baseBusPath, s.core.Index),
		baseBusIfc + ".SinkInput",
	}
}

func (a *Audio) update() {
	sinfo := a.core.GetServer()
	a.setPropDefaultSink(sinfo.DefaultSinkName)
	a.setPropDefaultSource(sinfo.DefaultSourceName)

	var sinks []*Sink
	var sources []*Source
	var sinkinputs []*SinkInput
	for _, s := range a.core.GetSinkList() {
		sinks = append(sinks, NewSink(s))
	}
	for _, s := range a.core.GetSourceList() {
		obj := NewSource(s)
		if len(obj.Ports) > 0 {
			sources = append(sources, obj)
		}
	}
	for _, s := range a.core.GetSinkInputList() {
		sinkinputs = append(sinkinputs, NewSinkInput(s))
	}

	a.setPropSinks(sinks)
	a.setPropSources(sources)
	a.setPropSinkInputs(sinkinputs)
}

func (s *Audio) setPropDefaultSink(v string) {
	if s.DefaultSink != v {
		s.DefaultSink = v
		dbus.NotifyChange(s, "DefaultSink")
	}
}
func (s *Audio) setPropDefaultSource(v string) {
	if s.DefaultSource != v {
		s.DefaultSource = v
		dbus.NotifyChange(s, "DefaultSource")
	}
}
func (s *Audio) setPropSinks(v []*Sink) {
	for _, o := range s.Sinks {
		dbus.UnInstallObject(o)
	}
	for _, o := range v {
		dbus.InstallOnSession(o)
	}
	s.Sinks = v
	dbus.NotifyChange(s, "Sinks")
}
func (s *Audio) setPropSources(v []*Source) {
	for _, o := range s.Sources {
		dbus.UnInstallObject(o)
	}
	for _, o := range v {
		dbus.InstallOnSession(o)
	}
	s.Sources = v
	dbus.NotifyChange(s, "Sources")
}
func (s *Audio) setPropSinkInputs(v []*SinkInput) {
	for _, o := range s.SinkInputs {
		dbus.UnInstallObject(o)
	}
	for _, o := range v {
		dbus.InstallOnSession(o)
	}
	s.SinkInputs = v
	dbus.NotifyChange(s, "SinkInputs")
}

func (s *Sink) setPropCanBalance(v bool) {
}
func (s *Sink) update() {
	s.Name = s.core.Name
	s.Description = s.core.Description

	s.BaseVolume = s.core.BaseVolume.ToPercent()

	s.setPropMute(s.core.Mute)
	s.setPropVolume(s.core.Volume.Avg())

	s.setPropSupportFade(false)
	s.setPropFade(s.core.Volume.Fade(s.core.ChannelMap))
	s.setPropSupportBalance(true)
	s.setPropBalance(s.core.Volume.Balance(s.core.ChannelMap))

	s.setPropActivePort(s.core.ActivePort.Name)
	var ports []string
	for _, p := range s.core.Ports {
		ports = append(ports, p.Name)
	}
	s.setPropPorts(ports)
}

func (s *Sink) setPropPorts(v []string) {
	s.Ports = v
	dbus.NotifyChange(s, "Ports")
}
func (s *Sink) setPropVolume(v float64) {
	if s.Volume != v {
		s.Volume = v
		dbus.NotifyChange(s, "Volume")
	}
}
func (s *Sink) setPropBalance(v float64) {
	if s.Volume != v {
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
func (s *Sink) setPropActivePort(v string) {
	if s.ActivePort != v {
		s.ActivePort = v
		dbus.NotifyChange(s, "ActivePort")
	}
}

func (s *Source) update() {
	s.Name = s.core.Name
	s.Description = s.core.Description

	s.BaseVolume = s.core.BaseVolume.ToPercent()

	s.setPropVolume(s.core.Volume.Avg())
	s.setPropMute(s.core.Mute)

	//TODO: handle this
	s.setPropSupportFade(false)
	s.setPropFade(s.core.Volume.Fade(s.core.ChannelMap))
	s.setPropSupportBalance(true)
	s.setPropBalance(s.core.Volume.Balance(s.core.ChannelMap))

	s.setPropActivePort(s.core.ActivePort.Name)

	var ports []string
	for _, p := range s.core.Ports {
		ports = append(ports, p.Name)
	}
	s.setPropPorts(ports)
}
func (s *Source) setPropPorts(v []string) {
	s.Ports = v
	dbus.NotifyChange(s, "Ports")
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
	if s.Volume != v {
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
func (s *Source) setPropActivePort(v string) {
	if s.ActivePort != v {
		s.ActivePort = v
		dbus.NotifyChange(s, "ActivePort")
	}
}

func (s *SinkInput) update() {
	s.Name = s.core.PropList[PropAppName]
	s.Icon = s.core.PropList[PropAppIconName]
	s.setPropVolume(s.core.Volume.Avg())
	s.setPropMute(s.core.Mute)

	s.setPropSupportFade(false)
	s.setPropFade(s.core.Volume.Fade(s.core.ChannelMap))
	s.setPropSupportBalance(true)
	s.setPropBalance(s.core.Volume.Balance(s.core.ChannelMap))
}

func (s *SinkInput) setPropVolume(v float64) {
	if s.Volume != v {
		s.Volume = v
		dbus.NotifyChange(s, "Volume")
	}
}
func (s *SinkInput) setPropMute(v bool) {
	if s.Mute != v {
		s.Mute = v
		dbus.NotifyChange(s, "Mute")
	}
}

func (s *SinkInput) setPropBalance(v float64) {
	if s.Volume != v {
		s.Balance = v
		dbus.NotifyChange(s, "Balance")
	}
}
func (s *SinkInput) setPropSupportBalance(v bool) {
	if s.SupportBalance != v {
		s.SupportBalance = v
		dbus.NotifyChange(s, "SupportBalance")
	}
}
func (s *SinkInput) setPropSupportFade(v bool) {
	if s.SupportFade != v {
		s.SupportFade = v
		dbus.NotifyChange(s, "SupportFade")
	}
}
func (s *SinkInput) setPropFade(v float64) {
	if s.Fade != v {
		s.Fade = v
		dbus.NotifyChange(s, "Fade")
	}
}
