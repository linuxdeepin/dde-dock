/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher

import (
	"gir/glib-2.0"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/xdg/basedir"
	"strings"
)

const (
	AppDirName                 = "applications"
	DirDefaultPerm os.FileMode = 0755
)

func isZH() bool {
	lang := gettext.QueryLang()
	return strings.HasPrefix(lang, "zh")
}

func getUserDesktopDir() string {
	return glib.GetUserSpecialDir(glib.UserDirectoryDirectoryDesktop)
}

// return $HOME/.local/share/applications
func getUserAppDir() string {
	userDataDir := basedir.GetUserDataDir()
	return filepath.Join(userDataDir, AppDirName)
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

// SaveKeyFile saves key file.
func SaveKeyFile(file *glib.KeyFile, path string) error {
	_, content, err := file.ToData()
	if err != nil {
		return err
	}

	stat, err := os.Lstat(path)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, []byte(content), stat.Mode())
	if err != nil {
		return err
	}
	return nil
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
