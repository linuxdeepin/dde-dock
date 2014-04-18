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

// #cgo pkg-config: x11 xtst glib-2.0
// #include "grab-xrecord.h"
// #include <stdlib.h>
import "C"

import (
        "dlib/dbus"
        "github.com/BurntSushi/xgb/xproto"
        "github.com/BurntSushi/xgbutil"
        "github.com/BurntSushi/xgbutil/keybind"
        "github.com/BurntSushi/xgbutil/mousebind"
        "github.com/BurntSushi/xgbutil/xevent"
        "strings"
        "unsafe"
)

type GrabManager struct {
        KeyReleaseEvent func(string)
}

const (
        _GRAB_PATH = "/com/deepin/daemon/GrabKey"
        _GRAB_IFC  = "com.deepin.daemon.GrabKey"
)

func (m *GrabManager) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                _BINDING_DEST,
                _GRAB_PATH,
                _GRAB_IFC,
        }
}

func (m *GrabManager) GrabShortcut(wid xproto.Window,
        shortcut, action string) bool {
        if wid == 0 {
                wid = X.RootWin()
        }

        key := getXGBShortcut(formatShortcut(shortcut))
        if len(key) <= 0 {
                return false
        }
        if !grabKeyPress(wid, key) {
                return false
        }
        keyInfo, ok := newKeyCodeInfo(key)
        if !ok {
                return false
        }
        GrabKeyBinds[keyInfo] = action

        return true
}

func (m *GrabManager) UngrabShortcut(wid xproto.Window,
        shortcut string) bool {
        if wid == 0 {
                wid = X.RootWin()
        }

        key := getXGBShortcut(formatShortcut(shortcut))
        return ungrabKey(wid, key)
}

func (m *GrabManager) GrabKeyboard() {
        go func() {
                X, err := xgbutil.NewConn()
                if err != nil {
                        logObj.Info("Get New Connection Failed:", err)
                        return
                }
                keybind.Initialize(X)
                mousebind.Initialize(X)

                err = keybind.GrabKeyboard(X, X.RootWin())
                if err != nil {
                        logObj.Info("Grab Keyboard Failed:", err)
                        return
                }

                GrabAllButton(X)

                xevent.ButtonPressFun(
                        func(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {
                                m.KeyReleaseEvent("")
                                UngrabAllButton(X)
                                keybind.UngrabKeyboard(X)
                                logObj.Info("Button Press Event")
                                xevent.Quit(X)
                        }).Connect(X, X.RootWin())

                xevent.KeyReleaseFun(
                        func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
                                modStr := keybind.ModifierString(e.State)
                                keyStr := strings.ToLower(
                                        keybind.LookupString(X,
                                                e.State, e.Detail))
                                if e.Detail == 65 {
                                        keyStr = "space"
                                }
                                value := ""
                                if len(modStr) > 0 {
                                        value = ConvertKeyFromMod(filterModStr(modStr)) + keyStr
                                } else {
                                        value = keyStr
                                }
                                m.KeyReleaseEvent(value)
                                UngrabAllButton(X)
                                keybind.UngrabKeyboard(X)
                                logObj.Infof("Key: %s\n", value)
                                xevent.Quit(X)
                        }).Connect(X, X.RootWin())

                xevent.Main(X)
        }()
}

func GrabAllButton(X *xgbutil.XUtil) {
        mousebind.Grab(X, X.RootWin(), 0, 1, false)
        mousebind.Grab(X, X.RootWin(), 0, 2, false)
        mousebind.Grab(X, X.RootWin(), 0, 3, false)
}

func UngrabAllButton(X *xgbutil.XUtil) {
        mousebind.Ungrab(X, X.RootWin(), 0, 1)
        mousebind.Ungrab(X, X.RootWin(), 0, 2)
        mousebind.Ungrab(X, X.RootWin(), 0, 3)
}

func ConvertKeyFromMod(mod string) string {
        values := ""
        vals := strings.Split(mod, "-")
        for _, v := range vals {
                if v == "mod1" || v == "mod2" ||
                        v == "mod4" || v == "lock" {
                        t, ok := _ModKeyMap[v]
                        if !ok {
                                logObj.Info("Get Key Failed From Modify")
                                return ""
                        }
                        values += t + "-"
                } else {
                        values += v + "-"
                }
        }

        return values
}

func (m *GrabManager) GrabSingleKey(key, action string) {
        GrabXRecordKey(key, action)
}

func GrabXRecordKey(key, action string) {
        if len(action) <= 0 {
                logObj.Info("action is null")
                return
        }

        mod, keys, err := keybind.ParseString(X, key)
        if err != nil {
                logObj.Info("ParseString Failed:", err)
                return
        }

        logObj.Infof("mod: %d, key: %d\n", mod, keys[0])
        if mod > 0 {
                logObj.Infof("Not single key\n")
                return
        }

        tmp := C.CString(action)
        defer C.free(unsafe.Pointer(tmp))
        C.grab_xrecord_key(C.int(keys[0]), tmp)
}

func (m *GrabManager) UngrabSingleKey(key string) {
        UngrabXRecordKey(key)
}

func UngrabXRecordKey(key string) {
        mod, keys, err := keybind.ParseString(X, key)
        if err != nil {
                logObj.Info("ParseString Failed:", err)
                return
        }

        if mod > 0 {
                logObj.Infof("Not single key\n")
                return
        }

        C.ungrab_xrecord_key(C.int(keys[0]))
}

/*
func (m *GrabManager) GrabSingleFinalize() {
	C.grab_xrecord_finalize()
}
*/
