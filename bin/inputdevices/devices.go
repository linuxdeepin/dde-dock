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
	"io/ioutil"
	"regexp"
	"strings"
)

func (op *MouseEntry) Reset() bool {
	list := mouseSettings.ListKeys()
	for _, key := range list {
		mouseSettings.Reset(key)
	}

	return true
}

func (OP *TPadEntry) Reset() bool {
	list := tpadSettings.ListKeys()
	for _, key := range list {
		tpadSettings.Reset(key)
	}

	return true
}

func (op *KbdEntry) Reset() bool {
	list := kbdSettings.ListKeys()
	for _, key := range list {
		kbdSettings.Reset(key)
	}

	return true
}

func (op *KbdEntry) LayoutList() map[string]string {
	defer func() {
		if err := recover(); err != nil {
			logObj.Warning("Receive Error: ", err)
		}
	}()

	if layoutDescMap == nil {
		datas := parseXML(_LAYOUT_XML_PATH)
		layoutDescMap = getLayoutList(datas)
	}

	return layoutDescMap
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

func (op *KbdEntry) AddLayoutOption(option string) {
	if len(option) < 1 {
		return
	}

	options := kbdSettings.GetStrv(KBD_KEY_LAYOUT_OPTIONS)
	if !utilObj.IsElementExist(option, options) {
		options = append(options, option)
		kbdSettings.SetStrv(KBD_KEY_LAYOUT_OPTIONS, options)
	}
}

func (op *KbdEntry) DeleteLayoutOption(option string) {
	if len(option) < 1 {
		return
	}

	options := kbdSettings.GetStrv(KBD_KEY_LAYOUT_OPTIONS)
	if utilObj.IsElementExist(option, options) {
		tmp := []string{}
		for _, v := range options {
			if v == option {
				continue
			}
			tmp = append(tmp, v)
		}
		kbdSettings.SetStrv(KBD_KEY_LAYOUT_OPTIONS, tmp)
	}
}

func (op *KbdEntry) ClearLayoutOption() {
	kbdSettings.Reset(KBD_KEY_LAYOUT_OPTIONS)
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
	m.MotionAcceleration = property.NewGSettingsFloatProperty(
		m, "MotionAcceleration",
		mouseSettings, MOUSE_KEY_ACCEL)
	m.MotionThreshold = property.NewGSettingsFloatProperty(
		m, "MotionThreshold",
		mouseSettings, MOUSE_KEY_THRES)
	m.DoubleClick = property.NewGSettingsIntProperty(
		m, "DoubleClick",
		mouseSettings, MOUSE_KEY_DOUBLE_CLICK)
	m.DragThreshold = property.NewGSettingsIntProperty(
		m, "DragThreshold",
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
	m.TapClick = property.NewGSettingsBoolProperty(
		m, "TapClick",
		tpadSettings, TPAD_KEY_TAP_CLICK)
	m.MotionAcceleration = property.NewGSettingsFloatProperty(
		m, "MotionAcceleration",
		tpadSettings, TPAD_KEY_ACCEL)
	m.MotionThreshold = property.NewGSettingsFloatProperty(
		m, "MotionThreshold",
		tpadSettings, TPAD_KEY_THRES)
	m.DoubleClick = property.NewGSettingsIntProperty(
		m, "DoubleClick",
		mouseSettings, MOUSE_KEY_DOUBLE_CLICK)
	m.DragThreshold = property.NewGSettingsIntProperty(
		m, "DragThreshold",
		mouseSettings, MOUSE_KEY_DRAG_THRES)
	m.deviceId = "TouchPad"

	return m
}

func NewKeyboard() *KbdEntry {
	m := &KbdEntry{}

	m.CurrentLayout = property.NewGSettingsStringProperty(
		m, "CurrentLayout",
		kbdSettings, KBD_KEY_LAYOUT)
	if len(m.CurrentLayout.GetValue().(string)) < 1 {
		v := getDefaultLayout()
		if _, ok := layoutDescMap[v]; ok {
			m.CurrentLayout.SetValue(v)
			kbdSettings.SetString(KBD_KEY_LAYOUT, v)
		} else {
			strs := strings.Split(v, LAYOUT_DELIM)
			m.CurrentLayout.SetValue(strs[0] + LAYOUT_DELIM)
			kbdSettings.SetString(KBD_KEY_LAYOUT, strs[0]+LAYOUT_DELIM)
			if len(strs[1]) > 0 {
				list := kbdSettings.GetStrv(KBD_KEY_LAYOUT_OPTIONS)
				list = append(list, strs[1])
				kbdSettings.SetStrv(KBD_KEY_LAYOUT_OPTIONS, list)
			}
		}
	}
	m.RepeatEnabled = property.NewGSettingsBoolProperty(
		m, "RepeatEnabled",
		kbdSettings, KBD_KEY_REPEAT_ENABLE)
	m.RepeatInterval = property.NewGSettingsUintProperty(
		m, "RepeatInterval",
		kbdSettings, KBD_KEY_REPEAT_INTERVAL)
	m.RepeatDelay = property.NewGSettingsUintProperty(
		m, "RepeatDelay",
		kbdSettings, KBD_KEY_DELAY)
	m.CursorBlink = property.NewGSettingsIntProperty(
		m, "CursorBlink",
		kbdSettings, KBD_CURSOR_BLINK_TIME)
	m.CapslockToggle = property.NewGSettingsBoolProperty(
		m, "CapslockToggle",
		kbdSettings, KBD_KEY_CAPSLOCK_TOGGLE)
	m.UserLayoutList = property.NewGSettingsStrvProperty(
		m, "UserLayoutList",
		kbdSettings, KBD_KEY_USER_LAYOUT_LIST)
	m.deviceId = "Keyboard"

	return m
}

func getDefaultLayout() string {
	layout := "us"
	option := ""
	contents, err := ioutil.ReadFile(KBD_DEFAULT_FILE)
	if err != nil {
		logObj.Warning("ReadFile Failed:", err)
		return layout + LAYOUT_DELIM + option
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if ok, _ := regexp.MatchString(`^XKBLAYOUT`, line); ok {
			layout = strings.Split(line, "=")[1]
		} else if ok, _ := regexp.MatchString(`^XKBOPTIONS`, line); ok {
			option = strings.Split(line, "=")[1]
		}
	}

	layout = strings.Trim(layout, "\"")
	option = strings.Trim(option, "\"")
	if len(layout) < 1 {
		layout = "us"
		option = ""
	}

	return layout + LAYOUT_DELIM + option
}
