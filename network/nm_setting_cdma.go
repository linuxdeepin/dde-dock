/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
