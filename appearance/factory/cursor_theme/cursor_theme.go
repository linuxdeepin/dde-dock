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

package cursor_theme

// #cgo pkg-config: gtk+-3.0
// #cgo CFLAGS: -Wall -g
// #include "cursor.h"
import "C"

import (
	"fmt"
	"os"
	"path"
	. "pkg.linuxdeepin.com/dde-daemon/appearance/utils"
	xsettings "pkg.linuxdeepin.com/dde-daemon/xsettings_wrapper"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	systemThemePath = "/usr/share/icons"
	userThemePath   = ".icons"

	systemThumbPath = "/usr/share/personalization/thumbnail/CursorThemes"
	userThumbPath   = ".local/share/personalization/thumbnail/CursorThemes"

	dirModePerm = 0755
)

var (
	themeConditions = []string{"cursors"}

	errWriteCursor  = fmt.Errorf("Write cursor theme failed")
	errInvalidTheme = fmt.Errorf("Invalid Theme")
)

type CursorTheme struct {
	infoList []PathInfo

	eventHandler func([]string)
	watcher      *dutils.WatchProxy
}

func NewCursorTheme(handler func([]string)) *CursorTheme {
	if handler == nil {
		return nil
	}

	cursor := &CursorTheme{}

	cursor.eventHandler = handler
	sysDirs, userDirs := getDirInfoList()
	cursor.infoList = getThemeList(sysDirs, userDirs)

	cursor.watcher = dutils.NewWatchProxy()
	if cursor.watcher != nil {
		cursor.watcher.SetFileList(getDirList())
		cursor.watcher.SetEventHandler(cursor.handleEvent)
		go cursor.watcher.StartWatch()
	}
	xsettings.InitXSettings()
	// handle cursor changed by gtk+
	C.handle_cursor_changed()

	return cursor
}

func (cursor *CursorTheme) IsValueValid(value string) bool {
	if IsNameInInfoList(value, cursor.infoList) {
		return true
	}

	return false
}

func (cursor *CursorTheme) GetInfoByName(name string) (PathInfo, error) {
	return GetInfoByName(name, cursor.infoList)
}

func (cursor *CursorTheme) Set(theme string) error {
	if !IsNameInInfoList(theme, cursor.infoList) {
		return errInvalidTheme
	}

	// Ignore xsettings error
	setThemeViaXSettings(theme)
	setGtk2Theme(GetUserGtk2Config(), theme)
	setGtk3Theme(GetUserGtk3Config(), theme)

	dir := path.Join(os.Getenv("HOME"), ".icons/default")
	if !dutils.IsFileExist(dir) {
		err := os.MkdirAll(dir, dirModePerm)
		if err != nil {
			return err
		}
	}
	filename := path.Join(dir, "index.theme")
	ok := dutils.WriteKeyToKeyFile(filename,
		"Icon Theme", "Inherits", theme)
	if !ok {
		return errWriteCursor
	}

	return nil
}

func (cursor *CursorTheme) Delete(name string) error {
	for _, info := range cursor.infoList {
		if info.BaseName == name {
			if info.FileFlag == FileFlagUserOwned {
				return os.RemoveAll(info.FilePath)
			}
		}
	}

	return nil
}

func (cursor *CursorTheme) Destroy() {
	if cursor.watcher != nil {
		return
	}

	cursor.watcher.EndWatch()
	cursor.watcher = nil
}

func (cursor *CursorTheme) GetNameStrList() []string {
	return GetBaseNameList(cursor.infoList)
}

func (cursor *CursorTheme) GetFlag(name string) int32 {
	return GetFileFlagByName(name, cursor.infoList)
}

func (cursor *CursorTheme) GetThumbnail(theme string) string {
	var thumb string
	for _, info := range cursor.infoList {
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

		thumb = GetThumbnail("--cursor", info.FilePath)
		break
	}

	return thumb
}

func setThemeViaXSettings(theme string) error {
	return xsettings.SetString(xsettings.GtkStringCursorTheme, theme)
}

func setGtk2Theme(config, theme string) error {
	return WriteUserGtk2Config(config, "gtk-cursor-theme-name", theme)
}

func setGtk3Theme(config, theme string) error {
	return WriteUserGtk3Config(config, "gtk-cursor-theme-name", theme)
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
		}}

	return sysDirs, userDirs
}

func getThemeList(sysDirs, userDirs []PathInfo) []PathInfo {
	cursorList := GetInfoListFromDirs(userDirs, themeConditions)
	sysList := GetInfoListFromDirs(sysDirs, themeConditions)
	for _, info := range sysList {
		if IsNameInInfoList(info.BaseName, cursorList) {
			continue
		}

		cursorList = append(cursorList, info)
	}

	return cursorList
}

func getDirList() []string {
	var list []string

	list = append(list, systemThemePath)
	list = append(list, path.Join(os.Getenv("HOME"), userThemePath))

	return list
}
