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
	"dlib/dbus"
	"dlib/gio-2.0"
	"net/url"
	"os"
	"strconv"
)

const (
	MANAGER_DEST = "com.deepin.daemon.Themes"
	MANAGER_PATH = "/com/deepin/daemon/ThemeManager"
	MANAGER_IFC  = "com.deepin.daemon.ThemeManager"

	PERSONALIZATION_ID       = "com.deepin.dde.personalization"
	GKEY_CURRENT_THEME       = "current-theme"
	GKEY_CURRENT_BACKGROUND  = "current-picture"
	GKEY_CURRENT_SOUND_THEME = "current-sound-theme"
	DEFAULT_THEME_NAME       = "Deepin"
)

var (
	personSettings  = gio.NewSettings(PERSONALIZATION_ID)
	gnomeBgSettings = gio.NewSettings("org.gnome.desktop.background")
)

func (op *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		MANAGER_DEST,
		MANAGER_PATH,
		MANAGER_IFC,
	}
}

func (op *Manager) copyBackgroundFile(name string) (string, bool) {
	if len(name) <= 0 {
		return "", false
	}

	name, _ = objUtil.PathToFileURI(name)
	if urlInfo, err := url.Parse(name); err != nil {
		logObject.Info("Parse rawurl failed:", err)
		return "", false
	} else if urlInfo != nil {
		name = urlInfo.Scheme + "://" + urlInfo.Path
	}
	if ok := objUtil.IsFileExist(name); !ok {
		logObject.Warningf("BG %s not exist", name)
		return "", false
	}

	if !objUtil.IsElementExist(name, op.BackgroundList) {
		// Copy name to custom dir
		dir := getHomeDir() + THUMB_LOCAL_THEME_PATH + "/Custom/" + PERSON_BG_DIR_NAME
		if ok := objUtil.IsFileExist(dir); !ok {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return "", false
			}
		}

		src, _ := objUtil.URIToPath(name)
		baseName, _ := objUtil.GetBaseName(src)
		path := dir + "/" + baseName
		if ok := objUtil.CopyFile(src, path); !ok {
			return "", false
		}
		name = path
	}

	return name, true
}

func (op *Manager) setPropName(propName string) {
	switch propName {
	case "ThemeList":
		list := getThemeList()
		//logObject.Infof("Theme List: %v", list)
		tmpMap := make(map[string]PathInfo)
		tmpNameMap := make(map[string]string)
		tmp := []string{}
		for _, l := range list {
			id := genId()
			idStr := strconv.FormatInt(int64(id), 10)
			path := THEME_PATH + idStr
			tmp = append(tmp, path)
			tmpMap[path] = l
			tmpNameMap[l.path] = path
		}
		op.ThemeList = tmp
		op.pathNameMap = tmpMap
		themeNamePathMap = tmpNameMap
		dbus.NotifyChange(op, propName)
	case "GtkThemeList":
		list := getGtkThemeList()
		//logObject.Infof("Gtk Theme List: %v\n", list)
		tmp := []string{}
		for _, l := range list {
			tmp = append(tmp, l.path)
		}
		op.GtkThemeList = tmp
		dbus.NotifyChange(op, propName)
	case "IconThemeList":
		list := getIconThemeList()
		//logObject.Infof("Icon Theme List: %v\n", list)
		tmp := []string{}
		for _, l := range list {
			tmp = append(tmp, l.path)
		}
		op.IconThemeList = tmp
		dbus.NotifyChange(op, propName)
	case "CursorThemeList":
		list := getCursorThemeList()
		//logObject.Infof("Cursor Theme List: %v\n", list)
		tmp := []string{}
		for _, l := range list {
			tmp = append(tmp, l.path)
		}
		op.CursorThemeList = tmp
		dbus.NotifyChange(op, propName)
	case "SoundThemeList":
		op.SoundThemeList = getSoundThemeList()
		dbus.NotifyChange(op, propName)
	case "BackgroundList":
		op.BackgroundList = getBackgroundList()
		dbus.NotifyChange(op, propName)
	case "CurrentTheme":
		value := personSettings.GetString(GKEY_CURRENT_THEME)
		if obj := op.getThemeObject(value); obj != nil {
			// Check if theme properties valid
			if !objUtil.IsElementExist(obj.GtkTheme, op.GtkThemeList) ||
				!objUtil.IsElementExist(obj.IconTheme, op.IconThemeList) ||
				!objUtil.IsElementExist(obj.CursorTheme, op.CursorThemeList) {
				println("---- Reset CurrentTheme: ", value)
				op.CurrentTheme = DEFAULT_THEME_NAME
				personSettings.SetString(GKEY_CURRENT_THEME, DEFAULT_THEME_NAME)
			} else {
				op.CurrentTheme = value
				obj.updateThemeInfo()
			}
		}
		//} else {
		//println("---- Reset CurrentTheme: ", value)
		//op.CurrentTheme = DEFAULT_THEME_NAME
		//personSettings.SetString(GKEY_CURRENT_THEME, DEFAULT_THEME_NAME)
		//}
		dbus.NotifyChange(op, propName)
	}
}

func (op *Manager) getThemeObject(name string) *Theme {
	for _, path := range op.ThemeList {
		o, ok := themeObjMap[path]
		if !ok {
			continue
		}
		if o.Name == name {
			return o
		}
	}

	return nil
}

func (op *Manager) updateAllProps() {
	// TODO similar to newManager()
	op.setPropName("ThemeList")
	op.setPropName("GtkThemeList")
	op.setPropName("IconThemeList")
	op.setPropName("CursorThemeList")
	op.setPropName("SoundThemeList")
	op.setPropName("BackgroundList")

	// the following properties should be configure at end for their values
	// depends on other property
	op.setPropName("CurrentTheme")

	updateThemeObj(op.pathNameMap)
}

func (op *Manager) updateGSettingsKey(name string, value interface{}) {
	if value == nil {
		return
	}
	logObject.Infof("Update GSettings Key: %s, value: %s",
		name, value.(string))
	switch name {
	case GKEY_CURRENT_THEME:
		str := value.(string)
		logObject.Info("Set Theme Value: ", str)
		v := personSettings.GetString(GKEY_CURRENT_THEME)
		logObject.Info("Cur Theme Value: ", v)
		if v != str {
			personSettings.SetString(GKEY_CURRENT_THEME, str)
		}
	case GKEY_CURRENT_BACKGROUND:
		str := value.(string)
		if v := personSettings.GetString(GKEY_CURRENT_BACKGROUND); v != str {
			personSettings.SetString(GKEY_CURRENT_BACKGROUND, str)
		}
	case GKEY_CURRENT_SOUND_THEME:
		str := value.(string)
		if v := personSettings.GetString(GKEY_CURRENT_SOUND_THEME); v != str {
			personSettings.SetString(GKEY_CURRENT_SOUND_THEME, str)
		}
	}
}

func (op *Manager) listenSettingsChanged() {
	personSettings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case GKEY_CURRENT_THEME:
			value := personSettings.GetString(key)
			logObject.Infof("Theme GSettings Changed: %s", value)
			obj := op.getThemeObject(value)
			if obj != nil {
				obj.setThemeViaXSettings()
				op.setPropName("CurrentTheme")
			}
		case GKEY_CURRENT_BACKGROUND: // TODO
			value := personSettings.GetString(key)
			gbg := gnomeBgSettings.GetString(key)
			logObject.Infof("DEEPIN Bg GSettings Changed: %s", value)
			obj := op.getThemeObject(op.CurrentTheme)
			if obj != nil && obj.BackgroundFile != value {
				path, ok := op.copyBackgroundFile(value)
				if !ok {
					return
				}
				value = path
				if name := op.setTheme(obj.GtkTheme, obj.IconTheme,
					obj.CursorTheme, obj.FontSize,
					value, obj.SoundTheme); name != op.CurrentTheme {
					op.updateGSettingsKey(GKEY_CURRENT_THEME, name)
				} else {
					//obj.BackgroundFile = value
					obj.updateThemeInfo()
					dbus.NotifyChange(obj, "BackgroundFile")
				}
			}

			if value != gbg {
				gnomeBgSettings.SetString("picture-uri", value)
			}
		case GKEY_CURRENT_SOUND_THEME: // TODO
			value := personSettings.GetString(key)
			logObject.Infof("Sound GSettings Changed: %s", value)
			obj := op.getThemeObject(op.CurrentTheme)
			if obj == nil {
				break
			}
			if obj != nil && obj.SoundTheme != value {
				if name := op.setTheme(obj.GtkTheme, obj.IconTheme,
					obj.CursorTheme, obj.FontSize,
					obj.BackgroundFile, value); name != op.CurrentTheme {
					op.updateGSettingsKey(GKEY_CURRENT_THEME, name)
				}
			}
		}
	})

	gnomeBgSettings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case "picture-uri":
			v := gnomeBgSettings.GetString(key)
			tmp := personSettings.GetString(GKEY_CURRENT_BACKGROUND)
			if v != tmp {
				personSettings.SetString(GKEY_CURRENT_BACKGROUND, v)
			}
		}
	})
}
