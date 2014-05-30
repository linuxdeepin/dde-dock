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
	"os"
	"strings"
)

const (
	BACKGROUND_DEFAULT_DIR = "/usr/share/backgrounds"
	PERSON_BG_DIR_NAME     = "wallpappers"
)

func getBackgroundList() []string {
	list := []string{}

	if defaultList, ok := getImagePath(BACKGROUND_DEFAULT_DIR); ok {
		list = append(list, defaultList...)
	}

	if dirs, ok := getDirPath(THUMB_THEME_PATH); ok {
		for _, d := range dirs {
			if l, ok := getImagePath(d); ok {
				list = append(list, l...)
			}
		}
	}

	homeDir := getHomeDir()
	if dirs, ok := getDirPath(homeDir + THUMB_LOCAL_THEME_PATH); ok {
		for _, d := range dirs {
			if l, ok := getImagePath(d); ok {
				list = append(list, l...)
			}
		}
	}

	//logObject.Infof("Background List: %v", list)
	return list
}

func getDirPath(dir string) ([]string, bool) {
	if ok := objUtil.IsFileExist(dir); !ok {
		return []string{}, false
	}

	f, err := os.Open(dir)
	if err != nil {
		logObject.Infof("Opne '%s' failed: %v", dir, err)
		return []string{}, false
	}
	defer f.Close()

	fi, err1 := f.Readdir(0)
	if err1 != nil {
		logObject.Infof("Readdir '%s' failed: %v", dir, err1)
		return []string{}, false
	}

	list := []string{}
	conditions := []string{PERSON_BG_DIR_NAME}
	for _, i := range fi {
		if i.IsDir() {
			path := dir + "/" + i.Name()
			if filterTheme(path, conditions) {
				list = append(list, path+"/"+PERSON_BG_DIR_NAME)
			}
		}
	}

	return list, true
}

func getImagePath(dir string) ([]string, bool) {
	if ok := objUtil.IsFileExist(dir); !ok {
		return []string{}, false
	}

	f, err := os.Open(dir)
	if err != nil {
		logObject.Infof("Opne '%s' failed: %v", dir, err)
		return []string{}, false
	}
	defer f.Close()

	fi, err1 := f.Readdir(0)
	if err1 != nil {
		logObject.Infof("Readdir '%s' failed: %v", dir, err1)
		return []string{}, false
	}

	list := []string{}
	for _, i := range fi {
		if i.Mode().IsRegular() {
			name := i.Name()
			if strings.Contains(name, "jpg") ||
				strings.Contains(name, "JPG") ||
				strings.Contains(name, "png") ||
				strings.Contains(name, "PNG") {
				path := dir + "/" + name
				if tmp, ok := objUtil.PathToFileURI(path); ok {
					list = append(list, tmp)
				}
			}
		}
	}

	return list, true
}
