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

const NM_SETTING_SERIAL_SETTING_NAME = "serial"

const (
	NM_SETTING_SERIAL_BAUD       = "baud"
	NM_SETTING_SERIAL_BITS       = "bits"
	NM_SETTING_SERIAL_PARITY     = "parity"
	NM_SETTING_SERIAL_STOPBITS   = "stopbits"
	NM_SETTING_SERIAL_SEND_DELAY = "send-delay"
)

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
