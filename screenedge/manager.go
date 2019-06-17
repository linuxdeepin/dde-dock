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
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
	"pkg.deepin.io/dde/daemon/common/dsync"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	TopLeft     = "left-up"
	TopRight    = "right-up"
	BottomLeft  = "left-down"
	BottomRight = "right-down"

	dbusServiceName = "com.deepin.daemon.Zone"
	dbusPath        = "/com/deepin/daemon/Zone"
	dbusInterface   = "com.deepin.daemon.Zone"

	wmDBusServiceName = "com.deepin.wm"
)

type Manager struct {
	service        *dbusutil.Service
	settings       *Settings
	wm             *wm.Wm
	sessionSigLoop *dbusutil.SignalLoop
	syncConfig     *dsync.Config

	methods *struct {
		EnableZoneDetected func() `in:"enabled"`
		SetTopLeft         func() `in:"value"`
		TopLeftAction      func() `out:"value"`
		SetBottomLeft      func() `in:"value"`
		BottomLeftAction   func() `out:"value"`
		SetTopRight        func() `in:"value"`
		TopRightAction     func() `out:"value"`
		SetBottomRight     func() `in:"value"`
		BottomRightAction  func() `out:"value"`
	}
}

func newManager(service *dbusutil.Service) *Manager {
	var m = new(Manager)
	m.service = service
	m.settings = NewSettings()
	m.wm = wm.NewWm(service.Conn())
	m.sessionSigLoop = dbusutil.NewSignalLoop(service.Conn(), 10)
	m.sessionSigLoop.Start()
	m.syncConfig = dsync.NewConfig("screen_edge", &syncConfig{m: m},
		m.sessionSigLoop, dbusPath, logger)
	return m
}

func (m *Manager) destroy() {
	m.settings.Destroy()
	m.sessionSigLoop.Stop()
	m.syncConfig.Destroy()
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}
