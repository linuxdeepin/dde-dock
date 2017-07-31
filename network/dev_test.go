//+build dev

/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"fmt"
	C "gopkg.in/check.v1"
	"os"
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gdkpixbuf"
	. "pkg.deepin.io/lib/gettext"
	"time"
)

func init() {
	gdkpixbuf.InitGdk()
}

func (*testWrapper) TestMain(c *C.C) {
	fmt.Println("Start service...")

	InitI18n()
	Textdomain("dde-daemon")

	daemon := NewDaemon(logger)
	daemon.Start()

	// manager = NewManager()
	// err := dbus.InstallOnSession(manager)
	// if err != nil {
	// 	logger.Error("register dbus interface failed: ", err)
	// 	manager = nil
	// 	os.Exit(1)
	// }
	// manager.initManager()

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		fmt.Printf("Lost dbus: %v", err)
		os.Exit(1)
	}

	logger.Info("dbus connection is closed by user")
	os.Exit(0)
}

func (*testWrapper) TestMaxMatchRules(c *C.C) {
	initDbusDaemon()

	// test dbus max match rules
	for i := 0; i < 1000; i++ {
		dev, err := nmNewDevice("/org/freedesktop/NetworkManager/Devices/0")
		if err == nil {
			nmDestroyDevice(dev)
		}
	}

	dbusDaemon.ConnectNameOwnerChanged(func(name, oldOwner, newOwner string) {
		fmt.Println("catch by dbus wrapper", name)
	})

	fmt.Println("Watch networkmanager restart...")
	dbus.Wait()
}

func (*testWrapper) TestWatchNmConnections(c *C.C) {
	initDbusObjects()
	nmSettings.ConnectNewConnection(func(cpath dbus.ObjectPath) {
		fmt.Println("connection added", cpath)
	})
	nmSettings.ConnectConnectionRemoved(func(cpath dbus.ObjectPath) {
		fmt.Println("connection removed", cpath)
	})
	fmt.Println("Watch nm connections...")
	dbus.Wait()
}

func (*testWrapper) TestWatchNmDeviceState(c *C.C) {
	devPath := dbus.ObjectPath("/org/freedesktop/NetworkManager/Devices/2")
	nmDev, _ := nmNewDevice(devPath)
	defer nmDestroyDevice(nmDev)

	nmDev.ConnectStateChanged(func(newState, oldState, reason uint32) {
		logger.Infof("device state changed, %d => %d, reason[%d]", oldState, newState, reason)
	})

	fmt.Println("Watch nm device state changed: ", devPath)
	dbus.Wait()
}

func (*testWrapper) TestNotify(c *C.C) {
	InitI18n()
	Textdomain("dde-daemon")
	initNmStateReasons()

	notify(notifyIconWiredConnected, Tr("Connected"), "Wired network")
	time.Sleep(1 * time.Second) // more time to wait for the first notification
	snapshotNotify("wired_connected")

	notify(notifyIconWiredDisconnected, Tr("Disconnected"), "Wired network")
	snapshotNotify("wired_disconnected")

	notify(notifyIconWiredLocal, Tr("Disconnected"), "Wired network local only")
	snapshotNotify("wired_local")

	notify(notifyIconWiredError, Tr("Disconnected"), deviceErrorTable[nm.NM_DEVICE_STATE_REASON_CONFIG_FAILED])
	snapshotNotify("wired_error")

	notifyWiredCableUnplugged()
	snapshotNotify("wired_unplugged")

	notify(notifyIconWirelessConnected, Tr("Connected"), "wireless-ssid")
	snapshotNotify("wireless_connected")

	notify(notifyIconWirelessDisconnected, Tr("Disconnected"), "wireless-ssid")
	snapshotNotify("wireless_disconnected")

	notify(notifyIconWirelessError, Tr("Disconnected"), deviceErrorTable[nm.NM_DEVICE_STATE_REASON_NO_SECRETS])
	snapshotNotify("wireless_error")

	notifyApModeNotSupport()
	snapshotNotify("apmode_error")

	notifyWirelessHardSwitchOff()
	snapshotNotify("wireless_hard_switch_off")

	notifyVpnConnected("vpn-pptp")
	snapshotNotify("vpn_connected")

	notifyVpnDisconnected("vpn-pptp")
	snapshotNotify("vpn_disconnected")

	notifyVpnFailed("vpn-pptp", nm.NM_VPN_CONNECTION_STATE_REASON_LOGIN_FAILED)
	snapshotNotify("vpn_error")

	notifyNetworkOffline()
	snapshotNotify("offline")

	notifyAirplanModeEnabled()
	snapshotNotify("airplanmode")

	notifyProxyEnabled()
	snapshotNotify("proxy_enabled")

	notifyProxyDisabled()
	snapshotNotify("proxy_disabled")
}
func snapshotNotify(suffix string) {
	time.Sleep(1 * time.Second)
	os.MkdirAll("testresult", 0755)
	resultScreenFile := "testresult/test_network_screen.png"
	resultClipFile := "testresult/test_network_" + suffix + ".png"
	gdkpixbuf.ScreenshotImage(resultScreenFile, gdkpixbuf.FormatPng)
	sw, _, _ := gdkpixbuf.GetImageSize(resultScreenFile)
	pading := 35
	w, h := 300-2, 70-2
	gdkpixbuf.ClipImage(resultScreenFile, resultClipFile, sw-pading-w, pading, w, h, gdkpixbuf.FormatPng)
}

func (*testWrapper) TestGetUdevDeviceVendor(c *C.C) {
	initDbusObjects()

	var syspaths []string
	syspaths = append(syspaths, "/sys/devices/pci0000:00/0000:00:01.0")
	for _, p := range nmGetDevices() {
		syspaths = append(syspaths, nmGetDeviceUdi(p))
	}
	syspaths = append(syspaths, "/sys/devices/wrong_format")
	for _, p := range syspaths {
		doPrintDeviceVendor(p)
	}
}
func doPrintDeviceVendor(syspath string) {
	vendor := udevGetDeviceVendor(syspath)
	fmt.Println("device syspath:", syspath)
	fmt.Println("device vendor:", vendor)
	fmt.Println("is usb device:", udevIsUsbDevice(syspath))
	fmt.Println("")
}

func (*testWrapper) TestLocalSupportedVpnTypes(c *C.C) {
	fmt.Println(getLocalSupportedVpnTypes())
}

func (*testWrapper) TestIsWirelessDeviceSuportHotspot(c *C.C) {
	fmt.Println("wlan0 support hotspot mode:", isWirelessDeviceSuportHotspot("wlan0"))
	fmt.Println("wlan1 support hotspot mode:", isWirelessDeviceSuportHotspot("wlan1"))
	fmt.Println("wlan3 support hotspot mode:", isWirelessDeviceSuportHotspot("wlan3"))
}
