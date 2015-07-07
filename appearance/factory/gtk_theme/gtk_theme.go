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

package gtk_theme

import (
	"fmt"
	"os"
	"path"
	. "pkg.deepin.io/dde-daemon/appearance/utils"
	"pkg.deepin.io/dde-daemon/xsettings"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	systemThemePath = "/usr/share/themes"
	userThemePath   = ".themes"

	systemThumbPath = "/usr/share/personalization/thumbnail/WindowThemes"
	userThumbPath   = ".local/share/personalization/thumbnail/WindowThemes"

	wmGSettingsSchema = "com.deepin.wrap.gnome.desktop.wm.preferences"
)

var (
	themeConditions = []string{"gtk-2.0", "gtk-3.0", "metacity-1"}
	errInvalidTheme = fmt.Errorf("Invalid Theme")
	errReadKey      = fmt.Errorf("Read key from keyfile failed")
	errWriteValue   = fmt.Errorf("Write value to keyfile failed")
)

type GtkTheme struct {
	infoList []PathInfo

	eventHandler func([]string)
	watcher      *dutils.WatchProxy
}

func NewGtkTheme(handler func([]string)) *GtkTheme {
	if handler == nil {
		return nil
	}

	gtk := &GtkTheme{}

	gtk.eventHandler = handler
	sysDirs, userDirs := getDirInfoList()
	gtk.infoList = getThemeList(sysDirs, userDirs)

	gtk.watcher = dutils.NewWatchProxy()
	if gtk.watcher != nil {
		gtk.watcher.SetFileList(getDirList())
		gtk.watcher.SetEventHandler(gtk.handleEvent)
		go gtk.watcher.StartWatch()
	}

	return gtk
}

func (gtk *GtkTheme) IsValueValid(value string) bool {
	if IsNameInInfoList(value, gtk.infoList) {
		return true
	}

	return false
}

func (gtk *GtkTheme) GetInfoByName(name string) (PathInfo, error) {
	return GetInfoByName(name, gtk.infoList)
}

func (gtk *GtkTheme) Set(theme string) error {
	if !IsNameInInfoList(theme, gtk.infoList) {
		return errInvalidTheme
	}

	err := setThemeViaXSettings(theme)
	if err != nil {
		return err
	}

	setCompizTheme(theme)
	setQt4Theme(GetUserQt4Config())

	return nil
}

func (gtk *GtkTheme) Delete(name string) error {
	for _, info := range gtk.infoList {
		if info.BaseName == name {
			if info.FileFlag == FileFlagUserOwned {
				return os.RemoveAll(info.FilePath)
			}
		}
	}

	return nil
}

func (gtk *GtkTheme) Destroy() {
	if gtk.watcher == nil {
		return
	}

	gtk.watcher.EndWatch()
	gtk.watcher = nil
}

func (gtk *GtkTheme) GetNameStrList() []string {
	return GetBaseNameList(gtk.infoList)
}

func (gtk *GtkTheme) GetFlag(theme string) int32 {
	return GetFileFlagByName(theme, gtk.infoList)
}

func (gtk *GtkTheme) GetThumbnail(theme string) string {
	var thumb string
	for _, info := range gtk.infoList {
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

		thumb = GetThumbnail("--gtk", info.FilePath)
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
	gtkList := GetInfoListFromDirs(userDirs, themeConditions)
	sysList := GetInfoListFromDirs(sysDirs, themeConditions)

	for _, info := range sysList {
		if IsNameInInfoList(info.BaseName, gtkList) {
			continue
		}

		gtkList = append(gtkList, info)
	}

	return gtkList
}

func setThemeViaXSettings(theme string) error {
	proxy, err := xsettings.NewXSProxy()
	if err != nil {
		return err
	}
	defer proxy.Free()
	return proxy.SetString(xsettings.NetThemeName, theme)
}

func setQt4Theme(config string) error {
	value, _ := dutils.ReadKeyFromKeyFile(config, "Qt", "style", "")

	if value == "GTK+" {
		return nil
	}
	ok := dutils.WriteKeyToKeyFile(config, "Qt", "style", "GTK+")
	if !ok {
		return errWriteValue
	}

	return nil
}

func setCompizTheme(theme string) {
	settings := CheckAndNewGSettings(wmGSettingsSchema)
	if settings == nil {
		return
	}

	settings.SetString("theme", theme)
	//Unref(settings)
}

func getDirList() []string {
	list := []string{}

	list = append(list, systemThemePath)
	userDir := path.Join(os.Getenv("HOME"), userThemePath)
	if !dutils.IsFileExist(userDir) {
		os.MkdirAll(userDir, 0755)
	}
	list = append(list, userDir)

	return list
}
