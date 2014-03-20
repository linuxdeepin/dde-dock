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

package main

import (
        "dlib/dbus"
        "dlib/gio-2.0"
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
        personSettings = gio.NewSettings(PERSONALIZATION_ID)
)

func (op *Manager) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                MANAGER_DEST,
                MANAGER_PATH,
                MANAGER_IFC,
        }
}

func (op *Manager) OnPropertiesChanged(propName string, old interface{}) {
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
                if _, ok := themeNamePathMap[value]; ok {
                        op.CurrentTheme = value
                } else {
                        op.CurrentTheme = DEFAULT_THEME_NAME
                        personSettings.SetString(GKEY_CURRENT_THEME, DEFAULT_THEME_NAME)
                }
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
        logObject.Infof("Update GSettings Key: %s", name)
        switch name {
        case GKEY_CURRENT_THEME:
                str := value.(string)
                if v := personSettings.GetString(GKEY_CURRENT_THEME); v != str {
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
                logObject.Infof("Theme GSettings Key Changed: %s", key)
                switch key {
                case GKEY_CURRENT_THEME:
                        value := personSettings.GetString(key)
                        if value == op.CurrentTheme {
                                break
                        }
                        obj := op.getThemeObject(value)
                        if obj != nil {
                                obj.setThemeViaXSettings()
                                op.setPropName("CurrentTheme")
                        }
                case GKEY_CURRENT_BACKGROUND: // TODO
                        value := personSettings.GetString(key)
                        obj := op.getThemeObject(op.CurrentTheme)
                        if obj != nil && obj.BackgroundFile != value {
                                if name := op.setTheme(obj.GtkTheme, obj.IconTheme,
                                        obj.CursorTheme, obj.FontName,
                                        value, obj.SoundTheme); name != op.CurrentTheme {
                                        op.updateGSettingsKey(GKEY_CURRENT_THEME, name)
                                }
                        }
                case GKEY_CURRENT_SOUND_THEME: // TODO
                        value := personSettings.GetString(key)
                        obj := op.getThemeObject(op.CurrentTheme)
                        if obj != nil && obj.SoundTheme != value {
                                if name := op.setTheme(obj.GtkTheme, obj.IconTheme,
                                        obj.CursorTheme, obj.FontName,
                                        obj.BackgroundFile, value); name != op.CurrentTheme {
                                        op.updateGSettingsKey(GKEY_CURRENT_THEME, name)
                                }
                        }
                }
        })
}
