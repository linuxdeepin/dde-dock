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
        "dlib/dbus"
        "dlib/gio-2.0"
        "fmt"
        "github.com/BurntSushi/xgbutil"
        "github.com/BurntSushi/xgbutil/keybind"
        "github.com/BurntSushi/xgbutil/xevent"
)

type Manager struct {
        MediaKeyList    []string
        AccelKeyChanged func(string, string)
}

const (
        KEY_PRESS   = "Press"
        KEY_RELEASE = "Release"

        MEDIA_DEST = "com.deepin.daemon.MediaKey"
        MEDIA_PATH = "/com/deepin/daemon/MediaKey"
        MEDIA_IFC  = "com.deepin.daemon.MediaKey"
)

func (op *Manager) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                MEDIA_DEST,
                MEDIA_PATH,
                MEDIA_IFC,
        }
}

func (op *Manager) listenKeyChanged() {
        xevent.KeyPressFun(
                func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
                        fmt.Printf("State: %d, Detail: %d\n",
                                e.State, e.Detail)
                        modStr := keybind.ModifierString(e.State)
                        keyStr := keybind.LookupString(X, e.State, e.Detail)
                        accel := ""
                        if len(modStr) > 0 {
                                accel = modStr + "-" + keyStr
                        } else {
                                accel = keyStr
                        }
                        fmt.Println("Accel Key:", accel)
                        op.AccelKeyChanged(KEY_PRESS, convertModToKey(accel))
                }).Connect(X, X.RootWin())

        /*
           xevent.KeyReleaseFun(
                   func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
                           fmt.Printf("State: %d, Detail: %d\n",
                                   e.State, e.Detail)
                           modStr := keybind.ModifierString(e.State)
                           keyStr := keybind.LookupString(X, e.State, e.Detail)
                           accel := ""
                           if len(modStr) > 0 {
                                   accel = modStr + "-" + keyStr
                           } else {
                                   accel = keyStr
                           }
                           fmt.Println("Accel Key:", accel)
                           op.AccelKeyChanged(KEY_RELEASE, accel)
                   }).Connect(X, X.RootWin())
        */
}

func (op *Manager) listenMediaKeyChanged() {
        mediaSettings.Connect("changed", func(s *gio.Settings, key string) {
                value := mediaSettings.GetString(key)
                v := mediaKeyMap[key]
                if v != value {
                        op.UnregisterAccelKey(v)
                        op.RegisterAccelKey(value)
                        mediaKeyMap[key] = value
                }
        })
}
