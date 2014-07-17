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
	"github.com/howeyc/fsnotify"
	"os"
	"path"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"regexp"
	"time"
)

var (
	gnomeIfcSettings = gio.NewSettings("org.gnome.desktop.interface")
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
				//t.setAllProps()
			} else {
				obj.setPropCurrentTheme(DEFAULT_THEME)
			}
		case GS_KEY_CURRENT_SOUND:
			value := themeSettings.GetString(key)
			obj.setSoundTheme(value)
		case GS_KEY_CURRENT_BG:
			value := themeSettings.GetString(key)
			if !obj.setBackground(decodeURI(value)) {
				obj.setBackground(DEFAULT_BG)
			}
		}
	})

	gnmSettings.Connect("changed::picture-uri", func(s *gio.Settings, key string) {
		value := gnmSettings.GetString("picture-uri")
		value = decodeURI(value)
		bg := themeSettings.GetString(GS_KEY_CURRENT_BG)
		bg = decodeURI(bg)
		if bg != value {
			value = encodeURI(value)
			themeSettings.SetString(GS_KEY_CURRENT_BG, value)
		}
	})

	gnomeIfcSettings.Connect("changed::font-name", func(s *gio.Settings, key string) {
		value := gnomeIfcSettings.GetString("font-name")
		setGtkFont(value)
		setQtFont(value)
	})
}

func (obj *Manager) startWatch() {
	if obj.watcher == nil {
		var err error
		obj.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			Logger.Warningf("New Watcher Failed: %v", err)
			panic(err)
		}
	}

	obj.watcher.Watch(THEME_SYS_PATH)

	errFlag := false
	dir := path.Join(homeDir, THEME_LOCAL_PATH)
	if !dutils.IsFileExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			errFlag = true
			Logger.Warningf("Mkdir '%s' failed: %v", dir, err)
		}
	}
	if !errFlag {
		obj.watcher.Watch(dir)
	}

	obj.watcher.Watch(ICON_SYS_PATH)

	errFlag = false
	dir = path.Join(homeDir, ICON_LOCAL_PATH)
	if !dutils.IsFileExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			errFlag = true
			Logger.Warningf("Mkdir '%s' failed: %v", dir, err)
		}
	}
	if !errFlag {
		obj.watcher.Watch(dir)
	}

	obj.watcher.Watch(SOUND_THEME_PATH)

	obj.watcher.Watch(PERSON_SYS_THEME_PATH)

	errFlag = false
	dir = path.Join(homeDir, PERSON_LOCAL_THEME_PATH)
	if !dutils.IsFileExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			errFlag = true
			Logger.Warningf("Mkdir '%s' failed: %v", dir, err)
		}
	}
	if !errFlag {
		obj.watcher.Watch(dir)
	}

	obj.watcher.Watch(PERSON_SYS_GREETER_PATH)
	errFlag = false
	dir = path.Join(homeDir, PERSON_LOCAL_GREETER_PATH)
	if !dutils.IsFileExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			errFlag = true
			Logger.Warningf("Mkdir '%s' failed: %v", dir, err)
		}
	}
	if !errFlag {
		obj.watcher.Watch(dir)
	}
}

func (obj *Manager) endWatch() {
	if obj.watcher == nil {
		return
	}

	obj.watcher.RemoveWatch(THEME_SYS_PATH)
	dir := path.Join(homeDir, THEME_LOCAL_PATH)
	if dutils.IsFileExist(dir) {
		obj.watcher.RemoveWatch(dir)
	}

	obj.watcher.RemoveWatch(ICON_SYS_PATH)
	dir = path.Join(homeDir, ICON_LOCAL_PATH)
	if dutils.IsFileExist(dir) {
		obj.watcher.RemoveWatch(dir)
	}

	obj.watcher.RemoveWatch(SOUND_THEME_PATH)

	obj.watcher.RemoveWatch(PERSON_SYS_THEME_PATH)
	dir = path.Join(homeDir, PERSON_LOCAL_THEME_PATH)
	if dutils.IsFileExist(dir) {
		obj.watcher.RemoveWatch(dir)
	}

	obj.watcher.RemoveWatch(PERSON_SYS_GREETER_PATH)
	dir = path.Join(homeDir, PERSON_LOCAL_GREETER_PATH)
	if dutils.IsFileExist(dir) {
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
				if obj.watcher != nil {
					obj.endWatch()
				}
				obj.startWatch()
				break
			}

			if ev == nil {
				break
			}

			Logger.Debugf("Manager Event: %v", ev)
			ok1 := false
			ok2 := false

			ok1, _ = regexp.MatchString(ICON_SYS_PATH, ev.Name)
			ok2, _ = regexp.MatchString(ICON_LOCAL_PATH, ev.Name)
			if ok1 || ok2 {
				Logger.Debugf("Update IconList")
				obj.setPropIconThemeList(obj.getIconStrList())
				obj.setPropCursorThemeList(obj.getCursorStrList())
				break
			}

			ok1, _ = regexp.MatchString(SOUND_THEME_PATH, ev.Name)
			if ok1 {
				Logger.Debugf("Update SoundTheme")
				obj.setPropSoundThemeList(obj.getSoundStrList())
				break
			}

			ok1, _ = regexp.MatchString(`greeter-theme`, ev.Name)
			if ok1 {
				obj.setPropGreeterList(obj.getGreeterStrList())
				break
			}

			ok1, _ = regexp.MatchString(PERSON_SYS_THEME_PATH, ev.Name)
			ok2, _ = regexp.MatchString(PERSON_LOCAL_THEME_PATH, ev.Name)
			if ok1 || ok2 {
				obj.rebuildThemes()
				obj.setPropThemeList(obj.getDThemeStrList())
				break
			}

			ok1, _ = regexp.MatchString(THEME_SYS_PATH, ev.Name)
			ok2, _ = regexp.MatchString(THEME_LOCAL_PATH, ev.Name)
			if ok1 || ok2 {
				Logger.Debugf("Update GtkList")
				obj.setPropGtkThemeList(obj.getGtkStrList())
				break
			}
		case err, ok := <-obj.watcher.Error:
			if !ok || err != nil {
				if obj.watcher != nil {
					obj.endWatch()
				}
				obj.startWatch()
			}
		}
	}
}

func (obj *Manager) startBgWatch() {
	if obj.bgWatcher == nil {
		var err error
		obj.bgWatcher, err = fsnotify.NewWatcher()
		if err != nil {
			Logger.Fatalf("New Watcher Failed: %v", err)
		}
	}

	obj.bgWatcher.Watch(DEFAULT_SYS_BG_DIR)
	obj.bgWatcher.Watch("/usr/share/personalization/thumbnail/autogen")

	pict := getUserPictureDir()
	userBG := path.Join(pict, "Wallpapers")
	Logger.Debugf("User Special Bg: %v", userBG)
	if !dutils.IsFileExist(userBG) {
		if err := os.MkdirAll(userBG, 0755); err != nil {
			return
		}
	}
	obj.bgWatcher.Watch(userBG)
}

func (obj *Manager) endBgWatch() {
	if obj.bgWatcher == nil {
		return
	}

	pict := getUserPictureDir()
	userBG := path.Join(pict, "Wallpapers")
	if dutils.IsFileExist(userBG) {
		obj.bgWatcher.RemoveWatch(userBG)
	}

	obj.bgWatcher.RemoveWatch(DEFAULT_SYS_BG_DIR)
	obj.bgWatcher.RemoveWatch("/usr/share/personalization/thumbnail/autogen")
}

func (obj *Manager) handleBgEvent() {
	preTimestamp := int64(0)
	for {
		select {
		case <-obj.bgQuitFlag:
			return
		case ev, ok := <-obj.bgWatcher.Event:
			if !ok {
				if obj.bgWatcher != nil {
					obj.endBgWatch()
				}
				obj.startBgWatch()
				break
			}

			if ev == nil {
				break
			}

			Logger.Debugf("Bg Event: %v", ev)
			curTimestamp := time.Now().Unix()
			if curTimestamp-preTimestamp <= 3 {
				break
			}
			preTimestamp = curTimestamp
			if ok, _ := regexp.MatchString(`(autogen)(\.png$)`, ev.Name); ok {
				go func() {
					<-time.After(time.Second * 3)
					obj.setPropGtkThemeList(obj.getGtkStrList())
					obj.setPropIconThemeList(obj.getIconStrList())
					obj.setPropBackgroundList(obj.getBgStrList())
					return
				}()
				break
			}

			obj.setPropBackgroundList(obj.getBgStrList())
		case err, ok := <-obj.bgWatcher.Error:
			if !ok || err != nil {
				if obj.bgWatcher != nil {
					obj.endBgWatch()
				}
				obj.startBgWatch()
			}
		}
	}
}

func (obj *Theme) startWatch() {
	if obj.Type == THEME_TYPE_SYSTEM {
		return
	}

	if obj.watcher == nil {
		var err error
		obj.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			Logger.Warningf("New Watcher Failed: %v", err)
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
				if obj.watcher != nil {
					obj.endWatch()
				}
				obj.startWatch()
				break
			}

			if ok, _ := regexp.MatchString(`\.swa?px?$`, ev.Name); ok {
				break
			}

			if ev.IsDelete() {
				if obj.watcher != nil {
					obj.endWatch()
				}
				obj.startWatch()
				break
			}

			if ev.IsModify() {
				obj.setAllProps()
				if GetManager().CurrentTheme.GetValue().(string) == obj.Name {
					obj.setAllThemes()
				}
			}
		case err, ok := <-obj.watcher.Error:
			if !ok || err != nil {
				if obj.watcher != nil {
					obj.endWatch()
				}
				obj.startWatch()
			}
		}
	}
}
