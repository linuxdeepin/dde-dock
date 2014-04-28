package main

const NM_SETTING_PPPOE_SETTING_NAME = "pppoe"

const (
	NM_SETTING_PPPOE_SERVICE        = "service"
	NM_SETTING_PPPOE_USERNAME       = "username"
	NM_SETTING_PPPOE_PASSWORD       = "password"
	NM_SETTING_PPPOE_PASSWORD_FLAGS = "password-flags"
)

func newPppoeConnection(id, username string) (uuid string) {
	logger.Debugf("new pppoe connection, id=%s", id)
	uuid = newUUID()
	data := newPppoeConnectionData(id, uuid)
	setSettingPppoeUsername(data, username)
	nmAddConnection(data)
	return
}

func newPppoeConnectionData(id, uuid string) (data connectionData) {
	data = make(connectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, typePppoe)

	addSettingField(data, fieldPppoe)

	addSettingField(data, fieldPpp)
	logicSetSettingVkPppEnableLcpEcho(data, true)

	addSettingField(data, fieldIpv4)
	setSettingIp4ConfigMethod(data, NM_SETTING_IP4_CONFIG_METHOD_AUTO)

	return
}

// Get available keys
func getSettingPppoeAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(keys, fieldPppoe, NM_SETTING_PPPOE_SERVICE)
	keys = appendAvailableKeys(keys, fieldPppoe, NM_SETTING_PPPOE_USERNAME)
	keys = appendAvailableKeys(keys, fieldPppoe, NM_SETTING_PPPOE_PASSWORD)
	return
}

// Get available values
func getSettingPppoeAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}

// Check whether the values are correct
func checkSettingPppoeValues(data connectionData) (errs FieldKeyErrors) {
	errs = make(map[string]string)

	// check username
	ensureSettingPppoeUsernameNoEmpty(data, errs)

	// TODO what about if no password?

	return
}
