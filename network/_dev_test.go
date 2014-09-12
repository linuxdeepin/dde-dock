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
	C "launchpad.net/gocheck"
	"os"
	"pkg.linuxdeepin.com/lib/gdkpixbuf"
	. "pkg.linuxdeepin.com/lib/gettext"
	"testing"
	"time"
)

func init() {
	gdkpixbuf.InitGdk()
}

func (*testWrapper) TestNotify(c *C.C) {
	InitI18n()
	Textdomain("dde-daemon")
	initNmStateReasons()

	notify(notifyIconEthernetConnected, Tr("Connected"), "有线网络")
	time.Sleep(1 * time.Second) // more time to wait for the first notification
	snapshotNotify("wired_connected")
	notify(notifyIconEthernetDisconnected, Tr("Disconnected"), "有线网络")
	snapshotNotify("wired_disconnected")
	notify(notifyIconEthernetDisconnected, Tr("Disconnected"), deviceErrorTable[NM_DEVICE_STATE_REASON_CONFIG_FAILED])
	snapshotNotify("wired_error")
	notify(notifyIconWirelessConnected, Tr("Connected"), "linuxdeepin-1")
	snapshotNotify("wireless_connected")
	notify(notifyIconWirelessDisconnected, Tr("Disconnected"), "linuxdeepin-1")
	snapshotNotify("wireless_disconnected")
	notify(notifyIconWirelessDisconnected, Tr("Disconnected"), deviceErrorTable[NM_DEVICE_STATE_REASON_NO_SECRETS])
	snapshotNotify("wireless_error")
	notifyVpnConnected("vpn-pptp")
	snapshotNotify("vpn_connected")
	notifyVpnDisconnected("vpn-pptp")
	snapshotNotify("vpn_disconnected")
	notifyVpnFailed("vpn-pptp", NM_VPN_CONNECTION_STATE_REASON_LOGIN_FAILED)
	snapshotNotify("vpn_error")
	notifyAirplanModeEnabled()
	snapshotNotify("airplanmode")
	notifyNetworkOffline()
	snapshotNotify("offline")
	notifyApModeNotSupport()
	snapshotNotify("apmode_error")
	notifyWirelessHardSwitchOff()
	snapshotNotify("wireless_hard_switch_off")
	notifyProxyEnabled()
	snapshotNotify("proxy_enabled")
	notifyProxyDisabled()
	snapshotNotify("proxy_disabled")
}
func snapshotNotify(suffix string) {
	time.Sleep(1 * time.Second)
	os.MkdirAll("testresult", 0755)
	resultScreenFile := "testdata/test_network_screen.png"
	resultClipFile := "testdata/test_network_" + suffix + ".png"
	gdkpixbuf.ScreenshotImage(resultScreenFile, gdkpixbuf.FormatPng)
	sw, _, _ := gdkpixbuf.GetImageSize(resultScreenFile)
	pading := 35
	w, h := 300-2, 70-2
	gdkpixbuf.ClipImage(resultScreenFile, resultClipFile, sw-pading-w, pading, w, h, gdkpixbuf.FormatPng)
}

func (*testWrapper) TestGetUdevDeviceVendor(c *C.C) {
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
