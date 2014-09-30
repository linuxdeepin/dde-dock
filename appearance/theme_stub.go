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
	"pkg.linuxdeepin.com/lib/glib-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
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

func (t *Theme) readFromFile() {
	kFile := glib.NewKeyFile()
	defer kFile.Free()

	_, err := kFile.LoadFromFile(path.Join(t.filePath, "theme.ini"),
		glib.KeyFileFlagsKeepComments|
			glib.KeyFileFlagsKeepTranslations)
	if err != nil {
		logger.Warningf("Load KeyFile '%s' Failed: %v",
			t.filePath, err)
		return
	}

	str, err := kFile.GetString(groupKeyTheme, themeKeyId)
	if err == nil {
		t.setPropString(&t.Name, "Name", str)
	}

	str, err = kFile.GetLocaleString(groupKeyTheme,
		themeKeyName, "\x00")
	if err == nil {
		t.setPropString(&t.DisplayName, "DisplayName", str)
	}

	str, _ = kFile.GetString(groupKeyComponent, themeKeyGtk)
	if len(str) == 0 {
		str = "Deepin"
	}
	t.setPropString(&t.GtkTheme, "GtkTheme", str)

	str, _ = kFile.GetString(groupKeyComponent, themeKeyIcon)
	if len(str) == 0 {
		str = "Deepin"
	}
	t.setPropString(&t.IconTheme, "IconTheme", str)

	str, _ = kFile.GetString(groupKeyComponent, themeKeySound)
	if len(str) == 0 {
		str = "LinuxDeepin"
	}
	t.setPropString(&t.SoundTheme, "SoundTheme", str)

	str, _ = kFile.GetString(groupKeyComponent, themeKeyCursor)
	if len(str) == 0 {
		str = "Deepin-Cursor"
	}
	t.setPropString(&t.CursorTheme, "CursorTheme", str)

	str, _ = kFile.GetString(groupKeyComponent, themeKeyFontName)
	if len(str) == 0 {
		str = "WenQuanYi Micro Hei"
	}
	t.setPropString(&t.FontName, "FontName", str)

	str, _ = kFile.GetString(groupKeyComponent, themeKeyFontMono)
	if len(str) == 0 {
		str = "WenQuanYi Micro Hei Mono"
	}
	t.setPropString(&t.FontMono, "FontMono", str)

	str, _ = kFile.GetString(groupKeyComponent, themeKeyBackground)
	if len(str) == 0 {
		str = "/usr/share/backgrounds/default_background.jpg"
	}
	str = dutils.EncodeURI(str, dutils.SCHEME_FILE)
	t.setPropString(&t.Background, "Background", str)

	var interval int
	interval, _ = kFile.GetInteger(groupKeyComponent, themeKeyFontSize)
	if interval == 0 {
		interval = 10
	}
	t.setPropFontSize(int32(interval))
}
