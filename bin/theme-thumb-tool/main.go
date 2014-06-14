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
	"dlib/utils"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
	"unsafe"
)

var (
	objUtil = utils.NewUtils()
)

const (
	_CMD_           = "theme-thumb-tool"
	_GTK_THUMB_CMD_ = "/usr/lib/deepin-daemon/gtk-thumb-tool"

	THUMB_CACHE_DIR = "cache"
)

func getThumbBg() string {
	list, _ := getImageList("/usr/share/personalization/thumb_bg")
	l := len(list)
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(l)

	return list[n]
}

func getThumbCachePath(t, src string, isSystem bool) string {
	src, _ = objUtil.URIToPath(src)
	md5Str, ok := getStrMd5(t + src)
	if !ok {
		return ""
	}

	if isSystem {
		filename := path.Join(PERSON_SYS_PATH, THUMB_CACHE_DIR)
		if !objUtil.IsFileExist(filename) {
			os.MkdirAll(filename, 0755)
		}
		filename = path.Join(filename, md5Str+".png")
		return filename
	} else {
		homeDir, ok := objUtil.GetHomeDir()
		if !ok {
			return ""
		}
		filename := path.Join(homeDir, PERSON_LOCAL_PATH, THUMB_CACHE_DIR)
		if !objUtil.IsFileExist(filename) {
			os.MkdirAll(filename, 0755)
		}
		filename = path.Join(filename, md5Str+".png")
		return filename
	}

	return ""
}

func genCursorThumbnail(info pathInfo, dest, bg string) bool {
	if len(dest) < 1 || len(bg) < 1 {
		return false
	}

	item1, item2, item3 := getCursorIcons(info)
	if len(item1) < 1 || len(item2) < 1 || len(item3) < 1 {
		fmt.Println("getCursorIcons Failed")
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
		fmt.Println("Generate Icon Thumbnail Error")
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
		fmt.Println("getIconTypeFile Failed")
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
		fmt.Println("Generate Icon Thumbnail Error")
		return false
	}

	return true
}

func printHelper() {
	fmt.Printf("Name\n\t%s: Theme Thumbnail Tool\n", _CMD_)
	fmt.Printf("Usage\n\t%s [Option] [Type]\n", _CMD_)
	fmt.Printf("Options:\n")
	fmt.Printf("\t--gtk:    Generate Gtk Theme Thumbnail\n")
	fmt.Printf("\t--icon:   Generate Icon Theme Thumbnail\n")
	fmt.Printf("\t--cursor: Generate Cursor Theme Thumbnail\n")
	fmt.Printf("\t--background: Generate Background Thumbnail\n")
	fmt.Printf("Types:\n")
	fmt.Printf("\tsystem: Save file to system dir\n")
	fmt.Printf("\tlocale: Save file to user dir\n")
}

func main() {
	if C.init_env() == 0 {
		fmt.Println("Can't generate thumbnails, try run this program under an X11 enviorment")
		return
	}

	if len(os.Args) < 2 {
		printHelper()
		return
	}

	op := os.Args[1]
	isSystem := false
	if len(os.Args) == 3 {
		if os.Args[2] == "system" {
			isSystem = true
		}
	}

	switch op {
	case "--gtk":
		list := getGtkList()
		for _, l := range list {
			dest := getThumbCachePath(op, l.Path, isSystem)
			if len(dest) < 1 || objUtil.IsFileExist(dest) {
				continue
			}

			name := path.Base(l.Path)
			out, err := exec.Command(_GTK_THUMB_CMD_, name, dest).Output()
			if err != nil || strings.Contains(string(out), "ERROR") {
				fmt.Printf("ERROR: Generate Gtk Thumbnail\n")
			}
		}
	case "--icon":
		list := getIconList()
		for _, l := range list {
			dest := getThumbCachePath(op, l.Path, isSystem)
			if len(dest) < 1 || objUtil.IsFileExist(dest) {
				continue
			}
			bg := getThumbBg()
			if !genIconThumbnail(l, dest, bg) {
				fmt.Printf("ERROR: Generate Icon Thumbnail\n")
			}
		}
	case "--cursor":
		list := getCursorList()
		for _, l := range list {
			dest := getThumbCachePath(op, l.Path, isSystem)
			if len(dest) < 1 || objUtil.IsFileExist(dest) {
				continue
			}
			bg := getThumbBg()
			if !genCursorThumbnail(l, dest, bg) {
				fmt.Printf("ERROR: Generate Cursor Thumbnail\n")
			}
		}
	case "--background":
		list := getBgList()
		for _, l := range list {
			dest := getThumbCachePath(op, l.Path, isSystem)
			if len(dest) < 1 || objUtil.IsFileExist(dest) {
				continue
			}
			err := graphic.ThumbnailImage(l.Path, dest, 130, 73, graphic.PNG)
			if err != nil {
				fmt.Println("ERROR:", err)
			}
		}
	default:
		fmt.Printf("Invalid option: %s\n\n", op)
		printHelper()
	}
}
