/**
 * Copyright (c) 2016 Deepin, Inc.
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package network

import (
	. "pkg.deepin.io/lib/gettext"
)

// For the NM <-> VPN plugin service
const (
	NM_DBUS_SERVICE_STRONGSWAN = "org.freedesktop.NetworkManager.strongswan"
	nmVpnStrongswanNameFile    = VPN_NAME_FILES_DIR + "nm-strongswan-service.name"
)

const (
	NM_SETTING_VPN_STRONGSWAN_KEY_ADDRESS        = "address"
	NM_SETTING_VPN_STRONGSWAN_KEY_CERTIFICATE    = "certificate"
	NM_SETTING_VPN_STRONGSWAN_KEY_METHOD         = "method"
	NM_SETTING_VPN_STRONGSWAN_KEY_USER           = "user"
	NM_SETTING_VPN_STRONGSWAN_KEY_USERCERT       = "usercert"
	NM_SETTING_VPN_STRONGSWAN_KEY_USERKEY        = "userkey"
	NM_SETTING_VPN_STRONGSWAN_KEY_PASSWORD       = "password"
	NM_SETTING_VPN_STRONGSWAN_KEY_PASSWORD_FLAGS = "password-flags"
	NM_SETTING_VPN_STRONGSWAN_KEY_VIRTUAL        = "virtual"
	NM_SETTING_VPN_STRONGSWAN_KEY_ENCAP          = "encap"
	NM_SETTING_VPN_STRONGSWAN_KEY_IPCOMP         = "ipcomp"
)

const (
	NM_STRONGSWAN_METHOD_KEY       = "key"
	NM_STRONGSWAN_METHOD_AGENT     = "agent"
	NM_STRONGSWAN_METHOD_SMARTCARD = "smartcard"
	NM_STRONGSWAN_METHOD_EAP       = "eap"
	NM_STRONGSWAN_METHOD_PSK       = "psk"
)

// new connection data
func newVpnStrongswanConnectionData(id, uuid string) (data connectionData) {
	data = newBasicVpnConnectionData(id, uuid)
	initSettingSectionVpnStrongswan(data)
	return
}

func initSettingSectionVpnStrongswan(data connectionData) {
	initBasicSettingSectionVpn(data, NM_DBUS_SERVICE_STRONGSWAN)
	setSettingVpnStrongswanKeyMethod(data, NM_STRONGSWAN_METHOD_KEY)
	setSettingVpnStrongswanKeyPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	setSettingVpnStrongswanKeyVirtual(data, true)
}

// strongswan
func getSettingVpnStrongswanAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_ADDRESS)
	keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_CERTIFICATE)

	keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_METHOD)
	switch getSettingVpnStrongswanKeyMethod(data) {
	case NM_STRONGSWAN_METHOD_KEY:
		keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_USERCERT)
		keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_USERKEY)
	case NM_STRONGSWAN_METHOD_AGENT:
		keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_USERCERT)
	case NM_STRONGSWAN_METHOD_SMARTCARD:
	case NM_STRONGSWAN_METHOD_EAP:
		keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_USER)
		keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_PASSWORD)
	case NM_STRONGSWAN_METHOD_PSK:
		keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_USER)
		keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_PASSWORD)
	}

	keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_VIRTUAL)
	keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_ENCAP)
	keys = appendAvailableKeys(data, keys, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_IPCOMP)
	return
}
func getSettingVpnStrongswanAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_STRONGSWAN_KEY_METHOD:
		values = []kvalue{
			kvalue{NM_STRONGSWAN_METHOD_KEY, Tr("Private Key")},
			kvalue{NM_STRONGSWAN_METHOD_AGENT, Tr("SSH Agent")},
			kvalue{NM_STRONGSWAN_METHOD_SMARTCARD, Tr("Smartcard")},
			kvalue{NM_STRONGSWAN_METHOD_EAP, Tr("EAP")},
			kvalue{NM_STRONGSWAN_METHOD_PSK, Tr("Pre-Shared Key")},
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
	ensureFileExists(errs, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_CERTIFICATE, value,
		".pem", ".crt", ".key", ".cer", ".p12")
}
func checkSettingVpnStrongswanKeyUsercert(data connectionData, errs sectionErrors) {
	if !isSettingVpnStrongswanKeyUsercertExists(data) {
		return
	}
	value := getSettingVpnStrongswanKeyUsercert(data)
	ensureFileExists(errs, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_USERCERT, value,
		".pem", ".crt", ".key", ".cer", ".p12")
}
func checkSettingVpnStrongswanKeyUserkey(data connectionData, errs sectionErrors) {
	if !isSettingVpnStrongswanKeyUserkeyExists(data) {
		return
	}
	value := getSettingVpnStrongswanKeyUserkey(data)
	ensureFileExists(errs, sectionVpnStrongswan, NM_SETTING_VPN_STRONGSWAN_KEY_USERKEY, value,
		".pem", ".crt", ".key", ".cer", ".p12")
}

// Logic setter
func logicSetSettingVpnStrongswanKeyMethod(data connectionData, value string) (err error) {
	switch value {
	case NM_STRONGSWAN_METHOD_KEY:
		removeSettingVpnStrongswanKeyUser(data)
	case NM_STRONGSWAN_METHOD_AGENT:
		removeSettingVpnStrongswanKeyUser(data)
		removeSettingVpnStrongswanKeyUserkey(data)
	case NM_STRONGSWAN_METHOD_SMARTCARD:
		removeSettingVpnStrongswanKeyUser(data)
		removeSettingVpnStrongswanKeyUsercert(data)
		removeSettingVpnStrongswanKeyUserkey(data)
	case NM_STRONGSWAN_METHOD_EAP:
		removeSettingVpnStrongswanKeyUsercert(data)
		removeSettingVpnStrongswanKeyUserkey(data)
	case NM_STRONGSWAN_METHOD_PSK:
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
