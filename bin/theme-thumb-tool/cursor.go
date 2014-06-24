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
	dutils "pkg.linuxdeepin.com/lib/utils"
	"os/exec"
	"path"
)

const (
	CMD_XCUR2PNG = "/usr/bin/xcur2png"

	LEFT_PTR       = "cursors/left_ptr"
	LEFT_PTR_WATCH = "cursors/left_ptr_watch"
	QUESTION_ARROW = "cursors/question_arrow"
)

func convertCursor2Png(filename string) string {
	_, err := exec.Command(CMD_XCUR2PNG,
		"-i", "24",
		"-c", "/tmp",
		"-d", "/tmp",
		"-q",
		filename).Output()
	if err != nil {
		return ""
	}

	dest := path.Join("/tmp", path.Base(filename)+"_024.png")
	if dutils.IsFileExist(dest) {
		return dest
	}

	return ""
}

func getCursorFiles(info pathInfo) (string, string, string) {
	s1 := path.Join(info.Path, LEFT_PTR)
	s2 := path.Join(info.Path, LEFT_PTR_WATCH)
	s3 := path.Join(info.Path, QUESTION_ARROW)

	return s1, s2, s3
}

func getCursorIcons(info pathInfo) (string, string, string) {
	s1, s2, s3 := getCursorFiles(info)
	if len(s1) < 1 || len(s2) < 1 || len(s3) < 1 {
		return "", "", ""
	}

	d1 := convertCursor2Png(s1)
	d2 := convertCursor2Png(s2)
	d3 := convertCursor2Png(s3)
	if len(d1) < 1 || len(d2) < 1 || len(d3) < 1 {
		return "", "", ""
	}

	return d1, d2, d3
}
