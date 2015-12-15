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

import (
	"os/user"
)

// Get available keys
func getSettingConnectionAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionConnection, NM_SETTING_CONNECTION_ID)

	// auto-connect only available for target connection types
	switch getSettingConnectionType(data) {
	case NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_PPPOE_SETTING_NAME, NM_SETTING_GSM_SETTING_NAME, NM_SETTING_CDMA_SETTING_NAME:
		keys = appendAvailableKeys(data, keys, sectionConnection, NM_SETTING_CONNECTION_AUTOCONNECT)
	case NM_SETTING_VPN_SETTING_NAME:
		keys = appendAvailableKeys(data, keys, sectionConnection, NM_SETTING_VK_VPN_AUTOCONNECT)
	}
	return
}

// Get available values
func getSettingConnectionAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}

// Check whether the values are correct
func checkSettingConnectionValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)

	// check id
	ensureSettingConnectionIdNoEmpty(data, errs)

	// if the connection is created manually, ensure the id is unique
	if isCreatedManuallyConnection(data) {
		id := getSettingConnectionId(data)
		uuid := getSettingConnectionUuid(data)
		if isStringInArray(id, nmGetOtherConnectionIds(uuid)) {
			rememberError(errs, sectionConnection, NM_SETTING_CONNECTION_ID, NM_KEY_ERROR_INVALID_VALUE)
		}
	}

	return
}

// Virtual key getter and setter
func getSettingVkConnectionNoPermission(data connectionData) (value bool) {
	permission := getSettingConnectionPermissions(data)
	if len(permission) > 0 {
		return false
	}
	return true
}
func logicSetSettingVkConnectionNoPermission(data connectionData, value bool) (err error) {
	if value {
		removeSettingConnectionPermissions(data)
	} else {
		currentUser, err2 := user.Current()
		if err2 != nil {
			logger.Error(err2)
			return
		}
		permission := "user:" + currentUser.Username + ":"
		setSettingConnectionPermissions(data, []string{permission})
	}
	return
}
