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

import (
	"pkg.deepin.io/dde/daemon/network/nm"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/utils"
	"strconv"
	"strings"
)

// Custom device types, use sting instead of number, used by front-end
const (
	deviceUnknown    = "unknown"
	deviceEthernet   = "wired"
	deviceWifi       = "wireless"
	deviceUnused1    = "unused1"
	deviceUnused2    = "unused2"
	deviceBt         = "bt"
	deviceOlpcMesh   = "olpc-mesh"
	deviceWimax      = "wimax"
	deviceModem      = "modem"
	deviceInfiniband = "infiniband"
	deviceBond       = "bond"
	deviceVlan       = "vlan"
	deviceAdsl       = "adsl"
	deviceBridge     = "bridge"
	deviceGeneric    = "generic"
	deviceTeam       = "team"
)

func getCustomDeviceType(devType uint32) (customDevType string) {
	switch devType {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		return deviceEthernet
	case nm.NM_DEVICE_TYPE_WIFI:
		return deviceWifi
	case nm.NM_DEVICE_TYPE_UNUSED1:
		return deviceUnused1
	case nm.NM_DEVICE_TYPE_UNUSED2:
		return deviceUnused2
	case nm.NM_DEVICE_TYPE_BT:
		return deviceBt
	case nm.NM_DEVICE_TYPE_OLPC_MESH:
		return deviceOlpcMesh
	case nm.NM_DEVICE_TYPE_WIMAX:
		return deviceWimax
	case nm.NM_DEVICE_TYPE_MODEM:
		return deviceModem
	case nm.NM_DEVICE_TYPE_INFINIBAND:
		return deviceInfiniband
	case nm.NM_DEVICE_TYPE_BOND:
		return deviceBond
	case nm.NM_DEVICE_TYPE_VLAN:
		return deviceVlan
	case nm.NM_DEVICE_TYPE_ADSL:
		return deviceAdsl
	case nm.NM_DEVICE_TYPE_BRIDGE:
		return deviceBridge
	case nm.NM_DEVICE_TYPE_GENERIC:
		return deviceGeneric
	case nm.NM_DEVICE_TYPE_TEAM:
		return deviceTeam
	case nm.NM_DEVICE_TYPE_UNKNOWN:
	default:
		logger.Error("unknown device type", devType)
	}
	return deviceUnknown
}

// Custom connection types
const (
	connectionUnknown         = "unknown"
	connectionWired           = "wired"
	connectionWireless        = "wireless"
	connectionWirelessAdhoc   = "wireless-adhoc"
	connectionWirelessHotspot = "wireless-hotspot"
	connectionPppoe           = "pppoe"
	connectionMobileGsm       = "mobile-gsm"
	connectionMobileCdma      = "mobile-cdma"
	connectionVpnL2tp         = "vpn-l2tp"
	connectionVpnOpenconnect  = "vpn-openconnect"
	connectionVpnOpenvpn      = "vpn-openvpn"
	connectionVpnStrongswan   = "vpn-strongswan"
	connectionVpnPptp         = "vpn-pptp"
	connectionVpnVpnc         = "vpn-vpnc"
)

// wrapper for custom connection types
const (
	connectionMobile = "mobile" // wrapper for gsm and cdma
	connectionVpn    = "vpn"    // wrapper for all vpn types
)

// key-map values for internationalization
type connectionType struct {
	Value, Text string
}

var supportedConnectionTypes = []string{
	connectionWired,
	connectionWireless,
	connectionWirelessAdhoc,
	connectionWirelessHotspot,
	connectionPppoe,
	connectionMobile,
	connectionMobileGsm,
	connectionMobileCdma,
	connectionVpn,
	connectionVpnL2tp,
	connectionVpnOpenconnect,
	connectionVpnOpenvpn,
	connectionVpnPptp,
	connectionVpnStrongswan,
	connectionVpnVpnc,
}

func getCustomConnectionTypeForUuid(uuid string) (connType string) {
	connType = connectionUnknown
	cpath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	return getCustomConnectionType(data)
}

// return custom connection type, and the wrapper types will be ignored, e.g. connectionMobile.
func getCustomConnectionType(data connectionData) (connType string) {
	t := getSettingConnectionType(data)
	switch t {
	case nm.NM_SETTING_WIRED_SETTING_NAME:
		connType = connectionWired
	case nm.NM_SETTING_WIRELESS_SETTING_NAME:
		if isSettingWirelessModeExists(data) {
			switch getSettingWirelessMode(data) {
			case nm.NM_SETTING_WIRELESS_MODE_INFRA:
				connType = connectionWireless
			case nm.NM_SETTING_WIRELESS_MODE_ADHOC:
				connType = connectionWirelessAdhoc
			case nm.NM_SETTING_WIRELESS_MODE_AP:
				connType = connectionWirelessHotspot
			}
		} else {
			connType = connectionWireless
		}
	case nm.NM_SETTING_PPPOE_SETTING_NAME:
		connType = connectionPppoe
	case nm.NM_SETTING_GSM_SETTING_NAME:
		connType = connectionMobileGsm
	case nm.NM_SETTING_CDMA_SETTING_NAME:
		connType = connectionMobileCdma
	case nm.NM_SETTING_VPN_SETTING_NAME:
		switch getSettingVpnServiceType(data) {
		case nm.NM_DBUS_SERVICE_L2TP:
			connType = connectionVpnL2tp
		case nm.NM_DBUS_SERVICE_OPENCONNECT:
			connType = connectionVpnOpenconnect
		case nm.NM_DBUS_SERVICE_OPENVPN:
			connType = connectionVpnOpenvpn
		case nm.NM_DBUS_SERVICE_PPTP:
			connType = connectionVpnPptp
		case nm.NM_DBUS_SERVICE_STRONGSWAN:
			connType = connectionVpnStrongswan
		case nm.NM_DBUS_SERVICE_VPNC:
			connType = connectionVpnVpnc
		}
	}
	if len(connType) == 0 {
		connType = connectionUnknown
	}
	return
}

func isWirelessConnection(data connectionData) (isWireless bool) {
	if getSettingConnectionType(data) == nm.NM_SETTING_WIRELESS_SETTING_NAME {
		return true
	}
	return false
}

func isVpnConnection(data connectionData) (isVpn bool) {
	if getSettingConnectionType(data) == nm.NM_SETTING_VPN_SETTING_NAME {
		return true
	}
	return false
}

func isCreatedManuallyConnection(data connectionData) (isCreateManual bool) {
	if isVpnConnection(data) {
		return true
	}
	switch getCustomConnectionType(data) {
	case connectionPppoe:
		return true
	}
	return false
}

// generate connection id when creating a new connection
func genConnectionId(connType string) (id string) {
	var idPrefix string
	switch connType {
	default:
		idPrefix = Tr("Connection")
	case connectionWired:
		idPrefix = Tr("Wired Connection")
	case connectionWireless:
		idPrefix = Tr("Wireless Connection")
	case connectionWirelessAdhoc:
		idPrefix = Tr("Wireless Ad-Hoc")
	case connectionWirelessHotspot:
		idPrefix = Tr("Wireless Ap-Hotspot")
	case connectionPppoe:
		idPrefix = Tr("PPPoE Connection")
	case connectionMobile:
		idPrefix = Tr("Mobile Connection")
	case connectionMobileGsm:
		idPrefix = Tr("Mobile GSM Connection")
	case connectionMobileCdma:
		idPrefix = Tr("Mobile CDMA Connection")
	case connectionVpn:
		idPrefix = Tr("VPN Connection")
	case connectionVpnL2tp:
		idPrefix = Tr("VPN L2TP")
	case connectionVpnOpenconnect:
		idPrefix = Tr("VPN OpenConnect")
	case connectionVpnOpenvpn:
		idPrefix = Tr("VPN OpenVPN")
	case connectionVpnPptp:
		idPrefix = Tr("VPN PPTP")
	case connectionVpnStrongswan:
		idPrefix = Tr("VPN StrongSwan")
	case connectionVpnVpnc:
		idPrefix = Tr("VPN VPNC")
	}
	allIds := nmGetConnectionIds()
	for i := 1; ; i++ {
		id = idPrefix + " " + strconv.Itoa(i)
		if !isStringInArray(id, allIds) {
			break
		}
	}
	return
}

func isConnectionAlwaysAsk(data connectionData, settingName string) (ask bool) {
	sectionData, ok := data[settingName]
	if !ok {
		ask = false
		return
	}
	// query all secret key flags that should be suffixed with
	// "-flags" and check if is marked as ask user always
	for key, variant := range sectionData {
		if strings.HasSuffix(key, "-flags") {
			value := variant.Value()
			if flag, ok := value.(uint32); ok {
				if flag == nm.NM_SETTING_SECRET_FLAG_NONE {
					ask = true
				}
			}
		}
	}
	return
}

const (
	nmKeyErrorMissingValue       = "missing value"
	nmKeyErrorEmptyValue         = "value is empty"
	nmKeyErrorInvalidValue       = "invalid value"
	nmKeyErrorIp4MethodConflict  = `%s cannot be used with the 'shared', 'link-local', or 'disabled' methods`
	nmKeyErrorIp4AddressesStruct = "echo IPv4 address structure is composed of 3 32-bit values, address, prefix and gateway"
	// nm.NM_KEY_ERROR_IP4_ADDRESSES_PREFIX = "IPv4 prefix's value should be 1-32"
	nmKeyErrorIp6MethodConflict     = `%s cannot be used with the 'shared', 'link-local', or 'ignore' methods`
	nmKeyErrorMissingSection        = "missing section %s"
	nmKeyErrorEmptySection          = "section %s is empty"
	nmKeyErrorMissingDependsKey     = "missing depends key %s"
	nmKeyErrorMissingDependsPackage = "missing depends package %s"
)

func rememberError(errs sectionErrors, section, key, errMsg string) {
	relatedVks := getRelatedVkeys(section, key)
	if len(relatedVks) > 0 {
		rememberVkError(errs, section, key, errMsg)
		return
	}
	doRememberError(errs, key, errMsg)
}

func rememberVkError(errs sectionErrors, section, key, errMsg string) {
	vks := getRelatedVkeys(section, key)
	for _, vk := range vks {
		if !isOptionalVkey(section, vk) {
			doRememberError(errs, vk, errMsg)
		}
	}
}

func doRememberError(errs sectionErrors, key, errMsg string) {
	if _, ok := errs[key]; ok {
		// error already exists for this key
		return
	}
	errs[key] = errMsg
}

// start with "file://", end with null byte
func ensureByteArrayUriPathExistsFor8021x(errs sectionErrors, section, key string, bytePath []byte) {
	path := byteArrayToStrPath(bytePath)
	if !utils.IsURI(path) {
		rememberError(errs, section, key, nmKeyErrorInvalidValue)
		return
	}
	ensureFileExists(errs, section, key, toLocalPathFor8021x(path))
}

func ensureFileExists(errs sectionErrors, section, key, file string) {
	file = toLocalPath(file)
	if !utils.IsFileExist(file) {
		rememberError(errs, section, key, nmKeyErrorInvalidValue)
	}
}

// Custom device types, use sting instead of number, used by front-end
const (
	passTypeGeneral           = "general"
	passTypeWifiWepKey        = "wifi-wep-key"
	passTypeWifiWepPassphrase = "wifi-wep-passphrase"
	passTypeWifiWpaPsk        = "wifi-wpa-psk"
)

func getSettingPassKeyType(data connectionData, settingName string) (passType string) {
	switch settingName {
	case nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME:
		switch getSettingVkWirelessSecurityKeyMgmt(data) {
		case "wep":
			wepKeyType := getSettingWirelessSecurityWepKeyType(data)
			if wepKeyType == nm.NM_WEP_KEY_TYPE_KEY {
				passType = passTypeWifiWepKey
			} else if wepKeyType == nm.NM_WEP_KEY_TYPE_PASSPHRASE {
				passType = passTypeWifiWepPassphrase
			}
		case "wpa-psk":
			passType = passTypeWifiWpaPsk
		}
	default:
		passType = passTypeGeneral
	}
	if len(passType) == 0 {
		logger.Errorf("Unknown password key type for data[%s]", settingName)
	}
	return
}

func isPasswordValid(passType, value string) (ok bool) {
	switch passType {
	case passTypeWifiWepKey:
		// If wep key type set to 1 and the keys are hexadecimal, they
		// must be either 10 or 26 characters in length. If set to 1
		// and the keys are ASCII keys, they must be either 5 or 13
		// characters in length.
		if len(value) != 10 && len(value) != 26 && len(value) != 5 && len(value) != 13 {
			ok = false
		} else {
			ok = true
		}
	case passTypeWifiWepPassphrase:
		// If wep key type set to 2, the passphrase is hashed using
		// the de-facto MD5 method to derive the actual WEP key.
		if len(value) == 0 {
			ok = false
		} else {
			ok = true
		}
	case passTypeWifiWpaPsk:
		// If the wpa-psk key is 64-characters long, it must contain
		// only hexadecimal characters and is interpreted as a
		// hexadecimal WPA key. Otherwise, the key must be between 8
		// and 63 ASCII characters (as specified in the 802.11i
		// standard) and is interpreted as a WPA passphrase
		if len(value) < 8 || len(value) > 64 {
			ok = false
		} else {
			ok = true
		}
	case passTypeGeneral:
		ok = true
	default:
		ok = true
	}
	return
}
