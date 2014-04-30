package main

import (
	"dlib"
)

const (
	NM_DBUS_SERVICE_PPTP   = "org.freedesktop.NetworkManager.pptp"
	NM_DBUS_INTERFACE_PPTP = "org.freedesktop.NetworkManager.pptp"
	NM_DBUS_PATH_PPTP      = "/org/freedesktop/NetworkManager/pptp"
)

const (
	NM_SETTING_VPN_PPTP_KEY_GATEWAY           = "gateway"
	NM_SETTING_VPN_PPTP_KEY_USER              = "user"
	NM_SETTING_VPN_PPTP_KEY_PASSWORD          = "password"
	NM_SETTING_VPN_PPTP_KEY_PASSWORD_FLAGS    = "password-flags"
	NM_SETTING_VPN_PPTP_KEY_DOMAIN            = "domain"
	NM_SETTING_VPN_PPTP_KEY_REFUSE_EAP        = "refuse-eap"
	NM_SETTING_VPN_PPTP_KEY_REFUSE_PAP        = "refuse-pap"
	NM_SETTING_VPN_PPTP_KEY_REFUSE_CHAP       = "refuse-chap"
	NM_SETTING_VPN_PPTP_KEY_REFUSE_MSCHAP     = "refuse-mschap"
	NM_SETTING_VPN_PPTP_KEY_REFUSE_MSCHAPV2   = "refuse-mschapv2"
	NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE      = "require-mppe"
	NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE_40   = "require-mppe-40"
	NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE_128  = "require-mppe-128"
	NM_SETTING_VPN_PPTP_KEY_MPPE_STATEFUL     = "mppe-stateful"
	NM_SETTING_VPN_PPTP_KEY_NOBSDCOMP         = "nobsdcomp"
	NM_SETTING_VPN_PPTP_KEY_NODEFLATE         = "nodeflate"
	NM_SETTING_VPN_PPTP_KEY_NO_VJ_COMP        = "no-vj-comp"
	NM_SETTING_VPN_PPTP_KEY_LCP_ECHO_FAILURE  = "lcp-echo-failure"
	NM_SETTING_VPN_PPTP_KEY_LCP_ECHO_INTERVAL = "lcp-echo-interval"
)

// vpn key descriptions
// static ValidProperty valid_properties[] = {
// 	{ NM_PPTP_KEY_GATEWAY,           G_TYPE_STRING, TRUE },
// 	{ NM_PPTP_KEY_USER,              G_TYPE_STRING, FALSE },
// 	{ NM_PPTP_KEY_DOMAIN,            G_TYPE_STRING, FALSE },
// 	{ NM_PPTP_KEY_REFUSE_EAP,        G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_REFUSE_PAP,        G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_REFUSE_CHAP,       G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_REFUSE_MSCHAP,     G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_REFUSE_MSCHAPV2,   G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_REQUIRE_MPPE,      G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_REQUIRE_MPPE_40,   G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_REQUIRE_MPPE_128,  G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_MPPE_STATEFUL,     G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_NOBSDCOMP,         G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_NODEFLATE,         G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_NO_VJ_COMP,        G_TYPE_BOOLEAN, FALSE },
// 	{ NM_PPTP_KEY_LCP_ECHO_FAILURE,  G_TYPE_UINT, FALSE },
// 	{ NM_PPTP_KEY_LCP_ECHO_INTERVAL, G_TYPE_UINT, FALSE },
// 	{ NM_PPTP_KEY_PASSWORD"-flags",  G_TYPE_UINT, FALSE },
// 	{ NULL,                          G_TYPE_NONE, FALSE }
// }
// static ValidProperty valid_secrets[] = {
// 	{ NM_PPTP_KEY_PASSWORD,          G_TYPE_STRING, FALSE },
// 	{ NULL,                          G_TYPE_NONE,   FALSE }
// };
// static ValidProperty valid_secrets[] = {
// 	{ NM_L2TP_KEY_PASSWORD,          G_TYPE_STRING, FALSE },
// 	{ NULL,                          G_TYPE_NONE,   FALSE }
// };

// Define secret flags
const (
	NM_PPTP_SECRET_FLAG_NONE         = 0
	NM_PPTP_SECRET_FLAG_AGENT_OWNED  = 1
	NM_PPTP_SECRET_FLAG_NOT_SAVED    = 3
	NM_PPTP_SECRET_FLAG_NOT_REQUIRED = 5
)

var availableValuesNMPptpSecretFlag = []kvalue{
	kvalue{NM_PPTP_SECRET_FLAG_NONE, dlib.Tr("Saved")}, // system saved
	kvalue{NM_PPTP_SECRET_FLAG_NOT_SAVED, dlib.Tr("Always Ask")},
	kvalue{NM_PPTP_SECRET_FLAG_NOT_REQUIRED, dlib.Tr("Not Required")},
}

func isVpnPptpNeedShowPassword(data connectionData) bool {
	flag := getSettingVpnPptpKeyPasswordFlags(data)
	if flag == NM_PPTP_SECRET_FLAG_NONE || flag == NM_PPTP_SECRET_FLAG_AGENT_OWNED {
		return true
	}
	return false
}

// new connection data
func newVpnPptpConnectionData(id, uuid string) (data connectionData) {
	data = newBasicVpnConnectionData(id, uuid, NM_DBUS_SERVICE_PPTP)
	setSettingVpnPptpKeyPasswordFlags(data, NM_PPTP_SECRET_FLAG_NONE)
	return
}

// vpn-pptp
func getSettingVpnPptpAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, fieldVpnPptp, NM_SETTING_VPN_PPTP_KEY_GATEWAY)
	keys = appendAvailableKeys(data, keys, fieldVpnPptp, NM_SETTING_VPN_PPTP_KEY_USER)
	keys = appendAvailableKeys(data, keys, fieldVpnPptp, NM_SETTING_VPN_PPTP_KEY_PASSWORD_FLAGS)
	if isVpnPptpNeedShowPassword(data) {
		keys = appendAvailableKeys(data, keys, fieldVpnPptp, NM_SETTING_VPN_PPTP_KEY_PASSWORD)
	}
	keys = appendAvailableKeys(data, keys, fieldVpnPptp, NM_SETTING_VPN_PPTP_KEY_DOMAIN)
	return
}
func getSettingVpnPptpAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_PPTP_KEY_PASSWORD_FLAGS:
		values = availableValuesNMPptpSecretFlag
	}
	return
}
func checkSettingVpnPptpValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	ensureSettingVpnPptpKeyGatewayNoEmpty(data, errs)
	return
}

// vpn-pptp-ppp
func getSettingVpnPptpPppAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REFUSE_EAP)
	keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REFUSE_PAP)
	keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REFUSE_CHAP)
	keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REFUSE_MSCHAP)
	keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REFUSE_MSCHAPV2)
	keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE)
	if getSettingVpnPptpKeyRequireMppe(data) {
		keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE_40)
		keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE_128)
		keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_MPPE_STATEFUL)
	}
	keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_NOBSDCOMP)
	keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_NODEFLATE)
	keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_NO_VJ_COMP)
	keys = appendAvailableKeys(data, keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_LCP_ECHO_FAILURE)
	return
}
func getSettingVpnPptpPppAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnPptpPppValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	return
}
func logicSetSettingVpnPptpKeyRequireMppe(data connectionData, value bool) (err error) {
	if !value {
		removeSettingVpnPptpKeyRequireMppe40(data)
		removeSettingVpnPptpKeyRequireMppe128(data)
		removeSettingVpnPptpKeyMppeStateful(data)
	}
	setSettingVpnPptpKeyRequireMppe(data, value)
	return
}

// Virtual key
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
