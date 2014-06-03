package network

import . "dlib/gettext"

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

// Define secret flags
const (
	NM_OPENVPN_SECRET_FLAG_SAVE   = 0
	NM_OPENVPN_SECRET_FLAG_ASK    = 2
	NM_OPENVPN_SECRET_FLAG_UNUSED = 4
)

var availableValuesNMOpenvpnSecretFlag = []kvalue{
	kvalue{NM_OPENVPN_SECRET_FLAG_SAVE, Tr("Saved")}, // system saved
	kvalue{NM_OPENVPN_SECRET_FLAG_ASK, Tr("Always Ask")},
	kvalue{NM_OPENVPN_SECRET_FLAG_UNUSED, Tr("Not Required")},
}

func isVpnOpenvpnRequireSecret(flag uint32) bool {
	if flag == NM_OPENVPN_SECRET_FLAG_SAVE {
		return true
	}
	return false
}

func isVpnOpenvpnNeedShowPassword(data connectionData) bool {
	return isVpnOpenvpnRequireSecret(getSettingVpnOpenvpnKeyPasswordFlags(data))
}

func isVpnOpenvpnNeedShowCertpass(data connectionData) bool {
	return isVpnOpenvpnRequireSecret(getSettingVpnOpenvpnKeyCertpassFlags(data))
}

func isVpnOpenvpnNeedShowHttpProxyPassword(data connectionData) bool {
	return isVpnOpenvpnRequireSecret(getSettingVpnOpenvpnKeyHttpProxyPasswordFlags(data))
}

// new connection data
func newVpnOpenvpnConnectionData(id, uuid string) (data connectionData) {
	data = newBasicVpnConnectionData(id, uuid, NM_DBUS_SERVICE_OPENVPN)
	setSettingVpnOpenvpnKeyConnectionType(data, "tls")
	setSettingVpnOpenvpnKeyCertpassFlags(data, NM_OPENVPN_SECRET_FLAG_SAVE)

	initSettingSectionIpv6(data)
	return
}

// openvpn
func getSettingVpnOpenvpnAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_REMOTE)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CONNECTION_TYPE)
	switch getSettingVpnOpenvpnKeyConnectionType(data) {
	case NM_OPENVPN_CONTYPE_TLS:
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CERT)
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CA)
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_KEY)
		if isVpnOpenvpnNeedShowCertpass(data) {
			keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CERTPASS)
		}
	case NM_OPENVPN_CONTYPE_PASSWORD:
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_USERNAME)
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS)
		if isVpnOpenvpnNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_PASSWORD)
		}
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CA)
	case NM_OPENVPN_CONTYPE_PASSWORD_TLS:
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_USERNAME)
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS)
		if isVpnOpenvpnNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_PASSWORD)
		}
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CERT)
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CA)
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_KEY)
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CERTPASS)
	case NM_OPENVPN_CONTYPE_STATIC_KEY:
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY)
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION)
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_REMOTE_IP)
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_LOCAL_IP)
	}
	return
}
func getSettingVpnOpenvpnAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_OPENVPN_KEY_CONNECTION_TYPE:
		values = []kvalue{
			kvalue{NM_OPENVPN_CONTYPE_TLS, Tr("Certificates (TLS)")},
			kvalue{NM_OPENVPN_CONTYPE_PASSWORD, Tr("Password")},
			kvalue{NM_OPENVPN_CONTYPE_PASSWORD_TLS, Tr("Password with Certificates (TLS)")},
			kvalue{NM_OPENVPN_CONTYPE_STATIC_KEY, Tr("Static Key")},
		}
	case NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION:
		values = []kvalue{
			kvalue{0, Tr("0")},
			kvalue{1, Tr("1")},
		}
	case NM_SETTING_VPN_OPENVPN_KEY_CERTPASS_FLAGS:
		values = availableValuesNMOpenvpnSecretFlag
	case NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS:
		values = availableValuesNMOpenvpnSecretFlag
	}
	return
}
func checkSettingVpnOpenvpnValues(data connectionData) (errs sectionErrors) {
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
		if isVpnOpenvpnNeedShowPassword(data) {
			ensureSettingVpnOpenvpnKeyPasswordNoEmpty(data, errs)
		}
	case NM_OPENVPN_CONTYPE_PASSWORD_TLS:
		ensureSettingVpnOpenvpnKeyUsernameNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCertNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCaNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyKeyNoEmpty(data, errs)
		if isVpnOpenvpnNeedShowCertpass(data) {
			ensureSettingVpnOpenvpnKeyCertpassNoEmpty(data, errs)
		}
		if isVpnOpenvpnNeedShowPassword(data) {
			ensureSettingVpnOpenvpnKeyPasswordNoEmpty(data, errs)
		}
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
func checkSettingVpnOpenvpnKeyCert(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenvpnKeyCertExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyCert(data)
	ensureFileExists(errs, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CERT, value,
		".pem", ".crt", ".key", ".cer", ".p12")
}
func checkSettingVpnOpenvpnKeyCa(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenvpnKeyCaExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyCa(data)
	ensureFileExists(errs, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_CA, value,
		".pem", ".crt", ".key", ".cer", ".p12")
}
func checkSettingVpnOpenvpnKeyKey(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenvpnKeyKeyExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyKey(data)
	ensureFileExists(errs, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_KEY, value,
		".pem", ".crt", ".key", ".cer", ".p12")
}
func checkSettingVpnOpenvpnKeyStaticKey(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenvpnKeyStaticKeyExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyStaticKey(data)
	ensureFileExists(errs, sectionVpnOpenvpn, NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY, value, ".key")
}

// openvpn-advanced general
func getSettingVpnOpenvpnAdvancedAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_PORT)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_RENEG_SECONDS)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_COMP_LZO)
	if !isSettingVpnOpenvpnKeyProxyTypeExists(data) {
		// when proxy enabled, use a tcp connection default
		keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_PROTO_TCP)
	}
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_TAP_DEV)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_TUNNEL_MTU)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_FRAGMENT_SIZE)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_MSSFIX)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnAdvanced, NM_SETTING_VPN_OPENVPN_KEY_REMOTE_RANDOM)
	return
}
func getSettingVpnOpenvpnAdvancedAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnOpenvpnAdvancedValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	return
}

// openvpn-security
func getSettingVpnOpenvpnSecurityAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnSecurity, NM_SETTING_VPN_OPENVPN_KEY_CIPHER)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnSecurity, NM_SETTING_VPN_OPENVPN_KEY_AUTH)
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
			kvalue{"", Tr("Default")},
			kvalue{"DES-CBC", Tr("DES-CBC")},
			kvalue{"RC2-CBC", Tr("RC2-CBC")},
			kvalue{"DES-EDE-CBC", Tr("DES-EDE-CBC")},
			kvalue{"DES-EDE3-CBC", Tr("DES-EDE3-CBC")},
			kvalue{"DESX-CBC", Tr("DESX-CBC")},
			kvalue{"BF-CBC", Tr("BF-CBC")},
			kvalue{"RC2-40-CBC", Tr("RC2-40-CBC")},
			kvalue{"CAST5-CBC", Tr("CAST5-CBC")},
			kvalue{"RC2-64-CBC", Tr("RC2-64-CBC")},
			kvalue{"AES-128-CBC", Tr("AES-128-CBC")},
			kvalue{"AES-192-CBC", Tr("AES-192-CBC")},
			kvalue{"AES-256-CBC", Tr("AES-256-CBC")},
			kvalue{"CAMELLIA-128-CBC", Tr("CAMELLIA-128-CBC")},
			kvalue{"CAMELLIA-192-CBC", Tr("CAMELLIA-192-CBC")},
			kvalue{"CAMELLIA-256-CBC", Tr("CAMELLIA-256-CBC")},
			kvalue{"SEED-CBC", Tr("SEED-CBC")},
		}
	case NM_SETTING_VPN_OPENVPN_KEY_AUTH:
		values = []kvalue{
			kvalue{"", Tr("Default")},
			kvalue{NM_OPENVPN_AUTH_NONE, Tr("None")},
			kvalue{NM_OPENVPN_AUTH_RSA_MD4, Tr("RSA MD-4")},
			kvalue{NM_OPENVPN_AUTH_MD5, Tr("MD-5")},
			kvalue{NM_OPENVPN_AUTH_SHA1, Tr("SHA-1")},
			kvalue{NM_OPENVPN_AUTH_SHA224, Tr("SHA-224")},
			kvalue{NM_OPENVPN_AUTH_SHA256, Tr("SHA-256")},
			kvalue{NM_OPENVPN_AUTH_SHA384, Tr("SHA-384")},
			kvalue{NM_OPENVPN_AUTH_SHA512, Tr("SHA-512")},
			kvalue{NM_OPENVPN_AUTH_RIPEMD160, Tr("RIPEMD-160")},
		}
	}
	return
}
func checkSettingVpnOpenvpnSecurityValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	return
}

// openvpn-tlsauth
func getSettingVpnOpenvpnTlsauthAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnTlsauth, NM_SETTING_VPN_OPENVPN_KEY_TLS_REMOTE)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnTlsauth, NM_SETTING_VPN_OPENVPN_KEY_REMOTE_CERT_TLS)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnTlsauth, NM_SETTING_VPN_OPENVPN_KEY_TA)
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnTlsauth, NM_SETTING_VPN_OPENVPN_KEY_TA_DIR)
	return
}
func getSettingVpnOpenvpnTlsauthAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_OPENVPN_KEY_REMOTE_CERT_TLS:
		values = []kvalue{
			kvalue{"", Tr("Default")}, // default
			kvalue{NM_OPENVPN_REM_CERT_TLS_CLIENT, Tr("Client")},
			kvalue{NM_OPENVPN_REM_CERT_TLS_SERVER, Tr("Server")},
		}
	case NM_SETTING_VPN_OPENVPN_KEY_TA_DIR:
		values = []kvalue{
			kvalue{0, Tr("0")},
			kvalue{1, Tr("1")},
		}
	}
	return
}
func checkSettingVpnOpenvpnTlsauthValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	checkSettingVpnOpenvpnKeyTa(data, errs)
	return
}
func checkSettingVpnOpenvpnKeyTa(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenvpnKeyTaExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyTa(data)
	ensureFileExists(errs, sectionVpnOpenvpnTlsauth, NM_SETTING_VPN_OPENVPN_KEY_TA, value)
}

// openvpn-proxies
func getSettingVpnOpenvpnProxiesAvailableKeys(data connectionData) (keys []string) {
	// proxies
	keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_TYPE)
	if isSettingVpnOpenvpnKeyProxyTypeExists(data) {
		switch getSettingVpnOpenvpnKeyProxyType(data) {
		case "httpect":
			keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_SERVER)
			keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_PORT)
			keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_RETRY)
			keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_USERNAME)
			if isVpnOpenvpnNeedShowHttpProxyPassword(data) {
				keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_PASSWORD)
			}
		case "socksct":
			keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_SERVER)
			keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_PORT)
			keys = appendAvailableKeys(data, keys, sectionVpnOpenvpnProxies, NM_SETTING_VPN_OPENVPN_KEY_PROXY_RETRY)
		}
	}
	return
}
func getSettingVpnOpenvpnProxiesAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_OPENVPN_KEY_PROXY_TYPE:
		values = []kvalue{
			kvalue{"none", Tr("Not Required")},
			kvalue{"httpect", Tr("HTTP")},
			kvalue{"socksct", Tr("SOCKS")},
		}
	}
	return
}
func checkSettingVpnOpenvpnProxiesValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	switch getSettingVpnOpenvpnKeyProxyType(data) {
	case "httpect":
		ensureSettingVpnOpenvpnKeyProxyServerNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyProxyPortNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyProxyRetryNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyHttpProxyUsernameNoEmpty(data, errs)
		if isVpnOpenvpnNeedShowHttpProxyPassword(data) {
			ensureSettingVpnOpenvpnKeyHttpProxyPasswordNoEmpty(data, errs)
		}
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
		NM_SETTING_VPN_OPENVPN_KEY_CERTPASS_FLAGS,
		NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY,
		NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION,
		NM_SETTING_VPN_OPENVPN_KEY_REMOTE_IP,
		NM_SETTING_VPN_OPENVPN_KEY_LOCAL_IP,
	}
	switch value {
	case NM_OPENVPN_CONTYPE_TLS:
		removeSettingKey(data, sectionVpnOpenvpn, stringArrayBut(allRelatedKeys,
			NM_SETTING_VPN_OPENVPN_KEY_CERT,
			NM_SETTING_VPN_OPENVPN_KEY_CA,
			NM_SETTING_VPN_OPENVPN_KEY_KEY,
			NM_SETTING_VPN_OPENVPN_KEY_CERTPASS_FLAGS,
			NM_SETTING_VPN_OPENVPN_KEY_CERTPASS,
		)...)
		setSettingVpnOpenvpnKeyCertpassFlags(data, NM_OPENVPN_SECRET_FLAG_SAVE)
	case NM_OPENVPN_CONTYPE_PASSWORD:
		removeSettingKey(data, sectionVpnOpenvpn, stringArrayBut(allRelatedKeys,
			NM_SETTING_VPN_OPENVPN_KEY_USERNAME,
			NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS,
			NM_SETTING_VPN_OPENVPN_KEY_PASSWORD,
			NM_SETTING_VPN_OPENVPN_KEY_CA,
		)...)
	case NM_OPENVPN_CONTYPE_PASSWORD_TLS:
		removeSettingKey(data, sectionVpnOpenvpn, stringArrayBut(allRelatedKeys,
			NM_SETTING_VPN_OPENVPN_KEY_USERNAME,
			NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS,
			NM_SETTING_VPN_OPENVPN_KEY_PASSWORD,
			NM_SETTING_VPN_OPENVPN_KEY_CERT,
			NM_SETTING_VPN_OPENVPN_KEY_CA,
			NM_SETTING_VPN_OPENVPN_KEY_KEY,
			NM_SETTING_VPN_OPENVPN_KEY_CERTPASS_FLAGS,
			NM_SETTING_VPN_OPENVPN_KEY_CERTPASS,
		)...)
	case NM_OPENVPN_CONTYPE_STATIC_KEY:
		removeSettingKey(data, sectionVpnOpenvpn, stringArrayBut(allRelatedKeys,
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
		removeSettingVpnOpenvpnKeyHttpProxyPasswordFlags(data)
		removeSettingVpnOpenvpnKeyProxyType(data)
		return
	}

	// when proxy enabled, use a tcp connection default
	setSettingVpnOpenvpnKeyProtoTcp(data, true)

	switch value {
	case "httpect":
		setSettingVpnOpenvpnKeyProxyRetry(data, false)
		setSettingVpnOpenvpnKeyHttpProxyPasswordFlags(data, NM_OPENVPN_SECRET_FLAG_SAVE)
	case "socksct":
		setSettingVpnOpenvpnKeyProxyRetry(data, false)
		removeSettingVpnOpenvpnKeyHttpProxyUsername(data)
		removeSettingVpnOpenvpnKeyHttpProxyPassword(data)
		removeSettingVpnOpenvpnKeyHttpProxyPasswordFlags(data)
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
