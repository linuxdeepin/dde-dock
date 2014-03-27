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

// Getter
func getSettingConnectionId(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ID, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ID))
	return
}
func getSettingConnectionUuid(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_UUID, getSettingConnectionKeyType(NM_SETTING_CONNECTION_UUID))
	return
}
func getSettingConnectionType(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TYPE, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TYPE))
	return
}
func getSettingConnectionAutoconnect(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_802_1X_SYSTEM_CA_CERTS, getSettingConnectionKeyType(NM_SETTING_CONNECTION_AUTOCONNECT))
	return
}
func getSettingConnectionTimestamp(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TIMESTAMP, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TIMESTAMP))
	return
}
func getSettingConnectionReadOnly(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_READ_ONLY, getSettingConnectionKeyType(NM_SETTING_CONNECTION_READ_ONLY))
	return
}
func getSettingConnectionPermissions(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_PERMISSIONS, getSettingConnectionKeyType(NM_SETTING_CONNECTION_PERMISSIONS))
	return
}
func getSettingConnectionZone(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ZONE, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ZONE))
	return
}
func getSettingConnectionMaster(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_MASTER, getSettingConnectionKeyType(NM_SETTING_CONNECTION_MASTER))
	return
}
func getSettingConnectionSlaveType(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SLAVE_TYPE, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SLAVE_TYPE))
	return
}
func getSettingConnectionSecondaries(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SECONDARIES, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SECONDARIES))
	return
}

// Setter
func setSettingConnectionId(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ID, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ID))
	return
}
func setSettingConnectionUuid(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_UUID, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_UUID))
	return
}
func setSettingConnectionType(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TYPE, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TYPE))
	return
}
func setSettingConnectionAutoconnect(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_AUTOCONNECT, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_AUTOCONNECT))
	return
}
func setSettingConnectionTimestamp(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_TIMESTAMP, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_TIMESTAMP))
	return
}
func setSettingConnectionReadOnly(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_READ_ONLY, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_READ_ONLY))
	return
}
func setSettingConnectionPermissions(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_PERMISSIONS, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_PERMISSIONS))
	return
}
func setSettingConnectionZone(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_ZONE, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_ZONE))
	return
}
func setSettingConnectionMaster(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_MASTER, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_MASTER))
	return
}
func setSettingConnectionSlaveType(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SLAVE_TYPE, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SLAVE_TYPE))
	return
}
func setSettingConnectionSecondaries(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_SECONDARIES, value, getSettingConnectionKeyType(NM_SETTING_CONNECTION_SECONDARIES))
	return
}

// TODO remove

// SettingConnectionId NM_SETTING_CONNECTION_ID
// SettingConnectionUuid NM_SETTING_CONNECTION_UUID
// SettingConnectionType NM_SETTING_CONNECTION_TYPE
// SettingConnectionAutoconnect NM_SETTING_CONNECTION_AUTOCONNECT
// SettingConnectionTimestamp NM_SETTING_CONNECTION_TIMESTAMP
// SettingConnectionReadOnly NM_SETTING_CONNECTION_READ_ONLY
// SettingConnectionPermissions NM_SETTING_CONNECTION_PERMISSIONS
// SettingConnectionZone NM_SETTING_CONNECTION_ZONE
// SettingConnectionMaster NM_SETTING_CONNECTION_MASTER
// SettingConnectionSlaveType NM_SETTING_CONNECTION_SLAVE_TYPE
// SettingConnectionSecondaries NM_SETTING_CONNECTION_SECONDARIES
