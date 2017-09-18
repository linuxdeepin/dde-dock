/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

// #cgo pkg-config: libudev
// #include <stdlib.h>
// #include "utils_udev.h"
import "C"
import "strings"

import "unsafe"

var deviceDescIgnoredWords = []string{
	"Semiconductor",
	"Components",
	"Corporation",
	"Communications",
	"Company",
	"Corp.",
	"Corp",
	"Co.",
	"Inc.",
	"Inc",
	"Incorporated",
	"Ltd.",
	"Limited.",
	"Intel?",
	"chipset",
	"adapter",
	"[hex]",
	"NDIS",
	"Module",
	"Technology",
	"(Motherboard)",
	"Fast",
}
var deviceDescIgnoredPhrases = []string{
	"Multiprotocol MAC/baseband processor",
	"Wireless LAN Controller",
	"Wireless LAN Adapter",
	"Wireless Adapter",
	"Network Connection",
	"Wireless Cardbus Adapter",
	"Wireless CardBus Adapter",
	"54 Mbps Wireless PC Card",
	"Wireless PC Card",
	"Wireless PC",
	"PC Card with XJACK(r) Antenna",
	"Wireless cardbus",
	"Wireless LAN PC Card",
	"Technology Group Ltd.",
	"Communication S.p.A.",
	"Business Mobile Networks BV",
	"Mobile Broadband Minicard Composite Device",
	"Mobile Communications AB",
	"(PC-Suite Mode)",
	"PCI Express",
	"Ethernet Controller",
	"Ethernet Adapter",
	"(Industrial Computer Source / ICS Advent)",
}

func udevGetDeviceDesc(syspath string) (desc string, ok bool) {
	vendor := fixupDeviceDesc(udevGetDeviceVendor(syspath))
	product := fixupDeviceDesc(udevGetDeviceProduct(syspath))
	if len(vendor) == 0 && len(product) == 0 {
		return "", false
	}

	// If all of the fixed up vendor string is found in product,
	// ignore the vendor.
	if strings.Contains(vendor, product) {
		desc = vendor
	} else {
		desc = vendor + " " + product
	}
	return desc, true
}

func udevGetDeviceVendor(syspath string) (vendor string) {
	cSyspath := C.CString(syspath)
	defer C.free(unsafe.Pointer(cSyspath))

	cVendor := C.get_device_vendor(cSyspath)
	defer C.free(unsafe.Pointer(cVendor))
	vendor = C.GoString(cVendor)
	return
}

func udevGetDeviceProduct(syspath string) (product string) {
	cSyspath := C.CString(syspath)
	defer C.free(unsafe.Pointer(cSyspath))

	cVendor := C.get_device_product(cSyspath)
	defer C.free(unsafe.Pointer(cVendor))
	product = C.GoString(cVendor)
	return
}

func udevIsUsbDevice(syspath string) bool {
	cSyspath := C.CString(syspath)
	defer C.free(unsafe.Pointer(cSyspath))

	ret := C.is_usb_device(cSyspath)
	if ret == 0 {
		return true
	}
	return false
}

// fixupDeviceDesc attempt to shorten description by ignoring certain
// phrases and individual words, such as "Corporation", "Inc".
func fixupDeviceDesc(desc string) (fixedDesc string) {
	desc = strings.Replace(desc, "_", " ", -1)
	desc = strings.Replace(desc, ",", " ", -1)

	for _, phrase := range deviceDescIgnoredPhrases {
		desc = strings.Replace(desc, phrase, "", -1)
	}

	words := strings.Split(desc, " ")
	for _, w := range words {
		if len(w) > 0 && !isStringInArray(w, deviceDescIgnoredWords) {
			if len(fixedDesc) == 0 {
				fixedDesc = w
			} else {
				fixedDesc = fixedDesc + " " + w
			}
		}
	}
	return
}
