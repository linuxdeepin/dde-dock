package main

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

// Getter
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

// Setter
func setSettingConnectionId(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ID, value)
}
func setSettingConnectionUuid(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_UUID, value)
}
func setSettingConnectionType(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TYPE, value)
}

// Getter JSON
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

// Setter JSON
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
