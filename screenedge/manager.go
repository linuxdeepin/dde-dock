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
