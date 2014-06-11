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
	"strings"
)

func (op *Manager) GetFlag(t, name string) int32 {
	t = strings.ToLower(t)
	switch t {
	case "theme":
	case "gtk":
	case "icon":
	case "sound":
	case "cursor":
	case "background":
	}

	return 0
}

func (op *Manager) Set(t, name, value string) {
	t = strings.ToLower(t)
	switch t {
	case "theme":
	case "gtk":
	case "icon":
	case "sound":
	case "cursor":
	case "fontsize":
	case "background":
	}
}

func (op *Manager) GetThumbnail(t, name string) string {
	t = strings.ToLower(t)
	switch t {
	case "theme":
	case "gtk":
	case "icon":
	case "sound":
	case "cursor":
	case "background":
	}

	return "/usr/share/personalization/thumbnail/IconThemes/Deepin/thumbnail.png"
}

func (obj *Manager) Delete(t, name string) {
	t = strings.ToLower(t)
	switch t {
	case "theme":
	case "gtk":
	case "icon":
	case "sound":
	case "cursor":
	case "background":
	}
}

func (obj *Manager) setGtkTheme(theme string) {
	for _, l := range obj.GtkThemeList {
		if theme == l {
			t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
			if !ok {
				obj.setPropCurrentTheme(DEFAULT_THEME)
				break
			}
			if t.GtkTheme == theme {
				return
			}
			name, ok1 := obj.isThemeExit(theme, t.IconTheme, t.SoundTheme,
				t.CursorTheme, t.Background, t.FontSize)
			if ok1 {
				obj.setPropCurrentTheme(name)
				break
			}

			obj.modifyTheme("Custom", theme, t.IconTheme, t.SoundTheme,
				t.CursorTheme, t.Background, t.FontSize)
			obj.setPropCurrentTheme("Custom")
			break
		}
	}
}

func (obj *Manager) setIconTheme(theme string) {
	for _, l := range obj.IconThemeList {
		if theme == l {
			t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
			if !ok {
				obj.setPropCurrentTheme(DEFAULT_THEME)
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

			obj.modifyTheme("Custom", t.GtkTheme, theme, t.SoundTheme,
				t.CursorTheme, t.Background, t.FontSize)
			obj.setPropCurrentTheme("Custom")
			break
		}
	}
}

func (obj *Manager) setCursorTheme(theme string) {
	for _, l := range obj.CursorThemeList {
		if theme == l {
			t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
			if !ok {
				obj.setPropCurrentTheme(DEFAULT_THEME)
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

			obj.modifyTheme("Custom", t.GtkTheme, t.IconTheme, t.SoundTheme,
				theme, t.Background, t.FontSize)
			obj.setPropCurrentTheme("Custom")
			break
		}
	}
}

func (obj *Manager) setSoundTheme(theme string) {
	for _, l := range obj.SoundThemeList {
		if theme == l {
			t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
			if !ok {
				obj.setPropCurrentTheme(DEFAULT_THEME)
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

			obj.modifyTheme("Custom", t.GtkTheme, t.IconTheme, theme,
				t.CursorTheme, t.Background, t.FontSize)
			obj.setPropCurrentTheme("Custom")
			break
		}
	}
}

func (obj *Manager) setBackground(bg string) {
	flag := false
	uri, _ := objUtil.PathToFileURI(bg)
	for _, l := range obj.BackgroundList {
		if uri == l {
			flag = true
			break
		}
	}

	// bg not in list
	if !flag {
		bg, _ = objUtil.URIToPath(bg)
		if !obj.appendBackground(bg) {
			return
		}
	}

	t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
	if !ok {
		obj.setPropCurrentTheme(DEFAULT_THEME)
		return
	}
	if t.Background == bg {
		return
	}
	name, ok1 := obj.isThemeExit(t.GtkTheme, t.IconTheme, t.SoundTheme,
		t.CursorTheme, bg, t.FontSize)
	if ok1 {
		obj.setPropCurrentTheme(name)
		return
	}

	obj.modifyTheme("Custom", t.GtkTheme, t.IconTheme, t.SoundTheme,
		t.CursorTheme, bg, t.FontSize)
	obj.setPropCurrentTheme("Custom")
}

func (obj *Manager) setFontSize(size int32) {
	if size > 20 || size < 7 {
		return
	}

	t, ok := obj.themeObjMap[obj.CurrentTheme.GetValue().(string)]
	if !ok {
		obj.setPropCurrentTheme(DEFAULT_THEME)
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

	obj.modifyTheme("Custom", t.GtkTheme, t.IconTheme, t.SoundTheme,
		t.CursorTheme, t.Background, size)
	obj.setPropCurrentTheme("Custom")
}

func (obj *Manager) deleteGtkTheme(theme string) {
	filepath := ""
	flag := false
	for _, l := range obj.GtkThemeList {
		if l == theme {
			if obj.GetFlag("gtk", l) == int32(THEME_TYPE_LOCAL) {
				flag = true
				//filepath = l.Path
			}
			break
		}
	}

	if flag {
		os.RemoveAll(filepath)
	}
}

func (obj *Manager) deleteIconTheme(theme string) {
	filepath := ""
	flag := false
	for _, l := range obj.IconThemeList {
		if l == theme {
			if obj.GetFlag("icon", l) == int32(THEME_TYPE_LOCAL) {
				flag = true
				//filepath = l.Path
			}
			break
		}
	}

	if flag {
		os.RemoveAll(filepath)
	}
}

func (obj *Manager) deleteSoundTheme(theme string) {
	filepath := ""
	flag := false
	for _, l := range obj.SoundThemeList {
		if l == theme {
			if obj.GetFlag("sound", l) == int32(THEME_TYPE_LOCAL) {
				flag = true
				//filepath = l.Path
			}
			break
		}
	}

	if flag {
		os.RemoveAll(filepath)
	}
}

func (obj *Manager) deleteCursorTheme(theme string) {
	filepath := ""
	flag := false
	for _, l := range obj.CursorThemeList {
		if l == theme {
			if obj.GetFlag("curosr", l) == int32(THEME_TYPE_LOCAL) {
				flag = true
				//filepath = l.Path
			}
			break
		}
	}

	if flag {
		os.RemoveAll(filepath)
	}
}

func (obj *Manager) appendBackground(bg string) bool {
	if isStrInList(bg, obj.BackgroundList) {
		return true
	}

	pict := getUserPictureDir()
	pict = path.Join(pict, "Wallpapers")
	if !objUtil.IsFileExist(pict) {
		os.MkdirAll(pict, 0755)
	}
	destPath := path.Join(pict, path.Base(bg))
	if !objUtil.CopyFile(bg, destPath) {
		return false
	}

	return true
}

func (obj *Manager) deleteBackground(bg string) {
	if !isStrInList(bg, obj.BackgroundList) {
		return
	}

	if obj.GetFlag("background", bg) == int32(THEME_TYPE_LOCAL) {
		os.RemoveAll(bg)
	}
}
