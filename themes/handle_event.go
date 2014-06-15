/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package themes

import (
	"dlib/gio-2.0"
	"github.com/howeyc/fsnotify"
	"os"
	"path"
	"regexp"
)

func (obj *Manager) listenGSettings() {
	themeSettings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case GS_KEY_CURRENT_THEME:
			value := themeSettings.GetString(key)
			//if obj.CurrentTheme.GetValue().(string) == value {
			//return
			//}
			if t, ok := obj.themeObjMap[value]; ok {
				t.setAllThemes()
			} else {
				obj.setPropCurrentTheme(DEFAULT_THEME)
			}
		case GS_KEY_CURRENT_SOUND:
			value := themeSettings.GetString(key)
			obj.setSoundTheme(value)
		case GS_KEY_CURRENT_BG:
			value := themeSettings.GetString(key)
			obj.setBackground(value)
		}
	})

	gnmSettings.Connect("changed::picture-uri", func(s *gio.Settings, key string) {
		value := gnmSettings.GetString("picture-uri")
		bg := themeSettings.GetString(GS_KEY_CURRENT_BG)
		if bg != value {
			themeSettings.SetString(GS_KEY_CURRENT_BG, value)
		}
	})
}

func (obj *Manager) startWatch() {
	if obj.watcher == nil {
		var err error
		obj.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			Logger.Errorf("New Watcher Failed: %v", err)
			panic(err)
		}
	}

	homeDir, _ := objUtil.GetHomeDir()
	obj.watcher.Watch(THEME_SYS_PATH)
	dir := path.Join(homeDir, THEME_LOCAL_PATH)
	if !objUtil.IsFileExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			Logger.Errorf("Mkdir '%s' failed: %v", dir, err)
		} else {
			obj.watcher.Watch(dir)
		}
	}
	obj.watcher.Watch(ICON_SYS_PATH)
	dir = path.Join(homeDir, ICON_LOCAL_PATH)
	if !objUtil.IsFileExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			Logger.Errorf("Mkdir '%s' failed: %v", dir, err)
		} else {
			obj.watcher.Watch(dir)
		}
	}
	obj.watcher.Watch(SOUND_THEME_PATH)
	obj.watcher.Watch(PERSON_SYS_THEME_PATH)
	dir = path.Join(homeDir, PERSON_LOCAL_THEME_PATH)
	if !objUtil.IsFileExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			Logger.Errorf("Mkdir '%s' failed: %v", dir, err)
		} else {
			obj.watcher.Watch(dir)
		}
	}
}

func (obj *Manager) endWatch() {
	if obj.watcher == nil {
		return
	}

	homeDir, _ := objUtil.GetHomeDir()
	obj.watcher.RemoveWatch(THEME_SYS_PATH)
	dir := path.Join(homeDir, THEME_LOCAL_PATH)
	if objUtil.IsFileExist(dir) {
		obj.watcher.RemoveWatch(dir)
	}
	obj.watcher.RemoveWatch(ICON_SYS_PATH)
	dir = path.Join(homeDir, ICON_LOCAL_PATH)
	if objUtil.IsFileExist(dir) {
		obj.watcher.RemoveWatch(dir)
	}
	obj.watcher.RemoveWatch(SOUND_THEME_PATH)
	obj.watcher.RemoveWatch(PERSON_SYS_THEME_PATH)
	dir = path.Join(homeDir, PERSON_LOCAL_THEME_PATH)
	if objUtil.IsFileExist(dir) {
		obj.watcher.RemoveWatch(dir)
	}
}

func (obj *Manager) handleEvent() {
	for {
		select {
		case <-obj.quitFlag:
			return
		case ev, ok := <-obj.watcher.Event:
			if !ok {
				obj.endWatch()
				obj.startWatch()
				break
			}

			if ev == nil {
				break
			}

			ok1, _ := regexp.MatchString(THEME_SYS_PATH, ev.Name)
			ok2, _ := regexp.MatchString(THEME_LOCAL_PATH, ev.Name)
			if ok1 || ok2 {
				obj.setPropGtkThemeList(obj.getGtkStrList())
			}

			ok1, _ = regexp.MatchString(ICON_SYS_PATH, ev.Name)
			ok2, _ = regexp.MatchString(ICON_LOCAL_PATH, ev.Name)
			if ok1 || ok2 {
				obj.setPropIconThemeList(obj.getIconStrList())
				obj.setPropCursorThemeList(obj.getCursorStrList())
			}

			ok1, _ = regexp.MatchString(SOUND_THEME_PATH, ev.Name)
			if ok1 {
				obj.setPropSoundThemeList(obj.getSoundStrList())
			}

			ok1, _ = regexp.MatchString(PERSON_SYS_THEME_PATH, ev.Name)
			ok2, _ = regexp.MatchString(PERSON_LOCAL_THEME_PATH, ev.Name)
			if ok1 || ok2 {
				obj.rebuildThemes()
				obj.setPropThemeList(obj.getDThemeStrList())
			}
		case err, ok := <-obj.watcher.Error:
			if !ok || err != nil {
				obj.endWatch()
				obj.startWatch()
			}
		}
	}
}

func (obj *Manager) handleBgEvent() {
	for {
		select {
		case <-obj.bgQuitFlag:
			return
		case ev, ok := <-obj.bgWatcher.Event:
			if !ok {
				obj.endBgWatch()
				obj.startBgWatch()
			}

			if ev == nil {
				break
			}

			if ev.IsDelete() || ev.IsCreate() {
				obj.setPropBackgroundList(obj.getBgStrList())
			}
		case err, ok := <-obj.bgWatcher.Error:
			if !ok || err != nil {
				obj.endBgWatch()
				obj.startBgWatch()
			}
		}
	}
}

func (obj *Manager) startBgWatch() {
	if obj.bgWatcher == nil {
		var err error
		obj.bgWatcher, err = fsnotify.NewWatcher()
		if err != nil {
			Logger.Errorf("New Watcher Failed: %v", err)
			panic(err)
		}
	}

	pict := getUserPictureDir()
	userBG := path.Join(pict, "Wallpapers")
	if !objUtil.IsFileExist(userBG) {
		os.MkdirAll(userBG, 0755)
	}

	obj.bgWatcher.Watch(userBG)
	obj.bgWatcher.Watch(DEFAULT_SYS_BG_DIR)
}

func (obj *Manager) endBgWatch() {
	if obj.bgWatcher == nil {
		return
	}

	pict := getUserPictureDir()
	userBG := path.Join(pict, "Wallpapers")
	if objUtil.IsFileExist(userBG) {
		obj.bgWatcher.RemoveWatch(userBG)
	}

	obj.bgWatcher.RemoveWatch(DEFAULT_SYS_BG_DIR)
}

func (obj *Theme) startWatch() {
	if obj.Type == THEME_TYPE_SYSTEM {
		return
	}

	if obj.watcher == nil {
		var err error
		obj.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			Logger.Errorf("New Watcher Failed: %v", err)
			panic(err)
		}
	}

	filename := path.Join(obj.filePath, "theme.ini")
	obj.watcher.Watch(filename)
}

func (obj *Theme) endWatch() {
	if obj.watcher == nil {
		return
	}

	filename := path.Join(obj.filePath, "theme.ini")
	obj.watcher.RemoveWatch(filename)
}

func (obj *Theme) handleEvent() {
	for {
		select {
		case <-obj.quitFlag:
			return
		case ev, ok := <-obj.watcher.Event:
			if !ok {
				obj.endWatch()
				obj.startWatch()
				break
			}

			if ok, _ := regexp.MatchString(`\.swa?px?$`, ev.Name); ok {
				break
			}

			if ev.IsDelete() {
				obj.endWatch()
				obj.startWatch()
				break
			}

			if ev.IsModify() {
				obj.setAllProps()
				obj.setAllThemes()
			}
		case err, ok := <-obj.watcher.Error:
			if !ok || err != nil {
				obj.endWatch()
				obj.startWatch()
			}
		}
	}
}
