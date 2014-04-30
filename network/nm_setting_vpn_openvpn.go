package main

import (
	"dlib"
)

const (
	NM_DBUS_SERVICE_OPENVPN   = "org.freedesktop.NetworkManager.openvpn"
	NM_DBUS_INTERFACE_OPENVPN = "org.freedesktop.NetworkManager.openvpn"
	NM_DBUS_PATH_OPENVPN      = "/org/freedesktop/NetworkManager/openvpn"
)

const (
	// openvpn
	NM_SETTING_VPN_OPENVPN_KEY_REMOTE               = "remote"
	NM_SETTING_VPN_OPENVPN_KEY_CONNECTION_TYPE      = "connection-type"
	NM_SETTING_VPN_OPENVPN_KEY_USERNAME             = "username"
	NM_SETTING_VPN_OPENVPN_KEY_PASSWORD             = "password"
	NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS       = "password-flags"
	NM_SETTING_VPN_OPENVPN_KEY_CERT                 = "cert"
	NM_SETTING_VPN_OPENVPN_KEY_CA                   = "ca"
	NM_SETTING_VPN_OPENVPN_KEY_KEY                  = "key"
	NM_SETTING_VPN_OPENVPN_KEY_CERTPASS             = "cert-pass"
	NM_SETTING_VPN_OPENVPN_KEY_CERTPASS_FLAGS       = "cert-pass-flags"
	NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY           = "static-key"
	NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION = "static-key-direction"
	NM_SETTING_VPN_OPENVPN_KEY_REMOTE_IP            = "remote-ip"
	NM_SETTING_VPN_OPENVPN_KEY_LOCAL_IP             = "local-ip"

	// advanced
	NM_SETTING_VPN_OPENVPN_KEY_PORT          = "port"
	NM_SETTING_VPN_OPENVPN_KEY_RENEG_SECONDS = "reneg-seconds"
	NM_SETTING_VPN_OPENVPN_KEY_COMP_LZO      = "comp-lzo"
	NM_SETTING_VPN_OPENVPN_KEY_PROTO_TCP     = "proto-tcp"
	NM_SETTING_VPN_OPENVPN_KEY_TAP_DEV       = "tap-dev"
	NM_SETTING_VPN_OPENVPN_KEY_TUNNEL_MTU    = "tunnel-mtu"
	NM_SETTING_VPN_OPENVPN_KEY_FRAGMENT_SIZE = "fragment-size"
	NM_SETTING_VPN_OPENVPN_KEY_MSSFIX        = "mssfix"
	NM_SETTING_VPN_OPENVPN_KEY_REMOTE_RANDOM = "remote-random"

	// security
	NM_SETTING_VPN_OPENVPN_KEY_CIPHER = "cipher"
	NM_SETTING_VPN_OPENVPN_KEY_AUTH   = "auth"

	// tls auth
	NM_SETTING_VPN_OPENVPN_KEY_TLS_REMOTE      = "tls-remote"
	NM_SETTING_VPN_OPENVPN_KEY_REMOTE_CERT_TLS = "remote-cert-tls"
	NM_SETTING_VPN_OPENVPN_KEY_TA              = "ta"
	NM_SETTING_VPN_OPENVPN_KEY_TA_DIR          = "ta-dir"

	// proxies
	NM_SETTING_VPN_OPENVPN_KEY_PROXY_TYPE                = "proxy-type"
	NM_SETTING_VPN_OPENVPN_KEY_PROXY_SERVER              = "proxy-server"
	NM_SETTING_VPN_OPENVPN_KEY_PROXY_PORT                = "proxy-port"
	NM_SETTING_VPN_OPENVPN_KEY_PROXY_RETRY               = "proxy-retry"
	NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_USERNAME       = "http-proxy-username"
	NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_PASSWORD       = "http-proxy-password"
	NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_PASSWORD_FLAGS = "http-proxy-password-flags"

	/* Internal auth-dialog -> service token indicating that no secrets are
	 * required for the connection.
	 */
	NM_SETTING_VPN_OPENVPN_KEY_NOSECRET = "no-secret"
)

const (
	NM_OPENVPN_SECRET_FLAG_SAVE   = 0
	NM_OPENVPN_SECRET_FLAG_ASK    = 2
	NM_OPENVPN_SECRET_FLAG_UNUSED = 4
)
const (
	NM_OPENVPN_CONTYPE_TLS          = "tls"
	NM_OPENVPN_CONTYPE_STATIC_KEY   = "static-key"
	NM_OPENVPN_CONTYPE_PASSWORD     = "password"
	NM_OPENVPN_CONTYPE_PASSWORD_TLS = "password-tls"
)
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
	data = newBasicVpnConnectionData(id, uuid, NM_DBUS_SERVICE_OPENVPN)
	setSettingVpnOpenvpnKeyConnectionType(data, "tls")
	setSettingVpnOpenvpnKeyCertpassFlags(data, 1)

	addSettingField(data, fieldIpv6)
	setSettingIp6ConfigMethod(data, NM_SETTING_IP6_CONFIG_METHOD_AUTO)

	return
}

// openvpn
func getSettingVpnOpenvpnAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_REMOTE)
	keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CONNECTION_TYPE)
	switch getSettingVpnOpenvpnKeyConnectionType(data) {
	case NM_OPENVPN_CONTYPE_TLS:
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CERT)
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CA)
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_KEY)
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CERTPASS)
	case NM_OPENVPN_CONTYPE_PASSWORD:
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_USERNAME)
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS)
		if getSettingVpnOpenvpnKeyPasswordFlags(data) == NM_OPENVPN_SECRET_FLAG_SAVE {
			keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_PASSWORD)
		}
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CA)
	case NM_OPENVPN_CONTYPE_PASSWORD_TLS:
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_USERNAME)
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS)
		if getSettingVpnOpenvpnKeyPasswordFlags(data) == NM_OPENVPN_SECRET_FLAG_SAVE {
			keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_PASSWORD)
		}
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CERT)
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CA)
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_KEY)
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CERTPASS)
	case NM_OPENVPN_CONTYPE_STATIC_KEY:
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY)
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION)
		if getSettingVkVpnOpenvpnKeyEnableStaticKeyDirection(data) {
			keys = append(keys, NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION)
		}
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_REMOTE_IP)
		keys = appendAvailableKeys(keys, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_LOCAL_IP)
	}
	return
}
func getSettingVpnOpenvpnAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_OPENVPN_KEY_CONNECTION_TYPE:
		values = []kvalue{
			kvalue{NM_OPENVPN_CONTYPE_TLS, dlib.Tr("Certificates (TLS)")},
			kvalue{NM_OPENVPN_CONTYPE_PASSWORD, dlib.Tr("Password")},
			kvalue{NM_OPENVPN_CONTYPE_PASSWORD_TLS, dlib.Tr("Password with Certificates (TLS)")},
			kvalue{NM_OPENVPN_CONTYPE_STATIC_KEY, dlib.Tr("Static Key")},
		}
	case NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION:
		values = []kvalue{
			kvalue{0, dlib.Tr("0")},
			kvalue{1, dlib.Tr("1")},
		}
	case NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS:
		values = []kvalue{
			kvalue{NM_OPENVPN_SECRET_FLAG_SAVE, dlib.Tr("Saved")},
			kvalue{NM_OPENVPN_SECRET_FLAG_ASK, dlib.Tr("Always Ask")},
			kvalue{NM_OPENVPN_SECRET_FLAG_UNUSED, dlib.Tr("Not Required")},
		}
	}
	return
}
func checkSettingVpnOpenvpnValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	ensureSettingVpnOpenvpnKeyRemoteNoEmpty(data, errs)
	ensureSettingVpnOpenvpnKeyConnectionTypeNoEmpty(data, errs)
	switch getSettingVpnOpenvpnKeyConnectionType(data) {
	case NM_OPENVPN_CONTYPE_TLS:
		ensureSettingVpnOpenvpnKeyCertNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCaNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyKeyNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCertpassNoEmpty(data, errs)
	case NM_OPENVPN_CONTYPE_PASSWORD:
		ensureSettingVpnOpenvpnKeyUsernameNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCaNoEmpty(data, errs)
	case NM_OPENVPN_CONTYPE_PASSWORD_TLS:
		ensureSettingVpnOpenvpnKeyUsernameNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCertNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCaNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyKeyNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCertpassNoEmpty(data, errs)
	case NM_OPENVPN_CONTYPE_STATIC_KEY:
		ensureSettingVpnOpenvpnKeyStaticKeyNoEmpty(data, errs)
		// TODO not sure the following keys
		// ensureSettingVpnOpenvpnKeyRemoteIpNoEmpty(data, errs)
		// ensureSettingVpnOpenvpnKeyLocalIpNoEmpty(data, errs)
	}
	checkSettingVpnOpenvpnKeyCert(data, errs)
	checkSettingVpnOpenvpnKeyCa(data, errs)
	checkSettingVpnOpenvpnKeyKey(data, errs)
	checkSettingVpnOpenvpnKeyStaticKey(data, errs)
	return
}
func checkSettingVpnOpenvpnKeyCert(data connectionData, errs fieldErrors) {
	if !isSettingVpnOpenvpnKeyCertExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyCert(data)
	ensureFileExists(errs, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CERT, value,
		".pem", ".crt", ".key", ".cer", ".p12")
}
func checkSettingVpnOpenvpnKeyCa(data connectionData, errs fieldErrors) {
	if !isSettingVpnOpenvpnKeyCaExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyCa(data)
	ensureFileExists(errs, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CA, value,
		".pem", ".crt", ".key", ".cer", ".p12")
}
func checkSettingVpnOpenvpnKeyKey(data connectionData, errs fieldErrors) {
	if !isSettingVpnOpenvpnKeyKeyExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyKey(data)
	ensureFileExists(errs, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_KEY, value,
		".pem", ".crt", ".key", ".cer", ".p12")
}
func checkSettingVpnOpenvpnKeyStaticKey(data connectionData, errs fieldErrors) {
	if !isSettingVpnOpenvpnKeyStaticKeyExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyStaticKey(data)
	ensureFileExists(errs, fieldVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY, value, ".key")
}

// openvpn-advanced general
func getSettingVpnOpenvpnAdvancedAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_PORT)
	if getSettingVkVpnOpenvpnKeyEnablePort(data) {
		keys = append(keys, NM_SETTING_VPN_OPENVPN_KEY_PORT)
	}
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_RENEG_SECONDS)
	if getSettingVkVpnOpenvpnKeyEnableRenegSeconds(data) {
		keys = append(keys, NM_SETTING_VPN_OPENVPN_KEY_RENEG_SECONDS)
	}
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_COMP_LZO)
	if !isSettingVpnOpenvpnKeyProxyTypeExists(data) {
		// when proxy enabled, use a tcp connection default
		keys = appendAvailableKeys(keys, fieldVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_PROTO_TCP)
	}
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_TAP_DEV)
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_TUNNEL_MTU)
	if getSettingVkVpnOpenvpnKeyEnableTunnelMtu(data) {
		keys = append(keys, NM_SETTING_VPN_OPENVPN_KEY_TUNNEL_MTU)
	}
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_FRAGMENT_SIZE)
	if getSettingVkVpnOpenvpnKeyEnableFragmentSize(data) {
		keys = append(keys, NM_SETTING_VPN_OPENVPN_KEY_FRAGMENT_SIZE)
	}
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_MSSFIX)
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_REMOTE_RANDOM)
	return
}
func getSettingVpnOpenvpnAdvancedAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnOpenvpnAdvancedValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	return
}

// openvpn-security
func getSettingVpnOpenvpnSecurityAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnSecurity, NM_SETTING_VPN_OPENVPN_KEY_CIPHER)
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnSecurity, NM_SETTING_VPN_OPENVPN_KEY_AUTH)
	return
}
func getSettingVpnOpenvpnSecurityAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_OPENVPN_KEY_CIPHER:
		// TODO get openvpn cipher
		// "/usr/sbin/openvpn"
		// "/sbin/openvpn"
		// --show-ciphers
		values = []kvalue{
			kvalue{"", dlib.Tr("Default")},
			kvalue{"DES-CBC", dlib.Tr("DES-CBC")},
			kvalue{"RC2-CBC", dlib.Tr("RC2-CBC")},
			kvalue{"DES-EDE-CBC", dlib.Tr("DES-EDE-CBC")},
			kvalue{"DES-EDE3-CBC", dlib.Tr("DES-EDE3-CBC")},
			kvalue{"DESX-CBC", dlib.Tr("DESX-CBC")},
			kvalue{"BF-CBC", dlib.Tr("BF-CBC")},
			kvalue{"RC2-40-CBC", dlib.Tr("RC2-40-CBC")},
			kvalue{"CAST5-CBC", dlib.Tr("CAST5-CBC")},
			kvalue{"RC2-64-CBC", dlib.Tr("RC2-64-CBC")},
			kvalue{"AES-128-CBC", dlib.Tr("AES-128-CBC")},
			kvalue{"AES-192-CBC", dlib.Tr("AES-192-CBC")},
			kvalue{"AES-256-CBC", dlib.Tr("AES-256-CBC")},
			kvalue{"CAMELLIA-128-CBC", dlib.Tr("CAMELLIA-128-CBC")},
			kvalue{"CAMELLIA-192-CBC", dlib.Tr("CAMELLIA-192-CBC")},
			kvalue{"CAMELLIA-256-CBC", dlib.Tr("CAMELLIA-256-CBC")},
			kvalue{"SEED-CBC", dlib.Tr("SEED-CBC")},
		}
	case NM_SETTING_VPN_OPENVPN_KEY_AUTH:
		values = []kvalue{
			kvalue{"", dlib.Tr("Default")},
			kvalue{NM_OPENVPN_AUTH_NONE, dlib.Tr("None")},
			kvalue{NM_OPENVPN_AUTH_RSA_MD4, dlib.Tr("RSA MD-4")},
			kvalue{NM_OPENVPN_AUTH_MD5, dlib.Tr("MD-5")},
			kvalue{NM_OPENVPN_AUTH_SHA1, dlib.Tr("SHA-1")},
			kvalue{NM_OPENVPN_AUTH_SHA224, dlib.Tr("SHA-224")},
			kvalue{NM_OPENVPN_AUTH_SHA256, dlib.Tr("SHA-256")},
			kvalue{NM_OPENVPN_AUTH_SHA384, dlib.Tr("SHA-384")},
			kvalue{NM_OPENVPN_AUTH_SHA512, dlib.Tr("SHA-512")},
			kvalue{NM_OPENVPN_AUTH_RIPEMD160, dlib.Tr("RIPEMD-160")},
		}
	}
	return
}
func checkSettingVpnOpenvpnSecurityValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	return
}

// openvpn-tlsauth
func getSettingVpnOpenvpnTlsauthAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnTlsauth, NM_SETTING_VPN_OPENVPN_KEY_TLS_REMOTE)
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnTlsauth, NM_SETTING_VPN_OPENVPN_KEY_REMOTE_CERT_TLS)
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnTlsauth, NM_SETTING_VPN_OPENVPN_KEY_TA)
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnTlsauth, NM_SETTING_VPN_OPENVPN_KEY_TA_DIR)
	if getSettingVkVpnOpenvpnKeyEnableTaDir(data) {
		keys = append(keys, NM_SETTING_VPN_OPENVPN_KEY_TA_DIR)
	}
	return
}
func getSettingVpnOpenvpnTlsauthAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_OPENVPN_KEY_REMOTE_CERT_TLS:
		values = []kvalue{
			kvalue{"", dlib.Tr("Default")}, // default
			kvalue{NM_OPENVPN_REM_CERT_TLS_CLIENT, dlib.Tr("Client")},
			kvalue{NM_OPENVPN_REM_CERT_TLS_SERVER, dlib.Tr("Server")},
		}
	case NM_SETTING_VPN_OPENVPN_KEY_TA_DIR:
		values = []kvalue{
			kvalue{0, dlib.Tr("0")},
			kvalue{1, dlib.Tr("1")},
		}
	}
	return
}
func checkSettingVpnOpenvpnTlsauthValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	checkSettingVpnOpenvpnKeyTa(data, errs)
	return
}
func checkSettingVpnOpenvpnKeyTa(data connectionData, errs fieldErrors) {
	if !isSettingVpnOpenvpnKeyTaExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyTa(data)
	ensureFileExists(errs, fieldVpnOpenvpnTlsauth, NM_SETTING_VPN_OPENVPN_KEY_TA, value)
}

// openvpn-proxies
func getSettingVpnOpenvpnProxiesAvailableKeys(data connectionData) (keys []string) {
	// proxies
	keys = appendAvailableKeys(keys, fieldVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_TYPE)
	if isSettingVpnOpenvpnKeyProxyTypeExists(data) {
		switch getSettingVpnOpenvpnKeyProxyType(data) {
		case "httpect":
			keys = appendAvailableKeys(keys, fieldVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_SERVER)
			keys = appendAvailableKeys(keys, fieldVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_PORT)
			keys = appendAvailableKeys(keys, fieldVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_RETRY)
			keys = appendAvailableKeys(keys, fieldVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_USERNAME)
			keys = appendAvailableKeys(keys, fieldVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_PASSWORD)
		case "socksct":
			keys = appendAvailableKeys(keys, fieldVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_SERVER)
			keys = appendAvailableKeys(keys, fieldVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_PORT)
			keys = appendAvailableKeys(keys, fieldVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_RETRY)
		}
	}
	return
}
func getSettingVpnOpenvpnProxiesAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_OPENVPN_KEY_PROXY_TYPE:
		values = []kvalue{
			kvalue{"none", dlib.Tr("Not Required")},
			kvalue{"httpect", dlib.Tr("HTTP")},
			kvalue{"socksct", dlib.Tr("SOCKS")},
		}
	}
	return
}
func checkSettingVpnOpenvpnProxiesValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	switch getSettingVpnOpenvpnKeyProxyType(data) {
	case "httpect":
		ensureSettingVpnOpenvpnKeyProxyServerNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyProxyPortNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyProxyRetryNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyHttpProxyUsernameNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyHttpProxyPasswordNoEmpty(data, errs)
	case "socksct":
		ensureSettingVpnOpenvpnKeyProxyServerNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyProxyPortNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyProxyRetryNoEmpty(data, errs)
	}
	return
}

// Logic setter
func logicSetSettingVpnOpenvpnKeyConnectionType(data connectionData, value string) (err error) {
	allRelatedKeys := []string{
		NM_SETTING_VPN_OPENVPN_KEY_USERNAME,
		NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS,
		NM_SETTING_VPN_OPENVPN_KEY_PASSWORD,
		NM_SETTING_VPN_OPENVPN_KEY_CERT,
		NM_SETTING_VPN_OPENVPN_KEY_CA,
		NM_SETTING_VPN_OPENVPN_KEY_KEY,
		NM_SETTING_VPN_OPENVPN_KEY_CERTPASS,
		NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY,
		NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION,
		NM_SETTING_VPN_OPENVPN_KEY_REMOTE_IP,
		NM_SETTING_VPN_OPENVPN_KEY_LOCAL_IP,
	}
	switch value {
	case NM_OPENVPN_CONTYPE_TLS:
		removeSettingKey(data, fieldVpnOpenvpn, stringArrayBut(allRelatedKeys,
			NM_SETTING_VPN_OPENVPN_KEY_CERT,
			NM_SETTING_VPN_OPENVPN_KEY_CA,
			NM_SETTING_VPN_OPENVPN_KEY_KEY,
			NM_SETTING_VPN_OPENVPN_KEY_CERTPASS,
		)...)
	case NM_OPENVPN_CONTYPE_PASSWORD:
		removeSettingKey(data, fieldVpnOpenvpn, stringArrayBut(allRelatedKeys,
			NM_SETTING_VPN_OPENVPN_KEY_USERNAME,
			NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS,
			NM_SETTING_VPN_OPENVPN_KEY_PASSWORD,
			NM_SETTING_VPN_OPENVPN_KEY_CA,
		)...)
	case NM_OPENVPN_CONTYPE_PASSWORD_TLS:
		removeSettingKey(data, fieldVpnOpenvpn, stringArrayBut(allRelatedKeys,
			NM_SETTING_VPN_OPENVPN_KEY_USERNAME,
			NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS,
			NM_SETTING_VPN_OPENVPN_KEY_PASSWORD,
			NM_SETTING_VPN_OPENVPN_KEY_CERT,
			NM_SETTING_VPN_OPENVPN_KEY_CA,
			NM_SETTING_VPN_OPENVPN_KEY_KEY,
			NM_SETTING_VPN_OPENVPN_KEY_CERTPASS,
		)...)
	case NM_OPENVPN_CONTYPE_STATIC_KEY:
		removeSettingKey(data, fieldVpnOpenvpn, stringArrayBut(allRelatedKeys,
			NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY,
			NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION,
			NM_SETTING_VPN_OPENVPN_KEY_REMOTE_IP,
			NM_SETTING_VPN_OPENVPN_KEY_LOCAL_IP,
		)...)
	}
	return
}
func logicSetSettingVpnOpenvpnKeyProxyType(data connectionData, value string) (err error) {
	if value == "none" {
		removeSettingVpnOpenvpnKeyProxyServer(data)
		removeSettingVpnOpenvpnKeyProxyPort(data)
		removeSettingVpnOpenvpnKeyProxyRetry(data)
		removeSettingVpnOpenvpnKeyHttpProxyUsername(data)
		removeSettingVpnOpenvpnKeyHttpProxyPassword(data)
		removeSettingVpnOpenvpnKeyProxyType(data)
		return
	}

	// when proxy enabled, use a tcp connection default
	setSettingVpnOpenvpnKeyProtoTcp(data, true)

	switch value {
	case "httpect":
		setSettingVpnOpenvpnKeyProxyRetry(data, false)
	case "socksct":
		setSettingVpnOpenvpnKeyProxyRetry(data, false)
		removeSettingVpnOpenvpnKeyHttpProxyUsername(data)
		removeSettingVpnOpenvpnKeyHttpProxyPassword(data)
	}
	setSettingVpnOpenvpnKeyProxyType(data, value)
	return
}
func logicSetSettingVpnOpenvpnKeyCert(data connectionData, value string) (err error) {
	setSettingVpnOpenvpnKeyCert(data, toLocalPath(value))
	return
}
func logicSetSettingVpnOpenvpnKeyCa(data connectionData, value string) (err error) {
	setSettingVpnOpenvpnKeyCa(data, toLocalPath(value))
	return
}
func logicSetSettingVpnOpenvpnKeyKey(data connectionData, value string) (err error) {
	setSettingVpnOpenvpnKeyKey(data, toLocalPath(value))
	return
}
func logicSetSettingVpnOpenvpnKeyStaticKey(data connectionData, value string) (err error) {
	setSettingVpnOpenvpnKeyStaticKey(data, toLocalPath(value))
	return
}
