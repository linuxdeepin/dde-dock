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

package main

import (
	"fmt"
	"os"
)

// "ktypeString" -> "EditLineTextInput", "ktypeBoolean" -> "EditLineSwitchButton"
func ToFrontEndWidget(keyName string) (widget string) {
	var ktype, customWidget string
	if isVk(keyName) {
		vkInfo := getVkInfo(keyName)
		ktype = vkInfo.Type
		customWidget = vkInfo.FrontEndWidget
	} else {
		keyInfo := getKeyInfo(keyName)
		ktype = keyInfo.Type
		customWidget = keyInfo.FrontEndWidget
	}
	if customWidget != "<default>" {
		return customWidget
	}
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		widget = "EditLineTextInput"
	// case "ktypeByte":
	case "ktypeInt32":
		widget = "EditLineSpinner"
	case "ktypeUint32":
		widget = "EditLineSpinner"
	case "ktypeUint64":
		widget = "EditLineSpinner"
	case "ktypeBoolean":
		widget = "EditLineSwitchButton"
	// case "ktypeArrayByte":
	case "ktypeArrayString":
		widget = "EditLineComboBox"
	// case "ktypeArrayUint32":
	// case "ktypeArrayArrayByte":
	// case "ktypeArrayArrayUint32":
	// case "ktypeDictStringString":
	// case "ktypeIpv6Addresses":
	// case "ktypeIpv6Routes":
	case "ktypeWrapperString":
		widget = "EditLineTextInput"
		// EditLineFileChooser
	case "ktypeWrapperMacAddress":
		widget = "EditLineTextInput"
		// case "ktypeWrapperIpv4Dns":
		// case "ktypeWrapperIpv4Addresses":
		// case "ktypeWrapperIpv4Routes":
		// case "ktypeWrapperIpv6Dns":
		// case "ktypeWrapperIpv6Addresses":
		// case "ktypeWrapperIpv6Routes":
	}
	return
}
