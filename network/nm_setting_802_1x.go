package main

import (
	. "dlib/gettext"
	"fmt"
)

// TODO doc

const NM_SETTING_802_1X_SETTING_NAME = "802-1x"

const (
	NM_SETTING_802_1X_EAP                               = "eap"
	NM_SETTING_802_1X_IDENTITY                          = "identity"
	NM_SETTING_802_1X_ANONYMOUS_IDENTITY                = "anonymous-identity"
	NM_SETTING_802_1X_PAC_FILE                          = "pac-file"
	NM_SETTING_802_1X_CA_CERT                           = "ca-cert"
	NM_SETTING_802_1X_CA_PATH                           = "ca-path"
	NM_SETTING_802_1X_SUBJECT_MATCH                     = "subject-match"
	NM_SETTING_802_1X_ALTSUBJECT_MATCHES                = "altsubject-matches"
	NM_SETTING_802_1X_CLIENT_CERT                       = "client-cert"
	NM_SETTING_802_1X_PHASE1_PEAPVER                    = "phase1-peapver"
	NM_SETTING_802_1X_PHASE1_PEAPLABEL                  = "phase1-peaplabel"
	NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING          = "phase1-fast-provisioning"
	NM_SETTING_802_1X_PHASE2_AUTH                       = "phase2-auth"
	NM_SETTING_802_1X_PHASE2_AUTHEAP                    = "phase2-autheap"
	NM_SETTING_802_1X_PHASE2_CA_CERT                    = "phase2-ca-cert"
	NM_SETTING_802_1X_PHASE2_CA_PATH                    = "phase2-ca-path"
	NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH              = "phase2-subject-match"
	NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES         = "phase2-altsubject-matches"
	NM_SETTING_802_1X_PHASE2_CLIENT_CERT                = "phase2-client-cert"
	NM_SETTING_802_1X_PASSWORD                          = "password"
	NM_SETTING_802_1X_PASSWORD_FLAGS                    = "password-flags"
	NM_SETTING_802_1X_PASSWORD_RAW                      = "password-raw"
	NM_SETTING_802_1X_PASSWORD_RAW_FLAGS                = "password-raw-flags"
	NM_SETTING_802_1X_PRIVATE_KEY                       = "private-key"
	NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD              = "private-key-password"
	NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS        = "private-key-password-flags"
	NM_SETTING_802_1X_PHASE2_PRIVATE_KEY                = "phase2-private-key"
	NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD       = "phase2-private-key-password"
	NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS = "phase2-private-key-password-flags"
	NM_SETTING_802_1X_PIN                               = "pin"
	NM_SETTING_802_1X_PIN_FLAGS                         = "pin-flags"
	NM_SETTING_802_1X_SYSTEM_CA_CERTS                   = "system-ca-certs"
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
	switch getCustomConnectionType(data) {
	case connectionWired:
		keys = []string{NM_SETTING_VK_802_1X_ENABLE}
		if !isSettingSectionExists(data, section8021x) {
			return
		}
	}
	keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_EAP)
	switch getSettingVk8021xEap(data) {
	case "tls":
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_CLIENT_CERT)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_CA_CERT)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PRIVATE_KEY)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS)
		if is8021xNeedShowPrivatePassword(data) {
			keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD)
		}
	case "md5":
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PASSWORD_FLAGS)
		if is8021xNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PASSWORD)
		}
	case "leap":
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PASSWORD_FLAGS)
		if is8021xNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PASSWORD)
		}
	case "fast":
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_ANONYMOUS_IDENTITY)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PAC_FILE)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PHASE2_AUTH)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PASSWORD_FLAGS)
		if is8021xNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PASSWORD)
		}
	case "ttls":
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_ANONYMOUS_IDENTITY)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_CA_CERT)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PHASE2_AUTH)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PASSWORD_FLAGS)
		if is8021xNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PASSWORD)
		}
	case "peap":
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_ANONYMOUS_IDENTITY)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_CA_CERT)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PHASE1_PEAPVER)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PHASE2_AUTH)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_IDENTITY)
		keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PASSWORD_FLAGS)
		if is8021xNeedShowPassword(data) {
			keys = appendAvailableKeys(data, keys, section8021x, NM_SETTING_802_1X_PASSWORD)
		}
	}
	return
}
func is8021xNeedShowPrivatePassword(data connectionData) bool {
	flag := getSetting8021xPrivateKeyPasswordFlags(data)
	if flag == NM_SETTING_SECRET_FLAG_NONE || flag == NM_SETTING_SECRET_FLAG_AGENT_OWNED {
		return true
	}
	return false
}
func is8021xNeedShowPassword(data connectionData) bool {
	flag := getSetting8021xPasswordFlags(data)
	if flag == NM_SETTING_SECRET_FLAG_NONE || flag == NM_SETTING_SECRET_FLAG_AGENT_OWNED {
		return true
	}
	return false
}

// Get available values
func getSetting8021xAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_802_1X_EAP:
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
	case NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING:
		values = []kvalue{
			kvalue{"0", Tr("Disabled")},      // Disabled
			kvalue{"1", Tr("Anonymous")},     // Anonymous, allow unauthenticated provisioning
			kvalue{"2", Tr("Authenticated")}, // Authenticated, allow authenticated provisioning
			kvalue{"3", Tr("Both")},          // Both, allow both authenticated and unauthenticated provisioning
		}
	case NM_SETTING_802_1X_PHASE1_PEAPVER:
		values = []kvalue{
			kvalue{"", Tr("Automatic")}, // auto mode
			kvalue{"0", Tr("Version 0")},
			kvalue{"1", Tr("Version 1")},
		}
	case NM_SETTING_802_1X_PHASE2_AUTH:
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
	case NM_SETTING_802_1X_PASSWORD_FLAGS:
		values = availableValuesNMSettingSecretFlag
	case NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS:
		values = availableValuesNMSettingSecretFlag
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
		rememberError(errs, section8021x, NM_SETTING_802_1X_EAP, NM_KEY_ERROR_INVALID_VALUE)
	case "tls":
		ensureSetting8021xIdentityNoEmpty(data, errs)
		ensureSetting8021xClientCertNoEmpty(data, errs)
		ensureSetting8021xCaCertNoEmpty(data, errs)
		ensureSetting8021xPrivateKeyNoEmpty(data, errs)
		ensureSetting8021xPasswordNoEmpty(data, errs)
		ensureSetting8021xSystemCaCertsNoEmpty(data, errs)
	case "md5":
		ensureSetting8021xIdentityNoEmpty(data, errs)
		ensureSetting8021xPasswordNoEmpty(data, errs)
	case "leap":
		ensureSetting8021xIdentityNoEmpty(data, errs)
		ensureSetting8021xPasswordNoEmpty(data, errs)
	case "fast":
		ensureSetting8021xPhase2AuthNoEmpty(data, errs)
		ensureSetting8021xIdentityNoEmpty(data, errs)
		ensureSetting8021xPasswordNoEmpty(data, errs)
		ensureSetting8021xSystemCaCertsNoEmpty(data, errs)
	case "ttls":
		ensureSetting8021xPhase2AuthNoEmpty(data, errs)
		ensureSetting8021xIdentityNoEmpty(data, errs)
		ensureSetting8021xPasswordNoEmpty(data, errs)
		ensureSetting8021xSystemCaCertsNoEmpty(data, errs)
	case "peap":
		ensureSetting8021xPhase2AuthNoEmpty(data, errs)
		ensureSetting8021xIdentityNoEmpty(data, errs)
		ensureSetting8021xPasswordNoEmpty(data, errs)
		ensureSetting8021xSystemCaCertsNoEmpty(data, errs)
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
	if !isLocalPath(value) {
		rememberError(errs, section8021x, NM_SETTING_802_1X_PAC_FILE, NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	ensureFileExists(errs, section8021x, NM_SETTING_802_1X_PAC_FILE, value, ".pac")
}
func checkSetting8021xClientCert(data connectionData, errs sectionErrors) {
	if !isSetting8021xClientCertExists(data) {
		return
	}
	value := getSetting8021xClientCert(data)
	ensureByteArrayUriPathExists(errs, section8021x, NM_SETTING_802_1X_CLIENT_CERT, value,
		".der", ".pem", ".crt", ".cer")
}
func checkSetting8021xCaCert(data connectionData, errs sectionErrors) {
	if !isSetting8021xCaCertExists(data) {
		return
	}
	value := getSetting8021xCaCert(data)
	ensureByteArrayUriPathExists(errs, section8021x, NM_SETTING_802_1X_CA_CERT, value,
		".der", ".pem", ".crt", ".cer")
}
func checkSetting8021xPrivateKey(data connectionData, errs sectionErrors) {
	if !isSetting8021xPrivateKeyExists(data) {
		return
	}
	value := getSetting8021xPrivateKey(data)
	ensureByteArrayUriPathExists(errs, section8021x, NM_SETTING_802_1X_PRIVATE_KEY, value,
		".der", ".pem", ".p12", ".key")
}

// Logic setter
func logicSetSetting8021xEap(data connectionData, value []string) (err error) {
	if len(value) == 0 {
		logger.Error("eap value is empty")
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	eap := value[0]
	switch eap {
	case "tls":
		removeSettingKeyBut(data, section8021x,
			NM_SETTING_802_1X_EAP,
			NM_SETTING_802_1X_IDENTITY,
			NM_SETTING_802_1X_CLIENT_CERT,
			NM_SETTING_802_1X_CA_CERT,
			NM_SETTING_802_1X_PRIVATE_KEY,
			NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD,
			NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS,
			NM_SETTING_802_1X_SYSTEM_CA_CERTS)
		setSetting8021xSystemCaCerts(data, true)
		setSetting8021xPrivateKeyPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	case "md5":
		removeSettingKeyBut(data, section8021x,
			NM_SETTING_802_1X_EAP,
			NM_SETTING_802_1X_IDENTITY,
			NM_SETTING_802_1X_PASSWORD,
			NM_SETTING_802_1X_PASSWORD_FLAGS,
			NM_SETTING_802_1X_SYSTEM_CA_CERTS)
		setSetting8021xSystemCaCerts(data, true)
		setSetting8021xPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	case "leap":
		removeSettingKeyBut(data, section8021x,
			NM_SETTING_802_1X_EAP,
			NM_SETTING_802_1X_IDENTITY,
			NM_SETTING_802_1X_PASSWORD,
			NM_SETTING_802_1X_PASSWORD_FLAGS,
			NM_SETTING_802_1X_SYSTEM_CA_CERTS)
		setSetting8021xSystemCaCerts(data, true)
		setSetting8021xPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	case "fast":
		removeSettingKeyBut(data, section8021x,
			NM_SETTING_802_1X_EAP,
			NM_SETTING_802_1X_IDENTITY,
			NM_SETTING_802_1X_ANONYMOUS_IDENTITY,
			NM_SETTING_802_1X_PAC_FILE,
			NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING,
			NM_SETTING_802_1X_PHASE2_AUTH,
			NM_SETTING_802_1X_PASSWORD,
			NM_SETTING_802_1X_PASSWORD_FLAGS,
			NM_SETTING_802_1X_SYSTEM_CA_CERTS)
		setSetting8021xPhase1FastProvisioning(data, "1")
		setSetting8021xPhase2Auth(data, "gtc")
		setSetting8021xSystemCaCerts(data, true)
		setSetting8021xPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	case "ttls":
		removeSettingKeyBut(data, section8021x,
			NM_SETTING_802_1X_EAP,
			NM_SETTING_802_1X_IDENTITY,
			NM_SETTING_802_1X_ANONYMOUS_IDENTITY,
			NM_SETTING_802_1X_CA_CERT,
			NM_SETTING_802_1X_PHASE2_AUTH,
			NM_SETTING_802_1X_PASSWORD,
			NM_SETTING_802_1X_PASSWORD_FLAGS,
			NM_SETTING_802_1X_SYSTEM_CA_CERTS)
		setSetting8021xPhase2Auth(data, "pap")
		setSetting8021xSystemCaCerts(data, true)
		setSetting8021xPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	case "peap":
		removeSettingKeyBut(data, section8021x,
			NM_SETTING_802_1X_EAP,
			NM_SETTING_802_1X_IDENTITY,
			NM_SETTING_802_1X_ANONYMOUS_IDENTITY,
			NM_SETTING_802_1X_CA_CERT,
			NM_SETTING_802_1X_PHASE1_PEAPVER,
			NM_SETTING_802_1X_PHASE2_AUTH,
			NM_SETTING_802_1X_PASSWORD,
			NM_SETTING_802_1X_PASSWORD_FLAGS,
			NM_SETTING_802_1X_SYSTEM_CA_CERTS)
		setSetting8021xPhase1Peapver(data, "")
		setSetting8021xPhase2Auth(data, "mschapv2")
		setSetting8021xSystemCaCerts(data, true)
		setSetting8021xPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	}
	setSetting8021xEap(data, value)
	return
}

// Virtual key getter
func getSettingVk8021xEnable(data connectionData) (value bool) {
	if isSettingSectionExists(data, section8021x) {
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
	pacFile := getSetting8021xPacFile(data)
	if len(pacFile) > 0 {
		value = toUriPath(pacFile)
	}
	return
}
func getSettingVk8021xCaCert(data connectionData) (value string) {
	caCert := getSetting8021xCaCert(data)
	value = byteArrayToStrPath(caCert)
	return
}
func getSettingVk8021xClientCert(data connectionData) (value string) {
	clientCert := getSetting8021xClientCert(data)
	value = byteArrayToStrPath(clientCert)
	return
}
func getSettingVk8021xPrivateKey(data connectionData) (value string) {
	privateKey := getSetting8021xPrivateKey(data)
	value = byteArrayToStrPath(privateKey)
	return
}

// Virtual key logic setter
func logicSetSettingVk8021xEnable(data connectionData, value bool) (err error) {
	if value {
		addSettingSection(data, section8021x)
		err = logicSetSettingVk8021xEap(data, "tls")
	} else {
		removeSettingSection(data, section8021x)
	}
	return
}
func logicSetSettingVk8021xEap(data connectionData, value string) (err error) {
	return logicSetSetting8021xEap(data, []string{value})
}
func logicSetSettingVk8021xPacFile(data connectionData, value string) (err error) {
	setSetting8021xPacFile(data, toLocalPath(value))
	return
}
func logicSetSettingVk8021xCaCert(data connectionData, value string) (err error) {
	setSetting8021xCaCert(data, strToByteArrayPath(toUriPath(value)))
	return
}
func logicSetSettingVk8021xClientCert(data connectionData, value string) (err error) {
	setSetting8021xClientCert(data, strToByteArrayPath(toUriPath(value)))
	return
}
func logicSetSettingVk8021xPrivateKey(data connectionData, value string) (err error) {
	setSetting8021xPrivateKey(data, strToByteArrayPath(toUriPath(value)))
	return
}
