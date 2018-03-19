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
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/utils"
)

func newBasicVpnConnectionData(id, uuid string) (data connectionData) {
	data = make(connectionData)

	addSetting(data, nm.NM_SETTING_CONNECTION_SETTING_NAME)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, nm.NM_SETTING_VPN_SETTING_NAME)

	initSettingSectionIpv4(data)
	setSettingIP4ConfigNeverDefault(data, false)
	return
}

func initBasicSettingSectionVpn(data connectionData, service string) {
	addSetting(data, nm.NM_SETTING_VPN_SETTING_NAME)
	setSettingVpnServiceType(data, service)
	setSettingVpnData(data, make(map[string]string))
	setSettingVpnSecrets(data, make(map[string]string))
}

func getLocalSupportedVpnTypes() (vpnTypes []string) {
	for _, vpnType := range []string{
		connectionVpnL2tp,
		connectionVpnOpenconnect,
		connectionVpnOpenvpn,
		connectionVpnPptp,
		connectionVpnStrongswan,
		connectionVpnVpnc,
	} {
		_, program, _, _ := parseVpnNameFile(getVpnNameFile(vpnType))
		if utils.IsFileExist(program) {
			vpnTypes = append(vpnTypes, vpnType)
		}
	}
	return
}
func getVpnAuthDialogBin(data connectionData) (authdialog string) {
	vpnType := getCustomConnectionType(data)
	return doGetVpnAuthDialogBin(vpnType)
}
func doGetVpnAuthDialogBin(vpnType string) (authdialog string) {
	_, _, authdialog, _ = parseVpnNameFile(getVpnNameFile(vpnType))
	return
}
func getVpnNameFile(vpnType string) (nameFile string) {
	var baseName string
	switch vpnType {
	case connectionVpnL2tp:
		baseName = nmVpnL2tpNameFile
	case connectionVpnOpenconnect:
		baseName = nmVpnOpenconnectNameFile
	case connectionVpnOpenvpn:
		baseName = nmVpnOpenvpnNameFile
	case connectionVpnPptp:
		baseName = nmVpnPptpNameFile
	case connectionVpnStrongswan:
		baseName = nmVpnStrongswanNameFile
	case connectionVpnVpnc:
		baseName = nmVpnVpncNameFile
	default:
		return ""
	}

	for _, dir := range []string{"/etc/NetworkManager/VPN", "/usr/lib/NetworkManager/VPN"} {
		nameFile = filepath.Join(dir, baseName)
		_, err := os.Stat(nameFile)
		if err == nil {
			return nameFile
		}
	}

	return ""
}

func parseVpnNameFile(nameFile string) (service, program, authdialog, properties string) {
	fileContent, err := ioutil.ReadFile(nameFile)
	if err != nil {
		// service file not exists
		return
	}
	return doParseVpnNameFile(string(fileContent))
}
func doParseVpnNameFile(fileContent string) (service, program, authdialog, properties string) {
	serviceReg := regexp.MustCompile("\nservice=(.*)\n")
	programReg := regexp.MustCompile("\nprogram=(.*)\n")
	authdialogReg := regexp.MustCompile("\nauth-dialog=(.*)\n")
	propertiesReg := regexp.MustCompile("\nproperties=(.*)\n")
	service = serviceReg.FindStringSubmatch(fileContent)[1]
	program = programReg.FindStringSubmatch(fileContent)[1]
	authdialog = authdialogReg.FindStringSubmatch(fileContent)[1]
	properties = propertiesReg.FindStringSubmatch(fileContent)[1]
	return
}

func getSettingVpnAvailableKeys(data connectionData) (keys []string) { return }
func getSettingVpnAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	return
}

func isSettingVpnPluginKey(section string) bool {
	// all keys in vpn virtual sections are vpn plugin key
	realSetting := getAliasSettingRealName(section)
	if realSetting == nm.NM_SETTING_VPN_SETTING_NAME && realSetting != section {
		return true
	}
	return false
}
func isSettingVpnPluginSecretKey(section, key string) (isSecret bool) {
	switch section {
	case nm.NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VPN_L2TP_KEY_PASSWORD, nm.NM_SETTING_VPN_L2TP_KEY_IPSEC_PSK:
			isSecret = true
		}
	case nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD, nm.NM_SETTING_VPN_OPENVPN_KEY_CERTPASS:
			isSecret = true
		}
	case nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_PASSWORD:
			isSecret = true
		}
	case nm.NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VPN_PPTP_KEY_PASSWORD:
			isSecret = true
		}
	case nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VPN_STRONGSWAN_KEY_PASSWORD:
			isSecret = true
		}
	case nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD, nm.NM_SETTING_VPN_VPNC_KEY_SECRET:
			isSecret = true
		}
	}
	return
}

// Basic getter and setter for vpn plugin keys
func getSettingVpnPluginKey(data connectionData, section, key string) (value interface{}) {
	value = generalGetSettingDefaultValue(section, key) // get default value firstly
	vpnPluginData, ok := getSettingVpnPluginData(data, section, key)
	if !ok {
		// not exists, just return nil
		logger.Errorf("invalid vpn plugin data: data[%s][%s]", section, key)
		return
	}
	valueStr, ok := vpnPluginData[key]
	if !ok {
		return
	}
	value = unmarshalVpnPluginKey(valueStr, generalGetSettingKeyType(section, key))
	logger.Debugf("getSettingVpnPluginKey: data[%s][%s]=%v", section, key, value)
	return
}
func setSettingVpnPluginKey(data connectionData, section, key string, value interface{}) {
	vpnPluginData, ok := getSettingVpnPluginData(data, section, key)
	if !ok {
		logger.Errorf("invalid vpn plugin data: data[%s][%s]", section, key)
		return
	}
	valueStr := marshalVpnPluginKey(value, generalGetSettingKeyType(section, key))
	vpnPluginData[key] = valueStr
	logger.Debugf("setSettingVpnPluginKey data[%s][%s]=%#v, valueStr=%s", section, key, value, valueStr)
}
func isSettingVpnPluginKeyExists(data connectionData, section, key string) (ok bool) {
	vpnPluginData, ok := getSettingVpnPluginData(data, section, key)
	if !ok {
		return
	}
	_, ok = vpnPluginData[key]
	return
}
func removeSettingVpnPluginKey(data connectionData, section string, keys ...string) {
	vpnSerectData, ok := doGetSettingVpnPluginData(data, true)
	if ok {
		doRemoveSettingVpnPluginKey(vpnSerectData, keys...)
	}
	vpnPluginData, ok := doGetSettingVpnPluginData(data, false)
	if ok {
		doRemoveSettingVpnPluginKey(vpnPluginData, keys...)
	}
}
func removeSettingVpnPluginKeyBut(data connectionData, section string, keys ...string) {
	vpnSerectData, ok := doGetSettingVpnPluginData(data, true)
	if ok {
		doRemoveSettingVpnPluginKeyBut(vpnSerectData, keys...)
	}
	vpnPluginData, ok := doGetSettingVpnPluginData(data, false)
	if ok {
		doRemoveSettingVpnPluginKeyBut(vpnPluginData, keys...)
	}
}
func doRemoveSettingVpnPluginKey(vpnPluginData map[string]string, keys ...string) {
	for _, k := range keys {
		delete(vpnPluginData, k)
	}
}
func doRemoveSettingVpnPluginKeyBut(vpnPluginData map[string]string, keys ...string) {
	for k := range vpnPluginData {
		if !isStringInArray(k, keys) {
			delete(vpnPluginData, k)
		}
	}
}

func getSettingVpnPluginData(data connectionData, section, key string) (vpnPluginData map[string]string, ok bool) {
	if isSettingVpnPluginSecretKey(section, key) {
		vpnPluginData, ok = doGetSettingVpnPluginData(data, true)
	} else {
		vpnPluginData, ok = doGetSettingVpnPluginData(data, false)
	}
	return
}
func doGetSettingVpnPluginData(data connectionData, isSecretKey bool) (vpnPluginData map[string]string, ok bool) {
	vpnSectionData, ok := data[nm.NM_SETTING_VPN_SETTING_NAME]
	if !ok {
		return
	}
	var variantValue dbus.Variant
	if isSecretKey {
		variantValue, ok = vpnSectionData[nm.NM_SETTING_VPN_SECRETS]
		if !ok {
			return
		}
	} else {
		variantValue, ok = vpnSectionData[nm.NM_SETTING_VPN_DATA]
		if !ok {
			return
		}
	}
	vpnPluginData = interfaceToDictStringString(variantValue.Value())
	return
}

// "string" -> "string", 123 -> "123", true -> "yes", false -> "no"
func marshalVpnPluginKey(value interface{}, t ktype) (valueStr string) {
	var err error
	switch t {
	default:
		logger.Error("unknown vpn plugin key type", t)
	case ktypeString:
		valueStr, _ = value.(string)
	case ktypeUint32:
		valueStr, err = marshalJSON(value)
	case ktypeBoolean:
		valueBoolean, _ := value.(bool)
		if valueBoolean {
			valueStr = "yes"
		} else {
			valueStr = "no"
		}
	}
	if err != nil {
		logger.Error(err)
	}
	return
}

// "string" -> "string", "123" -> 123, "yes" -> true, "no" -> false
func unmarshalVpnPluginKey(valueStr string, t ktype) (value interface{}) {
	var err error
	switch t {
	default:
		logger.Error("unknown vpn plugin key type", t)
	case ktypeString:
		value = valueStr
	case ktypeUint32:
		value, err = jsonToKeyValueUint32(valueStr)
	case ktypeBoolean:
		if valueStr == "yes" {
			value = true
		} else if valueStr == "no" {
			value = false
		} else {
			logger.Error("invalid vpn boolean key", valueStr)
		}
	}
	if err != nil {
		logger.Error(err)
	}
	return
}

// Virtual key
func getSettingVkVpnAutoconnect(data connectionData) (value bool) {
	return manager.config.isVpnConnectionAutoConnect(getSettingConnectionUuid(data))
}
func logicSetSettingVkVpnAutoconnect(data connectionData, value bool) (err error) {
	manager.config.setVpnConnectionAutoConnect(getSettingConnectionUuid(data), value)
	return
}
