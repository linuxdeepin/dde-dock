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
	"dlib/gio-2.0"
	"os/exec"
)

const (
	QT_CONFIG_FILE    = ".config/Trolltech.conf"
	DEFAULT_FONT_SIZE = " 11"

	QT_KEY_GROUP   = "Qt"
	QT_KEY_STYLE   = "stype"
	QT_STYLE_VALUE = "GTK+"
	QT_KEY_FONT    = "font"
	QT_FONT_ARGS   = ",-1,5,50,0,0,0,0,0"

	DEFAULT_FONT      = "WenQuanYi Micro Hei"
	DEFAULT_FONT_MONO = "WenQuanYi Micro Hei Mono"
)

var (
	wmPreSettings = gio.NewSettings("org.gnome.desktop.wm.preferences")
)

func (op *Theme) setThemeViaXSettings() {
	setGtkThemeViaXSettings(op.GtkTheme)
	setIconThemeViaXSettings(op.IconTheme)
	setCursorThemeViaXSettings(op.CursorTheme)
	setFontNameViaXSettings(DEFAULT_FONT, op.FontSize)

	bg := personSettings.GetString(GKEY_CURRENT_BACKGROUND)
	if bg != op.BackgroundFile {
		personSettings.SetString(GKEY_CURRENT_BACKGROUND,
			op.BackgroundFile)
	}
}

func setGtkThemeViaXSettings(name string) {
	objXSettings.SetString("Net/ThemeName", name)
	wmPreSettings.SetString("theme", name)
	homeDir := getHomeDir()
	if ok := objUtil.WriteKeyToKeyFile(homeDir+"/"+QT_CONFIG_FILE,
		QT_KEY_GROUP, QT_KEY_STYLE, QT_STYLE_VALUE); !ok {
		logObject.Infof("Write key: '%s', value: '%s', in file: '%s' failed", QT_KEY_STYLE, "GTK+", homeDir+"/"+QT_CONFIG_FILE)
	}
}

func setIconThemeViaXSettings(name string) {
	objXSettings.SetString("Net/IconThemeName", name)
}

func setCursorThemeViaXSettings(name string) {
	objXSettings.SetString("Gtk/CursorThemeName", name)
}

func setFontNameViaXSettings(name, size string) {
	//logObject.Infof("Set Font: %s\n", name)
	if len(name) <= 0 {
		name = DEFAULT_FONT
	}

	if len(size) <= 0 {
		size = DEFAULT_FONT_SIZE
	}

	objXSettings.SetString("Gtk/FontName", name+" "+size)
	wmPreSettings.SetString("titlebar-font", name+" Bold "+size)
	homeDir := getHomeDir()
	if ok := objUtil.WriteKeyToKeyFile(homeDir+"/"+QT_CONFIG_FILE,
		QT_KEY_GROUP, QT_KEY_FONT, "\""+name+","+size+QT_FONT_ARGS+"\""); !ok {
		logObject.Infof("Write key: '%s', value: '%s', in file: '%s' failed", QT_KEY_FONT, name, homeDir+"/"+QT_CONFIG_FILE)
	}
	setFontMono(DEFAULT_FONT_MONO, size)
}

func setFontMono(name, size string) {
	if len(name) <= 0 {
		return
	}

	if len(size) <= 0 {
		size = "10"
	}

	args := []string{}
	args = append(args, "-t")
	args = append(args, "string")
	args = append(args, "-s")
	args = append(args, "/desktop/gnome/interface/monospace_font_name")
	args = append(args, name+" "+size)

	exec.Command("/usr/bin/gconftool", args...).Run()
}
