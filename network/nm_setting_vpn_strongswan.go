/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *             Xu FaSheng <fasheng.xu@gmail.com>
 *
 * Maintainer: Xu FaSheng <fasheng.xu@gmail.com>
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
	. "pkg.deepin.io/lib/gettext"
)

const (
	nmVpnStrongswanNameFile = "nm-strongswan-service.name"
)

// new connection data
func newVpnStrongswanConnectionData(id, uuid string) (data connectionData) {
	data = newBasicVpnConnectionData(id, uuid)
	initSettingSectionVpnStrongswan(data)
	return
}

func initSettingSectionVpnStrongswan(data connectionData) {
	initBasicSettingSectionVpn(data, nm.NM_DBUS_SERVICE_STRONGSWAN)
	setSettingVpnStrongswanKeyMethod(data, nm.NM_STRONGSWAN_METHOD_KEY)
	setSettingVpnStrongswanKeyPasswordFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
	setSettingVpnStrongswanKeyVirtual(data, true)
}

// strongswan
func getSettingVpnStrongswanAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_ADDRESS)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_CERTIFICATE)

	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_METHOD)
	switch getSettingVpnStrongswanKeyMethod(data) {
	case nm.NM_STRONGSWAN_METHOD_KEY:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_USERCERT)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_USERKEY)
	case nm.NM_STRONGSWAN_METHOD_AGENT:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_USERCERT)
	case nm.NM_STRONGSWAN_METHOD_SMARTCARD:
	case nm.NM_STRONGSWAN_METHOD_EAP:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_USER)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_PASSWORD)
	case nm.NM_STRONGSWAN_METHOD_PSK:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_USER)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_PASSWORD)
	}

	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_VIRTUAL)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_ENCAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_IPCOMP)
	return
}
func getSettingVpnStrongswanAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_VPN_STRONGSWAN_KEY_METHOD:
		values = []kvalue{
			kvalue{nm.NM_STRONGSWAN_METHOD_KEY, Tr("Private Key")},
			kvalue{nm.NM_STRONGSWAN_METHOD_AGENT, Tr("SSH Agent")},
			kvalue{nm.NM_STRONGSWAN_METHOD_SMARTCARD, Tr("Smartcard")},
			kvalue{nm.NM_STRONGSWAN_METHOD_EAP, Tr("EAP")},
			kvalue{nm.NM_STRONGSWAN_METHOD_PSK, Tr("Pre-Shared Key")},
		}
	}
	return
}
func checkSettingVpnStrongswanValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	checkSettingVpnStrongswanKeyCertificate(data, errs)
	checkSettingVpnStrongswanKeyUsercert(data, errs)
	checkSettingVpnStrongswanKeyUserkey(data, errs)
	return
}
func checkSettingVpnStrongswanKeyCertificate(data connectionData, errs sectionErrors) {
	if !isSettingVpnStrongswanKeyCertificateExists(data) {
		return
	}
	value := getSettingVpnStrongswanKeyCertificate(data)
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_CERTIFICATE, value)
}
func checkSettingVpnStrongswanKeyUsercert(data connectionData, errs sectionErrors) {
	if !isSettingVpnStrongswanKeyUsercertExists(data) {
		return
	}
	value := getSettingVpnStrongswanKeyUsercert(data)
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_USERCERT, value)
}
func checkSettingVpnStrongswanKeyUserkey(data connectionData, errs sectionErrors) {
	if !isSettingVpnStrongswanKeyUserkeyExists(data) {
		return
	}
	value := getSettingVpnStrongswanKeyUserkey(data)
	ensureFileExists(errs, nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME, nm.NM_SETTING_VPN_STRONGSWAN_KEY_USERKEY, value)
}

// Logic setter
func logicSetSettingVpnStrongswanKeyMethod(data connectionData, value string) (err error) {
	switch value {
	case nm.NM_STRONGSWAN_METHOD_KEY:
		removeSettingVpnStrongswanKeyUser(data)
	case nm.NM_STRONGSWAN_METHOD_AGENT:
		removeSettingVpnStrongswanKeyUser(data)
		removeSettingVpnStrongswanKeyUserkey(data)
	case nm.NM_STRONGSWAN_METHOD_SMARTCARD:
		removeSettingVpnStrongswanKeyUser(data)
		removeSettingVpnStrongswanKeyUsercert(data)
		removeSettingVpnStrongswanKeyUserkey(data)
	case nm.NM_STRONGSWAN_METHOD_EAP:
		removeSettingVpnStrongswanKeyUsercert(data)
		removeSettingVpnStrongswanKeyUserkey(data)
	case nm.NM_STRONGSWAN_METHOD_PSK:
		removeSettingVpnStrongswanKeyUsercert(data)
		removeSettingVpnStrongswanKeyUserkey(data)
	}
	setSettingVpnStrongswanKeyMethod(data, value)
	return
}
func logicSetSettingVpnStrongswanKeyCertificate(data connectionData, value string) (err error) {
	setSettingVpnStrongswanKeyCertificate(data, toLocalPath(value))
	return
}
func logicSetSettingVpnStrongswanKeyUsercert(data connectionData, value string) (err error) {
	setSettingVpnStrongswanKeyUsercert(data, toLocalPath(value))
	return
}
func logicSetSettingVpnStrongswanKeyUserkey(data connectionData, value string) (err error) {
	setSettingVpnStrongswanKeyUserkey(data, toLocalPath(value))
	return
}
