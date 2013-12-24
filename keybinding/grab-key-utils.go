/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package main

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"os/exec"
	"strings"
)

func GrabKeyPress(wid xproto.Window, shortcut string) bool {
	if len(shortcut) <= 0 {
		return false
	}

	err := keybind.KeyPressFun(
		func(X *xgbutil.XUtil, ev xevent.KeyPressEvent) {
			tmpInfo := NewKeyInfo(ev.State, ev.Detail)
			if ok, v := GetExecAction(tmpInfo); ok {
				ExecCommand(v)
			}
		}).Connect(X, wid, shortcut, true)
	if err != nil {
		fmt.Println("Binding key failed:", err)
		return false
	}

	return true
}

func UngrabKey(wid xproto.Window, shortcut string) bool {
	if len(shortcut) <= 0 {
		return false
	}

	mod, keys, err := keybind.ParseString(X, shortcut)
	if err != nil {
		fmt.Println("Get key info failed:", err)
		return false
	}

	keybind.Ungrab(X, wid, mod, keys[0])

	return true
}

func ExecCommand(value string) {
	cmd := exec.Command(value)
	cmd.Run()
}

func GetXGBShortcut(shortcut string) string {
	/*str := FormatShortcut(shortcut)
	if len(str) <= 0 {
		return ""
	}*/

	value := ""
	array := strings.Split(shortcut, "-")
	for i, v := range array {
		if i != 0 {
			value += "-"
		}

		if v == "alt" || v == "super" ||
			v == "meta" || v == "num_lock" ||
			v == "caps_lock" || v == "hyper" {
			modStr, _ := _ModifierMap[v]
			value += modStr
		} else {
			value += v
		}
	}

	return value
}

/*
 * Input string format: '<Control><Alt>T'
 * Output string format: 'control-alt-t'
 */

func FormatShortcut(shortcut string) string {
	l := len(shortcut)

	if l <= 0 {
		fmt.Println("args null")
		return ""
	}

	str := strings.ToLower(shortcut)
	value := ""
	flag := false
	start := 0
	end := 0

	for i, ch := range str {
		if ch == '<' {
			flag = true
			start = i
		}

		if ch == '>' && flag {
			end = i
			flag = false
			if start != 0 {
				value += "-"
			}

			for j := start + 1; j < end; j++ {
				value += string(str[j])
			}
		}
	}

	if end != l {
		i := 0
		if end > 0 {
			i = end + 1
			value += "-"
		}
		for ; i < l; i++ {
			value += string(str[i])
		}
	}

	array := strings.Split(value, "-")
	value = ""
	for i, v := range array {
		if v == "primary" || v == "control" {
			if !strings.Contains(value, "control") {
				value += "control"
			}
			continue
		}

		if i != 0 {
			value += "-"
		}

		value += v
	}

	return value
}

func NewKeyInfo(state uint16, keycode xproto.Keycode) *GrabKeyInfo {
	return &GrabKeyInfo{State: state, Detail: keycode}
}

func GetExecAction(k1 *GrabKeyInfo) (bool, string) {
	for k, v := range GrabKeyBinds {
		if k1.State == k.State || k1.State == (k.State+2) {
			if k.Detail == k1.Detail {
				return true, v
			}
		}
	}

	return false, ""
}
