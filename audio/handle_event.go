/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import "pkg.deepin.io/lib/pulse"

func (a *Audio) initEventHandlers() {
	if !a.init {
		a.core.ConnectStateChanged(pulse.ContextStateFailed, func() {
			logger.Warning("Pulse context connection failed, try again")
			a.core = pulse.GetContextForced()
			a.update()
			a.init = false
			a.initEventHandlers()
		})

		a.core.Connect(pulse.FacilityCard, func(e int, idx uint32) {
			a.handleCardEvent(e, idx)
			a.setPropActiveSinkPort(a.getActiveSinkPort())
			a.setPropActiveSourcePort(a.getActiveSourcePort())
			a.saveConfig()
		})
		a.core.Connect(pulse.FacilitySink, func(e int, idx uint32) {
			a.handleSinkEvent(e, idx)
			a.saveConfig()
		})
		a.core.Connect(pulse.FacilitySource, func(e int, idx uint32) {
			a.handleSourceEvent(e, idx)
			a.saveConfig()
		})
		a.core.Connect(pulse.FacilitySinkInput, func(e int, idx uint32) {
			a.handleSinkInputEvent(e, idx)
		})
		a.core.Connect(pulse.FacilityServer, func(e int, idx uint32) {
			a.handleServerEvent()
			a.saveConfig()
		})
		a.init = true
	}
}

func (a *Audio) handleCardEvent(eType int, idx uint32) {
	switch eType {
	case pulse.EventTypeNew:
		logger.Debug("[Event] card added:", idx)
		card, err := a.core.GetCard(idx)
		if nil != err {
			logger.Warning("get card info failed: ", err)
			return
		}
		infos, added := a.cards.add(newCardInfo(card))
		if added {
			a.setPropCards(infos.string())
			a.cards = infos
		}
		selectNewCardProfile(card)
	case pulse.EventTypeRemove:
		logger.Debug("[Event] card removed:", idx)
		infos, deleted := a.cards.delete(idx)
		if deleted {
			a.setPropCards(infos.string())
			a.cards = infos
		}
	case pulse.EventTypeChange:
		logger.Debug("[Event] card changed:", idx)
		card, err := a.core.GetCard(idx)
		if nil != err {
			logger.Warning("get card info failed: ", err)
			return
		}
		info, _ := a.cards.get(idx)
		if info != nil {
			info.update(card)
			a.setPropCards(a.cards.string())
		}
	}
}
func (a *Audio) handleSinkEvent(eType int, idx uint32) {
	switch eType {
	case pulse.EventTypeNew, pulse.EventTypeRemove:
		logger.Debug("[Event] sink added:", (eType == pulse.EventTypeNew), idx)
		a.update()

	case pulse.EventTypeChange:
		logger.Debug("[Event] sink changed:", idx)
		for _, s := range a.Sinks {
			if s.index == idx {
				info, err := a.core.GetSink(idx)
				if err != nil {
					logger.Error(err)
					break
				}

				s.core = info
				s.update()
				break
			}
		}
		a.setPropActiveSinkPort(a.getActiveSinkPort())
	}
}

func (a *Audio) sinkInputPoller() {
	for {
		select {
		case handler, ok := <-a.siEventChan:
			if !ok {
				logger.Error("SinkInput event channel has been abnormally closed!")
				return
			}

			handler()
		case <-a.siPollerExit:
			return
		}
	}
}

func (a *Audio) handleSinkInputEvent(eType int, idx uint32) {
	switch eType {
	case pulse.EventTypeNew:
		a.siEventChan <- func() {
			a.addSinkInput(idx)
		}
	case pulse.EventTypeRemove:
		a.siEventChan <- func() {
			a.removeSinkInput(idx)
		}

	case pulse.EventTypeChange:
		a.siEventChan <- func() {
			for _, s := range a.SinkInputs {
				if s.index == idx {
					info, err := a.core.GetSinkInput(idx)
					if err != nil {
						logger.Warning(err)
						break
					}

					s.core = info
					s.update()
					break
				}
			}
		}
	}
}
func (a *Audio) handleSourceEvent(eType int, idx uint32) {
	switch eType {
	case pulse.EventTypeNew, pulse.EventTypeRemove:
		logger.Debug("[Event] source added:", (eType == pulse.EventTypeNew), idx)
		a.update()

	case pulse.EventTypeChange:
		logger.Debug("[Event] source changed:", idx)
		for _, s := range a.Sources {
			if s.index == idx {
				info, err := a.core.GetSource(idx)
				if err != nil {
					logger.Error(err)
					break
				}

				s.core = info
				s.update()
				break
			}
		}
		a.setPropActiveSourcePort(a.getActiveSourcePort())
	}
}

func (a *Audio) handleServerEvent() {
	sinfo, err := a.core.GetServer()
	if err != nil {
		logger.Error(err)
		return
	}

	a.setPropDefaultSink(sinfo.DefaultSinkName)
	a.setPropDefaultSource(sinfo.DefaultSourceName)
}
