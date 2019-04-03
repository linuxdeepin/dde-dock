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
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/notify"
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

var _notifier *notify.Notification

func sendNotify(summary, body, icon string) {
	if _notifier == nil {
		notify.Init("dde-daemon")
		_notifier = notify.NewNotification(summary, body, icon)
	} else {
		_notifier.Update(summary, body, icon)
	}
	err := _notifier.Show()
	if err != nil {
		logger.Warning("Failed to send notify:", summary, body)
	}
}

func handleSetScaleFactorDone() {
	sendNotify(gettext.Tr("Set successfully"),
		gettext.Tr("Please log back in to view the changes"), "dialog-window-scale")
}

func handleSetScaleFactorStarted() {
	sendNotify(gettext.Tr("Display scaling"),
		gettext.Tr("Setting display scaling"), "dialog-window-scale")
}

func (m *Manager) setScreenScaleFactors(factors map[string]float64) error {
	return m.xSettings.SetScreenScaleFactors(dbus.FlagNoAutoStart, factors)
}

func (m *Manager) getScreenScaleFactors() (map[string]float64, error) {
	return m.xSettings.GetScreenScaleFactors(dbus.FlagNoAutoStart)
}
