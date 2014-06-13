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
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	THEME_TYPE_SYSTEM = 0
	THEME_TYPE_LOCAL  = 1

	BG_CACHE_DIR    = "cache/wallpapers"
	THEME_CACHE_DIR = "ccache/themes"
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
			Logger.Errorf("Open '%s' failed: %v", dir.Path, err)
			continue
		}
		defer f.Close()

		var infos []os.FileInfo
		if infos, err = f.Readdir(0); err != nil {
			Logger.Errorf("Readdir '%s' failed: %v", dir.Path, err)
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
		Logger.Errorf("Open '%s' failed: %v", dir, err)
		return false
	}
	defer f.Close()

	names := []string{}
	if names, err = f.Readdirnames(0); err != nil {
		Logger.Errorf("Readdirnames '%s' failed: %v", dir, err)
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
	homeDir, _ := objUtil.GetHomeDir()
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
	homeDir, _ := objUtil.GetHomeDir()
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

	return localList
}

func getCursorThemeList() []ThemeInfo {
	homeDir, _ := objUtil.GetHomeDir()
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
	if !objUtil.IsFileExist(dir) {
		return list, false
	}

	f, err := os.Open(dir)
	if err != nil {
		Logger.Errorf("Open '%s' failed: %v", dir, err)
		return list, false
	}
	defer f.Close()

	if infos, err := f.Readdir(0); err != nil {
		Logger.Errorf("Readdir '%s' failed: %v", dir, err)
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
	if !objUtil.IsFileExist(dir) {
		return list, false
	}

	f, err := os.Open(dir)
	if err != nil {
		Logger.Errorf("Open '%s' failed: %v", dir, err)
		return list, false
	}
	defer f.Close()

	if infos, err := f.Readdir(0); err != nil {
		Logger.Errorf("Readdir '%s' failed: %v", dir, err)
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

func getBackgroundList() []ThemeInfo {
	bgList := []ThemeInfo{}

	if tmpList, ok := getImageList(DEFAULT_SYS_BG_DIR); ok {
		for _, l := range tmpList {
			tmp := ThemeInfo{}
			tmp.Name = path.Base(l)
			uri, _ := objUtil.PathToFileURI(l)
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
					tmp := ThemeInfo{}
					tmp.Name = path.Base(l)
					uri, _ := objUtil.PathToFileURI(l)
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
	homeDir, _ := objUtil.GetHomeDir()
	if dirs, ok := getBackgroundDir(path.Join(homeDir, PERSON_LOCAL_THEME_PATH)); ok {
		for _, d := range dirs {
			if list, ok := getImageList(d); ok {
				for _, l := range list {
					tmp := ThemeInfo{}
					tmp.Name = path.Base(l)
					uri, _ := objUtil.PathToFileURI(l)
					tmp.Path = uri
					tmp.T = THEME_TYPE_LOCAL
					//if !isBgInfoInList(tmp, list) {
					bgList = append(bgList, tmp)
					//}
				}
			}
		}
	}

	return bgList
}

func getBgCachePath(src string) string {
	src, _ = objUtil.URIToPath(src)
	homeDir, _ := objUtil.GetHomeDir()
	bgDir := path.Join(homeDir, PERSON_LOCAL_BASE_PATH, BG_CACHE_DIR)
	md5Str, _ := getStrMd5(src)
	filename := path.Join(bgDir, md5Str+".png")
	if objUtil.IsFileExist(filename) {
		return filename
	}

	filename = path.Join(PERSON_SYS_BASE_PATH, BG_CACHE_DIR, md5Str+".png")
	if objUtil.IsFileExist(filename) {
		return filename
	}

	return ""
}

func isBackgroundSame(bg1, bg2 string) bool {
	bg1, _ = objUtil.URIToPath(bg1)
	bg2, _ = objUtil.URIToPath(bg2)
	if bg1 == bg2 {
		return true
	}

	return false
}

func isBgInfoInList(bg ThemeInfo, list []ThemeInfo) bool {
	for _, l := range list {
		if isBackgroundSame(bg.Path, l.Path) {
			return true
		}
	}

	return false
}

func isBgInfoListEqual(list1, list2 []ThemeInfo) bool {
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
		Logger.Errorf("Open '%s' failed: %v", dir, err)
		return list
	}
	defer f.Close()

	if infos, err := f.Readdir(0); err != nil {
		Logger.Errorf("Readdir '%s' failed: %v", dir, err)
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
	homeDir, ok := objUtil.GetHomeDir()
	if !ok {
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

	if info1.Name == info2.Name && info1.T == info2.T {
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

func getThemeCachePath(src string) string {
	src, _ = objUtil.URIToPath(src)
	homeDir, _ := objUtil.GetHomeDir()
	dir := path.Join(homeDir, PERSON_LOCAL_BASE_PATH, THEME_CACHE_DIR)
	md5Str, _ := getStrMd5(src)
	filename := path.Join(dir, md5Str+".png")
	if objUtil.IsFileExist(filename) {
		return filename
	}

	filename = path.Join(PERSON_SYS_BASE_PATH, THEME_CACHE_DIR, md5Str+".png")
	if objUtil.IsFileExist(filename) {
		return filename
	}

	return ""
}

func getDThemeThumb(name string) string {
	if name == "Custom" {
		if t, ok := GetManager().themeObjMap[name]; ok {
			return GetManager().GetThumbnail("background", t.Background)
		}
		return getThumbBg()
	}

	list := getDThemeList()
	for _, l := range list {
		if name == l.Name {
			thumb := path.Join(l.Path, "thumbnail.png")
			if objUtil.IsFileExist(thumb) {
				return thumb
			}
			break
		}
	}

	return getThumbBg()
}

func getGtkThumb(name string) string {
	list := getGtkThemeList()

	for _, l := range list {
		if name == l.Name {
			dest := getThemeCachePath(l.Path)
			if len(dest) > 0 {
				return dest
			}
		}
	}

	return getThumbBg()
}

func getIconThumb(name string) string {
	list := getIconThemeList()

	for _, l := range list {
		if name == l.Name {
			dest := getThemeCachePath(l.Path)
			if len(dest) > 0 {
				return dest
			}
		}
	}

	return getThumbBg()
}

func getCursorThumb(name string) string {
	list := getCursorThemeList()

	for _, l := range list {
		if name == l.Name {
			dest := getThemeCachePath(l.Path)
			if len(dest) > 0 {
				return dest
			}
		}
	}

	return getThumbBg()
}

func getBgThumb(bg string) string {
	dest := getBgCachePath(bg)
	if len(dest) > 0 {
		return dest
	}

	return getThumbBg()
}

func getThumbBg() string {
	list, _ := getImageList("/usr/share/personalization/thumb_bg")
	l := len(list)
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(l)

	return list[n]
}
