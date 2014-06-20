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

package main

// #cgo pkg-config: gtk+-3.0
// #cgo CFLAGS: -Wall -g
// #include <stdlib.h>
// #include "common.h"
import "C"

import (
	"dlib/graphic"
	"dlib/logger"
	dutils "dlib/utils"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
	"unsafe"
)

var (
	Logger = logger.NewLogger("theme-thumb-tool")
)

const (
	_CMD_           = "theme-thumb-tool"
	_GTK_THUMB_CMD_ = "/usr/lib/deepin-daemon/gtk-thumb-tool"

	THUMB_CACHE_DIR = "autogen"
)

func getUserPictureDir() string {
	str := C.get_user_pictures_dir()
	//defer C.free(unsafe.Pointer(str))

	ret := C.GoString(str)

	return ret
}

func getThumbBg() string {
	list, _ := getImageList("/usr/share/personalization/thumbnail/autogen/thumb_bg")
	l := len(list)
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(l)

	return list[n]
}

func getThumbCachePath(t, src, outDir string) string {
	if len(outDir) < 1 {
		Logger.Debug("Output Dir Error")
		return ""
	}

	src = dutils.URIToPath(src)
	md5Str, ok := getStrMd5(t + src)
	if !ok {
		return ""
	}

	filename := path.Join(outDir, THUMB_CACHE_DIR)
	if !dutils.IsFileExist(filename) {
		err := os.MkdirAll(filename, 0755)
		if err != nil {
			return ""
		}
	}
	filename = path.Join(filename, md5Str+".png")
	return filename
}

func genCursorThumbnail(info pathInfo, dest, bg string) bool {
	if len(dest) < 1 || len(bg) < 1 {
		return false
	}

	item1, item2, item3 := getCursorIcons(info)
	if len(item1) < 1 || len(item2) < 1 || len(item3) < 1 {
		Logger.Debug("getCursorIcons Failed")
		return false
	}

	cBg := C.CString(bg)
	defer C.free(unsafe.Pointer(cBg))
	cDest := C.CString(dest)
	defer C.free(unsafe.Pointer(cDest))
	cItem1 := C.CString(item1)
	defer C.free(unsafe.Pointer(cItem1))
	cItem2 := C.CString(item2)
	defer C.free(unsafe.Pointer(cItem2))
	cItem3 := C.CString(item3)
	defer C.free(unsafe.Pointer(cItem3))

	ret := C.gen_icon_preview(cBg, cDest, cItem1, cItem2, cItem3)
	if int(ret) == -1 {
		Logger.Debug("Generate Cursor Thumbnail Error")
		return false
	}

	return true
}

func genIconThumbnail(info pathInfo, dest, bg string) bool {
	if len(dest) < 1 || len(bg) < 1 {
		return false
	}

	item1, item2, item3 := getIconTypeFile(info)
	if len(item1) < 1 || len(item2) < 1 || len(item3) < 1 {
		Logger.Debug("getIconTypeFile Failed")
		return false
	}

	cBg := C.CString(bg)
	defer C.free(unsafe.Pointer(cBg))
	cDest := C.CString(dest)
	defer C.free(unsafe.Pointer(cDest))
	cItem1 := C.CString(item1)
	defer C.free(unsafe.Pointer(cItem1))
	cItem2 := C.CString(item2)
	defer C.free(unsafe.Pointer(cItem2))
	cItem3 := C.CString(item3)
	defer C.free(unsafe.Pointer(cItem3))

	ret := C.gen_icon_preview(cBg, cDest, cItem1, cItem2, cItem3)
	if int(ret) == -1 {
		Logger.Debug("Generate Icon Thumbnail Error")
		return false
	}

	return true
}

func printHelper() {
	Logger.Debugf("Name\n\t%s: Theme Thumbnail Tool\n", _CMD_)
	Logger.Debugf("Usage\n\t%s [Option] [Output Dir]\n", _CMD_)
	Logger.Debugf("Options:\n")
	Logger.Debugf("\t--gtk:    Generate Gtk Theme Thumbnail\n")
	Logger.Debugf("\t--icon:   Generate Icon Theme Thumbnail\n")
	Logger.Debugf("\t--cursor: Generate Cursor Theme Thumbnail\n")
	Logger.Debugf("\t--background: Generate Background Thumbnail\n")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			Logger.Debugf("Error: %v\n", err)
			os.Exit(0)
		}
	}()

	if C.init_env() == 0 {
		Logger.Debug("Can't generate thumbnails, try run this program under an X11 enviorment")
		return
	}

	if len(os.Args) < 2 {
		printHelper()
		return
	}

	op := os.Args[1]
	outDir := ""
	homeDir := dutils.GetHomeDir()
	if len(homeDir) > 0 {
		outDir = path.Join(homeDir, PERSON_LOCAL_THUMB_PATH)
	}
	if len(os.Args) == 3 {
		outDir = os.Args[2]
	}

	switch op {
	case "--gtk":
		list := getGtkList()
		for _, l := range list {
			dest := getThumbCachePath(op, l.Path, outDir)
			if len(dest) < 1 || dutils.IsFileExist(dest) {
				continue
			}

			name := path.Base(l.Path)
			out, err := exec.Command(_GTK_THUMB_CMD_, name, dest).Output()
			if err != nil || strings.Contains(string(out), "ERROR") {
				Logger.Debugf("ERROR: Generate Gtk Thumbnail\n")
			}
		}
	case "--icon":
		list := getIconList()
		for _, l := range list {
			dest := getThumbCachePath(op, l.Path, outDir)
			if len(dest) < 1 || dutils.IsFileExist(dest) {
				continue
			}
			bg := getThumbBg()
			if !genIconThumbnail(l, dest, bg) {
				Logger.Debugf("ERROR: Generate Icon Thumbnail\n")
			}
		}
	case "--cursor":
		list := getCursorList()
		for _, l := range list {
			dest := getThumbCachePath(op, l.Path, outDir)
			if len(dest) < 1 || dutils.IsFileExist(dest) {
				continue
			}
			bg := getThumbBg()
			if !genCursorThumbnail(l, dest, bg) {
				Logger.Debugf("ERROR: Generate Cursor Thumbnail\n")
			}
		}
	case "--background":
		list := getBgList()
		for _, l := range list {
			dest := getThumbCachePath(op, l.Path, outDir)
			if len(dest) < 1 || dutils.IsFileExist(dest) {
				continue
			}
			err := graphic.ThumbnailImage(l.Path, dest, 128, 72, graphic.PNG)
			if err != nil {
				Logger.Debug("ERROR:", err)
			}
		}
	default:
		Logger.Debugf("Invalid option: %s\n\n", op)
		printHelper()
	}
}
