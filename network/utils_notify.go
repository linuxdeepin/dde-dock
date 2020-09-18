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

package network

import (
	"container/list"
	"sync"
	"time"

	"github.com/godbus/dbus"
	notifications "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"
	"pkg.deepin.io/dde/daemon/network/nm"
	. "pkg.deepin.io/lib/gettext"
)

const (
	notifyIconNetworkOffline            = "notification-network-offline"
	notifyIconWiredConnected            = "notification-network-wired-connected"
	notifyIconWiredDisconnected         = "notification-network-wired-disconnected"
	notifyIconWiredError                = notifyIconWiredDisconnected
	notifyIconWirelessConnected         = "notification-network-wireless-full"
	notifyIconWirelessDisconnected      = "notification-network-wireless-disconnected"
	notifyIconWirelessDisabled          = "notification-network-wireless-disabled"
	notifyIconWirelessError             = notifyIconWirelessDisconnected
	notifyIconVpnConnected              = "notification-network-vpn-connected"
	notifyIconVpnDisconnected           = "notification-network-vpn-disconnected"
	notifyIconProxyEnabled              = "notification-network-proxy-enabled"
	notifyIconProxyDisabled             = "notification-network-proxy-disabled"
	notifyIconNetworkConnected          = notifyIconWiredConnected
	notifyIconNetworkDisconnected       = notifyIconWiredDisconnected
	notifyIconMobile2gConnected         = "notification-network-mobile-2g-connected"
	notifyIconMobile2gDisconnected      = "notification-network-mobile-2g-disconnected"
	notifyIconMobile3gConnected         = "notification-network-mobile-3g-connected"
	notifyIconMobile3gDisconnected      = "notification-network-mobile-3g-disconnected"
	notifyIconMobile4gConnected         = "notification-network-mobile-4g-connected"
	notifyIconMobile4gDisconnected      = "notification-network-mobile-4g-disconnected"
	notifyIconMobileUnknownConnected    = "notification-network-mobile-unknown-connected"
	notifyIconMobileUnknownDisconnected = "notification-network-mobile-unknown-disconnected"
)

var (
	notifyEnabled       = true
	notification        *notifications.Notifications
	notifyId            uint32
	notifyIdMu          sync.Mutex
	globalNotifyManager *NotifyManager
)

func initNotifyManager() {
	globalNotifyManager = newNotifyManager()
	go globalNotifyManager.loop()
}

func enableNotify() {
	go func() {
		time.Sleep(5 * time.Second)
		notifyEnabled = true
	}()
}
func disableNotify() {
	notifyEnabled = false
}

type notifyMsgQueue list.List

func newNotifyMsgQueue() *notifyMsgQueue {
	l := list.New()
	return (*notifyMsgQueue)(l)
}

func (q *notifyMsgQueue) toList() *list.List {
	return (*list.List)(q)
}

type notifyMsg struct {
	icon    string
	summary string
	body    string
}

// 入队列
func (q *notifyMsgQueue) enqueue(val *notifyMsg) {
	q.toList().PushBack(val)
}

// 出队列
func (q *notifyMsgQueue) dequeue() *notifyMsg {
	l := q.toList()
	elem := l.Front()
	if elem == nil {
		return nil
	}
	val := l.Remove(elem)
	return val.(*notifyMsg)
}

type NotifyManager struct {
	queue *notifyMsgQueue
	mu    sync.Mutex
	cond  *sync.Cond
}

func newNotifyManager() *NotifyManager {
	sm := &NotifyManager{}
	sm.queue = newNotifyMsgQueue()
	sm.cond = sync.NewCond(&sm.mu)
	return sm
}

func (nm *NotifyManager) addMsg(msg *notifyMsg) {
	nm.mu.Lock()
	nm.queue.enqueue(msg)
	nm.cond.Signal()
	nm.mu.Unlock()
}

func (nm *NotifyManager) loop() {
	for {
		nm.mu.Lock()
		msg := nm.queue.dequeue()

		if msg == nil {
			nm.cond.Wait()
		} else {
			_notify(msg.icon, msg.summary, msg.body)
		}
		nm.mu.Unlock()

		if msg != nil {
			time.Sleep(1500 * time.Millisecond) // sleep 1.5 seconds
		}
	}
}

func notify(icon, summary, body string) {
	globalNotifyManager.addMsg(&notifyMsg{
		icon:    icon,
		summary: summary,
		body:    body,
	})
}

func _notify(icon, summary, body string) {
	logger.Debugf("notify icon: %q, summary: %q, body: %q", icon, summary, body)
	if !notifyEnabled {
		logger.Debug("notify disabled")
		return
	}
	notifyIdMu.Lock()
	nid := notifyId
	notifyIdMu.Unlock()

	nid, err := notification.Notify(0, "dde-control-center", nid,
		icon, summary, body, nil, nil, -1)
	if err != nil {
		logger.Warning(err)
		return
	}

	notifyIdMu.Lock()
	notifyId = nid
	notifyIdMu.Unlock()
}

func notifyAirplanModeEnabled() {
	notify(notifyIconNetworkOffline, Tr("Disconnected"), Tr("Airplane mode enabled."))
}

func notifyWiredCableUnplugged() {
	notify(notifyIconWiredError, Tr("Disconnected"), deviceErrorTable[CUSTOM_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED])
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

func getMobileConnectedNotifyIcon(mobileNetworkType string) (icon string) {
	switch mobileNetworkType {
	case moblieNetworkType4G:
		icon = notifyIconMobile4gConnected
	case moblieNetworkType3G:
		icon = notifyIconMobile3gConnected
	case moblieNetworkType2G:
		icon = notifyIconMobile2gConnected
	case moblieNetworkTypeUnknown:
		icon = notifyIconMobileUnknownConnected
	default:
		icon = notifyIconMobileUnknownConnected
	}
	return
}
func getMobileDisconnectedNotifyIcon(mobileNetworkType string) (icon string) {
	switch mobileNetworkType {
	case moblieNetworkType4G:
		icon = notifyIconMobile4gDisconnected
	case moblieNetworkType3G:
		icon = notifyIconMobile3gDisconnected
	case moblieNetworkType2G:
		icon = notifyIconMobile2gDisconnected
	case moblieNetworkTypeUnknown:
		icon = notifyIconMobileUnknownDisconnected
	default:
		icon = notifyIconMobileUnknownDisconnected
	}
	return
}

func generalGetNotifyConnectedIcon(devType uint32, devPath dbus.ObjectPath) (icon string) {
	switch devType {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		icon = notifyIconWiredConnected
	case nm.NM_DEVICE_TYPE_WIFI:
		icon = notifyIconWirelessConnected
	case nm.NM_DEVICE_TYPE_MODEM:
		var mobileNetworkType string
		dev := manager.getDevice(devPath)
		if dev != nil {
			manager.devicesLock.Lock()
			defer manager.devicesLock.Unlock()
			mobileNetworkType = dev.MobileNetworkType
		}
		icon = getMobileConnectedNotifyIcon(mobileNetworkType)
	default:
		icon = notifyIconNetworkConnected
	}
	return
}
func generalGetNotifyDisconnectedIcon(devType uint32, devPath dbus.ObjectPath) (icon string) {
	switch devType {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		icon = notifyIconWiredDisconnected
	case nm.NM_DEVICE_TYPE_WIFI:
		icon = notifyIconWirelessDisconnected
	case nm.NM_DEVICE_TYPE_MODEM:
		var mobileNetworkType string
		dev := manager.getDevice(devPath)
		if dev != nil {
			manager.devicesLock.Lock()
			mobileNetworkType = dev.MobileNetworkType
			manager.devicesLock.Unlock()
		}
		icon = getMobileDisconnectedNotifyIcon(mobileNetworkType)
	default:
		logger.Warning("lost default notify icon for device", getCustomDeviceType(devType))
		icon = notifyIconNetworkDisconnected
	}
	return
}

func notifyDeviceRemoved(devPath dbus.ObjectPath) {
	var devType uint32
	dev := manager.getDevice(devPath)
	if dev != nil {
		manager.devicesLock.Lock()
		devType = dev.nmDevType
		manager.devicesLock.Unlock()
	}
	if !isDeviceTypeValid(devType) {
		return
	}
	icon := generalGetNotifyDisconnectedIcon(devType, devPath)
	msg := deviceErrorTable[nm.NM_DEVICE_STATE_REASON_REMOVED]
	notify(icon, Tr("Disconnected"), msg)
}
