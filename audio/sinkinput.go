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

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/procfs"
	"pkg.deepin.io/lib/pulse"
)

const (
	PropAppIconName = "application.icon_name"
	PropAppName     = "application.name"
	PropAppPID      = "application.process.id"
)

type SinkInput struct {
	service          *dbusutil.Service
	PropsMu          sync.RWMutex
	core             *pulse.SinkInput
	index            uint32
	correctAppCalled bool
	// Name process name
	Name           string
	Icon           string
	Mute           bool
	Volume         float64
	Balance        float64
	SupportBalance bool
	Fade           float64
	SupportFade    bool

	methods *struct {
		SetVolume  func() `in:"value,isPlay"`
		SetBalance func() `in:"value,isPlay"`
		SetFade    func() `in:"value"`
		SetMute    func() `in:"value"`
	}
}

func NewSinkInput(core *pulse.SinkInput, service *dbusutil.Service) *SinkInput {
	if core == nil {
		return nil
	}
	s := &SinkInput{
		core:    core,
		service: service,
	}
	s.index = s.core.Index
	s.update()
	return s
}

func (s *SinkInput) SetVolume(v float64, isPlay bool) *dbus.Error {
	if !isVolumeValid(v) {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	if v == 0 {
		v = 0.001
	}
	s.core.SetVolume(s.core.Volume.SetAvg(v))
	if isPlay {
		playFeedback()
	}
	return nil
}

func (s *SinkInput) SetBalance(v float64, isPlay bool) *dbus.Error {
	if v < -1.00 || v > 1.00 {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
	if isPlay {
		playFeedback()
	}
	return nil
}

func (s *SinkInput) SetFade(v float64) *dbus.Error {
	if v < -1.00 || v > 1.00 {
		return dbusutil.ToError(fmt.Errorf("invalid volume value: %v", v))
	}

	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
	playFeedback()
	return nil
}

func (s *SinkInput) SetMute(v bool) *dbus.Error {
	s.core.SetMute(v)
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

// correct app's name and icon
func (s *SinkInput) correctApp() error {
	if s.correctAppCalled {
		return nil
	}
	s.correctAppCalled = true

	pidStr := s.core.PropList[PropAppPID]
	pid, err := strconv.ParseUint(pidStr, 10, 32)
	if err != nil {
		return err
	}
	process := procfs.Process(pid)
	cmdline, err := process.Cmdline()
	if err != nil {
		return err
	}
	logger.Debugf("cmdline: %#v", cmdline)
	cmd := strings.Join(cmdline, " ")

	switch {
	case strings.Contains(cmd, "deepin-movie"):
		s.Name = "Deepin Movie"
		s.Icon = "deepin-movie"
	case strings.Contains(cmd, "firefox"):
		s.Name = "Firefox"
		s.Icon = "firefox"
	case strings.Contains(cmd, "maxthon"):
		s.Name = "Maxthon"
		s.Icon = "maxthon-browser"
	case (strings.Contains(cmd, "chrome") && strings.Contains(cmd, "google")):
		s.Name = "Google Chrome"
		s.Icon = "google-chrome"
	case strings.Contains(cmd, "deepin-music-player"):
		s.Name = "Deepin Music"
		s.Icon = "deepin-music-player"
	case strings.Contains(cmd, "smplayer"):
		s.Name = "SMPlayer"
		s.Icon = "smplayer"
	}
	return nil
}

func (s *SinkInput) update() {
	s.Name = s.core.PropList[PropAppName]
	s.Icon = s.core.PropList[PropAppIconName]

	err := s.correctApp()
	if err != nil {
		logger.Warning(err)
	}

	if len(s.Icon) == 0 {
		// Using default media player icon
		s.Icon = "media-player"
	}

	s.PropsMu.Lock()
	s.setPropVolume(s.core.Volume.Avg())
	s.setPropMute(s.core.Mute)

	s.setPropSupportFade(false)
	s.setPropFade(s.core.Volume.Fade(s.core.ChannelMap))
	s.setPropSupportBalance(true)
	s.setPropBalance(s.core.Volume.Balance(s.core.ChannelMap))
	s.PropsMu.Unlock()
}
