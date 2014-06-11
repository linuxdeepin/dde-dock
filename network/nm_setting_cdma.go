/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

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
	setSettingCdmaNumber(data, "#777")
	setSettingCdmaPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	// TODO: for easy test
	setSettingCdmaUsername(data, "ctnet@mycdma.cn")
	setSettingCdmaPassword(data, "vnet.mobi")
}

// Get available keys
func getSettingCdmaAvailableKeys(data connectionData) (keys []string) {
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
