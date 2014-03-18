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
        "dlib/glib-2.0"
)

const (
        THEME_DEST = "com.deepin.daemon.Themes"
        THEME_PATH = "/com/deepin/daemon/Theme"
        //THEME_PATH         = "/com/deepin/daemon/Theme/Entry"
        THEME_IFC = "com.deepin.daemon.Theme"

        THEME_GROUP_THEME     = "Theme"
        THEME_KEY_NAME        = "Name"
        THEME_GROUP_COMPONENT = "Component"
        THEME_KEY_GTK         = "GtkTheme"
        THEME_KEY_ICONS       = "IconTheme"
        THEME_KEY_CURSOR      = "CursorTheme"
        THEME_KEY_GTK_FONT    = "FontName"
        THEME_KEY_BG          = "BackgroundFile"
        THEME_KEY_SOUND       = "SoundTheme"
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
        case "BackgroundFile":
                if v, ok := old.(string); ok && v != op.BackgroundFile {
                        personSettings.SetString(GKEY_CURRENT_PICTURE,
                                op.BackgroundFile)
                }
        }
}

func (op *Theme) updateThemeInfo() {
        filename := op.basePath + "/theme.ini"
        keyFile := glib.NewKeyFile()
        defer keyFile.Free()

        _, err := keyFile.LoadFromFile(filename,
                glib.KeyFileFlagsKeepComments)
        if err != nil {
                logObject.Info("LoadFile '%s' failed: %v",
                        filename, err)
                return
        }

        str, err1 := keyFile.GetString(THEME_GROUP_COMPONENT,
                THEME_KEY_GTK)
        if err1 != nil {
                logObject.Info("Get key '%s' value failed: %v",
                        THEME_KEY_GTK, err1)
                return
        }
        op.GtkTheme = str
        dbus.NotifyChange(op, "GtkTheme")

        str, err1 = keyFile.GetString(THEME_GROUP_COMPONENT,
                THEME_KEY_ICONS)
        if err1 != nil {
                logObject.Info("Get key '%s' value failed: %v",
                        THEME_KEY_CURSOR, err1)
                return
        }
        op.IconTheme = str
        dbus.NotifyChange(op, "IconTheme")

        str, err1 = keyFile.GetString(THEME_GROUP_COMPONENT,
                THEME_KEY_CURSOR)
        if err1 != nil {
                logObject.Info("Get key '%s' value failed: %v",
                        THEME_KEY_CURSOR, err1)
                return
        }
        op.CursorTheme = str
        dbus.NotifyChange(op, "CursorTheme")

        str, err1 = keyFile.GetString(THEME_GROUP_COMPONENT,
                THEME_KEY_GTK_FONT)
        if err1 != nil {
                logObject.Info("Get key '%s' value failed: %v",
                        THEME_KEY_GTK_FONT, err1)
                return
        }
        op.FontName = str
        dbus.NotifyChange(op, "FontName")

        str, err1 = keyFile.GetString(THEME_GROUP_COMPONENT,
                THEME_KEY_BG)
        if err1 != nil {
                logObject.Info("Get key '%s' value failed: %v",
                        THEME_KEY_BG, err1)
                return
        }
        op.BackgroundFile = str
        dbus.NotifyChange(op, "BackgroundFile")

        str, err1 = keyFile.GetString(THEME_GROUP_COMPONENT,
                THEME_KEY_SOUND)
        if err1 != nil {
                logObject.Info("Get key '%s' value failed: %v",
                        THEME_KEY_SOUND, err1)
                return
        }
        op.SoundThemeName = str
        dbus.NotifyChange(op, "SoundThemeName")
}
