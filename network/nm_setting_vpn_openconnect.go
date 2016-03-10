/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

const (
	NM_DBUS_SERVICE_OPENCONNECT   = "org.freedesktop.NetworkManager.openconnect"
	NM_DBUS_INTERFACE_OPENCONNECT = "org.freedesktop.NetworkManager.openconnect"
	NM_DBUS_PATH_OPENCONNECT      = "/org/freedesktop/NetworkManager/openconnect"
)

const (
	nmVpnOpenconnectNameFile = VPN_NAME_FILES_DIR + "nm-openconnect-service.name"
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
	data = newBasicVpnConnectionData(id, uuid)
	initSettingSectionVpnOpenconnect(data)
	initSettingSectionIpv6(data)
	return
}

func initSettingSectionVpnOpenconnect(data connectionData) {
	initBasicSettingSectionVpn(data, NM_DBUS_SERVICE_OPENCONNECT)

	setSettingVpnOpenconnectKeyCsdEnable(data, false)
	setSettingVpnOpenconnectKeyPemPassphraseFsid(data, false)
	setSettingVpnOpenconnectKeyStokenSource(data, "disabled")
	setSettingVpnOpenconnectKeyAuthtype(data, "password")

	if vpnPluginData, ok := doGetSettingVpnPluginData(data, false); ok {
		vpnPluginData["gwcert-flags"] = "2"
		vpnPluginData["cookie-flags"] = "2"
		vpnPluginData["gateway-flags"] = "2"

		vpnPluginData["xmlconfig-flags"] = "0"
		vpnPluginData["lasthost-flags"] = "0"
		vpnPluginData["autoconnect-flags"] = "0"
		vpnPluginData["certsigs-flags"] = "0"
	}
}

func getSettingVpnOpenconnectAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_GATEWAY)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_CACERT)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_PROXY)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_CSD_ENABLE)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_CSD_WRAPPER)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_USERCERT)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_PRIVKEY)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_PEM_PASSPHRASE_FSID)
	return
}
func getSettingVpnOpenconnectAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnOpenconnectValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	ensureSettingVpnOpenconnectKeyGatewayNoEmpty(data, errs)
	checkSettingVpnOpenconnectKeyCacert(data, errs)
	checkSettingVpnOpenconnectKeyUsercert(data, errs)
	checkSettingVpnOpenconnectKeyPrivkey(data, errs)
	return
}
func checkSettingVpnOpenconnectKeyCacert(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenconnectKeyCacertExists(data) {
		return
	}
	value := getSettingVpnOpenconnectKeyCacert(data)
	ensureFileExists(errs, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_CACERT, value,
		".pem", ".crt", ".key")
}
func checkSettingVpnOpenconnectKeyUsercert(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenconnectKeyUsercertExists(data) {
		return
	}
	value := getSettingVpnOpenconnectKeyUsercert(data)
	ensureFileExists(errs, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_USERCERT, value,
		".pem", ".crt", ".key")
}
func checkSettingVpnOpenconnectKeyPrivkey(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenconnectKeyPrivkeyExists(data) {
		return
	}
	value := getSettingVpnOpenconnectKeyPrivkey(data)
	ensureFileExists(errs, sectionVpnOpenconnect, NM_SETTING_VPN_OPENCONNECT_KEY_PRIVKEY, value,
		".pem", ".crt", ".key")
}
