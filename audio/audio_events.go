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
	"sort"
	"strconv"
	"strings"
	"time"

	dbus "github.com/godbus/dbus"
	bluez "github.com/linuxdeepin/go-dbus-factory/org.bluez"
	"pkg.deepin.io/lib/gsettings"
	"pkg.deepin.io/lib/pulse"
)

func (a *Audio) handleEvent() {
	for {
		select {
		case event := <-a.eventChan:
			switch event.Facility {
			case pulse.FacilityServer:
				a.handleServerEvent(event.Type)
				a.saveConfig()
			case pulse.FacilityCard:
				a.handleCardEvent(event.Type, event.Index)
				a.saveConfig()
			case pulse.FacilitySink:
				a.handleSinkEvent(event.Type, event.Index)
				a.saveConfig()
			case pulse.FacilitySource:
				a.handleSourceEvent(event.Type, event.Index)
				a.saveConfig()
			case pulse.FacilitySinkInput:
				a.handleSinkInputEvent(event.Type, event.Index)
			}

		case <-a.quit:
			logger.Debug("handleEvent return")
			return
		}
	}
}

func (a *Audio) handleStateChanged() {
	for {
		select {
		case state := <-a.stateChan:
			switch state {
			case pulse.ContextStateFailed:
				logger.Warning("pulseaudio context state failed")
				a.destroyCtxRelated()

				if !a.noRestartPulseAudio {
					logger.Debug("retry init")
					err := a.init()
					if err != nil {
						logger.Warning("failed to init:", err)
					}
					return
				} else {
					logger.Debug("do not restart pulseaudio")
				}
			}

		case <-a.quit:
			logger.Debug("handleStateChanged return")
			return
		}
	}
}

func isDeviceValid(deviceName string) bool {
	if strings.Contains(deviceName, "bluez") {
		systemBus, err := dbus.SystemBus()
		if err != nil {
			logger.Warning("[isDeviceValid] dbus connect failed:", err)
			return false
		}
		nameArray := strings.Split(deviceName, ".")
		if len(nameArray) < 2 {
			return false
		}
		var path string = "/org/bluez/hci0/dev_" + nameArray[1]
		bluezDevice, err := bluez.NewDevice(systemBus, dbus.ObjectPath(path))
		if err != nil {
			logger.Warning("[isDeviceValid] new device failed:", err)
			return false
		}
		icon, err := bluezDevice.Icon().Get(0)
		if err != nil {
			logger.Warning("[isDeviceValid] get icon failed:", err)
			return false
		}
		if icon == "computer" {
			return false
		}
		return true
	} else {
		return true
	}
}

func (a *Audio) isCardIdValid(cardId uint32) bool {
	for _, card := range a.cards {
		if card.Id == cardId {
			return true
		}
	}
	return false
}

func (a *Audio) needAutoSwitchInputPort() bool {
	if a.defaultSource == nil || !a.isCardIdValid(a.defaultSource.Card) {
		return true
	}

	cardName, portName := priorities.GetFirstInput()
	currentCardName := a.getCardNameById(a.defaultSource.Card)
	currentPortName := a.defaultSource.ActivePort.Name

	if cardName == currentCardName && portName == currentPortName {
		return false
	}

	firstPortType := GetPortType(cardName, portName)
	currentPortType := GetPortType(currentCardName, currentPortName)
	return priorities.IsInputTypeAfter(firstPortType, currentPortType)
}

func (a *Audio) needAutoSwitchOutputPort() bool {
	if a.defaultSink == nil || !a.isCardIdValid(a.defaultSink.Card) {
		return true
	}

	cardName, portName := priorities.GetFirstOutput()
	currentCardName := a.getCardNameById(a.defaultSink.Card)
	currentPortName := a.defaultSink.ActivePort.Name

	if cardName == currentCardName && portName == currentPortName {
		return false
	}

	firstPortType := GetPortType(cardName, portName)
	currentPortType := GetPortType(currentCardName, currentPortName)
	return priorities.IsOutputTypeAfter(firstPortType, currentPortType)
}

func (a *Audio) autoSwitchPort() {
	if a.needAutoSwitchInputPort() {
		cardName, portName := priorities.GetFirstInput()
		if cardName != "" && portName != "" {
			logger.Debugf("input port auto switch to %s %s", cardName, portName)
			card, err := a.cards.getByName(cardName)
			if err == nil {
				err = a.setPort(card.Id, portName, pulse.DirectionSource)
			}
			if err != nil {
				logger.Warning(err)
			}
		} else {
			logger.Debugf("no input port")
		}
	}

	if a.needAutoSwitchOutputPort() {
		cardName, portName := priorities.GetFirstOutput()
		if cardName != "" && portName != "" {
			logger.Debugf("output port auto switch to %s %s", cardName, portName)
			card, err := a.cards.getByName(cardName)
			if err == nil {
				err = a.setPort(card.Id, portName, pulse.DirectionSink)
			}
			if err != nil {
				logger.Warning(err)
			}
		} else {
			logger.Debugf("no output port")
		}
	}
}

func (a *Audio) handleCardEvent(eventType int, idx uint32) {
	switch eventType {
	case pulse.EventTypeNew:
		cardInfo, err := a.ctx.GetCard(idx)
		if nil != err {
			logger.Warning("get card info failed: ", err)
			return
		}
		if !isDeviceValid(cardInfo.Name) {
			return
		}
		logger.Debugf("[Event] card #%d added %s", idx, cardInfo.Name)
		cards, added := a.cards.add(newCard(cardInfo))
		if added {
			a.PropsMu.Lock()
			a.setPropCards(cards.string())
			a.setPropCardsWithoutUnavailable(cards.stringWithoutUnavailable())
			a.PropsMu.Unlock()
			a.cards = cards

			card, err := cards.get(idx)
			if err == nil {
				for _, port := range cardInfo.Ports {
					if port.Available == pulse.AvailableTypeNo {
						logger.Warningf("port(%s %s) available is no", card.Name, port.Name)
						continue
					}

					if port.Available == pulse.AvailableTypeUnknow {
						logger.Warningf("port(%s %s) available is unknown", card.Name, port.Name)
					}

					if port.Direction == pulse.DirectionSink {
						priorities.AddOutputPort(card.core.Name, port.Name)
					} else {
						priorities.AddInputPort(card.core.Name, port.Name)
					}
				}
				err = priorities.Save(globalPrioritiesFilePath)
				priorities.Print()
			}
			if err != nil {
				logger.Warning(err)
			}
		}
		// fix change profile not work
		time.AfterFunc(time.Millisecond*500, func() {
			selectNewCardProfile(cardInfo)
			logger.Debug("After select profile:", cardInfo.ActiveProfile.Name)
			if cardInfo.ActiveProfile.Name == "a2dp_sink" {
				a.disableBluezSourceIfProfileIsA2dp()
			}
			a.autoSwitchPort()
		})
	case pulse.EventTypeRemove:
		cardInfo, err := a.cards.get(idx)
		if nil != err {
			logger.Warning("get card info failed: ", err)
			return
		}
		cards, deleted := a.cards.delete(idx)
		logger.Debugf("[Event] card #%d removed %s", idx, cards.string())
		if deleted {
			a.PropsMu.Lock()
			a.setPropCards(cards.string())
			a.setPropCardsWithoutUnavailable(cards.stringWithoutUnavailable())
			a.PropsMu.Unlock()
			a.cards = cards
			priorities.RemoveCard(cardInfo.core.Name)
			err := priorities.Save(globalPrioritiesFilePath)
			priorities.Print()
			if err != nil {
				logger.Warning(err)
			}
		}
		a.autoSwitchPort()
	case pulse.EventTypeChange:
		cardInfo, err := a.ctx.GetCard(idx)
		if nil != err {
			logger.Warning("get card info failed: ", err)
			return
		}
		if !isDeviceValid(cardInfo.Name) {
			return
		}
		logger.Debugf("[Event] card #%d changed %s", idx, cardInfo.Name)
		a.mu.Lock()
		card, _ := a.cards.get(idx)
		if card != nil {
			card.update(cardInfo)
			a.PropsMu.Lock()
			a.setPropCards(a.cards.string())
			a.setPropCardsWithoutUnavailable(a.cards.stringWithoutUnavailable())
			a.PropsMu.Unlock()
		}
		//如果声卡配置文件是a2dp时,是不允许添加输入设备的
		if cardInfo.ActiveProfile.Name == "a2dp_sink" {
			a.enableSource = false
		}
		a.mu.Unlock()
		priorities.RemoveUnavailable(a.cards)
		priorities.AddAvailable(a.cards)
		priorities.Print()
		a.autoSwitchPort()
	}
}

func (a *Audio) addSink(sinkInfo *pulse.Sink) {
	sink := newSink(sinkInfo, a)

	a.mu.Lock()
	a.sinks[sinkInfo.Index] = sink
	a.mu.Unlock()

	sinkPath := sink.getPath()
	err := a.service.Export(sinkPath, sink)
	if err != nil {
		logger.Warningf("failed to export sink #%d: %v", sink.index, err)
		return
	}
	a.updatePropSinks()

	if sink.Name == a.defaultSinkName {
		a.defaultSink = sink
		a.PropsMu.Lock()
		a.setPropDefaultSink(sinkPath)
		a.PropsMu.Unlock()
		logger.Debug("set prop default sink:", sinkPath)
	}
}

func (a *Audio) handleSinkEvent(eventType int, idx uint32) {
	switch eventType {
	case pulse.EventTypeNew:
		sinkInfo, err := a.ctx.GetSink(idx)
		if err != nil {
			logger.Warning(err)
			return
		}
		logger.Debugf("[Event] sink #%d added %s", idx, sinkInfo.Name)
		if !isPhysicalDevice(sinkInfo.Name) {
			return
		}
		if !isDeviceValid(sinkInfo.Name) {
			return
		}
		a.mu.Lock()
		_, ok := a.sinks[idx]
		a.mu.Unlock()
		if ok {
			return
		}
		a.addSink(sinkInfo)

	case pulse.EventTypeRemove:
		a.mu.Lock()
		sink, ok := a.sinks[idx]
		if !ok {
			a.mu.Unlock()
			return
		}
		logger.Debugf("[Event] sink #%d removed %s", idx, sink.Name)
		delete(a.sinks, idx)
		a.mu.Unlock()
		a.updatePropSinks()

		err := a.service.StopExport(sink)
		if err != nil {
			logger.Warning(err)
		}

	case pulse.EventTypeChange:
		sinkInfo, err := a.ctx.GetSink(idx)
		if err != nil {
			logger.Warning(err)
			return
		}
		logger.Debugf("[Event] sink #%d changed %s", idx, sinkInfo.Name)
		if !isPhysicalDevice(sinkInfo.Name) {
			return
		}
		if !isDeviceValid(sinkInfo.Name) {
			return
		}
		a.mu.Lock()
		sink, ok := a.sinks[idx]
		a.mu.Unlock()
		if !ok {
			a.addSink(sinkInfo)
			return
		}
		sink.update(sinkInfo)
	}
}

func (a *Audio) handleSinkInputEvent(eType int, idx uint32) {
	switch eType {
	case pulse.EventTypeNew:
		logger.Debugf("[Event] sink-input #%d added", idx)
		a.handleSinkInputAdded(idx)
	case pulse.EventTypeRemove:
		logger.Debugf("[Event] sink-input #%d removed", idx)
		a.handleSinkInputRemoved(idx)
	case pulse.EventTypeChange:
		sinkInputInfo, err := a.ctx.GetSinkInput(idx)
		if err != nil {
			logger.Warning(err)
			return
		}
		logger.Debugf("[Event] sink-input #%d changed %s", idx, sinkInputInfo.Name)
		a.mu.Lock()
		sinkInput, ok := a.sinkInputs[idx]
		a.mu.Unlock()
		if !ok {
			return
		}
		sinkInput.update(sinkInputInfo)
	}
}

func (a *Audio) updateObjPathsProp(type0 string, ids []int, setFn func(value []dbus.ObjectPath) bool) {
	sort.Ints(ids)
	paths := make([]dbus.ObjectPath, len(ids))
	for idx, id := range ids {
		paths[idx] = dbus.ObjectPath(dbusPath + "/" + type0 + strconv.Itoa(id))
	}
	a.PropsMu.Lock()
	setFn(paths)
	a.PropsMu.Unlock()
}

func (a *Audio) updatePropSinks() {
	var ids []int
	a.mu.Lock()
	for _, sink := range a.sinks {
		ids = append(ids, int(sink.index))
	}
	a.mu.Unlock()
	a.updateObjPathsProp("Sink", ids, a.setPropSinks)
}

func (a *Audio) updatePropSources() {
	var ids []int
	a.mu.Lock()
	for _, source := range a.sources {
		ids = append(ids, int(source.index))
	}
	a.mu.Unlock()
	a.updateObjPathsProp("Source", ids, a.setPropSources)
}

func (a *Audio) updatePropSinkInputs() {
	var ids []int
	a.mu.Lock()
	for _, sinkInput := range a.sinkInputs {
		if sinkInput.visible {
			ids = append(ids, int(sinkInput.index))
		}
	}
	a.mu.Unlock()
	a.updateObjPathsProp("SinkInput", ids, a.setPropSinkInputs)
}

func (a *Audio) addSinkInput(sinkInputInfo *pulse.SinkInput) {
	sinkInput := newSinkInput(sinkInputInfo, a)
	a.mu.Lock()
	a.sinkInputs[sinkInputInfo.Index] = sinkInput
	a.mu.Unlock()

	sinkInputPath := sinkInput.getPath()

	if sinkInput.visible {
		err := a.service.Export(sinkInputPath, sinkInput)
		if err != nil {
			logger.Warning(err)
			return
		}
	}
	a.updatePropSinkInputs()

	logger.Debugf("sink-input (#%d) %s play with sink #%d", sinkInputInfo.Index,
		sinkInputInfo.Name, sinkInputInfo.Sink)
}

func (a *Audio) handleSinkInputAdded(idx uint32) {
	sinkInputInfo, err := a.ctx.GetSinkInput(idx)
	if err != nil {
		logger.Warning(err)
		return
	}
	logger.Debugf("[Event] sink-input #%d added %s", idx, sinkInputInfo.Name)
	a.mu.Lock()
	_, ok := a.sinkInputs[idx]
	a.mu.Unlock()
	if ok {
		return
	}

	a.addSinkInput(sinkInputInfo)
}

func (a *Audio) handleSinkInputRemoved(idx uint32) {
	a.mu.Lock()
	sinkInput, ok := a.sinkInputs[idx]
	if !ok {
		a.mu.Unlock()
		return
	}
	logger.Debugf("[Event] sink-input #%d removed %s", idx, sinkInput.Name)
	delete(a.sinkInputs, idx)
	a.mu.Unlock()

	if sinkInput.visible {
		err := a.service.StopExport(sinkInput)
		if err != nil {
			logger.Warning(err)
		}
	}

	a.updatePropSinkInputs()
}

func (a *Audio) addSource(sourceInfo *pulse.Source) {
	//如果不能启用输入源,说明声卡配置文件是"a2dp",此时不能添加a2dp输入设备
	if !a.enableSource {
		a.enableSource = true
		return
	}
	source := newSource(sourceInfo, a)

	a.mu.Lock()
	a.sources[sourceInfo.Index] = source
	a.mu.Unlock()

	sourcePath := source.getPath()
	err := a.service.Export(sourcePath, source)
	if err != nil {
		logger.Warning(err)
		return
	}

	a.updatePropSources()

	if a.defaultSourceName == source.Name {
		a.defaultSource = source
		a.PropsMu.Lock()
		a.setPropDefaultSource(sourcePath)
		a.PropsMu.Unlock()
	}
}

func (a *Audio) handleSourceEvent(eventType int, idx uint32) {
	switch eventType {
	case pulse.EventTypeNew:
		sourceInfo, err := a.ctx.GetSource(idx)
		if err != nil {
			logger.Warning(err)
			return
		}
		logger.Debugf("[Event] source #%d added %s", idx, sourceInfo.Name)
		if !isPhysicalDevice(sourceInfo.Name) {
			return
		}
		if !isDeviceValid(sourceInfo.Name) {
			return
		}
		a.mu.Lock()
		_, ok := a.sources[idx]
		a.mu.Unlock()
		if ok {
			return
		}
		a.sourceIdx = idx
		a.addSource(sourceInfo)

		_, portConfig := configKeeper.GetCardAndPortConfig(a.getCardNameById(sourceInfo.Card), sourceInfo.ActivePort.Name)
		err = a.setReduceNoise(portConfig.ReduceNoise)
		if err != nil {
			logger.Debug("reduce physical device noise failed:", err)
		}
		a.ReduceNoise.Set(portConfig.ReduceNoise)
	case pulse.EventTypeRemove:
		a.mu.Lock()
		source, ok := a.sources[idx]
		if !ok {
			a.mu.Unlock()
			return
		}
		logger.Debugf("[Event] source #%d removed %s", idx, source.Name)
		delete(a.sources, idx)
		a.mu.Unlock()
		a.updatePropSources()

		err := a.service.StopExport(source)
		if err != nil {
			logger.Warning(err)
			return
		}
		// 移除物理设备需要关闭虚拟通道，后面切换
		if isPhysicalDevice(source.Name) {
			err = a.setReduceNoise(false)
			if err != nil {
				logger.Warning("set reduce noise fail:", err)
			}
			a.ReduceNoise.Set(false)
		}
	case pulse.EventTypeChange:
		sourceInfo, err := a.ctx.GetSource(idx)
		if err != nil {
			logger.Warning(err)
			return
		}
		logger.Debugf("[Event] source #%d changed %s", idx, sourceInfo.Name)
		if !isPhysicalDevice(sourceInfo.Name) {
			return
		}
		if !isDeviceValid(sourceInfo.Name) {
			return
		}
		a.mu.Lock()
		source, ok := a.sources[idx]
		a.mu.Unlock()
		if !ok {
			// not found source
			a.addSource(sourceInfo)
			return
		}
		source.update(sourceInfo)
	}
}

func isPhysicalDevice(deviceName string) bool {
	for _, virtualDeviceKey := range []string{
		"echoCancelSource", "echo-cancel", "Echo-Cancel", // virtual key
	} {
		if strings.Contains(deviceName, virtualDeviceKey) {
			return false
		}
	}
	return true
}

func (a *Audio) handleServerEvent(eventType int) {
	switch eventType {
	case pulse.EventTypeChange:
		server, err := a.ctx.GetServer()
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Debugf("[Event] server changed: default sink: %s, default source: %s",
			server.DefaultSinkName, server.DefaultSourceName)

		a.defaultSinkName = server.DefaultSinkName
		a.defaultSourceName = server.DefaultSourceName
		a.updateDefaultSink(server.DefaultSinkName)
		a.updateDefaultSource(server.DefaultSourceName)
	}
}

func (a *Audio) listenGSettingVolumeIncreaseChanged() {
	gsettings.ConnectChanged(gsSchemaAudio, gsKeyVolumeIncrease, func(val string) {
		volInc := a.settings.GetBoolean(gsKeyVolumeIncrease)
		if volInc {
			a.MaxUIVolume = increaseMaxVolume
		} else {
			a.MaxUIVolume = normalMaxVolume
		}
		gMaxUIVolume = a.MaxUIVolume
		err := a.emitPropChangedMaxUIVolume(a.MaxUIVolume)
		if err != nil {
			logger.Warning("changed Max UI Volume failed: ", err)
		} else {
			sink := a.defaultSink
			configKeeper.SetIncreaseVolume(a.getCardNameById(sink.Card), sink.ActivePort.Name, volInc)
			err = configKeeper.Save(configKeeperFile)
			if err != nil {
				logger.Warning(err)
			}
		}
	})
}

func (a *Audio) listenGSettingReduceNoiseChanged() {
	gsettings.ConnectChanged(gsSchemaAudio, gsKeyReduceNoise, func(val string) {
		reduce := a.ReduceNoise.Get()
		err := a.setReduceNoise(reduce)
		if err != nil {
			logger.Warning("set Reduce Noise failed: ", err)
		} else {
			source := a.defaultSource
			configKeeper.SetReduceNoise(a.getCardNameById(source.Card), source.ActivePort.Name, reduce)
			err = configKeeper.Save(configKeeperFile)
			if err != nil {
				logger.Warning(err)
			}
		}
	})
}
