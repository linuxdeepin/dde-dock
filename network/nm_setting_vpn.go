package main

import (
	"dlib/dbus"
)

// TODO doc
const NM_SETTING_VPN_SETTING_NAME = "vpn"

const (
	NM_SETTING_VPN_SERVICE_TYPE = "service-type"
	NM_SETTING_VPN_USER_NAME    = "user-name"
	NM_SETTING_VPN_DATA         = "data"
	NM_SETTING_VPN_SECRETS      = "secrets"
)

// VPN connection states
const (
	NM_VPN_CONNECTION_STATE_UNKNOWN       = 0
	NM_VPN_CONNECTION_STATE_PREPARE       = 1
	NM_VPN_CONNECTION_STATE_NEED_AUTH     = 2
	NM_VPN_CONNECTION_STATE_CONNECT       = 3
	NM_VPN_CONNECTION_STATE_IP_CONFIG_GET = 4
	NM_VPN_CONNECTION_STATE_ACTIVATED     = 5
	NM_VPN_CONNECTION_STATE_FAILED        = 6
	NM_VPN_CONNECTION_STATE_DISCONNECTE   = 7
)

func newBasicVpnConnectionData(id, uuid, service string) (data connectionData) {
	data = make(connectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, NM_SETTING_VPN_SETTING_NAME)
	setSettingConnectionAutoconnect(data, false)
	logicSetSettingVkConnectionNoPermission(data, false)

	addSettingField(data, fieldVpn)
	setSettingVpnServiceType(data, service)
	setSettingVpnData(data, make(map[string]string))
	setSettingVpnSecrets(data, make(map[string]string))

	initSettingFieldIpv4(data)
	return
}

func getSettingVpnAvailableKeys(data connectionData) (keys []string) { return }
func getSettingVpnAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	return
}

func isSettingVpnPluginKey(field string) bool {
	// all keys in vpn virtual fields are vpn plugin key
	realField := getRealFieldName(field)
	if realField == fieldVpn && realField != field {
		return true
	}
	return false
}
func isSettingVpnPluginSecretKey(field, key string) bool {
	switch field {
	case fieldVpnL2tp:
		switch key {
		case NM_SETTING_VPN_L2TP_KEY_PASSWORD:
			return true
		}
	case fieldVpnOpenvpn:
		switch key {
		case NM_SETTING_VPN_OPENVPN_KEY_PASSWORD:
			return true
		case NM_SETTING_VPN_OPENVPN_KEY_CERTPASS:
			return true
		}
	case fieldVpnOpenvpnProxies:
		switch key {
		case NM_SETTING_VPN_OPENVPN_KEY_HTTP_PROXY_PASSWORD:
			return true
		}
	case fieldVpnPptp:
		switch key {
		case NM_SETTING_VPN_PPTP_KEY_PASSWORD:
			return true
		}
	case fieldVpnVpnc:
		switch key {
		case NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD:
			return true
		case NM_SETTING_VPN_VPNC_KEY_SECRET:
			return true
		}
	}
	return false
}

// Basic getter and setter for vpn plugin keys
func getSettingVpnPluginKey(data connectionData, field, key string) (value interface{}) {
	vpnData, ok := getSettingVpnPluginData(data, field, key)
	if !ok {
		// not exists, just return nil
		logger.Errorf("invalid vpn plugin data: data[%s][%s]", field, key)
		return
	}
	valueStr, ok := vpnData[key]
	if !ok {
		return
	}
	value = unmarshalVpnPluginKey(valueStr, generalGetSettingKeyType(field, key))
	return
}
func setSettingVpnPluginKey(data connectionData, field, key string, value interface{}) {
	vpnData, ok := getSettingVpnPluginData(data, field, key)
	if !ok {
		logger.Errorf("invalid vpn plugin data: data[%s][%s]", field, key)
		return
	}
	valueStr := marshalVpnPluginKey(value, generalGetSettingKeyType(field, key))
	vpnData[key] = valueStr
}
func isSettingVpnPluginKeyExists(data connectionData, field, key string) (ok bool) {
	vpnData, ok := getSettingVpnPluginData(data, field, key)
	if !ok {
		return
	}
	_, ok = vpnData[key]
	return
}
func removeSettingVpnPluginKey(data connectionData, field string, keys ...string) {
	vpnSerectData, ok := doGetSettingVpnPluginData(data, true)
	if ok {
		doRemoveSettingVpnPluginKey(vpnSerectData, keys...)
	}
	vpnData, ok := doGetSettingVpnPluginData(data, false)
	if ok {
		doRemoveSettingVpnPluginKey(vpnData, keys...)
	}
}
func removeSettingVpnPluginKeyBut(data connectionData, field string, keys ...string) {
	vpnSerectData, ok := doGetSettingVpnPluginData(data, true)
	if ok {
		doRemoveSettingVpnPluginKeyBut(vpnSerectData, keys...)
	}
	vpnData, ok := doGetSettingVpnPluginData(data, false)
	if ok {
		doRemoveSettingVpnPluginKeyBut(vpnData, keys...)
	}
}
func doRemoveSettingVpnPluginKey(vpnData map[string]string, keys ...string) {
	for _, k := range keys {
		delete(vpnData, k)
	}
}
func doRemoveSettingVpnPluginKeyBut(vpnData map[string]string, keys ...string) {
	for k := range vpnData {
		if !isStringInArray(k, keys) {
			delete(vpnData, k)
		}
	}
}

func getSettingVpnPluginData(data connectionData, field, key string) (vpnData map[string]string, ok bool) {
	if isSettingVpnPluginSecretKey(field, key) {
		vpnData, ok = doGetSettingVpnPluginData(data, true)
	} else {
		vpnData, ok = doGetSettingVpnPluginData(data, false)
	}
	return
}
func doGetSettingVpnPluginData(data connectionData, isSecretKey bool) (vpnData map[string]string, ok bool) {
	vpnFieldData, ok := data[fieldVpn]
	if !ok {
		return
	}
	var variantValue dbus.Variant
	if isSecretKey {
		variantValue, ok = vpnFieldData[NM_SETTING_VPN_SECRETS]
		if !ok {
			return
		}
	} else {
		variantValue, ok = vpnFieldData[NM_SETTING_VPN_DATA]
		if !ok {
			return
		}
	}
	vpnData, err := interfaceToDictStringString(variantValue.Value())
	if err != nil {
		ok = false
		logger.Error("invalid vpn plugin data:", err)
	} else {
		ok = true
	}
	return
}

// "string" -> "string", 123 -> "123", true -> "true"
func marshalVpnPluginKey(value interface{}, t ktype) (valueStr string) {
	var err error
	switch t {
	default:
		logger.Error("unknown vpn plugin key type", t)
	case ktypeString:
		valueStr, _ = value.(string)
	case ktypeUint32:
		valueStr, err = marshalJSON(value)
	case ktypeBoolean:
		valueBoolean, _ := value.(bool)
		if valueBoolean {
			valueStr = "yes"
		} else {
			valueStr = "no"
		}
	}
	if err != nil {
		logger.Error(err)
	}
	return
}

// "string" -> "string", "123" -> 123, "true" -> true
func unmarshalVpnPluginKey(valueStr string, t ktype) (value interface{}) {
	var err error
	switch t {
	default:
		logger.Error("unknown vpn plugin key type", t)
	case ktypeString:
		value = valueStr
	case ktypeUint32:
		value, err = jsonToKeyValueUint32(valueStr)
	case ktypeBoolean:
		if valueStr == "yes" {
			value = true
		} else if valueStr == "no" {
			value = false
		} else {
			logger.Error("invalid vpn boolean key", valueStr)
		}
	}
	if err != nil {
		logger.Error(err)
	}
	return
}
