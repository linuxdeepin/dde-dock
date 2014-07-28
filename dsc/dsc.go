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

package dsc

import (
	"os/exec"
	"path"
	"pkg.linuxdeepin.com/lib/log"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"time"
)

var logger = log.NewLogger("dde-session-daemon/dsc")

var quitFlag = make(chan bool)

const (
	DSC_CONFIG_PATH = ".config/deepin-software-center/config_info.ini"
)

func setDSCAutoUpdate(interval time.Duration) {
	logger.Info("AutoUpgrade interval:", interval)
	if interval <= 0 {
		return
	}

	timer := time.After(interval)
	select {
	case <-quitFlag:
		return
	case <-timer:
		go exec.Command("/usr/bin/dsc-daemon", []string{"--no-daemon"}...).Run()
		logger.Info("Running dsc-daemon.....")
		return
	}
}

func getDscConfInfo() (isUpdate bool, duration int32) {
	isUpdate = true
	duration = 3

	homeDir := dutils.GetHomeDir()
	filename := path.Join(homeDir, DSC_CONFIG_PATH)
	if !dutils.IsFileExist(filename) {
		return
	}

	str, _ := dutils.ReadKeyFromKeyFile(filename,
		"update", "auto", "")
	if v, ok := str.(string); ok && v == "False" {
		isUpdate = false
	}

	interval, ok1 := dutils.ReadKeyFromKeyFile(filename,
		"update", "interval", int32(0))
	if ok1 {
		if i, ok := interval.(int32); ok {
			duration = int32(i)
		}
	}

	return

}

func Start() {
	go setDSCAutoUpdate(time.Duration(time.Minute * 5))

	go func() {
		for {
			isUpdate, duration := getDscConfInfo()
			if isUpdate {
				setDSCAutoUpdate(time.Hour * time.Duration(duration))
			} else {
				return
			}
		}
	}()
}

func Stop() {
	quitFlag <- true
}
