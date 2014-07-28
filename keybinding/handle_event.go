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
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"strings"
)

var _mediaManager *MediaKeyManager

const (
	CMD_DDE_OSD = "/usr/lib/deepin-daemon/dde-osd "
)

func GetMediaManager() *MediaKeyManager {
	if _mediaManager == nil {
		_mediaManager = &MediaKeyManager{}
	}

	return _mediaManager
}

func (obj *MediaKeyManager) emitMediaSignal(modStr, keyStr string, press bool) bool {
	switch keyStr {
	case "XF86MonBrightnessUp":
		if press {
			go doAction(CMD_DDE_OSD + "--BrightnessUp")
		}
		obj.BrightnessUp(press)
	case "XF86MonBrightnessDown":
		if press {
			go doAction(CMD_DDE_OSD + "--BrightnessDown")
		}
		obj.BrightnessDown(press)
	case "XF86AudioMute":
		if press {
			go doAction(CMD_DDE_OSD + "--AudioMute")
		}
		obj.AudioMute(press)
	case "XF86AudioLowerVolume":
		if press {
			go doAction(CMD_DDE_OSD + "--AudioDown")
		}
		obj.AudioDown(press)
	case "XF86AudioRaiseVolume":
		if press {
			go doAction(CMD_DDE_OSD + "--AudioUp")
		}
		obj.AudioUp(press)
	case "Num_Lock":
		if strings.Contains(modStr, "mod2") {
			if press {
				go doAction(CMD_DDE_OSD + "--NumLockOff")
			}
			obj.NumLockOff(press)
		} else {
			if press {
				go doAction(CMD_DDE_OSD + "--NumLockOn")
			}
			obj.NumLockOn(press)
		}
	case "Caps_Lock":
		if strings.Contains(modStr, "lock") {
			if press {
				go doAction(CMD_DDE_OSD + "--CapsLockOff")
			}
			obj.CapsLockOff(press)
		} else {
			if press {
				go doAction(CMD_DDE_OSD + "--CapsLockOn")
			}
			obj.CapsLockOn(press)
		}
	case "XF86TouchPadOn":
		if press {
			go doAction(CMD_DDE_OSD + "--TouchPadOn")
		}
		obj.TouchPadOn(press)
	case "XF86TouchPadOff":
		if press {
			go doAction(CMD_DDE_OSD + "--TouchPadOff")
		}
		obj.TouchPadOff(press)
	case "XF86Display":
		obj.SwitchMonitors(press)
	case "XF86PowerOff":
		obj.PowerOff(press)
	case "XF86Sleep":
		obj.PowerSleep(press)
	case "p", "P":
		modStr = deleteSpecialMod(modStr)
		if strings.Contains(modStr, "-") {
			return false
		}
		if strings.Contains(modStr, "mod4") {
			if press {
				go doAction(CMD_DDE_OSD + "--SwitchMonitors")
			}
			obj.SwitchMonitors(press)
		} else {
			return false
		}
	case "space":
		modStr = deleteSpecialMod(modStr)
		if strings.Contains(modStr, "-") {
			return false
		}
		if strings.Contains(modStr, "mod4") {
			if press {
				go doAction(CMD_DDE_OSD + "--SwitchLayout")
			}
			obj.SwitchLayout(press)
		} else {
			return false
		}
	case "XF86AudioPlay":
		obj.AudioPlay(press)
	case "XF86AudioPause":
		obj.AudioPause(press)
	case "XF86AudioStop":
		obj.AudioStop(press)
	case "XF86AudioPrev":
		obj.AudioPrevious(press)
	case "XF86AudioNext":
		obj.AudioNext(press)
	case "XF86AudioRewind":
		obj.AudioRewind(press)
	case "XF86AudioForward":
		obj.AudioForward(press)
	case "XF86AudioRepeat":
		obj.AudioRepeat(press)
	case "XF86WWW":
		obj.LaunchBrowser(press)
	case "XF86Mail":
		obj.LaunchEmail(press)
	case "XF86Calculator":
		obj.LaunchCalculator(press)
	default:
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

	strs := strings.Split(action, " ")
	cmd := strs[0]
	args := []string{}

	if len(strs) > 1 {
		args = append(args, strs[1:]...)
	}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		logger.Errorf("Exec '%s' failed: %v", action, err)
	}
}

func (obj *Manager) listenKeyEvents() {
	xevent.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			modStr := keybind.ModifierString(e.State)
			keyStr := keybind.LookupString(X, e.State, e.Detail)
			if e.Detail == 65 {
				keyStr = "space"
			}
			logger.Infof("KeyStr: %s, modStr: %s", keyStr, modStr)
			if !GetMediaManager().emitMediaSignal(modStr, keyStr, true) {
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
			GetMediaManager().emitMediaSignal(modStr, keyStr, false)
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

	sysGSettings.Connect("changed", func(s *gio.Settings, key string) {
		if id, ok := getAccelIdByName(key); ok {
			invalidFlag := false
			if isInvalidConflict(id) {
				invalidFlag = true
			}

			shortcut := getSystemKeyValue(key, false)

			if id >= 0 && id < 300 {
				grabKeyPairs(PrevSystemPairs, false)
				grabKeyPairs(getSystemKeyPairs(), true)
			} else if id >= 600 && id < 1000 {
				if !invalidFlag {
					obj.setCompizSettings(id, key, shortcut)
				}
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
}
