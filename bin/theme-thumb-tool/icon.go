/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	ICON_CONFIG = "index.theme"
)

func getIconTypeDir(info pathInfo) (string, string, string) {
	filename := path.Join(info.Path, ICON_CONFIG)
	if !objUtil.IsFileExist(filename) {
		return "", "", ""
	}

	value, ok := objUtil.ReadKeyFromKeyFile("/usr/share/icons/Deepin/index.theme", "Icon Theme", "Directories", []string{})
	if !ok {
		return "", "", ""
	}

	list, ok := value.([]string)
	if !ok {
		return "", "", ""
	}

	appsDir := ""
	placesDir := ""
	devicesDir := ""
	for _, l := range list {
		strs := strings.Split(l, ",")
		for _, t := range strs {
			if strings.Contains(t, "48") {
				if strings.Contains(t, "places") {
					placesDir = path.Join(info.Path, t)
				} else if strings.Contains(t, "apps") {
					appsDir = path.Join(info.Path, t)
				} else if strings.Contains(t, "devices") {
					devicesDir = path.Join(info.Path, t)
				}
			}
		}
	}

	return devicesDir, placesDir, appsDir
}

func getPngFile(dir string) string {
	fp, err := os.Open(dir)
	if err != nil {
		return ""
	}
	defer fp.Close()

	infos, err1 := fp.Readdir(0)
	if err1 != nil {
		return ""
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if ok, _ := regexp.MatchString(`\.png$`, info.Name()); ok {
			return path.Join(dir, info.Name())
		}
	}

	return ""
}

func getIconTypeFile(info pathInfo) (string, string, string) {
	device, place, app := getIconTypeDir(info)
	if len(device) < 1 || len(place) < 1 || len(app) < 1 {
		return "", "", ""
	}

	dFile := getPngFile(device)
	pFile := getPngFile(place)
	aFile := getPngFile(app)
	if (len(dFile) < 1) || (len(pFile) < 1) || (len(aFile) < 1) {
		return "", "", ""
	}

	return dFile, pFile, aFile
}
