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

package wacom

import (
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/wrapper"
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

	DeviceList  []wrapper.XIDeviceInfo
	ActionInfos []KeyActionInfo
	Exist       bool

	logger   *log.Logger
	settings *gio.Settings
}

func (w *Wacom) Reset() {
	keys := w.settings.ListKeys()
	for _, key := range keys {
		w.settings.Reset(key)
	}
}

var _wacom *Wacom

func NewWacom(l *log.Logger) *Wacom {
	w := &Wacom{}

	w.settings = gio.NewSettings("com.deepin.dde.wacom")

	w.LeftHanded = property.NewGSettingsBoolProperty(
		w, "LeftHanded",
		w.settings, wacomKeyLeftHanded)
	w.CursorMode = property.NewGSettingsBoolProperty(
		w, "CursorMode",
		w.settings, wacomKeyCursorMode)

	w.KeyUpAction = property.NewGSettingsStringProperty(
		w, "KeyUpAction",
		w.settings, wacomKeyUpAction)
	w.KeyDownAction = property.NewGSettingsStringProperty(
		w, "KeyDownAction",
		w.settings, wacomKeyDownAction)

	w.DoubleDelta = property.NewGSettingsUintProperty(
		w, "DoubleDelta",
		w.settings, wacomKeyDoubleDelta)
	w.PressureSensitive = property.NewGSettingsUintProperty(
		w, "PressureSensitive",
		w.settings, wacomKeyPressureSensitive)

	_, _, wacomList := wrapper.GetDevicesList()
	w.setPropDeviceList(wacomList)
	if len(w.DeviceList) > 0 {
		w.setPropExist(true)
	} else {
		w.setPropExist(false)
	}

	w.logger = l
	w.ActionInfos = generateActionInfos()

	_wacom = w
	w.init()
	w.handleGSettings()

	return w
}

func HandleDeviceChanged(devList []wrapper.XIDeviceInfo) {
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

func (w *Wacom) init() {
	if !w.Exist {
		return
	}

	w.rotationAngle(w.LeftHanded.Get())
	w.cursorMode(w.CursorMode.Get())

	w.keyUpAction(w.KeyUpAction.Get())
	w.keyDownAction(w.KeyDownAction.Get())

	w.pressureSensitive(w.PressureSensitive.Get())
	w.doubleDelta(w.DoubleDelta.Get())
}
