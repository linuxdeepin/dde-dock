/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package screenedge

import (
	"errors"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

// Enable desktop edge zone detected
//
// 是否启用桌面边缘热区功能
func (m *Manager) EnableZoneDetected(enabled bool) *dbus.Error {
	has, err := m.service.NameHasOwner(wmDBusServiceName)
	if err != nil {
		return dbusutil.ToError(err)
	}

	if !has {
		return dbusutil.ToError(errors.New("deepin-wm is not running"))
	}

	err = m.wm.EnableZoneDetected(enabled)
	return dbusutil.ToError(err)
}

// Set left-top edge action
func (m *Manager) SetTopLeft(value string) *dbus.Error {
	m.settings.SetEdgeAction(TopLeft, value)
	return nil
}

// Get left-top edge action
func (m *Manager) TopLeftAction() (string, *dbus.Error) {
	return m.settings.GetEdgeAction(TopLeft), nil
}

// Set left-bottom edge action
func (m *Manager) SetBottomLeft(value string) *dbus.Error {
	m.settings.SetEdgeAction(BottomLeft, value)
	return nil
}

// Get left-bottom edge action
func (m *Manager) BottomLeftAction() (string, *dbus.Error) {
	return m.settings.GetEdgeAction(BottomLeft), nil
}

// Set right-top edge action
func (m *Manager) SetTopRight(value string) *dbus.Error {
	m.settings.SetEdgeAction(TopRight, value)
	return nil
}

// Get right-top edge action
func (m *Manager) TopRightAction() (string, *dbus.Error) {
	return m.settings.GetEdgeAction(TopRight), nil
}

// Set right-bottom edge action
func (m *Manager) SetBottomRight(value string) *dbus.Error {
	m.settings.SetEdgeAction(BottomRight, value)
	return nil
}

// Get right-bottom edge action
func (m *Manager) BottomRightAction() (string, *dbus.Error) {
	return m.settings.GetEdgeAction(BottomRight), nil
}
