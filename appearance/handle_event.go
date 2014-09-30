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

package appearance

import (
	"github.com/howeyc/fsnotify"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"regexp"
)

func (m *Manager) listenGSettings() {
	m.settings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case deepinGSKeyTheme:
			m.applyTheme(m.settings.GetString(key))
		case deepinGSKeyPicture:
			m.Set("background", m.settings.GetString(key))
		case deepinGSKeySound:
			m.Set("sound", m.settings.GetString(key))
		case deepinGSKeyGreeter:
			m.greeter.Set(m.settings.GetString(key))
		}
	})

	if m.gnomeSettings == nil {
		return
	}

	m.gnomeSettings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case "picture-uri":
			m.bg.Set(m.gnomeSettings.GetString(key))
		}
	})
}

func (t *Theme) handleEvent(ev *fsnotify.FileEvent) {
	if ev == nil {
		return
	}

	if ok, _ := regexp.MatchString(`\.swa?px?$`, ev.Name); ok {
		return
	}

	if ev.IsDelete() {
		t.watcher.ResetFileListWatch()
		return
	}

	t.readFromFile()
	t.setPropPreview(t.getPreviewList())
	if t.eventHandler != nil {
		//apply the change
		t.eventHandler(t.Name)
	}
}
