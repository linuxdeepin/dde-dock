/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package bluetooth

import (
	"dbus/org/freedesktop/notifications"
	"fmt"
	. "pkg.deepin.io/lib/gettext"
)

const (
	dbusNotifyDest = "org.freedesktop.Notifications"
	dbusNotifyPath = "/org/freedesktop/Notifications"
)

const (
	notifyIconBluetoothConnected     = "notification-bluetooth-connected"
	notifyIconBluetoothDisconnected  = "notification-bluetooth-disconnected"
	notifyIconBluetoothConnectFailed = "notification-bluetooth-error"
)

func notify(icon, summary, body string) {
	notifier, err := notifications.NewNotifier(dbusNotifyDest, dbusNotifyPath)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("notify", icon, summary, body)
	// use goroutine to fix dbus cycle call issue
	go func() {
		notifier.Notify("Bluetooth", 0, icon, summary, body, nil, nil, 0)
		notifications.DestroyNotifier(notifier)
	}()
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
