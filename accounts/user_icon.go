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
        "dlib/glib-2.0"
        "fmt"
        "io"
        "io/ioutil"
        "os"
        "strings"
)

const (
        _SYSTEM_ICON_PATH  = "/var/lib/AccountsService/icons/"
        _USER_ICON_PATH    = "/.config/deepin-system-settings/account/icons/"
        _ICON_HISTORY_FILE = "/.config/deepin-system-settings/account/account_icon_history.ini"
        _GROUP_KEY         = "Icon History"
)

func (u *User) AllAccountsIcons() []string {
        icons := AllSystemIcons()
        tmp := AllUserIcons()

        for _, v := range tmp {
                icons = append(icons, v)
        }

        return icons
}

func (u *User) AllHistoryIcons() []string {
        homeDir := os.Getenv("HOME")
        if !CheckDirExist(homeDir + _USER_ICON_PATH) {
                return nil
        }

        if !IsFileExist(homeDir + _ICON_HISTORY_FILE) {
                return nil
        }

        keyFile := glib.NewKeyFile()
        defer keyFile.Free()
        _, err := keyFile.LoadFromFile(homeDir+_ICON_HISTORY_FILE,
                glib.KeyFileFlagsNone)
        if err != nil {
                fmt.Println("Load Icon History File Failed:", err)
                return nil
        }

        _, keyArrays, err1 := keyFile.GetKeys(_GROUP_KEY)
        if err1 != nil {
                fmt.Println("Get Group Keys Failed:", err1)
                return nil
        }

        icons := []string{}
        for _, v := range keyArrays {
                value, err2 := keyFile.GetInteger(_GROUP_KEY, v)
                if err2 != nil {
                        fmt.Printf("Get Valus '%s' Failed: %s\n",
                                v, err2)
                        continue
                }

                if value > 0 {
                        icons = append(icons, v)
                }
        }

        return icons
}

func (u *User) AddIconToUserDir(filename string) bool {
        if len(filename) <= 0 {
                return false
        }

        state, err := os.Stat(filename)
        if err != nil {
                fmt.Printf("Get '%s' info failed: %s\n", filename, err)
                return false
        }
        if !state.Mode().IsRegular() {
                return false
        }

        strs := strings.Split(filename, "/")
        homeDir := os.Getenv("HOME")
        destFile := homeDir + _USER_ICON_PATH + strs[len(strs)-1]
        _, err = CopyFile(filename, destFile)
        if err != nil {
                fmt.Println("Copy File Failed:", err)
                return false
        }

        AddIconToHistory(destFile)

        return true
}

func (u *User) DeleteUserIcon(filename string) bool {
        homeDir := os.Getenv("HOME")
        keyFile := glib.NewKeyFile()
        defer keyFile.Free()
        _, err := keyFile.LoadFromFile(homeDir+_ICON_HISTORY_FILE,
                glib.KeyFileFlagsNone)
        if err != nil {
                fmt.Println("Load Icon History File Failed:", err)
                return false
        }

        keyFile.SetInteger(_GROUP_KEY, filename, 0)
        _, data, err2 := keyFile.ToData()
        if err2 != nil {
                fmt.Println("Key File To Data Failed:", err2)
                return false
        }

        if !WriteKeyFile(homeDir+_ICON_HISTORY_FILE, data) {
                return false
        }

        if !strings.Contains(filename, _SYSTEM_ICON_PATH) {
                err := os.Remove(filename)
                if err != nil {
                        fmt.Printf("Remove '%s' failed: %s\n", filename, err)
                        return false
                }
        }

        return true
}

func AddIconToHistory(filename string) bool {
        if len(filename) <= 0 {
                return false
        }

        homeDir := os.Getenv("HOME")
        if !CheckDirExist(homeDir + _USER_ICON_PATH) {
                return false
        }

        if !IsFileExist(homeDir + _ICON_HISTORY_FILE) {
                return false
        }

        keyFile := glib.NewKeyFile()
        defer keyFile.Free()
        _, err := keyFile.LoadFromFile(homeDir+_ICON_HISTORY_FILE,
                glib.KeyFileFlagsNone)
        if err != nil {
                fmt.Println("Load Icon History File Failed:", err)
                return false
        }

        count, err1 := keyFile.GetInteger(_GROUP_KEY, filename)
        if err1 != nil {
                keyFile.SetInteger(_GROUP_KEY, filename, 1)
        } else {
                keyFile.SetInteger(_GROUP_KEY, filename, count+1)
        }

        _, data, err2 := keyFile.ToData()
        if err2 != nil {
                fmt.Println("Key File To Data Failed:", err2)
                return false
        }

        if !WriteKeyFile(homeDir+_ICON_HISTORY_FILE, data) {
                return false
        }

        return true
}

func AllSystemIcons() []string {
        sysFileInfo, err := ioutil.ReadDir(_SYSTEM_ICON_PATH)
        if err != nil {
                fmt.Println("read system icon path failed:", err)
                return nil
        }

        icons := []string{}
        for _, v := range sysFileInfo {
                if v.IsDir() {
                        continue
                }

                icons = append(icons, _SYSTEM_ICON_PATH+v.Name())
        }

        return icons
}

func AllUserIcons() []string {
        homeDir := os.Getenv("HOME")
        if !CheckDirExist(homeDir + _USER_ICON_PATH) {
                return nil
        }

        userFileInfo, err := ioutil.ReadDir(homeDir + _USER_ICON_PATH)
        if err != nil {
                fmt.Println("read user history icon path failed:", err)
                return nil
        }

        icons := []string{}
        for _, v := range userFileInfo {
                if v.IsDir() {
                        continue
                }

                icons = append(icons, homeDir+_USER_ICON_PATH+v.Name())
        }

        return icons
}

func WriteKeyFile(filename, data string) bool {
        if len(filename) <= 0 {
                return false
        }

        f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0664)
        if err != nil {
                fmt.Println("Open History File Failed:", err)
                return false
        }
        defer f.Close()

        _, err = f.WriteString(data)
        if err != nil {
                fmt.Println("Write File Failed:", err)
                return false
        }

        return true
}

func IsFileExist(filename string) bool {
        _, err := os.Stat(filename)
        if os.IsNotExist(err) {
                f, err1 := os.Create(filename)
                if err1 != nil {
                        fmt.Println("Create File Failed:", err1)
                        return false
                }
                defer f.Close()
        }

        return true
}

func CheckDirExist(path string) bool {
        err := os.MkdirAll(path, os.FileMode(0775))
        if err != nil {
                fmt.Println("CheckDirExist Failed:", err)
                return false
        }

        return true
}

func CopyFile(src, dest string) (int64, error) {
        srcFile, err1 := os.Open(src)
        if err1 != nil {
                return -1, err1
        }
        defer srcFile.Close()

        destFile, err2 := os.Create(dest)
        if err2 != nil {
                return -1, err2
        }
        defer destFile.Close()

        return io.Copy(destFile, srcFile)
}
