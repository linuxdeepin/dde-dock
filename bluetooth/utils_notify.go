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

	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"
	"pkg.deepin.io/lib/dbus1"
	. "pkg.deepin.io/lib/gettext"
)

const (
	notifyIconBluetoothConnected     = "notification-bluetooth-connected"
	notifyIconBluetoothDisconnected  = "notification-bluetooth-disconnected"
	notifyIconBluetoothConnectFailed = "notification-bluetooth-error"
)

func notify(icon, summary, body string) {
	sessionConn, err := dbus.SessionBus()
	if err != nil {
		logger.Warning(err)
		return
	}

	notifier := notifications.NewNotifications(sessionConn)
	logger.Info("notify", icon, summary, body)
	notifier.GoNotify(dbus.FlagNoReplyExpected, nil, "Bluetooth", 0, icon,
		summary, body, nil, nil, 0)
	return
}

func notifyBluetoothConnected(alias string) {
	notify(notifyIconBluetoothConnected, Tr("Connected"), alias)
}
func notifyBluetoothDisconnected(alias string) {
	notify(notifyIconBluetoothDisconnected, Tr("Disconnected"), alias)
}
func notifyBluetoothConnectFailed(alias string) {
	format := Tr("Make sure %q is turned on and in range")
	notify(notifyIconBluetoothConnectFailed, Tr("Bluetooth connection failed"), fmt.Sprintf(format, alias))
}
func notifyBluetoothDeviceIgnored(alias string) {
	format := Tr("Failed to connect %q, automatically ignored")
	notify(notifyIconBluetoothConnectFailed, Tr("Bluetooth connection failed"), fmt.Sprintf(format, alias))
}
