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
	"math/rand"
	"os"
	"os/exec"
	"path"
	"pkg.deepin.io/lib/graphic"
	"pkg.deepin.io/lib/log"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
	"time"
	"unsafe"
)

var (
	forceFlag = false
	logger    = log.NewLogger("daemon/theme-thumb-tool")
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
		logger.Debug("Output Dir Error")
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

func getIconFiles(theme string) (string, string, string) {
	if len(theme) < 1 {
		return "", "", ""
	}

	cTheme := C.CString(theme)
	defer C.free(unsafe.Pointer(cTheme))

	cFileM := C.CString(FILE_MANAGER)
	defer C.free(unsafe.Pointer(cFileM))
	cFileManager := C.get_icon_filepath(cTheme, cFileM)
	fileManager := C.GoString(cFileManager)

	cFolderName := C.CString(FOLDER)
	defer C.free(unsafe.Pointer(cFolderName))
	cFolder := C.get_icon_filepath(cTheme, cFolderName)
	folderName := C.GoString(cFolder)

	cFullTrash := C.CString(USER_TRASH_FULL)
	defer C.free(unsafe.Pointer(cFullTrash))
	cTrash := C.get_icon_filepath(cTheme, cFullTrash)
	userTrash := C.GoString(cTrash)
	if len(userTrash) < 1 {
		cUserTrash := C.CString(USER_TRASH)
		defer C.free(unsafe.Pointer(cUserTrash))
		cTrash = C.get_icon_filepath(cTheme, cUserTrash)
		userTrash = C.GoString(cTrash)
	}

	return fileManager, folderName, userTrash
}

func genCursorThumbnail(info pathInfo, dest, bg string) bool {
	if len(dest) < 1 || len(bg) < 1 {
		return false
	}

	item1, item2, item3 := getCursorIcons(info)
	if len(item1) < 1 || len(item2) < 1 || len(item3) < 1 {
		logger.Debug("getCursorIcons Failed")
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
		logger.Debug("Generate Cursor Thumbnail Error")
		return false
	}

	return true
}

func genIconThumbnail(info pathInfo, dest, bg string) bool {
	if len(dest) < 1 || len(bg) < 1 {
		return false
	}

	theme := path.Base(info.Path)
	item1, item2, item3 := getIconFiles(theme)
	if len(item1) < 1 || len(item2) < 1 || len(item3) < 1 {
		logger.Debug("getIconTypeFile Failed")
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
		logger.Debug("Generate Icon Thumbnail Error")
		return false
	}

	return true
}

func printHelper() {
	logger.Debugf("Name\n\t%s: Theme Thumbnail Tool\n", _CMD_)
	logger.Debugf("Usage\n\t%s [Option] [Output Dir]\n", _CMD_)
	logger.Debugf("Options:\n")
	logger.Debugf("\t--gtk:    Generate Gtk Theme Thumbnail\n")
	logger.Debugf("\t--icon:   Generate Icon Theme Thumbnail\n")
	logger.Debugf("\t--cursor: Generate Cursor Theme Thumbnail\n")
	logger.Debugf("\t--background: Generate Background Thumbnail\n")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Debugf("Error: %v\n", err)
			os.Exit(0)
		}
	}()

	if C.init_env() == 0 {
		logger.Debug("Can't generate thumbnails, try run this program under an X11 enviorment")
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

	if !isVersionSame() {
		cleanThumbCache()
		forceFlag = true
	}

	switch op {
	case "-a":
		doGenGtkThumb(forceFlag, outDir)
		doGenIconThumb(forceFlag, outDir)
		doGenCursorThumb(forceFlag, outDir)
		doGenBgThumb(forceFlag, outDir)
	case "--gtk":
		doGenGtkThumb(forceFlag, outDir)
	case "--icon":
		doGenIconThumb(forceFlag, outDir)
	case "--cursor":
		doGenCursorThumb(forceFlag, outDir)
	case "--background":
		doGenBgThumb(forceFlag, outDir)
	default:
		logger.Debugf("Invalid option: %s\n\n", op)
		printHelper()
	}
}

func doGenGtkThumb(forceFlag bool, outDir string) {
	list := getGtkList()
	for _, l := range list {
		dest := getThumbCachePath("--gtk", l.Path, outDir)
		if len(dest) < 1 || (!forceFlag && dutils.IsFileExist(dest)) {
			continue
		}

		name := path.Base(l.Path)
		bg := getThumbBg()
		out, err := exec.Command(_GTK_THUMB_CMD_, name, dest, bg).Output()
		if err != nil || strings.Contains(string(out), "ERROR") {
			logger.Debugf("ERROR: Generate Gtk Thumbnail\n")
		}
	}
}

func doGenIconThumb(forceFlag bool, outDir string) {
	list := getIconList()
	for _, l := range list {
		dest := getThumbCachePath("--icon", l.Path, outDir)
		if len(dest) < 1 || (!forceFlag && dutils.IsFileExist(dest)) {
			continue
		}
		bg := getThumbBg()
		if !genIconThumbnail(l, dest, bg) {
			logger.Debugf("ERROR: Generate Icon Thumbnail\n")
		}
	}
}

func doGenCursorThumb(forceFlag bool, outDir string) {
	list := getCursorList()
	for _, l := range list {
		dest := getThumbCachePath("--cursor", l.Path, outDir)
		if len(dest) < 1 || (!forceFlag && dutils.IsFileExist(dest)) {
			continue
		}
		bg := getThumbBg()
		if !genCursorThumbnail(l, dest, bg) {
			logger.Debugf("ERROR: Generate Cursor Thumbnail\n")
		}
	}
	if dutils.IsFileExist(XCUR2PNG_OUTDIR) {
		os.RemoveAll(XCUR2PNG_OUTDIR)
	}
}

func doGenBgThumb(forceFlag bool, outDir string) {
	list := getBgList()
	for _, l := range list {
		dest := getThumbCachePath("--background", l.Path, outDir)
		if len(dest) < 1 || (!forceFlag && dutils.IsFileExist(dest)) {
			continue
		}
		err := graphic.ThumbnailImage(l.Path, dest, 128, 72, graphic.FormatPng)
		if err != nil {
			logger.Debug("ERROR:", err)
		}
	}
}
