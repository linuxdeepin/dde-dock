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
	"io/ioutil"
	"path"
	"strings"
)

func (s *SinkInput) correctAppName() {
	pid := s.core.PropList[PropAppPID]
	filePath := path.Join("/proc", pid, "cmdline")
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Warningf("ReadFile '%s' failed: %v", filePath, err)
		return
	}

	ctx := string(contents)
	switch {
	case strings.Contains(ctx, "deepin-movie"):
		s.Name = "Deepin Movie"
		s.Icon = "deepin-movie"
	case strings.Contains(ctx, "firefox"):
		s.Name = "Firefox"
		s.Icon = "firefox"
	case strings.Contains(ctx, "maxthon"):
		s.Name = "Maxthon"
		s.Icon = "maxthon-browser"
	case (strings.Contains(ctx, "chrome") && strings.Contains(ctx, "google")):
		s.Name = "Google Chrome"
		s.Icon = "google-chrome"
	case strings.Contains(ctx, "deepin-music-player"):
		s.Name = "Deepin Music"
		s.Icon = "deepin-music-player"
	case strings.Contains(ctx, "smplayer"):
		s.Name = "SMPlayer"
		s.Icon = "smplayer"
	}
}
