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
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/utils"
)

// Init available values
var availableValues8021xPhasesAuth = make(availableValues)

func initAvailableValues8021x() {
	// 'pap', 'chap', 'mschap', 'mschapv2', 'gtc', 'otp', 'md5', and 'tls'
	availableValues8021xPhasesAuth["pap"] = kvalue{"pap", Tr("PAP")}
	availableValues8021xPhasesAuth["chap"] = kvalue{"chap", Tr("CHAP")}
	availableValues8021xPhasesAuth["mschap"] = kvalue{"mschap", Tr("MSCHAP")}
	availableValues8021xPhasesAuth["mschapv2"] = kvalue{"mschapv2", Tr("MSCHAPV2")}
	availableValues8021xPhasesAuth["gtc"] = kvalue{"gtc", Tr("GTC")}
	availableValues8021xPhasesAuth["otp"] = kvalue{"otp", Tr("OTP")}
	availableValues8021xPhasesAuth["md5"] = kvalue{"md5", Tr("MD5")}
	availableValues8021xPhasesAuth["tls"] = kvalue{"tls", Tr("TLS")}
}

// Get available keys
func getSetting8021xAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_EAP)
	switch getSettingVk8021xEap(data) {
	case "tls":
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_CLIENT_CERT)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_CA_CERT)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PRIVATE_KEY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS)
		if is8021xNeedShowPrivatePassword(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD)
		}
	case "md5":
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD_FLAGS)
		if is8021xNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD)
		}
	case "leap":
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD_FLAGS)
		if is8021xNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD)
		}
	case "fast":
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_ANONYMOUS_IDENTITY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PAC_FILE)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PHASE2_AUTH)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD_FLAGS)
		if is8021xNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD)
		}
	case "ttls":
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_ANONYMOUS_IDENTITY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_CA_CERT)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PHASE2_AUTH)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD_FLAGS)
		if is8021xNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD)
		}
	case "peap":
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_ANONYMOUS_IDENTITY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_CA_CERT)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PHASE1_PEAPVER)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PHASE2_AUTH)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD_FLAGS)
		if is8021xNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD)
		}
	}
	return
}
func is8021xNeedShowPrivatePassword(data connectionData) bool {
	flag := getSetting8021xPrivateKeyPasswordFlags(data)
	if flag == nm.NM_SETTING_SECRET_FLAG_NONE || flag == nm.NM_SETTING_SECRET_FLAG_AGENT_OWNED {
		return true
	}
	return false
}
func is8021xNeedShowPassword(data connectionData) bool {
	flag := getSetting8021xPasswordFlags(data)
	if flag == nm.NM_SETTING_SECRET_FLAG_NONE || flag == nm.NM_SETTING_SECRET_FLAG_AGENT_OWNED {
		return true
	}
	return false
}

// Get available values
func getSetting8021xAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_802_1X_EAP:
		if getCustomConnectionType(data) == connectionWired {
			values = []kvalue{
				kvalue{"tls", Tr("TLS")},
				kvalue{"md5", Tr("MD5")},
				kvalue{"fast", Tr("FAST")},
				kvalue{"ttls", Tr("Tunneled TLS")},
				kvalue{"peap", Tr("Protected EAP")},
			}
		} else {
			values = []kvalue{
				kvalue{"tls", Tr("TLS")},
				kvalue{"leap", Tr("LEAP")},
				kvalue{"fast", Tr("FAST")},
				kvalue{"ttls", Tr("Tunneled TLS")},
				kvalue{"peap", Tr("Protected EAP")},
			}
		}
	case nm.NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING:
		values = []kvalue{
			kvalue{"0", Tr("Disabled")},      // Disabled
			kvalue{"1", Tr("Anonymous")},     // Anonymous, allow unauthenticated provisioning
			kvalue{"2", Tr("Authenticated")}, // Authenticated, allow authenticated provisioning
			kvalue{"3", Tr("Both")},          // Both, allow both authenticated and unauthenticated provisioning
		}
	case nm.NM_SETTING_802_1X_PHASE1_PEAPVER:
		values = []kvalue{
			kvalue{"", Tr("Automatic")}, // auto mode
			kvalue{"0", Tr("Version 0")},
			kvalue{"1", Tr("Version 1")},
		}
	case nm.NM_SETTING_802_1X_PHASE2_AUTH:
		switch getSettingVk8021xEap(data) {
		case "tls":
		case "md5":
		case "leap":
		case "fast":
			values = []kvalue{
				availableValues8021xPhasesAuth["gtc"],
				availableValues8021xPhasesAuth["mschapv2"],
			}
		case "ttls":
			values = []kvalue{
				availableValues8021xPhasesAuth["pap"],
				availableValues8021xPhasesAuth["mschap"],
				availableValues8021xPhasesAuth["mschapv2"],
				availableValues8021xPhasesAuth["chap"],
			}
		case "peap":
			values = []kvalue{
				availableValues8021xPhasesAuth["gtc"],
				availableValues8021xPhasesAuth["md5"],
				availableValues8021xPhasesAuth["mschapv2"],
			}
		}
	case nm.NM_SETTING_802_1X_PASSWORD_FLAGS:
		values = availableValuesSettingSecretFlags
	case nm.NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS:
		values = availableValuesSettingSecretFlags
	}
	return
}

// Check whether the values are correct
func checkSetting8021xValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)

	// check eap
	ensureSetting8021xEapNoEmpty(data, errs)
	switch getSettingVk8021xEap(data) {
	default:
		rememberError(errs, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_EAP, nmKeyErrorInvalidValue)
	case "tls":
		ensureSetting8021xIdentityNoEmpty(data, errs)
		ensureSetting8021xClientCertNoEmpty(data, errs)
		ensureSetting8021xPrivateKeyNoEmpty(data, errs)
		if isSettingRequireSecret(getSetting8021xPrivateKeyPasswordFlags(data)) {
			ensureSetting8021xPrivateKeyPasswordNoEmpty(data, errs)
		}
	case "md5":
		ensureSetting8021xIdentityNoEmpty(data, errs)
		if isSettingRequireSecret(getSetting8021xPasswordFlags(data)) {
			ensureSetting8021xPasswordNoEmpty(data, errs)
		}
	case "leap":
		ensureSetting8021xIdentityNoEmpty(data, errs)
		if isSettingRequireSecret(getSetting8021xPasswordFlags(data)) {
			ensureSetting8021xPasswordNoEmpty(data, errs)
		}
	case "fast":
		ensureSetting8021xPhase2AuthNoEmpty(data, errs)
		ensureSetting8021xIdentityNoEmpty(data, errs)
		if isSettingRequireSecret(getSetting8021xPasswordFlags(data)) {
			ensureSetting8021xPasswordNoEmpty(data, errs)
		}
	case "ttls":
		ensureSetting8021xPhase2AuthNoEmpty(data, errs)
		ensureSetting8021xIdentityNoEmpty(data, errs)
		if isSettingRequireSecret(getSetting8021xPasswordFlags(data)) {
			ensureSetting8021xPasswordNoEmpty(data, errs)
		}
	case "peap":
		ensureSetting8021xPhase2AuthNoEmpty(data, errs)
		ensureSetting8021xIdentityNoEmpty(data, errs)
		if isSettingRequireSecret(getSetting8021xPasswordFlags(data)) {
			ensureSetting8021xPasswordNoEmpty(data, errs)
		}
	}

	// check value of pac file
	checkSetting8021xPacFile(data, errs)

	// check value of client cert
	checkSetting8021xClientCert(data, errs)

	// check value of ca cert
	checkSetting8021xCaCert(data, errs)

	// check value of private key
	checkSetting8021xPrivateKey(data, errs)

	return
}

func checkSetting8021xPacFile(data connectionData, errs sectionErrors) {
	if !isSetting8021xPacFileExists(data) {
		return
	}
	value := getSetting8021xPacFile(data)
	if utils.IsURI(value) {
		rememberError(errs, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PAC_FILE, nmKeyErrorInvalidValue)
		return
	}
	ensureFileExists(errs, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PAC_FILE, value)
}
func checkSetting8021xClientCert(data connectionData, errs sectionErrors) {
	if !isSetting8021xClientCertExists(data) {
		return
	}
	value := getSetting8021xClientCert(data)
	ensureByteArrayUriPathExistsFor8021x(errs, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_CLIENT_CERT, value)
}
func checkSetting8021xCaCert(data connectionData, errs sectionErrors) {
	if !isSetting8021xCaCertExists(data) {
		return
	}
	value := getSetting8021xCaCert(data)
	ensureByteArrayUriPathExistsFor8021x(errs, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_CA_CERT, value)
}
func checkSetting8021xPrivateKey(data connectionData, errs sectionErrors) {
	if !isSetting8021xPrivateKeyExists(data) {
		return
	}
	value := getSetting8021xPrivateKey(data)
	ensureByteArrayUriPathExistsFor8021x(errs, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PRIVATE_KEY, value)
}

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
	if isSettingExists(data, nm.NM_SETTING_802_1X_SETTING_NAME) {
		return true
	}
	return false
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
func getSettingVk8021xPacFile(data connectionData) (value string) {
	value = getSetting8021xPacFile(data)
	return
}
func getSettingVk8021xCaCert(data connectionData) (value string) {
	caCert := getSetting8021xCaCert(data)
	value = toLocalPathFor8021x(byteArrayToStrPath(caCert))
	return
}
func getSettingVk8021xClientCert(data connectionData) (value string) {
	clientCert := getSetting8021xClientCert(data)
	value = toLocalPathFor8021x(byteArrayToStrPath(clientCert))
	return
}
func getSettingVk8021xPrivateKey(data connectionData) (value string) {
	privateKey := getSetting8021xPrivateKey(data)
	value = toLocalPathFor8021x(byteArrayToStrPath(privateKey))
	return
}

// Virtual key logic setter
func logicSetSettingVk8021xPacFile(data connectionData, value string) (err error) {
	setSetting8021xPacFile(data, toLocalPath(value))
	return
}
func logicSetSettingVk8021xCaCert(data connectionData, value string) (err error) {
	setSetting8021xCaCert(data, strToByteArrayPath(toUriPathFor8021x(value)))
	return
}
func logicSetSettingVk8021xClientCert(data connectionData, value string) (err error) {
	setSetting8021xClientCert(data, strToByteArrayPath(toUriPathFor8021x(value)))
	return
}
func logicSetSettingVk8021xPrivateKey(data connectionData, value string) (err error) {
	setSetting8021xPrivateKey(data, strToByteArrayPath(toUriPathFor8021x(value)))
	return
}
