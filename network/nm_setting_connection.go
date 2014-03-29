package main

// TODO doc

const NM_SETTING_CONNECTION_SETTING_NAME = "connection"

const (
	NM_SETTING_CONNECTION_ID          = "id"
	NM_SETTING_CONNECTION_UUID        = "uuid"
	NM_SETTING_CONNECTION_TYPE        = "type"
	NM_SETTING_CONNECTION_AUTOCONNECT = "autoconnect"
	NM_SETTING_CONNECTION_TIMESTAMP   = "timestamp"
	NM_SETTING_CONNECTION_READ_ONLY   = "read-only"
	NM_SETTING_CONNECTION_PERMISSIONS = "permissions"
	NM_SETTING_CONNECTION_ZONE        = "zone"
	NM_SETTING_CONNECTION_MASTER      = "master"
	NM_SETTING_CONNECTION_SLAVE_TYPE  = "slave-type"
	NM_SETTING_CONNECTION_SECONDARIES = "secondaries"
)

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
func generalGetSettingConnectionKey(data _ConnectionData, key string) (value string) {
	switch key {
	default:
		LOGGER.Error("generalGetSettingConnectionKey: invalide key", key)
	case NM_SETTING_CONNECTION_ID:
		value = getSettingConnectionId(data)
	case NM_SETTING_CONNECTION_UUID:
		value = getSettingConnectionUuid(data)
	case NM_SETTING_CONNECTION_TYPE:
		value = getSettingConnectionType(data)
	case NM_SETTING_CONNECTION_AUTOCONNECT:
		value = getSettingConnectionAutoconnect(data)
	case NM_SETTING_CONNECTION_TIMESTAMP:
		value = getSettingConnectionTimestamp(data)
	case NM_SETTING_CONNECTION_READ_ONLY:
		value = getSettingConnectionReadOnly(data)
	case NM_SETTING_CONNECTION_PERMISSIONS:
		value = getSettingConnectionPermissions(data)
	case NM_SETTING_CONNECTION_ZONE:
		value = getSettingConnectionZone(data)
	case NM_SETTING_CONNECTION_MASTER:
		value = getSettingConnectionMaster(data)
	case NM_SETTING_CONNECTION_SLAVE_TYPE:
		value = getSettingConnectionSlaveType(data)
	case NM_SETTING_CONNECTION_SECONDARIES:
		value = getSettingConnectionSecondaries(data)
	}
	return
}

// TODO use logic setter
func generalSetSettingConnectionKey(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSettingConnectionKey: invalide key", key)
	case NM_SETTING_CONNECTION_ID:
		setSettingConnectionId(data, value)
	case NM_SETTING_CONNECTION_UUID:
		setSettingConnectionUuid(data, value)
	case NM_SETTING_CONNECTION_TYPE:
		setSettingConnectionType(data, value)
	case NM_SETTING_CONNECTION_AUTOCONNECT:
		setSettingConnectionAutoconnect(data, value)
	case NM_SETTING_CONNECTION_TIMESTAMP:
		setSettingConnectionTimestamp(data, value)
	case NM_SETTING_CONNECTION_READ_ONLY:
		setSettingConnectionReadOnly(data, value)
	case NM_SETTING_CONNECTION_PERMISSIONS:
		setSettingConnectionPermissions(data, value)
	case NM_SETTING_CONNECTION_ZONE:
		setSettingConnectionZone(data, value)
	case NM_SETTING_CONNECTION_MASTER:
		setSettingConnectionMaster(data, value)
	case NM_SETTING_CONNECTION_SLAVE_TYPE:
		setSettingConnectionSlaveType(data, value)
	case NM_SETTING_CONNECTION_SECONDARIES:
		setSettingConnectionSecondaries(data, value)
	}
	return
}

// Getter
func getSettingConnectionId(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ID, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ID))
	return
}
func getSettingConnectionUuid(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_UUID, getSettingConnectionKeyType(NM_SETTING_CONNECTION_UUID))
	return
}
func getSettingConnectionType(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TYPE, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TYPE))
	return
}
func getSettingConnectionAutoconnect(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_802_1X_SYSTEM_CA_CERTS, getSettingConnectionKeyType(NM_SETTING_CONNECTION_AUTOCONNECT))
	return
}
func getSettingConnectionTimestamp(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TIMESTAMP, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TIMESTAMP))
	return
}
func getSettingConnectionReadOnly(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_READ_ONLY, getSettingConnectionKeyType(NM_SETTING_CONNECTION_READ_ONLY))
	return
}
func getSettingConnectionPermissions(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_PERMISSIONS, getSettingConnectionKeyType(NM_SETTING_CONNECTION_PERMISSIONS))
	return
}
func getSettingConnectionZone(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ZONE, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ZONE))
	return
}
func getSettingConnectionMaster(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_MASTER, getSettingConnectionKeyType(NM_SETTING_CONNECTION_MASTER))
	return
}
func getSettingConnectionSlaveType(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SLAVE_TYPE, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SLAVE_TYPE))
	return
}
func getSettingConnectionSecondaries(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SECONDARIES, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SECONDARIES))
	return
}

// Setter
func setSettingConnectionId(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ID, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ID))
}
func setSettingConnectionUuid(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_UUID, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_UUID))
}
func setSettingConnectionType(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TYPE, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TYPE))
}
func setSettingConnectionAutoconnect(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_AUTOCONNECT, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_AUTOCONNECT))
}
func setSettingConnectionTimestamp(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TIMESTAMP, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TIMESTAMP))
}
func setSettingConnectionReadOnly(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_READ_ONLY, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_READ_ONLY))
}
func setSettingConnectionPermissions(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_PERMISSIONS, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_PERMISSIONS))
}
func setSettingConnectionZone(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ZONE, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ZONE))
}
func setSettingConnectionMaster(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_MASTER, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_MASTER))
}
func setSettingConnectionSlaveType(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SLAVE_TYPE, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SLAVE_TYPE))
}
func setSettingConnectionSecondaries(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SECONDARIES, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SECONDARIES))
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
