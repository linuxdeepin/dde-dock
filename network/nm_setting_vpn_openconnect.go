package main

const (
	NM_DBUS_SERVICE_OPENCONNECT   = "org.freedesktop.NetworkManager.openconnect"
	NM_DBUS_INTERFACE_OPENCONNECT = "org.freedesktop.NetworkManager.openconnect"
	NM_DBUS_PATH_OPENCONNECT      = "/org/freedesktop/NetworkManager/openconnect"
)

const (
	NM_OPENCONNECT_KEY_GATEWAY             = "gateway"
	NM_OPENCONNECT_KEY_COOKIE              = "cookie"
	NM_OPENCONNECT_KEY_GWCERT              = "gwcert"
	NM_OPENCONNECT_KEY_AUTHTYPE            = "authtype"
	NM_OPENCONNECT_KEY_USERCERT            = "usercert"
	NM_OPENCONNECT_KEY_CACERT              = "cacert"
	NM_OPENCONNECT_KEY_PRIVKEY             = "userkey"
	NM_OPENCONNECT_KEY_MTU                 = "mtu"
	NM_OPENCONNECT_KEY_PEM_PASSPHRASE_FSID = "pem_passphrase_fsid"
	NM_OPENCONNECT_KEY_PROXY               = "proxy"
	NM_OPENCONNECT_KEY_CSD_ENABLE          = "enable_csd_trojan"
	NM_OPENCONNECT_KEY_CSD_WRAPPER         = "csd_wrapper"
	NM_OPENCONNECT_KEY_STOKEN_SOURCE       = "stoken_source"
	NM_OPENCONNECT_KEY_STOKEN_STRING       = "stoken_string"
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

func newVpnOpenconnectConnectionData(id, uuid string) (data _ConnectionData) {
	data = make(_ConnectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, typeVpn)

	// TODO

	return
}
