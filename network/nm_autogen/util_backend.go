package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// NM_SETTING_CONNECTION_SETTING_NAME -> ../nm_setting_connection_autogen.go
func getBackEndFilePath(fieldName string) (filePath string) {
	fieldName = strings.Replace(fieldName, "NM_SETTING_VF_", "NM_SETTING_", -1) // remove virtual field tag
	fileName := strings.TrimSuffix(fieldName, "_SETTING_NAME")
	fileName = strings.ToLower(fileName) + "_autogen.go"
	filePath = path.Join(backEndDir, fileName)
	return
}

// "general" -> "../../../dss/modules/network/components_autogen/EditSectionGeneral.qml"
func getFrontEndFilePath(pageName string) (filePath string) {
	fileName := "EditSection" + ToClassName(pageName) + ".qml"
	filePath = path.Join(frontEndDir, fileName)
	return
}

// "ktypeString" -> "String", "ktypeBoolean" -> "Boolean"
func ToKeyTypeShortName(ktype string) string {
	return strings.TrimPrefix(ktype, "ktype")
}

// "ktypeString" -> "string", "ktypeArrayByte" -> "[]byte"
func ToKeyTypeRealData(ktype string) (realData string) {
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		realData = "string"
	case "ktypeByte":
		realData = "byte"
	case "ktypeInt32":
		realData = "int32"
	case "ktypeUint32":
		realData = "uint32"
	case "ktypeUint64":
		realData = "uint64"
	case "ktypeBoolean":
		realData = "bool"
	case "ktypeArrayByte":
		realData = "[]byte"
	case "ktypeArrayString":
		realData = "[]string"
	case "ktypeArrayUint32":
		realData = "[]uint32"
	case "ktypeArrayArrayByte":
		realData = "[][]byte"
	case "ktypeArrayArrayUint32":
		realData = "[][]uint32"
	case "ktypeDictStringString":
		realData = "map[string]string"
	case "ktypeIpv6Addresses":
		realData = "ipv6Addresses"
	case "ktypeIpv6Routes":
		realData = "ipv6Routes"
	case "ktypeWrapperString":
		realData = "[]byte"
	case "ktypeWrapperMacAddress":
		realData = "[]byte"
	case "ktypeWrapperIpv4Dns":
		realData = "[]uint32"
	case "ktypeWrapperIpv4Addresses":
		realData = "[][]uint32"
	case "ktypeWrapperIpv4Routes":
		realData = "[][]uint32"
	case "ktypeWrapperIpv6Dns":
		realData = "[][]byte"
	case "ktypeWrapperIpv6Addresses":
		realData = "ipv6Addresses"
	case "ktypeWrapperIpv6Routes":
		realData = "ipv6Routes"
	}
	return
}

// "ktypeString" -> `""`, "ktypeBool" -> `false`
func ToKeyTypeDefaultValue(keyName string) (valueJSON string) {
	keyInfo := getKeyInfo(keyName)
	ktype := keyInfo.Type
	customValue := keyInfo.Default
	if customValue == "<null>" {
		if ktype == "ktypeString" {
			return `""`
		} else {
			return "nil"
		}
	} else if customValue != "<default>" {
		return customValue
	}
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		valueJSON = `""`
	case "ktypeByte":
		valueJSON = `0`
	case "ktypeInt32":
		valueJSON = `0`
	case "ktypeUint32":
		valueJSON = `0`
	case "ktypeUint64":
		valueJSON = `0`
	case "ktypeBoolean":
		valueJSON = `false`
	case "ktypeArrayByte":
		valueJSON = `""`
	case "ktypeArrayString":
		valueJSON = `nil`
	case "ktypeArrayUint32":
		valueJSON = `nil`
	case "ktypeArrayArrayByte":
		valueJSON = `nil`
	case "ktypeArrayArrayUint32":
		valueJSON = `nil`
	case "ktypeDictStringString":
		valueJSON = `nil`
	case "ktypeIpv6Addresses":
		valueJSON = `nil`
	case "ktypeIpv6Routes":
		valueJSON = `nil`
	case "ktypeWrapperString":
		valueJSON = `""`
	case "ktypeWrapperMacAddress":
		valueJSON = `""`
	case "ktypeWrapperIpv4Dns":
		valueJSON = `nil`
	case "ktypeWrapperIpv4Addresses":
		valueJSON = `nil`
	case "ktypeWrapperIpv4Routes":
		valueJSON = `nil`
	case "ktypeWrapperIpv6Dns":
		valueJSON = `nil`
	case "ktypeWrapperIpv6Addresses":
		valueJSON = `nil`
	case "ktypeWrapperIpv6Routes":
		valueJSON = `nil`
	}
	return
}

// test if need check value length to ensure value not empty
func IfNeedCheckValueLength(ktype string) (need string) {
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		need = "t"
	case "ktypeByte":
		need = ""
	case "ktypeInt32":
		need = ""
	case "ktypeUint32":
		need = ""
	case "ktypeUint64":
		need = ""
	case "ktypeBoolean":
		need = ""
	case "ktypeArrayByte":
		need = "t"
	case "ktypeArrayString":
		need = "t"
	case "ktypeArrayUint32":
		need = "t"
	case "ktypeArrayArrayByte":
		need = "t"
	case "ktypeArrayArrayUint32":
		need = "t"
	case "ktypeDictStringString":
		need = "t"
	case "ktypeIpv6Addresses":
		need = "t"
	case "ktypeIpv6Routes":
		need = "t"
	case "ktypeWrapperString":
		need = "t"
	case "ktypeWrapperMacAddress":
		need = "t"
	case "ktypeWrapperIpv4Dns":
		need = "t"
	case "ktypeWrapperIpv4Addresses":
		need = "t"
	case "ktypeWrapperIpv4Routes":
		need = "t"
	case "ktypeWrapperIpv6Dns":
		need = "t"
	case "ktypeWrapperIpv6Addresses":
		need = "t"
	case "ktypeWrapperIpv6Routes":
		need = "t"
	}
	return
}

func GetAllVkFields(nmSettingVks []NMSettingVkStruct) (fields []string) {
	for _, vk := range nmSettingVks {
		fields = appendStrArrayUnion(fields, vk.RelatedField)
	}
	return
}

// get all virtual keys in target field
func GetAllVkFieldKeys(nmSettingVks []NMSettingVkStruct, field string) (keys []string) {
	for _, vk := range nmSettingVks {
		if vk.RelatedField == field {
			keys = append(keys, vk.Name)
		}
	}
	return
}

// TODO remove
// func IsVkNeedLogicSetter(nmSettingVks []NMSettingVkStruct, keyName string) bool {
// 	for _, vk := range nmSettingVks {
// 		if vk.Name == keyName {
// 			return vk.LogicSet
// 		}
// 	}
// 	fmt.Println("invalid virtual key:", keyName)
// 	os.Exit(1)
// 	return false
// }
