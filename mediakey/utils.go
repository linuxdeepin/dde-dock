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
        "github.com/BurntSushi/xgb/xproto"
        //"github.com/BurntSushi/xgbutil"
        "github.com/BurntSushi/xgbutil/keybind"
        "strings"
)

var (
        keyToModMap = map[string]string{
                "caps_lock": "lock",
                "alt":       "mod1",
                "meta":      "mod1",
                "num_lock":  "mod2",
                "super":     "mod4",
                "hyper":     "mod4",
        }

        modToKeyMap = map[string]string{
                "mod1": "alt",
                "mod2": "num_lock",
                "mod4": "super",
                "lock": "caps_lock",
        }
)

func grapAccelKey(wid xproto.Window, accel string) bool {
        if len(accel) <= 0 {
                fmt.Println("grapAccelKey accel is null")
                return false
        }

        state, details, err := keybind.ParseString(X, accel)
        if err != nil {
                fmt.Printf("ParseString '%s' failed: %s\n", accel, err)
                //panic(err)
                return false
        }

        if len(details) < 1 {
                fmt.Printf("'%s' no details\n", accel)
                return false
        }

        err = keybind.GrabChecked(X, wid, state, details[0])
        if err != nil {
                fmt.Printf("GrabChecked '%s' failed: %s\n", accel, err)
                //panic(err)
                return false
        }

        return true
}

func ungrabAccelKey(wid xproto.Window, accel string) bool {
        if len(accel) <= 0 {
                fmt.Println("ungrabAccelKey accel is null")
                return false
        }

        state, details, err := keybind.ParseString(X, accel)
        if err != nil {
                fmt.Printf("ParseString '%s' failed: %s\n", accel, err)
                //panic(err)
                return false
        }

        if len(details) < 1 {
                fmt.Printf("'%s' no details\n", accel)
                return false
        }

        keybind.Ungrab(X, wid, state, details[0])
        return true
}

func convertKeyToMod(key string) string {
        if len(key) <= 0 {
                return ""
        }

        if !strings.Contains(key, "-") {
                return key
        }

        strs := strings.Split(key, "-")
        mod, ok := keyToModMap[strs[0]]
        if !ok {
                return key
        }

        tmp := mod + "-" + strs[1]
        fmt.Println("Mod Key:", tmp)
        return tmp
}

func convertModToKey(key string) string {
        if len(key) <= 0 {
                return ""
        }

        if !strings.Contains(key, "-") {
                return key
        }

        strs := strings.Split(key, "-")
        str, ok := modToKeyMap[strs[0]]
        if !ok {
                return key
        }

        tmp := str + "-" + strs[1]
        fmt.Println("String Key:", tmp)
        return tmp
}
