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

        PERSONALIZATION_ID = "com.deepin.dde.personalization"
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

func (op *Manager) getGtkPictPath(typePict, name string) string {
        path := ""

        t := getThemeType(name, op.ThemeList)
        if len(t) <= 0 {
                return ""
        }
        if t == PATH_TYPE_SYSTEM {
                path = THUMB_GTK_PATH + "/" + name
        } else if t == PATH_TYPE_LOCAL {
                homeDir := getHomeDir()
                path = homeDir + THUMB_LOCAL_GTK_PATH + "/" + name
        }

        switch typePict {
        case "thumbnail":
                path += "/thumbnail.png"
        case "preview":
                path += "/preview.png"
        }

        return path
}

func (op *Manager) getIconPictPath(typePict, name string) string {
        path := ""

        t := getThemeType(name, op.ThemeList)
        if len(t) <= 0 {
                return ""
        }
        if t == PATH_TYPE_SYSTEM {
                path = THUMB_ICON_PATH + "/" + name
        } else if t == PATH_TYPE_LOCAL {
                homeDir := getHomeDir()
                path = homeDir + THUMB_LOCAL_ICON_PATH + "/" + name
        }

        switch typePict {
        case "thumbnail":
                path += "/thumbnail.png"
        case "preview":
                path += "/preview.png"
        }

        return path
}

func (op *Manager) getCursorPictPath(typePict, name string) string {
        path := ""

        t := getThemeType(name, op.ThemeList)
        if len(t) <= 0 {
                return ""
        }
        if t == PATH_TYPE_SYSTEM {
                path = THUMB_CURSOR_PATH + "/" + name
        } else if t == PATH_TYPE_LOCAL {
                homeDir := getHomeDir()
                path = homeDir + THUMB_LOCAL_CURSOR_PATH + "/" + name
        }

        switch typePict {
        case "thumbnail":
                path += "/thumbnail.png"
        case "preview":
                path += "/preview.png"
        }

        return path
}

func (op *Manager) OnPropertiesChanged(propName string, old interface{}) {
        switch propName {
        case "CurrentTheme":
                if v, ok := old.(string); ok && v != op.CurrentTheme {
                        personSettings.SetString("current-theme",
                                op.CurrentTheme)
                }
        }
}

func (op *Manager) setPropName(propName string) {
        switch propName {
        case "ThemeList":
                list := getThemeList()
                for _, l := range list {
                        id := genId()
                        idStr := strconv.FormatInt(int64(id), 10)
                        path := THEME_PATH + idStr
                        op.ThemeList = append(op.ThemeList, path)
                        op.pathNameMap[path] = l
                }
                dbus.NotifyChange(op, propName)
        case "CurrentTheme":
                value := personSettings.GetString("current-theme")
                op.CurrentTheme = value
                dbus.NotifyChange(op, propName)
        case "GtkThemeList":
                list := getGtkThemeList()
                for _, l := range list {
                        op.GtkThemeList = append(op.GtkThemeList, l.path)
                }
                dbus.NotifyChange(op, propName)
        case "IconThemeList":
                list := getIconThemeList()
                for _, l := range list {
                        op.IconThemeList = append(op.IconThemeList, l.path)
                }
                dbus.NotifyChange(op, propName)
        case "CursorThemeList":
                list := getCursorThemeList()
                for _, l := range list {
                        op.CursorThemeList = append(op.CursorThemeList, l.path)
                }
                dbus.NotifyChange(op, propName)
        }
}

func (op *Manager) getCurrentThemeObject(name string) *Theme {
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
        op.setPropName("ThemeList")
        op.setPropName("CurrentTheme")
        op.setPropName("GtkThemeList")
        op.setPropName("IconThemeList")
        op.setPropName("CursorThemeList")

        updateThemeObj(op.pathNameMap)
}

func (op *Manager) listenSettingsChanged() {
        personSettings.Connect("changed", func(s *gio.Settings, key string) {
                switch key {
                case "current-picture":
                        value := personSettings.GetString(key)
                        obj := op.getCurrentThemeObject(op.CurrentTheme)
                        if obj != nil && obj.BackgroundFile != value {
                                op.SetTheme(obj.GtkTheme, obj.IconTheme,
                                        obj.GtkCursorTheme, obj.GtkFontName,
                                        value)
                        }
                case "current-theme":
                        value := personSettings.GetString(key)
                        if value == op.CurrentTheme {
                                break
                        }
                        obj := op.getCurrentThemeObject(value)
                        if obj != nil {
                                obj.setThemeViaXSettings()
                                op.CurrentTheme = value
                        }
                }
        })
}
