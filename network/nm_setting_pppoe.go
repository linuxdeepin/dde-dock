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

const NM_SETTING_PPPOE_SETTING_NAME = "pppoe"

const (
	NM_SETTING_PPPOE_SERVICE        = "service"
	NM_SETTING_PPPOE_USERNAME       = "username"
	NM_SETTING_PPPOE_PASSWORD       = "password"
	NM_SETTING_PPPOE_PASSWORD_FLAGS = "password-flags"
)

func newPppoeConnection(id, username string) (uuid string) {
	logger.Debugf("new pppoe connection, id=%s", id)
	uuid = genUuid()
	data := newPppoeConnectionData(id, uuid)
	setSettingPppoeUsername(data, username)
	nmAddConnection(data)
	return
}

func newPppoeConnectionData(id, uuid string) (data connectionData) {
	data = make(connectionData)

	addSettingSection(data, sectionConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, NM_SETTING_PPPOE_SETTING_NAME)
	setSettingConnectionAutoconnect(data, false)

	initSettingSectionWired(data)

	addSettingSection(data, sectionPppoe)

	addSettingSection(data, sectionPpp)
	logicSetSettingVkPppEnableLcpEcho(data, true)

	initSettingSectionIpv4(data)
	return
}

// Get available keys
func getSettingPppoeAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionPppoe, NM_SETTING_PPPOE_SERVICE)
	keys = appendAvailableKeys(data, keys, sectionPppoe, NM_SETTING_PPPOE_USERNAME)
	if isSettingRequireSecret(getSettingPppoePasswordFlags(data)) {
		keys = appendAvailableKeys(data, keys, sectionPppoe, NM_SETTING_PPPOE_PASSWORD)
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
