package main

const (
	NM_DBUS_SERVICE_OPENVPN   = "org.freedesktop.NetworkManager.openvpn"
	NM_DBUS_INTERFACE_OPENVPN = "org.freedesktop.NetworkManager.openvpn"
	NM_DBUS_PATH_OPENVPN      = "/org/freedesktop/NetworkManager/openvpn"
)

const (
	NM_SETTING_VPN_OPENVPN_KEY_AUTH                      = "auth"
	NM_SETTING_VPN_OPENVPN_KEY_CA                        = "ca"
	NM_SETTING_VPN_OPENVPN_KEY_CERT                      = "cert"
	NM_SETTING_VPN_OPENVPN_KEY_CIPHER                    = "cipher"
	NM_SETTING_VPN_OPENVPN_KEY_COMP_LZO                  = "comp-lzo"
	NM_SETTING_VPN_OPENVPN_KEY_CONNECTION_TYPE           = "connection-type"
	NM_SETTING_VPN_OPENVPN_KEY_FRAGMENT_SIZE             = "fragment-size"
	NM_SETTING_VPN_OPENVPN_KEY_KEY                       = "key"
	NM_SETTING_VPN_OPENVPN_KEY_LOCAL_IP                  = "local-ip" /* ??? */
	NM_SETTING_VPN_OPENVPN_KEY_MSSFIX                    = "mssfix"
	NM_SETTING_VPN_OPENVPN_KEY_PORT                      = "port"
	NM_SETTING_VPN_OPENVPN_KEY_PROTO_TCP                 = "proto-tcp"
	NM_SETTING_VPN_OPENVPN_KEY_PROXY_TYPE                = "proxy-type"
	NM_SETTING_VPN_OPENVPN_KEY_PROXY_SERVER              = "proxy-server"
	NM_SETTING_VPN_OPENVPN_KEY_PROXY_PORT                = "proxy-port"
	NM_SETTING_VPN_OPENVPN_KEY_PROXY_RETRY               = "proxy-retry"
	NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_USERNAME       = "http-proxy-username"
	NM_SETTING_VPN_OPENVPN_KEY_REMOTE                    = "remote"
	NM_SETTING_VPN_OPENVPN_KEY_REMOTE_RANDOM             = "remote-random"
	NM_SETTING_VPN_OPENVPN_KEY_REMOTE_IP                 = "remote-ip"
	NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY                = "static-key"
	NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION      = "static-key-direction"
	NM_SETTING_VPN_OPENVPN_KEY_TA                        = "ta"
	NM_SETTING_VPN_OPENVPN_KEY_TA_DIR                    = "ta-dir"
	NM_SETTING_VPN_OPENVPN_KEY_TUNNEL_MTU                = "tunnel-mtu"
	NM_SETTING_VPN_OPENVPN_KEY_USERNAME                  = "username"
	NM_SETTING_VPN_OPENVPN_KEY_TAP_DEV                   = "tap-dev"
	NM_SETTING_VPN_OPENVPN_KEY_TLS_REMOTE                = "tls-remote"
	NM_SETTING_VPN_OPENVPN_KEY_REMOTE_CERT_TLS           = "remote-cert-tls"
	NM_SETTING_VPN_OPENVPN_KEY_PASSWORD                  = "password"
	NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS            = "password-flags"
	NM_SETTING_VPN_OPENVPN_KEY_CERTPASS                  = "cert-pass"
	NM_SETTING_VPN_OPENVPN_KEY_CERTPASS_FLAGS            = "cert-pass-flags"
	NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_PASSWORD       = "http-proxy-password"
	NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_PASSWORD_FLAGS = "http-proxy-password-flags"
)

/* Internal auth-dialog -> service token indicating that no secrets are
 * required for the connection.
 */
const NM_OPENVPN_KEY_NOSECRET = "no-secret"

// #define NM_OPENVPN_KEY_RENEG_SECONDS "reneg-seconds"
const (
	NM_OPENVPN_AUTH_NONE      = "none"
	NM_OPENVPN_AUTH_RSA_MD4   = "RSA-MD4"
	NM_OPENVPN_AUTH_MD5       = "MD5"
	NM_OPENVPN_AUTH_SHA1      = "SHA1"
	NM_OPENVPN_AUTH_SHA224    = "SHA224"
	NM_OPENVPN_AUTH_SHA256    = "SHA256"
	NM_OPENVPN_AUTH_SHA384    = "SHA384"
	NM_OPENVPN_AUTH_SHA512    = "SHA512"
	NM_OPENVPN_AUTH_RIPEMD160 = "RIPEMD160"
)
const (
	NM_OPENVPN_CONTYPE_TLS          = "tls"
	NM_OPENVPN_CONTYPE_STATIC_KEY   = "static-key"
	NM_OPENVPN_CONTYPE_PASSWORD     = "password"
	NM_OPENVPN_CONTYPE_PASSWORD_TLS = "password-tls"
)

/* arguments of "--remote-cert-tls" */
const (
	NM_OPENVPN_REM_CERT_TLS_CLIENT = "client"
	NM_OPENVPN_REM_CERT_TLS_SERVER = "server"
)

// vpn key descriptions
// static ValidProperty valid_properties[] = {
// 	{ NM_OPENVPN_KEY_AUTH,                 G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_CA,                   G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_CERT,                 G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_CIPHER,               G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_COMP_LZO,             G_TYPE_BOOLEAN, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_CONNECTION_TYPE,      G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_FRAGMENT_SIZE,        G_TYPE_INT, 0, G_MAXINT, FALSE },
// 	{ NM_OPENVPN_KEY_KEY,                  G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_LOCAL_IP,             G_TYPE_STRING, 0, 0, TRUE },
// 	{ NM_OPENVPN_KEY_MSSFIX,               G_TYPE_BOOLEAN, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_PROTO_TCP,            G_TYPE_BOOLEAN, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_PORT,                 G_TYPE_INT, 1, 65535, FALSE },
// 	{ NM_OPENVPN_KEY_PROXY_TYPE,           G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_PROXY_SERVER,         G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_PROXY_PORT,           G_TYPE_INT, 1, 65535, FALSE },
// 	{ NM_OPENVPN_KEY_PROXY_RETRY,          G_TYPE_BOOLEAN, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_HTTP_PROXY_USERNAME,  G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_REMOTE,               G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_REMOTE_RANDOM,        G_TYPE_BOOLEAN, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_REMOTE_IP,            G_TYPE_STRING, 0, 0, TRUE },
// 	{ NM_OPENVPN_KEY_RENEG_SECONDS,        G_TYPE_INT, 0, G_MAXINT, FALSE },
// 	{ NM_OPENVPN_KEY_STATIC_KEY,           G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_STATIC_KEY_DIRECTION, G_TYPE_INT, 0, 1, FALSE },
// 	{ NM_OPENVPN_KEY_TA,                   G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_TA_DIR,               G_TYPE_INT, 0, 1, FALSE },
// 	{ NM_OPENVPN_KEY_TAP_DEV,              G_TYPE_BOOLEAN, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_TLS_REMOTE,           G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_REMOTE_CERT_TLS,      G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_TUNNEL_MTU,           G_TYPE_INT, 0, G_MAXINT, FALSE },
// 	{ NM_OPENVPN_KEY_USERNAME,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_PASSWORD"-flags",     G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_CERTPASS"-flags",     G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_NOSECRET,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_HTTP_PROXY_PASSWORD"-flags", G_TYPE_STRING, 0, 0, FALSE },
// 	{ NULL,                                G_TYPE_NONE, FALSE }
// }
// static ValidProperty valid_secrets[] = {
// 	{ NM_OPENVPN_KEY_PASSWORD,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_CERTPASS,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_NOSECRET,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_HTTP_PROXY_PASSWORD,  G_TYPE_STRING, 0, 0, FALSE },
// 	{ NULL,                                G_TYPE_NONE, FALSE }
// };

func newVpnOpenvpnConnectionData(id, uuid string) (data connectionData) {
	data = newVpnConnectionData(id, uuid, NM_DBUS_SERVICE_OPENVPN)
	// TODO
	return
}
