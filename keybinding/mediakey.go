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
        "strings"
)

type MediaKeyManager struct {
        AudioMute        func(bool)
        AudioUp          func(bool)
        AudioDown        func(bool)
        BrightnessUp     func(bool)
        BrightnessDown   func(bool)
        CapsLockOn       func(bool)
        CapsLockOff      func(bool)
        NumLockOn        func(bool)
        NumLockOff       func(bool)
        SwitchMonitors   func(bool)
        TouchPadOn       func(bool)
        TouchPadOff      func(bool)
        PowerOff         func(bool)
        PowerSleep       func(bool)
        SwitchLayout     func(bool)
        AudioPlay        func(bool)
        AudioPause       func(bool)
        AudioStop        func(bool)
        AudioPrevious    func(bool)
        AudioNext        func(bool)
        AudioRewind      func(bool)
        AudioForward     func(bool)
        AudioRepeat      func(bool)
        LaunchEmail      func(bool)
        LaunchBrowser    func(bool)
        LaunchCalculator func(bool)
}

const (
        MEDIA_KEY_DEST = "com.deepin.daemon.KeyBinding"
        MEDIA_KEY_PATH = "/com/deepin/daemon/MediaKey"
        MEDIA_KEY_IFC  = "com.deepin.daemon.MediaKey"

        MEDIA_KEY_SCHEMA_ID = "com.deepin.dde.key-binding.mediakey"
)

var (
        mediaKeySettings = gio.NewSettings(MEDIA_KEY_SCHEMA_ID)
        mediaKeyMap      = make(map[string]string)
)

func (op *MediaKeyManager) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                MEDIA_KEY_DEST,
                MEDIA_KEY_PATH,
                MEDIA_KEY_IFC,
        }
}

func initMediaKey() {
        defer func() {
                if err := recover(); err != nil {
                        fmt.Println("Recover Error:", err)
                }
        }()

        keyList := mediaKeySettings.ListKeys()
        for _, key := range keyList {
                value := mediaKeySettings.GetString(key)
                mediaKeyMap[key] = value
                grabKeyPress(X.RootWin(), convertKeyToMod(value))
        }
}

func (op *MediaKeyManager) listenMediaKey() {
        mediaKeySettings.Connect("changed", func(s *gio.Settings, key string) {
                value := mediaKeySettings.GetString(key)
                v := mediaKeyMap[key]
                if v != value {
                        ungrabKey(X.RootWin(), convertKeyToMod(v))
                        grabKeyPress(X.RootWin(), convertKeyToMod(value))
                        mediaKeyMap[key] = value
                }
        })
}

func (op *MediaKeyManager) emitSignal(modStr, keyStr string, press bool) bool {
        fmt.Printf("Emit mod: %s, key: %s\n", modStr, keyStr)
        switch keyStr {
        case "XF86MonBrightnessUp":
                op.BrightnessUp(press)
        case "XF86MonBrightnessDown":
                op.BrightnessDown(press)
        case "XF86AudioMute":
                op.AudioMute(press)
        case "XF86AudioLowerVolume":
                op.AudioDown(press)
        case "XF86AudioRaiseVolume":
                op.AudioUp(press)
        case "Num_Lock":
                if strings.Contains(modStr, "mod2") {
                        op.NumLockOff(press)
                } else {
                        op.NumLockOn(press)
                }
        case "Caps_Lock":
                if strings.Contains(modStr, "lock") {
                        op.CapsLockOff(press)
                } else {
                        op.CapsLockOn(press)
                }
        case "XF86TouchPadOn":
                op.TouchPadOn(press)
        case "XF86TouchPadOff":
                op.TouchPadOff(press)
        case "XF86Display":
                op.SwitchMonitors(press)
        case "XF86PowerOff":
                op.PowerOff(press)
        case "XF86Sleep":
                op.PowerSleep(press)
        case "p", "P":
                tmps := deleteSpecialMod(modStr)
                println("mod string after delete: ", tmps)
                if strings.Contains(tmps, "-") {
                        return false
                }
                if strings.Contains(tmps, "mod4") {
                        op.SwitchMonitors(press)
                }
        case "space":
                tmps := deleteSpecialMod(modStr)
                println("mod string after delete: ", tmps)
                if strings.Contains(tmps, "-") {
                        return false
                }
                if strings.Contains(modStr, "mod4") {
                        op.SwitchLayout(press)
                }
        case "XF86AudioPlay":
                op.AudioPlay(press)
        case "XF86AudioPause":
                op.AudioPause(press)
        case "XF86AudioStop":
                op.AudioStop(press)
        case "XF86AudioPrev":
                op.AudioPrevious(press)
        case "XF86AudioNext":
                op.AudioNext(press)
        case "XF86AudioRewind":
                op.AudioRewind(press)
        case "XF86AudioForward":
                op.AudioForward(press)
        case "XF86AudioRepeat":
                op.AudioRepeat(press)
        case "XF86WWW":
                op.LaunchBrowser(press)
        case "XF86Mail":
                op.LaunchEmail(press)
        case "XF86Calculator":
                op.LaunchCalculator(press)
        default:
                return false
        }

        return true
}

func startMediaKey() {
        m := &MediaKeyManager{}

        initMediaKey()
        m.listenMediaKey()
        m.listenKeyPressEvent()

        dbus.InstallOnSession(m)
}

/*
 * delete Num_Lock and Caps_Lock
 */
func deleteSpecialMod(modStr string) string {
        ret := ""
        strs := strings.Split(modStr, "-")
        l := len(strs)
        for i, s := range strs {
                if s == "lock" || s == "mod2" {
                        continue
                }

                if i == l-1 {
                        ret += s
                        break
                }
                ret += s + "-"
        }

        return ret
}
