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
	"dlib/utils"
	"os/exec"
	"path"
	"time"
)

var objUtil = utils.NewUtils()

const (
	DSC_CONFIG_PATH = ".config/deepin-software-center/config_info.ini"
)

func setDSCAutoUpdate(interval time.Duration) {
	if interval <= 0 {
		return
	}

	for {
		timer := time.After(time.Hour * interval)
		select {
		case <-timer:
			go exec.Command("/usr/bin/dsc-daemon", []string{"--no-daemon"}...).Run()
		}
	}
}

func dscAutoUpdate() {
	homeDir, ok := objUtil.GetHomeDir()
	if !ok {
		return
	}
	filename := path.Join(homeDir, DSC_CONFIG_PATH)
	if !objUtil.IsFileExist(filename) {
		return
	}

	interval, ok1 := objUtil.ReadKeyFromKeyFile(filename,
		"update", "interval", int32(0))
	if !ok1 {
		interval = 3
	}
	isUpdate, ok2 := objUtil.ReadKeyFromKeyFile(filename,
		"update", "auto", false)
	if !ok2 {
		isUpdate = true
	}
	if v, ok := isUpdate.(bool); ok && v {
		if i, ok := interval.(int32); ok {
			setDSCAutoUpdate(time.Duration(i))
		}
	}
}
