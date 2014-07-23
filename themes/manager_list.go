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

package themes

import (
	"os"
	"path"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"regexp"
	"strings"
)

const (
	THEME_TYPE_SYSTEM = 0
	THEME_TYPE_LOCAL  = 1
)

type pathInfo struct {
	Path string
	T    int32
}

type ThemeInfo struct {
	Name string
	Path string
	T    int32
}

type backgroundInfo ThemeInfo

const (
	THEME_THUMB  = "/usr/share/personalization/themes/Deepin/thumbnail.png"
	GTK_THUMB    = "/usr/share/personalization/thumbnail/WindowThemes/Deepin/thumbnail.png"
	ICON_THUMB   = "/usr/share/personalization/thumbnail/IconThemes/Deepin/thumbnail.png"
	CURSOR_THUMB = "/usr/share/personalization/thumbnail/CursorThemes/Deepin-Cursor/thumbnail.png"

	THEME_TYPE_GTK    = 1
	THEME_TYPE_ICON   = 2
	THEME_TYPE_CURSOR = 3
)

func getThemeList(dirs []pathInfo, conditions []string) []ThemeInfo {
	list := []ThemeInfo{}

	for _, dir := range dirs {
		var f *os.File
		var err error

		if f, err = os.Open(dir.Path); err != nil {
			Logger.Warningf("Open '%s' failed: %v", dir.Path, err)
			continue
		}
		defer f.Close()

		var infos []os.FileInfo
		if infos, err = f.Readdir(0); err != nil {
			Logger.Warningf("Readdir '%s' failed: %v", dir.Path, err)
			continue
		}

		for _, info := range infos {
			if !info.IsDir() {
				continue
			}

			filename := path.Join(dir.Path, info.Name())
			if filterTheme(filename, conditions) {
				tmp := ThemeInfo{}
				tmp.Name = info.Name()
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
		Logger.Warningf("Open '%s' failed: %v", dir, err)
		return false
	}
	defer f.Close()

	names := []string{}
	if names, err = f.Readdirnames(0); err != nil {
		Logger.Warningf("Readdirnames '%s' failed: %v", dir, err)
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

func getGtkThemeList() []ThemeInfo {
	localDir := path.Join(homeDir, THEME_LOCAL_PATH)
	sysDirs := []pathInfo{pathInfo{THEME_SYS_PATH, THEME_TYPE_SYSTEM}}
	localDirs := []pathInfo{pathInfo{localDir, THEME_TYPE_LOCAL}}
	conditions := []string{"gtk-2.0", "gtk-3.0", "metacity-1"}

	sysList := getThemeList(sysDirs, conditions)
	localList := getThemeList(localDirs, conditions)
	for _, l := range sysList {
		if isThemeInfoInList(l, localList) {
			continue
		}

		localList = append(localList, l)
	}

	return localList
}

func getIconThemeList() []ThemeInfo {
	localDir := path.Join(homeDir, ICON_LOCAL_PATH)
	sysDirs := []pathInfo{pathInfo{ICON_SYS_PATH, THEME_TYPE_SYSTEM}}
	localDirs := []pathInfo{pathInfo{localDir, THEME_TYPE_LOCAL}}
	conditions := []string{"index.theme"}

	sysList := getThemeList(sysDirs, conditions)
	localList := getThemeList(localDirs, conditions)
	for _, l := range sysList {
		if isThemeInfoInList(l, localList) {
			continue
		}

		localList = append(localList, l)
	}

	list := []ThemeInfo{}
	for _, l := range localList {
		filename := path.Join(l.Path, "index.theme")
		_, ok := dutils.ReadKeyFromKeyFile(filename,
			"Icon Theme", "Directories", []string{})
		value, _ := dutils.ReadKeyFromKeyFile(filename,
			"Icon Theme", "Hidden", false)
		v, ok1 := value.(bool)
		if !ok || (ok1 && v) {
			continue
		}

		list = append(list, l)
	}

	return list
}

func getCursorThemeList() []ThemeInfo {
	localDir := path.Join(homeDir, ICON_LOCAL_PATH)
	sysDirs := []pathInfo{pathInfo{ICON_SYS_PATH, THEME_TYPE_SYSTEM}}
	localDirs := []pathInfo{pathInfo{localDir, THEME_TYPE_LOCAL}}
	conditions := []string{"cursors", "cursor.theme"}

	sysList := getThemeList(sysDirs, conditions)
	localList := getThemeList(localDirs, conditions)
	for _, l := range sysList {
		if isThemeInfoInList(l, localList) {
			continue
		}

		localList = append(localList, l)
	}

	return localList
}

func getSoundThemeList() []ThemeInfo {
	sysDirs := []pathInfo{pathInfo{SOUND_THEME_PATH, THEME_TYPE_SYSTEM}}
	conditions := []string{"index.theme"}

	sysList := getThemeList(sysDirs, conditions)

	return sysList
}

func getBackgroundDir(dir string) ([]string, bool) {
	list := []string{}
	if !dutils.IsFileExist(dir) {
		return list, false
	}

	f, err := os.Open(dir)
	if err != nil {
		Logger.Warningf("Open '%s' failed: %v", dir, err)
		return list, false
	}
	defer f.Close()

	if infos, err := f.Readdir(0); err != nil {
		Logger.Warningf("Readdir '%s' failed: %v", dir, err)
		return list, false
	} else {
		conditions := []string{THEME_BG_NAME}
		for _, i := range infos {
			if !i.IsDir() {
				continue
			}
			filename := path.Join(dir, i.Name())
			if filterTheme(filename, conditions) {
				list = append(list, path.Join(filename, THEME_BG_NAME))
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
		Logger.Warningf("Open '%s' failed: %v", dir, err)
		return list, false
	}
	defer f.Close()

	if infos, err := f.Readdir(0); err != nil {
		Logger.Warningf("Readdir '%s' failed: %v", dir, err)
		return list, false
	} else {
		for _, i := range infos {
			if !i.Mode().IsRegular() {
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

func getBackgroundList() []backgroundInfo {
	bgList := []backgroundInfo{}

	if tmpList, ok := getImageList(DEFAULT_SYS_BG_DIR); ok {
		for _, l := range tmpList {
			tmp := backgroundInfo{}
			tmp.Name = path.Base(l)
			uri := dutils.PathToURI(l, dutils.SCHEME_FILE)
			tmp.Path = uri
			tmp.T = THEME_TYPE_SYSTEM
			bgList = append(bgList, tmp)
		}
	}

	// system personalization theme
	if dirs, ok := getBackgroundDir(PERSON_SYS_THEME_PATH); ok {
		for _, d := range dirs {
			if list, ok := getImageList(d); ok {
				for _, l := range list {
					tmp := backgroundInfo{}
					tmp.Name = path.Base(l)
					uri := dutils.PathToURI(l, dutils.SCHEME_FILE)
					tmp.Path = uri
					tmp.T = THEME_TYPE_SYSTEM
					//if !isBgInfoInList(tmp, list) {
					bgList = append(bgList, tmp)
					//}
				}
			}
		}
	}

	// local personalization theme
	if dirs, ok := getBackgroundDir(path.Join(homeDir, PERSON_LOCAL_THEME_PATH)); ok {
		for _, d := range dirs {
			if list, ok := getImageList(d); ok {
				for _, l := range list {
					tmp := backgroundInfo{}
					tmp.Name = path.Base(l)
					uri := dutils.PathToURI(l, dutils.SCHEME_FILE)
					tmp.Path = uri
					tmp.T = THEME_TYPE_LOCAL
					//if !isBgInfoInList(tmp, list) {
					bgList = append(bgList, tmp)
					//}
				}
			}
		}
	}
	pict := getUserPictureDir()
	userBG := path.Join(pict, "Wallpapers")
	if !dutils.IsFileExist(userBG) {
		return bgList
	}
	if tmpList, ok := getImageList(userBG); ok {
		for _, l := range tmpList {
			tmp := backgroundInfo{}
			tmp.Name = path.Base(l)
			uri := dutils.PathToURI(l, dutils.SCHEME_FILE)
			tmp.Path = uri
			tmp.T = THEME_TYPE_LOCAL
			bgList = append(bgList, tmp)
		}
	}

	return bgList
}

func isBackgroundSame(bg1, bg2 string) bool {
	bg1 = dutils.URIToPath(bg1)
	bg2 = dutils.URIToPath(bg2)
	if bg1 == bg2 {
		return true
	}

	return false
}

func isBgInfoInList(bg backgroundInfo, list []backgroundInfo) bool {
	for _, l := range list {
		if isBackgroundSame(bg.Path, l.Path) {
			return true
		}
	}

	return false
}

func isBgInfoListEqual(list1, list2 []backgroundInfo) bool {
	l1 := len(list1)
	l2 := len(list2)

	if l1 != l2 {
		return false
	}

	for i := 0; i < l1; i++ {
		if !isBackgroundSame(list1[i].Path, list2[i].Path) {
			return false
		}
	}

	return true
}

func getDThemeByDir(dir string, t int32) []ThemeInfo {
	if len(dir) < 1 {
		return []ThemeInfo{}
	}

	list := []ThemeInfo{}
	f, err := os.Open(dir)
	if err != nil {
		Logger.Warningf("Open '%s' failed: %v", dir, err)
		return list
	}
	defer f.Close()

	if infos, err := f.Readdir(0); err != nil {
		Logger.Warningf("Readdir '%s' failed: %v", dir, err)
		return list
	} else {
		conditions := []string{"theme.ini"}
		for _, i := range infos {
			if !i.IsDir() {
				continue
			}
			filename := path.Join(dir, i.Name())
			if filterTheme(filename, conditions) {
				tmp := ThemeInfo{}
				tmp.Name = i.Name()
				tmp.Path = filename
				tmp.T = t
				list = append(list, tmp)
			}
		}
	}

	return list
}

func getDThemeList() []ThemeInfo {
	if len(homeDir) < 1 {
		return []ThemeInfo{}
	}

	localList := getDThemeByDir(path.Join(homeDir, PERSON_LOCAL_THEME_PATH),
		THEME_TYPE_LOCAL)
	sysList := getDThemeByDir(PERSON_SYS_THEME_PATH, THEME_TYPE_SYSTEM)

	list := []ThemeInfo{}
	list = localList
	for _, l := range sysList {
		if isThemeInfoInList(l, list) {
			continue
		}

		list = append(list, l)
	}

	return list
}

func isThemeInfoSame(info1, info2 *ThemeInfo) bool {
	if info1 == nil || info2 == nil {
		return false
	}

	if info1.Name == info2.Name {
		return true
	}

	return false
}

func isThemeInfoInList(info ThemeInfo, list []ThemeInfo) bool {
	for _, l := range list {
		if isThemeInfoSame(&info, &l) {
			return true
		}
	}

	return false
}

func isThemeInfoListEqual(list1, list2 []ThemeInfo) bool {
	l1 := len(list1)
	l2 := len(list2)

	if l1 != l2 {
		return false
	}

	for i := 0; i < l1; i++ {
		if list1[i].Name != list2[i].Name ||
			list1[i].Path != list2[i].Path ||
			list1[i].T != list2[i].T {
			return false
		}
	}

	return true
}

func getGreeterThemeList() []ThemeInfo {
	homeDir := dutils.GetHomeDir()
	list := filterGreeterTheme(path.Join(homeDir, PERSON_LOCAL_GREETER_PATH), THEME_TYPE_LOCAL)
	tList := filterGreeterTheme(PERSON_SYS_GREETER_PATH, THEME_TYPE_SYSTEM)

	for _, l := range tList {
		if isThemeInfoInList(l, list) {
			continue
		}

		list = append(list, l)
	}

	return list
}

func filterGreeterTheme(dir string, t int32) []ThemeInfo {
	list := []ThemeInfo{}
	if !dutils.IsFileExist(dir) {
		return list
	}

	fp, err := os.Open(dir)
	if err != nil {
		Logger.Warningf("Open '%s' failed: %v", dir, err)
		return list
	}

	infos, err1 := fp.Readdir(0)
	if err1 != nil {
		Logger.Warningf("Readdir '%s' failed: %v", dir, err)
		return list
	}

	for _, info := range infos {
		if !info.Mode().IsDir() {
			continue
		}

		filepath := path.Join(dir, info.Name())
		if dutils.IsFileExist(path.Join(filepath, "thumb.png")) {
			tmp := ThemeInfo{info.Name(), filepath, t}
			list = append(list, tmp)
		}
	}

	return list
}
