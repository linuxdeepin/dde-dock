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

import (
	"crypto/md5"
	"io/ioutil"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"strconv"
	"strings"
)

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
		Logger.Debugf("ReadFile '%s' failed: %v", file, err)
		return "", false
	}

	md5Byte := md5.Sum(contents)
	md5Str := convertMd5ByteToStr(md5Byte)
	if len(md5Str) < 32 {
		return "", false
	}

	return md5Str, true
}

func isStrInList(str string, list []string) (string, bool) {
	for _, l := range list {
		if strings.Contains(l, str) {
			return l, true
		}
	}
	return "", false
}
