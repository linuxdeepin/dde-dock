/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"pkg.deepin.io/dde/daemon/network/nm"
	. "pkg.deepin.io/lib/gettext"
)

const (
	nmVpnVpncNameFile = VPN_NAME_FILES_DIR + "nm-vpnc-service.name"
)

var availableValuesNmVpncSecretFlags []kvalue

func initAvailableValuesNmVpncSecretFlags() {
	availableValuesNmVpncSecretFlags = []kvalue{
		kvalue{nm.NM_VPNC_SECRET_FLAG_NONE, Tr("Saved")},
		// kvalue{nm.NM_VPNC_SECRET_FLAG_SAVE, Tr("Saved")},
		kvalue{nm.NM_VPNC_SECRET_FLAG_ASK, Tr("Always Ask")},
		kvalue{nm.NM_VPNC_SECRET_FLAG_UNUSED, Tr("Not Required")},
	}
}

func isVpnVpncRequireSecret(flag uint32) bool {
	if flag == nm.NM_VPNC_SECRET_FLAG_NONE || flag == nm.NM_VPNC_SECRET_FLAG_SAVE {
		return true
	}
	return false
}

func isVpnVpncNeedShowSecret(data connectionData) bool {
	return isVpnVpncRequireSecret(getSettingVpnVpncKeySecretFlags(data))
}

func isVpnVpncNeedShowXauthPassword(data connectionData) bool {
	return isVpnVpncRequireSecret(getSettingVpnVpncKeyXauthPasswordFlags(data))
}

// new connection data
func newVpnVpncConnectionData(id, uuid string) (data connectionData) {
	data = newBasicVpnConnectionData(id, uuid)
	initSettingSectionVpnVpnc(data)
	return
}

func initSettingSectionVpnVpnc(data connectionData) {
	initBasicSettingSectionVpn(data, nm.NM_DBUS_SERVICE_VPNC)
	setSettingVpnVpncKeyNatTraversalMode(data, nm.NM_VPNC_NATT_MODE_NATT)
	logicSetSettingVpnVpncKeySecretFlags(data, nm.NM_VPNC_SECRET_FLAG_NONE)
	logicSetSettingVpnVpncKeyXauthPasswordFlags(data, nm.NM_VPNC_SECRET_FLAG_NONE)
	setSettingVpnVpncKeyVendor(data, nm.NM_VPNC_VENDOR_CISCO)
	setSettingVpnVpncKeyPerfectForward(data, nm.NM_VPNC_PFS_SERVER)
	setSettingVpnVpncKeyDhgroup(data, nm.NM_VPNC_DHGROUP_DH2)
	setSettingVpnVpncKeyLocalPort(data, 0)
}

// vpn-vpnc
func getSettingVpnVpncAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_GATEWAY)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_XAUTH_USER)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD_FLAGS)
	if isVpnVpncNeedShowXauthPassword(data) {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD)
	}
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_ID)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_SECRET_FLAGS)
	if isVpnVpncNeedShowSecret(data) {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_SECRET)
	}
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_AUTHMODE)
	if getSettingVkVpnVpncKeyHybridAuthmode(data) {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_CA_FILE)
	}
	return
}
func getSettingVpnVpncAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD_FLAGS:
		values = availableValuesNmVpncSecretFlags
	case nm.NM_SETTING_VPN_VPNC_KEY_SECRET_FLAGS:
		values = availableValuesNmVpncSecretFlags
	}
	return
}
func checkSettingVpnVpncValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	if isVpnVpncNeedShowXauthPassword(data) {
		ensureSettingVpnVpncKeyXauthPasswordNoEmpty(data, errs)
	}
	if isVpnVpncNeedShowSecret(data) {
		ensureSettingVpnVpncKeySecretNoEmpty(data, errs)
	}
	ensureSettingVpnVpncKeyGatewayNoEmpty(data, errs)
	ensureSettingVpnVpncKeyIdNoEmpty(data, errs)
	checkSettingVpnVpncCaFile(data, errs)
	return
}
func checkSettingVpnVpncCaFile(data connectionData, errs sectionErrors) {
	if !isSettingVpnVpncKeyCaFileExists(data) {
		return
	}
	value := getSettingVpnVpncKeyCaFile(data)
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_CA_FILE, value)
}

// vpn-vpnc-advanced
func getSettingVpnVpncAdvancedAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_DOMAIN)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_VENDOR)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_APP_VERSION)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_SINGLE_DES)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_NO_ENCRYPTION)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_NAT_TRAVERSAL_MODE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_DHGROUP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_PERFECT_FORWARD)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_LOCAL_PORT)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME, nm.NM_SETTING_VPN_VPNC_KEY_DPD_IDLE_TIMEOUT)
	return
}
func getSettingVpnVpncAdvancedAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_VPN_VPNC_KEY_VENDOR:
		values = []kvalue{
			kvalue{nm.NM_VPNC_VENDOR_CISCO, Tr("Cisco (default)")},
			kvalue{nm.NM_VPNC_VENDOR_NETSCREEN, Tr("Netscreen")},
		}
	case nm.NM_SETTING_VPN_VPNC_KEY_NAT_TRAVERSAL_MODE:
		values = []kvalue{
			kvalue{nm.NM_VPNC_NATT_MODE_NATT, Tr("NAT-T When Available (default)")},
			kvalue{nm.NM_VPNC_NATT_MODE_NATT_ALWAYS, Tr("NAT-T Always")},
			kvalue{nm.NM_VPNC_NATT_MODE_CISCO, Tr("Cisco UDP")},
			kvalue{nm.NM_VPNC_NATT_MODE_NONE, Tr("Disabled")},
		}
	case nm.NM_SETTING_VPN_VPNC_KEY_DHGROUP:
		values = []kvalue{
			kvalue{nm.NM_VPNC_DHGROUP_DH1, Tr("DH Group 1")},
			kvalue{nm.NM_VPNC_DHGROUP_DH2, Tr("DH Group 2 (default)")},
			kvalue{nm.NM_VPNC_DHGROUP_DH5, Tr("DH Group 5")},
		}
	case nm.NM_SETTING_VPN_VPNC_KEY_PERFECT_FORWARD:
		values = []kvalue{
			kvalue{nm.NM_VPNC_PFS_SERVER, Tr("Server (default)")},
			kvalue{nm.NM_VPNC_PFS_NOPFS, Tr("None")},
			kvalue{nm.NM_VPNC_PFS_DH1, Tr("DH Group 1")},
			kvalue{nm.NM_VPNC_PFS_DH2, Tr("DH Group 2")},
			kvalue{nm.NM_VPNC_PFS_DH5, Tr("DH Group 5")},
		}
	}
	return
}
func checkSettingVpnVpncAdvancedValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	return
}

// Logic setter
func logicSetSettingVpnVpncKeySecretFlags(data connectionData, value uint32) (err error) {
	switch value {
	case nm.NM_VPNC_SECRET_FLAG_NONE:
		setSettingVpnVpncKeySecretType(data, nm.NM_VPNC_PW_TYPE_SAVE)
	case nm.NM_VPNC_SECRET_FLAG_SAVE:
		setSettingVpnVpncKeySecretType(data, nm.NM_VPNC_PW_TYPE_SAVE)
	case nm.NM_VPNC_SECRET_FLAG_ASK:
		setSettingVpnVpncKeySecretType(data, nm.NM_VPNC_PW_TYPE_ASK)
	case nm.NM_VPNC_SECRET_FLAG_UNUSED:
		setSettingVpnVpncKeySecretType(data, nm.NM_VPNC_PW_TYPE_UNUSED)
	}
	setSettingVpnVpncKeySecretFlags(data, value)
	return
}
func logicSetSettingVpnVpncKeyXauthPasswordFlags(data connectionData, value uint32) (err error) {
	switch value {
	case nm.NM_VPNC_SECRET_FLAG_NONE:
		setSettingVpnVpncKeyXauthPasswordType(data, nm.NM_VPNC_PW_TYPE_SAVE)
	case nm.NM_VPNC_SECRET_FLAG_SAVE:
		setSettingVpnVpncKeyXauthPasswordType(data, nm.NM_VPNC_PW_TYPE_SAVE)
	case nm.NM_VPNC_SECRET_FLAG_ASK:
		setSettingVpnVpncKeyXauthPasswordType(data, nm.NM_VPNC_PW_TYPE_ASK)
	case nm.NM_VPNC_SECRET_FLAG_UNUSED:
		setSettingVpnVpncKeyXauthPasswordType(data, nm.NM_VPNC_PW_TYPE_UNUSED)
	}
	setSettingVpnVpncKeyXauthPasswordFlags(data, value)
	return
}

// Virtual key getter
func getSettingVkVpnVpncKeyHybridAuthmode(data connectionData) (value bool) {
	if isSettingVpnVpncKeyAuthmodeExists(data) {
		return true
	}
	return false
}
func getSettingVkVpnVpncKeyEncryptionMethod(data connectionData) (value string) {
	if getSettingVpnVpncKeySingleDes(data) {
		return "weak"
	} else if getSettingVpnVpncKeyNoEncryption(data) {
		return "none"
	}
	return "secure"
}
func getSettingVkVpnVpncKeyDisableDpd(data connectionData) (value bool) {
	if isSettingVpnVpncKeyDpdIdleTimeoutExists(data) && getSettingVpnVpncKeyDpdIdleTimeout(data) == 0 {
		return true
	}
	return false
}

// Virtual key logic setter, all virtual keys has a logic setter
func logicSetSettingVkVpnVpncKeyHybridAuthmode(data connectionData, value bool) (err error) {
	if value {
		setSettingVpnVpncKeyAuthmode(data, "hybrid")
	} else {
		removeSettingVpnVpncKeyAuthmode(data)
	}
	return
}
func logicSetSettingVkVpnVpncKeyEncryptionMethod(data connectionData, value string) (err error) {
	removeSettingVpnVpncKeySingleDes(data)
	removeSettingVpnVpncKeyNoEncryption(data)
	switch value {
	case "secure":
	case "weak":
		setSettingVpnVpncKeySingleDes(data, true)
	case "none":
		setSettingVpnVpncKeyNoEncryption(data, true)
	}
	return
}
func logicSetSettingVkVpnVpncKeyDisableDpd(data connectionData, value bool) (err error) {
	if value {
		setSettingVpnVpncKeyDpdIdleTimeout(data, 0)
	} else {
		removeSettingVpnVpncKeyDpdIdleTimeout(data)
	}
	return
}
