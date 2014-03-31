package main

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

// TODO Get available keys
func getSettingConnectionAvailableKeys(data _ConnectionData) (keys []string) {
	keys = []string{
		NM_SETTING_802_1X_EAP,
		NM_SETTING_802_1X_IDENTITY,
		NM_SETTING_802_1X_ANONYMOUS_IDENTITY,
		NM_SETTING_802_1X_PAC_FILE,
	}
	return
}

// TODO Check whether the values are correct
func checkSetting8021xValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)
	return
}

// Get and set key's value generally
func generalGetSetting8021xKeyJSON(data _ConnectionData, key string) (value string) {
	switch key {
	default:
		LOGGER.Error("generalGetSetting8021xKey: invalide key", key)
	case NM_SETTING_802_1X_EAP:
		value = getSetting8021xEapJSON(data)
	case NM_SETTING_802_1X_IDENTITY:
		value = getSetting8021xIdentityJSON(data)
	case NM_SETTING_802_1X_ANONYMOUS_IDENTITY:
		value = getSetting8021xAnonymousIdentityJSON(data)
	case NM_SETTING_802_1X_PAC_FILE:
		value = getSetting8021xPacFileJSON(data)
	case NM_SETTING_802_1X_CA_CERT:
		value = getSetting8021xCaCertJSON(data)
	case NM_SETTING_802_1X_CA_PATH:
		value = getSetting8021xCaPathJSON(data)
	case NM_SETTING_802_1X_SUBJECT_MATCH:
		value = getSetting8021xSubjectMatchJSON(data)
	case NM_SETTING_802_1X_ALTSUBJECT_MATCHES:
		value = getSetting8021xAltsubjectMatchesJSON(data)
	case NM_SETTING_802_1X_CLIENT_CERT:
		value = getSetting8021xClientCertJSON(data)
	case NM_SETTING_802_1X_PHASE1_PEAPVER:
		value = getSetting8021xPhase1PeapverJSON(data)
	case NM_SETTING_802_1X_PHASE1_PEAPLABEL:
		value = getSetting8021xPhase1PeaplabelJSON(data)
	case NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING:
		value = getSetting8021xPhase1FastProvisioningJSON(data)
	case NM_SETTING_802_1X_PHASE2_AUTH:
		value = getSetting8021xPhase2AuthJSON(data)
	case NM_SETTING_802_1X_PHASE2_AUTHEAP:
		value = getSetting8021xPhase2AutheapJSON(data)
	case NM_SETTING_802_1X_PHASE2_CA_CERT:
		value = getSetting8021xPhase2CaCertJSON(data)
	case NM_SETTING_802_1X_PHASE2_CA_PATH:
		value = getSetting8021xPhase2CaPathJSON(data)
	case NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH:
		value = getSetting8021xPhase2SubjectMatchJSON(data)
	case NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES:
		value = getSetting8021xPhase2AltsubjectMatchesJSON(data)
	case NM_SETTING_802_1X_PASSWORD:
		value = getSetting8021xPasswordJSON(data)
	case NM_SETTING_802_1X_PASSWORD_FLAGS:
		value = getSetting8021xPasswordFlagsJSON(data)
	case NM_SETTING_802_1X_PASSWORD_RAW:
		value = getSetting8021xPasswordRawJSON(data)
	case NM_SETTING_802_1X_PASSWORD_RAW_FLAGS:
		value = getSetting8021xPasswordRawFlagsJSON(data)
	case NM_SETTING_802_1X_PRIVATE_KEY:
		value = getSetting8021xPrivateKeyJSON(data)
	case NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD:
		value = getSetting8021xPrivateKeyPasswordJSON(data)
	case NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS:
		value = getSetting8021xPrivateKeyPasswordFlagsJSON(data)
	case NM_SETTING_802_1X_PHASE2_PRIVATE_KEY:
		value = getSetting8021xPhase2PrivateKeyJSON(data)
	case NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD:
		value = getSetting8021xPhase2PrivateKeyPasswordJSON(data)
	case NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS:
		value = getSetting8021xPhase2PrivateKeyPasswordFlagsJSON(data)
	case NM_SETTING_802_1X_PIN:
		value = getSetting8021xPinJSON(data)
	case NM_SETTING_802_1X_PIN_FLAGS:
		value = getSetting8021xPinFlagsJSON(data)
	case NM_SETTING_802_1X_SYSTEM_CA_CERTS:
		value = getSetting8021xSystemCaCertsJSON(data)
	}
	return
}

// TODO use logic setter
func generalSetSetting8021xKeyJSON(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSetting8021xKey: invalide key", key)
	case NM_SETTING_802_1X_EAP:
		setSetting8021xEapJSON(data, value)
	case NM_SETTING_802_1X_IDENTITY:
		setSetting8021xIdentityJSON(data, value)
	case NM_SETTING_802_1X_ANONYMOUS_IDENTITY:
		setSetting8021xAnonymousIdentityJSON(data, value)
	case NM_SETTING_802_1X_PAC_FILE:
		setSetting8021xPacFileJSON(data, value)
	case NM_SETTING_802_1X_CA_CERT:
		setSetting8021xCaCertJSON(data, value)
	case NM_SETTING_802_1X_CA_PATH:
		setSetting8021xCaPathJSON(data, value)
	case NM_SETTING_802_1X_SUBJECT_MATCH:
		setSetting8021xSubjectMatchJSON(data, value)
	case NM_SETTING_802_1X_ALTSUBJECT_MATCHES:
		setSetting8021xAltsubjectMatchesJSON(data, value)
	case NM_SETTING_802_1X_CLIENT_CERT:
		setSetting8021xClientCertJSON(data, value)
	case NM_SETTING_802_1X_PHASE1_PEAPVER:
		setSetting8021xPhase1PeapverJSON(data, value)
	case NM_SETTING_802_1X_PHASE1_PEAPLABEL:
		setSetting8021xPhase1PeaplabelJSON(data, value)
	case NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING:
		setSetting8021xPhase1FastProvisioningJSON(data, value)
	case NM_SETTING_802_1X_PHASE2_AUTH:
		setSetting8021xPhase2AuthJSON(data, value)
	case NM_SETTING_802_1X_PHASE2_AUTHEAP:
		setSetting8021xPhase2AutheapJSON(data, value)
	case NM_SETTING_802_1X_PHASE2_CA_CERT:
		setSetting8021xPhase2CaCertJSON(data, value)
	case NM_SETTING_802_1X_PHASE2_CA_PATH:
		setSetting8021xPhase2CaPathJSON(data, value)
	case NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH:
		setSetting8021xPhase2SubjectMatchJSON(data, value)
	case NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES:
		setSetting8021xPhase2AltsubjectMatchesJSON(data, value)
	case NM_SETTING_802_1X_PASSWORD:
		setSetting8021xPasswordJSON(data, value)
	case NM_SETTING_802_1X_PASSWORD_FLAGS:
		setSetting8021xPasswordFlagsJSON(data, value)
	case NM_SETTING_802_1X_PASSWORD_RAW:
		setSetting8021xPasswordRawJSON(data, value)
	case NM_SETTING_802_1X_PASSWORD_RAW_FLAGS:
		setSetting8021xPasswordRawFlagsJSON(data, value)
	case NM_SETTING_802_1X_PRIVATE_KEY:
		setSetting8021xPrivateKeyJSON(data, value)
	case NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD:
		setSetting8021xPrivateKeyPasswordJSON(data, value)
	case NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS:
		setSetting8021xPrivateKeyPasswordFlagsJSON(data, value)
	case NM_SETTING_802_1X_PHASE2_PRIVATE_KEY:
		setSetting8021xPhase2PrivateKeyJSON(data, value)
	case NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD:
		setSetting8021xPhase2PrivateKeyPasswordJSON(data, value)
	case NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS:
		setSetting8021xPhase2PrivateKeyPasswordFlagsJSON(data, value)
	case NM_SETTING_802_1X_PIN:
		setSetting8021xPinJSON(data, value)
	case NM_SETTING_802_1X_PIN_FLAGS:
		setSetting8021xPinFlagsJSON(data, value)
	case NM_SETTING_802_1X_SYSTEM_CA_CERTS:
		setSetting8021xSystemCaCertsJSON(data, value)
	}
	return
}
