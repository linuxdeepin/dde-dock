/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
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

import . "pkg.linuxdeepin.com/lib/gettext"

const (
	NM_DBUS_SERVICE_VPNC   = "org.freedesktop.NetworkManager.vpnc"
	NM_DBUS_INTERFACE_VPNC = "org.freedesktop.NetworkManager.vpnc"
	NM_DBUS_PATH_VPNC      = "/org/freedesktop/NetworkManager/vpnc"
)

const (
	nmVpnVpncServiceFile = VPN_NAME_FILES_DIR + "nm-vpnc-service.name"
	nmVpnVpncServiceBin  = "/usr/lib/NetworkManager/nm-vpnc-service"
	nmVpnVpncHelperBin   = "/usr/lib/NetworkManager/nm-vpnc-service-vpnc-helper"
	nmVpnVpncAuthDlgBin  = "/usr/lib/NetworkManager/nm-vpnc-auth-dialog"
)

const (
	NM_SETTING_VPN_VPNC_KEY_GATEWAY               = "IPSec gateway"
	NM_SETTING_VPN_VPNC_KEY_XAUTH_USER            = "Xauth username"
	NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD        = "Xauth password"
	NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD_TYPE   = "xauth-password-type"
	NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD_FLAGS  = "Xauth password-flags"
	NM_SETTING_VPN_VPNC_KEY_ID                    = "IPSec ID"
	NM_SETTING_VPN_VPNC_KEY_SECRET                = "IPSec secret"
	NM_SETTING_VPN_VPNC_KEY_SECRET_TYPE           = "ipsec-secret-type"
	NM_SETTING_VPN_VPNC_KEY_SECRET_FLAGS          = "IPSec secret-flags"
	NM_SETTING_VPN_VPNC_KEY_AUTHMODE              = "IKE Authmode"
	NM_SETTING_VPN_VPNC_KEY_CA_FILE               = "CA-File"
	NM_SETTING_VPN_VPNC_KEY_DOMAIN                = "Domain"
	NM_SETTING_VPN_VPNC_KEY_VENDOR                = "Vendor"
	NM_SETTING_VPN_VPNC_KEY_APP_VERSION           = "Application Version"
	NM_SETTING_VPN_VPNC_KEY_SINGLE_DES            = "Enable Single DES"
	NM_SETTING_VPN_VPNC_KEY_NO_ENCRYPTION         = "Enable no encryption"
	NM_SETTING_VPN_VPNC_KEY_NAT_TRAVERSAL_MODE    = "NAT Traversal Mode"
	NM_SETTING_VPN_VPNC_KEY_DHGROUP               = "IKE DH Group"
	NM_SETTING_VPN_VPNC_KEY_PERFECT_FORWARD       = "Perfect Forward Secrecy"
	NM_SETTING_VPN_VPNC_KEY_LOCAL_PORT            = "Local Port"
	NM_SETTING_VPN_VPNC_KEY_DPD_IDLE_TIMEOUT      = "DPD idle timeout (our side)"
	NM_SETTING_VPN_VPNC_KEY_CISCO_UDP_ENCAPS_PORT = "Cisco UDP Encapsulation Port"
)

const (
	NM_VPNC_NATT_MODE_NATT        = "natt"
	NM_VPNC_NATT_MODE_NONE        = "none"
	NM_VPNC_NATT_MODE_NATT_ALWAYS = "force-natt"
	NM_VPNC_NATT_MODE_CISCO       = "cisco-udp"
)
const (
	NM_VPNC_PW_TYPE_SAVE   = "save"   // -> flags 1
	NM_VPNC_PW_TYPE_ASK    = "ask"    // -> flags 3
	NM_VPNC_PW_TYPE_UNUSED = "unused" // -> flags 5
)
const (
	NM_VPNC_DHGROUP_DH1 = "dh1"
	NM_VPNC_DHGROUP_DH2 = "dh2"
	NM_VPNC_DHGROUP_DH5 = "dh5"
)
const (
	NM_VPNC_PFS_SERVER = "server"
	NM_VPNC_PFS_NOPFS  = "nopfs"
	NM_VPNC_PFS_DH1    = "dh1"
	NM_VPNC_PFS_DH2    = "dh2"
	NM_VPNC_PFS_DH5    = "dh5"
)
const (
	NM_VPNC_VENDOR_CISCO     = "cisco"
	NM_VPNC_VENDOR_NETSCREEN = "netscreen"
)

// vpn key descriptions
// static ValidProperty valid_properties[] = {
// 	{ NM_VPNC_KEY_GATEWAY,               ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_ID,                    ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_XAUTH_USER,            ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_DOMAIN,                ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_DHGROUP,               ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_PERFECT_FORWARD,       ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_VENDOR,                ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_APP_VERSION,           ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_SINGLE_DES,            ITEM_TYPE_BOOLEAN, 0, 0 },
// 	{ NM_VPNC_KEY_NO_ENCRYPTION,         ITEM_TYPE_BOOLEAN, 0, 0 },
// 	{ NM_VPNC_KEY_DPD_IDLE_TIMEOUT,      ITEM_TYPE_INT, 0, 86400 },
// 	{ NM_VPNC_KEY_NAT_TRAVERSAL_MODE,    ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_CISCO_UDP_ENCAPS_PORT, ITEM_TYPE_INT, 0, 65535 },
// 	{ NM_VPNC_KEY_LOCAL_PORT,            ITEM_TYPE_INT, 0, 65535 },
// 	/* Hybrid Auth */
// 	{ NM_VPNC_KEY_AUTHMODE,              ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_CA_FILE,               ITEM_TYPE_PATH, 0, 0 },
// 	/* Ignored option for internal use */
// 	{ NM_VPNC_KEY_SECRET_TYPE,           ITEM_TYPE_IGNORED, 0, 0 },
// 	{ NM_VPNC_KEY_XAUTH_PASSWORD_TYPE,   ITEM_TYPE_IGNORED, 0, 0 },
// 	{ NM_VPNC_KEY_SECRET"-flags",        ITEM_TYPE_IGNORED, 0, 0 },
// 	{ NM_VPNC_KEY_XAUTH_PASSWORD"-flags",ITEM_TYPE_IGNORED, 0, 0 },
// 	/* Legacy options that are ignored */
// 	{ LEGACY_NAT_KEEPALIVE,              ITEM_TYPE_STRING, 0, 0 },
// 	{ NULL,                              ITEM_TYPE_UNKNOWN, 0, 0 }
// }
// static ValidProperty valid_secrets[] = {
// 	{ NM_OPENVPN_KEY_PASSWORD,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_CERTPASS,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_NOSECRET,             G_TYPE_STRING, 0, 0, FALSE },
// 	{ NM_OPENVPN_KEY_HTTP_PROXY_PASSWORD,  G_TYPE_STRING, 0, 0, FALSE },
// 	{ NULL,                                G_TYPE_NONE, FALSE }
// };
// static ValidProperty valid_secrets[] = {
// 	{ NM_VPNC_KEY_SECRET,                ITEM_TYPE_STRING, 0, 0 },
// 	{ NM_VPNC_KEY_XAUTH_PASSWORD,        ITEM_TYPE_STRING, 0, 0 },
// 	{ NULL,                              ITEM_TYPE_UNKNOWN, 0, 0 }
// };

// Define secret flags
const (
	NM_VPNC_SECRET_FLAG_NONE   = 0
	NM_VPNC_SECRET_FLAG_SAVE   = 1
	NM_VPNC_SECRET_FLAG_ASK    = 3
	NM_VPNC_SECRET_FLAG_UNUSED = 5
)

var availableValuesNMVpncSecretFlag = []kvalue{
	kvalue{NM_VPNC_SECRET_FLAG_NONE, Tr("Saved")},
	// kvalue{NM_VPNC_SECRET_FLAG_SAVE, Tr("Saved")},
	kvalue{NM_VPNC_SECRET_FLAG_ASK, Tr("Always Ask")},
	kvalue{NM_VPNC_SECRET_FLAG_UNUSED, Tr("Not Required")},
}

func isVpnVpncRequireSecret(flag uint32) bool {
	if flag == NM_VPNC_SECRET_FLAG_NONE || flag == NM_VPNC_SECRET_FLAG_SAVE {
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
	initSettingSectionVpnPptp(data)
	return
}

func initSettingSectionVpnVpnc(data connectionData) {
	initBasicSettingSectionVpn(data, NM_DBUS_SERVICE_VPNC)
	setSettingVpnVpncKeyNatTraversalMode(data, NM_VPNC_NATT_MODE_NATT)
	logicSetSettingVpnVpncKeySecretFlags(data, NM_VPNC_SECRET_FLAG_NONE)
	logicSetSettingVpnVpncKeyXauthPasswordFlags(data, NM_VPNC_SECRET_FLAG_NONE)
	setSettingVpnVpncKeyVendor(data, NM_VPNC_VENDOR_CISCO)
	setSettingVpnVpncKeyPerfectForward(data, NM_VPNC_PFS_SERVER)
	setSettingVpnVpncKeyDhgroup(data, NM_VPNC_DHGROUP_DH2)
	setSettingVpnVpncKeyLocalPort(data, 0)
}

// vpn-vpnc
func getSettingVpnVpncAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionVpnVpnc, NM_SETTING_VPN_VPNC_KEY_GATEWAY)
	keys = appendAvailableKeys(data, keys, sectionVpnVpnc, NM_SETTING_VPN_VPNC_KEY_XAUTH_USER)
	keys = appendAvailableKeys(data, keys, sectionVpnVpnc, NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD_FLAGS)
	if isVpnVpncNeedShowXauthPassword(data) {
		keys = appendAvailableKeys(data, keys, sectionVpnVpnc, NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD)
	}
	keys = appendAvailableKeys(data, keys, sectionVpnVpnc, NM_SETTING_VPN_VPNC_KEY_ID)
	keys = appendAvailableKeys(data, keys, sectionVpnVpnc, NM_SETTING_VPN_VPNC_KEY_SECRET_FLAGS)
	if isVpnVpncNeedShowSecret(data) {
		keys = appendAvailableKeys(data, keys, sectionVpnVpnc, NM_SETTING_VPN_VPNC_KEY_SECRET)
	}
	keys = appendAvailableKeys(data, keys, sectionVpnVpnc, NM_SETTING_VPN_VPNC_KEY_AUTHMODE)
	if getSettingVkVpnVpncKeyHybridAuthmode(data) {
		keys = appendAvailableKeys(data, keys, sectionVpnVpnc, NM_SETTING_VPN_VPNC_KEY_CA_FILE)
	}
	return
}
func getSettingVpnVpncAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD_FLAGS:
		values = availableValuesNMVpncSecretFlag
	case NM_SETTING_VPN_VPNC_KEY_SECRET_FLAGS:
		values = availableValuesNMVpncSecretFlag
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
	ensureFileExists(errs, sectionVpnVpnc, NM_SETTING_VPN_VPNC_KEY_CA_FILE, value,
		".pem", ".crt", ".cer")
}

// vpn-vpnc-advanced
func getSettingVpnVpncAdvancedAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_DOMAIN)
	keys = appendAvailableKeys(data, keys, sectionVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_VENDOR)
	keys = appendAvailableKeys(data, keys, sectionVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_APP_VERSION)
	keys = appendAvailableKeys(data, keys, sectionVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_SINGLE_DES)
	keys = appendAvailableKeys(data, keys, sectionVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_NO_ENCRYPTION)
	keys = appendAvailableKeys(data, keys, sectionVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_NAT_TRAVERSAL_MODE)
	keys = appendAvailableKeys(data, keys, sectionVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_DHGROUP)
	keys = appendAvailableKeys(data, keys, sectionVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_PERFECT_FORWARD)
	keys = appendAvailableKeys(data, keys, sectionVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_LOCAL_PORT)
	keys = appendAvailableKeys(data, keys, sectionVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_DPD_IDLE_TIMEOUT)
	return
}
func getSettingVpnVpncAdvancedAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_VPNC_KEY_VENDOR:
		values = []kvalue{
			kvalue{NM_VPNC_VENDOR_CISCO, Tr("Cisco (default)")},
			kvalue{NM_VPNC_VENDOR_NETSCREEN, Tr("Netscreen")},
		}
	case NM_SETTING_VPN_VPNC_KEY_NAT_TRAVERSAL_MODE:
		values = []kvalue{
			kvalue{NM_VPNC_NATT_MODE_NATT, Tr("NAT-T When Available (default)")},
			kvalue{NM_VPNC_NATT_MODE_NATT_ALWAYS, Tr("NAT-T Always")},
			kvalue{NM_VPNC_NATT_MODE_CISCO, Tr("Cisco UDP")},
			kvalue{NM_VPNC_NATT_MODE_NONE, Tr("Disabled")},
		}
	case NM_SETTING_VPN_VPNC_KEY_DHGROUP:
		values = []kvalue{
			kvalue{NM_VPNC_DHGROUP_DH1, Tr("DH Group 1")},
			kvalue{NM_VPNC_DHGROUP_DH2, Tr("DH Group 2 (default)")},
			kvalue{NM_VPNC_DHGROUP_DH5, Tr("DH Group 5")},
		}
	case NM_SETTING_VPN_VPNC_KEY_PERFECT_FORWARD:
		values = []kvalue{
			kvalue{NM_VPNC_PFS_SERVER, Tr("Server (default)")},
			kvalue{NM_VPNC_PFS_NOPFS, Tr("None")},
			kvalue{NM_VPNC_PFS_DH1, Tr("DH Group 1")},
			kvalue{NM_VPNC_PFS_DH2, Tr("DH Group 2")},
			kvalue{NM_VPNC_PFS_DH5, Tr("DH Group 5")},
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
	case NM_VPNC_SECRET_FLAG_NONE:
		setSettingVpnVpncKeySecretType(data, NM_VPNC_PW_TYPE_SAVE)
	case NM_VPNC_SECRET_FLAG_SAVE:
		setSettingVpnVpncKeySecretType(data, NM_VPNC_PW_TYPE_SAVE)
	case NM_VPNC_SECRET_FLAG_ASK:
		setSettingVpnVpncKeySecretType(data, NM_VPNC_PW_TYPE_ASK)
	case NM_VPNC_SECRET_FLAG_UNUSED:
		setSettingVpnVpncKeySecretType(data, NM_VPNC_PW_TYPE_UNUSED)
	}
	setSettingVpnVpncKeySecretFlags(data, value)
	return
}
func logicSetSettingVpnVpncKeyXauthPasswordFlags(data connectionData, value uint32) (err error) {
	switch value {
	case NM_VPNC_SECRET_FLAG_NONE:
		setSettingVpnVpncKeyXauthPasswordType(data, NM_VPNC_PW_TYPE_SAVE)
	case NM_VPNC_SECRET_FLAG_SAVE:
		setSettingVpnVpncKeyXauthPasswordType(data, NM_VPNC_PW_TYPE_SAVE)
	case NM_VPNC_SECRET_FLAG_ASK:
		setSettingVpnVpncKeyXauthPasswordType(data, NM_VPNC_PW_TYPE_ASK)
	case NM_VPNC_SECRET_FLAG_UNUSED:
		setSettingVpnVpncKeyXauthPasswordType(data, NM_VPNC_PW_TYPE_UNUSED)
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
