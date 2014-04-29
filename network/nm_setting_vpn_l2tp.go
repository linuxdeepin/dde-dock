package main

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

func newVpnL2tpConnectionData(id, uuid string) (data connectionData) {
	data = newBasicVpnConnectionData(id, uuid, NM_DBUS_SERVICE_L2TP)
	setSettingVpnL2tpKeyPasswordFlags(data, 1)
	return
}

// vpn-l2tp
func getSettingVpnL2tpAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnL2tp, NM_SETTING_VPN_L2TP_KEY_GATEWAY)
	keys = appendAvailableKeys(keys, fieldVpnL2tp, NM_SETTING_VPN_L2TP_KEY_USER)
	keys = appendAvailableKeys(keys, fieldVpnL2tp, NM_SETTING_VPN_L2TP_KEY_PASSWORD)
	keys = appendAvailableKeys(keys, fieldVpnL2tp, NM_SETTING_VPN_L2TP_KEY_DOMAIN)
	return
}
func getSettingVpnL2tpAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnL2tpValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	ensureSettingVpnL2tpKeyGatewayNoEmpty(data, errs)
	// TODO
	return
}

// vpn-l2tp-ppp
func getSettingVpnL2tpPppAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REFUSE_EAP)
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REFUSE_PAP)
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REFUSE_CHAP)
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REFUSE_MSCHAP)
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REFUSE_MSCHAPV2)
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE)
	if getSettingVpnL2tpKeyRequireMppe(data) {
		keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE_40)
		keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_REQUIRE_MPPE_128)
		keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_MPPE_STATEFUL)
	}
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_NOBSDCOMP)
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_NODEFLATE)
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_NO_VJ_COMP)
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_NO_PCOMP)
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_NO_ACCOMP)
	keys = appendAvailableKeys(keys, fieldVpnL2tpPpp, NM_SETTING_VPN_L2TP_KEY_LCP_ECHO_FAILURE)
	return
}
func getSettingVpnL2tpPppAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnL2tpPppValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	return
}
func logicSetSettingVpnL2tpKeyRequireMppe(data connectionData, value bool) (err error) {
	if !value {
		removeSettingVpnL2tpKeyRequireMppe40(data)
		removeSettingVpnL2tpKeyRequireMppe128(data)
		removeSettingVpnL2tpKeyMppeStateful(data)
	}
	setSettingVpnL2tpKeyRequireMppe(data, value)
	return
}

// vpn-l2tp-ipsec
func getSettingVpnL2tpIpsecAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnL2tpIpsec, NM_SETTING_VPN_L2TP_KEY_IPSEC_ENABLE)
	if getSettingVpnL2tpKeyIpsecEnable(data) {
		keys = appendAvailableKeys(keys, fieldVpnL2tpIpsec, NM_SETTING_VPN_L2TP_KEY_IPSEC_GATEWAY_ID)
		keys = appendAvailableKeys(keys, fieldVpnL2tpIpsec, NM_SETTING_VPN_L2TP_KEY_IPSEC_GROUP_NAME)
		keys = appendAvailableKeys(keys, fieldVpnL2tpIpsec, NM_SETTING_VPN_L2TP_KEY_IPSEC_PSK)
	}
	return
}
func getSettingVpnL2tpIpsecAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnL2tpIpsecValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	// TODO
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
