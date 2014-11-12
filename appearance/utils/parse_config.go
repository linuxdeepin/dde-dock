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

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"regexp"
	"strings"
)

const (
	regularFileMode = 0644
	dirFileMode     = 0755
)

var (
	errInvalidArgs  = fmt.Errorf("Invalid args")
	errInvalidKey   = fmt.Errorf("Invalid line key")
	errWriteKeyFile = fmt.Errorf("Write key to keyfile failed")
)

func GetUserGtk3Config() string {
	dir := path.Join(dutils.GetConfigDir(), "gtk-3.0")
	if !dutils.IsFileExist(dir) {
		err := os.MkdirAll(dir, dirFileMode)
		if err != nil {
			return ""
		}
	}
	return path.Join(dir, "settings.ini")
}

func GetUserGtk2Config() string {
	return path.Join(os.Getenv("HOME"), ".gtkrc-2.0")
}

func GetUserQt4Config() string {
	dir := dutils.GetConfigDir()
	return path.Join(dir, "Trolltech.conf")
}

func WriteUserGtk3Config(filename, key, value string) error {
	if len(filename) == 0 || len(key) == 0 ||
		len(value) == 0 {
		return errInvalidArgs
	}

	ok := dutils.WriteKeyToKeyFile(filename, "Settings", key, value)
	if !ok {
		return errWriteKeyFile
	}

	return nil
}

func WriteUserGtk2Config(filename, key, value string) error {
	if len(filename) == 0 || len(key) == 0 ||
		len(value) == 0 {
		return errInvalidArgs
	}

	var line string
	switch key {
	case "gtk-theme-name", "gtk-icon-theme-name",
		"gtk-font-name", "gtk-cursor-theme-name",
		"gtk-xft-hintstyle", "gtk-xft-rgba":
		line = key + "=\"" + value + "\""
	default:
		line = key + "=" + value
	}

	if !dutils.IsFileExist(filename) {
		return ioutil.WriteFile(filename, []byte(line), regularFileMode)
	}

	return writeLineToFile(filename, "^"+key+"=", line)
}

func writeLineToFile(filename, key, value string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var (
		tmpStr string
		found  bool
	)
	lines := strings.Split(string(contents), "\n")
	for i, line := range lines {
		if i != 0 {
			tmpStr += "\n"
		}
		if ok, _ := regexp.MatchString(key, line); ok {
			tmpStr += value
			found = true
			continue
		}
		tmpStr += line
	}

	if !found {
		return errInvalidKey
	}

	return ioutil.WriteFile(filename, []byte(tmpStr), regularFileMode)
}
