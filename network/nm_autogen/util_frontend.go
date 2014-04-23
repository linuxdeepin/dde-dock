package main

import (
	"fmt"
	// "os"
)

// func GenFrontEndWidgetInfo(ktype string) (widget string) {

// TODO
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
		// os.Exit(1)
	case "ktypeString":
		widget = "EditLineTextInput"
	// case "ktypeByte":
	// case "ktypeInt32":
	// case "ktypeUint32":
	// case "ktypeUint64":
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
		// case "ktypeWrapperMacAddress":
		// case "ktypeWrapperIpv4Dns":
		// case "ktypeWrapperIpv4Addresses":
		// case "ktypeWrapperIpv4Routes":
		// case "ktypeWrapperIpv6Dns":
		// case "ktypeWrapperIpv6Addresses":
		// case "ktypeWrapperIpv6Routes":
	}
	return
}
