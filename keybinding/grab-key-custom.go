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
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
)

func BindingCustomKeys(accelPairs map[string]string, grab bool) {
	for k, v := range accelPairs {
		key := GetXGBShortcut(k)
		mod, keycodes, _ := keybind.ParseString(X, key)
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

func InitGrabKey() {
	var err error

	X, err = xgbutil.NewConn()
	if err != nil {
		fmt.Println("Get New Connection Failed:", err)
		return
	}
	keybind.Initialize(X)

	GrabKeyBinds = make(map[*GrabKeyInfo]string)

	BindingCustomKeys(GetCustomPairs(), true)
}
