/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

// #cgo pkg-config: glib-2.0 gdk-pixbuf-2.0
// #cgo LDFLAGS: -lm
// #include <stdlib.h>
// #include "blur-pict.h"
import "C"
import "unsafe"

import (
	"crypto/md5"
	"dbus/org/freedesktop/accounts"
	"dlib/dbus"
	"fmt"
	"os"
	"os/user"
	"strconv"
)

type BlurPictManager struct {
	BlurPictChanged func(string, string)
}

type BlurResult struct {
	Success  bool
	PictPath string
}

const (
	_BLUR_PICT_DEST = "com.deepin.daemon.BlurPictManager"
	_BLUR_PICT_PATH = "/com/deepin/daemon/BlurPictManager"
	_BLUR_PICT_IFC  = "com.deepin.daemon.BlurPictManager"

	_ACCOUNTS_PATH          = "/org/freedesktop/Accounts/User"
	_BG_BLUR_PICT_CACHE_DIR = "gaussian-background"
)

func (blur *BlurPictManager) BlurPictDestPath(uid string, sigma float64, numsteps int64) *BlurResult {
	homeDir, err := GetHomeDirById(uid)
	if err != nil {
		fmt.Println("get home dir failed")
		return &BlurResult{Success: false, PictPath: ""}
	}

	srcPath := GetCurrentSrcPath(uid)
	destPath := GenerateDestPath(srcPath, homeDir)
	if IsFileValid(srcPath, destPath) {
		return &BlurResult{Success: true, PictPath: destPath}
	}

	go GenerateBlurPict(blur, srcPath, destPath, uid, sigma, numsteps)

	return &BlurResult{Success: false, PictPath: srcPath}
}

func (blur *BlurPictManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_BLUR_PICT_DEST,
		_BLUR_PICT_PATH,
		_BLUR_PICT_IFC,
	}
}

func GetHomeDirById(uid string) (string, error) {
	users, err := user.LookupId(uid)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return users.HomeDir, nil
}

func GetCurrentSrcPath(uid string) string {
	userPath := _ACCOUNTS_PATH + uid
	accountsUser := accounts.GetUser(userPath)
	bgPath := accountsUser.BackgroundFile
	srcPath := bgPath.GetValue().(string)

	return srcPath
}

func GenerateBlurPict(blur *BlurPictManager, srcPath, destPath, uid string, sigma float64, numsteps int64) bool {
	if len(uid) <= 0 && len(srcPath) <= 0 && len(destPath) <= 0 {
		fmt.Println("args failed")
		return false
	}

	homeDir, err := GetHomeDirById(uid)
	if err != nil {
		fmt.Println("get home dir failed")
		return false
	}

	src := C.CString(srcPath)
	defer C.free(unsafe.Pointer(src))
	dest := C.CString(destPath)
	defer C.free(unsafe.Pointer(dest))

	pictPath := homeDir + "/.cache/" + _BG_BLUR_PICT_CACHE_DIR
	err = os.MkdirAll(pictPath, os.FileMode(0700))
	if err != nil {
		fmt.Println(err)
		return false
	}

	is_ok := C.generate_blur_pict(src, dest, C.double(sigma), C.long(numsteps))
	if is_ok == 0 {
		fmt.Println("generate gaussian picture failed")
		return false
	}
	blur.BlurPictChanged(uid, destPath)

	return true
}

func GenerateDestPath(srcPath, homeDir string) string {
	if len(homeDir) <= 0 && len(srcPath) <= 0 {
		fmt.Println("args failed")
		return ""
	}

	md5Sum := md5.Sum([]byte(srcPath))
	md5Str := ""
	for _, b := range md5Sum {
		s := strconv.FormatInt(int64(b), 16)
		if len(s) == 1 {
			md5Str += "0" + s
		} else {
			md5Str += s
		}
	}

	destPath := homeDir + "/.cache/" + _BG_BLUR_PICT_CACHE_DIR + "/" + md5Str + ".png"
	return destPath
}

func IsFileValid(srcPath, destPath string) bool {
	if len(srcPath) <= 0 && len(destPath) <= 0 {
		fmt.Println("args failed")
		return false
	}

	_, err := os.Stat(destPath)
	if os.IsNotExist(err) {
		fmt.Println("file is not exist")
		return false
	}

	src := C.CString(srcPath)
	defer C.free(unsafe.Pointer(src))
	dest := C.CString(destPath)
	defer C.free(unsafe.Pointer(dest))
	if C.blur_pict_is_valid(src, dest) == 0 {
		fmt.Println("file invalid")
		return false
	}

	return true
}

func main() {
	blur := &BlurPictManager{}
	dbus.InstallOnSession(blur)

	select {}
}
