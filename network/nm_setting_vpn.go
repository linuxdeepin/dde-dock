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
	"dlib/dbus"
)

const VPN_NAME_FILES_DIR = "/etc/NetworkManager/VPN/"

// TODO doc
const NM_SETTING_VPN_SETTING_NAME = "vpn"

const (
	NM_SETTING_VPN_SERVICE_TYPE = "service-type"
	NM_SETTING_VPN_USER_NAME    = "user-name"
	NM_SETTING_VPN_DATA         = "data"
	NM_SETTING_VPN_SECRETS      = "secrets"
)

// VPN connection states
const (
	NM_VPN_CONNECTION_STATE_UNKNOWN       = 0
	NM_VPN_CONNECTION_STATE_PREPARE       = 1
	NM_VPN_CONNECTION_STATE_NEED_AUTH     = 2
	NM_VPN_CONNECTION_STATE_CONNECT       = 3
	NM_VPN_CONNECTION_STATE_IP_CONFIG_GET = 4
	NM_VPN_CONNECTION_STATE_ACTIVATED     = 5
	NM_VPN_CONNECTION_STATE_FAILED        = 6
	NM_VPN_CONNECTION_STATE_DISCONNECTE   = 7
)

func newBasicVpnConnectionData(id, uuid string) (data connectionData) {
	data = make(connectionData)

	addSettingSection(data, sectionConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, NM_SETTING_VPN_SETTING_NAME)
	setSettingConnectionAutoconnect(data, false)

	initSettingSectionIpv4(data)
	return
}

func initBasicSettingSectionVpn(data connectionData, service string) {
	addSettingSection(data, sectionVpn)
	setSettingVpnServiceType(data, service)
	setSettingVpnData(data, make(map[string]string))
	setSettingVpnSecrets(data, make(map[string]string))
}

func getLocalSupportedVpnTypes() (vpnTypes []string) {
	if isFileExists(nmVpnL2tpServiceFile) {
		vpnTypes = append(vpnTypes, connectionVpnL2tp)
	}
	if isFileExists(nmVpnOpenconnectServiceFile) {
		vpnTypes = append(vpnTypes, connectionVpnOpenconnect)
	}
	if isFileExists(nmVpnOpenvpnServiceFile) {
		vpnTypes = append(vpnTypes, connectionVpnOpenvpn)
	}
	if isFileExists(nmVpnPptpServiceFile) {
		vpnTypes = append(vpnTypes, connectionVpnPptp)
	}
	if isFileExists(nmVpnVpncServiceFile) {
		vpnTypes = append(vpnTypes, connectionVpnVpnc)
	}
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
	realSection := getRealSectionName(section)
	if realSection == sectionVpn && realSection != section {
		return true
	}
	return false
}
func isSettingVpnPluginSecretKey(section, key string) bool {
	switch section {
	case sectionVpnL2tp:
		switch key {
		case NM_SETTING_VPN_L2TP_KEY_PASSWORD:
			return true
		}
	case sectionVpnOpenvpn:
		switch key {
		case NM_SETTING_VPN_OPENVPN_KEY_PASSWORD:
			return true
		case NM_SETTING_VPN_OPENVPN_KEY_CERTPASS:
			return true
		}
	case sectionVpnOpenvpnProxies:
		switch key {
		case NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_PASSWORD:
			return true
		}
	case sectionVpnPptp:
		switch key {
		case NM_SETTING_VPN_PPTP_KEY_PASSWORD:
			return true
		}
	case sectionVpnVpnc:
		switch key {
		case NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD:
			return true
		case NM_SETTING_VPN_VPNC_KEY_SECRET:
			return true
		}
	}
	return false
}

// Basic getter and setter for vpn plugin keys
func getSettingVpnPluginKey(data connectionData, section, key string) (value interface{}) {
	value = generalGetSettingDefaultValue(section, key) // get default value firstly
	vpnData, ok := getSettingVpnPluginData(data, section, key)
	if !ok {
		// not exists, just return nil
		logger.Errorf("invalid vpn plugin data: data[%s][%s]", section, key)
		return
	}
	valueStr, ok := vpnData[key]
	if !ok {
		return
	}
	value = unmarshalVpnPluginKey(valueStr, generalGetSettingKeyType(section, key))
	return
}
func setSettingVpnPluginKey(data connectionData, section, key string, value interface{}) {
	vpnData, ok := getSettingVpnPluginData(data, section, key)
	if !ok {
		logger.Errorf("invalid vpn plugin data: data[%s][%s]", section, key)
		return
	}
	valueStr := marshalVpnPluginKey(value, generalGetSettingKeyType(section, key))
	vpnData[key] = valueStr
}
func isSettingVpnPluginKeyExists(data connectionData, section, key string) (ok bool) {
	vpnData, ok := getSettingVpnPluginData(data, section, key)
	if !ok {
		return
	}
	_, ok = vpnData[key]
	return
}
func removeSettingVpnPluginKey(data connectionData, section string, keys ...string) {
	vpnSerectData, ok := doGetSettingVpnPluginData(data, true)
	if ok {
		doRemoveSettingVpnPluginKey(vpnSerectData, keys...)
	}
	vpnData, ok := doGetSettingVpnPluginData(data, false)
	if ok {
		doRemoveSettingVpnPluginKey(vpnData, keys...)
	}
}
func removeSettingVpnPluginKeyBut(data connectionData, section string, keys ...string) {
	vpnSerectData, ok := doGetSettingVpnPluginData(data, true)
	if ok {
		doRemoveSettingVpnPluginKeyBut(vpnSerectData, keys...)
	}
	vpnData, ok := doGetSettingVpnPluginData(data, false)
	if ok {
		doRemoveSettingVpnPluginKeyBut(vpnData, keys...)
	}
}
func doRemoveSettingVpnPluginKey(vpnData map[string]string, keys ...string) {
	for _, k := range keys {
		delete(vpnData, k)
	}
}
func doRemoveSettingVpnPluginKeyBut(vpnData map[string]string, keys ...string) {
	for k := range vpnData {
		if !isStringInArray(k, keys) {
			delete(vpnData, k)
		}
	}
}

func getSettingVpnPluginData(data connectionData, section, key string) (vpnData map[string]string, ok bool) {
	if isSettingVpnPluginSecretKey(section, key) {
		vpnData, ok = doGetSettingVpnPluginData(data, true)
	} else {
		vpnData, ok = doGetSettingVpnPluginData(data, false)
	}
	return
}
func doGetSettingVpnPluginData(data connectionData, isSecretKey bool) (vpnData map[string]string, ok bool) {
	vpnSectionData, ok := data[sectionVpn]
	if !ok {
		return
	}
	var variantValue dbus.Variant
	if isSecretKey {
		variantValue, ok = vpnSectionData[NM_SETTING_VPN_SECRETS]
		if !ok {
			return
		}
	} else {
		variantValue, ok = vpnSectionData[NM_SETTING_VPN_DATA]
		if !ok {
			return
		}
	}
	vpnData = interfaceToDictStringString(variantValue.Value())
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
