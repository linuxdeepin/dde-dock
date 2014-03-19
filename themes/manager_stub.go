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
        DEFAULT_SOUND_THEME_NAME = "LinuxDeepin"
        DEFAULT_BACKGROUND_FILE  = "file:///usr/share/backgrounds/default_background.jpg"

        SOUND_THEME_PATH      = "/usr/share/sounds/"
        SOUND_THEME_MAIN_FILE = "index.theme"
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
        switch propName {
        case "CurrentTheme":
                if v, ok := old.(string); ok && v != op.CurrentTheme {
                        personSettings.SetString(GKEY_CURRENT_THEME,
                                op.CurrentTheme)
                        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                                obj.setThemeViaXSettings()
                        }
                }
        case "CurrentSoundTheme": // TODO
                if v, ok := old.(string); ok && v != op.CurrentSoundTheme {
                        personSettings.SetString(GKEY_CURRENT_SOUND_THEME,
                                op.CurrentSoundTheme)
                }
        case "CurrentBackground": // TODO
                if v, ok := old.(string); ok && v != op.CurrentBackground {
                        personSettings.SetString(GKEY_CURRENT_BACKGROUND,
                                op.CurrentBackground)
                }
        }
}

func (op *Manager) setPropName(propName string) {
        switch propName {
        case "ThemeList":
                list := getThemeList()
                //logObject.Info("Theme List: %v", list)
                for _, l := range list {
                        id := genId()
                        idStr := strconv.FormatInt(int64(id), 10)
                        path := THEME_PATH + idStr
                        op.ThemeList = append(op.ThemeList, path)
                        op.pathNameMap[path] = l
                        themeNamePathMap[l.path] = path
                }
                dbus.NotifyChange(op, propName)
        case "GtkThemeList":
                list := getGtkThemeList()
                //logObject.Info("Gtk Theme List: %v", list)
                for _, l := range list {
                        op.GtkThemeList = append(op.GtkThemeList, l.path)
                }
                dbus.NotifyChange(op, propName)
        case "IconThemeList":
                list := getIconThemeList()
                //logObject.Info("Icon Theme List: %v", list)
                for _, l := range list {
                        op.IconThemeList = append(op.IconThemeList, l.path)
                }
                dbus.NotifyChange(op, propName)
        case "CursorThemeList":
                list := getCursorThemeList()
                //logObject.Info("Cursor Theme List: %v", list)
                for _, l := range list {
                        op.CursorThemeList = append(op.CursorThemeList, l.path)
                }
                dbus.NotifyChange(op, propName)
        case "SoundThemeList":
                op.SoundThemeList = getSoundThemeList()
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
        case "CurrentSoundTheme": // TODO
                value := personSettings.GetString(GKEY_CURRENT_SOUND_THEME)
                if isStringInArray(value, op.SoundThemeList) {
                        op.CurrentSoundTheme = value
                } else {
                        op.CurrentSoundTheme = DEFAULT_SOUND_THEME_NAME
                        personSettings.SetString(GKEY_CURRENT_SOUND_THEME, DEFAULT_SOUND_THEME_NAME)
                }
                dbus.NotifyChange(op, propName)
        case "CurrentBackground": // TODO
                value := personSettings.GetString(GKEY_CURRENT_BACKGROUND)
                if isStringInArray(value, op.BackgroundList) {
                        op.CurrentBackground = value
                } else {
                        op.CurrentBackground = DEFAULT_BACKGROUND_FILE
                        personSettings.SetString(GKEY_CURRENT_BACKGROUND, DEFAULT_BACKGROUND_FILE)
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

        // the following properties should be setup at end for their values
        // depends on other property
        op.setPropName("CurrentTheme")
        op.setPropName("CurrentSoundTheme")
        op.setPropName("CurrentBackground")

        updateThemeObj(op.pathNameMap)
}

func (op *Manager) updateCurrentTheme(name string) {
        logObject.Info("Update Current Theme: %s", name)
        if v := personSettings.GetString(GKEY_CURRENT_THEME); v != name {
                personSettings.SetString(GKEY_CURRENT_THEME, name)
        }
}

func (op *Manager) listenSettingsChanged() {
        personSettings.Connect("changed", func(s *gio.Settings, key string) {
                logObject.Info("Theme GSettings Key Changed: %s", key)
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
                        if value == op.CurrentBackground {
                                break
                        }

                        op.setPropName("CurrentBackground")

                        obj := op.getThemeObject(op.CurrentTheme)
                        if obj != nil && obj.BackgroundFile != value {
                                if name := op.setTheme(obj.GtkTheme, obj.IconTheme,
                                        obj.CursorTheme, obj.FontName,
                                        value, obj.SoundThemeName); name != op.CurrentTheme {
                                        op.updateCurrentTheme(name)
                                }
                        }
                case GKEY_CURRENT_SOUND_THEME: // TODO
                        value := personSettings.GetString(key)
                        if value == op.CurrentSoundTheme {
                                break
                        }
                        op.setPropName("SoundThemeList")
                }
        })
}
