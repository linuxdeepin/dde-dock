package main

const NM_SETTING_PPPOE_SETTING_NAME = "pppoe"

const (
	NM_SETTING_PPPOE_SERVICE        = "service"
	NM_SETTING_PPPOE_USERNAME       = "username"
	NM_SETTING_PPPOE_PASSWORD       = "password"
	NM_SETTING_PPPOE_PASSWORD_FLAGS = "password-flags"
)

func newPppoeConnection(id string) (uuid string) {
	Logger.Debugf("new pppoe connection, id=%s", id)
	uuid = newUUID()
	data := newPppoeConnectionData(id, uuid)
	nmAddConnection(data)
	return
}

func newPppoeConnectionData(id, uuid string) (data _ConnectionData) {
	data = make(_ConnectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, typePppoe)

	addSettingField(data, fieldPppoe)

	addSettingField(data, fieldPpp)
	setSettingPppLcpEchoFailure(data, 5)
	setSettingPppLcpEchoInterval(data, 30)

	addSettingField(data, fieldIPv4)
	setSettingIp4ConfigMethod(data, NM_SETTING_IP4_CONFIG_METHOD_AUTO)

	return
}

// Get available keys
func getSettingPppoeAvailableKeys(data _ConnectionData) (keys []string) {
	keys = []string{
		NM_SETTING_PPPOE_SERVICE,
		NM_SETTING_PPPOE_USERNAME,
		NM_SETTING_PPPOE_PASSWORD,
	}
	return
}

// Get available values
func getSettingPppoeAvailableValues(key string) (values []string, customizable bool) {
	customizable = true
	return
}

// Check whether the values are correct
func checkSettingPppoeValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)

	// check username
	ensureSettingPppoeUsernameNoEmpty(data, errs)

	return
}
