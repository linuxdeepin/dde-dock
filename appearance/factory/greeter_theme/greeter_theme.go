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

package greeter_theme

import (
	"dbus/com/deepin/api/greeterutils"
	"fmt"
	"os"
	"path"
	. "pkg.deepin.io/dde/daemon/appearance/utils"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	systemThemePath = "/usr/share/personalization/greeter-theme"
	userThemePath   = ".local/share/personalization/greeter-theme"
)

var (
	themeConditions = []string{"thumb.png"}
	errInvalidUser  = fmt.Errorf("Invalid User")
	errInvalidTheme = fmt.Errorf("Invalid Theme")
)

type GreeterTheme struct {
	infoList []PathInfo

	eventHandler func([]string)
	watcher      *dutils.WatchProxy
}

func NewGreeterTheme(handler func([]string)) *GreeterTheme {
	if handler == nil {
		return nil
	}

	greeter := &GreeterTheme{}

	sysDirs, userDirs := getDirInfoList()
	greeter.infoList = getThemeList(sysDirs, userDirs)
	greeter.eventHandler = handler

	greeter.watcher = dutils.NewWatchProxy()
	if greeter.watcher != nil {
		greeter.watcher.SetFileList(getDirList())
		greeter.watcher.SetEventHandler(greeter.handleEvent)
		go greeter.watcher.StartWatch()
	}

	return greeter
}

func (greeter *GreeterTheme) IsValueValid(value string) bool {
	if IsNameInInfoList(value, greeter.infoList) {
		return true
	}

	return false
}

func (greeter *GreeterTheme) GetInfoByName(name string) (PathInfo, error) {
	return GetInfoByName(name, greeter.infoList)
}

func (greeter *GreeterTheme) Set(name string) error {
	if !IsNameInInfoList(name, greeter.infoList) {
		return errInvalidTheme
	}

	username := dutils.GetUserName()
	homeDir := os.Getenv("HOME")
	if homeDir == path.Join("/tmp", username) {
		return errInvalidUser
	}

	greeterObj, err := greeterutils.NewGreeterUtils(
		"com.deepin.api.GreeterUtils",
		"/com/deepin/api/GreeterUtils")
	if err != nil {
		return err
	}

	return greeterObj.SetGreeterTheme(username, name)
}

func (greeter *GreeterTheme) Delete(name string) error {
	for _, info := range greeter.infoList {
		if info.BaseName == name {
			if info.FileFlag == FileFlagUserOwned {
				return os.RemoveAll(info.FilePath)
			}
			break
		}
	}

	return nil
}

func (greeter *GreeterTheme) Destroy() {
	if greeter.watcher == nil {
		return
	}

	greeter.watcher.EndWatch()
	greeter.watcher = nil
}

func (greeter *GreeterTheme) GetNameStrList() []string {
	return GetBaseNameList(greeter.infoList)
}

func (greeter *GreeterTheme) GetFlag(name string) int32 {
	return GetFileFlagByName(name, greeter.infoList)
}

func (greeter *GreeterTheme) GetThumbnail(theme string) string {
	var thumb string
	for _, info := range greeter.infoList {
		if theme == info.BaseName {
			thumb = path.Join(info.FilePath, "thumb.png")
			thumb = dutils.EncodeURI(thumb, dutils.SCHEME_FILE)
			break
		}
	}

	return thumb
}

func getDirInfoList() ([]PathInfo, []PathInfo) {
	sysDirs := []PathInfo{
		{
			BaseName: "",
			FilePath: systemThemePath,
			FileFlag: FileFlagSystemOwned,
		},
	}
	userPath := path.Join(os.Getenv("HOME"), userThemePath)
	userDirs := []PathInfo{
		{
			BaseName: "",
			FilePath: userPath,
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
	list = append(list, path.Join(os.Getenv("HOME"), userThemePath))

	return list
}
