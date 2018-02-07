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
	nmVpnPptpNameFile = VPN_NAME_FILES_DIR + "nm-pptp-service.name"
)

var availableValuesNmPptpSecretFlags []kvalue

func initAvailableValuesNmPptpSecretFlags() {
	availableValuesNmPptpSecretFlags = []kvalue{
		kvalue{nm.NM_PPTP_SECRET_FLAG_NONE, Tr("Saved")}, // system saved
		kvalue{nm.NM_PPTP_SECRET_FLAG_NOT_SAVED, Tr("Always Ask")},
		kvalue{nm.NM_PPTP_SECRET_FLAG_NOT_REQUIRED, Tr("Not Required")},
	}
}

func isVpnPptpRequireSecret(flag uint32) bool {
	if flag == nm.NM_PPTP_SECRET_FLAG_NONE || flag == nm.NM_PPTP_SECRET_FLAG_AGENT_OWNED {
		return true
	}
	return false
}

func isVpnPptpNeedShowPassword(data connectionData) bool {
	return isVpnPptpRequireSecret(getSettingVpnPptpKeyPasswordFlags(data))
}

// new connection data
func newVpnPptpConnectionData(id, uuid string) (data connectionData) {
	data = newBasicVpnConnectionData(id, uuid)
	initSettingSectionVpnPptp(data)
	return
}

func initSettingSectionVpnPptp(data connectionData) {
	initBasicSettingSectionVpn(data, nm.NM_DBUS_SERVICE_PPTP)
	setSettingVpnPptpKeyPasswordFlags(data, nm.NM_PPTP_SECRET_FLAG_NONE)
	logicSetSettingVkVpnPptpRequireMppe(data, true)
}

// vpn-pptp
func getSettingVpnPptpAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_GATEWAY)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_USER)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_PASSWORD_FLAGS)
	if isVpnPptpNeedShowPassword(data) {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_PASSWORD)
	}
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_DOMAIN)
	return
}
func getSettingVpnPptpAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_VPN_PPTP_KEY_PASSWORD_FLAGS:
		values = availableValuesNmPptpSecretFlags
	}
	return
}
func checkSettingVpnPptpValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	ensureSettingVpnPptpKeyGatewayNoEmpty(data, errs)
	if isVpnPptpNeedShowPassword(data) {
		ensureSettingVpnPptpKeyPasswordNoEmpty(data, errs)
	}
	return
}

// vpn-pptp-ppp
func getSettingVpnPptpPppAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_REFUSE_EAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_REFUSE_PAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_REFUSE_CHAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_REFUSE_MSCHAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_REFUSE_MSCHAPV2)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE)
	if getSettingVkVpnPptpRequireMppe(data) {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE_40)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE_128)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_MPPE_STATEFUL)
	}
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_NOBSDCOMP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_NODEFLATE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_NO_VJ_COMP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_LCP_ECHO_FAILURE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME, nm.NM_SETTING_VPN_PPTP_KEY_LCP_ECHO_INTERVAL)
	return
}
func getSettingVpnPptpPppAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnPptpPppValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	return
}

// Virtual key
func getSettingVkVpnPptpRequireMppe(data connectionData) (value bool) {
	if getSettingVpnPptpKeyRequireMppe(data) ||
		getSettingVpnPptpKeyRequireMppe128(data) ||
		getSettingVpnPptpKeyRequireMppe40(data) {
		value = true
	}
	return
}
func logicSetSettingVkVpnPptpRequireMppe(data connectionData, value bool) (err error) {
	if !value {
		// if disable mppe, remove related keys
		removeSettingVpnPptpKeyRequireMppe40(data)
		removeSettingVpnPptpKeyRequireMppe128(data)
		removeSettingVpnPptpKeyMppeStateful(data)
	}
	setSettingVpnPptpKeyRequireMppe(data, value)
	return
}

func getSettingVkVpnPptpMppeSecurity(data connectionData) (value string) {
	if getSettingVpnPptpKeyRequireMppe128(data) {
		value = "128-bit"
	} else if getSettingVpnPptpKeyRequireMppe40(data) {
		value = "40-bit"
	} else if getSettingVpnPptpKeyRequireMppe(data) {
		value = "default"
	} else {
		logger.Warning("get pptp mppe security failed")
	}
	return
}
func logicSetSettingVkVpnPptpMppeSecurity(data connectionData, value string) (err error) {
	if !getSettingVpnPptpKeyRequireMppe(data) {
		err = fmt.Errorf(nmKeyErrorMissingDependsKey, nm.NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE)
		return
	}
	switch value {
	case "default":
		removeSettingVpnPptpKeyRequireMppe40(data)
		removeSettingVpnPptpKeyRequireMppe128(data)
	case "128-bit":
		removeSettingVpnPptpKeyRequireMppe40(data)
		setSettingVpnPptpKeyRequireMppe128(data, true)
	case "40-bit":
		setSettingVpnPptpKeyRequireMppe40(data, true)
		removeSettingVpnPptpKeyRequireMppe128(data)
	}
	return
}

func getSettingVkVpnPptpEnableLcpEcho(data connectionData) (value bool) {
	if isSettingVpnPptpKeyLcpEchoFailureExists(data) && isSettingVpnPptpKeyLcpEchoIntervalExists(data) {
		return true
	}
	return false
}
func logicSetSettingVkVpnPptpEnableLcpEcho(data connectionData, value bool) (err error) {
	if value {
		setSettingVpnPptpKeyLcpEchoFailure(data, 5)
		setSettingVpnPptpKeyLcpEchoInterval(data, 30)
	} else {
		removeSettingVpnPptpKeyLcpEchoFailure(data)
		removeSettingVpnPptpKeyLcpEchoInterval(data)
	}
	return
}
