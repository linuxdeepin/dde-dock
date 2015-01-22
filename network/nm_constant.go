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

// Device types
const (
	NM_DEVICE_TYPE_UNKNOWN    = uint32(0)
	NM_DEVICE_TYPE_ETHERNET   = uint32(1)
	NM_DEVICE_TYPE_WIFI       = uint32(2)
	NM_DEVICE_TYPE_UNUSED1    = uint32(3)
	NM_DEVICE_TYPE_UNUSED2    = uint32(4)
	NM_DEVICE_TYPE_BT         = uint32(5)
	NM_DEVICE_TYPE_OLPC_MESH  = uint32(6)
	NM_DEVICE_TYPE_WIMAX      = uint32(7)
	NM_DEVICE_TYPE_MODEM      = uint32(8)
	NM_DEVICE_TYPE_INFINIBAND = uint32(9)
	NM_DEVICE_TYPE_BOND       = uint32(10)
	NM_DEVICE_TYPE_VLAN       = uint32(11)
	NM_DEVICE_TYPE_ADSL       = uint32(12)
	NM_DEVICE_TYPE_BRIDGE     = uint32(13)
)

// Device states
const (
	NM_DEVICE_STATE_UNKNOWN      = 0   // The device is in an unknown state.
	NM_DEVICE_STATE_UNMANAGED    = 10  // The device is recognized but not managed by NetworkManager.
	NM_DEVICE_STATE_UNAVAILABLE  = 20  // The device cannot be used (carrier off, rfkill, etc).
	NM_DEVICE_STATE_DISCONNECTED = 30  // The device is not connected.
	NM_DEVICE_STATE_PREPARE      = 40  // The device is preparing to connect.
	NM_DEVICE_STATE_CONFIG       = 50  // The device is being configured.
	NM_DEVICE_STATE_NEED_AUTH    = 60  // The device is awaiting secrets necessary to continue connection.
	NM_DEVICE_STATE_IP_CONFIG    = 70  // The IP settings of the device are being requested and configured.
	NM_DEVICE_STATE_IP_CHECK     = 80  // The device's IP connectivity ability is being determined.
	NM_DEVICE_STATE_SECONDARIES  = 90  // The device is waiting for secondary connections to be activated.
	NM_DEVICE_STATE_ACTIVATED    = 100 // The device is active.
	NM_DEVICE_STATE_DEACTIVATING = 110 // The device's network connection is being torn down.
	NM_DEVICE_STATE_FAILED       = 120 // The device is in a failure state following an attempt to activate it.
)

func isDeviceStateManaged(state uint32) bool {
	if state > NM_DEVICE_STATE_UNMANAGED {
		return true
	}
	return false
}
func isDeviceStateAvailable(state uint32) bool {
	if state > NM_DEVICE_STATE_UNAVAILABLE {
		return true
	}
	return false
}
func isDeviceStateActivated(state uint32) bool {
	if state == NM_DEVICE_STATE_ACTIVATED {
		return true
	}
	return false
}
func isDeviceStateInActivating(state uint32) bool {
	if state >= NM_DEVICE_STATE_PREPARE && state <= NM_DEVICE_STATE_ACTIVATED {
		return true
	}
	return false
}

// Device state reasons
const (
	NM_DEVICE_STATE_REASON_UNKNOWN                        = 0
	NM_DEVICE_STATE_REASON_NONE                           = 1
	NM_DEVICE_STATE_REASON_NOW_MANAGED                    = 2
	NM_DEVICE_STATE_REASON_NOW_UNMANAGED                  = 3
	NM_DEVICE_STATE_REASON_CONFIG_FAILED                  = 4
	NM_DEVICE_STATE_REASON_CONFIG_UNAVAILABLE             = 5
	NM_DEVICE_STATE_REASON_CONFIG_EXPIRED                 = 6
	NM_DEVICE_STATE_REASON_NO_SECRETS                     = 7
	NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT          = 8
	NM_DEVICE_STATE_REASON_SUPPLICANT_CONFIG_FAILED       = 9
	NM_DEVICE_STATE_REASON_SUPPLICANT_FAILED              = 10
	NM_DEVICE_STATE_REASON_SUPPLICANT_TIMEOUT             = 11
	NM_DEVICE_STATE_REASON_PPP_START_FAILED               = 12
	NM_DEVICE_STATE_REASON_PPP_DISCONNECT                 = 13
	NM_DEVICE_STATE_REASON_PPP_FAILED                     = 14
	NM_DEVICE_STATE_REASON_DHCP_START_FAILED              = 15
	NM_DEVICE_STATE_REASON_DHCP_ERROR                     = 16
	NM_DEVICE_STATE_REASON_DHCP_FAILED                    = 17
	NM_DEVICE_STATE_REASON_SHARED_START_FAILED            = 18
	NM_DEVICE_STATE_REASON_SHARED_FAILED                  = 19
	NM_DEVICE_STATE_REASON_AUTOIP_START_FAILED            = 20
	NM_DEVICE_STATE_REASON_AUTOIP_ERROR                   = 21
	NM_DEVICE_STATE_REASON_AUTOIP_FAILED                  = 22
	NM_DEVICE_STATE_REASON_MODEM_BUSY                     = 23
	NM_DEVICE_STATE_REASON_MODEM_NO_DIAL_TONE             = 24
	NM_DEVICE_STATE_REASON_MODEM_NO_CARRIER               = 25
	NM_DEVICE_STATE_REASON_MODEM_DIAL_TIMEOUT             = 26
	NM_DEVICE_STATE_REASON_MODEM_DIAL_FAILED              = 27
	NM_DEVICE_STATE_REASON_MODEM_INIT_FAILED              = 28
	NM_DEVICE_STATE_REASON_GSM_APN_FAILED                 = 29
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_NOT_SEARCHING = 30
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_DENIED        = 31
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_TIMEOUT       = 32
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_FAILED        = 33
	NM_DEVICE_STATE_REASON_GSM_PIN_CHECK_FAILED           = 34
	NM_DEVICE_STATE_REASON_FIRMWARE_MISSING               = 35
	NM_DEVICE_STATE_REASON_REMOVED                        = 36
	NM_DEVICE_STATE_REASON_SLEEPING                       = 37
	NM_DEVICE_STATE_REASON_CONNECTION_REMOVED             = 38
	NM_DEVICE_STATE_REASON_USER_REQUESTED                 = 39
	NM_DEVICE_STATE_REASON_CARRIER                        = 40
	NM_DEVICE_STATE_REASON_CONNECTION_ASSUMED             = 41
	NM_DEVICE_STATE_REASON_SUPPLICANT_AVAILABLE           = 42
	NM_DEVICE_STATE_REASON_MODEM_NOT_FOUND                = 43
	NM_DEVICE_STATE_REASON_BT_FAILED                      = 44
	NM_DEVICE_STATE_REASON_GSM_SIM_NOT_INSERTED           = 45
	NM_DEVICE_STATE_REASON_GSM_SIM_PIN_REQUIRED           = 46
	NM_DEVICE_STATE_REASON_GSM_SIM_PUK_REQUIRED           = 47
	NM_DEVICE_STATE_REASON_GSM_SIM_WRONG                  = 48
	NM_DEVICE_STATE_REASON_INFINIBAND_MODE                = 49
	NM_DEVICE_STATE_REASON_DEPENDENCY_FAILED              = 50
	NM_DEVICE_STATE_REASON_BR2684_FAILED                  = 51
	NM_DEVICE_STATE_REASON_MODEM_MANAGER_UNAVAILABLE      = 52
	NM_DEVICE_STATE_REASON_SSID_NOT_FOUND                 = 53
	NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED    = 54
)

func isDeviceStateReasonInvalid(reason uint32) bool {
	switch reason {
	case NM_DEVICE_STATE_REASON_UNKNOWN, NM_DEVICE_STATE_REASON_NONE:
		return true
	}
	return false
}

// custom device state reasons
const (
	GUESS_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED = iota + NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED + 1
	GUESS_NM_DEVICE_STATE_REASON_WIRELESS_DISABLED
	GUESS_NM_DEVICE_STATE_REASON_MODEM_NO_SIGNAL
	GUESS_NM_DEVICE_STATE_REASON_MODEM_WRONG_PLAN
)

// modem capabilities
const (
	NM_DEVICE_MODEM_CAPABILITY_NONE      = 0x00000000
	NM_DEVICE_MODEM_CAPABILITY_POTS      = 0x00000001
	NM_DEVICE_MODEM_CAPABILITY_CDMA_EVDO = 0x00000002
	NM_DEVICE_MODEM_CAPABILITY_GSM_UMTS  = 0x00000004
	NM_DEVICE_MODEM_CAPABILITY_LTE       = 0x00000008
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

// check if vpn connection activating or activated
func isVpnConnectionStateInActivating(state uint32) bool {
	if state >= NM_VPN_CONNECTION_STATE_PREPARE &&
		state <= NM_VPN_CONNECTION_STATE_ACTIVATED {
		return true
	}
	return false
}
func isVpnConnectionStateActivated(state uint32) bool {
	if state == NM_VPN_CONNECTION_STATE_ACTIVATED {
		return true
	}
	return false
}
func isVpnConnectionStateDeactivate(state uint32) bool {
	if state == NM_VPN_CONNECTION_STATE_DISCONNECTE {
		return true
	}
	return false
}
func isVpnConnectionStateFailed(state uint32) bool {
	if state == NM_VPN_CONNECTION_STATE_FAILED {
		return true
	}
	return false
}

// VPN connection state reason
const (
	//don't use iota, the value is defined by networkmanager
	NM_VPN_CONNECTION_STATE_REASON_UNKNOWN               = 0
	NM_VPN_CONNECTION_STATE_REASON_NONE                  = 1
	NM_VPN_CONNECTION_STATE_REASON_USER_DISCONNECTED     = 2
	NM_VPN_CONNECTION_STATE_REASON_DEVICE_DISCONNECTED   = 3
	NM_VPN_CONNECTION_STATE_REASON_SERVICE_STOPPED       = 4
	NM_VPN_CONNECTION_STATE_REASON_IP_CONFIG_INVALID     = 5
	NM_VPN_CONNECTION_STATE_REASON_CONNECT_TIMEOUT       = 6
	NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_TIMEOUT = 7
	NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_FAILED  = 8
	NM_VPN_CONNECTION_STATE_REASON_NO_SECRETS            = 9
	NM_VPN_CONNECTION_STATE_REASON_LOGIN_FAILED          = 10
	NM_VPN_CONNECTION_STATE_REASON_CONNECTION_REMOVED    = 11
)

// Networking states
const (
	NM_STATE_UNKNOWN          = uint32(0)
	NM_STATE_ASLEEP           = uint32(10)
	NM_STATE_DISCONNECTED     = uint32(20)
	NM_STATE_DISCONNECTING    = uint32(30)
	NM_STATE_CONNECTING       = uint32(40)
	NM_STATE_CONNECTED_LOCAL  = uint32(50)
	NM_STATE_CONNECTED_SITE   = uint32(60)
	NM_STATE_CONNECTED_GLOBAL = uint32(70)
)

// Access point secret flags
//https://projects.gnome.org/NetworkManager/developers/api/09/spec.html#type-NM_802_11_AP_SEC
const (
	NM_802_11_AP_SEC_NONE            = uint32(0x0)
	NM_802_11_AP_SEC_PAIR_WEP40      = uint32(0x1)
	NM_802_11_AP_SEC_PAIR_WEP104     = uint32(0x2)
	NM_802_11_AP_SEC_PAIR_TKIP       = uint32(0x4)
	NM_802_11_AP_SEC_PAIR_CCMP       = uint32(0x8)
	NM_802_11_AP_SEC_GROUP_WEP40     = uint32(0x10)
	NM_802_11_AP_SEC_GROUP_WEP104    = uint32(0x20)
	NM_802_11_AP_SEC_GROUP_TKIP      = uint32(0x40)
	NM_802_11_AP_SEC_GROUP_CCMP      = uint32(0x80)
	NM_802_11_AP_SEC_KEY_MGMT_PSK    = uint32(0x100)
	NM_802_11_AP_SEC_KEY_MGMT_802_1X = uint32(0x200)
)
const (
	NM_802_11_AP_FLAGS_NONE    = uint32(0x0)
	NM_802_11_AP_FLAGS_PRIVACY = uint32(0x1)
)

// Agent secret flags
const (
	// No special behavior; by default no user interaction is allowed
	// and requests for secrets are fulfilled from persistent storage,
	// or if no secrets are available an error is returned.
	NM_SECRET_AGENT_GET_SECRETS_FLAG_NONE = 0x0

	// Allows the request to interact with the user, possibly
	// prompting via UI for secrets if any are required, or if none
	// are found in persistent storage.
	NM_SECRET_AGENT_GET_SECRETS_FLAG_ALLOW_INTERACTION = 0x1

	// Explicitly prompt for new secrets from the user. This flag
	// signals that NetworkManager thinks any existing secrets are
	// invalid or wrong. This flag implies that interaction is
	// allowed.
	NM_SECRET_AGENT_GET_SECRETS_FLAG_REQUEST_NEW = 0x2

	// Set if the request was initiated by user-requested action via
	// the D-Bus interface, as opposed to automatically initiated by
	// NetworkManager in response to (for example) scan results or
	// carrier changes.
	NM_SECRET_AGENT_GET_SECRETS_FLAG_USER_REQUESTED = 0x4
)

// Active connection states
const (
	NM_ACTIVE_CONNECTION_STATE_UNKNOWN      = 0
	NM_ACTIVE_CONNECTION_STATE_ACTIVATING   = 1
	NM_ACTIVE_CONNECTION_STATE_ACTIVATED    = 2
	NM_ACTIVE_CONNECTION_STATE_DEACTIVATING = 3
	NM_ACTIVE_CONNECTION_STATE_DEACTIVATE   = 4
)

// check if connection activating or activated
func isConnectionStateInActivating(state uint32) bool {
	if state == NM_ACTIVE_CONNECTION_STATE_ACTIVATING ||
		state == NM_ACTIVE_CONNECTION_STATE_ACTIVATED {
		return true
	}
	return false
}
func isConnectionStateActivated(state uint32) bool {
	if state == NM_ACTIVE_CONNECTION_STATE_ACTIVATED {
		return true
	}
	return false
}
func isConnectionStateInDeactivating(state uint32) bool {
	if state == NM_ACTIVE_CONNECTION_STATE_DEACTIVATING ||
		state == NM_ACTIVE_CONNECTION_STATE_DEACTIVATE {
		return true
	}
	return false
}
func isConnectionStateDeactivate(state uint32) bool {
	if state == NM_ACTIVE_CONNECTION_STATE_DEACTIVATE {
		return true
	}
	return false
}

// Connection secret flags
const (
	NM_SETTING_SECRET_FLAG_NONE         = 0x00000000
	NM_SETTING_SECRET_FLAG_AGENT_OWNED  = 0x00000001
	NM_SETTING_SECRET_FLAG_NOT_SAVED    = 0x00000002
	NM_SETTING_SECRET_FLAG_NOT_REQUIRED = 0x00000004
)

var availableValuesSettingSecretFlags []kvalue

func initAvailableValuesSecretFlags() {
	availableValuesSettingSecretFlags = []kvalue{
		kvalue{NM_SETTING_SECRET_FLAG_NONE, Tr("Saved")}, // system saved
		// kvalue{NM_SETTING_SECRET_FLAG_AGENT_OWNED, Tr("Saved")},
		kvalue{NM_SETTING_SECRET_FLAG_NOT_SAVED, Tr("Always Ask")},
		kvalue{NM_SETTING_SECRET_FLAG_NOT_REQUIRED, Tr("Not Required")},
	}
}

func isSettingRequireSecret(flag uint32) bool {
	if flag == NM_SETTING_SECRET_FLAG_NONE || flag == NM_SETTING_SECRET_FLAG_AGENT_OWNED {
		return true
	}
	return false
}
