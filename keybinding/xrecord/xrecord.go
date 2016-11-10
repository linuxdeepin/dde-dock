/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package xrecord

// #cgo pkg-config: x11 xtst glib-2.0
// #include "xrecord.h"
import "C"
import (
//"fmt"
)

var enabled = true
var keyReleaseCb func(code int)

func Initialize() {
	C.xrecord_grab_init()
}

func Finalize() {
	C.xrecord_grab_finalize()
}

func SetKeyReleaseCallback(fn func(code int)) {
	keyReleaseCb = fn
}

func Enable(val bool) {
	enabled = val
}

//export handleSingleKeyEvent
func handleSingleKeyEvent(keycode, pressed int) {
	if !enabled {
		return
	}

	//Don't anything if pressed
	if pressed == 1 {
		return
	}
	//fmt.Println("handleSingleKeyEvent keycode:", keycode)
	if keyReleaseCb != nil {
		keyReleaseCb(keycode)
	}

}
