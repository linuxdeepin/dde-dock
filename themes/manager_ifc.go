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

package themes

import (
	"os"
	"path"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"strconv"
	"strings"
)

func (obj *Manager) GetFlag(t, name string) int32 {
	t = strings.ToLower(t)
	switch t {
	case "theme":
		list := getDThemeList()
		for _, l := range list {
			if name == l.Name {
				return l.T
			}
		}
	case "gtk":
		list := getGtkThemeList()
		for _, l := range list {
			if name == l.Name {
				return l.T
			}
		}
	case "icon":
		list := getIconThemeList()
		for _, l := range list {
			if name == l.Name {
				return l.T
			}
		}
	case "sound":
		list := getSoundThemeList()
		for _, l := range list {
			if name == l.Name {
				return l.T
			}
		}
	case "cursor":
		list := getCursorThemeList()
		for _, l := range list {
			if name == l.Name {
				return l.T
			}
		}
	case "background":
		list := getBackgroundList()
		for _, l := range list {
			if name == l.Path {
				return l.T
			}
		}
	case "greeter":
		list := getGreeterThemeList()
		for _, l := range list {
			if name == l.Path {
				return l.T
			}
		}
	}

	return -1
}

func (obj *Manager) Set(t, value string) {
	t = strings.ToLower(t)
	switch t {
	case "theme":
		obj.setPropCurrentTheme(value)
	case "gtk":
		obj.setGtkTheme(value)
	case "icon":
		obj.setIconTheme(value)
	case "sound":
		obj.setSoundTheme(value)
	case "cursor":
		obj.setCursorTheme(value)
	case "fontsize":
		s, _ := strconv.ParseInt(value, 10, 64)
		obj.setFontSize(int32(s))
	case "background":
		logger.Info("SET - set bg:", value)
		obj.setBackground(value)
	case "greeter":
		if value == obj.GreeterTheme.GetValue().(string) {
			break
		}

		for _, n := range obj.GreeterThemeList {
			if value == n {
				themeSettings.SetString(GS_KEY_CURRENT_GREETER, value)
				break
			}
		}
	}
}

func (op *Manager) GetThumbnail(t, name string) string {
	t = strings.ToLower(t)
	switch t {
	case "theme":
		return getDThemeThumb(name)
	case "gtk":
		return getGtkThumb(name)
	case "icon":
		return getIconThumb(name)
	case "cursor":
		return getCursorThumb(name)
	case "background":
		return getBgThumb(name)
	case "greeter":
		return getGreeterThumb(name)
	}

	return ""
}

func (obj *Manager) Delete(t, name string) {
	t = strings.ToLower(t)
	switch t {
	case "theme":
		obj.deleteDTheme(name)
	case "gtk":
		obj.deleteGtkTheme(name)
	case "icon":
		obj.deleteIconTheme(name)
	case "sound":
		obj.deleteSoundTheme(name)
	case "cursor":
		obj.deleteCursorTheme(name)
	case "background":
		obj.deleteBackground(name)
	case "greeter":
		list := getGreeterThemeList()
		for _, l := range list {
			if name == l.Name {
				if l.T == THEME_TYPE_LOCAL {
					rmAllFile(l.Path)
				}
			}
			break
		}
	}
}

func (obj *Manager) setGtkTheme(theme string) {
	if len(theme) < 1 {
		return
	}

	for _, l := range obj.GtkThemeList {
		if theme == l {
			t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
			if !ok {
				logger.Warning("Current Theme Invalid:", theme)
				obj.setPropCurrentTheme(DEFAULT_THEME_ID)
				break
			}
			if t.GtkTheme == theme {
				logger.Warning("Gtk Theme Same:", theme)
				return
			}
			name, ok1 := obj.isThemeExit(theme, t.IconTheme,
				t.SoundTheme, t.CursorTheme,
				t.Background, t.FontSize)
			if ok1 {
				logger.Warning("Exist Theme:", name)
				obj.setPropCurrentTheme(name)
				break
			}

			logger.Info("Theme: ", theme, t.IconTheme, t.SoundTheme, t.CursorTheme, t.Background, t.FontSize)
			obj.newCustomTheme(theme, t.IconTheme, t.SoundTheme,
				t.CursorTheme, t.Background, t.FontSize)
			obj.setPropCurrentTheme(CUSTOM_THEME_ID)
			break
		}
	}
}

func (obj *Manager) setIconTheme(theme string) {
	if len(theme) < 1 {
		return
	}

	for _, l := range obj.IconThemeList {
		if theme == l {
			t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
			if !ok {
				obj.setPropCurrentTheme(DEFAULT_THEME_ID)
				break
			}
			if t.IconTheme == theme {
				return
			}
			name, ok1 := obj.isThemeExit(t.GtkTheme, theme, t.SoundTheme,
				t.CursorTheme, t.Background, t.FontSize)
			if ok1 {
				obj.setPropCurrentTheme(name)
				break
			}

			obj.newCustomTheme(t.GtkTheme, theme, t.SoundTheme,
				t.CursorTheme, t.Background, t.FontSize)
			obj.setPropCurrentTheme(CUSTOM_THEME_ID)
			break
		}
	}
}

func (obj *Manager) setCursorTheme(theme string) {
	if len(theme) < 1 {
		return
	}

	for _, l := range obj.CursorThemeList {
		if theme == l {
			t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
			if !ok {
				obj.setPropCurrentTheme(DEFAULT_THEME_ID)
				break
			}
			if t.CursorTheme == theme {
				return
			}
			name, ok1 := obj.isThemeExit(t.GtkTheme, t.IconTheme, t.SoundTheme,
				theme, t.Background, t.FontSize)
			if ok1 {
				obj.setPropCurrentTheme(name)
				break
			}

			obj.newCustomTheme(t.GtkTheme, t.IconTheme, t.SoundTheme,
				theme, t.Background, t.FontSize)
			obj.setPropCurrentTheme(CUSTOM_THEME_ID)
			break
		}
	}
}

func (obj *Manager) setSoundTheme(theme string) {
	if len(theme) < 1 {
		return
	}

	for _, l := range obj.SoundThemeList {
		if theme == l {
			t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
			if !ok {
				obj.setPropCurrentTheme(DEFAULT_THEME_ID)
				break
			}
			if t.SoundTheme == theme {
				break
			}
			name, ok1 := obj.isThemeExit(t.GtkTheme, t.IconTheme, theme,
				t.CursorTheme, t.Background, t.FontSize)
			if ok1 {
				obj.setPropCurrentTheme(name)
				break
			}

			obj.newCustomTheme(t.GtkTheme, t.IconTheme, theme,
				t.CursorTheme, t.Background, t.FontSize)
			obj.setPropCurrentTheme(CUSTOM_THEME_ID)
			break
		}
	}
}

func (obj *Manager) setBackground(bg string) bool {
	if len(bg) < 1 {
		logger.Warning("setBackground invalid bg:", bg)
		return false
	}
	flag := false
	uri := dutils.EncodeURI(bg, dutils.SCHEME_FILE)
	for _, l := range obj.BackgroundList {
		if uri == l {
			flag = true
			break
		}
	}

	// bg not in list
	if !flag && uri != DEFAULT_BG_URI {
		bg = dutils.DecodeURI(bg)
		if !obj.appendBackground(bg) {
			logger.Warning("Append bg failed:", bg)
			return false
		}
	}

	t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
	if !ok {
		obj.setPropCurrentTheme(DEFAULT_THEME_ID)
		logger.Warning("Get current themObj failed. Set current themObj to default")
		return true
	}
	if t.Background == uri {
		logger.Warning("Bg is same:", bg)
		return true
	}
	name, ok1 := obj.isThemeExit(t.GtkTheme, t.IconTheme, t.SoundTheme,
		t.CursorTheme, bg, t.FontSize)
	if ok1 {
		obj.setPropCurrentTheme(name)
		return true
	}

	if !obj.newCustomTheme(t.GtkTheme, t.IconTheme, t.SoundTheme,
		t.CursorTheme, bg, t.FontSize) {
		logger.Warning("New custom theme failed")
		return false
	}
	obj.setPropCurrentTheme(CUSTOM_THEME_ID)

	return true
}

func (obj *Manager) setFontSize(size int32) {
	if size > 20 || size < 7 {
		return
	}

	t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
	if !ok {
		obj.setPropCurrentTheme(DEFAULT_THEME_ID)
		return
	}
	if t.FontSize == size {
		return
	}
	name, ok1 := obj.isThemeExit(t.GtkTheme, t.IconTheme, t.SoundTheme,
		t.CursorTheme, t.Background, size)
	if ok1 {
		obj.setPropCurrentTheme(name)
		return
	}

	obj.newCustomTheme(t.GtkTheme, t.IconTheme, t.SoundTheme,
		t.CursorTheme, t.Background, size)
	obj.setPropCurrentTheme(CUSTOM_THEME_ID)
}

func (obj *Manager) deleteGtkTheme(theme string) {
	for _, l := range obj.GtkThemeList {
		if l == theme {
			if obj.GetFlag("gtk", l) == int32(THEME_TYPE_LOCAL) {
				list := getGtkThemeList()
				for _, t := range list {
					if theme == t.Name {
						rmAllFile(t.Path)
					}
				}
			}
			break
		}
	}
}

func (obj *Manager) deleteIconTheme(theme string) {
	for _, l := range obj.IconThemeList {
		if l == theme {
			if obj.GetFlag("icon", l) == int32(THEME_TYPE_LOCAL) {
				list := getIconThemeList()
				for _, t := range list {
					if theme == t.Name {
						rmAllFile(t.Path)
					}
				}
			}
			break
		}
	}
}

func (obj *Manager) deleteSoundTheme(theme string) {
	for _, l := range obj.SoundThemeList {
		if l == theme {
			if obj.GetFlag("sound", l) == int32(THEME_TYPE_LOCAL) {
				list := getSoundThemeList()
				for _, t := range list {
					if theme == t.Name {
						rmAllFile(t.Path)
					}
				}
			}
			break
		}
	}
}

func (obj *Manager) deleteCursorTheme(theme string) {
	for _, l := range obj.CursorThemeList {
		if l == theme {
			if obj.GetFlag("curosr", l) == int32(THEME_TYPE_LOCAL) {
				list := getCursorThemeList()
				for _, t := range list {
					if theme == t.Name {
						rmAllFile(t.Path)
					}
				}
			}
			break
		}
	}
}

func (obj *Manager) appendBackground(bg string) bool {
	if isStrInList(bg, obj.BackgroundList) {
		return true
	}

	pict := getUserPictureDir()
	pict = path.Join(pict, "Wallpapers")
	if !dutils.IsFileExist(pict) {
		os.MkdirAll(pict, 0755)
	}
	destPath := path.Join(pict, path.Base(bg))
	if err := dutils.CopyFile(bg, destPath); err != nil {
		return false
	}

	return true
}

func (obj *Manager) deleteBackground(bg string) {
	if !isStrInList(bg, obj.BackgroundList) {
		return
	}

	if obj.GetFlag("background", bg) == int32(THEME_TYPE_LOCAL) {
		rmAllFile(bg)
	}
}

func (obj *Manager) deleteDTheme(name string) {
	list := getDThemeList()
	for _, t := range list {
		if name == t.Name {
			if t.T == int32(THEME_TYPE_LOCAL) {
				if obj.CurrentTheme.GetValue().(string) == name {
					obj.setPropCurrentTheme(DEFAULT_THEME_ID)
				}
				rmAllFile(t.Path)
			}
			break
		}
	}
}
