package audio

import "pkg.linuxdeepin.com/lib/dbus"
import "pkg.linuxdeepin.com/lib/pulse"
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
		Dest:       baseBusName,
		ObjectPath: baseBusPath,
		Interface:  baseBusIfc,
	}
}

func (m *Meter) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       baseBusName,
		ObjectPath: baseBusPath + "/Meter" + m.id,
		Interface:  baseBusIfc + ".Meter",
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
		Dest:       baseBusName,
		ObjectPath: fmt.Sprintf("%s/Sink%d", baseBusPath, s.index),
		Interface:  baseBusIfc + ".Sink",
	}
}

func (s *Source) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       baseBusName,
		ObjectPath: fmt.Sprintf("%s/Source%d", baseBusPath, s.index),
		Interface:  baseBusIfc + ".Source",
	}
}

func (s *SinkInput) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       baseBusName,
		ObjectPath: fmt.Sprintf("%s/SinkInput%d", baseBusPath, s.index),
		Interface:  baseBusIfc + ".SinkInput",
	}
}

func filterSinkInput(c *pulse.SinkInput) bool {
	switch c.PropList[pulse.PA_PROP_MEDIA_ROLE] {
	case "video", "music", "game":
		return false
	case "animation", "production", "phone":
		//TODO: what's the meaning of this type? Should we filter this SinkInput?
		return false
	default:
		return false

	case "event", "a11y", "test":
		//Filter this SinkInput
		return true
	}
}

func (a *Audio) rebuildSinkInputList() {
	var sinkinputs []*SinkInput
	for _, s := range a.core.GetSinkInputList() {
		if s == nil || filterSinkInput(s) {
			continue
		}
		si := NewSinkInput(s)
		sinkinputs = append(sinkinputs, si)
	}
	a.setPropSinkInputs(sinkinputs)
}

func (a *Audio) addSinkInput(idx uint32) {
	for _, si := range a.SinkInputs {
		if si.index == idx {
			return
		}
	}

	core, err := a.core.GetSinkInput(idx)
	if err != nil {
		logger.Warning(err)
		return
	}
	if filterSinkInput(core) {
		return
	}

	si := NewSinkInput(core)
	err = dbus.InstallOnSession(si)
	if err != nil {
		logger.Error(err)
		return
	}

	a.SinkInputs = append(a.SinkInputs, si)
	dbus.NotifyChange(a, "SinkInputs")
}
func (a *Audio) removeSinkInput(idx uint32) {
	var tryRemoveSinkInput *SinkInput
	var newSinkInputList []*SinkInput
	for _, si := range a.SinkInputs {
		if si.index == idx {
			tryRemoveSinkInput = si
		} else {
			newSinkInputList = append(newSinkInputList, si)
		}
	}

	if tryRemoveSinkInput != nil {
		dbus.UnInstallObject(tryRemoveSinkInput)
		a.SinkInputs = newSinkInputList
		dbus.NotifyChange(a, "SinkInputs")
	}
}

func (a *Audio) rebuildSinkList() {
	var sinks []*Sink
	for _, s := range a.core.GetSinkList() {
		if s == nil {
			continue
		}
		sinks = append(sinks, NewSink(s))
	}
	a.setPropSinks(sinks)
}

func (a *Audio) rebuildSourceList() {
	var sources []*Source
	for _, s := range a.core.GetSourceList() {
		if s == nil {
			continue
		}
		obj := NewSource(s)
		if len(obj.Ports) > 0 {
			sources = append(sources, obj)
		}
	}
	a.setPropSources(sources)
}
func (a *Audio) update() {
	sinfo, _ := a.core.GetServer()
	if sinfo != nil {
		a.setPropDefaultSink(sinfo.DefaultSinkName)
		a.setPropDefaultSource(sinfo.DefaultSourceName)
	}

	a.rebuildSinkList()
	a.rebuildSourceList()
	a.rebuildSinkInputList()
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

	s.setPropActivePort(toPort(s.core.ActivePort))
	var ports []Port
	for _, p := range s.core.Ports {
		ports = append(ports, toPort(p))
	}
	s.setPropPorts(ports)
}
func toPort(v pulse.PortInfo) Port {
	return Port{
		Name:        v.Name,
		Description: v.Description,
		Available:   byte(v.Available),
	}
}

func (s *Sink) setPropActivePort(v Port) {
	if s.ActivePort != v {
		s.ActivePort = v
		dbus.NotifyChange(s, "ActivePort")
	}
}
func (s *Sink) setPropPorts(v []Port) {
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

	s.setPropActivePort(toPort(s.core.ActivePort))

	var ports []Port
	for _, p := range s.core.Ports {
		ports = append(ports, toPort(p))
	}
	s.setPropPorts(ports)
}
func (s *Source) setPropPorts(v []Port) {
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
func (s *Source) setPropActivePort(v Port) {
	if s.ActivePort != v {
		s.ActivePort = v
		dbus.NotifyChange(s, "ActivePort")
	}
}

func (s *SinkInput) update() {
	s.Name = s.core.PropList[PropAppName]
	s.Icon = s.core.PropList[PropAppIconName]

	// Correct app name and icon
	s.correctAppName()

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
