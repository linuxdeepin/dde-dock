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

// #cgo pkg-config: gdk-3.0 x11 xi glib-2.0
// #cgo CFLAGS: -Wall -g
// #cgo LDFLAGS: -lm
// #include <stdlib.h>
// #include "devices.h"
import "C"

import (
	"dlib/gio-2.0"
	"os/exec"
	"strings"
	"unsafe"
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
)

var (
	tpadSettings    = gio.NewSettings("com.deepin.dde.touchpad")
	mouseSettings   = gio.NewSettings("com.deepin.dde.mouse")
	kbdSettings     = gio.NewSettings("com.deepin.dde.keyboard")
	tpadTypingChan  = make(chan bool)
	tpadTypingState = false
)

func enableTPadWhileTyping() {
	if tpadTypingState {
		println("syndaemon has running...")
		return
	}
	cmd := "/usr/bin/syndaemon"
	args := []string{}

	args = append(args, "-i")
	args = append(args, "1")
	args = append(args, "-K")
	args = append(args, "-R")

	tpadTypingState = true
	exec.Command(cmd, args...).Run()
	select {
	case <-tpadTypingChan:
		tpadTypingState = false
		return
	}
}

func setLayout(layout, option string) {
	args := []string{}
	args = append(args, "-layout")
	args = append(args, layout)
	args = append(args, "-option")
	args = append(args, option)
	exec.Command("/usr/bin/setxkbmap", args...).Run()
}

func initGdkEnv() {
	C.init_gdk_env()
}

func disableTPadWhileTyping(enable bool) {
	if tpadEnable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !tpadEnable {
		return
	}

	if !enable {
		tpadTypingChan <- true
	}

	go enableTPadWhileTyping()
}

func setQtCursorBlink(rate uint32) {
	if configDir, ok := utilObj.GetConfigDir(); ok {
		qtPath := configDir + "/Trolltech.conf"
		utilObj.WriteKeyToKeyFile(qtPath, "Qt",
			"cursorFlashTime", rate)
	}
}

func listenDevsSettings() {
	tpadSettings.Connect("changed", func(s *gio.Settings, key string) {
		println("TPad Settings Changed: ", key)
		switch key {
		case TPAD_KEY_ENABLE:
			if enable := tpadSettings.GetBoolean(key); enable {
				C.set_tpad_enable(C.TRUE)
			} else {
				C.set_tpad_enable(C.FALSE)
			}
		case TPAD_KEY_LEFT_HAND:
			if enable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !enable {
				return
			}
			if left := tpadSettings.GetBoolean(key); left {
				C.set_tab_to_click(C.int(1), C.TRUE)
			} else {
				C.set_tab_to_click(C.int(1), C.FALSE)
			}
		case TPAD_KEY_W_TYPING:
			if enable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !enable {
				return
			}
			ok := tpadSettings.GetBoolean(key)
			disableTPadWhileTyping(ok)
		case TPAD_KEY_NATURAL_SCROLL:
			if enable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !enable {
				return
			}
			if ok := tpadSettings.GetBoolean(key); ok {
				C.set_natural_scroll(C.TRUE)
			} else {
				C.set_natural_scroll(C.FALSE)
			}
		case TPAD_KEY_EDGE_SCROLL:
			if enable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !enable {
				return
			}
			if ok := tpadSettings.GetBoolean(key); ok {
				C.set_edge_scroll(C.TRUE)
			} else {
				C.set_edge_scroll(C.FALSE)
			}
		case TPAD_KEY_HORIZ_SCROLL:
			if enable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !enable {
				return
			}
			vert := C.int(0)
			if ok := tpadSettings.GetBoolean(TPAD_KEY_VERT_SCROLL); ok {
				vert = C.int(1)
			}
			if ok := tpadSettings.GetBoolean(key); ok {
				C.set_two_finger_scroll(vert, C.TRUE)
			} else {
				C.set_two_finger_scroll(vert, C.FALSE)
			}
		case TPAD_KEY_VERT_SCROLL:
			if enable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !enable {
				return
			}
			horiz := C.int(0)
			if ok := tpadSettings.GetBoolean(TPAD_KEY_HORIZ_SCROLL); ok {
				horiz = C.int(1)
			}
			if ok := tpadSettings.GetBoolean(key); ok {
				C.set_two_finger_scroll(C.TRUE, horiz)
			} else {
				C.set_two_finger_scroll(C.FALSE, horiz)
			}
		case TPAD_KEY_ACCEL, TPAD_KEY_THRES:
			if enable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !enable {
				return
			}
			thres := int(tpadSettings.GetDouble(TPAD_KEY_THRES))
			accel := tpadSettings.GetDouble(TPAD_KEY_ACCEL)
			tpadName := C.CString("touchpad")
			defer C.free(unsafe.Pointer(tpadName))
			C.set_motion(tpadName, C.double(accel), C.int(thres))
		}
	})

	mouseSettings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case MOUSE_KEY_LEFT_HAND:
			if ok := mouseSettings.GetBoolean(key); ok {
				C.set_left_handed(C.TRUE)
			} else {
				C.set_left_handed(C.FALSE)
			}
		case MOUSE_KEY_MID_BUTTON:
			if ok := mouseSettings.GetBoolean(key); ok {
				C.set_middle_button(C.TRUE)
			} else {
				C.set_middle_button(C.FALSE)
			}
		case TPAD_KEY_ACCEL, TPAD_KEY_THRES:
			thres := int(tpadSettings.GetDouble(TPAD_KEY_THRES))
			accel := tpadSettings.GetDouble(TPAD_KEY_ACCEL)
			mouseName := C.CString("mouse")
			defer C.free(unsafe.Pointer(mouseName))
			C.set_motion(mouseName, C.double(accel), C.int(thres))
		case MOUSE_KEY_DOUBLE_CLICK:
			value := mouseSettings.GetInt(key)
			xsObj.SetInterger("Net/DoubleClickTime", uint32(value))
		case MOUSE_KEY_DRAG_THRES:
			value := mouseSettings.GetInt(key)
			xsObj.SetInterger("Net/DndDragThreshold", uint32(value))
		}
	})

	kbdSettings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case KBD_KEY_REPEAT_ENABLE, KBD_KEY_DELAY, KBD_KEY_REPEAT_INTERVAL:
			enable := kbdSettings.GetBoolean(KBD_KEY_REPEAT_ENABLE)
			delay := kbdSettings.GetUint(KBD_KEY_DELAY)
			interval := kbdSettings.GetUint(KBD_KEY_REPEAT_INTERVAL)

			if enable {
				C.set_keyboard_repeat(C.int(1), C.uint(interval), C.uint(delay))
			} else {
				C.set_keyboard_repeat(C.int(0), C.uint(interval), C.uint(delay))
			}
		case KBD_KEY_LAYOUT:
			layout := kbdSettings.GetString(KBD_KEY_LAYOUT)
			if len(layout) < 1 || !strings.Contains(layout, ";") {
				setLayout("us", "")
			} else {
				strs := strings.Split(layout, ";")
				setLayout(strs[0], strs[1])

				list := kbdSettings.GetStrv(KBD_KEY_USER_LAYOUT_LIST)
				if !utilObj.IsElementExist(layout, list) {
					list = append(list, layout)
					kbdSettings.SetStrv(KBD_KEY_USER_LAYOUT_LIST, list)
				}
			}
		case KBD_CURSOR_BLINK_TIME:
			value := kbdSettings.GetInt(key)
			xsObj.SetInterger("Net/CursorBlinkTime", uint32(value))
			setQtCursorBlink(uint32(value))
		}
	})
}
