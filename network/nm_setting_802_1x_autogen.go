package main

// Get key type
func getSetting8021xKeyType(key string) (t ktype) {
	switch key {
	default:
		t = ktypeUnknown
	case NM_SETTING_802_1X_EAP:
		t = ktypeArrayString
	case NM_SETTING_802_1X_IDENTITY:
		t = ktypeString
	case NM_SETTING_802_1X_ANONYMOUS_IDENTITY:
		t = ktypeString
	case NM_SETTING_802_1X_PAC_FILE:
		t = ktypeString
	case NM_SETTING_802_1X_CA_CERT:
		t = ktypeArrayByte
	case NM_SETTING_802_1X_CA_PATH:
		t = ktypeString
	case NM_SETTING_802_1X_SUBJECT_MATCH:
		t = ktypeString
	case NM_SETTING_802_1X_ALTSUBJECT_MATCHES:
		t = ktypeArrayString
	case NM_SETTING_802_1X_CLIENT_CERT:
		t = ktypeArrayByte
	case NM_SETTING_802_1X_PHASE1_PEAPVER:
		t = ktypeString
	case NM_SETTING_802_1X_PHASE1_PEAPLABEL:
		t = ktypeString
	case NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING:
		t = ktypeString
	case NM_SETTING_802_1X_PHASE2_AUTH:
		t = ktypeString
	case NM_SETTING_802_1X_PHASE2_AUTHEAP:
		t = ktypeString
	case NM_SETTING_802_1X_PHASE2_CA_CERT:
		t = ktypeArrayByte
	case NM_SETTING_802_1X_PHASE2_CA_PATH:
		t = ktypeString
	case NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH:
		t = ktypeString
	case NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES:
		t = ktypeArrayString
	case NM_SETTING_802_1X_PHASE2_CLIENT_CERT:
		t = ktypeArrayByte
	case NM_SETTING_802_1X_PASSWORD:
		t = ktypeString
	case NM_SETTING_802_1X_PASSWORD_FLAGS:
		t = ktypeUint32
	case NM_SETTING_802_1X_PASSWORD_RAW:
		t = ktypeArrayByte
	case NM_SETTING_802_1X_PASSWORD_RAW_FLAGS:
		t = ktypeUint32
	case NM_SETTING_802_1X_PRIVATE_KEY:
		t = ktypeArrayByte
	case NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD:
		t = ktypeString
	case NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS:
		t = ktypeUint32
	case NM_SETTING_802_1X_PHASE2_PRIVATE_KEY:
		t = ktypeArrayByte
	case NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD:
		t = ktypeString
	case NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS:
		t = ktypeUint32
	case NM_SETTING_802_1X_PIN:
		t = ktypeString
	case NM_SETTING_802_1X_PIN_FLAGS:
		t = ktypeUint32
	case NM_SETTING_802_1X_SYSTEM_CA_CERTS:
		t = ktypeBoolean
	}
	return t
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

// Getter
func getSetting8021xEapJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_EAP, getSetting8021xKeyType(NM_SETTING_802_1X_EAP))
	return
}
func getSetting8021xIdentityJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_IDENTITY, getSetting8021xKeyType(NM_SETTING_802_1X_IDENTITY))
	return
}
func getSetting8021xAnonymousIdentityJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_ANONYMOUS_IDENTITY, getSetting8021xKeyType(NM_SETTING_802_1X_ANONYMOUS_IDENTITY))
	return
}
func getSetting8021xPacFileJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PAC_FILE, getSetting8021xKeyType(NM_SETTING_802_1X_PAC_FILE))
	return
}
func getSetting8021xCaCertJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CA_CERT, getSetting8021xKeyType(NM_SETTING_802_1X_CA_CERT))
	return
}
func getSetting8021xCaPathJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CA_PATH, getSetting8021xKeyType(NM_SETTING_802_1X_CA_PATH))
	return
}
func getSetting8021xSubjectMatchJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_SUBJECT_MATCH, getSetting8021xKeyType(NM_SETTING_802_1X_SUBJECT_MATCH))
	return
}
func getSetting8021xAltsubjectMatchesJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_ALTSUBJECT_MATCHES, getSetting8021xKeyType(NM_SETTING_802_1X_ALTSUBJECT_MATCHES))
	return
}
func getSetting8021xClientCertJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CLIENT_CERT, getSetting8021xKeyType(NM_SETTING_802_1X_CLIENT_CERT))
	return
}
func getSetting8021xPhase1PeapverJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_PEAPVER, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_PEAPVER))
	return
}
func getSetting8021xPhase1PeaplabelJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_PEAPLABEL, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_PEAPLABEL))
	return
}
func getSetting8021xPhase1FastProvisioningJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING))
	return
}
func getSetting8021xPhase2AuthJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_AUTH, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_AUTH))
	return
}
func getSetting8021xPhase2AutheapJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_AUTHEAP, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_AUTHEAP))
	return
}
func getSetting8021xPhase2CaCertJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CA_CERT, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CA_CERT))
	return
}
func getSetting8021xPhase2CaPathJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CA_PATH, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CA_PATH))
	return
}
func getSetting8021xPhase2SubjectMatchJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH))
	return
}
func getSetting8021xPhase2AltsubjectMatchesJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES))
	return
}
func getSetting8021xPhase2ClientCertJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CLIENT_CERT, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CLIENT_CERT))
	return
}
func getSetting8021xPasswordJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD))
	return
}
func getSetting8021xPasswordFlagsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_FLAGS, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_FLAGS))
	return
}
func getSetting8021xPasswordRawJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_RAW, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_RAW))
	return
}
func getSetting8021xPasswordRawFlagsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_RAW_FLAGS, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_RAW_FLAGS))
	return
}
func getSetting8021xPrivateKeyJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY))
	return
}
func getSetting8021xPrivateKeyPasswordJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD))
	return
}
func getSetting8021xPrivateKeyPasswordFlagsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS))
	return
}
func getSetting8021xPhase2PrivateKeyJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY))
	return
}
func getSetting8021xPhase2PrivateKeyPasswordJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD))
	return
}
func getSetting8021xPhase2PrivateKeyPasswordFlagsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS))
	return
}
func getSetting8021xPinJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PIN, getSetting8021xKeyType(NM_SETTING_802_1X_PIN))
	return
}
func getSetting8021xPinFlagsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PIN_FLAGS, getSetting8021xKeyType(NM_SETTING_802_1X_PIN_FLAGS))
	return
}
func getSetting8021xSystemCaCertsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_SYSTEM_CA_CERTS, getSetting8021xKeyType(NM_SETTING_802_1X_SYSTEM_CA_CERTS))
	return
}

// Setter
func setSetting8021xEapJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_EAP, value, getSetting8021xKeyType(NM_SETTING_802_1X_EAP))
}
func setSetting8021xIdentityJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_IDENTITY, value, getSetting8021xKeyType(NM_SETTING_802_1X_IDENTITY))
}
func setSetting8021xAnonymousIdentityJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_ANONYMOUS_IDENTITY, value, getSetting8021xKeyType(NM_SETTING_802_1X_ANONYMOUS_IDENTITY))
}
func setSetting8021xPacFileJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PAC_FILE, value, getSetting8021xKeyType(NM_SETTING_802_1X_PAC_FILE))
}
func setSetting8021xCaCertJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CA_CERT, value, getSetting8021xKeyType(NM_SETTING_802_1X_CA_CERT))
}
func setSetting8021xCaPathJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CA_PATH, value, getSetting8021xKeyType(NM_SETTING_802_1X_CA_PATH))
}
func setSetting8021xSubjectMatchJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_SUBJECT_MATCH, value, getSetting8021xKeyType(NM_SETTING_802_1X_SUBJECT_MATCH))
}
func setSetting8021xAltsubjectMatchesJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_ALTSUBJECT_MATCHES, value, getSetting8021xKeyType(NM_SETTING_802_1X_ALTSUBJECT_MATCHES))
}
func setSetting8021xClientCertJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CLIENT_CERT, value, getSetting8021xKeyType(NM_SETTING_802_1X_CLIENT_CERT))
}
func setSetting8021xPhase1PeapverJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_PEAPVER, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_PEAPVER))
}
func setSetting8021xPhase1PeaplabelJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_PEAPLABEL, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_PEAPLABEL))
}
func setSetting8021xPhase1FastProvisioningJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING))
}
func setSetting8021xPhase2AuthJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_AUTH, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_AUTH))
}
func setSetting8021xPhase2AutheapJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_AUTHEAP, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_AUTHEAP))
}
func setSetting8021xPhase2CaCertJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CA_CERT, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CA_CERT))
}
func setSetting8021xPhase2CaPathJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CA_PATH, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CA_PATH))
}
func setSetting8021xPhase2SubjectMatchJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH))
}
func setSetting8021xPhase2AltsubjectMatchesJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES))
}
func setSetting8021xPhase2ClientCertJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CLIENT_CERT, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CLIENT_CERT))
}
func setSetting8021xPasswordJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD, value, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD))
}
func setSetting8021xPasswordFlagsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_FLAGS, value, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_FLAGS))
}
func setSetting8021xPasswordRawJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_RAW, value, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_RAW))
}
func setSetting8021xPasswordRawFlagsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_RAW_FLAGS, value, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_RAW_FLAGS))
}
func setSetting8021xPrivateKeyJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY, value, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY))
}
func setSetting8021xPrivateKeyPasswordJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD, value, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD))
}
func setSetting8021xPrivateKeyPasswordFlagsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS, value, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS))
}
func setSetting8021xPhase2PrivateKeyJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY))
}
func setSetting8021xPhase2PrivateKeyPasswordJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD))
}
func setSetting8021xPhase2PrivateKeyPasswordFlagsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS))
}
func setSetting8021xPinJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PIN, value, getSetting8021xKeyType(NM_SETTING_802_1X_PIN))
}
func setSetting8021xPinFlagsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PIN_FLAGS, value, getSetting8021xKeyType(NM_SETTING_802_1X_PIN_FLAGS))
}
func setSetting8021xSystemCaCertsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_SYSTEM_CA_CERTS, value, getSetting8021xKeyType(NM_SETTING_802_1X_SYSTEM_CA_CERTS))
}

// Remover
func removeSetting8021xEapJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_EAP)
}
func removeSetting8021xIdentityJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_IDENTITY)
}
func removeSetting8021xAnonymousIdentityJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_ANONYMOUS_IDENTITY)
}
func removeSetting8021xPacFileJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PAC_FILE)
}
func removeSetting8021xCaCertJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CA_CERT)
}
func removeSetting8021xCaPathJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CA_PATH)
}
func removeSetting8021xSubjectMatchJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_SUBJECT_MATCH)
}
func removeSetting8021xAltsubjectMatchesJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_ALTSUBJECT_MATCHES)
}
func removeSetting8021xClientCertJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CLIENT_CERT)
}
func removeSetting8021xPhase1PeapverJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_PEAPVER)
}
func removeSetting8021xPhase1PeaplabelJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_PEAPLABEL)
}
func removeSetting8021xPhase1FastProvisioningJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING)
}
func removeSetting8021xPhase2AuthJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_AUTH)
}
func removeSetting8021xPhase2AutheapJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_AUTHEAP)
}
func removeSetting8021xPhase2CaCertJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CA_CERT)
}
func removeSetting8021xPhase2CaPathJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CA_PATH)
}
func removeSetting8021xPhase2SubjectMatchJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH)
}
func removeSetting8021xPhase2AltsubjectMatchesJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES)
}
func removeSetting8021xPhase2ClientCertJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CLIENT_CERT)
}
func removeSetting8021xPasswordJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD)
}
func removeSetting8021xPasswordFlagsJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_FLAGS)
}
func removeSetting8021xPasswordRawJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_RAW)
}
func removeSetting8021xPasswordRawFlagsJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_RAW_FLAGS)
}
func removeSetting8021xPrivateKeyJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY)
}
func removeSetting8021xPrivateKeyPasswordJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD)
}
func removeSetting8021xPrivateKeyPasswordFlagsJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS)
}
func removeSetting8021xPhase2PrivateKeyJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY)
}
func removeSetting8021xPhase2PrivateKeyPasswordJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD)
}
func removeSetting8021xPhase2PrivateKeyPasswordFlagsJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS)
}
func removeSetting8021xPinJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PIN)
}
func removeSetting8021xPinFlagsJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PIN_FLAGS)
}
func removeSetting8021xSystemCaCertsJSON(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_SYSTEM_CA_CERTS)
}
