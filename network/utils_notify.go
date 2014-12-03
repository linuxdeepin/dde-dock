/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package network

import (
	"dbus/org/freedesktop/notifications"
	. "pkg.linuxdeepin.com/lib/gettext"
	"time"
)

const (
	dbusNotifyDest = "org.freedesktop.Notifications"
	dbusNotifyPath = "/org/freedesktop/Notifications"
)

const (
	notifyIconNetworkOffline       = "notification-network-offline"
	notifyIconWiredConnected       = "notification-network-wired-connected"
	notifyIconWiredDisconnected    = "notification-network-wired-disconnected"
	notifyIconWiredLocal           = "notification-network-wired-local"
	notifyIconWiredError           = notifyIconWiredDisconnected
	notifyIconWirelessConnected    = "notification-network-wireless-full"
	notifyIconWirelessDisconnected = "notification-network-wireless-disconnected"
	notifyIconWirelessDisabled     = "notification-network-wireless-disabled"
	notifyIconWirelessLocal        = "notification-network-wireless-local"
	notifyIconWirelessError        = notifyIconWirelessDisconnected
	notifyIconVpnConnected         = "notification-network-vpn-connected"
	notifyIconVpnDisconnected      = "notification-network-vpn-disconnected"
	notifyIconProxyEnabled         = "notification-network-proxy-enabled"
	notifyIconProxyDisabled        = "notification-network-proxy-disabled"
	notifyIconNetworkConnected     = notifyIconWiredConnected
	notifyIconNetworkDisconnected  = notifyIconWiredDisconnected
)

var (
	notifyEnabled = true
)

func enableNotify() {
	go func() {
		time.Sleep(5 * time.Second)
		notifyEnabled = true
	}()
}
func disableNotify() {
	notifyEnabled = false
}

func notify(icon, summary, body string) {
	notifier, err := notifications.NewNotifier(dbusNotifyDest, dbusNotifyPath)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("notify", icon, summary, body)
	// use goroutine to fix dbus cycle call issue
	go func() {
		notifier.Notify("Network", 0, icon, summary, body, nil, nil, 0)
		notifications.DestroyNotifier(notifier)
	}()
	return
}

func notifyNetworkOffline() {
	notify(notifyIconNetworkOffline, Tr("Disconnected"), Tr("You are now offline."))
}

func notifyAirplanModeEnabled() {
	// TODO: airplan icon
	notify(notifyIconNetworkOffline, Tr("Disconnected"), Tr("Airplan mode enabled."))
}

func notifyWiredCableUnplugged() {
	notify(notifyIconWiredError, Tr("Disconnected"), deviceErrorTable[GUESS_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED])
}

func notifyApModeNotSupport() {
	notify(notifyIconWirelessError, Tr("Disconnected"), Tr("Access Point mode is not supported by this device."))
}

func notifyWirelessHardSwitchOff() {
	notify(notifyIconWirelessDisabled, Tr("Network"), Tr("The hardware switch of WLAN Card is off, please switch on as necessary."))
}

func notifyProxyEnabled() {
	notify(notifyIconProxyEnabled, Tr("Network"), Tr("System proxy is set successfully."))
}
func notifyProxyDisabled() {
	notify(notifyIconProxyDisabled, Tr("Network"), Tr("System proxy has been cancelled."))
}

func notifyVpnConnected(id string) {
	notify(notifyIconVpnConnected, Tr("Connected"), id)
}
func notifyVpnDisconnected(id string) {
	notify(notifyIconVpnDisconnected, Tr("Disconnected"), id)
}
func notifyVpnFailed(id string, reason uint32) {
	notify(notifyIconVpnDisconnected, Tr("Disconnected"), vpnErrorTable[reason])
}
