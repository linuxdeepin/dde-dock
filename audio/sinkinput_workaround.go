/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import (
	"pkg.deepin.io/lib/procfs"
	"strconv"
	"strings"
)

func (s *SinkInput) correctAppName() error {
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
