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

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"os/exec"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	"strings"
)

const (
	CMD_DDE_OSD = "/usr/lib/deepin-daemon/dde-osd "
)

func (obj *MediaKeyManager) emitMediaSignal(modStr, keyStr string, press bool) bool {
	switch keyStr {
	case "XF86MonBrightnessUp":
		if press {
			go doAction(CMD_DDE_OSD + "--BrightnessUp")
		}
		dbus.Emit(obj, "BrightnessUp", press)
	case "XF86MonBrightnessDown":
		if press {
			go doAction(CMD_DDE_OSD + "--BrightnessDown")
		}
		dbus.Emit(obj, "BrightnessDown", press)
	case "XF86AudioMute":
		if press {
			go doAction(CMD_DDE_OSD + "--AudioMute")
		}
		dbus.Emit(obj, "AudioMute", press)
	case "XF86AudioLowerVolume":
		if press {
			go doAction(CMD_DDE_OSD + "--AudioDown")
		}
		dbus.Emit(obj, "AudioDown", press)
	case "XF86AudioRaiseVolume":
		if press {
			go doAction(CMD_DDE_OSD + "--AudioUp")
		}
		dbus.Emit(obj, "AudioUp", press)
	case "Num_Lock":
		if strings.Contains(modStr, "mod2") {
			if press {
				go doAction(CMD_DDE_OSD + "--NumLockOff")
			}
			dbus.Emit(obj, "NumLockOff", press)
		} else {
			if press {
				go doAction(CMD_DDE_OSD + "--NumLockOn")
			}
			dbus.Emit(obj, "NumLockOn", press)
		}
	case "Caps_Lock":
		if strings.Contains(modStr, "lock") {
			if press {
				go doAction(CMD_DDE_OSD + "--CapsLockOff")
			}
			dbus.Emit(obj, "CapsLockOff", press)
		} else {
			if press {
				go doAction(CMD_DDE_OSD + "--CapsLockOn")
			}
			dbus.Emit(obj, "CapsLockOn", press)
		}
	case "XF86TouchpadToggle":
		if press {
			go doAction(CMD_DDE_OSD + "--TouchpadToggle")
		}
		dbus.Emit(obj, "TouchpadToggle", press)
	case "XF86TouchpadOn":
		if press {
			go doAction(CMD_DDE_OSD + "--TouchpadOn")
		}
		dbus.Emit(obj, "TouchpadOn", press)
	case "XF86TouchpadOff":
		if press {
			go doAction(CMD_DDE_OSD + "--TouchpadOff")
		}
		dbus.Emit(obj, "TouchpadOff", press)
	case "XF86Display":
		dbus.Emit(obj, "SwitchMonitors", press)
	case "XF86PowerOff":
		dbus.Emit(obj, "PowerOff", press)
	case "XF86Sleep":
		dbus.Emit(obj, "PowerSleep", press)
	case "p", "P":
		modStr = deleteSpecialMod(modStr)
		if strings.Contains(modStr, "-") {
			return false
		}
		if strings.Contains(modStr, "mod4") {
			if press {
				go doAction(CMD_DDE_OSD + "--SwitchMonitors")
			}
			dbus.Emit(obj, "SwitchMonitors", press)
		} else {
			return false
		}
	case "XF86AudioPlay":
		dbus.Emit(obj, "AudioPlay", press)
	case "XF86AudioPause":
		dbus.Emit(obj, "AudioPause", press)
	case "XF86AudioStop":
		dbus.Emit(obj, "AudioStop", press)
	case "XF86AudioPrev":
		dbus.Emit(obj, "AudioPrevious", press)
	case "XF86AudioNext":
		dbus.Emit(obj, "AudioNext", press)
	case "XF86AudioRewind":
		dbus.Emit(obj, "AudioRewind", press)
	case "XF86AudioForward":
		dbus.Emit(obj, "AudioForward", press)
	case "XF86AudioRepeat":
		dbus.Emit(obj, "AudioRepeat", press)
	case "XF86WWW":
		dbus.Emit(obj, "LaunchBrowser", press)
	case "XF86Mail":
		dbus.Emit(obj, "LaunchEmail", press)
	case "XF86Calculator":
		dbus.Emit(obj, "LaunchCalculator", press)
	default:
		shortcut := ""
		modStr = deleteSpecialMod(modStr)
		if len(modStr) > 1 {
			shortcut = modStr + "-" + keyStr
		} else {
			shortcut = keyStr
		}
		sLayoutList := sysGSettings.GetStrv("switch-layout")
		if len(sLayoutList) < 1 {
			return false
		}
		sLayout := sLayoutList[0]
		sLayout = formatXGBShortcut(sLayout)
		logger.Debugf("GSettings slayout: %v, input shortcut: %v",
			sLayout, shortcut)
		if strings.ToLower(shortcut) == sLayout {
			if press {
				go doAction(CMD_DDE_OSD + "--SwitchLayout")

			}
			dbus.Emit(obj, "SwitchLayout", press)
			return true
		}
		return false
	}

	return true
}

func getExecCommand(info KeycodeInfo) (string, bool) {
	for k, v := range grabKeyBindsMap {
		if info.State == k.State && info.Detail == k.Detail {
			return v, true
		}
	}

	return "", false
}

func doAction(action string) {
	if len(action) < 1 {
		return
	}

	err := exec.Command("/bin/sh", "-c", action).Run()
	if err != nil {
		logger.Debugf("Exec '%s' failed: %v", action, err)
	}
}

func (m *Manager) listenKeyEvents() {
	xevent.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			modStr := keybind.ModifierString(e.State)
			keyStr := keybind.LookupString(X, e.State, e.Detail)
			if e.Detail == 65 {
				keyStr = "space"
			}
			logger.Infof("KeyStr: %s, modStr: %s", keyStr, modStr)
			if !m.mediaKey.emitMediaSignal(modStr, keyStr, true) {
				modStr = deleteSpecialMod(modStr)
				value := ""
				if len(modStr) < 1 {
					value = keyStr
				} else {
					value = modStr + "-" + keyStr
				}

				info, ok := newKeycodeInfo(value)
				if !ok {
					return
				}

				if v, ok := getExecCommand(info); ok {
					// 不然按键会阻塞，直到程序推出
					go doAction(v)
				}
			}
		}).Connect(X, X.RootWin())

	xevent.KeyReleaseFun(
		func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			modStr := keybind.ModifierString(e.State)
			keyStr := keybind.LookupString(X, e.State, e.Detail)
			if e.Detail == 65 {
				keyStr = "space"
			}
			//modStr = deleteSpecialMod(modStr)
			m.mediaKey.emitMediaSignal(modStr, keyStr, false)
		}).Connect(X, X.RootWin())
}

func isKeyNameExist(name string) bool {
	if _, ok := getAccelIdByName(name); ok {
		return true
	}

	return false
}

func updateSystemSettings(key, shortcut string) {
	values := sysGSettings.GetStrv(key)
	if len(values) < 1 {
		sysGSettings.SetStrv(key, []string{})
	} else {
		if values[0] != shortcut {
			values[0] = shortcut
			sysGSettings.SetStrv(key, values)
		}
	}
}

func (obj *Manager) listenSettings() {
	bindGSettings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case BIND_KEY_VALID_LIST:
			obj.setPropConflictValid(getValidConflictList())
		case BIND_KEY_INVALID_LIST:
			obj.setPropConflictInvalid(getInvalidConflictList())
		case BIND_KEY_CUSTOM_LIST:
			obj.setPropCustomList(getCustomListInfo())
		}
	})
	bindGSettings.GetStrv(BIND_KEY_VALID_LIST)

	sysGSettings.Connect("changed", func(s *gio.Settings, key string) {
		if id, ok := getAccelIdByName(key); ok {
			//invalidFlag := false
			if isInvalidConflict(id) {
				//invalidFlag = true
			}

			//shortcut := getSystemKeyValue(key, false)

			if id >= 0 && id < 300 {
				grabKeyPairs(PrevSystemPairs, false)
				grabKeyPairs(getSystemKeyPairs(), true)
			}

			if isIdInSystemList(id) {
				obj.setPropSystemList(getSystemListInfo())
			} else if isIdInWindowList(id) {
				obj.setPropWindowList(getWindowListInfo())
			} else if isIdInWorkspaceList(id) {
				obj.setPropWorkspaceList(getWorkspaceListInfo())
			}
		}
	})
	sysGSettings.GetStrv("file-manager")
}

func (obj *Manager) listenAllCustomSettings() {
	idList := getCustomIdList()

	for _, id := range idList {
		obj.listenCustomSettings(id)
	}
}

func (obj *Manager) listenCustomSettings(id int32) {
	gs := getSettingsById(id)
	if gs == nil {
		return
	}

	// Prevent gs is released
	obj.idSettingsMap[id] = gs

	gs.Connect("changed", func(s *gio.Settings, key string) {
		logger.Infof("'%s' changed", key)
		if key != CUSTOM_KEY_NAME {
			grabKeyPairs(PrevCustomPairs, false)
			grabKeyPairs(getCustomKeyPairs(), true)
		}

		obj.setPropCustomList(getCustomListInfo())
	})
	gs.GetString("name")
}
