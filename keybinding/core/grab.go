/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
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
	"strings"
)

type HandleType func(uint16, int, bool)

var (
	_xu      *xgbutil.XUtil
	handlers Handlers

	loopRun bool
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
	listenKeyEvent()
	loopRun = true
	xevent.Main(_xu)
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
	news := formatAccelToXGB(s)
	if isValidSingleKey(news) {
		return doGrabSingleKey(news)
	}

	return doGrabAccel(news)
}

func ungrabAccel(s string) {
	news := formatAccelToXGB(s)
	if isValidSingleKey(news) {
		doUngrabSingleKey(news)
	}

	doUngrabAccel(news)
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
	h := handlers.GetHandler(s)
	if h != nil {
		return fmt.Errorf("''%s' has been grabed", s)
	}
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

// check single key valid
func isValidSingleKey(key string) bool {
	tmp := strings.ToLower(key)
	switch tmp {
	case "super_l", "super_r":
		return true
	}
	return false
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
