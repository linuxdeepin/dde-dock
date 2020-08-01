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

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"time"

	"pkg.deepin.io/dde/api/themes"
	"pkg.deepin.io/dde/api/thumbnails/cursor"
	"pkg.deepin.io/dde/api/thumbnails/gtk"
	"pkg.deepin.io/dde/api/thumbnails/icon"
	"pkg.deepin.io/dde/api/thumbnails/images"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/lib/graphic"
)

const (
	thumbBgDir = "/var/cache/appearance/thumbnail/background"

	defaultWidth  = 128
	defaultHeight = 72
)

func genAllThumbnails(force bool) []string {
	var ret []string
	ret = append(ret, genGtkThumbnails(force)...)
	ret = append(ret, genIconThumbnails(force)...)
	ret = append(ret, genCursorThumbnails(force)...)
	ret = append(ret, genBgThumbnails(force)...)
	return ret
}

func genGtkThumbnails(force bool) []string {
	var ret []string
	list := themes.ListGtkTheme()
	for _, v := range list {
		thumb, err := gtk.ThumbnailForTheme(path.Join(v, "index.theme"),
			getThumbBg(), defaultWidth, defaultHeight, force)
		if err != nil {
			fmt.Printf("Gen '%s' thumbnail failed: %v\n", v, err)
			continue
		}
		ret = append(ret, thumb)
	}
	return ret
}

func genIconThumbnails(force bool) []string {
	var ret []string
	list := themes.ListIconTheme()
	for _, v := range list {
		thumb, err := icon.ThumbnailForTheme(path.Join(v, "index.theme"),
			getThumbBg(), defaultWidth, defaultHeight, force)
		if err != nil {
			fmt.Printf("Gen '%s' thumbnail failed: %v\n", v, err)
			continue
		}
		ret = append(ret, thumb)
	}
	return ret
}

func genCursorThumbnails(force bool) []string {
	var ret []string
	list := themes.ListCursorTheme()
	for _, v := range list {
		thumb, err := cursor.ThumbnailForTheme(path.Join(v, "cursor.theme"),
			getThumbBg(), defaultWidth, defaultHeight, force)
		if err != nil {
			fmt.Printf("Gen '%s' thumbnail failed: %v\n", v, err)
			continue
		}
		ret = append(ret, thumb)
	}
	return ret
}

func genBgThumbnails(force bool) []string {
	var ret []string
	infos := background.ListBackground()
	for _, info := range infos {
		thumb, err := images.ThumbnailForTheme(info.Id,
			defaultWidth, defaultHeight, force)
		if err != nil {
			fmt.Printf("Gen '%s' thumbnail failed: %v\n", info.Id, err)
			continue
		}
		ret = append(ret, thumb)
	}
	return ret
}

func getThumbBg() string {
	var imgs = getImagesInDir()
	if len(imgs) == 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(imgs))
	return imgs[idx]
}

func getImagesInDir() []string {
	finfos, err := ioutil.ReadDir(thumbBgDir)
	if err != nil {
		return nil
	}

	var imgs []string
	for _, finfo := range finfos {
		tmp := path.Join(thumbBgDir, finfo.Name())
		if !graphic.IsSupportedImage(tmp) {
			continue
		}
		imgs = append(imgs, tmp)
	}
	return imgs
}
