package audio

import "pkg.linuxdeepin.com/lib/pulse"

func (a *Audio) initEventHandlers() {
	if !a.init {
		a.core.Connect(pulse.FacilitySink, func(e int, idx uint32) {
			a.handleSinkEvent(e, idx)
		})
		a.core.Connect(pulse.FacilitySource, func(e int, idx uint32) {
			a.handleSourceEvent(e, idx)
		})
		a.core.Connect(pulse.FacilitySinkInput, func(e int, idx uint32) {
			a.handleSinkInputEvent(e, idx)
		})
		a.core.Connect(pulse.FacilityServer, func(e int, idx uint32) {
			a.handleServerEvent()
		})
		a.init = true
	}
}

func (a *Audio) handleSinkEvent(eType int, idx uint32) {
	switch eType {
	case pulse.EventTypeNew, pulse.EventTypeRemove:
		a.rebuildSinkList()

	case pulse.EventTypeChange:
		for _, s := range a.Sinks {
			if s.core.Index == idx {
				s.core = a.core.GetSink(idx)
				if s.core != nil {
					s.update()
				}
				break
			}
		}
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
		for _, s := range a.SinkInputs {
			if s.core.Index == idx {
				s.core = a.core.GetSinkInput(idx)
				if s.core != nil {
					s.update()
				}
				break
			}
		}
	}
}
func (a *Audio) handleSourceEvent(eType int, idx uint32) {
	switch eType {
	case pulse.EventTypeNew, pulse.EventTypeRemove:
		a.rebuildSourceList()

	case pulse.EventTypeChange:
		for _, s := range a.Sources {
			if s.core.Index == idx {
				s.core = a.core.GetSource(idx)
				if s.core != nil {
					s.update()
				}
				break
			}
		}
	}
}

func (a *Audio) handleServerEvent() {
	sinfo := a.core.GetServer()
	if sinfo != nil {
		a.setPropDefaultSink(sinfo.DefaultSinkName)
		a.setPropDefaultSource(sinfo.DefaultSourceName)
	}
}
