/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import (
	"pkg.deepin.io/lib/pulse"
	"time"
)

func (a *Audio) applyConfig() {
	info, err := readConfigInfo()
	if err != nil {
		logger.Warning("Read config info failed:", err)
		return
	}

	if !a.isConfigValid(info) {
		logger.Warning("Invalid config:", info.string())
		return
	}

	for _, card := range a.core.GetCardList() {
		v, ok := info.Profiles[card.Name]
		if !ok {
			continue
		}

		if card.ActiveProfile.Name != v {
			card.SetProfile(v)
			time.Sleep(time.Microsecond * 300)
		}
	}

	var sinkValidity = true
	for _, s := range a.core.GetSinkList() {
		if s.Name == info.Sink {
			if len(info.SinkPort) == 0 {
				sinkValidity = false
				break
			}
			port := pulse.PortInfos(s.Ports).Get(info.SinkPort)
			// if port invalid, nothing to do.
			// TODO: some device port can play sound when state is 'NO', how to fix?
			if port == nil || port.Available == pulse.AvailableTypeNo {
				sinkValidity = false
				break
			}

			if s.ActivePort.Name != info.SinkPort {
				s.SetPort(info.SinkPort)
				time.Sleep(time.Microsecond * 50)
			}
			s.SetVolume(s.Volume.SetAvg(info.SinkVolume))
			time.Sleep(time.Microsecond * 50)
			break
		}
	}
	if sinkValidity {
		a.core.SetDefaultSink(info.Sink)
		time.Sleep(time.Microsecond * 50)
	}

	var sourceValidity = true
	for _, s := range a.core.GetSourceList() {
		if s.Name == info.Source {
			if len(info.SourcePort) == 0 {
				sourceValidity = false
				continue
			}
			port := pulse.PortInfos(s.Ports).Get(info.SourcePort)
			if port == nil || port.Available == pulse.AvailableTypeNo {
				sourceValidity = false
				continue
			}
			if s.ActivePort.Name != info.SourcePort {
				s.SetPort(info.SourcePort)
				time.Sleep(time.Microsecond * 50)
			}
			s.SetVolume(s.Volume.SetAvg(info.SourceVolume))
			time.Sleep(time.Microsecond * 50)
			break
		}
	}
	if sourceValidity {
		a.core.SetDefaultSource(info.Source)
		time.Sleep(time.Microsecond * 50)
	}
}

func (a *Audio) saveConfig() {
	a.saverLocker.Lock()
	if a.isSaving {
		a.saverLocker.Unlock()
		return
	}

	a.isSaving = true
	a.saverLocker.Unlock()

	time.AfterFunc(time.Second*1, func() {
		a.doSaveConfig()

		a.saverLocker.Lock()
		a.isSaving = false
		a.saverLocker.Unlock()
	})
}

func (a *Audio) doSaveConfig() {
	var info = configInfo{
		Profiles: make(map[string]string),
	}
	for _, card := range a.core.GetCardList() {
		info.Profiles[card.Name] = card.ActiveProfile.Name
	}

	for _, s := range a.core.GetSinkList() {
		if a.DefaultSink == nil || s.Name != a.DefaultSink.Name {
			continue
		}
		info.Sink = s.Name
		info.SinkPort = s.ActivePort.Name
		info.SinkVolume = s.Volume.Avg()
		break
	}

	for _, s := range a.core.GetSourceList() {
		if a.DefaultSource == nil || s.Name != a.DefaultSource.Name {
			continue
		}
		info.Source = s.Name
		info.SourcePort = s.ActivePort.Name
		info.SourceVolume = s.Volume.Avg()
		break
	}

	err := saveConfigInfo(&info)
	if err != nil {
		logger.Warning("Save config file failed:", info.string(), err)
	}
}

func (a *Audio) isConfigValid(cfg *configInfo) bool {
	// check cfg.Profiles
	var validProfileCount int
	for _, card := range a.core.GetCardList() {
		cardProfile, ok := cfg.Profiles[card.Name]
		if !ok {
			continue
		}
		// find cardProfile in card.Profiles
		var found bool
		for _, profile := range card.Profiles {
			if profile.Name == cardProfile {
				found = true
				break
			}
		}

		if found {
			validProfileCount++
		} else {
			// cardProfile is invalid
			return false
		}
	}
	if validProfileCount != len(cfg.Profiles) {
		return false
	}

	// check cfg.Sink and cfg.SinkPort
	var sinkValid bool
	for _, sink := range a.core.GetSinkList() {
		if sink.Name != cfg.Sink {
			continue
		}

		if len(cfg.SinkPort) == 0 {
			sinkValid = true
			break
		}

		for _, port := range sink.Ports {
			if port.Name == cfg.SinkPort {
				sinkValid = true
			}
		}
		break
	}
	if !sinkValid {
		return false
	}

	// check cfg.Source and cfg.SourcePort
	var sourceValid bool
	for _, source := range a.core.GetSourceList() {
		if source.Name != cfg.Source {
			continue
		}
		if len(cfg.SourcePort) == 0 {
			sourceValid = true
			break
		}

		for _, port := range source.Ports {
			if port.Name == cfg.SourcePort {
				sourceValid = true
			}
		}
		break
	}
	return sourceValid
}
