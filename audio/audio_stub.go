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
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/pulse"
)

const (
	baseBusName = "com.deepin.daemon.Audio"
	baseBusPath = "/com/deepin/daemon/Audio"
	baseBusIfc  = baseBusName
)

func (*Audio) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       baseBusName,
		ObjectPath: baseBusPath,
		Interface:  baseBusIfc,
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

func (a *Audio) update() {
	sinfo, _ := a.core.GetServer()
	if sinfo != nil {
		a.setPropDefaultSink(a.getDefaultSink(sinfo.DefaultSinkName))
		a.setPropDefaultSource(a.getDefaultSource(sinfo.DefaultSourceName))
	}

	a.rebuildSinkInputList()
	a.cards = newCardInfos(a.core.GetCardList())
	a.setPropCards(a.cards.string())
	a.setPropActiveSinkPort(a.getActiveSinkPort())
	a.setPropActiveSourcePort(a.getActiveSourcePort())
}

func (s *Audio) setPropDefaultSink(v *Sink) {
	if v == nil || toJSON(s.DefaultSink) != toJSON(v) {
		s.DefaultSink = v
		dbus.NotifyChange(s, "DefaultSink")
	}
}

func (s *Audio) setPropDefaultSource(v *Source) {
	if v == nil || toJSON(s.DefaultSource) != toJSON(v) {
		s.DefaultSource = v
		dbus.NotifyChange(s, "DefaultSource")
	}
}

func (a *Audio) setPropActiveSinkPort(port string) {
	if len(port) == 0 || a.ActiveSinkPort == port {
		return
	}
	a.ActiveSinkPort = port
	dbus.NotifyChange(a, "ActiveSinkPort")
}

func (a *Audio) setPropActiveSourcePort(port string) {
	if len(port) == 0 || a.ActiveSourcePort == port {
		return
	}
	a.ActiveSourcePort = port
	dbus.NotifyChange(a, "ActiveSourcePort")
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

func (a *Audio) setPropCards(v string) {
	if a.Cards == v {
		return
	}
	a.Cards = v
	dbus.NotifyChange(a, "Cards")
}
