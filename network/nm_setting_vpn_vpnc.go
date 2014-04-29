package main

import (
	"dlib"
)

const (
	NM_DBUS_SERVICE_VPNC   = "org.freedesktop.NetworkManager.vpnc"
	NM_DBUS_INTERFACE_VPNC = "org.freedesktop.NetworkManager.vpnc"
	NM_DBUS_PATH_VPNC      = "/org/freedesktop/NetworkManager/vpnc"
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
	NM_VPNC_SECRET_FLAG_SAVE   = 1
	NM_VPNC_SECRET_FLAG_ASK    = 3
	NM_VPNC_SECRET_FLAG_UNUSED = 5
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

// Initialize available values
var availableValuesNMVpncSecretFlag = []kvalue{
	kvalue{NM_VPNC_SECRET_FLAG_SAVE, dlib.Tr("Saved")},
	kvalue{NM_VPNC_SECRET_FLAG_ASK, dlib.Tr("Always Ask")},
	kvalue{NM_VPNC_SECRET_FLAG_UNUSED, dlib.Tr("Not Required")},
}

func newVpnVpncConnectionData(id, uuid string) (data connectionData) {
	data = newBasicVpnConnectionData(id, uuid, NM_DBUS_SERVICE_VPNC)

	setSettingVpnVpncKeyNatTraversalMode(data, NM_VPNC_NATT_MODE_NATT)
	logicSetSettingVpnVpncKeySecretFlags(data, NM_VPNC_SECRET_FLAG_ASK)
	logicSetSettingVpnVpncKeyXauthPasswordFlags(data, NM_VPNC_SECRET_FLAG_ASK)
	setSettingVpnVpncKeyVendor(data, NM_VPNC_VENDOR_CISCO)
	setSettingVpnVpncKeyPerfectForward(data, NM_VPNC_PFS_SERVER)
	setSettingVpnVpncKeyDhgroup(data, NM_VPNC_DHGROUP_DH2)
	setSettingVpnVpncKeyLocalPort(data, 0)

	return
}

// vpn-vpnc
func getSettingVpnVpncAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnVpnc, NM_SETTING_VPN_VPNC_KEY_GATEWAY)
	keys = appendAvailableKeys(keys, fieldVpnVpnc, NM_SETTING_VPN_VPNC_KEY_XAUTH_USER)
	keys = appendAvailableKeys(keys, fieldVpnVpnc, NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD_FLAGS)
	if getSettingVpnVpncKeyXauthPasswordFlags(data) == NM_VPNC_SECRET_FLAG_SAVE {
		keys = appendAvailableKeys(keys, fieldVpnVpnc, NM_SETTING_VPN_VPNC_KEY_XAUTH_PASSWORD)
	}
	keys = appendAvailableKeys(keys, fieldVpnVpnc, NM_SETTING_VPN_VPNC_KEY_ID)
	keys = appendAvailableKeys(keys, fieldVpnVpnc, NM_SETTING_VPN_VPNC_KEY_SECRET_FLAGS)
	if getSettingVpnVpncKeySecretFlags(data) == NM_VPNC_SECRET_FLAG_SAVE {
		keys = appendAvailableKeys(keys, fieldVpnVpnc, NM_SETTING_VPN_VPNC_KEY_SECRET)
	}
	keys = appendAvailableKeys(keys, fieldVpnVpnc, NM_SETTING_VPN_VPNC_KEY_AUTHMODE)
	if getSettingVkVpnVpncKeyHybridAuthmode(data) {
		keys = appendAvailableKeys(keys, fieldVpnVpnc, NM_SETTING_VPN_VPNC_KEY_CA_FILE)
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
func checkSettingVpnVpncValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	ensureSettingVpnVpncKeyGatewayNoEmpty(data, errs)
	ensureSettingVpnVpncKeyIdNoEmpty(data, errs)
	// TODO
	return
}

// vpn-vpnc-advanced
func getSettingVpnVpncAdvancedAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_DOMAIN)
	keys = appendAvailableKeys(keys, fieldVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_VENDOR)
	keys = appendAvailableKeys(keys, fieldVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_APP_VERSION)
	keys = appendAvailableKeys(keys, fieldVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_SINGLE_DES)
	keys = appendAvailableKeys(keys, fieldVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_NAT_TRAVERSAL_MODE)
	keys = appendAvailableKeys(keys, fieldVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_DHGROUP)
	keys = appendAvailableKeys(keys, fieldVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_PERFECT_FORWARD)
	keys = appendAvailableKeys(keys, fieldVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_LOCAL_PORT)
	keys = appendAvailableKeys(keys, fieldVpnVpncAdvanced, NM_SETTING_VPN_VPNC_KEY_DPD_IDLE_TIMEOUT)
	return
}
func getSettingVpnVpncAdvancedAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_VPN_VPNC_KEY_VENDOR:
		values = []kvalue{
			kvalue{NM_VPNC_VENDOR_CISCO, dlib.Tr("Cisco (default)")},
			kvalue{NM_VPNC_VENDOR_NETSCREEN, dlib.Tr("Netscreen")},
		}
	case NM_SETTING_VPN_VPNC_KEY_NAT_TRAVERSAL_MODE:
		values = []kvalue{
			kvalue{NM_VPNC_NATT_MODE_NATT, dlib.Tr("NAT-T When Available (default)")},
			kvalue{NM_VPNC_NATT_MODE_NATT_ALWAYS, dlib.Tr("NAT-T Always")},
			kvalue{NM_VPNC_NATT_MODE_CISCO, dlib.Tr("Cisco UDP")},
			kvalue{NM_VPNC_NATT_MODE_NONE, dlib.Tr("Disabled")},
		}
	case NM_SETTING_VPN_VPNC_KEY_DHGROUP:
		values = []kvalue{
			kvalue{NM_VPNC_DHGROUP_DH1, dlib.Tr("DH Group 1")},
			kvalue{NM_VPNC_DHGROUP_DH2, dlib.Tr("DH Group 2 (default)")},
			kvalue{NM_VPNC_DHGROUP_DH5, dlib.Tr("DH Group 5")},
		}
	case NM_SETTING_VPN_VPNC_KEY_PERFECT_FORWARD:
		values = []kvalue{
			kvalue{NM_VPNC_PFS_SERVER, dlib.Tr("Server (default)")},
			kvalue{NM_VPNC_PFS_NOPFS, dlib.Tr("None")},
			kvalue{NM_VPNC_PFS_DH1, dlib.Tr("DH Group 1")},
			kvalue{NM_VPNC_PFS_DH2, dlib.Tr("DH Group 2")},
			kvalue{NM_VPNC_PFS_DH5, dlib.Tr("DH Group 5")},
		}
	}
	return
}
func checkSettingVpnVpncAdvancedValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	// TODO
	return
}

// Logic setter
func logicSetSettingVpnVpncKeySecretFlags(data connectionData, value uint32) (ok bool, errMsg string) {
	ok = true
	switch value {
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
func logicSetSettingVpnVpncKeyXauthPasswordFlags(data connectionData, value uint32) (ok bool, errMsg string) {
	ok = true
	switch value {
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
