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

package bluetooth

import (
	"fmt"
	"sync"

	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"
	dbus "pkg.deepin.io/lib/dbus1"
	. "pkg.deepin.io/lib/gettext"
)

const (
	notifyIconBluetoothConnected     = "notification-bluetooth-connected"
	notifyIconBluetoothDisconnected  = "notification-bluetooth-disconnected"
	notifyIconBluetoothConnectFailed = "notification-bluetooth-error"
)

var globalNotifications *notifications.Notifications
var globalNotifyId uint32
var globalNotifyMu sync.Mutex

func initNotifications() error {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	globalNotifications = notifications.NewNotifications(sessionBus)
	return nil
}

func notify(icon, summary, body string) {
	logger.Info("notify", icon, summary, body)

	globalNotifyMu.Lock()
	nid := globalNotifyId
	globalNotifyMu.Unlock()

	nid, err := globalNotifications.Notify(0, "dde-control-center", nid, icon,
		summary, body, nil, nil, -1)
	if err != nil {
		logger.Warning(err)
		return
	}
	globalNotifyMu.Lock()
	globalNotifyId = nid
	globalNotifyMu.Unlock()
}

func notifyConnected(alias string) {
	format := Tr("Connect %q successfully")
	notify(notifyIconBluetoothConnected, "", fmt.Sprintf(format, alias))
}
func notifyDisconnected(alias string) {
	notify(notifyIconBluetoothDisconnected, Tr("Disconnected"), alias)
}

func notifyConnectFailedHostDown(alias string) {
	format := Tr("Make sure %q is turned on and in range")
	notifyConnectFailedAux(alias, format)
}

func notifyConnectFailedAux(alias, format string) {
	notify(notifyIconBluetoothConnectFailed, Tr("Bluetooth connection failed"), fmt.Sprintf(format, alias))
}
