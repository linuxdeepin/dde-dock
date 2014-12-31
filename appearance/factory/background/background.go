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

package background

import (
	"fmt"
	"os"
	"path"
	. "pkg.linuxdeepin.com/dde-daemon/appearance/utils"
	"pkg.linuxdeepin.com/lib/glib-2.0"
	"pkg.linuxdeepin.com/lib/graphic"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"regexp"
	"strings"
)

const (
	settingsKeyPictureURI = "picture-uri"
	defaultBgFile         = "/usr/share/backgrounds/default_background.jpg"

	dirModePerm = 0755
)

var (
	errUserPict = fmt.Errorf("Get User Picture Failed")
)

type Background struct {
	infoList []PathInfo

	eventHandler func([]string)
	watcher      *dutils.WatchProxy
}

func NewBackground(handler func([]string)) *Background {
	if handler == nil {
		return nil
	}

	bg := &Background{}

	bg.eventHandler = handler
	sysList, userList := getDirInfoList()
	bg.infoList = getBackgroundInfoList(sysList, userList)

	bg.watcher = dutils.NewWatchProxy()
	if bg.watcher != nil {
		bg.watcher.SetFileList(getBgDirList(bg.infoList))
		bg.watcher.SetEventHandler(bg.handleEvent)
		go bg.watcher.StartWatch()
	}

	return bg
}

func (bg *Background) IsValueValid(value string) bool {
	filename := dutils.DecodeURI(value)
	if dutils.IsFileExist(filename) &&
		graphic.IsSupportedImage(filename) {
		return true
	}

	return false
}

func (bg *Background) GetInfoByName(name string) (PathInfo, error) {
	return GetInfoByName(name, bg.infoList)
}

func (bg *Background) Set(uri string) error {
	settings := NewGSettings("com.deepin.dde.personalization")
	defer Unref(settings)

	value := settings.GetString(settingsKeyPictureURI)
	if isBackgroundSame(uri, value) {
		return nil
	}

	if isBackgroundInInfoList(uri, bg.infoList) {
		settings.SetString(settingsKeyPictureURI, uri)
		return nil
	}

	// cp uri to user wallpapers
	src := dutils.DecodeURI(uri)
	if !dutils.IsFileExist(src) {
		src = defaultBgFile
	}
	dir := getUserWallpaper()
	if len(dir) == 0 {
		return errUserPict
	}

	filename := path.Join(dir, path.Base(src))
	err := dutils.CopyFile(src, filename)
	if err != nil {
		return err
	}
	settings.SetString(settingsKeyPictureURI,
		dutils.EncodeURI(filename, dutils.SCHEME_FILE))

	return nil
}

func (bg *Background) Delete(uri string) error {
	for _, info := range bg.infoList {
		if info.FilePath == uri {
			if info.FileFlag == FileFlagUserOwned {
				return os.RemoveAll(dutils.DecodeURI(info.FilePath))
			}
			break
		}
	}

	return nil
}

func (bg *Background) Destroy() {
	if bg.watcher == nil {
		return
	}

	bg.watcher.EndWatch()
	bg.watcher = nil
}

func (bg *Background) GetNameStrList() []string {
	var list []string

	for _, info := range bg.infoList {
		list = append(list, info.FilePath)
	}

	return list
}

func (bg *Background) GetFlag(name string) int32 {
	for _, info := range bg.infoList {
		if name == info.FilePath {
			return info.FileFlag
		}
	}

	return -1
}

func (bg *Background) GetThumbnail(src string) string {
	thumb := GetThumbnail(BackgroundThumbSeed, src)
	if len(thumb) != 0 {
		return thumb
	}

	GenerateThumbnail()
	return dutils.DecodeURI(src)
}

func getBgDirList(infoList []PathInfo) []string {
	var list []string

	for _, info := range infoList {
		dir := dutils.DecodeURI(info.FilePath)
		list = append(list, dir)
	}

	return list
}

func getBackgroundInfoList(sysList, userList []PathInfo) []PathInfo {
	var list []PathInfo
	for _, dirInfo := range userList {
		images := getImageList(dirInfo.FilePath)
		for _, img := range images {
			var info PathInfo
			info.BaseName = path.Base(img)
			info.FilePath = dutils.EncodeURI(img, dutils.SCHEME_FILE)
			info.FileFlag = dirInfo.FileFlag
			list = append(list, info)
		}
	}

	for _, dirInfo := range sysList {
		images := getImageList(dirInfo.FilePath)
		for _, img := range images {
			if isBackgroundInInfoList(img, list) {
				continue
			}
			var info PathInfo
			info.BaseName = path.Base(img)
			info.FilePath = dutils.EncodeURI(img, dutils.SCHEME_FILE)
			info.FileFlag = dirInfo.FileFlag
			list = append(list, info)
		}
	}

	return list
}

func getDirInfoList() ([]PathInfo, []PathInfo) {
	sysList := []PathInfo{
		{
			BaseName: "",
			FilePath: "/usr/share/backgrounds",
			FileFlag: FileFlagSystemOwned,
		},
	}
	tmpList := GetInfoListFromDirs([]PathInfo{
		{
			BaseName: "",
			FilePath: "/usr/share/personalization/themes",
			FileFlag: FileFlagSystemOwned,
		},
	}, []string{"wallpapers"})
	for _, info := range tmpList {
		info.FilePath = path.Join(info.FilePath, "wallpapers")
		sysList = append(sysList, info)
	}

	homeDir := os.Getenv("HOME")
	userList := []PathInfo{
		{
			BaseName: "",
			FilePath: path.Join(homeDir, ".backgrounds"),
			FileFlag: FileFlagUserOwned,
		},
	}
	tmpList = GetInfoListFromDirs([]PathInfo{
		{
			BaseName: "",
			FilePath: path.Join(homeDir, ".local/personalization/themes"),
			FileFlag: FileFlagSystemOwned,
		},
	}, []string{"wallpapers"})
	for _, info := range tmpList {
		info.FilePath = path.Join(info.FilePath, "wallpapers")
		userList = append(userList, info)
	}
	dir := getUserWallpaper()
	if len(dir) != 0 {
		userList = append(userList, PathInfo{
			BaseName: "",
			FilePath: dir,
			FileFlag: FileFlagUserOwned,
		})
	}

	return sysList, userList
}

func getUserWallpaper() string {
	pict := glib.GetUserSpecialDir(
		glib.UserDirectoryDirectoryPictures)
	dir := path.Join(pict, "Wallpapers")
	if !dutils.IsFileExist(dir) {
		err := os.MkdirAll(dir, dirModePerm)
		if err != nil {
			return ""
		}
	}

	return dir
}

func getImageList(dir string) []string {
	var list []string

	fp, err := os.Open(dir)
	if err != nil {
		return list
	}

	infos, err := fp.Readdir(0)
	fp.Close()
	if err != nil {
		return list
	}

	for _, info := range infos {
		if !info.Mode().IsRegular() {
			continue
		}

		name := strings.ToLower(info.Name())
		ok, _ := regexp.MatchString(`\.jpe?g$|\.png$`, name)
		if !ok {
			continue
		}

		list = append(list, path.Join(dir, info.Name()))
	}

	return list
}

func isBackgroundInInfoList(bg string, infos []PathInfo) bool {
	for _, info := range infos {
		if isBackgroundSame(bg, info.FilePath) {
			return true
		}
	}

	return false
}

func isBackgroundSame(uri1, uri2 string) bool {
	bg1 := dutils.DecodeURI(uri1)
	bg2 := dutils.DecodeURI(uri2)

	md51, ok := dutils.SumFileMd5(bg1)
	if !ok {
		return false
	}
	md52, ok := dutils.SumFileMd5(bg2)
	if !ok {
		return false
	}
	if md51 == md52 {
		return true
	}

	return false
}
