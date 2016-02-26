/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

// Get available keys
func getSettingSerialAvailableKeys(data connectionData) (keys []string) {
	return
}

// Get available values
func getSettingSerialAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}

// Check whether the values are correct
func checkSettingSerialValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	return
}
