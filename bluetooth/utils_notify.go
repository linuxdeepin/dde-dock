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

	. "pkg.deepin.io/lib/gettext"
	libnotify "pkg.deepin.io/lib/notify"
)

const (
	notifyIconBluetoothConnected     = "notification-bluetooth-connected"
	notifyIconBluetoothDisconnected  = "notification-bluetooth-disconnected"
	notifyIconBluetoothConnectFailed = "notification-bluetooth-error"
)

var globalNotification *libnotify.Notification

func init() {
	globalNotification = libnotify.NewNotification("", "", "")
}

func notify(icon, summary, body string) {
	logger.Info("notify", icon, summary, body)
	globalNotification.Update(summary, body, icon)
	err := globalNotification.Show()
	if err != nil {
		logger.Warning(err)
	}
	return
}

func notifyConnected(alias string) {
	notify(notifyIconBluetoothConnected, Tr("Connected"), alias)
}
func notifyDisconnected(alias string) {
	notify(notifyIconBluetoothDisconnected, Tr("Disconnected"), alias)
}

func notifyConnectFailedHostDown(alias string) {
	format := Tr("Make sure %q is turned on and in range")
	notifyConnectFailedAux(alias, format)
}

func notifyBluetoothDeviceIgnored(alias string) {
	format := Tr("Failed to connect %q, automatically ignored")
	notifyConnectFailedAux(alias, format)
}

func notifyConnectFailedResourceUnavailable(alias string) {
	format := Tr("Failed to connect %q, resource temporarily unavailable")
	notifyConnectFailedAux(alias, format)
}

func notifyConnectFailedSoftwareCaused(alias string) {
	format := Tr("Failed to connect %q, software caused connection abort")
	notifyConnectFailedAux(alias, format)
}

func notifyConnectFailedOther(alias string) {
	format := Tr("Failed to connect %q")
	notifyConnectFailedAux(alias, format)
}

func notifyConnectFailedPairing(alias string) {
	format := Tr("Failed to connect %q, pairing failed")
	notifyConnectFailedAux(alias, format)
}

func notifyConnectFailedAux(alias, format string) {
	notify(notifyIconBluetoothConnectFailed, Tr("Bluetooth connection failed"), fmt.Sprintf(format, alias))
}
