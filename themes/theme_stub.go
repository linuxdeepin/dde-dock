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
)

const (
        THEME_DEST = "com.deepin.deamon.ThemeManager"
        THEME_PATH = "/com/deepin/daemon/Theme"
        //THEME_PATH         = "/com/deepin/daemon/Theme/Entry"
        THEME_IFC          = "com.deepin.daemon.Theme"
        PERSONALIZATION_ID = "com.deepin.dde.personalization"
)

var (
        personSettings = gio.NewSettings(PERSONALIZATION_ID)
)

func (op *Theme) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                THEME_DEST,
                op.path,
                THEME_IFC,
        }
}

func (op *Theme) OnPropertiesChanged(propName string, old interface{}) {
        switch propName {
        case "GtkTheme":
                if v, ok := old.(string); ok && v != op.GtkTheme {
                }
        case "IconTheme":
                if v, ok := old.(string); ok && v != op.IconTheme {
                }
        case "CursorTheme":
                if v, ok := old.(string); ok && v != op.CursorTheme {
                }
        case "FontName":
                if v, ok := old.(string); ok && v != op.FontName {
                }
        case "BackgroundFile":
                if v, ok := old.(string); ok && v != op.BackgroundFile {
                        personSettings.SetString("current-picture",
                                op.BackgroundFile)
                }
        }
}

func (op *Theme) setPropName(propName string) {
        switch propName {
        case "GtkTheme":
                dbus.NotifyChange(op, propName)
        case "IconTheme":
                dbus.NotifyChange(op, propName)
        case "CursorTheme":
                dbus.NotifyChange(op, propName)
        case "FontName":
                dbus.NotifyChange(op, propName)
        case "BackgroundFile":
                op.BackgroundFile = personSettings.GetString("current-picture")
                dbus.NotifyChange(op, propName)
        }
}

func (op *Theme) listenSettingsChanged() {
        personSettings.Connect("changed", func(s *gio.Settings, key string) {
                switch key {
                case "current-picture":
                        op.setPropName("BackgroundFile")
                }
        })
}
