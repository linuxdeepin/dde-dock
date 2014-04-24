package main

// TODO doc
const NM_SETTING_VPN_SETTING_NAME = "vpn"

const (
	NM_SETTING_VPN_SERVICE_TYPE = "service-type"
	NM_SETTING_VPN_USER_NAME    = "user-name"
	NM_SETTING_VPN_DATA         = "data"
	NM_SETTING_VPN_SECRETS      = "secrets"
)

func getSettingVpnAvailableKeys(data _ConnectionData) (keys []string) { return }
func getSettingVpnAvailableValues(data _ConnectionData, key string) (values []string, customizable bool) {
	return
}
func checkSettingVpnValues(data _ConnectionData) (errs FieldKeyErrors) {
	errs = make(map[string]string)
	return
}
