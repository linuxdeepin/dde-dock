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

package deepin_theme

import (
	"fmt"
	"os"
	"path"
	. "pkg.deepin.io/dde/daemon/appearance/utils"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	systemThemePath = "/usr/share/personalization/themes"
	userThemePath   = ".local/share/personalization/themes"
)

var (
	themeConditions = []string{"theme.ini"}
	errInvalidTheme = fmt.Errorf("Invalid Theme")
)

type DeepinTheme struct {
	infoList []PathInfo

	eventHandler func([]string)
	watcher      *dutils.WatchProxy
}

func NewDeepinTheme(handler func([]string)) *DeepinTheme {
	if handler == nil {
		return nil
	}

	dtheme := &DeepinTheme{}

	sysDirs, userDirs := getDirInfoList()
	dtheme.infoList = getThemeList(sysDirs, userDirs)
	dtheme.eventHandler = handler
	dtheme.watcher = dutils.NewWatchProxy()
	if dtheme.watcher != nil {
		dtheme.watcher.SetFileList(getDirList())
		dtheme.watcher.SetEventHandler(dtheme.handleEvent)
		go dtheme.watcher.StartWatch()
	}

	return dtheme
}

func (dtheme *DeepinTheme) IsValueValid(value string) bool {
	for _, info := range dtheme.infoList {
		if info.BaseName == value {
			return true
		}
	}

	return false
}

func (dtheme *DeepinTheme) GetInfoByName(name string) (PathInfo, error) {
	return GetInfoByName(name, dtheme.infoList)
}

func (dtheme *DeepinTheme) Set(name string) error {
	if !dtheme.IsValueValid(name) {
		return errInvalidTheme
	}

	settings := NewGSettings("com.deepin.dde.personalization")
	defer Unref(settings)

	value := settings.GetString("current-theme")
	if value == name {
		return nil
	}

	settings.SetString("current-theme", name)
	return nil
}

func (dtheme *DeepinTheme) Delete(name string) error {
	for _, info := range dtheme.infoList {
		if info.BaseName == name {
			if info.FileFlag == FileFlagUserOwned {
				return os.RemoveAll(info.FilePath)
			}
			break
		}
	}

	return nil
}

func (dtheme *DeepinTheme) Destroy() {
	if dtheme.watcher == nil {
		return
	}

	dtheme.watcher.EndWatch()
	dtheme.watcher = nil
}

func (dtheme *DeepinTheme) GetNameStrList() []string {
	return GetBaseNameList(dtheme.infoList)
}

func (dtheme *DeepinTheme) GetFlag(value string) int32 {
	return GetFileFlagByName(value, dtheme.infoList)
}

func (dtheme *DeepinTheme) GetThumbnail(value string) string {
	var thumb string

	for _, info := range dtheme.infoList {
		if info.BaseName == value {
			thumb = path.Join(info.FilePath, "thumbnail.png")
			break
		}
	}

	if dutils.IsFileExist(thumb) {
		return thumb
	}

	return ""
}

func getDirInfoList() ([]PathInfo, []PathInfo) {
	sysDirs := []PathInfo{
		{
			BaseName: "",
			FilePath: systemThemePath,
			FileFlag: FileFlagSystemOwned,
		},
	}
	userDirs := []PathInfo{
		{
			BaseName: "",
			FilePath: path.Join(os.Getenv("HOME"), userThemePath),
			FileFlag: FileFlagUserOwned,
		},
	}

	return sysDirs, userDirs
}

func getThemeList(sysDirs, userDirs []PathInfo) []PathInfo {
	infoList := GetInfoListFromDirs(userDirs, themeConditions)
	sysList := GetInfoListFromDirs(sysDirs, themeConditions)
	for _, info := range sysList {
		if IsNameInInfoList(info.BaseName, infoList) {
			continue
		}

		infoList = append(infoList, info)
	}

	return infoList
}

func getDirList() []string {
	var list []string

	list = append(list, systemThemePath)
	userDir := path.Join(os.Getenv("HOME"), userThemePath)
	if !dutils.IsFileExist(userDir) {
		os.MkdirAll(userDir, 0755)
	}
	list = append(list, userDir)

	return list
}
