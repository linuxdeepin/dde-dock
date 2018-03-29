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

package appearance

import (
	"os"
	"path"
	"strings"
	"time"

	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
)

var (
	gtkDirs  []string
	iconDirs []string
	bgDirs   []string
)

var prevTimestamp int64

func (m *Manager) handleThemeChanged() {
	if m.watcher == nil {
		return
	}

	m.watchGtkDirs()
	m.watchIconDirs()
	m.watchBgDirs()

	for {
		select {
		case <-m.endWatcher:
			logger.Debug("[Fsnotify] quit watch")
			return
		case err := <-m.watcher.Error:
			logger.Warning("Receive file watcher error:", err)
			return
		case ev, ok := <-m.watcher.Event:
			if !ok {
				logger.Error("Invalid event:", ev)
				return
			}

			timestamp := time.Now().UnixNano()
			tmp := timestamp - prevTimestamp
			logger.Debug("[Fsnotify] timestamp:", prevTimestamp, timestamp, tmp, ev)
			prevTimestamp = timestamp
			// Filter time duration < 100ms's event
			if tmp > 100000000 {
				<-time.After(time.Millisecond * 100)
				file := ev.Name
				logger.Debug("[Fsnotify] changed file:", file)
				switch {
				case hasEventOccurred(file, bgDirs):
					logger.Debug("fs event in bgDirs")
					background.NotifyChanged()
				case hasEventOccurred(file, gtkDirs):
					logger.Debug("fs event in gtkDirs")
					// Wait for theme copy finished
					<-time.After(time.Millisecond * 700)
					subthemes.RefreshGtkThemes()
					m.emitSignalRefreshed(TypeGtkTheme)
				case hasEventOccurred(file, iconDirs):
					// Wait for theme copy finished
					logger.Debug("fs event in iconDirs")
					<-time.After(time.Millisecond * 700)
					subthemes.RefreshIconThemes()
					subthemes.RefreshCursorThemes()
					m.emitSignalRefreshed(TypeIconTheme)
					m.emitSignalRefreshed(TypeCursorTheme)
				}
			}
		}
	}
}

func (m *Manager) watchGtkDirs() {
	var home = os.Getenv("HOME")
	gtkDirs = []string{
		path.Join(home, ".local/share/themes"),
		path.Join(home, ".themes"),
		"/usr/local/share/themes",
		"/usr/share/themes",
	}

	m.watchDirs(gtkDirs)
}

func (m *Manager) watchIconDirs() {
	var home = os.Getenv("HOME")
	iconDirs = []string{
		path.Join(home, ".local/share/icons"),
		path.Join(home, ".icons"),
		"/usr/local/share/icons",
		"/usr/share/icons",
	}

	m.watchDirs(iconDirs)
}

func (m *Manager) watchBgDirs() {
	bgDirs = background.ListDirs()
	m.watchDirs(bgDirs)
}

func (m *Manager) watchDirs(dirs []string) {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Debugf("Mkdir '%s' failed: %v", dir, err)
		}

		err = m.watcher.Watch(dir)
		if err != nil {
			logger.Debugf("Watch dir '%s' failed: %v", dir, err)
		}
	}
}

func hasEventOccurred(ev string, list []string) bool {
	for _, v := range list {
		if strings.Contains(ev, v) {
			return true
		}
	}
	return false
}

func (m *Manager) emitSignalRefreshed(type0 string) {
	m.service.Emit(m, "Refreshed", type0)
}
