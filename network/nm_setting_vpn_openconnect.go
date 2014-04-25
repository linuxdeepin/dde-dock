package main

const (
	NM_DBUS_SERVICE_OPENCONNECT   = "org.freedesktop.NetworkManager.openconnect"
	NM_DBUS_INTERFACE_OPENCONNECT = "org.freedesktop.NetworkManager.openconnect"
	NM_DBUS_PATH_OPENCONNECT      = "/org/freedesktop/NetworkManager/openconnect"
)

const (
	NM_SETTING_VPN_OPENCONNECT_KEY_GATEWAY             = "gateway"
	NM_SETTING_VPN_OPENCONNECT_KEY_COOKIE              = "cookie"
	NM_SETTING_VPN_OPENCONNECT_KEY_GWCERT              = "gwcert"
	NM_SETTING_VPN_OPENCONNECT_KEY_AUTHTYPE            = "authtype"
	NM_SETTING_VPN_OPENCONNECT_KEY_USERCERT            = "usercert"
	NM_SETTING_VPN_OPENCONNECT_KEY_CACERT              = "cacert"
	NM_SETTING_VPN_OPENCONNECT_KEY_PRIVKEY             = "userkey"
	NM_SETTING_VPN_OPENCONNECT_KEY_MTU                 = "mtu"
	NM_SETTING_VPN_OPENCONNECT_KEY_PEM_PASSPHRASE_FSID = "pem_passphrase_fsid"
	NM_SETTING_VPN_OPENCONNECT_KEY_PROXY               = "proxy"
	NM_SETTING_VPN_OPENCONNECT_KEY_CSD_ENABLE          = "enable_csd_trojan"
	NM_SETTING_VPN_OPENCONNECT_KEY_CSD_WRAPPER         = "csd_wrapper"
	NM_SETTING_VPN_OPENCONNECT_KEY_STOKEN_SOURCE       = "stoken_source"
	NM_SETTING_VPN_OPENCONNECT_KEY_STOKEN_STRING       = "stoken_string"
)

// vpn key descriptions
// static ValidProperty valid_properties[] = {
// 	{ NM_OPENCONNECT_KEY_GATEWAY,     G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_CACERT,      G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_AUTHTYPE,    G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_USERCERT,    G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_PRIVKEY,     G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_MTU,         G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_PEM_PASSPHRASE_FSID, G_TYPE_BOOLEAN, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_PROXY,       G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_CSD_ENABLE,  G_TYPE_BOOLEAN, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_CSD_WRAPPER, G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_STOKEN_SOURCE, G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_STOKEN_STRING, G_TYPE_STRING, 0, 0 },
// 	{ NULL,                           G_TYPE_NONE, 0, 0 }
// };
// static ValidProperty valid_secrets[] = {
// 	{ NM_OPENCONNECT_KEY_COOKIE,  G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_GATEWAY, G_TYPE_STRING, 0, 0 },
// 	{ NM_OPENCONNECT_KEY_GWCERT,  G_TYPE_STRING, 0, 0 },
// 	{ NULL,                       G_TYPE_NONE, 0, 0 }
// };

func newVpnOpenconnectConnectionData(id, uuid string) (data connectionData) {
	data = make(connectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, typeVpn)
	setSettingConnectionAutoconnect(data, false)
	logicSetSettingVkConnectionNoPermission(data, false)

	addSettingField(data, fieldVpn)
	setSettingVpnServiceType(data, NM_DBUS_SERVICE_OPENCONNECT)
	setSettingVpnOpenconnectKeyCsdEnable(data, false)
	setSettingKey(data, fieldVpn, "xmlconfig-flags", uint32(0))
	setSettingVpnOpenconnectKeyPemPassphraseFsid(data, false)
	setSettingKey(data, fieldVpn, "gwcert-flags", uint32(2))
	setSettingKey(data, fieldVpn, "gateway-flags", uint32(2))
	setSettingKey(data, fieldVpn, "autoconnect-flags", uint32(0))
	setSettingKey(data, fieldVpn, "lasthost-flags", uint32(0))
	setSettingVpnOpenconnectKeyStokenSource(data, "disabled")
	setSettingKey(data, fieldVpn, "certsigs-flags", uint32(0))
	setSettingKey(data, fieldVpn, "cookie-flags", uint32(2))
	setSettingVpnOpenconnectKeyAuthtype(data, "password")

	addSettingField(data, fieldIpv4)
	setSettingIp4ConfigMethod(data, NM_SETTING_IP4_CONFIG_METHOD_AUTO)

	addSettingField(data, fieldIpv6)
	setSettingIp6ConfigMethod(data, NM_SETTING_IP6_CONFIG_METHOD_AUTO)

	return
}

func getSettingVpnOpenconnectAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_GATEWAY)
	keys = appendAvailableKeys(keys, fieldVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_CACERT)
	keys = appendAvailableKeys(keys, fieldVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_PROXY)
	keys = appendAvailableKeys(keys, fieldVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_CSD_ENABLE)
	keys = appendAvailableKeys(keys, fieldVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_CSD_WRAPPER)
	keys = appendAvailableKeys(keys, fieldVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_USERCERT)
	keys = appendAvailableKeys(keys, fieldVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_PRIVKEY)
	keys = appendAvailableKeys(keys, fieldVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_PEM_PASSPHRASE_FSID)
	return
}
func getSettingVpnOpenconnectAvailableValues(data connectionData, key string) (values []string, customizable bool) {
	return
}
func checkSettingVpnOpenconnectValues(data connectionData) (errs FieldKeyErrors) {
	errs = make(map[string]string)
	ensureSettingVpnOpenconnectKeyGatewayNoEmpty(data, errs)
	// TODO
	return
}
