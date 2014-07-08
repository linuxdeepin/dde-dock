/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
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
	if err == nil {
		if strings.Contains(string(contents), "deepin-movie") {
			s.Name = "Deepin Movie"
			s.Icon = "deepin-movie"
		} else if strings.Contains(string(contents), "firefox") {
			s.Name = "Firefox"
			s.Icon = "firefox"
		} else if strings.Contains(string(contents), "maxthon") {
			s.Name = "Maxthon"
			s.Icon = "maxthon-browser"
		} else if strings.Contains(string(contents), "chrome") &&
			strings.Contains(string(contents), "google") {
			s.Name = "Google Chrome"
			s.Icon = "google-chrome"
		}
	}
}
