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
	"dbus/com/deepin/daemon/thememanager"
	"dbus/com/deepin/sessionmanager"
	"dlib/utils"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"unsafe"
)

var (
	objXS   *sessionmanager.XSettings
	objTM   *thememanager.ThemeManager
	objUtil = utils.NewUtils()
)

const (
	_CMD_           = "theme-thumb-tool"
	_GTK_THUMB_CMD_ = "/usr/lib/deepin-daemon/gtk-thumb-tool"
)

const (
	ICON_DEVICE     = "devices/48/block-device.png"
	ICON_PLACE      = "places/48/folder_home.png"
	ICON_APP        = "apps/48/google-chrome.png"
	ICON_DEEPIN_APP = "apps/48/deepin-software-center.png"
)

func genCursorThumbnail(theme, dest, bg string) bool {
	if len(theme) < 1 || len(dest) < 1 || len(bg) < 1 {
		return false
	}

	// Get Current Theme
	curTheme, _, err := objXS.GetString("Gtk/CursorThemeName")
	if err != nil {
		fmt.Printf("Get Current Cusrsor Theme Failed: %v\n", err)
		return false
	}
	// Set CursorTheme To theme
	if theme != curTheme {
		objXS.SetString("Gtk/CursorThemeName", theme)
	}
	cBg := C.CString(bg)
	defer C.free(unsafe.Pointer(cBg))
	cDest := C.CString(dest)
	defer C.free(unsafe.Pointer(cDest))
	cRet := C.gen_cursor_preview(cBg, cDest)
	// Set CursorTheme To curTheme
	if theme != curTheme {
		objXS.SetString("Gtk/CursorThemeName", curTheme)
	}
	if int(cRet) == -1 {
		return false
	}

	return true
}

/**
 * IconList Struct:
 *	Name string
 *	Path string
 *	Type int32
 *	Thumbnail string
**/
func genIconThumbnail(theme, dest, bg string) bool {
	if len(theme) < 1 || len(dest) < 1 || len(bg) < 1 {
		return false
	}

	item1 := ""
	item2 := ""
	item3 := ""

	home, _ := objUtil.GetHomeDir()
	for _, l := range objTM.IconThemeList.Get() {
		if l == theme {
			if f, err := objTM.GetFlag("icon", theme); err == nil {
				dir := ""
				if f == 0 {
					dir = path.Join("/usr/share/icons", theme)
				} else {
					dir = path.Join(home, ".icons", theme)
				}
				item1 = path.Join(dir, ICON_DEVICE)
				if !objUtil.IsFileExist(item1) {
					return false
				}

				item2 = path.Join(dir, ICON_PLACE)
				if !objUtil.IsFileExist(item2) {
					return false
				}

				item3 = path.Join(dir, ICON_DEEPIN_APP)
				if !objUtil.IsFileExist(item3) {
					item3 = path.Join(dir, ICON_APP)
					if !objUtil.IsFileExist(item3) {
						return false
					}
				}
			}
			break
		}
	}

	if len(item1) < 1 {
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
		return false
	}

	return true
}

func printHelper() {
	fmt.Printf("Name\n\t%s: Theme Thumbnail Tool\n", _CMD_)
	fmt.Printf("Usage\n\t%s [Option] [Theme Name] [Dest Path] [Background]\n", _CMD_)
	fmt.Printf("Options:\n")
	fmt.Printf("\t--gtk:    Generate Gtk Theme Thumbnail\n")
	fmt.Printf("\t--icon:   Generate Icon Theme Thumbnail\n")
	fmt.Printf("\t--cursor: Generate Cursor Theme Thumbnail\n")
}

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("ERROR\n")
		printHelper()
		return
	}

	var err error
	if objTM, err = thememanager.NewThemeManager("com.deepin.daemon.ThemeManager", "/com/deepin/daemon/ThemeManager"); err != nil {
		fmt.Printf("ERROR\n")
		fmt.Printf("New ThemeManager Failed: %v\n", err)
		return
	}

	if objXS, err = sessionmanager.NewXSettings("com.deepin.SessionManager", "/com/deepin/XSettings"); err != nil {
		fmt.Printf("ERROR\n")
		fmt.Printf("New XSettings Failed: %v\n", err)
		return
	}

	C.init_env()

	op := os.Args[1]
	theme := os.Args[2]
	dest := os.Args[3]
	bg := os.Args[4]

	switch op {
	case "--gtk":
		out, err := exec.Command(_GTK_THUMB_CMD_, []string{theme, dest}...).Output()
		if err != nil || strings.Contains(string(out), "ERROR") {
			fmt.Printf("ERROR\n")
		}
	case "--icon":
		if !genIconThumbnail(theme, dest, bg) {
			fmt.Printf("ERROR\n")
		}
	case "--cursor":
		if !genCursorThumbnail(theme, dest, bg) {
			fmt.Printf("ERROR\n")
		}
	default:
		fmt.Printf("ERROR\n")
		fmt.Printf("Invalid option: %s\n\n", op)
		printHelper()
	}
}
