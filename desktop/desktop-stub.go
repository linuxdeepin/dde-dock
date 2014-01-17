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
	"dlib/dbus"
	"dlib/gio-2.0"
)

func (desk *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_DESKTOP_DEST, _DESKTOP_PATH, _DESKTOP_IFC}
}

func (desk *Manager) listenCompizGSettings() {
	_compizIntegrated.Connect("changed::command-11", func(s *gio.Settings, name string) {
		_runCommand11 = s.GetString("command-11")
		desk.getEdgeAction()
	})
	_compizIntegrated.Connect("changed::command-12", func(s *gio.Settings, name string) {
		_runCommand12 = s.GetString("command-12")
		desk.getEdgeAction()
	})
	_compizCommand.Connect("changed::run-command10-edge", func(s *gio.Settings, name string) {
		_runCommandEdge10 = s.GetString("run-command10-edge")
		desk.getEdgeAction()
	})
	_compizCommand.Connect("changed::run-command11-edge", func(s *gio.Settings, name string) {
		_runCommandEdge11 = s.GetString("run-command11-edge")
		desk.getEdgeAction()
	})
	_compizScale.Connect("changed::initiate-edge", func(s *gio.Settings, name string) {
		_scale = s.GetString("initiate-edge")
		desk.getEdgeAction()
	})
}

func (desk *Manager) getEdgeAction() {
	if _runCommand11 == "" && _runCommandEdge10 == "" && _scale == "" {
		desk.TopLeft = ACTION_NONE
	} else if _scale == "TopLeft" && _runCommandEdge10 == "" {
		desk.TopLeft = ACTION_OPENED_WINDOWS
	} else if _runCommand11 == "launcher" && _runCommandEdge10 == "TopLeft" {
		desk.TopLeft = ACTION_LAUNCHER
	}

	if _runCommand12 == "" && _runCommandEdge11 == "" && _scale == "" {
		desk.BottomRight = ACTION_NONE
	} else if _scale == "BottomRight" && _runCommand12 == "" {
		desk.BottomRight = ACTION_OPENED_WINDOWS
	} else if _runCommand12 == "launcher" && _runCommandEdge11 == "BottomRight" {
		desk.BottomRight = ACTION_LAUNCHER
	}

	dbus.NotifyChange(desk, "TopLeft")
	dbus.NotifyChange(desk, "BottomRight")
}
