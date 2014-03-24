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
        "dlib/gio-2.0"
)

const (
        QT_CONFIG_FILE    = ".config/Trolltech.conf"
        FONT_DEFAULT_SIZE = " 11"
        TITLE_FONT_SIZE   = " 10"

        QT_KEY_GROUP   = "Qt"
        QT_KEY_STYLE   = "stype"
        QT_STYLE_VALUE = "GTK+"
        QT_KEY_FONT    = "font"
        QT_FONT_ARGS   = ",11,-1,5,50,0,0,0,0,0"
)

var (
        wmPreSettings = gio.NewSettings("org.gnome.desktop.wm.preferences")
)

func (op *Theme) setThemeViaXSettings() {
        setGtkThemeViaXSettings(op.GtkTheme)
        setIconThemeViaXSettings(op.IconTheme)
        setCursorThemeViaXSettings(op.CursorTheme)
        setGtkFontThemeViaXSettings(op.FontName)
}

func setGtkThemeViaXSettings(name string) {
        objXSettings.SetString("Net/ThemeName", name)
        wmPreSettings.SetString("theme", name)
        homeDir := getHomeDir()
        logObject.Info("File Path: ", homeDir+"/"+QT_CONFIG_FILE)
        if ok := objUtil.WriteKeyToKeyFile(homeDir+"/"+QT_CONFIG_FILE,
                QT_KEY_GROUP, QT_KEY_STYLE, QT_STYLE_VALUE); !ok {
                logObject.Infof("Write key: '%s', value: '%s', in file: '%s' failed", QT_KEY_STYLE, "GTK+", homeDir+"/"+QT_CONFIG_FILE)
        }
        logObject.Info("Set QT Style End...")
}

func setIconThemeViaXSettings(name string) {
        objXSettings.SetString("Net/IconThemeName", name)
}

func setCursorThemeViaXSettings(name string) {
        objXSettings.SetString("Gtk/CursorThemeName", name)
}

func setGtkFontThemeViaXSettings(name string) {
        //logObject.Infof("Set Font: %s\n", name)
        objXSettings.SetString("Gtk/FontName", name+FONT_DEFAULT_SIZE)
        wmPreSettings.SetString("titlebar-font", name+TITLE_FONT_SIZE)
        homeDir := getHomeDir()
        if ok := objUtil.WriteKeyToKeyFile(homeDir+"/"+QT_CONFIG_FILE,
                QT_KEY_GROUP, QT_KEY_FONT, "\""+name+QT_FONT_ARGS+"\""); !ok {
                logObject.Infof("Write key: '%s', value: '%s', in file: '%s' failed", QT_KEY_FONT, name, homeDir+"/"+QT_CONFIG_FILE)
        }
}
