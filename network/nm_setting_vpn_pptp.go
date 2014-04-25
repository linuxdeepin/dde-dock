package main

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

func newVpnPptpConnectionData(id, uuid string) (data _ConnectionData) {
	data = make(_ConnectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, typeVpn)
	setSettingConnectionAutoconnect(data, false)
	logicSetSettingVkConnectionNoPermission(data, false)

	addSettingField(data, fieldVpn)
	setSettingVpnServiceType(data, NM_DBUS_SERVICE_PPTP)
	setSettingVpnPptpKeyPasswordFlags(data, 1)

	addSettingField(data, fieldIpv4)
	setSettingIp4ConfigMethod(data, NM_SETTING_IP4_CONFIG_METHOD_AUTO)

	return
}

// TODO
// vpn-pptp
func getSettingVpnPptpAvailableKeys(data _ConnectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnPptp, NM_SETTING_VPN_PPTP_KEY_GATEWAY)
	keys = appendAvailableKeys(keys, fieldVpnPptp, NM_SETTING_VPN_PPTP_KEY_USER)
	keys = appendAvailableKeys(keys, fieldVpnPptp, NM_SETTING_VPN_PPTP_KEY_PASSWORD)
	keys = appendAvailableKeys(keys, fieldVpnPptp, NM_SETTING_VPN_PPTP_KEY_DOMAIN)
	return
}
func getSettingVpnPptpAvailableValues(data _ConnectionData, key string) (values []string, customizable bool) {
	return
}
func checkSettingVpnPptpValues(data _ConnectionData) (errs FieldKeyErrors) {
	errs = make(map[string]string)
	ensureSettingVpnPptpKeyGatewayNoEmpty(data, errs)
	return
}

// vpn-pptp-ppp
func getSettingVpnPptpPppAvailableKeys(data _ConnectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REFUSE_EAP)
	keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REFUSE_PAP)
	keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REFUSE_CHAP)
	keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REFUSE_MSCHAP)
	keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REFUSE_MSCHAPV2)
	keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE)
	if getSettingVpnPptpKeyRequireMppe(data) {
		keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE_40)
		keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_REQUIRE_MPPE_128)
		keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_MPPE_STATEFUL)
	}
	keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_NOBSDCOMP)
	keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_NODEFLATE)
	keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_NO_VJ_COMP)
	keys = appendAvailableKeys(keys, fieldVpnPptpPpp, NM_SETTING_VPN_PPTP_KEY_LCP_ECHO_FAILURE)
	return
}
func getSettingVpnPptpPppAvailableValues(data _ConnectionData, key string) (values []string, customizable bool) {
	return
}
func checkSettingVpnPptpPppValues(data _ConnectionData) (errs FieldKeyErrors) {
	errs = make(map[string]string)
	return
}
func logicSetSettingVpnPptpKeyRequireMppeJSON(data _ConnectionData, valueJSON string) {
	setSettingVpnPptpKeyRequireMppeJSON(data, valueJSON)
	value := getSettingVpnPptpKeyRequireMppe(data)
	logicSetSettingVpnPptpKeyRequireMppe(data, value)
}
func logicSetSettingVpnPptpKeyRequireMppe(data _ConnectionData, value bool) {
	if !value {
		removeSettingVpnPptpKeyRequireMppe40(data)
		removeSettingVpnPptpKeyRequireMppe128(data)
		removeSettingVpnPptpKeyMppeStateful(data)
	}
	setSettingVpnPptpKeyRequireMppe(data, value)
}
