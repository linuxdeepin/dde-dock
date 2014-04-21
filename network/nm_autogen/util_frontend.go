package main

import (
	"fmt"
	"os"
)

// TODO
// "ktypeString" -> "EditLineTextInput", "ktypeBoolean" -> "EditLineSwitchButton"
func genFrontEndWidgetInfo(ktype string) (widget string) {
	// EditLineSwitchButton {
	//     key: "{{$key.Value}}"
	//     text: dsTr("{{$key.Value}}") // TODO
	// }
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		widget = "EditLineTextInput"
	case "ktypeByte":
		// widget = "byte"
	case "ktypeInt32":
		// widget = "int32"
	case "ktypeUint32":
		// widget = "uint32"
	case "ktypeUint64":
		// widget = "uint64"
	case "ktypeBoolean":
		widget = "EditLineSwitchButton"
	case "ktypeArrayByte":
		// widget = "[]byte"
	case "ktypeArrayString":
		// widget = "[]string"
	case "ktypeArrayUint32":
		// widget = "[]uint32"
	case "ktypeArrayArrayByte":
		// widget = "[][]byte"
	case "ktypeArrayArrayUint32":
		// widget = "[][]uint32"
	case "ktypeDictStringString":
		// widget = "map[string]string"
	case "ktypeIpv6Addresses":
		// widget = "Ipv6Addresses"
	case "ktypeIpv6Routes":
		// widget = "Ipv6Routes"
	case "ktypeWrapperString":
		widget = "EditLineTextInput"
	case "ktypeWrapperMacAddress":
		widget = "EditLineIpv4"
	case "ktypeWrapperIpv4Dns":
		widget = "EditLineIpv4"
	case "ktypeWrapperIpv4Addresses":
		// widget = "[][]uint32"
	case "ktypeWrapperIpv4Routes":
		// widget = "[][]uint32"
	case "ktypeWrapperIpv6Dns":
		// widget = "[][]byte"
	case "ktypeWrapperIpv6Addresses":
		widget = "EditLineTextInput"
	case "ktypeWrapperIpv6Routes":
		widget = "EditLineTextInput"
	}
	return
}
