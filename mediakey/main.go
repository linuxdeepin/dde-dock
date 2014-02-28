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
        "dlib"
        "dlib/dbus"
        "dlib/gio-2.0"
        "fmt"
        "github.com/BurntSushi/xgbutil"
        "github.com/BurntSushi/xgbutil/keybind"
        "github.com/BurntSushi/xgbutil/xevent"
)

var (
        X             *xgbutil.XUtil
        mediaSettings *gio.Settings
        mediaKeyMap   map[string]string
)

const (
        MEDIA_KEY_SCHEMA_ID = "com.deepin.dde.key-binding.mediakey"
)

func (op *Manager) RegisterAccelKey(accel string) {
        defer func() {
                if err := recover(); err != nil {
                        fmt.Println("Recover Error In RegisterAccelKey:", err)
                }
        }()

        grapAccelKey(X.RootWin(), convertKeyToMod(accel))
}

func (op *Manager) UnregisterAccelKey(accel string) {
        defer func() {
                if err := recover(); err != nil {
                        fmt.Println("Recover Error In UnregisterAccelKey:", err)
                }
        }()

        ungrabAccelKey(X.RootWin(), convertKeyToMod(accel))
}

func initMediaKey() {
        defer func() {
                if err := recover(); err != nil {
                        fmt.Println("Recover Error In initMediaKey:", err)
                }
        }()

        keyList := mediaSettings.ListKeys()
        for _, key := range keyList {
                value := mediaSettings.GetString(key)
                mediaKeyMap[key] = value
                grapAccelKey(X.RootWin(), value)
        }
}

func (op *Manager) ChangeMediaKey(key, value string) bool {
        defer func() {
                if err := recover(); err != nil {
                        fmt.Println("Recover Error In ChangeMediaKey:", err)
                }
        }()

        v, ok := mediaKeyMap[key]
        if !ok {
                fmt.Printf("'%s' is not in MediaKeyList\n", key)
                return false
        }

        if v == value {
                return false
        }

        op.RegisterAccelKey(value)
        mediaKeyMap[key] = value

        return true
}

func newManager() *Manager {
        m := &Manager{}
        m.MediaKeyList = mediaSettings.ListKeys()
        m.listenKeyChanged()
        m.listenMediaKeyChanged()

        return m
}

func main() {
        defer func() {
                if err := recover(); err != nil {
                        fmt.Println("Recover Error:", err)
                }
        }()

        var err error

        X, err = xgbutil.NewConn()
        if err != nil {
                fmt.Println("New XUtil Connection Failed:", err)
                panic(err)
        }
        keybind.Initialize(X)
        mediaSettings = gio.NewSettings(MEDIA_KEY_SCHEMA_ID)
        mediaKeyMap = make(map[string]string)

        m := newManager()
        initMediaKey()
        dbus.InstallOnSession(m)
        dbus.DealWithUnhandledMessage()

        go dlib.StartLoop()
        xevent.Main(X)
}
