/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/procfs"
	"pkg.deepin.io/lib/pulse"
	"strconv"
	"strings"
)

const (
	PropAppIconName = "application.icon_name"
	PropAppName     = "application.name"
	PropAppPID      = "application.process.id"
)

type SinkInput struct {
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
}

func NewSinkInput(core *pulse.SinkInput) *SinkInput {
	if core == nil {
		return nil
	}
	s := &SinkInput{core: core}
	s.index = s.core.Index
	s.update()
	return s
}

func (s *SinkInput) SetVolume(v float64, isPlay bool) error {
	if !isVolumeValid(v) {
		return fmt.Errorf("Invalid volume value: %v", v)
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

func (s *SinkInput) SetBalance(v float64, isPlay bool) error {
	if v < -1.00 || v > 1.00 {
		return fmt.Errorf("Invalid volume value: %v", v)
	}

	s.core.SetVolume(s.core.Volume.SetBalance(s.core.ChannelMap, v))
	if isPlay {
		playFeedback()
	}
	return nil
}

func (s *SinkInput) SetFade(v float64) error {
	if v < -1.00 || v > 1.00 {
		return fmt.Errorf("Invalid volume value: %v", v)
	}

	s.core.SetVolume(s.core.Volume.SetFade(s.core.ChannelMap, v))
	playFeedback()
	return nil
}

func (s *SinkInput) SetMute(v bool) {
	s.core.SetMute(v)
	if !v {
		playFeedback()
	}
}

func (s *SinkInput) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       baseBusName,
		ObjectPath: fmt.Sprintf("%s/SinkInput%d", baseBusPath, s.index),
		Interface:  baseBusIfc + ".SinkInput",
	}
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

	s.setPropVolume(s.core.Volume.Avg())
	s.setPropMute(s.core.Mute)

	s.setPropSupportFade(false)
	s.setPropFade(s.core.Volume.Fade(s.core.ChannelMap))
	s.setPropSupportBalance(true)
	s.setPropBalance(s.core.Volume.Balance(s.core.ChannelMap))

}

func (s *SinkInput) setPropVolume(v float64) {
	if s.Volume != v {
		s.Volume = v
		dbus.NotifyChange(s, "Volume")
	}
}

func (s *SinkInput) setPropMute(v bool) {
	if s.Mute != v {
		s.Mute = v
		dbus.NotifyChange(s, "Mute")
	}
}

func (s *SinkInput) setPropBalance(v float64) {
	if s.Balance != v {
		s.Balance = v
		dbus.NotifyChange(s, "Balance")
	}
}

func (s *SinkInput) setPropSupportBalance(v bool) {
	if s.SupportBalance != v {
		s.SupportBalance = v
		dbus.NotifyChange(s, "SupportBalance")
	}
}

func (s *SinkInput) setPropSupportFade(v bool) {
	if s.SupportFade != v {
		s.SupportFade = v
		dbus.NotifyChange(s, "SupportFade")
	}
}

func (s *SinkInput) setPropFade(v float64) {
	if s.Fade != v {
		s.Fade = v
		dbus.NotifyChange(s, "Fade")
	}
}
