package main

import "dlib/pulse"

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
		var sinks []*Sink
		for _, s := range a.core.GetSinkList() {
			sinks = append(sinks, NewSink(s))
		}
		a.setPropSinks(sinks)

	case pulse.EventTypeChange:
		for _, s := range a.Sinks {
			if s.core.Index == idx {
				s.core = a.core.GetSink(idx)
				s.update()
				break
			}
		}
	}
}
func (a *Audio) handleSinkInputEvent(eType int, idx uint32) {
	switch eType {
	case pulse.EventTypeNew, pulse.EventTypeRemove:
		var sinkinputs []*SinkInput
		for _, s := range a.core.GetSinkInputList() {
			sinkinputs = append(sinkinputs, NewSinkInput(s))
		}
		a.setPropSinkInputs(sinkinputs)

	case pulse.EventTypeChange:
		for _, s := range a.SinkInputs {
			if s.core.Index == idx {
				s.core = a.core.GetSinkInput(idx)
				s.update()
				break
			}
		}
	}
}
func (a *Audio) handleSourceEvent(eType int, idx uint32) {
	switch eType {
	case pulse.EventTypeNew, pulse.EventTypeRemove:
		var sources []*Source
		for _, s := range a.core.GetSourceList() {
			sources = append(sources, NewSource(s))
		}
		a.setPropSources(sources)

	case pulse.EventTypeChange:
		for _, s := range a.Sources {
			if s.core.Index == idx {
				s.core = a.core.GetSource(idx)
				s.update()
				break
			}
		}
	}
}

func (a *Audio) handleServerEvent() {
	sinfo := a.core.GetServer()
	a.setPropDefaultSink(sinfo.DefaultSinkName)
	a.setPropDefaultSource(sinfo.DefaultSourceName)
}
