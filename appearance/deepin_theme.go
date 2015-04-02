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
	"os"
	"path"
	"pkg.linuxdeepin.com/dde-daemon/appearance/fonts"
	"pkg.linuxdeepin.com/lib/glib-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	systemThemePath = "/usr/share/personalization/themes"
	userThemePath   = ".local/share/personalization/themes"

	customDThemeId = "Custom"

	tempDThemeCustom = "/usr/share/dde-daemon/template/theme_custom"
	dirPermMode      = 0755
)

func (m *Manager) setTheme(themeType string, value interface{}) {
	dinfo := m.newThemeInfo(themeType, value)
	if dinfo == nil {
		return
	}
	if theme, ok := m.isThemeExist(dinfo); ok {
		m.Set("theme", theme)
		return
	}

	if m.modifyCustomTheme(dinfo) {
		m.settings.SetString(deepinGSKeyTheme, customDThemeId)
	}
}

func (m *Manager) applyTheme(name string) {
	if name != m.CurrentTheme.Get() {
		return
	}

	t, ok := m.themeObjMap[name]
	if !ok {
		return
	}

	m.gtk.Set(t.GtkTheme)
	m.icon.Set(t.IconTheme)
	m.cursor.Set(t.CursorTheme)
	m.bg.Set(t.Background)
	m.font.Set(fonts.FontTypeStandard, t.FontName, t.FontSize)
	m.font.Set(fonts.FontTypeMonospaced, t.FontMono, t.FontSize)
	m.sound.Set(t.SoundTheme)
	m.greeter.Set(m.settings.GetString(deepinGSKeyGreeter))
}

func (m *Manager) modifyCustomTheme(info *Theme) bool {
	homeDir := dutils.GetHomeDir()
	dir := path.Join(homeDir, userThemePath, "Custom")
	if !dutils.IsFileExist(dir) {
		err := os.MkdirAll(dir, dirPermMode)
		if err != nil {
			logger.Debug(err)
			return false
		}
	}

	filename := path.Join(dir, "theme.ini")
	var newFile bool
	if !dutils.IsFileExist(filename) {
		err := dutils.CopyFile(tempDThemeCustom, filename)
		if err != nil {
			logger.Debug(err)
			return false
		}
		newFile = true
	}

	kFile := glib.NewKeyFile()
	defer kFile.Free()
	_, err := kFile.LoadFromFile(filename,
		glib.KeyFileFlagsKeepComments|
			glib.KeyFileFlagsKeepTranslations)
	if err != nil {
		return false
	}

	kFile.SetString(groupKeyComponent, themeKeyGtk, info.GtkTheme)
	kFile.SetString(groupKeyComponent, themeKeyIcon, info.IconTheme)
	kFile.SetString(groupKeyComponent, themeKeyCursor, info.CursorTheme)
	kFile.SetString(groupKeyComponent, themeKeySound, info.SoundTheme)
	kFile.SetString(groupKeyComponent, themeKeyBackground, info.Background)
	kFile.SetString(groupKeyComponent, themeKeyFontName, info.FontName)
	kFile.SetString(groupKeyComponent, themeKeyFontMono, info.FontMono)
	kFile.SetInteger(groupKeyComponent, themeKeyFontSize, info.FontSize)

	_, contents, err := kFile.ToData()
	if err != nil {
		return false
	}

	m.wLocker.Lock()
	defer m.wLocker.Unlock()
	ok := dutils.WriteStringToKeyFile(filename, contents)
	if newFile && ok {
		touchFile()
	}

	return ok
}

func (m *Manager) isThemeExist(info *Theme) (string, bool) {
	for n, t := range m.themeObjMap {
		if isThemeInfoSame(info, t) {
			return n, true
		}
	}

	return "", false
}

func (m *Manager) newThemeInfo(themeType string, value interface{}) *Theme {
	t, ok := m.themeObjMap[m.CurrentTheme.Get()]
	if !ok {
		logger.Debug("Get Current Theme Object Failed")
		return nil
	}

	switch themeType {
	case "gtk":
		return &Theme{
			GtkTheme:    value.(string),
			IconTheme:   t.IconTheme,
			CursorTheme: t.CursorTheme,
			SoundTheme:  t.SoundTheme,
			Background:  t.Background,
			FontName:    t.FontName,
			FontMono:    t.FontMono,
			FontSize:    t.FontSize,
		}
	case "icon":
		return &Theme{
			GtkTheme:    t.GtkTheme,
			IconTheme:   value.(string),
			CursorTheme: t.CursorTheme,
			SoundTheme:  t.SoundTheme,
			Background:  t.Background,
			FontName:    t.FontName,
			FontMono:    t.FontMono,
			FontSize:    t.FontSize,
		}
	case "cursor":
		return &Theme{
			GtkTheme:    t.GtkTheme,
			IconTheme:   t.IconTheme,
			CursorTheme: value.(string),
			SoundTheme:  t.SoundTheme,
			Background:  t.Background,
			FontName:    t.FontName,
			FontMono:    t.FontMono,
			FontSize:    t.FontSize,
		}
	case "sound":
		return &Theme{
			GtkTheme:    t.GtkTheme,
			IconTheme:   t.IconTheme,
			CursorTheme: t.CursorTheme,
			SoundTheme:  value.(string),
			Background:  t.Background,
			FontName:    t.FontName,
			FontMono:    t.FontMono,
			FontSize:    t.FontSize,
		}
	case "background":
		return &Theme{
			GtkTheme:    t.GtkTheme,
			IconTheme:   t.IconTheme,
			CursorTheme: t.CursorTheme,
			SoundTheme:  t.SoundTheme,
			Background:  value.(string),
			FontName:    t.FontName,
			FontMono:    t.FontMono,
			FontSize:    t.FontSize,
		}
	case "font-standard":
		return &Theme{
			GtkTheme:    t.GtkTheme,
			IconTheme:   t.IconTheme,
			CursorTheme: t.CursorTheme,
			SoundTheme:  t.SoundTheme,
			Background:  t.Background,
			FontName:    value.(string),
			FontMono:    t.FontMono,
			FontSize:    t.FontSize,
		}
	case "font-mono":
		return &Theme{
			GtkTheme:    t.GtkTheme,
			IconTheme:   t.IconTheme,
			CursorTheme: t.CursorTheme,
			SoundTheme:  t.SoundTheme,
			Background:  t.Background,
			FontName:    t.FontName,
			FontMono:    value.(string),
			FontSize:    t.FontSize,
		}
	case "font-size":
		return &Theme{
			GtkTheme:    t.GtkTheme,
			IconTheme:   t.IconTheme,
			CursorTheme: t.CursorTheme,
			SoundTheme:  t.SoundTheme,
			Background:  t.Background,
			FontName:    t.FontName,
			FontMono:    t.FontMono,
			FontSize:    value.(int32),
		}
	}

	return nil
}

func getThemeObjectList(list []string) []string {
	list = sortNameByDeepin(list)
	var tmp []string
	for _, name := range list {
		name = path.Base(name)
		tmp = append(tmp, themeDBusPath+name)
	}

	return tmp
}

func touchFile() {
	dir := path.Join(os.Getenv("HOME"), userThemePath)
	filename := path.Join(dir, "__emit-signal__")

	if !dutils.IsFileExist(filename) {
		os.MkdirAll(filename, dirPermMode)
	} else {
		os.RemoveAll(filename)
	}
}
