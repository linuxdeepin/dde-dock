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

func (*SourceOutputTest) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		baseBusName,
		baseBusPath + "/SourceOutputTest",
		baseBusIfc + ".SoureOutputTest",
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
	s.Sinks = v
	dbus.NotifyChange(s, "Sinks")
}
func (s *Audio) setPropSources(v []*Source) {
	for _, o := range s.Sources {
		dbus.UnInstallObject(o)
	}
	s.Sources = v
	dbus.NotifyChange(s, "Sources")
}
func (s *Audio) setPropSinkInputs(v []*SinkInput) {
	for _, o := range s.SinkInputs {
		dbus.UnInstallObject(o)
	}
	s.SinkInputs = v
	dbus.NotifyChange(s, "SinkInputs")
}

func (s *Sink) update() {
	s.Name = s.core.Name
	s.Description = s.core.Description

	//s.BaseVolume = s.core.BaseVolume
	s.setPropVolume(s.core.Volume.Avg())
	s.setPropBalance(s.core.Volume.Balance(s.core.ChannelMap))
	s.setPropMute(s.core.Mute)

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
		s.Volume = v
		dbus.NotifyChange(s, "Balance")
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

	//s.BaseVolume = s.core.BaseVolume
	s.setPropVolume(s.core.Volume.Avg())
	s.setPropBalance(s.core.Volume.Balance(s.core.ChannelMap))
	s.setPropMute(s.core.Mute)

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
func (s *Source) setPropBalance(v float64) {
	if s.Volume != v {
		s.Volume = v
		dbus.NotifyChange(s, "Balance")
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

func (s *SinkInput) update() {
	s.Name = s.core.Name
	s.Icon = s.core.PropList[PropAppIconName]
	s.setPropVolume(s.core.GetAvgVolume())
	s.setPropMute(s.core.Mute)
}
