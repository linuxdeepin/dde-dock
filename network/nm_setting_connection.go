package main

// TODO doc

const NM_SETTING_CONNECTION_SETTING_NAME = "connection"

const (
	// User-readable connection identifier/name. Must be one or more
	// characters and may change over the lifetime of the connection
	// if the user decides to rename it.
	NM_SETTING_CONNECTION_ID = "id"

	// Universally unique connection identifier. Must be in the format
	// '2815492f-7e56-435e-b2e9-246bd7cdc664' (ie, contains only
	// hexadecimal characters and '-'). The UUID should be assigned
	// when the connection is created and never changed as long as the
	// connection still applies to the same network. For example, it
	// should not be changed when the user changes the connection's
	// 'id', but should be recreated when the WiFi SSID, mobile
	// broadband network provider, or the connection type changes.
	NM_SETTING_CONNECTION_UUID = "uuid"

	// Base type of the connection. For hardware-dependent
	// connections, should contain the setting name of the
	// hardware-type specific setting (ie, '802-3-ethernet' or
	// '802-11-wireless' or 'bluetooth', etc), and for non-hardware
	// dependent connections like VPN or otherwise, should contain the
	// setting name of that setting type (ie, 'vpn' or 'bridge', etc).
	NM_SETTING_CONNECTION_TYPE = "type"

	// If TRUE, NetworkManager will activate this connection when its
	// network resources are available. If FALSE, the connection must
	// be manually activated by the user or some other mechanism.
	NM_SETTING_CONNECTION_AUTOCONNECT = "autoconnect"

	// Timestamp (in seconds since the Unix Epoch) that the connection
	// was last successfully activated. Settings services should
	// update the connection timestamp periodically when the
	// connection is active to ensure that an active connection has
	// the latest timestamp.
	NM_SETTING_CONNECTION_TIMESTAMP = "timestamp" // TODO

	// If TRUE, the connection is read-only and cannot be changed by
	// the user or any other mechanism. This is normally set for
	// system connections whose plugin cannot yet write updated
	// connections back out.
	NM_SETTING_CONNECTION_READ_ONLY = "read-only"

	// An array of strings defining what access a given user has to
	// this connection. If this is NULL or empty, all users are
	// allowed to access this connection. Otherwise a user is allowed
	// to access this connection if and only if they are in this
	// array. Each entry is of the form "[type]:[id]:[reserved]", for
	// example: "user:dcbw:blah" At this time only the 'user' [type]
	// is allowed. Any other values are ignored and reserved for
	// future use. [id] is the username that this permission refers
	// to, which may not contain the ':' character. Any [reserved]
	// information (if present) must be ignored and is reserved for
	// future use. All of [type], [id], and [reserved] must be valid
	// UTF-8.
	NM_SETTING_CONNECTION_PERMISSIONS = "permissions"

	// The trust level of a the connection.Free form case-insensitive
	// string (for example "Home", "Work", "Public"). NULL or
	// unspecified zone means the connection will be placed in the
	// default zone as defined by the firewall.
	NM_SETTING_CONNECTION_ZONE = "zone"

	// Interface name of the master device or UUID of the master
	// connection
	NM_SETTING_CONNECTION_MASTER = "master"

	// Setting name describing the type of slave this connection is
	// (ie, 'bond') or NULL if this connection is not a slave.
	NM_SETTING_CONNECTION_SLAVE_TYPE = "slave-type"

	// List of connection UUIDs that should be activated when the base
	// connection itself is activated.
	NM_SETTING_CONNECTION_SECONDARIES = "secondaries"
)

// TODO Get available keys
func getSettingConnectionAvailableKeys(data _ConnectionData) (keys []string) {
	keys = []string{
		NM_SETTING_CONNECTION_ID,
		NM_SETTING_CONNECTION_AUTOCONNECT,
		NM_SETTING_CONNECTION_PERMISSIONS,
	}
	return
}

// TODO Check whether the values are correct
func checkSettingConnectionValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)
	return
}

// TODO use logic setter
func generalSetSettingConnectionKeyJSON(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSettingConnectionKey: invalide key", key)
	case NM_SETTING_CONNECTION_ID:
		setSettingConnectionIdJSON(data, value)
	case NM_SETTING_CONNECTION_UUID:
		setSettingConnectionUuidJSON(data, value)
	case NM_SETTING_CONNECTION_TYPE:
		setSettingConnectionTypeJSON(data, value)
	case NM_SETTING_CONNECTION_AUTOCONNECT:
		setSettingConnectionAutoconnectJSON(data, value)
	case NM_SETTING_CONNECTION_TIMESTAMP:
		setSettingConnectionTimestampJSON(data, value)
	case NM_SETTING_CONNECTION_READ_ONLY:
		setSettingConnectionReadOnlyJSON(data, value)
	case NM_SETTING_CONNECTION_PERMISSIONS:
		setSettingConnectionPermissionsJSON(data, value)
	case NM_SETTING_CONNECTION_ZONE:
		setSettingConnectionZoneJSON(data, value)
	case NM_SETTING_CONNECTION_MASTER:
		setSettingConnectionMasterJSON(data, value)
	case NM_SETTING_CONNECTION_SLAVE_TYPE:
		setSettingConnectionSlaveTypeJSON(data, value)
	case NM_SETTING_CONNECTION_SECONDARIES:
		setSettingConnectionSecondariesJSON(data, value)
	}
	return
}
