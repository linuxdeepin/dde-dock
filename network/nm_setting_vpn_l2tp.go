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
	"fmt"
	"pkg.deepin.io/dde/daemon/network/nm"
	. "pkg.deepin.io/lib/gettext"
)

const (
	nmVpnL2tpNameFile = VPN_NAME_FILES_DIR + "nm-l2tp-service.name"
)

var availableValuesNmL2tpSecretFlags []kvalue

func initAvailableValuesNmL2tpSecretFlags() {
	availableValuesNmL2tpSecretFlags = []kvalue{
		kvalue{nm.NM_L2TP_SECRET_FLAG_NONE, Tr("Saved")},
		kvalue{nm.NM_L2TP_SECRET_FLAG_NOT_SAVED, Tr("Always Ask")},
		kvalue{nm.NM_L2TP_SECRET_FLAG_NOT_REQUIRED, Tr("Not Required")},
	}
}

func isVpnL2tpRequireSecret(flag uint32) bool {
	if flag == nm.NM_L2TP_SECRET_FLAG_NONE || flag == nm.NM_L2TP_SECRET_FLAG_AGENT_OWNED {
		return true
	}
	return false
}

func isVpnL2tpNeedShowPassword(data connectionData) bool {
	return isVpnL2tpRequireSecret(getSettingVpnL2tpKeyPasswordFlags(data))
}

// new connection data
func newVpnL2tpConnectionData(id, uuid string) (data connectionData) {
	data = newBasicVpnConnectionData(id, uuid)
	initSettingSectionVpnL2tp(data)
	return
}

func initSettingSectionVpnL2tp(data connectionData) {
	initBasicSettingSectionVpn(data, nm.NM_DBUS_SERVICE_L2TP)
	setSettingVpnL2tpKeyPasswordFlags(data, nm.NM_L2TP_SECRET_FLAG_NONE)
	logicSetSettingVkVpnL2tpRequireMppe(data, true)
}

// vpn-l2tp
func getSettingVpnL2tpAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_GATEWAY)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_USER)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_PASSWORD_FLAGS)
	if isVpnL2tpNeedShowPassword(data) {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_PASSWORD)
	}
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_DOMAIN)
	return
}
func getSettingVpnL2tpAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_VPN_L2TP_KEY_PASSWORD_FLAGS:
		values = availableValuesNmL2tpSecretFlags
	}
	return
}
func checkSettingVpnL2tpValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	ensureSettingVpnL2tpKeyGatewayNoEmpty(data, errs)
	if isVpnL2tpNeedShowPassword(data) {
		ensureSettingVpnL2tpKeyPasswordNoEmpty(data, errs)
	}
	return
}

// vpn-l2tp-ppp
func getSettingVpnL2tpPppAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_REFUSE_EAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_REFUSE_PAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_REFUSE_CHAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_REFUSE_MSCHAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_REFUSE_MSCHAPV2)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE)
	if getSettingVkVpnL2tpRequireMppe(data) {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE_40)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE_128)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_MPPE_STATEFUL)
	}
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_NOBSDCOMP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_NODEFLATE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_NO_VJ_COMP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_NO_PCOMP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_NO_ACCOMP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_LCP_ECHO_FAILURE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_LCP_ECHO_INTERVAL)
	return
}
func getSettingVpnL2tpPppAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnL2tpPppValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	return
}

// vpn-l2tp-ipsec
func getSettingVpnL2tpIpsecAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_IPSEC_ENABLE)
	if getSettingVpnL2tpKeyIpsecEnable(data) {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_IPSEC_GATEWAY_ID)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_IPSEC_GROUP_NAME)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_IPSEC_PSK)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_IPSEC_IKE)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME, nm.NM_SETTING_VPN_L2TP_KEY_IPSEC_ESP)
	}
	return
}
func getSettingVpnL2tpIpsecAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnL2tpIpsecValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	return
}
func logicSetSettingVpnL2tpKeyIpsecEnable(data connectionData, value bool) (err error) {
	if !value {
		removeSettingVpnL2tpKeyIpsecGatewayId(data)
		removeSettingVpnL2tpKeyIpsecGroupName(data)
		removeSettingVpnL2tpKeyIpsecPsk(data)
	}
	setSettingVpnL2tpKeyIpsecEnable(data, value)
	return
}

// Virtual key
func getSettingVkVpnL2tpRequireMppe(data connectionData) (value bool) {
	if getSettingVpnL2tpKeyRequireMppe(data) ||
		getSettingVpnL2tpKeyRequireMppe128(data) ||
		getSettingVpnL2tpKeyRequireMppe40(data) {
		value = true
	}
	return
}
func logicSetSettingVkVpnL2tpRequireMppe(data connectionData, value bool) (err error) {
	if !value {
		// if disable mppe, remove related keys
		removeSettingVpnL2tpKeyRequireMppe40(data)
		removeSettingVpnL2tpKeyRequireMppe128(data)
		removeSettingVpnL2tpKeyMppeStateful(data)
	}
	setSettingVpnL2tpKeyRequireMppe(data, value)
	return
}

func getSettingVkVpnL2tpMppeSecurity(data connectionData) (value string) {
	if getSettingVpnL2tpKeyRequireMppe128(data) {
		value = "128-bit"
	} else if getSettingVpnL2tpKeyRequireMppe40(data) {
		value = "40-bit"
	} else if getSettingVpnL2tpKeyRequireMppe(data) {
		value = "default"
	} else {
		logger.Warning("get l2tp mppe security failed")
	}
	return
}
func logicSetSettingVkVpnL2tpMppeSecurity(data connectionData, value string) (err error) {
	if !getSettingVpnL2tpKeyRequireMppe(data) {
		err = fmt.Errorf(nmKeyErrorMissingDependsKey, nm.NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE)
		return
	}
	switch value {
	case "default":
		removeSettingVpnL2tpKeyRequireMppe40(data)
		removeSettingVpnL2tpKeyRequireMppe128(data)
	case "128-bit":
		removeSettingVpnL2tpKeyRequireMppe40(data)
		setSettingVpnL2tpKeyRequireMppe128(data, true)
	case "40-bit":
		setSettingVpnL2tpKeyRequireMppe40(data, true)
		removeSettingVpnL2tpKeyRequireMppe128(data)
	}
	return
}

func getSettingVkVpnL2tpEnableLcpEcho(data connectionData) (value bool) {
	if isSettingVpnL2tpKeyLcpEchoFailureExists(data) && isSettingVpnL2tpKeyLcpEchoIntervalExists(data) {
		return true
	}
	return false
}
func logicSetSettingVkVpnL2tpEnableLcpEcho(data connectionData, value bool) (err error) {
	if value {
		setSettingVpnL2tpKeyLcpEchoFailure(data, 5)
		setSettingVpnL2tpKeyLcpEchoInterval(data, 30)
	} else {
		removeSettingVpnL2tpKeyLcpEchoFailure(data)
		removeSettingVpnL2tpKeyLcpEchoInterval(data)
	}
	return
}
