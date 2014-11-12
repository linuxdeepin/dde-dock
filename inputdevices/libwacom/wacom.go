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

package libwacom

import (
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libwrapper"
	"pkg.linuxdeepin.com/lib/dbus/property"
	. "pkg.linuxdeepin.com/lib/gettext"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/log"
)

const (
	wacomKeyLeftHanded        = "left-handed"
	wacomKeyCursorMode        = "cursor-mode"
	wacomKeyUpAction          = "keyup-action"
	wacomKeyDownAction        = "keydown-action"
	wacomKeyDoubleDelta       = "double-delta"
	wacomKeyPressureSensitive = "pressure-sensitive"
)

type KeyActionInfo struct {
	Id   string
	Desc string
}

var descActionMap = map[string]string{
	"LeftClick":   "1",
	"MiddleClick": "2",
	"RightClick":  "3",
	"PageUp":      "key KP_Page_Up",
	"PageDown":    "key KP_Page_Down",
}

type Wacom struct {
	LeftHanded *property.GSettingsBoolProperty `access:"readwrite"`
	CursorMode *property.GSettingsBoolProperty `access:"readwrite"`

	KeyUpAction   *property.GSettingsStringProperty `access:"readwrite"`
	KeyDownAction *property.GSettingsStringProperty `access:"readwrite"`

	DoubleDelta       *property.GSettingsUintProperty `access:"readwrite"`
	PressureSensitive *property.GSettingsUintProperty `access:"readwrite"`

	DeviceList  []libwrapper.XIDeviceInfo
	ActionInfos []KeyActionInfo
	Exist       bool

	logger   *log.Logger
	settings *gio.Settings
}

func (wacom *Wacom) Reset() {
	keys := wacom.settings.ListKeys()
	for _, key := range keys {
		wacom.settings.Reset(key)
	}
}

var _wacom *Wacom

func NewWacom(l *log.Logger) *Wacom {
	wacom := &Wacom{}

	wacom.settings = gio.NewSettings("com.deepin.dde.wacom")

	wacom.LeftHanded = property.NewGSettingsBoolProperty(
		wacom, "LeftHanded",
		wacom.settings, wacomKeyLeftHanded)
	wacom.CursorMode = property.NewGSettingsBoolProperty(
		wacom, "CursorMode",
		wacom.settings, wacomKeyCursorMode)

	wacom.KeyUpAction = property.NewGSettingsStringProperty(
		wacom, "KeyUpAction",
		wacom.settings, wacomKeyUpAction)
	wacom.KeyDownAction = property.NewGSettingsStringProperty(
		wacom, "KeyDownAction",
		wacom.settings, wacomKeyDownAction)

	wacom.DoubleDelta = property.NewGSettingsUintProperty(
		wacom, "DoubleDelta",
		wacom.settings, wacomKeyDoubleDelta)
	wacom.PressureSensitive = property.NewGSettingsUintProperty(
		wacom, "PressureSensitive",
		wacom.settings, wacomKeyPressureSensitive)

	_, _, wacomList := libwrapper.GetDevicesList()
	wacom.setPropDeviceList(wacomList)
	if len(wacom.DeviceList) > 0 {
		wacom.setPropExist(true)
	} else {
		wacom.setPropExist(false)
	}

	wacom.logger = l
	wacom.ActionInfos = generateActionInfos()

	_wacom = wacom
	wacom.init()
	wacom.handleGSettings()

	return wacom
}

func HandleDeviceChanged(devList []libwrapper.XIDeviceInfo) {
	if _wacom == nil {
		return
	}

	_wacom.setPropDeviceList(devList)
	if len(devList) == 0 {
		_wacom.setPropExist(false)
	} else {
		_wacom.setPropExist(true)
		_wacom.init()
	}
}

/**
 * TODO:
 *	HandleDeviceAdded
 *	HandleDeviceRemoved
 **/

/**
 * KeyAction: PageUp/PageDown/LeftClick/RightClick/MiddleClick
 */
func generateActionInfos() []KeyActionInfo {
	return []KeyActionInfo{
		{
			Id:   "LeftClick",
			Desc: Tr("Left Click"),
		},
		{
			Id:   "MiddleClick",
			Desc: Tr("Middle Click"),
		},
		{
			Id:   "RightClick",
			Desc: Tr("Right Click"),
		},
		{
			Id:   "PageUp",
			Desc: Tr("Page Up"),
		},
		{
			Id:   "PageDown",
			Desc: Tr("Page Down"),
		},
	}
}

func (wacom *Wacom) init() {
	if !wacom.Exist {
		return
	}

	wacom.rotationAngle(wacom.LeftHanded.Get())
	wacom.cursorMode(wacom.CursorMode.Get())

	wacom.keyUpAction(wacom.KeyUpAction.Get())
	wacom.keyDownAction(wacom.KeyDownAction.Get())

	wacom.pressureSensitive(wacom.PressureSensitive.Get())
	wacom.doubleDelta(wacom.DoubleDelta.Get())
}
