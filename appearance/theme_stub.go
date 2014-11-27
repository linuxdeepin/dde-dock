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
	"path"
	"pkg.linuxdeepin.com/lib/dbus"
)

const (
	themeDBusPath = "/com/deepin/daemon/Theme/"
	themeDBusIFC  = "com.deepin.daemon.Theme"
)

const (
	groupKeyTheme      = "Theme"
	groupKeyComponent  = "Component"
	themeKeyId         = "Id"
	themeKeyName       = "Name"
	themeKeyGtk        = "GtkTheme"
	themeKeyIcon       = "IconTheme"
	themeKeySound      = "SoundTheme"
	themeKeyCursor     = "CursorTheme"
	themeKeyBackground = "BackgroundFile"
	themeKeyFontName   = "FontName"
	themeKeyFontMono   = "FontMono"
	themeKeyFontSize   = "FontSize"
)

func (t *Theme) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: t.objectPath,
		Interface:  themeDBusIFC,
	}
}

func (t *Theme) setPropString(handle *string, name, value string) {
	if *handle == value {
		return
	}
	*handle = value
	dbus.NotifyChange(t, name)
}

func (t *Theme) setPropFontSize(s int32) {
	if t.FontSize != s {
		t.FontSize = s
		dbus.NotifyChange(t, "FontSize")
	}
}

func (t *Theme) setPropType(ty int32) {
	if t.Type != ty {
		t.Type = ty
		dbus.NotifyChange(t, "Type")
	}
}

func (t *Theme) setPropPreview(list []string) {
	t.Preview = list
	dbus.NotifyChange(t, "Preview")
}

func (t *Theme) setPropsFromFile() {
	filename := path.Join(t.filePath, "theme.ini")
	info, err := getThemeInfoFromFile(filename)
	if err != nil {
		return
	}

	t.setPropString(&t.Name, "Name", info.Name)
	if len(info.DisplayName) == 0 {
		info.DisplayName = info.Name
	}
	t.setPropString(&t.DisplayName, "DisplayName", info.DisplayName)

	if len(info.GtkTheme) == 0 {
		info.GtkTheme = "Deepin"
	}
	t.setPropString(&t.GtkTheme, "GtkTheme", info.GtkTheme)

	if len(info.IconTheme) == 0 {
		info.IconTheme = "Deepin"
	}
	t.setPropString(&t.IconTheme, "IconTheme", info.IconTheme)

	if len(info.SoundTheme) == 0 {
		info.SoundTheme = "LinuxDeepin"
	}
	t.setPropString(&t.SoundTheme, "SoundTheme", info.SoundTheme)

	if len(info.CursorTheme) == 0 {
		info.CursorTheme = "Deepin-Cursor"
	}
	t.setPropString(&t.CursorTheme, "CursorTheme", info.CursorTheme)

	if len(info.FontName) == 0 {
		info.FontName = "Source Han Sans SC"
	}
	t.setPropString(&t.FontName, "FontName", info.FontName)

	if len(info.FontMono) == 0 {
		info.FontMono = "DejaVu Sans Mono"
	}
	t.setPropString(&t.FontMono, "FontMono", info.FontMono)

	if len(info.Background) == 0 {
		info.Background = "file:///usr/share/backgrounds/default_background.jpg"
	}
	t.setPropString(&t.Background, "Background", info.Background)

	if info.FontSize == 0 {
		info.FontSize = 10
	}
	t.setPropFontSize(int32(info.FontSize))
}
