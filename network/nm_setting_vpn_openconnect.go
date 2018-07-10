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

import (
	"pkg.deepin.io/dde/daemon/network/nm"
)

const (
	nmVpnOpenconnectNameFile = "nm-openconnect-service.name"
)

func newVpnOpenconnectConnectionData(id, uuid string) (data connectionData) {
	data = newBasicVpnConnectionData(id, uuid)
	initSettingSectionVpnOpenconnect(data)
	initSettingSectionIpv6(data)
	return
}

func initSettingSectionVpnOpenconnect(data connectionData) {
	initBasicSettingSectionVpn(data, nm.NM_DBUS_SERVICE_OPENCONNECT)

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
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_GATEWAY)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_CACERT)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_PROXY)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_CSD_ENABLE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_CSD_WRAPPER)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_USERCERT)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_PRIVKEY)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_PEM_PASSPHRASE_FSID)
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
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_CACERT, value)
}
func checkSettingVpnOpenconnectKeyUsercert(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenconnectKeyUsercertExists(data) {
		return
	}
	value := getSettingVpnOpenconnectKeyUsercert(data)
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_USERCERT, value)
}
func checkSettingVpnOpenconnectKeyPrivkey(data connectionData, errs sectionErrors) {
	if !isSettingVpnOpenconnectKeyPrivkeyExists(data) {
		return
	}
	value := getSettingVpnOpenconnectKeyPrivkey(data)
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME, nm.NM_SETTING_VPN_OPENCONNECT_KEY_PRIVKEY, value)
}
