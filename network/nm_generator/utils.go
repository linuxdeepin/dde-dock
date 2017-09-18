/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func yamlUnmarshalFile(file string, value interface{}) {
	yamlContent, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlContent, value)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func writeOutputFile(file, content string) {
	// write to .go file and execute gofmt
	err := ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		fmt.Println("error, write file failed:", err)
		return
	}
	fmt.Println("GEN " + file)
}

func mergeOverrideKeys() {
	for _, okey := range nmOverrideKeys {
		found := false
	out:
		for _, setting := range nmConsts.NMSettings {
			for _, key := range setting.Keys {
				if key.KeyName == okey.KeyName {
					found = true
					if len(okey.Value) != 0 {
						key.Value = okey.Value
					}
					if len(okey.CapcaseName) != 0 {
						key.CapcaseName = okey.CapcaseName
					}
					if len(okey.Type) != 0 {
						key.Type = okey.Type
					}
					if len(okey.DefaultValue) != 0 {
						key.DefaultValue = okey.DefaultValue
					}
					break out
				}
			}
		}
		if !found {
			fmt.Println("invalid override key", okey.KeyName)
			os.Exit(1)
		}
	}
}

// convert setting key's default value from yaml string to go interface
func getSettingKeyFixedDefaultValue(ktype, defaultValueYAML string) (fixedDefaultValue interface{}) {
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		var fixedValue string
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `""`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeByte":
		var fixedValue byte
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `0`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeInt32":
		var fixedValue int32
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `0`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeUint32":
		var fixedValue uint32
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `0`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeInt64":
		var fixedValue int64
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `0`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeUint64":
		var fixedValue uint64
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `0`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeBoolean":
		var fixedValue bool
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `False`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeArrayByte", "ktypeWrapperString", "ktypeWrapperMacAddress":
		var fixedValue []byte
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `[]`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeArrayString":
		var fixedValue []string
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `[]`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeArrayUint32", "ktypeWrapperIpv4Dns":
		var fixedValue []uint32
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `[]`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeArrayArrayByte", "ktypeWrapperIpv6Dns":
		var fixedValue [][]byte
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `[]`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeArrayArrayUint32", "ktypeWrapperIpv4Addresses", "ktypeWrapperIpv4Routes":
		var fixedValue [][]uint32
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `[]`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeDictStringString":
		var fixedValue map[string]string
		if len(defaultValueYAML) == 0 {
			defaultValueYAML = `{}`
		}
		yaml.Unmarshal([]byte(defaultValueYAML), &fixedValue)
		fixedDefaultValue = fixedValue
	case "ktypeIpv6Addresses", "ktypeIpv6Routes", "ktypeWrapperIpv6Addresses", "ktypeWrapperIpv6Routes":
		// ignore the combined structure here and it will be filled in GetKeyDefaultValue
	}
	return
}

// get target key's default value in go syntax
func GetKeyDefaultValue(name string) (gocode string) {
	// query networkmanager original keys
	for _, setting := range nmConsts.NMSettings {
		for _, key := range setting.Keys {
			if name == key.KeyName {
				return doGetKeyDefaultValue(getSettingKeyFixedDefaultValue(key.Type, key.DefaultValue), name, key.Type)
			}
		}
	}

	// query virtual keys
	for _, vsection := range nmVirtualSections {
		for _, key := range vsection.Keys {
			if len(key.VKeyInfo.VirtualKeyName) != 0 {
				if name == key.VKeyInfo.VirtualKeyName {
					return doGetKeyDefaultValue(getSettingKeyFixedDefaultValue(key.VKeyInfo.Type, ""), name, key.VKeyInfo.Type)
				}
			}
		}
	}
	fmt.Println("invalid key:", name)
	os.Exit(1)
	return
}

// "ktypeString" -> `""`, "ktypeBool" -> `false`
func doGetKeyDefaultValue(fixedDefaultValue interface{}, name, ktype string) (gocode string) {
	switch ktype {
	default:
		gocode = UnwrapInterface(fixedDefaultValue)
		// convert to byte/int32/int64 types
		switch ktype {
		case "ktypeByte":
			gocode = "byte(" + gocode + ")"
		case "ktypeInt32":
			gocode = "int32(" + gocode + ")"
		case "ktypeUint32":
			gocode = "uint32(" + gocode + ")"
		case "ktypeInt64":
			gocode = "int64(" + gocode + ")"
		case "ktypeUint64":
			gocode = "uint64(" + gocode + ")"
		}
	case "ktypeIpv6Addresses", "ktypeWrapperIpv6Addresses":
		gocode = `make(ipv6Addresses, 0)`
	case "ktypeIpv6Routes", "ktypeWrapperIpv6Routes":
		gocode = `make(ipv6Routes, 0)`
	}
	return
}

func UnwrapInterface(ifc interface{}) (value string) {
	value = fmt.Sprintf("%#v", ifc)
	return
}

// "ktypeString" -> "string", "ktypeArrayByte" -> "[]byte"
func GetKeyTypeGoSyntax(ktype string) (goSyntax string) {
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		goSyntax = "string"
	case "ktypeByte":
		goSyntax = "byte"
	case "ktypeInt32":
		goSyntax = "int32"
	case "ktypeUint32":
		goSyntax = "uint32"
	case "ktypeInt64":
		goSyntax = "int64"
	case "ktypeUint64":
		goSyntax = "uint64"
	case "ktypeBoolean":
		goSyntax = "bool"
	case "ktypeArrayByte", "ktypeWrapperString", "ktypeWrapperMacAddress":
		goSyntax = "[]byte"
	case "ktypeArrayString":
		goSyntax = "[]string"
	case "ktypeArrayUint32", "ktypeWrapperIpv4Dns":
		goSyntax = "[]uint32"
	case "ktypeArrayArrayByte", "ktypeWrapperIpv6Dns":
		goSyntax = "[][]byte"
	case "ktypeArrayArrayUint32", "ktypeWrapperIpv4Addresses", "ktypeWrapperIpv4Routes":
		goSyntax = "[][]uint32"
	case "ktypeDictStringString":
		goSyntax = "map[string]string"
	case "ktypeIpv6Addresses", "ktypeWrapperIpv6Addresses":
		goSyntax = "ipv6Addresses"
	case "ktypeIpv6Routes", "ktypeWrapperIpv6Routes":
		goSyntax = "ipv6Routes"
	}
	return
}

// "ktypeString" -> interfaceToString, "ktypeBool" -> interfaceToBoolean
func GetKeyTypeGoIfcConverterFunc(ktype string) (converter string) {
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		converter = "interfaceToString"
	case "ktypeByte":
		converter = "interfaceToByte"
	case "ktypeInt32":
		converter = "interfaceToInt32"
	case "ktypeUint32":
		converter = "interfaceToUint32"
	case "ktypeInt64":
		converter = "interfaceToInt64"
	case "ktypeUint64":
		converter = "interfaceToUint64"
	case "ktypeBoolean":
		converter = "interfaceToBoolean"
	case "ktypeArrayByte", "ktypeWrapperString", "ktypeWrapperMacAddress":
		converter = "interfaceToArrayByte"
	case "ktypeArrayString":
		converter = "interfaceToArrayString"
	case "ktypeArrayUint32", "ktypeWrapperIpv4Dns":
		converter = "interfaceToArrayUint32"
	case "ktypeArrayArrayByte", "ktypeWrapperIpv6Dns":
		converter = "interfaceToArrayArrayByte"
	case "ktypeArrayArrayUint32", "ktypeWrapperIpv4Addresses", "ktypeWrapperIpv4Routes":
		converter = "interfaceToArrayArrayUint32"
	case "ktypeDictStringString":
		converter = "interfaceToDictStringString"
	case "ktypeIpv6Addresses", "ktypeWrapperIpv6Addresses":
		converter = "interfaceToIpv6Addresses"
	case "ktypeIpv6Routes", "ktypeWrapperIpv6Routes":
		converter = "interfaceToIpv6Routes"
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
	case "ktypeInt64":
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

// check if target key need a logic setter
func IsLogicSetKey(keyID string) (logicSet string) {
	for _, k := range nmLogicSetKeys {
		if keyID == k {
			logicSet = "t"
		}
	}
	return
}

// NM_SETTING_CONNECTION_ID -> SettingConnectionId
func GetKeyFuncBaseName(name string) (funcName string) {
	funcName = strings.TrimPrefix(name, "NM_")
	funcName = strings.Replace(funcName, "_", " ", -1)
	funcName = ToCaplitalize(funcName)
	funcName = strings.Replace(funcName, " ", "", -1)
	return
}

// "hello world" -> "Hello World", "HELLO WORLD" -> "Hello World"
func ToCaplitalize(str string) (capstr string) {
	scaner := bufio.NewScanner(strings.NewReader(str))
	scaner.Split(bufio.ScanWords)
	for scaner.Scan() {
		word := scaner.Text()
		if len(word) > 1 {
			capstr = capstr + " " + strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		} else if len(word) == 1 {
			capstr = capstr + " " + strings.ToUpper(word)
		}
	}
	capstr = strings.TrimSpace(capstr)
	return
}

// "ktypeString" -> "String", "ktypeBoolean" -> "Boolean"
func GetKeyTypeShortName(ktype string) string {
	return strings.TrimPrefix(ktype, "ktype")
}

// GetVsRelatedSettings get all related setting values for target virtual setting
func GetVsRelatedSettings(vsname string) (relatedSettings []string) {
	for _, vsetting := range nmVirtualSections {
		if vsname == vsetting.VirtaulSectionName {
			for _, key := range vsetting.Keys {
				if !isStringInArray(key.Section, relatedSettings) {
					relatedSettings = append(relatedSettings, key.Section)
				}
			}
			break
		}
	}
	return
}

func isStringInArray(s string, list []string) bool {
	for _, i := range list {
		if i == s {
			return true
		}
	}
	return false
}
