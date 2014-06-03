package network

const NM_SETTING_PPP_SETTING_NAME = "ppp"

const (
	NM_SETTING_PPP_NOAUTH            = "noauth"
	NM_SETTING_PPP_REFUSE_EAP        = "refuse-eap"
	NM_SETTING_PPP_REFUSE_PAP        = "refuse-pap"
	NM_SETTING_PPP_REFUSE_CHAP       = "refuse-chap"
	NM_SETTING_PPP_REFUSE_MSCHAP     = "refuse-mschap"
	NM_SETTING_PPP_REFUSE_MSCHAPV2   = "refuse-mschapv2"
	NM_SETTING_PPP_NOBSDCOMP         = "nobsdcomp"
	NM_SETTING_PPP_NODEFLATE         = "nodeflate"
	NM_SETTING_PPP_NO_VJ_COMP        = "no-vj-comp"
	NM_SETTING_PPP_REQUIRE_MPPE      = "require-mppe"
	NM_SETTING_PPP_REQUIRE_MPPE_128  = "require-mppe-128"
	NM_SETTING_PPP_MPPE_STATEFUL     = "mppe-stateful"
	NM_SETTING_PPP_CRTSCTS           = "crtscts"
	NM_SETTING_PPP_BAUD              = "baud"
	NM_SETTING_PPP_MRU               = "mru"
	NM_SETTING_PPP_MTU               = "mtu"
	NM_SETTING_PPP_LCP_ECHO_FAILURE  = "lcp-echo-failure"
	NM_SETTING_PPP_LCP_ECHO_INTERVAL = "lcp-echo-interval"
)

// TODO Get available keys
func getSettingPppAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_REFUSE_EAP)
	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_REFUSE_PAP)
	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_REFUSE_CHAP)
	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_REFUSE_MSCHAP)
	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_REFUSE_MSCHAPV2)

	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_REQUIRE_MPPE)
	if getSettingPppRequireMppe(data) {
		keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_REQUIRE_MPPE_128)
		keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_MPPE_STATEFUL)
	}

	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_NOBSDCOMP)
	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_NODEFLATE)
	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_NO_VJ_COMP)

	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_LCP_ECHO_FAILURE)
	keys = appendAvailableKeys(data, keys, sectionPpp, NM_SETTING_PPP_LCP_ECHO_INTERVAL)
	return
}

// Get available values
func getSettingPppAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}

// Check whether the values are correct
func checkSettingPppValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	return
}

// Logic setter
func logicSetSettingPppRequireMppe(data connectionData, value bool) (err error) {
	if !value {
		removeSettingPppRequireMppe128(data)
		removeSettingPppMppeStateful(data)
	}
	setSettingPppRequireMppe(data, value)
	return
}

// Virtual key
func getSettingVkPppEnableLcpEcho(data connectionData) (value bool) {
	if isSettingPppLcpEchoFailureExists(data) && isSettingPppLcpEchoIntervalExists(data) {
		return true
	}
	return false
}
func logicSetSettingVkPppEnableLcpEcho(data connectionData, value bool) (err error) {
	if value {
		setSettingPppLcpEchoFailure(data, 5)
		setSettingPppLcpEchoInterval(data, 30)
	} else {
		removeSettingPppLcpEchoFailure(data)
		removeSettingPppLcpEchoInterval(data)
	}
	return
}
