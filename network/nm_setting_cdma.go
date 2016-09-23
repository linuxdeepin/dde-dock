/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"pkg.deepin.io/dde/daemon/network/nm"
)

func initSettingSectionCdma(data connectionData) {
	setSettingConnectionType(data, nm.NM_SETTING_CDMA_SETTING_NAME)
	addSetting(data, nm.NM_SETTING_CDMA_SETTING_NAME)
	setSettingCdmaNumber(data, "#777")
	setSettingCdmaPasswordFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
}

// Get available keys
func getSettingCdmaAvailableKeys(data connectionData) (keys []string) {
	if getSettingVkMobileProvider(data) == mobileProviderValueCustom {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_CDMA_SETTING_NAME, nm.NM_SETTING_CDMA_NUMBER)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_CDMA_SETTING_NAME, nm.NM_SETTING_CDMA_USERNAME)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_CDMA_SETTING_NAME, nm.NM_SETTING_CDMA_PASSWORD)
	}
	return
}

// Get available values
func getSettingCdmaAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}

// Check whether the values are correct
func checkSettingCdmaValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	ensureSettingCdmaNumberNoEmpty(data, errs)
	return
}
