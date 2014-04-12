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

func doGetSetting8021xEap(data _ConnectionData) (eap string) {
	eaps := getSetting8021xEap(data)
	if len(eaps) == 0 {
		LOGGER.Error("eap value is empty")
		return
	}
	eap = eaps[0]
	return
}

// Get available keys
func getSetting8021xAvailableKeys(data _ConnectionData) (keys []string) {
	keys = []string{NM_SETTING_802_1X_EAP}
	switch doGetSetting8021xEap(data) {
	case "tls":
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_IDENTITY)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_CLIENT_CERT)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_CA_CERT)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PRIVATE_KEY)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PASSWORD)
	case "leap":
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_IDENTITY)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PASSWORD)
	case "fast":
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_ANONYMOUS_IDENTITY)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PAC_FILE)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PHASE2_AUTH)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_IDENTITY)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PASSWORD)
	case "ttls":
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_ANONYMOUS_IDENTITY)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_CA_CERT)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PHASE2_AUTH)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_IDENTITY)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PASSWORD)
	case "peap":
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_ANONYMOUS_IDENTITY)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_CA_CERT)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PHASE2_AUTH)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_IDENTITY)
		keys = appendStrArrayUnion(keys, NM_SETTING_802_1X_PASSWORD)
	}
	return
}

// Get available values
func getSetting8021xAvailableValues(data _ConnectionData, key string) (values []string, customizable bool) {
	customizable = true
	switch key {
	case NM_SETTING_802_1X_EAP:
		values = []string{"tls", "leap", "fast", "ttls", "peap"}
	case NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING:
		values = []string{
			"0", // disabled
			"1", // allow unauthenticated provisioning
			"2", // allow authenticated provisioning
			"3", // allow both authenticated and unauthenticated provisioning
		}
	case NM_SETTING_802_1X_PHASE2_AUTH:
		// 'pap', 'chap', 'mschap', 'mschapv2', 'gtc', 'otp', 'md5', and 'tls'
		switch doGetSetting8021xEap(data) {
		case "tls":
		case "leap":
		case "fast":
			values = []string{"gtc", "mschapv2"}
		case "ttls":
			values = []string{"pap", "mschap", "mschapv2", "chap"}
		case "peap":
			values = []string{"gtc", "md5", "mschapv2"}
		}
	case NM_SETTING_802_1X_PASSWORD_FLAGS: // TODO available values not string
		// values = []string{
		// 	NM_SETTING_SECRET_FLAG_NONE,
		// 	NM_SETTING_SECRET_FLAG_AGENT_OWNED,
		// 	NM_SETTING_SECRET_FLAG_NOT_SAVED,
		// 	NM_SETTING_SECRET_FLAG_NOT_REQUIRED,
		// }
	}
	return
}

// Check whether the values are correct
func checkSetting8021xValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)

	// check eap
	ensureSetting8021xEapNoEmpty(data, errs)
	switch doGetSetting8021xEap(data) {
	default:
		rememberError(errs, NM_SETTING_802_1X_EAP, NM_KEY_ERROR_INVALID_VALUE)
	case "tls":
		ensureSetting8021xIdentityNoEmpty(data, errs)
		ensureSetting8021xClientCertNoEmpty(data, errs)
		ensureSetting8021xCaCertNoEmpty(data, errs)
		ensureSetting8021xPrivateKeyNoEmpty(data, errs)
		ensureSetting8021xPasswordNoEmpty(data, errs)
		ensureSetting8021xSystemCaCertsNoEmpty(data, errs)
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

	// TODO check value of client cert
	// TODO check value of ca cert
	// TODO check value of private key

	return
}

// Set JSON value generally
// TODO use logic setter
func generalSetSetting8021xKeyJSON(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSetting8021xKey: invalide key", key)
	case NM_SETTING_802_1X_EAP:
		// TODO
		// logicSetSetting8021xEapJSON(data, value)
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

// Logic setter
// TODO
