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
	"os/user"
	"pkg.deepin.io/dde/daemon/network/nm"
)

// Get available keys
func getSettingConnectionAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_CONNECTION_SETTING_NAME, nm.NM_SETTING_CONNECTION_ID)

	// auto-connect only available for target connection types
	switch getSettingConnectionType(data) {
	case nm.NM_SETTING_WIRED_SETTING_NAME, nm.NM_SETTING_WIRELESS_SETTING_NAME, nm.NM_SETTING_PPPOE_SETTING_NAME, nm.NM_SETTING_GSM_SETTING_NAME, nm.NM_SETTING_CDMA_SETTING_NAME, nm.NM_SETTING_VPN_SETTING_NAME:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_CONNECTION_SETTING_NAME, nm.NM_SETTING_CONNECTION_AUTOCONNECT)
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
			rememberError(errs, nm.NM_SETTING_CONNECTION_SETTING_NAME, nm.NM_SETTING_CONNECTION_ID, nmKeyErrorInvalidValue)
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

func getSettingVkConnectionAutoconnect(data connectionData) (value bool) {
	if isVpnConnection(data) {
		value = getSettingVkVpnAutoconnect(data)
	} else {
		value = getSettingConnectionAutoconnect(data)
	}
	return
}
func logicSetSettingVkConnectionAutoconnect(data connectionData, value bool) (err error) {
	if isVpnConnection(data) {
		err = logicSetSettingVkVpnAutoconnect(data, value)
	} else {
		setSettingConnectionAutoconnect(data, value)
	}
	return
}
func getSettingDummyAvailableKeys(data connectionData) (keys []string) {
	return nil
}
func getSettingDummyAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingDummyValues(data connectionData) (errs sectionErrors) {
	return
}
func getSettingUserAvailableKeys(data connectionData) (keys []string) {
	return nil
}
func getSettingUserAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}
func checkSettingUserValues(data connectionData) (errs sectionErrors) {
	return
}
