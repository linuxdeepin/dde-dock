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

package appearance

import (
	"sync/atomic"

	"github.com/godbus/dbus"
	notifications "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"
	"pkg.deepin.io/lib/gettext"
)

func (m *Manager) getScaleFactor() float64 {
	scaleFactor, err := m.xSettings.GetScaleFactor(dbus.FlagNoAutoStart)
	if err != nil {
		logger.Warning("failed to get scale factor:", err)
		scaleFactor = 1.0
	}
	return scaleFactor
}

func (m *Manager) setScaleFactor(scale float64) error {
	err := m.xSettings.SetScaleFactor(dbus.FlagNoAutoStart, scale)
	if err != nil {
		logger.Warning("failed to set scale factor:", err)
	}
	return err
}

var notifyId uint32

func sendNotify(summary, body, icon string) error {
	sessionConn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	nid := atomic.LoadUint32(&notifyId)
	notifier := notifications.NewNotifications(sessionConn)
	nid, err = notifier.Notify(0, "dde-control-center", nid,
		icon, summary, body,
		nil, nil, -1)
	if err != nil {
		atomic.StoreUint32(&notifyId, nid)
	}
	return err
}

func handleSetScaleFactorDone() {
	err := sendNotify(gettext.Tr("Set successfully"),
		gettext.Tr("Please log back in to view the changes"), "dialog-window-scale")
	if err != nil {
		logger.Warning(err)
	}
}

func handleSetScaleFactorStarted() {
	err := sendNotify(gettext.Tr("Display scaling"),
		gettext.Tr("Setting display scaling"), "dialog-window-scale")
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) setScreenScaleFactors(factors map[string]float64) error {
	return m.xSettings.SetScreenScaleFactors(dbus.FlagNoAutoStart, factors)
}

func (m *Manager) getScreenScaleFactors() (map[string]float64, error) {
	return m.xSettings.GetScreenScaleFactors(dbus.FlagNoAutoStart)
}
