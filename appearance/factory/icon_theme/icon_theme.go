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

package icon_theme

import (
	"fmt"
	"os"
	"path"
	. "pkg.linuxdeepin.com/dde-daemon/appearance/utils"
	"pkg.linuxdeepin.com/dde-daemon/xsettings"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	systemThemePath = "/usr/share/icons"
	userThemePath   = ".icons"

	systemThumbPath = "/usr/share/personalization/thumbnail/IconThemes"
	userThumbPath   = ".local/share/personalization/thumbnail/IconThemes"
)

var (
	themeConditions = []string{"index.theme"}
	errInvalidTheme = fmt.Errorf("Invalid Theme")
)

type IconTheme struct {
	infoList []PathInfo

	eventHandler func([]string)
	watcher      *dutils.WatchProxy
}

func NewIconTheme(handler func([]string)) *IconTheme {
	if handler == nil {
		return nil
	}

	icon := &IconTheme{}

	icon.eventHandler = handler
	sysDirs, userDirs := getDirInfoList()
	icon.infoList = getThemeList(sysDirs, userDirs)

	icon.watcher = dutils.NewWatchProxy()
	if icon.watcher != nil {
		icon.watcher.SetFileList(getDirList())
		icon.watcher.SetEventHandler(icon.handleEvent)
		go icon.watcher.StartWatch()
	}

	return icon
}

func (icon *IconTheme) IsValueValid(value string) bool {
	if IsNameInInfoList(value, icon.infoList) {
		return true
	}

	return false
}

func (icon *IconTheme) GetInfoByName(name string) (PathInfo, error) {
	return GetInfoByName(name, icon.infoList)
}

func (icon *IconTheme) Set(theme string) error {
	if !IsNameInInfoList(theme, icon.infoList) {
		return errInvalidTheme
	}

	err := setThemeViaXSettings(theme)
	if err != nil {
		return err
	}

	return nil
}

func (icon *IconTheme) Delete(name string) error {
	for _, info := range icon.infoList {
		if info.BaseName == name {
			if info.FileFlag == FileFlagUserOwned {
				return os.RemoveAll(info.FilePath)
			}
		}
	}

	return nil
}

func (icon *IconTheme) Destroy() {
	if icon.watcher == nil {
		return
	}

	icon.watcher.EndWatch()
	icon.watcher = nil
}

func (icon *IconTheme) GetNameStrList() []string {
	return GetBaseNameList(icon.infoList)
}

func (icon *IconTheme) GetFlag(name string) int32 {
	return GetFileFlagByName(name, icon.infoList)
}

func (icon *IconTheme) GetThumbnail(theme string) string {
	var thumb string
	for _, info := range icon.infoList {
		if theme != info.BaseName {
			continue
		}

		if info.FileFlag == FileFlagSystemOwned {
			thumb = path.Join(systemThumbPath,
				theme+"-thumbnail.png")
		} else {
			thumb = path.Join(os.Getenv("HOME"),
				userThumbPath,
				theme+"-thumbnail.png")
		}
		if dutils.IsFileExist(thumb) {
			break
		}

		thumb = GetThumbnail("--icon", info.FilePath)
		break
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
	localList := GetInfoListFromDirs(userDirs, themeConditions)
	sysList := GetInfoListFromDirs(sysDirs, themeConditions)

	for _, info := range sysList {
		if IsNameInInfoList(info.BaseName, localList) {
			continue
		}

		localList = append(localList, info)
	}

	var iconList []PathInfo
	for _, info := range localList {
		filename := path.Join(info.FilePath, "index.theme")
		if isIconHidden(filename) || !hasDirectories(filename) {
			continue
		}

		iconList = append(iconList, info)
	}

	return iconList
}

func setThemeViaXSettings(theme string) error {
	proxy, err := xsettings.NewXSProxy()
	if err != nil {
		return err
	}
	defer proxy.Free()
	return proxy.SetString(xsettings.NetIconTheme, theme)
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

func isIconHidden(filename string) bool {
	if !dutils.IsFileExist(filename) {
		return true
	}

	v, ok := dutils.ReadKeyFromKeyFile(filename,
		"Icon Theme", "Hidden", false)
	if !ok {
		return false
	}

	return v.(bool)
}

/**
 * If icon theme index.theme has 'Directories', it's valid.
 **/
func hasDirectories(filename string) bool {
	if !dutils.IsFileExist(filename) {
		return false
	}

	v, ok := dutils.ReadKeyFromKeyFile(filename,
		"Icon Theme", "Directories", "")
	if !ok || len(v.(string)) == 0 {
		return false
	}

	return true
}
