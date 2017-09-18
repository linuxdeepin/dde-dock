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

// Get available keys
func getSettingPppAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_REFUSE_EAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_REFUSE_PAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_REFUSE_CHAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_REFUSE_MSCHAP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_REFUSE_MSCHAPV2)

	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_REQUIRE_MPPE)
	if getSettingPppRequireMppe(data) {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_REQUIRE_MPPE_128)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_MPPE_STATEFUL)
	}

	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_NOBSDCOMP)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_NODEFLATE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_NO_VJ_COMP)

	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_LCP_ECHO_FAILURE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPP_SETTING_NAME, nm.NM_SETTING_PPP_LCP_ECHO_INTERVAL)
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
