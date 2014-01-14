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

const (
	ICON_KEY_GROUP   = "Icon Theme"
	ICON_KEY_NAME    = "Name"
	ICON_DIR         = "/usr/share/icons"
	ICON_FILE_INDEX  = "index.theme"
	ICON_FILE_CURSOR = "cursor.theme"
)

var (
	icons = []string{}
)

func ReadIconDir(dir string) {
	rets, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Read '%s' Failed: %s\n", dir, err)
		return
	}

	for _, v := range rets {
		name := v.Name()
		//fmt.Printf("Get Name: %s\n", name)
		if v.IsDir() {
			ReadIconDir(dir + "/" + name)
		} else if name == ICON_FILE_INDEX || name == ICON_FILE_CURSOR {
			ReadIconFile(dir + "/" + name)
		}
	}
}

func ReadIconFile(filename string) {
	conf := glib.NewKeyFile()
	_, err := conf.LoadFromFile(filename, glib.KeyFileFlagsNone)
	if err != nil {
		fmt.Printf("Key File Load File Failed: %s\n", err)
		return
	}

	name, err1 := conf.GetString(ICON_KEY_GROUP, ICON_KEY_NAME)
	if err1 != nil {
                fmt.Printf("Get '%s : %s' Failed: %s\n", 
                filename,ICON_KEY_NAME, err1)
		return
	}
	icons = append(icons, name)
}
