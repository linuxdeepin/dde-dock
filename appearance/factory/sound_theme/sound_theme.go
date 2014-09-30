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

package sound_theme

import (
	"fmt"
	"os"
	. "pkg.linuxdeepin.com/dde-daemon/appearance/utils"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	systemThemePath       = "/usr/share/sounds"
	settingsKeySoundTheme = "current-sound-theme"
)

var (
	themeConditions = []string{"index.theme"}
	errInvalidTheme = fmt.Errorf("Invalid Theme")
)

type SoundTheme struct {
	infoList []PathInfo

	eventHandler func([]string)
	watcher      *dutils.WatchProxy
}

func NewSoundTheme(handler func([]string)) *SoundTheme {
	if handler == nil {
		return nil
	}
	sound := &SoundTheme{}

	sysDirs, userDirs := getDirInfoList()
	sound.infoList = getThemeList(sysDirs, userDirs)
	sound.eventHandler = handler
	sound.watcher = dutils.NewWatchProxy()
	if sound.watcher != nil {
		sound.watcher.SetFileList(getDirList())
		sound.watcher.SetEventHandler(sound.handleEvent)
		go sound.watcher.StartWatch()
	}

	return sound
}

func (sound *SoundTheme) IsValueValid(value string) bool {
	if IsNameInInfoList(value, sound.infoList) {
		return true
	}

	return false
}

func (sound *SoundTheme) GetInfoByName(name string) (PathInfo, error) {
	return GetInfoByName(name, sound.infoList)
}

func (sound *SoundTheme) Set(name string) error {
	if !IsNameInInfoList(name, sound.infoList) {
		return errInvalidTheme
	}

	settings := gio.NewSettings("com.deepin.dde.personalization")
	value := settings.GetString(settingsKeySoundTheme)
	if value == name {
		settings.Unref()
		return nil
	}

	for _, info := range sound.infoList {
		if name == info.BaseName {
			settings.SetString(settingsKeySoundTheme, name)
			break
		}
	}

	settings.Unref()
	return nil
}

func (sound *SoundTheme) Delete(name string) error {
	for _, info := range sound.infoList {
		if info.BaseName == name {
			if info.FileFlag == FileFlagUserOwned {
				return os.RemoveAll(info.FilePath)
			}
			break
		}
	}

	return nil
}

func (sound *SoundTheme) Destroy() {
	if sound.watcher == nil {
		return
	}

	sound.watcher.EndWatch()
	sound.watcher = nil
}

func (sound *SoundTheme) GetNameStrList() []string {
	return GetBaseNameList(sound.infoList)
}

func (sound *SoundTheme) GetFlag(name string) int32 {
	return GetFileFlagByName(name, sound.infoList)
}

// Adapter interface
func (sound *SoundTheme) GetThumbnail(name string) string {
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

	return sysDirs, nil
}

func getThemeList(sysDirs, userDirs []PathInfo) []PathInfo {
	return GetInfoListFromDirs(sysDirs, themeConditions)
}

func getDirList() []string {
	var list []string

	list = append(list, systemThemePath)
	return list
}
