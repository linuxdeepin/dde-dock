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
)

func (obj *Manager) GetObjectPath(name string) string {
	for _, l := range obj.ThemeList {
		if name == l.Name {
			return THEME_PATH + name
		}
	}

	return ""
}

func (obj *Manager) SetGtkTheme(theme string) {
	for _, l := range obj.GtkThemeList {
		if theme == l.Name {
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

func (obj *Manager) SetIconTheme(theme string) {
	for _, l := range obj.IconThemeList {
		if theme == l.Name {
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

func (obj *Manager) SetCursorTheme(theme string) {
	for _, l := range obj.CursorThemeList {
		if theme == l.Name {
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

func (obj *Manager) SetSoundTheme(theme string) {
	for _, l := range obj.SoundThemeList {
		if theme == l.Name {
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

func (obj *Manager) SetBackground(bg string) {
	flag := false
	uri, _ := objUtil.PathToFileURI(bg)
	for _, l := range obj.BackgroundList {
		if uri == l.Name {
			flag = true
		}
	}

	// bg not in list
	if !flag {
		bg, _ = objUtil.URIToPath(bg)
		if !obj.AppendBackground(bg) {
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

func (obj *Manager) SetFontSize(size int32) {
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

func (obj *Manager) DeleteGtkTheme(theme string) {
	filepath := ""
	flag := false
	for _, l := range obj.GtkThemeList {
		if l.Name == theme {
			if l.T == THEME_TYPE_LOCAL {
				flag = true
				filepath = l.Path
			}
			break
		}
	}

	if flag {
		os.RemoveAll(filepath)
	}
}

func (obj *Manager) DeleteIconTheme(theme string) {
	filepath := ""
	flag := false
	for _, l := range obj.IconThemeList {
		if l.Name == theme {
			if l.T == THEME_TYPE_LOCAL {
				flag = true
				filepath = l.Path
			}
			break
		}
	}

	if flag {
		os.RemoveAll(filepath)
	}
}

func (obj *Manager) DeleteSoundTheme(theme string) {
	filepath := ""
	flag := false
	for _, l := range obj.SoundThemeList {
		if l.Name == theme {
			if l.T == THEME_TYPE_LOCAL {
				flag = true
				filepath = l.Path
			}
			break
		}
	}

	if flag {
		os.RemoveAll(filepath)
	}
}

func (obj *Manager) DeleteCursorTheme(theme string) {
	filepath := ""
	flag := false
	for _, l := range obj.CursorThemeList {
		if l.Name == theme {
			if l.T == THEME_TYPE_LOCAL {
				flag = true
				filepath = l.Path
			}
			break
		}
	}

	if flag {
		os.RemoveAll(filepath)
	}
}

func (obj *Manager) AppendBackground(bg string) bool {
	tmp := BgInfo{}
	tmp.Name = path.Base(bg)
	tmp.Path = bg
	tmp.T = THEME_TYPE_LOCAL

	if isBgInfoInList(tmp, obj.BackgroundList) {
		return true
	}

	pict := getUserPictureDir()
	pict = path.Join(pict, "Wallpapers")
	if !objUtil.IsFileExist(pict) {
		os.MkdirAll(pict, 0755)
	}
	destPath := path.Join(pict, tmp.Name)
	if !objUtil.CopyFile(bg, destPath) {
		return false
	}

	return true
}

func (obj *Manager) DeleteBackground(bg string) {
	flag := false
	//tmpList := []BgInfo{}
	for _, l := range obj.BackgroundList {
		if isBackgroundSame(bg, l.Path) {
			if l.T == THEME_TYPE_LOCAL {
				flag = true
			}
			break
			//continue
		}
		//tmpList = append(tmpList, l)
	}

	if flag {
		os.RemoveAll(bg)
		//obj.setPropBackgroundList(tmpList)
	}
}
