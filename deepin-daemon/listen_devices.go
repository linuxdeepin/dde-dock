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
        "dbus/com/deepin/daemon/inputdevices"
        "dlib/dbus"
        "github.com/howeyc/fsnotify"
        "io/ioutil"
        "strings"
)

//var (
//watcher *fsnotify.Watcher
//)

const (
        DEV_DIR_PATH     = "/dev/input/by-path"
        DEVICE_FILE_PATH = "/proc/bus/input/devices"
        PROC_KEY_NAME    = "N: Name"

        INPUT_DEV_DEST = "com.deepin.daemon.InputDevices"
        INPUT_DEV_PATH = "/com/deepin/daemon/InputDevices"
        TPAD_PATH_KEY  = "/com/deepin/daemon/InputDevice/TouchPad"
        MOUSE_KEY      = "mouse"
)

func listenDevices() {
        if !utilsObj.IsFileExist(DEV_DIR_PATH) {
                logObj.Warningf("File: '%s' not exist in listenDevices!",
                        DEV_DIR_PATH)
                return
        }

        watcher, err := fsnotify.NewWatcher()
        if err != nil {
                logObj.Warning("New Watch Failed: ", err)
                return
        }

        if err := watcher.Watch(DEV_DIR_PATH); err != nil {
                logObj.Warningf("Listen file: '%s' failed: %v",
                        DEV_DIR_PATH, err)
                return
        }

        go func() {
                defer watcher.Close()
                for {
                        select {
                        case ev := <-watcher.Event:
                                logObj.Info("Watch Event: ", ev)
                                enableTouchPad()
                        case err := <-watcher.Error:
                                logObj.Info("Watch Error: ", err)
                                break
                        }
                }
        }()
}

func getDeviceNames() ([]string, bool) {
        if !utilsObj.IsFileExist(DEVICE_FILE_PATH) {
                logObj.Warningf("Device File: '%s' not exist!",
                        DEVICE_FILE_PATH)
                return []string{}, false
        }

        contents, err := ioutil.ReadFile(DEVICE_FILE_PATH)
        if err != nil {
                logObj.Warningf("Read file: '%s' failed: %v",
                        DEVICE_FILE_PATH, err)
                return []string{}, false
        }

        lines := strings.Split(string(contents), "\n")
        nameList := []string{}
        for _, line := range lines {
                if strings.Contains(line, PROC_KEY_NAME) {
                        nameList = append(nameList, strings.ToLower(line))
                }
        }

        return nameList, true
}

func enableTouchPad() {
        names, ok := getDeviceNames()
        if !ok {
                return
        }
        if utilsObj.IsElementExist(MOUSE_KEY, names) {
                return
        }

        inputObj, err := inputdevices.NewInputDevices(INPUT_DEV_DEST,
                INPUT_DEV_PATH)
        if err != nil {
                logObj.Warning("New Input Devices Failed: ", err)
                return
        }

        devList := inputObj.DevInfoList.Get()
        for _, dev := range devList {
                path := string(dev[0].(string))
                if strings.Contains(path, TPAD_PATH_KEY) {
                        tpadObj, err := inputdevices.NewTouchPad(
                                INPUT_DEV_DEST, dbus.ObjectPath(path))
                        if err != nil {
                                logObj.Warning("New TPad Failed: ", err)
                                return
                        }
                        tpadObj.TPadEnable.Set(true)
                }
        }
}
