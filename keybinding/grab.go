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

package keybinding

// #cgo pkg-config: x11 xtst glib-2.0
// #include "record.h"
// #include <stdlib.h>
import "C"

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"pkg.linuxdeepin.com/lib/dbus"
	"strings"
	"unsafe"
)

const (
	ACCEL_DELIM = "-"
)

func convertKeyToMod(key string) string {
	if v, ok := keyToModMap[key]; ok {
		return v
	}

	return key
}

func convertModToKey(mod string) string {
	if v, ok := modToKeyMap[mod]; ok {
		return v
	}

	return mod
}

func convertKeysToMods(keys string) string {
	array := strings.Split(keys, ACCEL_DELIM)
	ret := ""
	for i, v := range array {
		if i != 0 {
			ret += ACCEL_DELIM
		}
		tmp := convertKeyToMod(v)
		ret += tmp
	}

	return ret
}

func convertModsToKeys(mods string) string {
	array := strings.Split(mods, ACCEL_DELIM)
	ret := ""
	for i, v := range array {
		if i != 0 {
			ret += ACCEL_DELIM
		}
		tmp := convertModToKey(v)
		ret += tmp
	}

	return ret
}

/**
 * Input: <control><alt>t
 * Output: modx-modx-t
 */
func formatXGBShortcut(shortcut string) string {
	if len(shortcut) < 1 {
		return ""
	}

	ret := formatShortcut(shortcut)
	return convertKeysToMods(ret)
}

/**
 * Input: <control><alt>t
 * Output: control-alt-t
 */
func formatShortcut(shortcut string) string {
	l := len(shortcut)
	if l < 1 {
		logger.Debug("formatShortcut args error")
		return ""
	}

	str := strings.ToLower(shortcut)
	ret := ""
	flag := false
	start := 0
	end := 0

	for i, ch := range str {
		if ch == '<' {
			flag = true
			start = i
			continue
		}

		if ch == '>' && flag {
			end = i
			flag = false

			for j := start + 1; j < end; j++ {
				ret += string(str[j])
			}
			ret += ACCEL_DELIM
			continue
		}

		if !flag {
			ret += string(ch)
		}
	}

	// parse 'primary' to 'control'
	array := strings.Split(ret, ACCEL_DELIM)
	ret = ""
	for i, v := range array {
		if v == "primary" || v == "control" {
			// multi control
			if !strings.Contains(ret, "control") {
				if i != 0 {
					ret += ACCEL_DELIM
				}
				ret += "control"
			}
			continue
		}

		if i != 0 {
			ret += ACCEL_DELIM
		}
		ret += v
	}

	return ret
}

/**
 * delete Num_Lock and Caps_Lock
 */
func deleteSpecialMod(modStr string) string {
	if !strings.Contains(modStr, ACCEL_DELIM) {
		if modStr == "lock" || modStr == "mod2" {
			return ""
		}

		return modStr
	}

	ret := ""
	strs := strings.Split(modStr, ACCEL_DELIM)
	l := len(strs)
	for i, s := range strs {
		if s == "lock" || s == "mod2" {
			continue
		}

		if i == l-1 {
			ret += s
			break
		}
		ret += s + ACCEL_DELIM
	}

	if length := len(ret); length > 1 {
		if ret[length-1] == '-' {
			ret = ret[0 : length-1]
		}
	}

	return ret
}

func getSystemKeyPairs() map[string]string {
	systemPairs := make(map[string]string)
	for _, info := range systemIdDescList {
		if info.Id >= 0 && info.Id < 300 {
			if isInvalidConflict(info.Id) {
				continue
			}
			shortcut := getSystemKeyValue(info.Name, false)
			action := getSystemKeyValue(info.Name, true)
			systemPairs[shortcut] = action
		}
	}
	PrevSystemPairs = systemPairs
	return systemPairs
}

func getCustomKeyPairs() map[string]string {
	customPairs := make(map[string]string)
	customList := getCustomIdList()

	for _, i := range customList {
		if isInvalidConflict(i) {
			logger.Warningf("%d is invalid conflict", i)
			continue
		}

		gs := getSettingsById(i)
		if gs == nil {
			logger.Warningf("Get Settings For '%d' Failed", i)
			continue
		}
		shortcut := gs.GetString(CUSTOM_KEY_SHORTCUT)
		action := gs.GetString(CUSTOM_KEY_ACTION)
		customPairs[shortcut] = action
	}

	PrevCustomPairs = customPairs
	return customPairs
}

func grabKeyPress(wid xproto.Window, shortcut string) bool {
	if len(shortcut) < 1 {
		logger.Warning("grabKeyPress args error...")
		return false
	}

	shortcut = convertKeysym2Weird(shortcut)
	mod, keys, err := keybind.ParseString(X, shortcut)
	if err != nil {
		logger.Warning("In GrabKey Parse shortcut failed:", err)
		return false
	}

	for _, k := range keys {
		if err = keybind.GrabChecked(X, wid, mod, k); err != nil {
			logger.Warningf("Grab '%s' failed: %v", shortcut, err)
			ungrabKey(wid, shortcut)
			return false
		}
	}

	return true
}

func ungrabKey(wid xproto.Window, shortcut string) bool {
	if len(shortcut) < 1 {
		logger.Warning("Ungrab args failed...")
		return false
	}

	shortcut = convertKeysym2Weird(shortcut)
	mod, keys, err := keybind.ParseString(X, shortcut)
	if err != nil {
		logger.Warning("In UngrabKey Parse shortcut failed:", err)
		return false
	}

	for _, k := range keys {
		keybind.Ungrab(X, wid, mod, k)
	}

	return true
}

func grabKeyPairs(pairs map[string]string, isGrab bool) {
	for k, v := range pairs {
		if len(k) < 1 {
			continue
		}

		if strings.ToLower(k) == "super" {
			grabSignalShortcut("Super_L", v, isGrab)
			grabSignalShortcut("Super_R", v, isGrab)
			continue
		}

		shortcut := formatXGBShortcut(formatShortcut(k))
		keyInfo, ok := newKeycodeInfo(shortcut)
		if !ok {
			logger.Warningf("New Keycode Info Failed. Key: %s, Value: %s", k, v)
			continue
		}

		if isGrab {
			if grabKeyPress(X.RootWin(), shortcut) {
				grabKeyBindsMap[keyInfo] = v
			}
		} else {
			if ungrabKey(X.RootWin(), shortcut) {
				delete(grabKeyBindsMap, keyInfo)
			}
		}
	}
}

func grabMediaKeys() {
	keyList := mediaGSettings.ListKeys()
	for _, key := range keyList {
		value := mediaGSettings.GetString(key)
		grabKeyPress(X.RootWin(), convertKeysToMods(value))
	}
}

var pressKeyStr string

func grabKeyboardAndMouse() {
	//go func() {
	X, err := xgbutil.NewConn()
	if err != nil {
		logger.Info("Get New Connection Failed:", err)
		return
	}
	keybind.Initialize(X)
	mousebind.Initialize(X)

	err = keybind.GrabKeyboard(X, X.RootWin())
	if err != nil {
		logger.Info("Grab Keyboard Failed:", err)
		return
	}

	grabAllMouseButton(X)

	xevent.ButtonPressFun(
		func(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {
			dbus.Emit(GetManager(), "KeyReleaseEvent", "")
			ungrabAllMouseButton(X)
			keybind.UngrabKeyboard(X)
			logger.Info("Button Press Event")
			xevent.Quit(X)
		}).Connect(X, X.RootWin())

	xevent.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			value := parseKeyEnvent(X, e.State, e.Detail)
			pressKeyStr = value
			dbus.Emit(GetManager(), "KeyPressEvent", value)
		}).Connect(X, X.RootWin())

	xevent.KeyReleaseFun(
		func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			if strings.ToLower(pressKeyStr) == "super_l" ||
				strings.ToLower(pressKeyStr) == "super_r" {
				pressKeyStr = "Super"
			}

			dbus.Emit(GetManager(), "KeyReleaseEvent", pressKeyStr)
			pressKeyStr = ""
			ungrabAllMouseButton(X)
			keybind.UngrabKeyboard(X)
			logger.Infof("Key: %s\n", pressKeyStr)
			xevent.Quit(X)
		}).Connect(X, X.RootWin())

	xevent.Main(X)
	//}()
}

func grabAllMouseButton(X *xgbutil.XUtil) {
	mousebind.Grab(X, X.RootWin(), 0, 1, false)
	mousebind.Grab(X, X.RootWin(), 0, 2, false)
	mousebind.Grab(X, X.RootWin(), 0, 3, false)
}

func ungrabAllMouseButton(X *xgbutil.XUtil) {
	mousebind.Ungrab(X, X.RootWin(), 0, 1)
	mousebind.Ungrab(X, X.RootWin(), 0, 2)
	mousebind.Ungrab(X, X.RootWin(), 0, 3)
}

func grabSignalShortcut(shortcut, action string, isGrab bool) {
	if len(shortcut) < 1 {
		logger.Error("grabSignalKey args error")
		return
	}

	shortcut = convertKeysym2Weird(shortcut)
	mod, keys, err := keybind.ParseString(X, shortcut)
	if err != nil {
		logger.Errorf("ParseString error: %v", err)
		return
	}

	if mod > 0 || len(keys) < 1 {
		return
	}

	if isGrab {
		if len(action) < 1 {
			return
		}
		tmp := C.CString(action)
		defer C.free(unsafe.Pointer(tmp))
		C.grab_xrecord_key(C.int(keys[0]), tmp)
	} else {
		C.ungrab_xrecord_key(C.int(keys[0]))
	}
}

func ungrabSignalShortcut(shortcut string) {
	if len(shortcut) < 1 {
		return
	}

	shortcut = convertKeysym2Weird(shortcut)
	mod, keys, err := keybind.ParseString(X, shortcut)
	if err != nil {
		logger.Errorf("ParseString error: %v", err)
		return
	}

	if mod > 0 || len(keys) < 1 {
		return
	}

	C.ungrab_xrecord_key(C.int(keys[0]))
}

func initXRecord() {
	C.grab_xrecord_init()
}

func stopXRecord() {
	C.grab_xrecord_finalize()
}

func parseKeyEnvent(X *xgbutil.XUtil, state uint16, detail xproto.Keycode) string {
	modStr := keybind.ModifierString(state)
	keyStr := keybind.LookupString(X, state, detail)

	if detail == 65 {
		keyStr = "Space"
	}

	if keyStr == "L1" {
		keyStr = "F11"
	}

	if keyStr == "L2" {
		keyStr = "F12"
	}

	value := ""
	modStr = deleteSpecialMod(modStr)

	if len(modStr) > 0 {
		value = convertModsToKeys(modStr) + ACCEL_DELIM + keyStr
	} else {
		value = keyStr
	}

	return value
}

func parseModifierKey(mod, key string) string {
	tmp := strings.ToLower(key)
	if strings.Contains(tmp, "control") {
		return "control"
	} else if strings.Contains(tmp, "alt") {
		return "alt"
	} else if strings.Contains(tmp, "shift") {
		return "shift"
	} else if strings.Contains(tmp, "super") {
		if len(mod) > 0 {
			return "super"
		}
	}

	return key
}
