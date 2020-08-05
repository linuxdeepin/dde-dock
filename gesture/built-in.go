/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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
	"os/exec"
)

const (
	wmActionShowWorkspace int32 = iota + 1
	wmActionToggleMaximize
	wmActionMinimize
	wmActionShowWindow    = 6
	wmActionShowAllWindow = 7
)

const (
	wmTileDirectionLeft uint32 = iota + 1
	wmTileDirectionRight
)

func (m *Manager) initBuiltinSets() {
	m.builtinSets = map[string]func() error{
		"Handle4Or5FingersSwipeUp":   m.doHandle4Or5FingersSwipeUp,
		"Handle4Or5FingersSwipeDown": m.doHandle4Or5FingersSwipeDown,
		"ToggleMaximize":             m.doToggleMaximize,
		"Minimize":                   m.doMinimize,
		"ShowWindow":                 m.doShowWindow,
		"ShowAllWindow":              m.doShowAllWindow,
		"SwitchApplication":          m.doSwitchApplication,
		"ReverseSwitchApplication":   m.doReverseSwitchApplication,
		"SwitchWorkspace":            m.doSwitchWorkspace,
		"ReverseSwitchWorkspace":     m.doReverseSwitchWorkspace,
		"SplitWindowLeft":            m.doTileActiveWindowLeft,
		"SplitWindowRight":           m.doTileActiveWindowRight,
		"MoveWindow":                 m.doMoveActiveWindow,
	}
}

func (m *Manager) toggleShowDesktop() error {
	return exec.Command("/usr/lib/deepin-daemon/desktop-toggle").Run()
}

func (m *Manager) toggleShowMultiTasking() error {
	return m.wm.PerformAction(0, wmActionShowWorkspace)
}

func (m *Manager) doHandle4Or5FingersSwipeUp() error {
	isShowDesktop, err := m.wm.GetIsShowDesktop(0)
	if err != nil {
		return err
	}
	isShowMultiTask, err := m.wm.GetMultiTaskingStatus(0)
	if err != nil {
		return err
	}
	if !isShowDesktop && !isShowMultiTask {
		return m.toggleShowMultiTasking()
	}
	if isShowDesktop && !isShowMultiTask {
		return m.toggleShowDesktop()
	}
	return nil
}

func (m *Manager) doHandle4Or5FingersSwipeDown() error {
	isShowDesktop, err := m.wm.GetIsShowDesktop(0)
	if err != nil {
		return err
	}
	isShowMultiTask, err := m.wm.GetMultiTaskingStatus(0)
	if err != nil {
		return err
	}
	if isShowMultiTask {
		return m.toggleShowMultiTasking()
	}
	if !isShowDesktop {
		return m.toggleShowDesktop()
	}
	return nil
}

func (m *Manager) doToggleMaximize() error {
	return m.wm.PerformAction(0, wmActionToggleMaximize)
}

func (m *Manager) doMinimize() error {
	return m.wm.PerformAction(0, wmActionMinimize)
}

func (m *Manager) doShowWindow() error {
	return m.wm.PerformAction(0, wmActionShowWindow)
}

func (m *Manager) doShowAllWindow() error {
	return m.wm.PerformAction(0, wmActionShowAllWindow)
}

func (m *Manager) doSwitchApplication() error {
	return m.wm.SwitchApplication(0, false)
}

func (m *Manager) doReverseSwitchApplication() error {
	return m.wm.SwitchApplication(0, true)
}

func (m *Manager) doSwitchWorkspace() error {
	return m.wm.SwitchToWorkspace(0, false)
}

func (m *Manager) doReverseSwitchWorkspace() error {
	return m.wm.SwitchToWorkspace(0, true)
}

func (m *Manager) doTileActiveWindowLeft() error {
	return m.wm.TileActiveWindow(0, wmTileDirectionLeft)
}

func (m *Manager) doTileActiveWindowRight() error {
	return m.wm.TileActiveWindow(0, wmTileDirectionRight)
}

func (m *Manager) doMoveActiveWindow() error {
	return m.wm.BeginToMoveActiveWindow(0)
}
