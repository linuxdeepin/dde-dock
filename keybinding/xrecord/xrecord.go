/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

package xrecord

// #cgo pkg-config: xcb xcb-record
// #include "xrecord.h"
import "C"
import (
	"fmt"
)

var enabled = true

type KeyEventFunc func(pressed bool, keycode uint8, state uint16)
type ButtonEventFunc func(pressed bool)

var KeyEventCallback KeyEventFunc
var ButtonEventCallback ButtonEventFunc

func Initialize() error {
	state := C.xrecord_grab_init()
	if state == 0 {
		go C.xrecord_grab_event_loop_start()
		return nil
	}
	return fmt.Errorf("xrecord init failed code %d", state)
}

func Finalize() {
	C.xrecord_grab_finalize()
}

func Enable(val bool) {
	enabled = val
}

//export handleKeyEvent
func handleKeyEvent(pressed int, keycode uint8, state uint16) {
	if !enabled {
		return
	}

	if KeyEventCallback != nil {
		KeyEventCallback(pressed == 1, keycode, state)
	}
}

//export handleButtonEvent
func handleButtonEvent(pressed int) {
	if !enabled {
		return
	}

	if ButtonEventCallback != nil {
		ButtonEventCallback(pressed == 1)
	}
}
