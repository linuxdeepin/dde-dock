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

package themes

// #cgo pkg-config: glib-2.0
// #cgo CFLAGS: -Wall -g
// #include <stdlib.h>
// #include "user_dir.h"
import "C"

import (
	"crypto/md5"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"strconv"
	"strings"
)

var thumbTool = false

func encodeURI(str string) string {
	filepath := dutils.URIToPath(str)
	if len(filepath) < 1 {
		return ""
	}

	u := url.URL{}
	u.Path = filepath
	v := dutils.PathToURI(u.String(), dutils.SCHEME_FILE)
	return v
}

func decodeURI(str string) string {
	u, err := url.Parse(str)
	if err != nil {
		Logger.Warningf("Url parse '%s' failed: %v", str, err)
		return ""
	}

	v := u.Scheme + "://" + u.Path
	return v
}

func changeUserThemeDir() {
	filename := path.Join(homeDir, PERSON_LOCAL_THEME_PATH, "test-dir")
	if !dutils.IsFileExist(filename) {
		os.MkdirAll(filename, 0755)
	} else {
		rmAllFile(filename)
	}
}

func rmAllFile(name string) {
	name = dutils.URIToPath(name)
	os.RemoveAll(name)
}

func getUserPictureDir() string {
	str := C.get_user_pictures_dir()
	//defer C.free(unsafe.Pointer(str))

	ret := C.GoString(str)

	Logger.Debug("User Pictures Dir:", ret)
	return ret
}

func convertMd5ByteToStr(bytes [16]byte) string {
	str := ""

	for _, b := range bytes {
		s := strconv.FormatInt(int64(b), 16)
		if len(s) == 1 {
			str += "0" + s
		} else {
			str += s
		}
	}

	return str
}

func getStrMd5(str string) (string, bool) {
	if len(str) < 1 {
		return "", false
	}

	md5Byte := md5.Sum([]byte(str))
	md5Str := convertMd5ByteToStr(md5Byte)
	if len(md5Str) < 32 {
		return "", false
	}

	return md5Str, true
}

func getFileMd5(file string) (string, bool) {
	if !dutils.IsFileExist(file) {
		return "", false
	}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		Logger.Warningf("ReadFile '%s' failed: %v", file, err)
		return "", false
	}

	md5Byte := md5.Sum(contents)
	md5Str := convertMd5ByteToStr(md5Byte)
	if len(md5Str) < 32 {
		return "", false
	}

	return md5Str, true
}

func isStrInList(str string, list []string) bool {
	for _, l := range list {
		if str == l {
			return true
		}
	}

	return false
}

func isStrListEqual(list1, list2 []string) bool {
	l1 := len(list1)
	l2 := len(list2)

	if l1 != l2 {
		return false
	}

	for i := 0; i < l1; i++ {
		if list1[i] != list2[i] {
			return false
		}
	}

	return true
}

func writeStringToKeyFile(filename, contents string) bool {
	if len(filename) <= 0 {
		return false
	}

	f, err := os.Create(filename + "~")
	if err != nil {
		Logger.Warningf("OpenFile '%s' failed: %v",
			filename+"~", err)
		return false
	}
	defer f.Close()

	if _, err = f.WriteString(contents); err != nil {
		Logger.Warningf("WriteString '%s' failed: %v",
			filename, err)
		return false
	}
	f.Sync()
	os.Rename(filename+"~", filename)

	return true
}

func genThemeThumb() {
	if thumbTool {
		Logger.Infof("There has a theme-thumb-tool running")
		return
	}

	go func() {
		thumbTool = true
		cmd := "/usr/lib/deepin-daemon/theme-thumb-tool"
		m := GetManager()
		exec.Command(cmd, "--background").Run()
		m.setPropBackgroundList(m.getBgStrList())
		//exec.Command(cmd, "--gtk").Run()
		//m.setPropGtkThemeList(m.getGtkStrList())
		//exec.Command(cmd, "--icon").Run()
		//m.setPropIconThemeList(m.getIconStrList())
		//exec.Command(cmd, "--cursor").Run()
		//m.setPropCursorThemeList(m.getCursorStrList())
		thumbTool = false
	}()
}

func sortByDeepin(list []string) []string {
	deepinList := []string{}
	tmpList := []string{}

	for _, l := range list {
		t := strings.ToLower(l)
		if strings.Contains(t, "deepin") {
			deepinList = append(deepinList, l)
			continue
		}

		tmpList = append(tmpList, l)
	}

	deepinList = append(deepinList, tmpList...)

	return deepinList
}
