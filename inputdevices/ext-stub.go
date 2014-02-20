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
        "strings"
)

func (dev *ExtDevManager) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{_EXT_DEV_NAME, _EXT_DEV_PATH, _EXT_DEV_IFC}
}

func (keyboard *KeyboardEntry) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                _EXT_DEV_NAME,
                _EXT_ENTRY_PATH + keyboard.DeviceID,
                _EXT_ENTRY_IFC + keyboard.DeviceID,
        }
}

func (mouse *MouseEntry) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                _EXT_DEV_NAME,
                _EXT_ENTRY_PATH + mouse.DeviceID,
                _EXT_ENTRY_IFC + mouse.DeviceID,
        }
}

func (tpad *TPadEntry) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                _EXT_DEV_NAME,
                _EXT_ENTRY_PATH + tpad.DeviceID,
                _EXT_ENTRY_IFC + tpad.DeviceID,
        }
}

func (keyboard *KeyboardEntry) listenLayoutChanged() {
        _layoutGSettings.Connect("changed", func(s *gio.Settings, key string) {
                keyboard.getPropName("CurrentLayout")
        })
}

func (keyboard *KeyboardEntry) OnPropertiesChanged(propName string, old interface{}) {
        switch propName {
        case "CurrentLayout":
                if v, ok := old.(string); ok && v != keyboard.CurrentLayout {
                        keyboard.setPropName(propName)
                }
        }
}

func (keyboard *KeyboardEntry) setPropName(propName string) {
        switch propName {
        case "CurrentLayout":
                strs := strings.Split(keyboard.CurrentLayout, LAYOUT_DELIM)
                switch len(strs) {
                case 0:
                        _layoutGSettings.SetStrv("layouts", []string{strs[0]})
                        _layoutGSettings.SetStrv("options", []string{})
                case 1:
                        _layoutGSettings.SetStrv("layouts", []string{strs[0]})
                        _layoutGSettings.SetStrv("options", []string{})
                case 2:
                        _layoutGSettings.SetStrv("layouts", []string{strs[0]})
                        _layoutGSettings.SetStrv("options", []string{strs[1]})
                }
                dbus.NotifyChange(keyboard, propName)
        case "UserLayoutList":
                _keyRepeatGSettings.SetStrv("user-layout-list", keyboard.UserLayoutList)
                dbus.NotifyChange(keyboard, propName)
        }
}

func (keyboard *KeyboardEntry) getPropName(propName string) {
        switch propName {
        case "CurrentLayout":
                layout := _layoutGSettings.GetStrv("layouts")
                option := _layoutGSettings.GetStrv("options")
                if len(layout) >= 1 {
                        keyboard.CurrentLayout = layout[0] + LAYOUT_DELIM
                        if len(option) >= 1 {
                                keyboard.CurrentLayout += option[0]
                        }
                } else {
                        keyboard.CurrentLayout = LAYOUT_DELIM
                }
                dbus.NotifyChange(keyboard, propName)
        case "UserLayoutList":
                keyboard.UserLayoutList = _keyRepeatGSettings.GetStrv("user-layout-list")
                dbus.NotifyChange(keyboard, propName)
        }
}

func (keyboard *KeyboardEntry) appendUserLayout(str string) {
        if len(str) <= 0 {
                str = "us;"
        } else if !strings.Contains(str, LAYOUT_DELIM) {
                return
        }
        if stringIsExist(str, keyboard.UserLayoutList) {
                return
        }

        keyboard.UserLayoutList = append(keyboard.UserLayoutList, str)
        keyboard.setPropName("UserLayoutList")
}

func (keyboard *KeyboardEntry) deleteUserLayout(str string) {
        if !strings.Contains(str, LAYOUT_DELIM) {
                return
        }
        if !stringIsExist(str, keyboard.UserLayoutList) {
                return
        }

        tmps := []string{}
        for _, v := range keyboard.UserLayoutList {
                if v == str {
                        continue
                } else {
                        tmps = append(tmps, v)
                }
        }

        keyboard.UserLayoutList = tmps
        keyboard.setPropName("UserLayoutList")
}

func stringIsExist(str string, strs []string) bool {
        for _, v := range strs {
                if str == v {
                        return true
                }
        }

        return false
}

func stringArrayIsEqual(array1, array2 []string) bool {
        l1 := len(array1)
        l2 := len(array2)

        if l1 != l2 {
                return false
        }

        for i := 0; i < l1; i++ {
                if array1[i] != array2[i] {
                        return false
                }
        }
        return true
}
