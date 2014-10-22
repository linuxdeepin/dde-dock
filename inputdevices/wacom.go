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

package inputdevices

import (
	"fmt"
	"os/exec"
	"pkg.linuxdeepin.com/lib/dbus/property"
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

const (
	WACOM_KEY_LEFT_HANDED        = "left-handed"
	WACOM_KEY_CURSOR_MODE        = "cursor-mode"
	WACOM_KEY_UP_ACTION          = "keyup-action"
	WACOM_KEY_DOWN_ACTION        = "keydown-action"
	WACOM_KEY_DOUBLE_DELTA       = "double-delta"
	WACOM_KEY_PRESSURE_SENSITIVE = "pressure-sensitive"
)

type WacomManager struct {
	LeftHanded *property.GSettingsBoolProperty `access:"readwrite"`
	CursorMode *property.GSettingsBoolProperty `access:"readwrite"`

	KeyUpAction   *property.GSettingsStringProperty `access:"readwrite"`
	KeyDownAction *property.GSettingsStringProperty `access:"readwrite"`

	DoubleDelta       *property.GSettingsUintProperty `access:"readwrite"`
	PressureSensitive *property.GSettingsUintProperty `access:"readwrite"`

	DeviceList []PointerDeviceInfo
	Exist      bool

	settings       *gio.Settings
	listenFlag     bool
	_descActionMap map[string]string
}

func (wacom *WacomManager) Reset() {
	keys := wacom.settings.ListKeys()
	for _, key := range keys {
		wacom.settings.Reset(key)
	}
}

/**
 * xsetwacom set deviceName Rotate half/none
 * default: none
 */
func (wManager *WacomManager) setRotate(leftHanded bool) {
	for _, info := range wManager.DeviceList {
		cmdline := "xsetwacom set \"" + info.Name + "\" Rotate "
		if leftHanded {
			cmdline += "half"
		} else {
			cmdline += "none"
		}

		wManager.doCommand(cmdline)
	}
}

/**
 * xsetwacom set deviceName mode Relative/Absolute
 * default: Absolute for  stylus,  eraser  and  tablet  PC  touch;
 *          Relative for cursor and tablet touch.
 */
func (wManager *WacomManager) setCursorMode(cursorMode bool) {
	for _, info := range wManager.DeviceList {
		cmdline := "xsetwacom set \"" + info.Name + "\" mode "
		if cursorMode {
			cmdline += "Relative"
		} else {
			cmdline += "Absolute"
		}

		wManager.doCommand(cmdline)
	}
}

/**
 * xsetwacom set deviceName Button 3 3/"KP_Page_Up"
 * default: 3
 */
func (wManager *WacomManager) setKeyUp(action string) {
	value, ok := wManager._descActionMap[action]
	if !ok {
		return
	}

	for _, info := range wManager.DeviceList {
		cmdline := "xsetwacom set \"" + info.Name + "\" Button 3 " + value
		wManager.doCommand(cmdline)
	}
}

/**
 * xsetwacom set deviceName Button 2 2/"KP_Page_Down"
 * default: 2
 */
func (wManager *WacomManager) setKeyDown(action string) {
	value, ok := wManager._descActionMap[action]
	if !ok {
		return
	}

	for _, info := range wManager.DeviceList {
		cmdline := "xsetwacom set \"" + info.Name + "\" Button 2 " + value
		wManager.doCommand(cmdline)
	}
}

/**
 * xsetwacom set deviceName Suppress 0-100
 * default: 2
 */
func (wManager *WacomManager) setDoubleDelta(delta uint32) {
	for _, info := range wManager.DeviceList {
		cmdline := fmt.Sprintf("xsetwacom set \"%s\" Suppress %v",
			info.Name, delta)
		wManager.doCommand(cmdline)
	}
}

/**
 * xsetwacom set deviceName Threshold 0-2047
 * default: 27
 */
func (wManager *WacomManager) setPressureSensitive(pressure uint32) {
	if pressure > 2047 {
		return
	}

	for _, info := range wManager.DeviceList {
		cmdline := fmt.Sprintf("xsetwacom set \"%s\" Threshold %v",
			info.Name, pressure)
		wManager.doCommand(cmdline)
	}
}

func (wManager *WacomManager) doCommand(cmdline string) (err error) {
	if len(cmdline) < 1 {
		logger.Debug("doCommand args is nil")
		return fmt.Errorf("doCommand args is nil")
	}

	logger.Debug(cmdline)
	err = exec.Command("/bin/sh", "-c", cmdline).Run()
	logger.Debug(err)

	return
}

var _wManager *WacomManager

func GetWacomManager() *WacomManager {
	if _wManager == nil {
		_wManager = newWacomManager()
	}

	return _wManager
}

func newWacomManager() *WacomManager {
	wManager := &WacomManager{}

	wManager.settings = gio.NewSettings("com.deepin.dde.wacom")

	wManager.LeftHanded = property.NewGSettingsBoolProperty(
		wManager, "LeftHanded",
		wManager.settings, WACOM_KEY_LEFT_HANDED)
	wManager.CursorMode = property.NewGSettingsBoolProperty(
		wManager, "CursorMode",
		wManager.settings, WACOM_KEY_CURSOR_MODE)

	wManager.KeyUpAction = property.NewGSettingsStringProperty(
		wManager, "KeyUpAction",
		wManager.settings, WACOM_KEY_UP_ACTION)
	wManager.KeyDownAction = property.NewGSettingsStringProperty(
		wManager, "KeyDownAction",
		wManager.settings, WACOM_KEY_DOWN_ACTION)

	wManager.DoubleDelta = property.NewGSettingsUintProperty(
		wManager, "DoubleDelta",
		wManager.settings, WACOM_KEY_DOUBLE_DELTA)
	wManager.PressureSensitive = property.NewGSettingsUintProperty(
		wManager, "PressureSensitive",
		wManager.settings, WACOM_KEY_PRESSURE_SENSITIVE)

	_, _, wacomList := getPointerDeviceList()
	wManager.setPropDeviceList(wacomList)
	if len(wManager.DeviceList) > 0 {
		wManager.setPropExist(true)
	} else {
		wManager.setPropExist(false)
	}

	wManager._descActionMap = make(map[string]string)
	wManager.initDescActionMap()

	wManager.listenFlag = false
	wManager.init()

	return wManager
}

/**
 * KeyAction: PageUp/PageDown/LeftClick/RightClick/MiddleClick
 */
func (wManager *WacomManager) initDescActionMap() {
	if wManager._descActionMap == nil {
		wManager._descActionMap = make(map[string]string)
	}

	wManager._descActionMap = map[string]string{
		"LeftClick":   "1",
		"MiddleClick": "2",
		"RightClick":  "3",
		"PageUp":      "key KP_Page_Up",
		"PageDown":    "key KP_Page_Down",
	}
}

func (wManager *WacomManager) init() {
	if !wManager.Exist {
		return
	}

	wManager.setRotate(wManager.LeftHanded.Get())
	wManager.setCursorMode(wManager.CursorMode.Get())

	wManager.setKeyUp(wManager.KeyUpAction.Get())
	wManager.setKeyDown(wManager.KeyDownAction.Get())

	wManager.setPressureSensitive(wManager.PressureSensitive.Get())
	wManager.setDoubleDelta(wManager.DoubleDelta.Get())

	wManager.listenGSettings()
}
