/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/lib/dbus"
)

const (
	TopLeft     = "left-up"
	TopRight    = "right-up"
	BottomLeft  = "left-down"
	BottomRight = "right-down"
)

type Manager struct {
	settings *Settings
}

func NewManager() (*Manager, error) {
	var m = new(Manager)
	m.settings = NewSettings()
	err := dbus.InstallOnSession(m)
	if err != nil {
		m.destroy()
		return nil, err
	}
	return m, nil
}

func (m *Manager) destroy() {
	dbus.UnInstallObject(m)
}
