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

// Get key type
func getSettingConnectionKeyType(key string) (t ktype) {
	switch key {
	default:
		t = ktypeUnknown
	case NM_SETTING_CONNECTION_ID:
		t = ktypeString
	case NM_SETTING_CONNECTION_UUID:
		t = ktypeString
	case NM_SETTING_CONNECTION_TYPE:
		t = ktypeString
	case NM_SETTING_CONNECTION_PERMISSIONS:
		t = ktypeArrayString
	case NM_SETTING_CONNECTION_AUTOCONNECT:
		t = ktypeBoolean
	case NM_SETTING_CONNECTION_TIMESTAMP:
		t = ktypeUint64
	case NM_SETTING_CONNECTION_READ_ONLY:
		t = ktypeBoolean
	case NM_SETTING_CONNECTION_ZONE:
		t = ktypeString
	case NM_SETTING_CONNECTION_MASTER:
		t = ktypeString
	case NM_SETTING_CONNECTION_SLAVE_TYPE:
		t = ktypeString
	case NM_SETTING_CONNECTION_SECONDARIES:
		t = ktypeArrayString
	}
	return
}

// Get and set key's value generally
func generalGetSettingConnectionKeyJSON(data _ConnectionData, key string) (value string) {
	switch key {
	default:
		LOGGER.Error("generalGetSettingConnectionKey: invalide key", key)
	case NM_SETTING_CONNECTION_ID:
		value = getSettingConnectionIdJSON(data)
	case NM_SETTING_CONNECTION_UUID:
		value = getSettingConnectionUuidJSON(data)
	case NM_SETTING_CONNECTION_TYPE:
		value = getSettingConnectionTypeJSON(data)
	case NM_SETTING_CONNECTION_AUTOCONNECT:
		value = getSettingConnectionAutoconnectJSON(data)
	case NM_SETTING_CONNECTION_TIMESTAMP:
		value = getSettingConnectionTimestampJSON(data)
	case NM_SETTING_CONNECTION_READ_ONLY:
		value = getSettingConnectionReadOnlyJSON(data)
	case NM_SETTING_CONNECTION_PERMISSIONS:
		value = getSettingConnectionPermissionsJSON(data)
	case NM_SETTING_CONNECTION_ZONE:
		value = getSettingConnectionZoneJSON(data)
	case NM_SETTING_CONNECTION_MASTER:
		value = getSettingConnectionMasterJSON(data)
	case NM_SETTING_CONNECTION_SLAVE_TYPE:
		value = getSettingConnectionSlaveTypeJSON(data)
	case NM_SETTING_CONNECTION_SECONDARIES:
		value = getSettingConnectionSecondariesJSON(data)
	}
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

// TODO tmp
func getSettingConnectionId(data _ConnectionData) (value string) {
	value, _ = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ID).(string)
	return
}
func getSettingConnectionUuid(data _ConnectionData) (value string) {
	value, _ = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_UUID).(string)
	return
}
func getSettingConnectionType(data _ConnectionData) (value string) {
	value, _ = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TYPE).(string)
	return
}
func setSettingConnectionId(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ID, value)
}
func setSettingConnectionUuid(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_UUID, value)
}
func setSettingConnectionType(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TYPE, value)
}

// Getter
func getSettingConnectionIdJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ID, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ID))
	return
}
func getSettingConnectionUuidJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_UUID, getSettingConnectionKeyType(NM_SETTING_CONNECTION_UUID))
	return
}
func getSettingConnectionTypeJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TYPE, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TYPE))
	return
}
func getSettingConnectionAutoconnectJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_802_1X_SYSTEM_CA_CERTS, getSettingConnectionKeyType(NM_SETTING_CONNECTION_AUTOCONNECT))
	return
}
func getSettingConnectionTimestampJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TIMESTAMP, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TIMESTAMP))
	return
}
func getSettingConnectionReadOnlyJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_READ_ONLY, getSettingConnectionKeyType(NM_SETTING_CONNECTION_READ_ONLY))
	return
}
func getSettingConnectionPermissionsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_PERMISSIONS, getSettingConnectionKeyType(NM_SETTING_CONNECTION_PERMISSIONS))
	return
}
func getSettingConnectionZoneJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ZONE, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ZONE))
	return
}
func getSettingConnectionMasterJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_MASTER, getSettingConnectionKeyType(NM_SETTING_CONNECTION_MASTER))
	return
}
func getSettingConnectionSlaveTypeJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SLAVE_TYPE, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SLAVE_TYPE))
	return
}
func getSettingConnectionSecondariesJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SECONDARIES, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SECONDARIES))
	return
}

// Setter
func setSettingConnectionIdJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ID, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ID))
}
func setSettingConnectionUuidJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_UUID, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_UUID))
}
func setSettingConnectionTypeJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TYPE, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TYPE))
}
func setSettingConnectionAutoconnectJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_AUTOCONNECT, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_AUTOCONNECT))
}
func setSettingConnectionTimestampJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TIMESTAMP, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TIMESTAMP))
}
func setSettingConnectionReadOnlyJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_READ_ONLY, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_READ_ONLY))
}
func setSettingConnectionPermissionsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_PERMISSIONS, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_PERMISSIONS))
}
func setSettingConnectionZoneJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ZONE, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ZONE))
}
func setSettingConnectionMasterJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_MASTER, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_MASTER))
}
func setSettingConnectionSlaveTypeJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SLAVE_TYPE, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SLAVE_TYPE))
}
func setSettingConnectionSecondariesJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SECONDARIES, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SECONDARIES))
}

// Remover
func removeSettingConnectionId(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ID)
}
func removeSettingConnectionUuid(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_UUID)
}
func removeSettingConnectionType(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TYPE)
}
func removeSettingConnectionAutoconnect(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_802_1X_SYSTEM_CA_CERTS)
}
func removeSettingConnectionTimestamp(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TIMESTAMP)
}
func removeSettingConnectionReadOnly(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_READ_ONLY)
}
func removeSettingConnectionPermissions(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_PERMISSIONS)
}
func removeSettingConnectionZone(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ZONE)
}
func removeSettingConnectionMaster(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_MASTER)
}
func removeSettingConnectionSlaveType(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SLAVE_TYPE)
}
func removeSettingConnectionSecondaries(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SECONDARIES)
}
