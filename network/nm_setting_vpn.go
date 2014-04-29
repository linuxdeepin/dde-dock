package main

// TODO doc
const NM_SETTING_VPN_SETTING_NAME = "vpn"

const (
	NM_SETTING_VPN_SERVICE_TYPE = "service-type"
	NM_SETTING_VPN_USER_NAME    = "user-name"
	NM_SETTING_VPN_DATA         = "data"
	NM_SETTING_VPN_SECRETS      = "secrets"
)

func newBasicVpnConnectionData(id, uuid, service string) (data connectionData) {
	data = make(connectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, typeVpn)
	setSettingConnectionAutoconnect(data, false)
	logicSetSettingVkConnectionNoPermission(data, false)

	addSettingField(data, fieldVpn)
	setSettingVpnServiceType(data, service)

	addSettingField(data, fieldIpv4)
	setSettingIp4ConfigMethod(data, NM_SETTING_IP4_CONFIG_METHOD_AUTO)

	return
}

func getSettingVpnAvailableKeys(data connectionData) (keys []string) { return }
func getSettingVpnAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingVpnValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	return
}
