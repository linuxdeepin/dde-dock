/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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
	"time"

	soundthemeplayer "github.com/linuxdeepin/go-dbus-factory/com.deepin.api.soundthemeplayer"
	"pkg.deepin.io/lib/asound"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/pulse"
)

func (a *Audio) applyConfig() {
	cfg, err := readConfig()
	if err != nil {
		logger.Warning("Read config info failed:", err)
		return
	}

	if !a.isConfigValid(cfg) {
		logger.Warning("Invalid config:", cfg.string())
		a.trySelectBestPort()
		return
	}

	for _, card := range a.ctx.GetCardList() {
		profileName, ok := cfg.Profiles[card.Name]
		if !ok {
			continue
		}

		if card.ActiveProfile.Name != profileName {
			card.SetProfile(profileName)
		}
	}

	var sinkValidity = true
	for _, s := range a.ctx.GetSinkList() {
		if s.Name == cfg.Sink {
			if len(cfg.SinkPort) == 0 {
				sinkValidity = false
				break
			}
			port := pulse.PortInfos(s.Ports).Get(cfg.SinkPort)
			// if port invalid, nothing to do.
			// TODO: some device port can play sound when state is 'NO', how to fix?
			if port == nil {
				sinkValidity = false
				break
			}

			if s.ActivePort.Name != cfg.SinkPort {
				a.ctx.SetSinkPortByIndex(s.Index, cfg.SinkPort)
			}
			cv := s.Volume.SetAvg(cfg.SinkVolume)
			a.ctx.SetSinkVolumeByIndex(s.Index, cv)
			break
		}
	}
	logger.Debug("Audio config sink validity:", sinkValidity, cfg.Sink)
	if sinkValidity {
		a.ctx.SetDefaultSink(cfg.Sink)
	}

	var sourceValidity = true
	for _, s := range a.ctx.GetSourceList() {
		if s.Name == cfg.Source {
			if len(cfg.SourcePort) == 0 {
				sourceValidity = false
				continue
			}
			port := pulse.PortInfos(s.Ports).Get(cfg.SourcePort)
			if port == nil {
				sourceValidity = false
				continue
			}
			if s.ActivePort.Name != cfg.SourcePort {
				a.ctx.SetSourcePortByIndex(s.Index, cfg.SourcePort)
			}
			cv := s.Volume.SetAvg(cfg.SourceVolume)
			a.ctx.SetSourceVolumeByIndex(s.Index, cv)
			break
		}
	}
	logger.Debug("Audio config source validity:", sourceValidity, cfg.Source)
	if sourceValidity {
		a.ctx.SetDefaultSource(cfg.Source)
	}
}

func (a *Audio) trySelectBestPort() {
	logger.Debug("trySelectBestPort")
	cardId, sinkPort := a.cards.getPassablePort(pulse.DirectionSink)
	if sinkPort != nil {
		logger.Debugf("switch to sink port %s, avail: %s",
			sinkPort.Name, portAvailToString(sinkPort.Available))
		err := a.setPort(cardId, sinkPort.Name, sinkPort.Direction)
		if err != nil {
			logger.Warningf("failed to switch to sink port %s: %v",
				sinkPort.Name, err)
		}
	}

	cardId, sourcePort := a.cards.getPassablePort(pulse.DirectionSource)
	if sourcePort != nil {
		logger.Debugf("switch to source port %s, avail: %s",
			sourcePort.Name, portAvailToString(sourcePort.Available))
		err := a.setPort(cardId, sourcePort.Name, pulse.DirectionSource)
		if err != nil {
			logger.Warningf("failed to switch to source port %s: %v",
				sourcePort.Name, err)
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
	var info = config{
		Profiles: make(map[string]string),
	}

	ctx := a.context()
	for _, card := range ctx.GetCardList() {
		info.Profiles[card.Name] = card.ActiveProfile.Name
	}

	for _, sinkInfo := range ctx.GetSinkList() {
		if a.getDefaultSinkName() != sinkInfo.Name {
			continue
		}

		info.Sink = sinkInfo.Name
		info.SinkPort = sinkInfo.ActivePort.Name
		info.SinkVolume = sinkInfo.Volume.Avg()
		break
	}

	for _, sourceInfo := range ctx.GetSourceList() {
		if a.getDefaultSourceName() != sourceInfo.Name {
			continue
		}

		info.Source = sourceInfo.Name
		info.SourcePort = sourceInfo.ActivePort.Name
		info.SourceVolume = sourceInfo.Volume.Avg()
		break
	}

	readConfig()
	err := saveConfig(&info)
	if err != nil {
		logger.Warning("Save config file failed:", info.string(), err)
	}

	err = a.saveAudioState()
	if err != nil {
		logger.Warning(err)
	}
}

func (a *Audio) saveAudioState() error {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	sink := a.getDefaultSink()
	sink.PropsMu.RLock()
	device := sink.props["alsa.device"]
	card := sink.props["alsa.card"]
	mute := sink.Mute
	sink.PropsMu.RUnlock()

	cardId, err := toALSACardId(card)
	if err != nil {
		return err
	}

	activePlayback := map[string]dbus.Variant{
		"card":   dbus.MakeVariant(cardId),
		"device": dbus.MakeVariant(device),
		"mute":   dbus.MakeVariant(mute),
	}

	player := soundthemeplayer.NewSoundThemePlayer(sysBus)
	err = player.SaveAudioState(0, activePlayback)
	return err
}

func toALSACardId(idx string) (cardId string, err error) {
	ctl, err := asound.CTLOpen("hw:"+idx, 0)
	if err != nil {
		return
	}
	defer ctl.Close()

	cardInfo, err := asound.NewCTLCardInfo()
	if err != nil {
		return
	}
	defer cardInfo.Free()

	err = ctl.CardInfo(cardInfo)
	if err != nil {
		return
	}

	cardId = cardInfo.GetID()
	return
}

func (a *Audio) isConfigValid(cfg *config) bool {
	if len(cfg.Profiles) == 0 {
		return false
	}

	// check cfg.Profiles
	var validProfileCount int
	for _, card := range a.ctx.GetCardList() {
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
	for _, sink := range a.ctx.GetSinkList() {
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
	for _, source := range a.ctx.GetSourceList() {
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
