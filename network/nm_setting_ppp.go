package main

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
func getSettingPppAvailableKeys(data _ConnectionData) (keys []string) {
	keys = []string{}
	return
}

// TODO Get available values
func getSettingPppAvailableValues(key string) (values []string, customizable bool) {
	customizable = true
	return
}

// TODO Check whether the values are correct
func checkSettingPppValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)
	return
}
