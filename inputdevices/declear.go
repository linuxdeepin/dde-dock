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
)

const (
	TPAD_KEY_ENABLE         = "touchpad-enabled"
	TPAD_KEY_LEFT_HAND      = "left-handed"
	TPAD_KEY_W_TYPING       = "disable-while-typing"
	TPAD_KEY_NATURAL_SCROLL = "natural-scroll"
	TPAD_KEY_EDGE_SCROLL    = "edge-scroll-enabled"
	TPAD_KEY_HORIZ_SCROLL   = "horiz-scroll-enabled"
	TPAD_KEY_VERT_SCROLL    = "vert-scroll-enabled"
	TPAD_KEY_ACCEL          = "motion-acceleration"
	TPAD_KEY_THRES          = "motion-threshold"

	MOUSE_KEY_LEFT_HAND    = "left-handed"
	MOUSE_KEY_MID_BUTTON   = "middle-button-enabled"
	MOUSE_KEY_ACCEL        = "motion-acceleration"
	MOUSE_KEY_THRES        = "motion-threshold"
	MOUSE_KEY_DOUBLE_CLICK = "double-click"
	MOUSE_KEY_DRAG_THRES   = "drag-threshold"

	KBD_KEY_REPEAT_ENABLE    = "repeat-enabled"
	KBD_KEY_REPEAT_INTERVAL  = "repeat-interval"
	KBD_KEY_DELAY            = "delay"
	KBD_KEY_LAYOUT           = "layout"
	KBD_KEY_LAYOUT_MODEL     = "layout-model"
	KBD_KEY_LAYOUT_OPTION    = "layout-option"
	KBD_KEY_USER_LAYOUT_LIST = "user-layout-list"
	KBD_CURSOR_BLINK_TIME    = "cursor-blink-time"
	KBD_DEFAULT_FILE         = "/etc/default/keyboard"

	DEVICE_DEST  = "com.deepin.daemon.InputDevices"
	DEVICE_PATH  = "/com/deepin/daemon/InputDevice/"
	DEVICE_IFC   = "com.deepin.daemon.InputDevice."
	MANAGER_PATH = "/com/deepin/daemon/InputDevices"
	MANAGER_IFC  = "com.deepin.daemon.InputDevices"
)

type MouseEntry struct {
	LeftHanded         *property.GSettingsBoolProperty  `access:"readwrite"`
	MotionAcceleration *property.GSettingsFloatProperty `access:"readwrite"`
	MotionThreshold    *property.GSettingsFloatProperty `access:"readwrite"`
	DoubleClick        *property.GSettingsIntProperty   `access:"readwrite"`
	DragThreshold      *property.GSettingsIntProperty   `access:"readwrite"`
	deviceId           string
}

type TPadEntry struct {
	TPadEnable         *property.GSettingsBoolProperty  `access:"readwrite"`
	LeftHanded         *property.GSettingsBoolProperty  `access:"readwrite"`
	DisableIfTyping    *property.GSettingsBoolProperty  `access:"readwrite"`
	NaturalScroll      *property.GSettingsBoolProperty  `access:"readwrite"`
	EdgeScroll         *property.GSettingsBoolProperty  `access:"readwrite"`
	HorizScroll        *property.GSettingsBoolProperty  `access:"readwrite"`
	VertScroll         *property.GSettingsBoolProperty  `access:"readwrite"`
	MotionAcceleration *property.GSettingsFloatProperty `access:"readwrite"`
	MotionThreshold    *property.GSettingsFloatProperty `access:"readwrite"`
	DoubleClick        *property.GSettingsIntProperty   `access:"readwrite"`
	DragThreshold      *property.GSettingsIntProperty   `access:"readwrite"`
	deviceId           string
}

type KbdEntry struct {
	RepeatEnabled  *property.GSettingsBoolProperty   `access:"readwrite"`
	RepeatInterval *property.GSettingsUintProperty   `access:"readwrite"`
	RepeatDelay    *property.GSettingsUintProperty   `access:"readwrite"`
	CurrentLayout  *property.GSettingsStringProperty `access:"readwrite"`
	CursorBlink    *property.GSettingsIntProperty    `access:"readwrite"`
	UserLayoutList *property.GSettingsStrvProperty
	deviceId       string
}
