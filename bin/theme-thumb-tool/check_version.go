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
	"errors"
	"io/ioutil"
	"os"
	"path"
	dutils "pkg.deepin.io/lib/utils"
	"regexp"
)

const (
	VERSION              = "0.2"
	VERSION_FILE         = "version"
	THUMB_CACHE_SYS_PATH = "/usr/share/personalization/thumbnail/autogen"
)

func isVersionSame() bool {
	versionFile := path.Join(THUMB_CACHE_SYS_PATH, VERSION_FILE)
	if !dutils.IsFileExist(versionFile) {
		err := writeVersionFile(versionFile)
		if err != nil {
			logger.Warning("writeVersionFile failed:", err)
		}
		return false
	}

	contents, err := ioutil.ReadFile(versionFile)
	if err != nil {
		logger.Warning("Read version file failed:", err)
		return false
	}

	if string(contents) == VERSION {
		return true
	}
	if err := writeVersionFile(versionFile); err != nil {
		logger.Warning("writeVersionFile failed:", err)
	}

	return false
}

func writeVersionFile(filename string) error {
	if !dutils.IsFileExist(THUMB_CACHE_SYS_PATH) {
		if err := os.MkdirAll(THUMB_CACHE_SYS_PATH, 0755); err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		logger.Warning("Open version file failed:", err)
		return err
	}
	defer file.Close()

	file.WriteString(VERSION)
	file.Sync()

	return nil
}

func cleanThumbCache() error {
	if !dutils.IsFileExist(THUMB_CACHE_SYS_PATH) {
		return errors.New("file not exist")
	}

	var (
		err   error
		file  *os.File
		infos []os.FileInfo
		reg   *regexp.Regexp
	)

	file, err = os.Open(THUMB_CACHE_SYS_PATH)
	if err != nil {
		logger.Warning("Open thumbnail cache dir failed:", err)
		return err
	}
	defer file.Close()

	infos, err = file.Readdir(0)
	if err != nil {
		return err
	}

	reg, err = regexp.Compile(`(?i)(\.png$)`)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}

		if reg.MatchString(info.Name()) {
			os.Remove(path.Join(THUMB_CACHE_SYS_PATH, info.Name()))
		}
	}

	file.Sync()
	return nil
}
