/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package background

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"pkg.deepin.io/lib/imgutil"

	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
	"pkg.deepin.io/lib/xdg/basedir"
)

var (
	backgroundsCache   Backgrounds
	backgroundsCacheMu sync.Mutex
	fsChanged          bool

	CustomWallpapersConfigDir     string
	customWallpaperDeleteCallback func(file string)
	logger                        *log.Logger
)

const customWallpapersLimit = 10

func SetLogger(value *log.Logger) {
	logger = value
}

func SetCustomWallpaperDeleteCallback(fn func(file string)) {
	customWallpaperDeleteCallback = fn
}

func init() {
	CustomWallpapersConfigDir = filepath.Join(basedir.GetUserConfigDir(),
		"deepin/dde-daemon/appearance/custom-wallpapers")
	err := os.MkdirAll(CustomWallpapersConfigDir, 0755)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

type Background struct {
	Id        string
	Deletable bool
}

type Backgrounds []*Background

func refreshBackground() {
	logger.Debug("refresh background")
	var bgs Backgrounds
	// add custom
	for _, file := range getCustomBgFiles() {
		bgs = append(bgs, &Background{
			Id:        dutils.EncodeURI(file, dutils.SCHEME_FILE),
			Deletable: true,
		})
	}

	// add system
	for _, file := range getSysBgFiles() {
		bgs = append(bgs, &Background{
			Id:        dutils.EncodeURI(file, dutils.SCHEME_FILE),
			Deletable: false,
		})
	}

	backgroundsCache = bgs
	fsChanged = false
}

func ListBackground() Backgrounds {
	backgroundsCacheMu.Lock()
	defer backgroundsCacheMu.Unlock()

	if len(backgroundsCache) == 0 || fsChanged {
		refreshBackground()
	}
	return backgroundsCache
}

func NotifyChanged() {
	backgroundsCacheMu.Lock()
	fsChanged = true
	backgroundsCacheMu.Unlock()
}

var uiSupportedFormats = strv.Strv([]string{"jpeg", "png", "bmp", "tiff", "gif"})

func IsBackgroundFile(file string) bool {
	file = dutils.DecodeURI(file)
	format, err := imgutil.SniffFormat(file)
	if err != nil {
		return false
	}

	if uiSupportedFormats.Contains(format) {
		return true
	}
	return false
}

func (bgs Backgrounds) Get(uri string) *Background {
	uri = dutils.EncodeURI(uri, dutils.SCHEME_FILE)
	for _, info := range bgs {
		if uri == info.Id {
			return info
		}
	}
	return nil
}

func (bgs Backgrounds) ListGet(uris []string) Backgrounds {
	var ret Backgrounds
	for _, uri := range uris {
		v := bgs.Get(uri)
		if v == nil {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

func (bgs Backgrounds) Delete(uri string) error {
	info := bgs.Get(uri)
	if info == nil {
		return fmt.Errorf("not found '%s'", uri)
	}

	return info.Delete()
}

func (bgs Backgrounds) Thumbnail(uri string) (string, error) {
	return "", errors.New("not supported")
}

func (info *Background) Delete() error {
	if !info.Deletable {
		return fmt.Errorf("not custom")
	}

	file := dutils.DecodeURI(info.Id)
	err := os.Remove(file)

	if customWallpaperDeleteCallback != nil {
		customWallpaperDeleteCallback(file)
	}
	return err
}

func (info *Background) Thumbnail() (string, error) {
	return "", errors.New("not supported")
}

func Prepare(file string) (string, error) {
	file = dutils.DecodeURI(file)
	if isFileInDirs(file, systemWallpapersDir) {
		logger.Debug("is system")
		return file, nil
	}

	logger.Debug("is custom")
	return prepare(file)
}
