/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package screenedge

import (
	"github.com/BurntSushi/xgbutil/mousebind"
	"regexp"
	"strings"
)

var ddeLauncherCommandRegexp = regexp.MustCompile(`(?i)dde.*launcher`)

func (m *Manager) handleSettingsChanged() {
	m.settings.ConnectChanged(func(key string) {
		if edge, ok := m.edges[key]; ok {
			edge.Command = m.settings.GetEdgeAction(key)
			logger.Debugf("changed %q => %q", key, edge.Command)
		}
	})
}

func (m *Manager) handleDBusSignal() {
	logger.Debug("handleDBusSignal")

	m.display.ConnectPrimaryChanged(func(argv []interface{}) {
		logger.Debug("event PrimaryChanged")
		m.unregisterEdgeAreas()
		m.setEdgeAreas()
		m.registerEdgeAreas()
	})

	m.xmousearea.ConnectCursorInto(func(x, y int32, id string) {
		m.handleCursorSignal(x, y, id, true)
	})

	m.xmousearea.ConnectCursorOut(func(x, y int32, id string) {
		m.handleCursorSignal(x, y, id, false)
	})

	m.xmousearea.ConnectCancelAllArea(func() {
		logger.Debug("event CancelAllArea")
		m.unregisterEdgeAreas()
		m.registerEdgeAreas()
	})

}

func (m *Manager) tryGrabPointer() (bool, error) {
	win := m.xu.RootWin()
	ok, err := mousebind.GrabPointer(m.xu, win, 0, 0)
	defer mousebind.UngrabPointer(m.xu)
	return ok, err
}

// return true 过滤掉，没有动作
func (m *Manager) filterCursorSignal(id string) bool {
	if id != m.areaId {
		logger.Debug("id not eq m.areaId")
		return true
	}

	canGrabPointer, err := m.tryGrabPointer()
	if err == nil {
		if !canGrabPointer {
			logger.Debug("can not grab pointer")
			return true
		}
	} else {
		logger.Warning(err)
	}

	activeWindow, err := getActiveWindow()
	if err != nil {
		return false
	}

	pid, err := getWindowPid(activeWindow)
	if err != nil {
		return false
	}

	blackList := m.settings.GetBlackList()
	if isAppInList(pid, blackList) {
		logger.Debug("active window app in blacklist")
		return true
	}

	isActiveWindowFullscreen, err := isWindowFullscreen(activeWindow)
	if err != nil {
		return false
	}

	if isActiveWindowFullscreen {
		whiteList := m.settings.GetWhiteList()
		if isAppInList(pid, whiteList) {
			logger.Debug("active window is fullscreen, and in whiteList")
			return false
		}
		logger.Debug("active window is fullscreen, and not in whiteList")
		return true
	}

	return false
}

func (m *Manager) handleCursorSignal(x, y int32, id string, into bool) {
	logger.Debugf("handleCursorSignal x: %v, y: %v,id %q\n, into %v", x, y, id, into)
	if !into {
		// mouse move out
		m.timer.Stop()
		return
	}

	if m.filterCursorSignal(id) {
		return
	}

	isLauncherShowing := false
	activeWindow, err := getActiveWindow()
	if err != nil {
		logger.Debugf("getActiveWindow failed %v, but still execAction", err)
	} else {
		windowName, err := getWindowName(activeWindow)
		if err == nil {
			logger.Debugf("active window name is %q", windowName)
			if windowName == "dde-launcher" {
				isLauncherShowing = true
			}
		}
	}

	for _, edge := range m.edges {
		if edge.Area.Contains(x, y) {
			if isLauncherShowing {
				if ddeLauncherCommandRegexp.MatchString(edge.Command) {
					m.execAction(edge)
					return
				}
				logger.Debug("launcher is showing, do not exec action")
				return
			}
			m.execAction(edge)
			return
		}
	}
}

func (m *Manager) execAction(edge *edge) {
	if strings.Contains(edge.Command, "dde.ControlCenter.Toggle") {
		timeout := m.settings.GetDelay()
		logger.Debug("delay execute ", timeout)
		m.timer.Start(timeout, func() {
			edge.ExecAction()
		})
	} else {
		edge.ExecAction()
	}
}
