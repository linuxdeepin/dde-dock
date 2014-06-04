package network

// TODO doc

const NM_SETTING_CDMA_SETTING_NAME = "cdma"

const (
	NM_SETTING_CDMA_NUMBER         = "number"
	NM_SETTING_CDMA_USERNAME       = "username"
	NM_SETTING_CDMA_PASSWORD       = "password"
	NM_SETTING_CDMA_PASSWORD_FLAGS = "password-flags"
)

func initSettingSectionCdma(data connectionData) {
	setSettingConnectionType(data, NM_SETTING_CDMA_SETTING_NAME)
	addSettingSection(data, sectionCdma)
	setSettingCdmaPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	// TODO: for easy test
	setSettingCdmaUsername(data, "ctnet@mycdma.cn")
	setSettingCdmaPassword(data, "vnet.mobi")
}

// Get available keys
func getSettingCdmaAvailableKeys(data connectionData) (keys []string) {
	// TODO remove
	// keys = append(keys, NM_SETTING_VK_MOBILE_SERVICE_TYPE)
	keys = appendAvailableKeys(data, keys, sectionCdma, NM_SETTING_CDMA_NUMBER)
	keys = appendAvailableKeys(data, keys, sectionCdma, NM_SETTING_CDMA_USERNAME)
	keys = appendAvailableKeys(data, keys, sectionCdma, NM_SETTING_CDMA_PASSWORD)
	return
}

// Get available values
func getSettingCdmaAvailableValues(data connectionData, key string) (values []kvalue) {
	// TODO
	return
}

// Check whether the values are correct
func checkSettingCdmaValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	// TODO
	ensureSettingCdmaNumberNoEmpty(data, errs)
	if isSettingRequireSecret(getSettingCdmaPasswordFlags(data)) {
		ensureSettingCdmaPasswordNoEmpty(data, errs)
	}
	return
}
