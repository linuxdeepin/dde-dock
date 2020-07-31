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
	"fmt"
	"strconv"
	"strings"
	"sync"

	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/procfs"
	"pkg.deepin.io/lib/pulse"
)

const (
	PropAppIconName      = "application.icon_name"
	PropAppName          = "application.name"
	PropAppProcessID     = "application.process.id"
	PropAppProcessBinary = "application.process.binary"
)

type SinkInput struct {
	audio             *Audio
	service           *dbusutil.Service
	PropsMu           sync.RWMutex
	index             uint32
	correctIconCalled bool
	correctedIcon     string
	visible           bool
	cVolume           pulse.CVolume
	channelMap        pulse.ChannelMap
	// Name process name
	Name           string
	Icon           string
	Mute           bool
	Volume         float64
	Balance        float64
	SupportBalance bool
	Fade           float64
	SupportFade    bool
	SinkIndex      uint32

	// nolint
	methods *struct {
		SetVolume  func() `in:"value,isPlay"`
		SetBalance func() `in:"value,isPlay"`
		SetFade    func() `in:"value"`
		SetMute    func() `in:"value"`
	}
}

func newSinkInput(sinkInputInfo *pulse.SinkInput, audio *Audio) *SinkInput {
	if sinkInputInfo == nil {
		return nil
	}
	sinkInput := &SinkInput{
		audio:   audio,
		service: audio.service,
		index:   sinkInputInfo.Index,
		visible: getSinkInputVisible(sinkInputInfo),
	}
	sinkInput.update(sinkInputInfo)
	return sinkInput
}

func (s *SinkInput) getPropSinkIndex() uint32 {
	s.PropsMu.RLock()
	v := s.SinkIndex
	s.PropsMu.RUnlock()
	return v
}

func getSinkInputVisible(sinkInputInfo *pulse.SinkInput) bool {
	appName := sinkInputInfo.PropList[pulse.PA_PROP_APPLICATION_NAME]
	switch appName {
	case "com.deepin.SoundEffect", "deepin-notifications":
		return false
	}

	switch sinkInputInfo.PropList[pulse.PA_PROP_MEDIA_ROLE] {
	case "video", "music", "game":
		return true
	case "animation", "production", "phone":
		//TODO: what's the meaning of this type? Should we filter this SinkInput?
		return true
	case "event", "a11y", "test", "filter":
		return false
	default:
		return true
	}
}

func (s *SinkInput) SetVolume(v float64, isPlay bool) *dbus.Error {
	if !isVolumeValid(v) {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	if v == 0 {
		v = 0.001
	}
	s.PropsMu.RLock()
	cv := s.cVolume.SetAvg(v)
	s.PropsMu.RUnlock()
	s.audio.context().SetSinkInputVolume(s.index, cv)

	if isPlay {
		playFeedback()
	}
	return nil
}

func (s *SinkInput) SetBalance(v float64, isPlay bool) *dbus.Error {
	if v < -1.00 || v > 1.00 {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	s.PropsMu.RLock()
	cv := s.cVolume.SetBalance(s.channelMap, v)
	s.PropsMu.RUnlock()
	s.audio.context().SetSinkInputVolume(s.index, cv)

	if isPlay {
		playFeedback()
	}
	return nil
}

func (s *SinkInput) SetFade(v float64) *dbus.Error {
	if v < -1.00 || v > 1.00 {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	s.PropsMu.RLock()
	cv := s.cVolume.SetFade(s.channelMap, v)
	s.PropsMu.RUnlock()
	s.audio.context().SetSinkInputVolume(s.index, cv)

	playFeedback()
	return nil
}

func (s *SinkInput) SetMute(v bool) *dbus.Error {
	s.audio.context().SetSinkInputMute(s.index, v)
	if !v {
		playFeedback()
	}
	return nil
}

func (s *SinkInput) getPath() dbus.ObjectPath {
	return dbus.ObjectPath(dbusPath + "/SinkInput" + strconv.Itoa(int(s.index)))
}

func (*SinkInput) GetInterfaceName() string {
	return dbusInterface + ".SinkInput"
}

func getProcessParentCmdline(pidStr string) ([]string, error) {
	pid, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		return nil, err
	}

	p := procfs.Process(pid)
	status, err := p.Status()
	if err != nil {
		return nil, err
	}
	ppid, err := status.PPid()
	if err != nil {
		return nil, err
	}
	pp := procfs.Process(ppid)
	return pp.Cmdline()
}

func isProcessParentFirefox(pidStr string) (bool, error) {
	cmdline, err := getProcessParentCmdline(pidStr)
	if err != nil {
		return false, err
	}
	if len(cmdline) > 0 && strings.Contains(cmdline[0], "firefox") {
		return true, nil
	}
	return false, nil
}

func isProcessParentSMPlayer(pidStr string) (bool, error) {
	cmdline, err := getProcessParentCmdline(pidStr)
	if err != nil {
		return false, err
	}
	if len(cmdline) > 0 && strings.Contains(cmdline[0], "smplayer") {
		return true, nil
	}
	return false, nil
}

// correct icon
func (s *SinkInput) correctIcon(sinkInputInfo *pulse.SinkInput) (string, error) {
	if s.correctIconCalled {
		return s.correctedIcon, nil
	}
	s.correctIconCalled = true

	processBin := sinkInputInfo.PropList[PropAppProcessBinary]
	appPid := sinkInputInfo.PropList[PropAppProcessID]
	appName := sinkInputInfo.PropList[PropAppName]

	var icon string
	switch processBin {
	case "firefox":
		icon = "firefox"
	case "plugin-container":
		// may be flash player embed in firefox
		is, err := isProcessParentFirefox(appPid)
		if err != nil {
			logger.Warning(err)
			break
		}
		if is {
			icon = "firefox"
		}
	case "mpv":
		if appName == "SMPlayer" {
			icon = "smplayer"
		}

	case "mplayer":
		is, err := isProcessParentSMPlayer(appPid)
		if err != nil {
			logger.Warning(err)
			break
		}
		if is {
			icon = "smplayer"
		}

	case "wine-preloader":
		if appName == "foobar2000.exe" {
			icon = "apps.org.foobar2000"
		}

	case "python2.7":
		if appName == "foobnix" {
			icon = "foobnix"
		}
	case "python3.5":
		if appName == "com.github.geigi.cozy" {
			icon = "com.github.geigi.cozy"
		}

	case "cocomusic":
		icon = "cocomusic"
	case "Museeks":
		icon = "museeks"
	case "cumulonimbus":
		icon = "cumulonimbus"
	case "yarock":
		icon = "application-x-yarock"
	case "mixnode":
		icon = "mixnode"
	case "headset":
		icon = "headset"
	case "electron-xiami":
		icon = "electron_xiami"
	}

	s.correctedIcon = icon
	if icon != "" {
		logger.Debugf("correct icon of sink-input #%d to %q", sinkInputInfo.Index, icon)
	}
	return icon, nil
}

func (s *SinkInput) update(sinkInputInfo *pulse.SinkInput) {
	s.PropsMu.Lock()
	defer s.PropsMu.Unlock()

	if !s.visible {
		s.SinkIndex = sinkInputInfo.Sink
		return
	}

	s.cVolume = sinkInputInfo.Volume
	s.channelMap = sinkInputInfo.ChannelMap
	s.setPropSinkIndex(sinkInputInfo.Sink)
	name := sinkInputInfo.PropList[PropAppName]
	s.setPropName(name)
	icon := sinkInputInfo.PropList[PropAppIconName]
	correctedIcon, err := s.correctIcon(sinkInputInfo)
	if err != nil {
		logger.Warning(err)
	}
	if correctedIcon != "" {
		icon = correctedIcon
	}
	if icon == "" {
		// Using default media player icon
		icon = "media-player"
	}
	s.setPropIcon(icon)

	s.setPropVolume(sinkInputInfo.Volume.Avg())
	s.setPropMute(sinkInputInfo.Mute)

	s.setPropSupportFade(false)
	s.setPropFade(sinkInputInfo.Volume.Fade(sinkInputInfo.ChannelMap))
	s.setPropSupportBalance(true)
	s.setPropBalance(sinkInputInfo.Volume.Balance(sinkInputInfo.ChannelMap))
}
