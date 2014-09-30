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
	"os"
	"os/exec"
	"path"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	GtkThumbSeed        = "--gtk"
	IconThumbSeed       = "--icon"
	CursorThumbSeed     = "--cursor"
	BackgroundThumbSeed = "--background"
)

const (
	userThumbFileDir   = ".local/share/personalization/thumbnail/autogen"
	systemThumbFileDir = "/usr/share/personalization/thumbnail/autogen"
)

const genThumbCommand = "/usr/lib/deepin-daemon/theme-thumb-tool"

var _process *os.Process

func GenerateThumbnail() error {
	if _process != nil {
		return nil
	}

	cmd := exec.Command("/bin/sh", "-c",
		genThumbCommand+" -a")
	err := cmd.Start()
	if err != nil {
		return err
	}

	_process = cmd.Process
	return nil
}

func GetThumbnail(seed, uri string) string {
	if !isThumbSeedValid(seed) {
		return ""
	}

	src := dutils.DecodeURI(uri)
	dir := path.Join(os.Getenv("HOME"), userThumbFileDir)
	md5Str, _ := dutils.SumStrMd5(seed + src)
	dest := path.Join(dir, md5Str+".png")
	if dutils.IsFileExist(dest) {
		//return dutils.EncodeURI(dest, dutils.SCHEME_FILE)
		return dest
	}

	dest = path.Join(systemThumbFileDir, md5Str+".png")
	if dutils.IsFileExist(dest) {
		//return dutils.EncodeURI(dest, dutils.SCHEME_FILE)
		return dest
	}

	return ""
}

func isThumbSeedValid(seed string) bool {
	switch seed {
	case GtkThumbSeed, IconThumbSeed,
		CursorThumbSeed, BackgroundThumbSeed:
		return true
	}

	return false
}
