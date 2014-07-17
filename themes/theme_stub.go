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
	"path"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/glib-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"regexp"
)

func (obj *Theme) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		MANAGER_DEST,
		obj.objectPath,
		THEME_IFC,
	}
}

func (obj *Theme) setPropName(name string) {
	if obj.Name != name {
		obj.Name = name
		dbus.NotifyChange(obj, "Name")
	}
}

func (obj *Theme) setPropGtkTheme(theme string) {
	if obj.GtkTheme != theme {
		obj.GtkTheme = theme
		dbus.NotifyChange(obj, "GtkTheme")
	}
}

func (obj *Theme) setPropIconTheme(theme string) {
	if obj.IconTheme != theme {
		obj.IconTheme = theme
		dbus.NotifyChange(obj, "IconTheme")
	}
}

func (obj *Theme) setPropSoundTheme(theme string) {
	if obj.SoundTheme != theme {
		obj.SoundTheme = theme
		dbus.NotifyChange(obj, "SoundTheme")
	}
}

func (obj *Theme) setPropCursorTheme(theme string) {
	if obj.CursorTheme != theme {
		obj.CursorTheme = theme
		dbus.NotifyChange(obj, "CursorTheme")
	}
}

func (obj *Theme) setPropBackground(bg string) {
	if obj.Background != bg {
		obj.Background = bg
		dbus.NotifyChange(obj, "Background")
	}
}

func (obj *Theme) setPropFontSize(s int32) {
	if obj.FontSize != s {
		obj.FontSize = s
		dbus.NotifyChange(obj, "FontSize")
	}
}

func (obj *Theme) setPropType(t int32) {
	if obj.Type != t {
		obj.Type = t
		dbus.NotifyChange(obj, "Type")
	}
}

func (obj *Theme) setAllProps() {
	kf := glib.NewKeyFile()
	defer kf.Free()

	var err error
	_, err = kf.LoadFromFile(path.Join(obj.filePath, "theme.ini"),
		glib.KeyFileFlagsKeepComments)
	if err != nil {
		Logger.Warningf("Load KeyFile '%s' Failed: %v", obj.filePath, err)
		return
	}

	var str string
	str, err = kf.GetString(THEME_GROUP_THEME, THEME_KEY_NAME)
	if err != nil {
		Logger.Warningf("Get '%s' failed: %v", THEME_KEY_NAME, err)
		return
	}
	obj.setPropName(str)

	str, err = kf.GetString(THEME_GROUP_COMPONENT, THEME_KEY_GTK)
	if err != nil {
		Logger.Warningf("Get '%s' failed: %v", THEME_KEY_GTK, err)
		return
	}
	obj.setPropGtkTheme(str)

	str, err = kf.GetString(THEME_GROUP_COMPONENT, THEME_KEY_ICON)
	if err != nil {
		Logger.Warningf("Get '%s' failed: %v", THEME_KEY_ICON, err)
		return
	}
	obj.setPropIconTheme(str)

	str, err = kf.GetString(THEME_GROUP_COMPONENT, THEME_KEY_SOUND)
	if err != nil {
		Logger.Warningf("Get '%s' failed: %v", THEME_KEY_SOUND, err)
		return
	}
	obj.setPropSoundTheme(str)

	str, err = kf.GetString(THEME_GROUP_COMPONENT, THEME_KEY_CURSOR)
	if err != nil {
		Logger.Warningf("Get '%s' failed: %v", THEME_KEY_CURSOR, err)
		return
	}
	obj.setPropCursorTheme(str)

	str, err = kf.GetString(THEME_GROUP_COMPONENT, THEME_KEY_BACKGROUND)
	if err != nil {
		Logger.Warningf("Get '%s' failed: %v", THEME_KEY_BACKGROUND, err)
		return
	}
	if ok, _ := regexp.MatchString(`^/`, str); !ok {
		if ok, _ = regexp.MatchString(`^file://`, str); !ok {
			str = path.Join(obj.filePath, THEME_BG_NAME, str)
			str = dutils.PathToURI(str, dutils.SCHEME_FILE)
		}
	} else {
		str = dutils.PathToURI(str, dutils.SCHEME_FILE)
	}
	obj.setPropBackground(str)

	var interval int
	interval, err = kf.GetInteger(THEME_GROUP_COMPONENT, THEME_KEY_FONT_SIZE)
	if err != nil {
		Logger.Warningf("Get '%s' failed: %v", THEME_KEY_FONT_SIZE, err)
		return
	}
	obj.setPropFontSize(int32(interval))
}
