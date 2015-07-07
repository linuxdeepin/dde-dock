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
	dutils "pkg.deepin.io/lib/utils"
	"regexp"
	"strings"
)

const (
	THEME_TYPE_SYS   = 0
	THEME_TYPE_LOCAL = 1
)

const (
	PERSON_SYS_PATH         = "/usr/share/personalization"
	PERSON_LOCAL_PATH       = ".local/share/personalization"
	PERSON_SYS_THEME_PATH   = PERSON_SYS_PATH + "/themes"
	PERSON_LOCAL_THEME_PATH = PERSON_LOCAL_PATH + "/themes"

	THEME_SYS_PATH   = "/usr/share/themes"
	THEME_LOCAL_PATH = ".themes"
	ICON_SYS_PATH    = "/usr/share/icons"
	ICON_LOCAL_PATH  = ".icons"

	BG_DEFAULT_PATH = "/usr/share/backgrounds"
	PERSON_BG_NAME  = "wallpapers"
)

var (
	PERSON_SYS_THUMB_PATH   = path.Join(PERSON_SYS_PATH, "thumbnail")
	PERSON_LOCAL_THUMB_PATH = path.Join(PERSON_LOCAL_PATH, "thumbnail")
)

type pathInfo struct {
	Path string
	T    int32
}

func getThemeList(dirs []pathInfo, conditions []string) []pathInfo {
	list := []pathInfo{}

	for _, dir := range dirs {
		var f *os.File
		var err error

		if f, err = os.Open(dir.Path); err != nil {
			logger.Debugf("Open '%s' failed: %v", dir.Path, err)
			continue
		}
		defer f.Close()

		var infos []os.FileInfo
		if infos, err = f.Readdir(0); err != nil {
			logger.Debugf("Readdir '%s' failed: %v", dir.Path, err)
			continue
		}

		for _, info := range infos {
			if !info.IsDir() {
				continue
			}

			filename := path.Join(dir.Path, info.Name())
			if filterTheme(filename, conditions) {
				tmp := pathInfo{}
				tmp.Path = filename
				tmp.T = dir.T
				list = append(list, tmp)
			}
		}
	}

	return list
}

func filterTheme(dir string, conditions []string) bool {
	var err error

	var f *os.File
	if f, err = os.Open(dir); err != nil {
		logger.Debugf("Open '%s' failed: %v", dir, err)
		return false
	}
	defer f.Close()

	names := []string{}
	if names, err = f.Readdirnames(0); err != nil {
		logger.Debugf("Readdirnames '%s' failed: %v", dir, err)
		return false
	}

	cnt := 0
	for _, name := range names {
		for _, c := range conditions {
			if name == c {
				cnt++
				break
			}
		}
	}

	if cnt == len(conditions) {
		return true
	}

	return false
}

func getDThemeList() []pathInfo {
	list := []pathInfo{}
	sysDirs := []pathInfo{pathInfo{PERSON_SYS_THEME_PATH, THEME_TYPE_SYS}}
	conditions := []string{"theme.ini"}
	sysList := getThemeList(sysDirs, conditions)

	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		list = sysList
		return list
	}

	dir := path.Join(homeDir, PERSON_LOCAL_THEME_PATH)
	localDirs := []pathInfo{pathInfo{dir, THEME_TYPE_LOCAL}}
	localList := getThemeList(localDirs, conditions)

	list = localList
	for _, l := range sysList {
		if isPathThemeInList(l, list) {
			continue
		}
		list = append(list, l)
	}

	return list
}

func getGtkList() []pathInfo {
	list := []pathInfo{}
	sysDirs := []pathInfo{pathInfo{THEME_SYS_PATH, THEME_TYPE_SYS}}
	conditions := []string{"gtk-2.0", "gtk-3.0", "metacity-1"}
	sysList := getThemeList(sysDirs, conditions)

	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		list = sysList
		return list
	}

	dir := path.Join(homeDir, THEME_LOCAL_PATH)
	localDirs := []pathInfo{pathInfo{dir, THEME_TYPE_LOCAL}}
	localList := getThemeList(localDirs, conditions)

	list = localList
	for _, l := range sysList {
		if isPathThemeInList(l, list) {
			continue
		}
		list = append(list, l)
	}

	return list
}

func getIconList() []pathInfo {
	list := []pathInfo{}
	sysDirs := []pathInfo{pathInfo{ICON_SYS_PATH, THEME_TYPE_SYS}}
	conditions := []string{"index.theme"}
	sysList := getThemeList(sysDirs, conditions)

	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		list = sysList
		return list
	}

	dir := path.Join(homeDir, ICON_LOCAL_PATH)
	localDirs := []pathInfo{pathInfo{dir, THEME_TYPE_LOCAL}}
	localList := getThemeList(localDirs, conditions)

	list = localList
	for _, l := range sysList {
		if isPathThemeInList(l, list) {
			continue
		}
		list = append(list, l)
	}

	return list
}

func getCursorList() []pathInfo {
	list := []pathInfo{}
	sysDirs := []pathInfo{pathInfo{ICON_SYS_PATH, THEME_TYPE_SYS}}
	conditions := []string{"cursors"}
	sysList := getThemeList(sysDirs, conditions)

	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		list = sysList
		return list
	}

	dir := path.Join(homeDir, ICON_LOCAL_PATH)
	localDirs := []pathInfo{pathInfo{dir, THEME_TYPE_LOCAL}}
	localList := getThemeList(localDirs, conditions)

	list = localList
	for _, l := range sysList {
		if isPathThemeInList(l, list) {
			continue
		}
		list = append(list, l)
	}

	return list
}

func isPathThemeEqual(info1, info2 *pathInfo) bool {
	if info1 == nil || info2 == nil {
		return false
	}

	if path.Base(info1.Path) == path.Base(info2.Path) {
		return true
	}

	return false
}

func isPathThemeInList(info pathInfo, list []pathInfo) bool {
	for _, l := range list {
		if isPathThemeEqual(&info, &l) {
			return true
		}
	}

	return false
}

func isPathBgEqual(info1, info2 *pathInfo) bool {
	if info1 == nil || info2 == nil {
		return false
	}

	if info1.Path == info2.Path {
		return true
	}

	return false
}

func isPathBgInList(info pathInfo, list []pathInfo) bool {
	for _, l := range list {
		if isPathBgEqual(&info, &l) {
			return true
		}
	}

	return false
}

func getBgDir(dir string) ([]string, bool) {
	list := []string{}
	if !dutils.IsFileExist(dir) {
		return list, false
	}

	f, err := os.Open(dir)
	if err != nil {
		logger.Debugf("Open '%s' failed: %v", dir, err)
		return list, false
	}
	defer f.Close()

	if infos, err := f.Readdir(0); err != nil {
		logger.Debugf("Readdir '%s' failed: %v", dir, err)
		return list, false
	} else {
		conditions := []string{PERSON_BG_NAME}
		for _, i := range infos {
			if !i.IsDir() {
				continue
			}
			filename := path.Join(dir, i.Name())
			if filterTheme(filename, conditions) {
				list = append(list, path.Join(filename, PERSON_BG_NAME))
			}
		}
	}

	return list, true
}

func getImageList(dir string) ([]string, bool) {
	list := []string{}
	if !dutils.IsFileExist(dir) {
		return list, false
	}

	f, err := os.Open(dir)
	if err != nil {
		logger.Debugf("Open '%s' failed: %v", dir, err)
		return list, false
	}
	defer f.Close()

	if infos, err := f.Readdir(0); err != nil {
		logger.Debugf("Readdir '%s' failed: %v", dir, err)
		return list, false
	} else {
		for _, i := range infos {
			if i.IsDir() {
				continue
			}
			name := strings.ToLower(i.Name())
			ok1, _ := regexp.MatchString(`\.jpe?g$`, name)
			ok2, _ := regexp.MatchString(`\.png$`, name)
			if ok1 || ok2 {
				list = append(list, path.Join(dir, i.Name()))
			}
		}
	}

	return list, true
}

func getBgList() []pathInfo {
	list := []pathInfo{}

	if tmp, ok := getImageList(BG_DEFAULT_PATH); ok {
		for _, l := range tmp {
			t := pathInfo{}
			t.Path = l
			t.T = THEME_TYPE_SYS
			list = append(list, t)
		}
	}

	// system bg
	if dirs, ok := getBgDir(PERSON_SYS_THEME_PATH); ok {
		for _, d := range dirs {
			if tmp, ok := getImageList(d); ok {
				for _, l := range tmp {
					t := pathInfo{}
					t.Path = l
					t.T = THEME_TYPE_SYS
					if !isPathBgInList(t, list) {
						list = append(list, t)
					}
				}
			}
		}
	}

	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		return list
	}
	// local bg
	if dirs, ok := getBgDir(path.Join(homeDir, PERSON_LOCAL_THEME_PATH)); ok {
		for _, d := range dirs {
			if tmp, ok := getImageList(d); ok {
				for _, l := range tmp {
					t := pathInfo{}
					t.Path = l
					t.T = THEME_TYPE_LOCAL
					if !isPathBgInList(t, list) {
						list = append(list, t)
					}
				}
			}
		}
	}
	pict := getUserPictureDir()
	userBG := path.Join(pict, "Wallpapers")
	if !dutils.IsFileExist(userBG) {
		return list
	}
	if tmpList, ok := getImageList(userBG); ok {
		for _, l := range tmpList {
			tmp := pathInfo{}
			tmp.Path = l
			tmp.T = THEME_TYPE_LOCAL
			list = append(list, tmp)
		}
	}

	return list
}
