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

// #cgo pkg-config: gdk-3.0 x11 xi glib-2.0 gio-2.0
// #cgo CFLAGS: -Wall -g
// #cgo LDFLAGS: -lm
// #include <stdlib.h>
// #include "devices.h"
import "C"

import (
	"os/exec"
	"path"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"strings"
	"unsafe"
)

var (
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
	go exec.Command(cmd, args...).Run()
	select {
	case <-tpadTypingChan:
		exec.Command("/usr/bin/killall", "/usr/bin/syndaemon")
		tpadTypingState = false
		return
	}
}

func setLayoutOptions() {
	options := kbdSettings.GetStrv(KBD_KEY_LAYOUT_OPTIONS)
	setLayoutOption("")
	for _, v := range options {
		setLayoutOption(v)
	}
}

func setLayoutOption(option string) bool {
	args := []string{}
	args = append(args, "-option")
	args = append(args, option)
	if err := exec.Command("/usr/bin/setxkbmap", args...).Run(); err != nil {
		logger.Warningf("Set option '%s' failed: %v", option, err)
		return false
	}

	return true
}

func setLayout(key string) {
	layout := ""
	option := ""
	if len(key) < 1 || !strings.Contains(key, ";") {
		layout = "us"
		option = ""
	} else {
		strs := strings.Split(key, ";")
		if len(strs[0]) < 1 {
			layout = "us"
			option = ""
		} else {
			layout = strs[0]
			option = strs[1]
		}
	}

	setLayoutOptions()

	args := []string{}
	args = append(args, "-layout")
	args = append(args, layout)
	args = append(args, "-option")
	args = append(args, option)
	if err := exec.Command("/usr/bin/setxkbmap", args...).Run(); err != nil {
		logger.Warningf("Set Layout: %s - %s Failed: %v",
			layout, option, err)
		return
	}

	key = layout + LAYOUT_DELIM + option
	if len(key) > 0 {
		list := kbdSettings.GetStrv(KBD_KEY_USER_LAYOUT_LIST)
		if !isStrInList(key, list) {
			list = append(list, key)
			kbdSettings.SetStrv(KBD_KEY_USER_LAYOUT_LIST, list)
		}

		username := dutils.GetUserName()
		homeDir := dutils.GetHomeDir()
		if homeDir != path.Join("/tmp", username) {
			if len(username) > 0 {
				greeterObj.SetKbdLayout(username, key)
			}
		}
	}
}

func disableTPadWhileTyping(enable bool) {
	if tpadEnable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !tpadEnable {
		if tpadTypingState {
			tpadTypingChan <- true
		}
		return
	}

	if !enable {
		if tpadTypingState {
			tpadTypingChan <- true
		}
	} else {
		go enableTPadWhileTyping()
	}
}

func setQtCursorBlink(rate uint32) {
	if configDir := dutils.GetConfigDir(); len(configDir) > 0 {
		qtPath := configDir + "/Trolltech.conf"
		dutils.WriteKeyToKeyFile(qtPath, "Qt",
			"cursorFlashTime", rate)
	}
}

func listenDevsSettings() {
	C.listen_device_changed()
	tpadSettings.Connect("changed", func(s *gio.Settings, key string) {
		println("TPad Settings Changed: ", key)
		switch key {
		case TPAD_KEY_ENABLE:
			logger.Infof("%s changed", key)
			if enable := tpadSettings.GetBoolean(key); enable {
				C.set_tpad_enable(C.TRUE)
				ok := tpadSettings.GetBoolean(TPAD_KEY_W_TYPING)
				disableTPadWhileTyping(ok)
			} else {
				if !isDeviceExist("mouse") {
					tpadSettings.SetBoolean(TPAD_KEY_ENABLE, true)
					return
				}
				C.set_tpad_enable(C.FALSE)
				disableTPadWhileTyping(false)
			}
		case TPAD_KEY_LEFT_HAND, TPAD_KEY_TAP_CLICK:
			if enable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !enable {
				return
			}

			tapEnable := C.int(0)
			if ok := tpadSettings.GetBoolean(TPAD_KEY_TAP_CLICK); ok {
				tapEnable = C.int(1)
			}

			leftHand := C.int(0)
			if left := tpadSettings.GetBoolean(TPAD_KEY_LEFT_HAND); left {
				leftHand = C.int(1)
			}
			C.set_tab_to_click(tapEnable, leftHand)
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
		case MOUSE_KEY_DISABLE_TPAD:
			if mouseSettings.GetBoolean(MOUSE_KEY_DISABLE_TPAD) {
				disableTPadWhenMouse()
			} else {
				tpadSettings.SetBoolean(TPAD_KEY_ENABLE, true)
			}
		case MOUSE_KEY_MID_BUTTON:
			if ok := mouseSettings.GetBoolean(key); ok {
				C.set_middle_button(C.TRUE)
			} else {
				C.set_middle_button(C.FALSE)
			}
		case MOUSE_KEY_ACCEL, MOUSE_KEY_THRES:
			thres := int(mouseSettings.GetDouble(MOUSE_KEY_THRES))
			accel := mouseSettings.GetDouble(MOUSE_KEY_ACCEL)
			//logger.Infof("accel: %v, thres: %v", accel, thres)
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
			setLayout(layout)
		case KBD_CURSOR_BLINK_TIME:
			value := kbdSettings.GetInt(key)
			xsObj.SetInterger("Net/CursorBlinkTime", uint32(value))
			setQtCursorBlink(uint32(value))
		case KBD_KEY_LAYOUT_OPTIONS:
			setLayoutOptions()
			layout := kbdSettings.GetString(KBD_KEY_LAYOUT)
			setLayout(layout)
		case KBD_KEY_USER_LAYOUT_LIST:
			list := kbdSettings.GetStrv(KBD_KEY_USER_LAYOUT_LIST)
			username := dutils.GetUserName()
			homeDir := dutils.GetHomeDir()
			if homeDir != path.Join("/tmp", username) {
				if len(username) > 0 {
					greeterObj.SetKbdLayoutList(username, list)
				}
			}
		}
	})
}

func initMouseSettings() {
	// init mouse gsettings value
	if ok := mouseSettings.GetBoolean(MOUSE_KEY_LEFT_HAND); ok {
		C.set_left_handed(C.TRUE)
	} else {
		C.set_left_handed(C.FALSE)
	}

	if ok := mouseSettings.GetBoolean(MOUSE_KEY_MID_BUTTON); ok {
		C.set_middle_button(C.TRUE)
	} else {
		C.set_middle_button(C.FALSE)
	}

	thres := int(mouseSettings.GetDouble(MOUSE_KEY_THRES))
	accel := mouseSettings.GetDouble(MOUSE_KEY_ACCEL)
	mouseName := C.CString("mouse")
	defer C.free(unsafe.Pointer(mouseName))
	C.set_motion(mouseName, C.double(accel), C.int(thres))

	value := mouseSettings.GetInt(MOUSE_KEY_DOUBLE_CLICK)
	xsObj.SetInterger("Net/DoubleClickTime", uint32(value))

	value = mouseSettings.GetInt(MOUSE_KEY_DRAG_THRES)
	xsObj.SetInterger("Net/DndDragThreshold", uint32(value))
}

func initTPadSettings(tpadFlag bool) {
	// init touchpad gsettings value
	if !tpadFlag {
		return
	}
	if enable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); enable {
		C.set_tpad_enable(C.TRUE)
	} else {
		C.set_tpad_enable(C.FALSE)
		return
	}

	tapEnable := C.int(0)
	if ok := tpadSettings.GetBoolean(TPAD_KEY_TAP_CLICK); ok {
		tapEnable = C.int(1)
	}

	leftHand := C.int(0)
	if left := tpadSettings.GetBoolean(TPAD_KEY_LEFT_HAND); left {
		leftHand = C.int(1)
	}
	C.set_tab_to_click(tapEnable, leftHand)

	ok := tpadSettings.GetBoolean(TPAD_KEY_W_TYPING)
	disableTPadWhileTyping(ok)

	if ok := tpadSettings.GetBoolean(TPAD_KEY_NATURAL_SCROLL); ok {
		C.set_natural_scroll(C.TRUE)
	} else {
		C.set_natural_scroll(C.FALSE)
	}

	if ok := tpadSettings.GetBoolean(TPAD_KEY_EDGE_SCROLL); ok {
		C.set_edge_scroll(C.TRUE)
	} else {
		C.set_edge_scroll(C.FALSE)
	}

	vert := C.int(0)
	horiz := C.int(0)
	if ok := tpadSettings.GetBoolean(TPAD_KEY_VERT_SCROLL); ok {
		vert = C.int(1)
	}
	if ok := tpadSettings.GetBoolean(TPAD_KEY_HORIZ_SCROLL); ok {
		horiz = C.int(1)
	}
	C.set_two_finger_scroll(vert, horiz)

	thres := int(tpadSettings.GetDouble(TPAD_KEY_THRES))
	accel := tpadSettings.GetDouble(TPAD_KEY_ACCEL)
	tpadName := C.CString("touchpad")
	defer C.free(unsafe.Pointer(tpadName))
	C.set_motion(tpadName, C.double(accel), C.int(thres))
}

func initKbdSettings() {
	// init keyyboard gsettings value
	layout := kbdSettings.GetString(KBD_KEY_LAYOUT)
	setLayout(layout)

	blinkTime := kbdSettings.GetInt(KBD_CURSOR_BLINK_TIME)
	xsObj.SetInterger("Net/CursorBlinkTime", uint32(blinkTime))
	setQtCursorBlink(uint32(blinkTime))

	enable := kbdSettings.GetBoolean(KBD_KEY_REPEAT_ENABLE)
	delay := kbdSettings.GetUint(KBD_KEY_DELAY)
	interval := kbdSettings.GetUint(KBD_KEY_REPEAT_INTERVAL)
	if enable {
		C.set_keyboard_repeat(C.int(1), C.uint(interval), C.uint(delay))
	} else {
		C.set_keyboard_repeat(C.int(0), C.uint(interval), C.uint(delay))
	}
}

func initGSettingsSet(tpadFlag bool) {
	//logger.Info("Init devices start...")

	initMouseSettings()
	initTPadSettings(tpadFlag)
	initKbdSettings()

	//logger.Info("Init devices end...")
}

func isDeviceExist(device string) bool {
	names := getDeviceNames()

	for _, name := range names {
		if strings.Contains(name, device) {
			return true
		}
	}

	return false
}

func disableTPadWhenMouse() {
	if mouseSettings.GetBoolean(MOUSE_KEY_DISABLE_TPAD) {
		if !isDeviceExist("mouse") {
			mouseSettings.SetBoolean(MOUSE_KEY_DISABLE_TPAD, false)
			return
		}

		//C.set_tpad_enable(C.FALSE)
		tpadSettings.SetBoolean(TPAD_KEY_ENABLE, false)
	}
}
