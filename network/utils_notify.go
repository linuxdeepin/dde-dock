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
)

const (
	dbusNotifyDest = "org.freedesktop.Notifications"
	dbusNotifyPath = "/org/freedesktop/Notifications"
)

const (
	notifyIconNetworkConnected     = "network-transmit-receive"
	notifyIconNetworkDisconnected  = "network-error"
	notifyIconNetworkOffline       = "network-offline"
	notifyIconEthernetConnected    = "notification-network-ethernet-connected"
	notifyIconEthernetDisconnected = "notification-network-ethernet-disconnected"
	notifyIconWirelessConnected    = "notification-network-wireless-full"
	notifyIconWirelessDisconnected = "notification-network-wireless-disconnected"
	notifyIconVpnConnected         = "notification-network-vpn-connected"
	notifyIconVpnDisconnected      = "notification-network-vpn-disconnected"
)

var notifier, _ = notifications.NewNotifier(dbusNotifyDest, dbusNotifyPath)

func notify(icon, summary, body string) (err error) {
	if notifier == nil {
		logger.Error("connect to org.freedesktop.Notifications failed")
		return
	}
	_, err = notifier.Notify("Network", 0, icon, summary, body, nil, nil, 0)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func notifyVpnConnected(id string) {
	icon := notifyIconVpnConnected
	notify(icon, Tr("Connected"), id)
}
func notifyVpnDisconnected(id string) {
	icon := notifyIconVpnDisconnected
	notify(icon, Tr("Disconnected"), id)
}
func notifyVpnFailed(id string, reason uint32) {
	icon := notifyIconVpnDisconnected
	msg := vpnErrorTable[reason]
	notify(icon, Tr("Disconnected"), msg)
}

func notifyApModeNotSupport() {
	icon := notifyIconWirelessDisconnected
	notify(icon, Tr("Disconnected"), Tr("Access Point mode is not supported by this device."))
}

func notifyProxyEnabled() {
	// TODO: icon
	notify("", Tr("Network"), Tr("System proxy is set successfully."))
}
func notifyProxyDisabled() {
	// TODO: icon
	notify("", Tr("Network"), Tr("System proxy has been cancelled."))
}

func notifyWirelessHardSwitchOff() {
	icon := notifyIconWirelessDisconnected
	notify(icon, Tr("Disconnected"), Tr("The hardware switch of WLAN Card is off, please switch on as necessary."))
}
