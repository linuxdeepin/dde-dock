/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package network

import . "pkg.deepin.io/lib/gettext"
import "pkg.deepin.io/dde/daemon/network/nm"

const (
	nmVpnOpenvpnNameFile = "nm-openvpn-service.name"
)

var availableValuesNmOpenvpnSecretFlags []kvalue

func initAvailableValuesNmOpenvpnSecretFlags() {
	availableValuesNmOpenvpnSecretFlags = []kvalue{
		kvalue{nm.NM_OPENVPN_SECRET_FLAG_SAVE, Tr("Saved")}, // system saved
		kvalue{nm.NM_OPENVPN_SECRET_FLAG_ASK, Tr("Always Ask")},
		kvalue{nm.NM_OPENVPN_SECRET_FLAG_UNUSED, Tr("Not Required")},
	}
}

func isVpnOpenvpnRequireSecret(flag uint32) bool {
	if flag == nm.NM_OPENVPN_SECRET_FLAG_SAVE {
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
	data = newBasicVpnConnectionData(id, uuid)
	initSettingSectionVpnOpenvpn(data)
	initSettingSectionIpv6(data)
	return
}

func initSettingSectionVpnOpenvpn(data connectionData) {
	initBasicSettingSectionVpn(data, nm.NM_DBUS_SERVICE_OPENVPN)
	setSettingVpnOpenvpnKeyConnectionType(data, "tls")
	setSettingVpnOpenvpnKeyCertpassFlags(data, nm.NM_OPENVPN_SECRET_FLAG_SAVE)
}

// openvpn
func getSettingVpnOpenvpnAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_REMOTE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CONNECTION_TYPE)
	switch getSettingVpnOpenvpnKeyConnectionType(data) {
	case nm.NM_OPENVPN_CONTYPE_TLS:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CERT)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CA)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_KEY)
		if isVpnOpenvpnNeedShowCertpass(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CERTPASS)
		}
	case nm.NM_OPENVPN_CONTYPE_PASSWORD:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_USERNAME)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS)
		if isVpnOpenvpnNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD)
		}
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CA)
	case nm.NM_OPENVPN_CONTYPE_PASSWORD_TLS:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_USERNAME)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS)
		if isVpnOpenvpnNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD)
		}
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CERT)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CA)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_KEY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CERTPASS)
	case nm.NM_OPENVPN_CONTYPE_STATIC_KEY:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_REMOTE_IP)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_LOCAL_IP)
	}
	return
}
func getSettingVpnOpenvpnAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_VPN_OPENVPN_KEY_CONNECTION_TYPE:
		values = []kvalue{
			kvalue{nm.NM_OPENVPN_CONTYPE_TLS, Tr("Certificates (TLS)")},
			kvalue{nm.NM_OPENVPN_CONTYPE_PASSWORD, Tr("Password")},
			kvalue{nm.NM_OPENVPN_CONTYPE_PASSWORD_TLS, Tr("Certificates with Password (TLS)")},
			kvalue{nm.NM_OPENVPN_CONTYPE_STATIC_KEY, Tr("Static Key")},
		}
	case nm.NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION:
		values = []kvalue{
			kvalue{0, Tr("0")},
			kvalue{1, Tr("1")},
		}
	case nm.NM_SETTING_VPN_OPENVPN_KEY_CERTPASS_FLAGS:
		values = availableValuesNmOpenvpnSecretFlags
	case nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS:
		values = availableValuesNmOpenvpnSecretFlags
	}
	return
}
func checkSettingVpnOpenvpnValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	ensureSettingVpnOpenvpnKeyRemoteNoEmpty(data, errs)
	ensureSettingVpnOpenvpnKeyConnectionTypeNoEmpty(data, errs)
	switch getSettingVpnOpenvpnKeyConnectionType(data) {
	case nm.NM_OPENVPN_CONTYPE_TLS:
		ensureSettingVpnOpenvpnKeyCertNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCaNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyKeyNoEmpty(data, errs)
	case nm.NM_OPENVPN_CONTYPE_PASSWORD:
		ensureSettingVpnOpenvpnKeyUsernameNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCaNoEmpty(data, errs)
		if isVpnOpenvpnNeedShowPassword(data) {
			ensureSettingVpnOpenvpnKeyPasswordNoEmpty(data, errs)
		}
	case nm.NM_OPENVPN_CONTYPE_PASSWORD_TLS:
		ensureSettingVpnOpenvpnKeyUsernameNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCertNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyCaNoEmpty(data, errs)
		ensureSettingVpnOpenvpnKeyKeyNoEmpty(data, errs)
		if isVpnOpenvpnNeedShowPassword(data) {
			ensureSettingVpnOpenvpnKeyPasswordNoEmpty(data, errs)
		}
	case nm.NM_OPENVPN_CONTYPE_STATIC_KEY:
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
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CERT, value)
}
func checkSettingVpnOpenvpnKeyCa(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenvpnKeyCaExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyCa(data)
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CA, value)
}
func checkSettingVpnOpenvpnKeyKey(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenvpnKeyKeyExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyKey(data)
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_KEY, value)
}
func checkSettingVpnOpenvpnKeyStaticKey(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenvpnKeyStaticKeyExists(data) {
		return
	}
	value := getSettingVpnOpenvpnKeyStaticKey(data)
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY, value)
}

// openvpn-advanced general
func getSettingVpnOpenvpnAdvancedAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PORT)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_RENEG_SECONDS)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_COMP_LZO)
	if !isSettingVpnOpenvpnKeyProxyTypeExists(data) {
		// when proxy enabled, use a tcp connection default
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PROTO_TCP)
	}
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_TAP_DEV)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_TUNNEL_MTU)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_FRAGMENT_SIZE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_MSSFIX)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_REMOTE_RANDOM)
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
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SECURITY_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_CIPHER)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SECURITY_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_AUTH)
	return
}
func getSettingVpnOpenvpnSecurityAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_VPN_OPENVPN_KEY_CIPHER:
		// TODO get openvpn cipher "/usr/sbin/openvpn" "/sbin/openvpn" --show-ciphers
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
	case nm.NM_SETTING_VPN_OPENVPN_KEY_AUTH:
		values = []kvalue{
			kvalue{"", Tr("Default")},
			kvalue{nm.NM_OPENVPN_AUTH_NONE, Tr("None")},
			kvalue{nm.NM_OPENVPN_AUTH_RSA_MD4, Tr("RSA MD-4")},
			kvalue{nm.NM_OPENVPN_AUTH_MD5, Tr("MD-5")},
			kvalue{nm.NM_OPENVPN_AUTH_SHA1, Tr("SHA-1")},
			kvalue{nm.NM_OPENVPN_AUTH_SHA224, Tr("SHA-224")},
			kvalue{nm.NM_OPENVPN_AUTH_SHA256, Tr("SHA-256")},
			kvalue{nm.NM_OPENVPN_AUTH_SHA384, Tr("SHA-384")},
			kvalue{nm.NM_OPENVPN_AUTH_SHA512, Tr("SHA-512")},
			kvalue{nm.NM_OPENVPN_AUTH_RIPEMD160, Tr("RIPEMD-160")},
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
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_TLSAUTH_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_TLS_REMOTE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_TLSAUTH_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_REMOTE_CERT_TLS)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_TLSAUTH_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_TA)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_TLSAUTH_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_TA_DIR)
	return
}
func getSettingVpnOpenvpnTlsauthAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_VPN_OPENVPN_KEY_REMOTE_CERT_TLS:
		values = []kvalue{
			kvalue{"", Tr("Default")}, // default
			kvalue{nm.NM_OPENVPN_REM_CERT_TLS_CLIENT, Tr("Client")},
			kvalue{nm.NM_OPENVPN_REM_CERT_TLS_SERVER, Tr("Server")},
		}
	case nm.NM_SETTING_VPN_OPENVPN_KEY_TA_DIR:
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
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_OPENVPN_TLSAUTH_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_TA, value)
}

// openvpn-proxies
func getSettingVpnOpenvpnProxiesAvailableKeys(data connectionData) (keys []string) {
	// proxies
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PROXY_TYPE)
	if isSettingVpnOpenvpnKeyProxyTypeExists(data) {
		switch getSettingVpnOpenvpnKeyProxyType(data) {
		case "httpect":
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PROXY_SERVER)
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PROXY_PORT)
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PROXY_RETRY)
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_USERNAME)
			if isVpnOpenvpnNeedShowHttpProxyPassword(data) {
				keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_PASSWORD)
			}
		case "socksct":
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PROXY_SERVER)
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PROXY_PORT)
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME, nm.NM_SETTING_VPN_OPENVPN_KEY_PROXY_RETRY)
		}
	}
	return
}
func getSettingVpnOpenvpnProxiesAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_VPN_OPENVPN_KEY_PROXY_TYPE:
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
		nm.NM_SETTING_VPN_OPENVPN_KEY_USERNAME,
		nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS,
		nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD,
		nm.NM_SETTING_VPN_OPENVPN_KEY_CERT,
		nm.NM_SETTING_VPN_OPENVPN_KEY_CA,
		nm.NM_SETTING_VPN_OPENVPN_KEY_KEY,
		nm.NM_SETTING_VPN_OPENVPN_KEY_CERTPASS,
		nm.NM_SETTING_VPN_OPENVPN_KEY_CERTPASS_FLAGS,
		nm.NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY,
		nm.NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION,
		nm.NM_SETTING_VPN_OPENVPN_KEY_REMOTE_IP,
		nm.NM_SETTING_VPN_OPENVPN_KEY_LOCAL_IP,
	}
	switch value {
	case nm.NM_OPENVPN_CONTYPE_TLS:
		removeSettingKey(data, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, stringArrayBut(allRelatedKeys,
			nm.NM_SETTING_VPN_OPENVPN_KEY_CERT,
			nm.NM_SETTING_VPN_OPENVPN_KEY_CA,
			nm.NM_SETTING_VPN_OPENVPN_KEY_KEY,
			nm.NM_SETTING_VPN_OPENVPN_KEY_CERTPASS_FLAGS,
			nm.NM_SETTING_VPN_OPENVPN_KEY_CERTPASS,
		)...)
		setSettingVpnOpenvpnKeyCertpassFlags(data, nm.NM_OPENVPN_SECRET_FLAG_SAVE)
	case nm.NM_OPENVPN_CONTYPE_PASSWORD:
		removeSettingKey(data, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, stringArrayBut(allRelatedKeys,
			nm.NM_SETTING_VPN_OPENVPN_KEY_USERNAME,
			nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS,
			nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD,
			nm.NM_SETTING_VPN_OPENVPN_KEY_CA,
		)...)
	case nm.NM_OPENVPN_CONTYPE_PASSWORD_TLS:
		removeSettingKey(data, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, stringArrayBut(allRelatedKeys,
			nm.NM_SETTING_VPN_OPENVPN_KEY_USERNAME,
			nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD_FLAGS,
			nm.NM_SETTING_VPN_OPENVPN_KEY_PASSWORD,
			nm.NM_SETTING_VPN_OPENVPN_KEY_CERT,
			nm.NM_SETTING_VPN_OPENVPN_KEY_CA,
			nm.NM_SETTING_VPN_OPENVPN_KEY_KEY,
			nm.NM_SETTING_VPN_OPENVPN_KEY_CERTPASS_FLAGS,
			nm.NM_SETTING_VPN_OPENVPN_KEY_CERTPASS,
		)...)
	case nm.NM_OPENVPN_CONTYPE_STATIC_KEY:
		removeSettingKey(data, nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME, stringArrayBut(allRelatedKeys,
			nm.NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY,
			nm.NM_SETTING_VPN_OPENVPN_KEY_STATIC_KEY_DIRECTION,
			nm.NM_SETTING_VPN_OPENVPN_KEY_REMOTE_IP,
			nm.NM_SETTING_VPN_OPENVPN_KEY_LOCAL_IP,
		)...)
	}
	setSettingVpnOpenvpnKeyConnectionType(data, value)
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
		setSettingVpnOpenvpnKeyHttpProxyPasswordFlags(data, nm.NM_OPENVPN_SECRET_FLAG_SAVE)
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
