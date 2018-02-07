/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package shortcuts

import (
	"strings"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/keysyms"
)

type XRecordEventHandler struct {
	keySymbols           *keysyms.KeySymbols
	pressedMods          uint16
	historyPressedMods   uint16
	nonModKeyPressed     bool
	modKeyReleasedCb     func(uint8, uint16)
	allModKeysReleasedCb func()
}

func NewXRecordEventHandler(keySymbols *keysyms.KeySymbols) *XRecordEventHandler {
	return &XRecordEventHandler{
		keySymbols: keySymbols,
	}
}

func (h *XRecordEventHandler) logPressedMods(title string) {
	logger.Debug(title, "pressedMods:", Modifiers(h.pressedMods))
}

func (h *XRecordEventHandler) handleButtonEvent(pressed bool) {
	if h.pressedMods > 0 {
		h.nonModKeyPressed = true
	}
}

func (h *XRecordEventHandler) handleKeyEvent(pressed bool, keycode uint8, state uint16) {
	keystr, _ := h.keySymbols.LookupString(x.Keycode(keycode), state)
	//var pr string
	//if pressed {
	//	pr = "PRESS"
	//} else {
	//	pr = "RELEASE"
	//}
	//logger.Debugf("%s keycode: [%d|%s], state: %v\n", pr, keycode, keystr, Modifiers(state))

	if pressed {
		mod, ok := key2Mod(keystr)
		if ok {
			h.pressedMods |= mod
			h.historyPressedMods |= mod
		} else {
			//logger.Debug("non-mod key pressed")
			if h.pressedMods > 0 {
				h.nonModKeyPressed = true
			}
		}
		//h.logPressedMods("pressed")

	} else {
		// release
		//h.logPressedMods("before release")
		mod, ok := key2Mod(keystr)
		if !ok {
			return
		}
		if h.pressedMods == h.historyPressedMods && !h.nonModKeyPressed {
			if h.modKeyReleasedCb != nil {
				logger.Debugf("modKeyReleased keycode %d historyPressedMods: %s",
					keycode, Modifiers(h.historyPressedMods))
				h.modKeyReleasedCb(keycode, h.historyPressedMods)
			}
		}
		h.pressedMods &^= mod
		//h.logPressedMods("after release")

		if h.pressedMods == 0 {
			h.historyPressedMods = 0
			h.nonModKeyPressed = false
			if h.allModKeysReleasedCb != nil {
				logger.Debug("allModKeysReleased")
				h.allModKeysReleasedCb()
			}
		}
	}
}

func keys2Mod(keys []string) (uint16, bool) {
	var ret uint16
	for _, key := range keys {
		mod, ok := key2Mod(key)
		if !ok {
			return 0, false
		}
		ret |= mod
	}
	return ret, true
}

func key2Mod(key string) (uint16, bool) {
	key = strings.ToLower(key)
	// caps_lock and num_lock
	if key == "caps_lock" {
		return keysyms.ModMaskCapsLock, true
	} else if key == "num_lock" {
		return keysyms.ModMaskNumLock, true
	}

	// control/alt/meta/shift/super _ l/r
	parts := strings.Split(key, "_")
	if len(parts) != 2 {
		return 0, false
	}

	if parts[1] != "l" && parts[1] != "r" {
		return 0, false
	}

	switch parts[0] {
	case "shift":
		return keysyms.ModMaskShift, true
	case "control":
		return keysyms.ModMaskControl, true
	case "super":
		return keysyms.ModMaskSuper, true
	case "alt", "meta":
		return keysyms.ModMaskAlt, true
	}
	return 0, false
}
