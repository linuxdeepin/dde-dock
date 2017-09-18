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
	"pkg.deepin.io/lib/utils"
)

func newPppoeConnection(id, username string) (uuid string) {
	logger.Debugf("new pppoe connection, id=%s", id)
	uuid = utils.GenUuid()
	data := newPppoeConnectionData(id, uuid)
	setSettingPppoeUsername(data, username)
	nmAddConnection(data)
	return
}

func newPppoeConnectionData(id, uuid string) (data connectionData) {
	data = make(connectionData)

	addSetting(data, nm.NM_SETTING_CONNECTION_SETTING_NAME)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, nm.NM_SETTING_PPPOE_SETTING_NAME)
	setSettingConnectionAutoconnect(data, true)

	initSettingSectionWired(data)

	addSetting(data, nm.NM_SETTING_PPPOE_SETTING_NAME)

	addSetting(data, nm.NM_SETTING_PPP_SETTING_NAME)
	logicSetSettingVkPppEnableLcpEcho(data, true)

	initSettingSectionIpv4(data)
	return
}

// Get available keys
func getSettingPppoeAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPPOE_SETTING_NAME, nm.NM_SETTING_PPPOE_SERVICE)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPPOE_SETTING_NAME, nm.NM_SETTING_PPPOE_USERNAME)
	if isSettingRequireSecret(getSettingPppoePasswordFlags(data)) {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_PPPOE_SETTING_NAME, nm.NM_SETTING_PPPOE_PASSWORD)
	}
	return
}

// Get available values
func getSettingPppoeAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}

// Check whether the values are correct
func checkSettingPppoeValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	ensureSettingPppoeUsernameNoEmpty(data, errs)
	if isSettingRequireSecret(getSettingPppoePasswordFlags(data)) {
		ensureSettingPppoePasswordNoEmpty(data, errs)
	}
	return
}
