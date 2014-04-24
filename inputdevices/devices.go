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
        "dlib/dbus/property"
        "strings"
)

func (op *KbdEntry) LayoutList() map[string]string {
        defer func() {
                if err := recover(); err != nil {
                        logObj.Warning("Receive Error: ", err)
                }
        }()

        datas := parseXML(_LAYOUT_XML_PATH)
        layouts := getLayoutList(datas)

        return layouts
}

func (op *KbdEntry) GetLayoutLocale(layout string) string {
        if len(layout) < 1 || !strings.Contains(layout, LAYOUT_DELIM) {
                layout = "us" + LAYOUT_DELIM
        }
        listMap := op.LayoutList()
        desc, ok := listMap[layout]
        if !ok {
                logObj.Warningf("Invalid layout: '%s'", layout)
                return ""
        }

        return desc
}

func (op *KbdEntry) AddUserLayout(layout string) bool {
        if len(layout) < 1 || !strings.Contains(layout, LAYOUT_DELIM) {
                layout = "us" + LAYOUT_DELIM
        }

        list := op.UserLayoutList.GetValue().([]string)
        if utilObj.IsElementExist(layout, list) {
                return false
        }

        list = append(list, layout)
        kbdSettings.SetStrv(KBD_KEY_USER_LAYOUT_LIST, list)
        return true
}

func (op *KbdEntry) DeleteUserLayout(layout string) bool {
        if len(layout) < 1 || !strings.Contains(layout, LAYOUT_DELIM) {
                return false
        }

        list := op.UserLayoutList.GetValue().([]string)
        if !utilObj.IsElementExist(layout, list) {
                return false
        }

        tmp := []string{}
        for _, l := range list {
                if l == layout {
                        continue
                }
                tmp = append(tmp, l)
        }

        kbdSettings.SetStrv(KBD_KEY_USER_LAYOUT_LIST, tmp)
        return true
}

func NewMouse() *MouseEntry {
        m := &MouseEntry{}

        m.LeftHanded = property.NewGSettingsBoolProperty(
                m, "LeftHanded",
                mouseSettings, MOUSE_KEY_LEFT_HAND)
        m.MotionAccel = property.NewGSettingsFloatProperty(
                m, "MotionAccel",
                mouseSettings, MOUSE_KEY_ACCEL)
        m.MotionThres = property.NewGSettingsFloatProperty(
                m, "MotionThres",
                mouseSettings, MOUSE_KEY_THRES)
        m.DoubleClick = property.NewGSettingsIntProperty(
                m, "DoubleClick",
                mouseSettings, MOUSE_KEY_DOUBLE_CLICK)
        m.DragThres = property.NewGSettingsIntProperty(
                m, "DragThres",
                mouseSettings, MOUSE_KEY_DRAG_THRES)
        m.deviceId = "Mouse"

        return m
}

func NewTPad() *TPadEntry {
        m := &TPadEntry{}

        m.TPadEnable = property.NewGSettingsBoolProperty(
                m, "TPadEnable",
                tpadSettings, TPAD_KEY_ENABLE)
        m.LeftHanded = property.NewGSettingsBoolProperty(
                m, "LeftHanded",
                tpadSettings, TPAD_KEY_LEFT_HAND)
        m.DisableIfTyping = property.NewGSettingsBoolProperty(
                m, "DisableIfTyping",
                tpadSettings, TPAD_KEY_W_TYPING)
        m.NaturalScroll = property.NewGSettingsBoolProperty(
                m, "NaturalScroll",
                tpadSettings, TPAD_KEY_NATURAL_SCROLL)
        m.EdgeScroll = property.NewGSettingsBoolProperty(
                m, "EdgeScroll",
                tpadSettings, TPAD_KEY_EDGE_SCROLL)
        m.HorizScroll = property.NewGSettingsBoolProperty(
                m, "HorizScroll",
                tpadSettings, TPAD_KEY_HORIZ_SCROLL)
        m.VertScroll = property.NewGSettingsBoolProperty(
                m, "VertScroll",
                tpadSettings, TPAD_KEY_VERT_SCROLL)
        m.MotionAccel = property.NewGSettingsFloatProperty(
                m, "MotionAccel",
                tpadSettings, TPAD_KEY_ACCEL)
        m.MotionThres = property.NewGSettingsFloatProperty(
                m, "MotionThres",
                tpadSettings, TPAD_KEY_THRES)
        m.DoubleClick = property.NewGSettingsIntProperty(
                m, "DoubleClick",
                mouseSettings, MOUSE_KEY_DOUBLE_CLICK)
        m.DragThres = property.NewGSettingsIntProperty(
                m, "DragThres",
                mouseSettings, MOUSE_KEY_DRAG_THRES)
        m.deviceId = "TouchPad"

        return m
}

func NewKeyboard() *KbdEntry {
        m := &KbdEntry{}

        m.CurrentLayout = property.NewGSettingsStringProperty(
                m, "CurrentLayout",
                kbdSettings, KBD_KEY_LAYOUT)
        logObj.Debug("CurrentLayout: ", m.CurrentLayout.GetValue().(string))
        m.RepeatEnabled = property.NewGSettingsBoolProperty(
                m, "RepeatEnabled",
                kbdSettings, KBD_KEY_REPEAT_ENABLE)
        logObj.Debug("RepeatEnabled: ", m.RepeatEnabled.GetValue().(bool))
        m.RepeatInterval = property.NewGSettingsUintProperty(
                m, "RepeatInterval",
                kbdSettings, KBD_KEY_REPEAT_INTERVAL)
        logObj.Debug("RepeatInterval: ", m.RepeatInterval.GetValue().(uint32))
        m.RepeatDelay = property.NewGSettingsUintProperty(
                m, "RepeatDelay",
                kbdSettings, KBD_KEY_DELAY)
        logObj.Debug("RepeatDelay: ", m.RepeatDelay.GetValue().(uint32))
        m.CursorBlink = property.NewGSettingsIntProperty(
                m, "CursorBlink",
                kbdSettings, KBD_CURSOR_BLINK_TIME)
        logObj.Debug("CursorBlink: ", m.CursorBlink.GetValue().(int32))
        m.UserLayoutList = property.NewGSettingsStrvProperty(
                m, "UserLayoutList",
                kbdSettings, KBD_KEY_USER_LAYOUT_LIST)
        logObj.Debug("UserLayoutList: ", m.UserLayoutList.GetValue().([]string))
        m.deviceId = "Keyboard"
        logObj.Debug("deviceId: ", m.deviceId)

        return m
}
