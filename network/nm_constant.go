package main

const (
	fieldConnection       = "connection"
	fieldWireless         = "802-11-wireless"
	fieldWirelessSecurity = "802-11-wireless-security"
	fieldIPv4             = "ipv4"
	fieldIPv6             = "ipv6"
)

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

const (
	NM_DEVICE_STATE_UNKNOWN      = uint32(0)
	NM_DEVICE_STATE_UNMANAGED    = uint32(10)
	NM_DEVICE_STATE_UNAVAILABLE  = uint32(20)
	NM_DEVICE_STATE_DISCONNECTED = uint32(30)
	NM_DEVICE_STATE_PREPARE      = uint32(40)
	NM_DEVICE_STATE_CONFIG       = uint32(50)
	NM_DEVICE_STATE_NEED_AUTH    = uint32(60)
	NM_DEVICE_STATE_IP_CONFIG    = uint32(70)
	NM_DEVICE_STATE_IP_CHECK     = uint32(80)
	NM_DEVICE_STATE_SECONDARIES  = uint32(90)
	NM_DEVICE_STATE_ACTIVATED    = uint32(100)
	NM_DEVICE_STATE_DEACTIVATING = uint32(110)
	NM_DEVICE_STATE_FAILED       = uint32(120)
)

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
