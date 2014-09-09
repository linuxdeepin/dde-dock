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
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

func (mManager *MouseManager) listenGSettings() {
	if mManager.listenFlag {
		return
	}
	mManager.listenFlag = true

	mouseSettings.Connect("changed", func(s *gio.Settings, key string) {
		logger.Debug("mouse gsettings changed:", key)
		switch key {
		case MOUSE_KEY_LEFT_HAND:
			leftEnable := mouseSettings.GetBoolean(key)
			mManager.leftHanded(leftEnable)
		case MOUSE_KEY_NATURAL_SCROLL:
			naturalEnable := mouseSettings.GetBoolean(key)
			mManager.naturalScroll(naturalEnable)
		case MOUSE_KEY_DISABLE_TPAD:
			mManager.disableTouchpad(
				mouseSettings.GetBoolean(key))
		case MOUSE_KEY_ACCEL:
			accel := mouseSettings.GetDouble(key)
			mManager.motionAcceleration(accel)
		case MOUSE_KEY_THRES:
			thres := mouseSettings.GetDouble(key)
			mManager.motionThreshold(thres)
		case MOUSE_KEY_DOUBLE_CLICK:
			value := mouseSettings.GetInt(key)
			mManager.doubleClick(uint32(value))
		case MOUSE_KEY_DRAG_THRES:
			value := mouseSettings.GetInt(key)
			mManager.dragThreshold(uint32(value))
		}
	})
}

func (tManager *TouchpadManager) listenGSettings() {
	if tManager.listenFlag {
		return
	}
	tManager.listenFlag = true

	tpadSettings.Connect("changed", func(s *gio.Settings, key string) {
		logger.Debug("touchpad gsettings changed:", key)
		switch key {
		case TPAD_KEY_ENABLE:
			if enable := tpadSettings.GetBoolean(key); enable {
				tManager.init()
				//tManager.enable(enable)
				//leftEnable := tpadSettings.GetBoolean(
				//TPAD_KEY_LEFT_HAND)
				//tManager.leftHanded(leftEnable)
				//tapEnable := tpadSettings.GetBoolean(
				//TPAD_KEY_TAP_CLICK)
				//tManager.tapToClick(tapEnable, leftEnable)
			} else {
				if !GetMouseManager().Exist {
					tpadSettings.SetBoolean(TPAD_KEY_ENABLE, true)
					return
				}

				tManager.enable(false)
				tManager.leftHanded(false)
				tManager.tapToClick(false, false)

				if tManager.typingState {
					tManager.typingExitChan <- true
				}
			}
		case TPAD_KEY_LEFT_HAND:
			if enable := tpadSettings.GetBoolean(
				TPAD_KEY_ENABLE); !enable {
				return
			}

			leftEnable := tpadSettings.GetBoolean(key)
			tapEnable := tpadSettings.GetBoolean(
				TPAD_KEY_TAP_CLICK)

			tManager.leftHanded(leftEnable)
			tManager.tapToClick(tapEnable, leftEnable)
		case TPAD_KEY_TAP_CLICK:
			if enable := tpadSettings.GetBoolean(
				TPAD_KEY_ENABLE); !enable {
				return
			}

			leftEnable := tpadSettings.GetBoolean(
				TPAD_KEY_LEFT_HAND)
			tapEnable := tpadSettings.GetBoolean(key)

			tManager.tapToClick(tapEnable, leftEnable)
		case TPAD_KEY_W_TYPING:
			if enable := tpadSettings.GetBoolean(
				TPAD_KEY_ENABLE); !enable {
				if tManager.typingState {
					tManager.typingExitChan <- true
				}
				return
			}

			typingEnable := tpadSettings.GetBoolean(key)
			tManager.disableTPadWhileTyping(typingEnable)
		case TPAD_KEY_NATURAL_SCROLL, TPAD_KEY_DELTA:
			if enable := tpadSettings.GetBoolean(
				TPAD_KEY_ENABLE); !enable {
				return
			}

			naturalEnable := tpadSettings.GetBoolean(
				TPAD_KEY_NATURAL_SCROLL)
			delta := tpadSettings.GetInt(TPAD_KEY_DELTA)
			tManager.naturalScroll(naturalEnable, int32(delta))
		case TPAD_KEY_EDGE_SCROLL:
			if enable := tpadSettings.GetBoolean(
				TPAD_KEY_ENABLE); !enable {
				return
			}

			edgeEnable := tpadSettings.GetBoolean(key)
			tManager.edgeScroll(edgeEnable)
		case TPAD_KEY_VERT_SCROLL, TPAD_KEY_HORIZ_SCROLL:
			if enable := tpadSettings.GetBoolean(
				TPAD_KEY_ENABLE); !enable {
				return
			}

			vertEnable := tpadSettings.GetBoolean(
				TPAD_KEY_VERT_SCROLL)
			horizEnable := tpadSettings.GetBoolean(
				TPAD_KEY_HORIZ_SCROLL)
			tManager.twoFingerScroll(vertEnable, horizEnable)
		case TPAD_KEY_ACCEL:
			if enable := tpadSettings.GetBoolean(
				TPAD_KEY_ENABLE); !enable {
				return
			}

			accel := tpadSettings.GetDouble(key)
			tManager.motionAcceleration(accel)
		case TPAD_KEY_THRES:
			if enable := tpadSettings.GetBoolean(
				TPAD_KEY_ENABLE); !enable {
				return
			}

			thres := tpadSettings.GetDouble(key)
			tManager.motionThreshold(thres)
		}
	})
}

func (kbdManager *KeyboardManager) listenGSettings() {
	if kbdManager.listenFlag {
		return
	}
	kbdManager.listenFlag = true

	kbdSettings.Connect("changed", func(s *gio.Settings, key string) {
		logger.Debug("keyboard gsettings changed:", key)
		switch key {
		case KBD_KEY_REPEAT_ENABLE, KBD_KEY_DELAY,
			KBD_KEY_REPEAT_INTERVAL:
			enable := kbdSettings.GetBoolean(KBD_KEY_REPEAT_ENABLE)
			delay := kbdSettings.GetUint(KBD_KEY_DELAY)
			interval := kbdSettings.GetUint(KBD_KEY_REPEAT_INTERVAL)
			setKeyboardRepeat(enable,
				uint32(delay), uint32(interval))
		case KBD_KEY_LAYOUT:
			layout := kbdSettings.GetString(key)
			kbdManager.setLayout(layout)
		case KBD_CURSOR_BLINK_TIME:
			value := kbdSettings.GetInt(key)
			kbdManager.setCursorBlink(uint32(value))
		case KBD_KEY_LAYOUT_OPTIONS:
			kbdManager.setLayoutOptions()
			layout := kbdSettings.GetString(KBD_KEY_LAYOUT)
			kbdManager.setLayout(layout)
		case KBD_KEY_USER_LAYOUT_LIST:
			list := kbdSettings.GetStrv(KBD_KEY_USER_LAYOUT_LIST)
			kbdManager.setGreeterLayoutList(list)
		}
	})
}
