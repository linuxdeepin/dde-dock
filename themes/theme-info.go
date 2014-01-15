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
	"dlib/glib-2.0"
	"fmt"
	"io/ioutil"
)

type ThemeInfo struct {
	GtkTheme    string
	WindowTheme string //MetacityTheme
	IconTheme   string
	CursorTheme string
}

const (
	THEME_KEY_GTK    = "GtkTheme"
	THEME_KEY_WINDOW = "MetacityTheme"
	THEME_KEY_ICON   = "IconTheme"
	THEME_KEY_CURSOR = "CursorTheme"
	THEME_KEY_GROUP  = "X-GNOME-Metatheme"

	THEME_DIR        = "/usr/share/themes"
	THEME_FILE_INDEX = "index.theme"
)

var (
	systemThemes = []*ThemeInfo{}
)

func ReadThemeDir(dir string) {
	rets, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Read '%s' Failed: %s\n", dir, err)
		return
	}

	for _, v := range rets {
		name := v.Name()
		//fmt.Printf("Get Name: %s\n", name)
		if v.IsDir() {
			ReadThemeDir(dir + "/" + name)
		} else if name == THEME_FILE_INDEX {
			ReadThemeFile(dir + "/" + name)
		}
	}
}

func ReadThemeFile(filename string) {
	//fmt.Printf("Read File: %s\n", filename)
	conf := glib.NewKeyFile()
	_, err := conf.LoadFromFile(filename, glib.KeyFileFlagsNone)
	if err != nil {
		fmt.Printf("Key File Load File Failed: %s\n", err)
		return
	}
	info := &ThemeInfo{}
	info.GtkTheme, err = conf.GetString(THEME_KEY_GROUP, THEME_KEY_GTK)
	if err != nil {
		fmt.Printf("Get '%s : %s' Failed: %s\n",
			filename, THEME_KEY_GTK, err)
	}
	info.WindowTheme, err = conf.GetString(THEME_KEY_GROUP, THEME_KEY_WINDOW)
	if err != nil {
		fmt.Printf("Get '%s' Failed: %s\n", THEME_KEY_WINDOW, err)
	}
	info.IconTheme, err = conf.GetString(THEME_KEY_GROUP, THEME_KEY_ICON)
	if err != nil {
		fmt.Printf("Get '%s' Failed: %s\n", THEME_KEY_ICON, err)
	}
	info.CursorTheme, err = conf.GetString(THEME_KEY_GROUP, THEME_KEY_CURSOR)
	if err != nil {
		fmt.Printf("Get '%s' Failed: %s\n", THEME_KEY_CURSOR, err)
	}
	systemThemes = append(systemThemes, info)
}
