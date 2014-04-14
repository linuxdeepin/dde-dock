package main

const NM_SETTING_PPPOE_SETTING_NAME = "pppoe"

const (
	NM_SETTING_PPPOE_SERVICE        = "service"
	NM_SETTING_PPPOE_USERNAME       = "username"
	NM_SETTING_PPPOE_PASSWORD       = "password"
	NM_SETTING_PPPOE_PASSWORD_FLAGS = "password-flags"
)

// TODO Get available keys
func getSettingPppoeAvailableKeys(data _ConnectionData) (keys []string) {
	keys = []string{
	// NM_SETTING_CONNECTION_ID,
	// NM_SETTING_CONNECTION_AUTOCONNECT,
	}
	return
}

// TODO Get available values
func getSettingPppoeAvailableValues(key string) (values []string, customizable bool) {
	customizable = true
	return
}

// TODO Check whether the values are correct
func checkSettingPppoeValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)

	// TODO check id
	// ensureSettingConnectionIdNoEmpty(data, errs)

	return
}
