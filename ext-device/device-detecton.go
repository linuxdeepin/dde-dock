/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	_PROC_DEVICE_PATH = "/proc/bus/input/devices"
	_PROC_DEV_KEY     = "N: Name"
)

func IsFileNotExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return true
	}

	return false
}

func GetProcDeviceNameList() (bool, []string) {
	if IsFileNotExist(_PROC_DEVICE_PATH) {
		fmt.Printf("%s not exist\n", _PROC_DEVICE_PATH)
		return false, []string{}
	}

	contents, err := ioutil.ReadFile(_PROC_DEVICE_PATH)
	if err != nil {
		fmt.Println(err)
		return false, []string{}
	}

	lines := strings.Split(string(contents), "\n")
	nameList := []string{}
	for _, line := range lines {
		if strings.Contains(line, _PROC_DEV_KEY) {
			nameList = append(nameList, line)
		}
	}

	return true, nameList
}

func DeviceIsExist(nameList []string, device string) bool {
	for _, name := range nameList {
		if strings.Contains(name, device) {
			return true
		}
	}

	return false
}
