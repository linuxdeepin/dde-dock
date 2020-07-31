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
	"fmt"

	"pkg.deepin.io/dde/daemon/network/nm"
)

// Logic setter
func logicSetSetting8021xEap(data connectionData, value []string) (err error) {
	if len(value) == 0 {
		logger.Error("eap value is empty")
		err = fmt.Errorf(nmKeyErrorInvalidValue)
		return
	}
	eap := value[0]
	switch eap {
	case "tls":
		removeSettingKeyBut(data, nm.NM_SETTING_802_1X_SETTING_NAME,
			nm.NM_SETTING_802_1X_EAP,
			nm.NM_SETTING_802_1X_IDENTITY,
			nm.NM_SETTING_802_1X_CLIENT_CERT,
			nm.NM_SETTING_802_1X_CA_CERT,
			nm.NM_SETTING_802_1X_PRIVATE_KEY,
			nm.NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD,
			nm.NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS)
		setSetting8021xPrivateKeyPasswordFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
	case "md5":
		removeSettingKeyBut(data, nm.NM_SETTING_802_1X_SETTING_NAME,
			nm.NM_SETTING_802_1X_EAP,
			nm.NM_SETTING_802_1X_IDENTITY,
			nm.NM_SETTING_802_1X_PASSWORD,
			nm.NM_SETTING_802_1X_PASSWORD_FLAGS)
		setSetting8021xPasswordFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
	case "leap":
		removeSettingKeyBut(data, nm.NM_SETTING_802_1X_SETTING_NAME,
			nm.NM_SETTING_802_1X_EAP,
			nm.NM_SETTING_802_1X_IDENTITY,
			nm.NM_SETTING_802_1X_PASSWORD,
			nm.NM_SETTING_802_1X_PASSWORD_FLAGS)
		setSetting8021xPasswordFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
	case "fast":
		removeSettingKeyBut(data, nm.NM_SETTING_802_1X_SETTING_NAME,
			nm.NM_SETTING_802_1X_EAP,
			nm.NM_SETTING_802_1X_IDENTITY,
			nm.NM_SETTING_802_1X_ANONYMOUS_IDENTITY,
			nm.NM_SETTING_802_1X_PAC_FILE,
			nm.NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING,
			nm.NM_SETTING_802_1X_PHASE2_AUTH,
			nm.NM_SETTING_802_1X_PASSWORD,
			nm.NM_SETTING_802_1X_PASSWORD_FLAGS)
		setSetting8021xPhase1FastProvisioning(data, "1")
		setSetting8021xPhase2Auth(data, "gtc")
		setSetting8021xPasswordFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
	case "ttls":
		removeSettingKeyBut(data, nm.NM_SETTING_802_1X_SETTING_NAME,
			nm.NM_SETTING_802_1X_EAP,
			nm.NM_SETTING_802_1X_IDENTITY,
			nm.NM_SETTING_802_1X_ANONYMOUS_IDENTITY,
			nm.NM_SETTING_802_1X_CA_CERT,
			nm.NM_SETTING_802_1X_PHASE2_AUTH,
			nm.NM_SETTING_802_1X_PASSWORD,
			nm.NM_SETTING_802_1X_PASSWORD_FLAGS)
		setSetting8021xPhase2Auth(data, "pap")
		setSetting8021xPasswordFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
	case "peap":
		removeSettingKeyBut(data, nm.NM_SETTING_802_1X_SETTING_NAME,
			nm.NM_SETTING_802_1X_EAP,
			nm.NM_SETTING_802_1X_IDENTITY,
			nm.NM_SETTING_802_1X_ANONYMOUS_IDENTITY,
			nm.NM_SETTING_802_1X_CA_CERT,
			nm.NM_SETTING_802_1X_PHASE1_PEAPVER,
			nm.NM_SETTING_802_1X_PHASE2_AUTH,
			nm.NM_SETTING_802_1X_PASSWORD,
			nm.NM_SETTING_802_1X_PASSWORD_FLAGS)
		removeSetting8021xPhase1Peapver(data)
		setSetting8021xPhase2Auth(data, "mschapv2")
		setSetting8021xPasswordFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
	}
	setSetting8021xEap(data, value)
	return
}

// Virtual key getter
func getSettingVk8021xEnable(data connectionData) (value bool) {
	return isSettingExists(data, nm.NM_SETTING_802_1X_SETTING_NAME)
}

func getSettingVk8021xEap(data connectionData) (eap string) {
	eaps := getSetting8021xEap(data)
	if len(eaps) == 0 {
		logger.Error("eap value is empty")
		return
	}
	eap = eaps[0]
	return
}
