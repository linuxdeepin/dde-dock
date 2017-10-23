/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

package gesture

import (
	"dbus/com/deepin/wm"
)

const (
	wmActionShowWorkspace int32 = iota + 1
	wmActionToggleMaximize
	wmActionMinimize
	wmActionShowWindow    = 6
	wmActionShowAllWindow = 7
)

var builtinSets = map[string]func() error{
	"ShowWorkspace":            doShowWorkspace,
	"ToggleMaximize":           doToggleMaximize,
	"Minimize":                 doMinimize,
	"ShowWindow":               doShowWindow,
	"ShowAllWindow":            doShowAllWindow,
	"SwitchApplication":        doSwitchApplication,
	"ReverseSwitchApplication": doReverseSwitchApplication,
	"SwitchWorkspace":          doSwitchWorkspace,
	"ReverseSwitchWorkspace":   doReverseSwitchWorkspace,
}

var wmHandler, _ = wm.NewWm("com.deepin.wm", "/com/deepin/wm")

func doShowWorkspace() error {
	return wmHandler.PerformAction(wmActionShowWorkspace)
}

func doToggleMaximize() error {
	return wmHandler.PerformAction(wmActionToggleMaximize)
}

func doMinimize() error {
	return wmHandler.PerformAction(wmActionMinimize)
}

func doShowWindow() error {
	return wmHandler.PerformAction(wmActionShowWindow)
}

func doShowAllWindow() error {
	return wmHandler.PerformAction(wmActionShowAllWindow)
}

func doSwitchApplication() error {
	return wmHandler.SwitchApplication(false)
}

func doReverseSwitchApplication() error {
	return wmHandler.SwitchApplication(true)
}

func doSwitchWorkspace() error {
	return wmHandler.SwitchToWorkspace(false)
}

func doReverseSwitchWorkspace() error {
	return wmHandler.SwitchToWorkspace(true)
}
