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
)

type deepinThemeInfo struct {
	GtkTheme    string
	IconTheme   string
	CursorTheme string
	SoundTheme  string
	Background  string
	FontName    string
	FontMono    string
	FontSize    int32
}

func (m *Manager) setTheme(themeType string, value interface{}) {
	dinfo := m.newDeepinThemeInfo(themeType, value)
	if dinfo == nil {
		return
	}
	if theme, ok := m.isThemeExist(dinfo); ok {
		m.Set("theme", theme)
		return
	}

	m.modifyCustomTheme(dinfo)
	m.Set("theme", customDThemeId)
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

func (m *Manager) modifyCustomTheme(info *deepinThemeInfo) bool {
	homeDir := dutils.GetHomeDir()
	dir := path.Join(homeDir, userThemePath, "Custom")
	if !dutils.IsFileExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Debug(err)
			return false
		}
	}

	filename := path.Join(dir, "theme.ini")
	if !dutils.IsFileExist(filename) {
		if err := dutils.CopyFile(tempDThemeCustom, filename); err != nil {
			logger.Debug(err)
			return false
		}
	}

	kFile := glib.NewKeyFile()
	defer kFile.Free()
	if _, err := kFile.LoadFromFile(filename, glib.KeyFileFlagsKeepComments|
		glib.KeyFileFlagsKeepTranslations); err != nil {
		return false
	}

	kFile.SetString(groupKeyComponent, themeKeyGtk, info.GtkTheme)
	kFile.SetString(groupKeyComponent, themeKeyIcon, info.IconTheme)
	kFile.SetString(groupKeyComponent, themeKeyCursor, info.CursorTheme)
	kFile.SetString(groupKeyComponent, themeKeySound, info.SoundTheme)
	kFile.SetString(groupKeyComponent, themeKeyBackground, info.Background)
	kFile.SetString(groupKeyComponent, themeKeyFontName, info.FontName)
	kFile.SetString(groupKeyComponent, themeKeyFontMono, info.FontMono)
	kFile.SetInteger(groupKeyComponent, themeKeyFontSize, int(info.FontSize))

	_, contents, err := kFile.ToData()
	if err != nil {
		return false
	}
	return dutils.WriteStringToKeyFile(filename, contents)
}

func (m *Manager) isThemeExist(info *deepinThemeInfo) (string, bool) {
	for n, t := range m.themeObjMap {
		if t.GtkTheme == info.GtkTheme &&
			t.IconTheme == info.IconTheme &&
			t.CursorTheme == info.CursorTheme &&
			t.SoundTheme == info.SoundTheme &&
			t.Background == info.Background &&
			t.FontName == info.FontName &&
			t.FontMono == info.FontMono &&
			t.FontSize == info.FontSize {
			return n, true
		}
	}

	return "", false
}

func (m *Manager) newDeepinThemeInfo(themeType string, value interface{}) *deepinThemeInfo {
	t, ok := m.themeObjMap[m.CurrentTheme.Get()]
	if !ok {
		logger.Debug("Get Current Theme Object Failed")
		return nil
	}

	switch themeType {
	case "gtk":
		return &deepinThemeInfo{
			value.(string),
			t.IconTheme,
			t.CursorTheme,
			t.SoundTheme,
			t.Background,
			t.FontName,
			t.FontMono,
			t.FontSize,
		}
	case "icon":
		return &deepinThemeInfo{
			t.GtkTheme,
			value.(string),
			t.CursorTheme,
			t.SoundTheme,
			t.Background,
			t.FontName,
			t.FontMono,
			t.FontSize,
		}
	case "cursor":
		return &deepinThemeInfo{
			t.GtkTheme,
			t.IconTheme,
			value.(string),
			t.SoundTheme,
			t.Background,
			t.FontName,
			t.FontMono,
			t.FontSize,
		}
	case "sound":
		return &deepinThemeInfo{
			t.GtkTheme,
			t.IconTheme,
			t.CursorTheme,
			value.(string),
			t.Background,
			t.FontName,
			t.FontMono,
			t.FontSize,
		}
	case "background":
		return &deepinThemeInfo{
			t.GtkTheme,
			t.IconTheme,
			t.CursorTheme,
			t.SoundTheme,
			value.(string),
			t.FontName,
			t.FontMono,
			t.FontSize,
		}
	case "font-name":
		return &deepinThemeInfo{
			t.GtkTheme,
			t.IconTheme,
			t.CursorTheme,
			t.SoundTheme,
			t.Background,
			value.(string),
			t.FontMono,
			t.FontSize,
		}
	case "font-mono":
		return &deepinThemeInfo{
			t.GtkTheme,
			t.IconTheme,
			t.CursorTheme,
			t.SoundTheme,
			t.Background,
			t.FontName,
			value.(string),
			t.FontSize,
		}
	case "font-size":
		return &deepinThemeInfo{
			t.GtkTheme,
			t.IconTheme,
			t.CursorTheme,
			t.SoundTheme,
			t.Background,
			t.FontName,
			t.FontMono,
			value.(int32),
		}
	}

	return nil
}

func getThemeObjectList(list []string) []string {
	list = sortNameByDeepin(list)
	var tmp []string
	for _, name := range list {
		tmp = append(tmp, themeDBusPath+name)
	}

	return tmp
}
