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
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/pulse"
)

const (
	dbusServiceName = "com.deepin.daemon.Audio"
	dbusPath        = "/com/deepin/daemon/Audio"
	dbusInterface   = dbusServiceName
)

func (*Audio) GetInterfaceName() string {
	return dbusInterface
}

func filterSinkInput(c *pulse.SinkInput) bool {
	appName := c.PropList[pulse.PA_PROP_APPLICATION_NAME]
	if appName == "com.deepin.SoundEffect" {
		return true
	}

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
		si := NewSinkInput(s, a.service)
		if si == nil {
			continue
		}
		sinkinputs = append(sinkinputs, si)
	}

	for _, o := range a.sinkInputs {
		a.service.StopExport(o)
	}

	sinkInputPaths := make([]dbus.ObjectPath, len(sinkinputs))
	for idx, o := range sinkinputs {
		sinkInputPath := o.getPath()
		a.service.Export(sinkInputPath, o)
		sinkInputPaths[idx] = sinkInputPath
	}
	a.sinkInputs = sinkinputs
	a.PropsMu.Lock()
	a.setPropSinkInputs(sinkInputPaths)
	a.PropsMu.Unlock()
}

func (a *Audio) addSinkInput(idx uint32) {
	for _, si := range a.sinkInputs {
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

	si := NewSinkInput(core, a.service)
	if si == nil {
		return
	}
	sinkInputPath := si.getPath()
	err = a.service.Export(sinkInputPath, si)
	if err != nil {
		logger.Error(err)
		return
	}

	a.sinkInputs = append(a.sinkInputs, si)
	a.PropsMu.Lock()
	a.SinkInputs = append(a.SinkInputs, sinkInputPath)
	a.emitPropChangedSinkInputs(a.SinkInputs)
	a.PropsMu.Unlock()
	logger.Debugf("addSinkInput idx: %d, si: %#v", idx, si)
}

func (a *Audio) removeSinkInput(idx uint32) {
	var tryRemoveSinkInput *SinkInput
	var newSinkInputList []*SinkInput
	for _, si := range a.sinkInputs {
		if si.index == idx {
			tryRemoveSinkInput = si
		} else {
			newSinkInputList = append(newSinkInputList, si)
		}
	}

	if tryRemoveSinkInput != nil {
		logger.Debugf("removeSinkInput idx: %d, si: %#v", idx, tryRemoveSinkInput)
		a.service.StopExport(tryRemoveSinkInput)
		a.sinkInputs = newSinkInputList
		sinkInputPaths := make([]dbus.ObjectPath, len(newSinkInputList))
		for idx, si := range newSinkInputList {
			sinkInputPaths[idx] = si.getPath()
		}
		a.PropsMu.Lock()
		a.setPropSinkInputs(sinkInputPaths)
		a.PropsMu.Unlock()
	}
}

func (a *Audio) update() {
	logger.Debug("Audio.update")
	sinfo, _ := a.core.GetServer()
	if sinfo != nil {
		a.updateDefaultSink(sinfo.DefaultSinkName, true)
		a.updateDefaultSource(sinfo.DefaultSourceName, true)
	}

	a.rebuildSinkInputList()
	a.cards = newCardInfos(a.core.GetCardList())

	a.PropsMu.Lock()
	a.setPropCards(a.cards.string())
	a.PropsMu.Unlock()

	if a.defaultSink != nil {
		a.moveSinkInputsToSink(a.defaultSink.index)
	}
}

func (a *Audio) updateDefaultSink(name string, force bool) {
	if !force && a.defaultSink != nil && a.defaultSink.Name == name {
		// default source no changed
		a.defaultSink.update()
		return
	}
	// default sink changed
	for _, o := range a.core.GetSinkList() {
		if o.Name != name {
			continue
		}

		if a.defaultSink != nil {
			a.service.StopExport(a.defaultSink)
		}
		a.defaultSink = NewSink(o, a.service)
		defaultSinkPath := a.defaultSink.getPath()
		a.service.Export(defaultSinkPath, a.defaultSink)

		a.PropsMu.Lock()
		a.setPropDefaultSink(defaultSinkPath)
		a.PropsMu.Unlock()

		logger.Debugf("Audio.DefaultSink change to #%d %s",
			a.defaultSink.index, a.defaultSink.Name)
		return
	}
}

func (a *Audio) updateDefaultSource(name string, force bool) {
	if !force && a.defaultSource != nil && a.defaultSource.Name == name {
		// default source no changed
		a.defaultSource.update()
		return
	}
	// default source changed
	for _, o := range a.core.GetSourceList() {
		if o.Name != name {
			continue
		}

		if a.defaultSource != nil {
			a.service.StopExport(a.defaultSource)
		}
		a.defaultSource = NewSource(o, a.service)
		defaultSourcePath := a.defaultSource.getPath()
		a.service.Export(defaultSourcePath, a.defaultSource)

		a.PropsMu.Lock()
		a.setPropDefaultSource(defaultSourcePath)
		a.PropsMu.Unlock()

		logger.Debugf("Audio.DefaultSource change to #%d %s", a.defaultSource.index, a.defaultSource.Name)
		return
	}
}
