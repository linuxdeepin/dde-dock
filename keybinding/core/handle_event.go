/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package core

import "C"

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

//export handleSingleKeyEvent
func handleSingleKeyEvent(keycode, pressed int) {
	if !xrecordEnabled {
		return
	}

	//Don't anything if pressed
	if pressed == 1 {
		return
	}

	key, _ := FormatKeyEvent(0, keycode)
	key = formatAccelToXGB(key)
	if !IsValidSingleKey(key) {
		return
	}

	handler := handlers.GetHandlerByKeycode(0, keycode)
	// Not found
	if handler == nil {
		return
	}

	if isKbdGrabed() {
		return
	}

	// No register handler
	if handler.Handler == nil {
		return
	}

	// key release event as press event
	//handler.Handler(0, keycode, false)
	handler.Handler(0, keycode, true)
}

func handleKeyEvent(state uint16, detail int, pressed bool) {
	key, _ := FormatKeyEvent(state, detail)
	key = formatAccelToXGB(key)
	h := handlers.GetHandler(key)
	if h == nil {
		return
	}

	if h.Handler == nil {
		return
	}

	h.Handler(state, detail, pressed)
}

func listenKeyEvent() {
	xevent.KeyPressFun(
		func(x *xgbutil.XUtil, ev xevent.KeyPressEvent) {
			handleKeyEvent(ev.State, int(ev.Detail), true)
		}).Connect(_xu, _xu.RootWin())

	xevent.KeyReleaseFun(
		func(x *xgbutil.XUtil, ev xevent.KeyReleaseEvent) {
			handleKeyEvent(ev.State, int(ev.Detail), false)
		}).Connect(_xu, _xu.RootWin())
}
