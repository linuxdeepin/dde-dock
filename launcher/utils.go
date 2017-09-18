/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

package launcher

import (
	"os"
	"path"
	"path/filepath"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/xdg/basedir"
	"pkg.deepin.io/lib/xdg/userdir"
	"strings"
)

const (
	AppDirName                 = "applications"
	DirDefaultPerm os.FileMode = 0755
)

// appInDesktop returns the destination when the desktop file is
// sent to the user's desktop direcotry.
func appInDesktop(appId string) string {
	appId = strings.Replace(appId, "/", "-", -1)
	return filepath.Join(getUserDesktopDir(), appId+desktopExt)
}

func isZH() bool {
	lang := gettext.QueryLang()
	return strings.HasPrefix(lang, "zh")
}

func getUserDesktopDir() string {
	return userdir.Get(userdir.Desktop)
}

// return $HOME/.local/share/applications
func getUserAppDir() string {
	userDataDir := basedir.GetUserDataDir()
	return filepath.Join(userDataDir, AppDirName)
}

func getDataDirsForWatch() []string {
	userDataDir := basedir.GetUserDataDir()
	sysDataDirs := basedir.GetSystemDataDirs()
	return append(sysDataDirs, userDataDir)
}

// The default applications module of the DDE Control Center
// creates the desktop file with the file name beginning with
// "deepin-custom" in the applications directory under the XDG
// user data directory.
func isDeepinCustomDesktopFile(file string) bool {
	dir := filepath.Dir(file)
	base := filepath.Base(file)
	userAppDir := getUserAppDir()

	return dir == userAppDir && strings.HasPrefix(base, "deepin-custom-")
}

func getAppDirs() []string {
	dataDirs := basedir.GetSystemDataDirs()
	dataDirs = append(dataDirs, basedir.GetUserDataDir())
	var dirs []string
	for _, dir := range dataDirs {
		dirs = append(dirs, path.Join(dir, AppDirName))
	}
	return dirs
}

func getAppIdByFilePath(file string, appDirs []string) string {
	file = filepath.Clean(file)
	var desktopId string
	for _, dir := range appDirs {
		if strings.HasPrefix(file, dir) {
			desktopId, _ = filepath.Rel(dir, file)
			break
		}
	}
	if desktopId == "" {
		return ""
	}
	return strings.TrimSuffix(desktopId, desktopExt)
}

func runeSliceToStringSlice(runes []rune) []string {
	var list []string
	for _, v := range runes {
		list = append(list, string(v))
	}
	return list
}

func runeSliceDiff(key, current []rune) (popCount int, runesPush []rune) {
	var i int
	kLen := len(key)
	cLen := len(current)
	if kLen == 0 {
		popCount = cLen
		return
	}
	if cLen == 0 {
		runesPush = key
		return
	}

	for {
		k := key[i]
		c := current[i]
		//logger.Debugf("[%v] k %v c %v", i, k, c)

		if k == c {
			i++
			if i == kLen {
				//logger.Debug("i == key len")
				break
			}
			if i == cLen {
				//logger.Debug("i == current len")
				break
			}

		} else {
			break
		}
	}
	popCount = cLen - i
	for j := i; j < kLen; j++ {
		runesPush = append(runesPush, key[j])
	}

	//logger.Debug("i:", i)
	return
}
