/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package core

// #cgo pkg-config: x11 xtst glib-2.0
// #include "record.h"
import "C"

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"regexp"
	"strings"
	"sync"
)

type HandleType func(uint16, int, bool)

var (
	_xu      *xgbutil.XUtil
	handlers Handlers

	loopRun bool

	xrecordEnabled bool
	xrecordLocker  sync.Mutex
)

/**
 * Must be call before calling other functions
 **/
func Initialize() (*xgbutil.XUtil, error) {
	if _xu != nil {
		return _xu, nil
	}

	var err error
	_xu, err = xgbutil.NewConn()
	if err != nil {
		_xu = nil
		return nil, err
	}
	keybind.Initialize(_xu)

	return _xu, nil
}

func Finalize() {
	if _xu == nil {
		return
	}

	C.xrecord_grab_finalize()
	xevent.Quit(_xu)
	loopRun = false
	_xu = nil
}

/**
 * Block
 **/
func StartLoop() {
	if _xu == nil || loopRun {
		return
	}

	C.xrecord_grab_init()
	XRecordEnable(true)

	listenKeyEvent()
	loopRun = true
	xevent.Main(_xu)
}

func XRecordEnable(enabled bool) {
	xrecordLocker.Lock()
	defer xrecordLocker.Unlock()

	xrecordEnabled = enabled
}

func IsShortcutValid(s string) bool {
	tmp := formatAccelToXGB(s)
	if IsValidSingleKey(tmp) {
		return true
	}
	return isValidShortcut(tmp)
}

func GrabAccels(accels []string, handler HandleType) error {
	var grabed []string
	for _, s := range accels {
		err := grabAccel(s)
		if err != nil {
			UngrabAccels(grabed)
			return err
		}

		grabed = append(grabed, s)
	}

	for _, s := range grabed {
		handlers = handlers.AddHandler(NewHandler(s, handler))
	}

	return nil
}

func UngrabAccels(accels []string) {
	for _, s := range accels {
		ungrabAccel(s)
		handlers.DeleteHandler(s)
	}
}

func grabAccel(s string) error {
	if len(s) == 0 {
		return nil
	}

	news := formatAccelToXGB(s)
	if IsValidSingleKey(news) {
		return doGrabSingleKey(news)
	}

	if !isValidShortcut(news) {
		return fmt.Errorf("Invalid shortcut: %v", news)
	}

	return doGrabAccel(news)
}

func ungrabAccel(s string) {
	if len(s) == 0 {
		return
	}

	news := formatAccelToXGB(s)
	if IsValidSingleKey(news) {
		doUngrabSingleKey(news)
	}

	if isValidShortcut(news) {
		doUngrabAccel(news)
	}
}

func doGrabAccel(s string) error {
	mod, codes, err := keybind.ParseString(_xu, s)
	if err != nil {
		return err
	}

	for _, code := range codes {
		err := keybind.GrabChecked(_xu, _xu.RootWin(), mod, code)
		if err != nil {
			doUngrabAccel(s)
			return err
		}
	}
	return nil
}

func doGrabSingleKey(s string) error {
	// h := handlers.GetHandler(s)
	// if h != nil {
	// return fmt.Errorf("'%s' has been grabed", s)
	// }
	return nil
}

func doUngrabSingleKey(s string) {
	return
}

func doUngrabAccel(s string) {
	mod, codes, err := keybind.ParseString(_xu, s)
	if err != nil {
		return
	}

	for _, code := range codes {
		keybind.Ungrab(_xu, _xu.RootWin(), mod, code)
	}
}

// IsValidSingleKey check single key whether valid
func IsValidSingleKey(key string) bool {
	tmp := strings.ToLower(key)
	switch tmp {
	case "super_l", "super_r":
		return true
	}
	return false
}

var mediaReg = regexp.MustCompile(`^xf86`)

func isValidShortcut(shortcut string) bool {
	shortcut = strings.ToLower(shortcut)
	keys := strings.Split(shortcut, accelDelim)
	if len(keys) == 1 {
		if mediaReg.MatchString(shortcut) {
			return true
		}

		switch shortcut {
		case "f1", "f2", "f3", "f4", "f5", "f6",
			"f7", "f8", "f9", "f10", "f11", "f12",
			"caps_lock", "num_lock", "print",
			"backspace", "delete":
			return true
		}
		return false
	}

	key := keys[len(keys)-1]
	// The last key don't contain accel.
	var list = []string{
		"control",
		"shift",
		"super",
		"alt",
		"meta",
		"hyper",
	}

	for _, v := range list {
		if strings.Contains(key, v) {
			return false
		}
	}

	// filter 'shift-xxx'
	if len(keys) == 2 && keys[0] == "shift" {
		return false
	}

	return true
}

func isKeycodesEqual(list1, list2 []xproto.Keycode) bool {
	l1 := len(list1)
	l2 := len(list2)
	if l1 != l2 {
		return false
	}

	for i, v := range list1 {
		if v != list2[i] {
			return false
		}
	}

	return true
}

func isKbdGrabed() bool {
	return (C.is_grabbed() == 1)
}
