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
			time.Sleep(time.Microsecond * 500)
		}
	}

	a.core.SetDefaultSink(info.Sink)
	a.core.SetDefaultSource(info.Source)

	for _, s := range a.core.GetSinkList() {
		if s.Name == info.Sink {
			if len(info.SinkPort) != 0 {
				port := pulse.PortInfos(s.Ports).Get(info.SinkPort)
				if port != nil && port.Available != pulse.AvailableTypeNo &&
					s.ActivePort.Name != info.SinkPort {
					s.SetPort(info.SinkPort)
				}
			}
			s.SetVolume(s.Volume.SetAvg(info.SinkVolume))
			break
		}
	}

	for _, s := range a.core.GetSourceList() {
		if s.Name == info.Source {
			if len(info.SourcePort) != 0 {
				port := pulse.PortInfos(s.Ports).Get(info.SourcePort)
				if port != nil && port.Available != pulse.AvailableTypeNo &&
					s.ActivePort.Name != info.SourcePort {
					s.SetPort(info.SourcePort)
				}
			}
			s.SetVolume(s.Volume.SetAvg(info.SourceVolume))
			break
		}
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

func (a *Audio) isConfigValid(info *configInfo) bool {
	var (
		cardNumber  int
		sinkValid   bool
		sourceValid bool
	)

	for _, card := range a.core.GetCardList() {
		v, ok := info.Profiles[card.Name]
		if !ok {
			continue
		}

		for _, profile := range card.Profiles {
			if profile.Name == v {
				cardNumber += 1
				break
			}
		}
	}
	if cardNumber != len(info.Profiles) {
		return false
	}

	for _, sink := range a.core.GetSinkList() {
		if sink.Name == info.Sink {
			if len(info.SinkPort) == 0 {
				sinkValid = true
				break
			}

			for _, port := range sink.Ports {
				if port.Name == info.SinkPort {
					sinkValid = true
				}
			}
			break
		}
	}
	if !sinkValid {
		return false
	}

	for _, source := range a.core.GetSourceList() {
		if source.Name == info.Source {
			if len(info.SourcePort) == 0 {
				sourceValid = true
				break
			}

			for _, port := range source.Ports {
				if port.Name == info.SourcePort {
					sourceValid = true
				}
			}
			break
		}
	}
	return sourceValid
}
