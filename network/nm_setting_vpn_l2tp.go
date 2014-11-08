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
	"fmt"
	. "pkg.linuxdeepin.com/lib/gettext"
)

// For the NM <-> VPN plugin service
const (
	NM_DBUS_SERVICE_L2TP   = "org.freedesktop.NetworkManager.l2tp"
	NM_DBUS_INTERFACE_L2TP = "org.freedesktop.NetworkManager.l2tp"
	NM_DBUS_PATH_L2TP      = "/org/freedesktop/NetworkManager/l2tp"
)

const (
	NM_DBUS_SERVICE_L2TP_PPP   = "org.freedesktop.NetworkManager.l2tp-ppp"
	NM_DBUS_PATH_L2TP_PPP      = "/org/freedesktop/NetworkManager/l2tp/ppp"
	NM_DBUS_INTERFACE_L2TP_PPP = "org.freedesktop.NetworkManager.l2tp.ppp"
)

const (
	nmVpnL2tpNameFile = VPN_NAME_FILES_DIR + "nm-l2tp-service.name"
)

const (
	NM_SETTING_VPN_L2TP_KEY_GATEWAY           = "gateway"
	NM_SETTING_VPN_L2TP_KEY_USER              = "user"
	NM_SETTING_VPN_L2TP_KEY_PASSWORD          = "password"
	NM_SETTING_VPN_L2TP_KEY_PASSWORD_FLAGS    = "password-flags"
	NM_SETTING_VPN_L2TP_KEY_DOMAIN            = "domain"
	NM_SETTING_VPN_L2TP_KEY_REFUSE_EAP        = "refuse-eap"
	NM_SETTING_VPN_L2TP_KEY_REFUSE_PAP        = "refuse-pap"
	NM_SETTING_VPN_L2TP_KEY_REFUSE_CHAP       = "refuse-chap"
	NM_SETTING_VPN_L2TP_KEY_REFUSE_MSCHAP     = "refuse-mschap"
	NM_SETTING_VPN_L2TP_KEY_REFUSE_MSCHAPV2   = "refuse-mschapv2"
	NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE      = "require-mppe"
	NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE_40   = "require-mppe-40"
	NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE_128  = "require-mppe-128"
	NM_SETTING_VPN_L2TP_KEY_MPPE_STATEFUL     = "mppe-stateful"
	NM_SETTING_VPN_L2TP_KEY_NOBSDCOMP         = "nobsdcomp"
	NM_SETTING_VPN_L2TP_KEY_NODEFLATE         = "nodeflate"
	NM_SETTING_VPN_L2TP_KEY_NO_VJ_COMP        = "no-vj-comp"
	NM_SETTING_VPN_L2TP_KEY_NO_PCOMP          = "nopcomp"
	NM_SETTING_VPN_L2TP_KEY_NO_ACCOMP         = "noaccomp"
	NM_SETTING_VPN_L2TP_KEY_LCP_ECHO_FAILURE  = "lcp-echo-failure"
	NM_SETTING_VPN_L2TP_KEY_LCP_ECHO_INTERVAL = "lcp-echo-interval"
	NM_SETTING_VPN_L2TP_KEY_IPSEC_ENABLE      = "ipsec-enabled"
	NM_SETTING_VPN_L2TP_KEY_IPSEC_GATEWAY_ID  = "ipsec-gateway-id"
	NM_SETTING_VPN_L2TP_KEY_IPSEC_GROUP_NAME  = "ipsec-group-name"
	NM_SETTING_VPN_L2TP_KEY_IPSEC_PSK         = "ipsec-psk"
)

// vpn key descriptions
// sta_VPNtic ValidProperty valid_properties[] = {
// 	{ NM_L2TP_KEY_GATEWAY,           G_TYPE_STRING, TRUE },
// 	{ NM_L2TP_KEY_USER,              G_TYPE_STRING, FALSE },
// 	{ NM_L2TP_KEY_DOMAIN,            G_TYPE_STRING, FALSE },
// 	{ NM_L2TP_KEY_REFUSE_EAP,        G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_REFUSE_PAP,        G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_REFUSE_CHAP,       G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_REFUSE_MSCHAP,     G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_REFUSE_MSCHAPV2,   G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_REQUIRE_MPPE,      G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_REQUIRE_MPPE_40,   G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_REQUIRE_MPPE_128,  G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_MPPE_STATEFUL,     G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_NOBSDCOMP,         G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_NODEFLATE,         G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_NO_VJ_COMP,        G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_NO_PCOMP,          G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_LCP_ECHO_FAILURE,  G_TYPE_UINT, FALSE },
// 	{ NM_L2TP_KEY_LCP_ECHO_INTERVAL, G_TYPE_UINT, FALSE },
// 	{ NM_L2TP_KEY_PASSWORD"-flags",  G_TYPE_UINT, FALSE },
// 	{ NM_L2TP_KEY_IPSEC_ENABLE,      G_TYPE_BOOLEAN, FALSE },
// 	{ NM_L2TP_KEY_IPSEC_GATEWAY_ID,  G_TYPE_STRING, FALSE },
// 	{ NM_L2TP_KEY_IPSEC_GROUP_NAME,  G_TYPE_STRING, FALSE },
// 	{ NM_L2TP_KEY_IPSEC_PSK,         G_TYPE_STRING, FALSE },
// 	{ NULL,                          G_TYPE_NONE,   FALSE }
// }

// Define secret flags
const (
	NM_L2TP_SECRET_FLAG_NONE         = 0 // system saved
	NM_L2TP_SECRET_FLAG_AGENT_OWNED  = 1
	NM_L2TP_SECRET_FLAG_NOT_SAVED    = 3
	NM_L2TP_SECRET_FLAG_NOT_REQUIRED = 5
)

var availableValuesNmL2tpSecretFlags []kvalue

func initAvailableValuesNmL2tpSecretFlags() {
	availableValuesNmL2tpSecretFlags = []kvalue{
		kvalue{NM_L2TP_SECRET_FLAG_NONE, Tr("Saved")},
		kvalue{NM_L2TP_SECRET_FLAG_NOT_SAVED, Tr("Always Ask")},
		kvalue{NM_L2TP_SECRET_FLAG_NOT_REQUIRED, Tr("Not Required")},
	}
}

func isVpnL2tpRequireSecret(flag uint32) bool {
	if flag == NM_L2TP_SECRET_FLAG_NONE || flag == NM_L2TP_SECRET_FLAG_AGENT_OWNED {
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
	initBasicSettingSectionVpn(data, NM_DBUS_SERVICE_L2TP)
	setSettingVpnL2tpKeyPasswordFlags(data, NM_L2TP_SECRET_FLAG_NONE)
	logicSetSettingVkVpnL2tpRequireMppe(data, true)
}

// vpn-l2tp
func getSettingVpnL2tpAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionVpnL2tp, NM_SETTING_VPN_L2TP_KEY_GATEWAY)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tp, NM_SETTING_VPN_L2TP_KEY_USER)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tp, NM_SETTING_VPN_L2TP_KEY_PASSWORD_FLAGS)
	if isVpnL2tpNeedShowPassword(data) {
		keys = appendAvailableKeys(data, keys, sectionVpnL2tp, NM_SETTING_VPN_L2TP_KEY_PASSWORD)
	}
	keys = appendAvailableKeys(data, keys, sectionVpnL2tp, NM_SETTING_VPN_L2TP_KEY_DOMAIN)
	return
}
func getSettingVpnL2tpAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_L2TP_KEY_PASSWORD_FLAGS:
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
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REFUSE_EAP)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REFUSE_PAP)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REFUSE_CHAP)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REFUSE_MSCHAP)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REFUSE_MSCHAPV2)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE)
	if getSettingVkVpnL2tpRequireMppe(data) {
		keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE_40)
		keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE_128)
		keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_MPPE_STATEFUL)
	}
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_NOBSDCOMP)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_NODEFLATE)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_NO_VJ_COMP)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_NO_PCOMP)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_NO_ACCOMP)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_LCP_ECHO_FAILURE)
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_LCP_ECHO_INTERVAL)
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
	keys = appendAvailableKeys(data, keys, sectionVpnL2tpIpsec, NM_SETTING_VPN_L2TP_KEY_IPSEC_ENABLE)
	if getSettingVpnL2tpKeyIpsecEnable(data) {
		keys = appendAvailableKeys(data, keys, sectionVpnL2tpIpsec, NM_SETTING_VPN_L2TP_KEY_IPSEC_GATEWAY_ID)
		keys = appendAvailableKeys(data, keys, sectionVpnL2tpIpsec, NM_SETTING_VPN_L2TP_KEY_IPSEC_GROUP_NAME)
		keys = appendAvailableKeys(data, keys, sectionVpnL2tpIpsec, NM_SETTING_VPN_L2TP_KEY_IPSEC_PSK)
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
	if value {
		// if require mppe, refuse some authentications
		setSettingVpnL2tpKeyRefuseChap(data, true)
		setSettingVpnL2tpKeyRefuseEap(data, true)
		setSettingVpnL2tpKeyRefusePap(data, true)
	} else {
		// if disable mppe, remove related keys
		removeSettingVpnL2tpKeyRequireMppe40(data)
		removeSettingVpnL2tpKeyRequireMppe128(data)
		removeSettingVpnL2tpKeyMppeStateful(data)

		removeSettingVpnL2tpKeyRefuseChap(data)
		removeSettingVpnL2tpKeyRefuseEap(data)
		removeSettingVpnL2tpKeyRefusePap(data)
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
		err = fmt.Errorf(NM_KEY_ERROR_MISSING_DEPENDS_KEY, NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE)
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
