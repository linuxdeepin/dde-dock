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
        "github.com/BurntSushi/xgbutil"
        "github.com/BurntSushi/xgbutil/keybind"
        "github.com/BurntSushi/xgbutil/xevent"
        "os/exec"
        "strings"
)

func getSystemPairs() map[string]string {
        systemPairs := make(map[string]string)
        for i, k := range SystemIdNameMap {
                if i >= 0 && i < 300 {
                        if isInvalidConflict(i) {
                                continue
                        }
                        shortcut := getSystemValue(k, false)
                        action := getSystemValue(k, true)
                        systemPairs[shortcut] = action
                }
        }

        SystemPrevPairs = systemPairs
        return systemPairs
}

func getCustomPairs() map[string]string {
        customPairs := make(map[string]string)
        customList := getCustomList()

        for _, i := range customList {
                if isInvalidConflict(i) {
                        continue
                }
                gs := newGSettingsById(i)
                if gs == nil {
                        continue
                }
                shortcut := gs.GetString(_CUSTOM_KEY_SHORTCUT)
                action := gs.GetString(_CUSTOM_KEY_ACTION)
                customPairs[shortcut] = action
        }

        CustomPrevPairs = customPairs
        return customPairs
}

func grabKeyPairs(accelPairs map[string]string, grab bool) {
        for k, v := range accelPairs {
                if len(k) <= 0 {
                        continue
                }
                if k == "super" {
                        if grab {
                                GrabXRecordKey("Super_L", v)
                                GrabXRecordKey("Super_R", v)
                        } else {
                                UngrabXRecordKey("Super_L")
                                UngrabXRecordKey("Super_R")
                        }
                        continue
                }

                shortcut := getXGBShortcut(formatShortcut(k))
                keyInfo := newKeyCodeInfo(shortcut)
                if keyInfo == nil {
                        fmt.Printf("Failed: key: %s, value: %s\n", k, v)
                        continue
                }

                if grab {
                        if grabKeyPress(X.RootWin(), shortcut) {
                                GrabKeyBinds[keyInfo] = v
                        }
                } else {
                        ungrabKey(X.RootWin(), shortcut)
                        delete(GrabKeyBinds, keyInfo)
                }
        }
}

func grabKeyPress(wid xproto.Window, shortcut string) bool {
        if len(shortcut) <= 0 {
                fmt.Println("grab key args failed...")
                return false
        }

        mod, keys, err := keybind.ParseString(X, shortcut)
        if err != nil {
                fmt.Printf("grab parse shortcut string failed: %s\n", err)
                return false
        }

        if len(keys) < 1 {
                fmt.Printf("'%s' no details\n", shortcut)
                return false
        }

        err = keybind.GrabChecked(X, wid, mod, keys[0])
        if err != nil {
                fmt.Printf("Grab '%s' Failed: %s\n", shortcut, err)
                return false
        }

        return true
}

func ungrabKey(wid xproto.Window, shortcut string) bool {
        if len(shortcut) <= 0 {
                return false
        }

        mod, keys, err := keybind.ParseString(X, shortcut)
        if err != nil {
                fmt.Printf("ungrab parse shortcut string failed: %s\n", err)
                return false
        }

        if len(keys) < 1 {
                fmt.Printf("'%s' no details\n", shortcut)
                return false
        }

        keybind.Ungrab(X, wid, mod, keys[0])

        return true
}

func getExecAction(k1 *KeyCodeInfo) (bool, string) {
        for k, v := range GrabKeyBinds {
                if k1.State == k.State && k.Detail == k1.Detail {
                        return true, v
                }
        }

        return false, ""
}

func execCommand(value string) {
        var cmd *exec.Cmd
        vals := strings.Split(value, " ")
        l := len(vals)

        if l > 0 {
                args := []string{}
                for i := 1; i < l; i++ {
                        args = append(args, vals[i])
                }
                /*fmt.Println("args: ", args)*/
                cmd = exec.Command(vals[0], args...)
        } else {
                cmd = exec.Command(value)
        }
        _, err := cmd.Output()
        if err != nil {
                fmt.Println("Exec command failed:", err)
        }
}

func listenKeyPressEvent() {
        xevent.KeyPressFun(
                func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
                        fmt.Printf("State: %d, Detail: %d\n",
                                e.State, e.Detail)
                        modStr := keybind.ModifierString(e.State)
                        keyStr := keybind.LookupString(X, e.State, e.Detail)
                        value := ""
                        if len(modStr) > 0 {
                                value = modStr + "-" + keyStr
                        } else {
                                value = keyStr
                        }
                        fmt.Printf("%s pressed...\n", value)
                        tmpInfo := newKeyCodeInfo(value)
                        if ok, v := getExecAction(tmpInfo); ok {
                                // 不然按键会阻塞，直到程序推出
                                go execCommand(v)
                        }
                }).Connect(X, X.RootWin())
}
