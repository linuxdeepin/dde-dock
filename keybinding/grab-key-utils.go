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

func InitGrabKey() {
	var err error

	X, err = xgbutil.NewConn()
	if err != nil {
		fmt.Println("Get New Connection Failed:", err)
		return
	}
	keybind.Initialize(X)

	GrabKeyBinds = make(map[*GrabKeyInfo]string)

	BindingKeysPairs(GetCustomPairs(), true)
	ListenKeyPressEvent()
	BindingKeysPairs(GetGSDPairs(), true)
}

func BindingKeysPairs(accelPairs map[string]string, grab bool) {
	for k, v := range accelPairs {
		/*fmt.Printf("Binding Pairs: key -- %s, value -- %s\n", k, v)*/
		if k == "super" {
			if grab {
				GrabXRecordKey("Super_L", v)
				GrabXRecordKey("Super_R", v)
			} else {
				UngrabXRecordKey("Super_L")
				UngrabXRecordKey("Super_R")
			}
			continue
		}

		key := GetXGBShortcut(k)
		mod, keycodes, _ := keybind.ParseString(X, key)
		/*fmt.Printf("grab mod: %d, key: %d\n", mod, keycodes[0])*/
		if len(keycodes) <= 0 {
			continue
		}
		keyInfo := NewKeyInfo(mod, keycodes[0])
		if grab {
			if GrabKeyPress(X.RootWin(), key) {
				/*fmt.Println("grab mod:", mod)*/
				/*fmt.Println("grab keycodes:", keycodes)*/
				fmt.Println("grab key:", key)
				GrabKeyBinds[keyInfo] = v
			}
		} else {
			/*fmt.Println("ungrab mod:", mod)*/
			/*fmt.Println("ungrab keycodes:", keycodes)*/
			fmt.Println("ungrab key:", key)
			UngrabKey(X.RootWin(), key)
			delete(GrabKeyBinds, keyInfo)
		}
	}
}

func GrabKeyPress(wid xproto.Window, shortcut string) bool {
	if len(shortcut) <= 0 {
		fmt.Println("grab key args failed...")
		return false
	}

	mod, keys, err := keybind.ParseString(X, shortcut)
	if err != nil {
		fmt.Printf("grab parse shortcut string failed: %s\n", err)
		return false
	}

	err = keybind.GrabChecked(X, wid, mod, keys[0])
	if err != nil {
		fmt.Printf("Grab '%s' Failed: %s\n", shortcut, err)
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
		fmt.Printf("ungrab parse shortcut string failed: %s\n", err)
		return false
	}

	keybind.Ungrab(X, wid, mod, keys[0])

	return true
}

func ListenKeyPressEvent() {
	xevent.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			fmt.Printf("State: %d, Detail: %d\n",
				e.State, e.Detail)
			tmpInfo := NewKeyInfo(e.State, e.Detail)
			if ok, v := GetExecAction(tmpInfo); ok {
				// 不然按键会阻塞，直到程序推出
				go ExecCommand(v)
			}
		}).Connect(X, X.RootWin())
}

func ExecCommand(value string) {
	var cmd *exec.Cmd
	vals := strings.Split(value, " ")
	l := len(vals)

	if l > 0 {
		args := []string{}
		for i := 1; i < l; i++ {
			args = append(args, vals[i])
		}
		/*fmt.Println("args: ", args)*/
		cmd = exec.Command(vals[0], args...)
	} else {
		cmd = exec.Command(value)
	}
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Exec command failed:", err)
	}
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
		fmt.Println("format args null")
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
				if i != 0 {
					value += "-"
				}

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
	mod, key := keybind.DeduceKeyInfo(state, keycode)
	return &GrabKeyInfo{State: mod, Detail: key}
}

func CompareKeyInfo(t1, t2 *GrabKeyInfo) bool {
	if t1.State == t2.State && t1.Detail == t2.Detail {
		return true
	}

	return false
}

func GetExecAction(k1 *GrabKeyInfo) (bool, string) {
	for k, v := range GrabKeyBinds {
		if k1.State == k.State && k.Detail == k1.Detail {
			return true, v
		}
	}

	return false, ""
}
