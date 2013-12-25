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
	/*"dlib"*/
	"dlib/dbus"
	"github.com/BurntSushi/xgbutil/xevent"
)

func NewKeyBinding() *Manager {
	m := &Manager{}
	m.CustomBindList = GetCustomIdList()
	m.gsdAccelMap = GetGSDPairs()
	m.customAccelMap = GetCustomPairs()

	ListenKeyList(m)
	ListenCustomKey(m)
	/*ListenGSDKeyChanged(m)*/

	return m
}

func main() {
	binding := NewKeyBinding()
	err := dbus.InstallOnSession(binding)
	if err != nil {
		panic("Binding Get Session Bus Connect Failed")
	}

	kbd := &GrabManager{}
	err = dbus.InstallOnSession(kbd)
	if err != nil {
		panic("kbd Get Session Bus Connect Failed")
	}

	InitGrabKey()

	xevent.Main(X)
}
