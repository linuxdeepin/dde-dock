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
        "io/ioutil"
        "strings"
)

type Manager struct {
        Infos []deviceInfo
}

type deviceInfo struct {
        Path    string
        Id      string
}

const (
        _PROC_DEVICE_PATH = "/proc/bus/input/devices"
        _PROC_KEY_NAME    = "N: Name"
)

func getDeviceNames() []string {
        names := []string{}

        contents, err := ioutil.ReadFile(_PROC_DEVICE_PATH)
        if err != nil {
                logObj.Warningf("ReadFile '%s' failed: %v",
                        _PROC_DEVICE_PATH, err)
                return names
        }

        lines := strings.Split(string(contents), "\n")
        for _, line := range lines {
                if strings.Contains(line, _PROC_KEY_NAME) {
                        names = append(names, strings.ToLower(line))
                }
        }

        return names
}

func NewManager() *Manager {
        m := &Manager{}

        names := getDeviceNames()
        tmps := []deviceInfo{}
        for _, name := range names {
                if strings.Contains(name, "mouse") {
                        info := deviceInfo{DEVICE_PATH + "Mouse", "mouse"}
                        tmps = append(tmps, info)
                } else if strings.Contains(name, "touchpad") {
                        info := deviceInfo{DEVICE_PATH + "TouchPad", "touchpad"}
                        tmps = append(tmps, info)
                } else if strings.Contains(name, "keyboard") {
                        info := deviceInfo{DEVICE_PATH + "Keyboard", "keyboard"}
                        tmps = append(tmps, info)
                }
        }

        m.Infos = tmps

        return m
}
