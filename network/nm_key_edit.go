/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	dbus "pkg.deepin.io/lib/dbus1"
)

func getSettingKey(data connectionData, section, key string) (value interface{}) {
	if !isSettingKeyExists(data, section, key) {
		// if key not exists, return the default value
		return generalGetSettingDefaultValue(section, key)
	}
	return doGetSettingKey(data, section, key)
}

func doGetSettingKey(data connectionData, section, key string) (value interface{}) {
	sectionData, ok := data[section]
	if !ok {
		logger.Errorf("invalid section: data[%s]", section)
		return
	}
	variant, ok := sectionData[key]
	if !ok {
		// not exists, just return nil
		return
	}

	value = variant.Value()
	// only debug for develop
	// logger.Debugf("getSettingKey: data[%s][%s]=%v", section, key, value)
	if isInterfaceNil(value) {
		// variant exists, but the value is nil, so we give an error
		// message
		logger.Errorf("getSettingKey: data[%s][%s] is nil", section, key)
	}

	return
}

func setSettingKey(data connectionData, section, key string, value interface{}) {
	var sectionData map[string]dbus.Variant
	sectionData, ok := data[section]
	if !ok {
		logger.Errorf(`set connection data failed, section "%s" is not exits yet`, section)
		return
	}
	sectionData[key] = dbus.MakeVariant(value)
}

func removeSettingKey(data connectionData, section string, keys ...string) {
	logger.Debugf("removeSettingKey data[%s], %s", section, keys)
	sectionData, ok := data[section]
	if !ok {
		return
	}

	for _, k := range keys {
		delete(sectionData, k)
	}
}

func removeSettingKeyBut(data connectionData, section string, keys ...string) {
	sectionData, ok := data[section]
	if !ok {
		return
	}

	for k := range sectionData {
		if !isStringInArray(k, keys) {
			delete(sectionData, k)
		}
	}
}

func isSettingKeyExists(data connectionData, section, key string) bool {
	sectionData, ok := data[section]
	if !ok {
		return false
	}

	_, ok = sectionData[key]
	return ok
}

func addSetting(data connectionData, setting string) {
	var settingData map[string]dbus.Variant
	_, ok := data[setting]
	if !ok {
		// add setting if not exists
		settingData = make(map[string]dbus.Variant)
		data[setting] = settingData
	}
}

func removeSetting(data connectionData, setting string) {
	_, ok := data[setting]
	if ok {
		// remove setting if exists
		delete(data, setting)
	}
}

func isSettingExists(data connectionData, setting string) bool {
	_, ok := data[setting]
	return ok
}
