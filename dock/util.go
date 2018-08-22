/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package dock

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pkg.deepin.io/lib/xdg/basedir"
)

var xdgAutostartDirs []string

func init() {
	configDirs := make([]string, 0, 3)
	configDirs = append(configDirs, basedir.GetUserConfigDir())
	sysConfigDirs := basedir.GetSystemConfigDirs()
	configDirs = append(configDirs, sysConfigDirs...)

	for idx, configDir := range configDirs {
		configDirs[idx] = filepath.Join(configDir, "autostart")
	}
	xdgAutostartDirs = configDirs
}

func isInAutostartDir(file string) bool {
	dir := filepath.Dir(file)
	for _, adir := range xdgAutostartDirs {
		if adir == dir {
			return true
		}
	}
	return false
}

func dataUriToFile(dataUri, path string) (string, error) {
	// dataUri starts with string "data:image/png;base64,"
	commaIndex := strings.Index(dataUri, ",")
	img, err := base64.StdEncoding.DecodeString(dataUri[commaIndex+1:])
	if err != nil {
		return path, err
	}

	return path, ioutil.WriteFile(path, img, 0644)
}

func strSliceEqual(sa, sb []string) bool {
	if len(sa) != len(sb) {
		return false
	}
	for i, va := range sa {
		vb := sb[i]
		if va != vb {
			return false
		}
	}
	return true
}

func uniqStrSlice(slice []string) []string {
	newSlice := make([]string, 0)
	for _, e := range slice {
		if !strSliceContains(newSlice, e) {
			newSlice = append(newSlice, e)
		}
	}
	return newSlice
}

func strSliceContains(slice []string, v string) bool {
	for _, e := range slice {
		if e == v {
			return true
		}
	}
	return false
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func getCurrentTimestamp() uint32 {
	return uint32(time.Now().Unix())
}

func toLocalPath(in string) string {
	u, err := url.Parse(in)
	if err != nil {
		return ""
	}
	if u.Scheme == "file" {
		return u.Path
	}
	return in
}
