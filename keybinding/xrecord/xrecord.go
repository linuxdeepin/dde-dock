/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
