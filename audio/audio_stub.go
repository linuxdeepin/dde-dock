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
	logger.Debug("rebuildSinkInputList")
	var sinkinputs []*SinkInput
	for _, s := range a.core.GetSinkInputList() {
		if s == nil || filterSinkInput(s) {
			continue
		}
		si := NewSinkInput(s)
		if si == nil {
			continue
		}
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
	if si == nil {
		return
	}
	err = dbus.InstallOnSession(si)
	if err != nil {
		logger.Error(err)
		return
	}

	a.SinkInputs = append(a.SinkInputs, si)
	dbus.NotifyChange(a, "SinkInputs")
	logger.Debugf("addSinkInput idx: %d, si: %#v", idx, si)
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
		logger.Debugf("removeSinkInput idx: %d, si: %#v", idx, tryRemoveSinkInput)
		dbus.UnInstallObject(tryRemoveSinkInput)
		a.SinkInputs = newSinkInputList
		dbus.NotifyChange(a, "SinkInputs")
	}
}

func (a *Audio) update() {
	logger.Debug("Audio.update")
	sinfo, _ := a.core.GetServer()
	if sinfo != nil {
		a.updateDefaultSink(sinfo.DefaultSinkName)
		a.updateDefaultSource(sinfo.DefaultSourceName)
	}

	a.rebuildSinkInputList()
	a.cards = newCardInfos(a.core.GetCardList())
	a.setPropCards(a.cards.string())
	if a.DefaultSink != nil {
		a.moveSinkInputsToSink(a.DefaultSink.index)
	}
}

func (a *Audio) updateDefaultSink(name string) {
	if a.DefaultSink != nil && a.DefaultSink.Name == name {
		// default source no changed
		return
	}
	// default sink changed
	for _, o := range a.core.GetSinkList() {
		if o.Name != name {
			continue
		}

		if a.DefaultSink != nil {
			dbus.UnInstallObject(a.DefaultSink)
		}
		a.DefaultSink = NewSink(o)
		dbus.InstallOnSession(a.DefaultSink)
		dbus.NotifyChange(a, "DefaultSink")
		logger.Debugf("Audio.DefaultSink change to #%d %s", a.DefaultSink.index, a.DefaultSink.Name)
		return
	}
}

func (a *Audio) updateDefaultSource(name string) {
	if a.DefaultSource != nil && a.DefaultSource.Name == name {
		// default source no changed
		return
	}
	// default source changed
	for _, o := range a.core.GetSourceList() {
		if o.Name != name {
			continue
		}

		if a.DefaultSource != nil {
			dbus.UnInstallObject(a.DefaultSource)
		}
		a.DefaultSource = NewSource(o)
		dbus.InstallOnSession(a.DefaultSource)
		dbus.NotifyChange(a, "DefaultSource")
		logger.Debugf("Audio.DefaultSource change to #%d %s", a.DefaultSource.index, a.DefaultSource.Name)
		return
	}
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
