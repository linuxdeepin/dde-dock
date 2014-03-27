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

// Getter
func getSetting8021xEap(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_EAP, getSetting8021xKeyType(NM_SETTING_802_1X_EAP))
	return
}
func getSetting8021xIdentity(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_IDENTITY, getSetting8021xKeyType(NM_SETTING_802_1X_IDENTITY))
	return
}
func getSetting8021xAnonymousIdentity(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_ANONYMOUS_IDENTITY, getSetting8021xKeyType(NM_SETTING_802_1X_ANONYMOUS_IDENTITY))
	return
}
func getSetting8021xPacFile(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PAC_FILE, getSetting8021xKeyType(NM_SETTING_802_1X_PAC_FILE))
	return
}
func getSetting8021xCaCert(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CA_CERT, getSetting8021xKeyType(NM_SETTING_802_1X_CA_CERT))
	return
}
func getSetting8021xCaPath(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CA_PATH, getSetting8021xKeyType(NM_SETTING_802_1X_CA_PATH))
	return
}
func getSetting8021xSubjectMatch(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_SUBJECT_MATCH, getSetting8021xKeyType(NM_SETTING_802_1X_SUBJECT_MATCH))
	return
}
func getSetting8021xAltsubjectMatches(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_ALTSUBJECT_MATCHES, getSetting8021xKeyType(NM_SETTING_802_1X_ALTSUBJECT_MATCHES))
	return
}
func getSetting8021xClientCert(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CLIENT_CERT, getSetting8021xKeyType(NM_SETTING_802_1X_CLIENT_CERT))
	return
}
func getSetting8021xPhase1Peapver(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_PEAPVER, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_PEAPVER))
	return
}
func getSetting8021xPhase1Peaplabel(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_PEAPLABEL, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_PEAPLABEL))
	return
}
func getSetting8021xPhase1FastProvisioning(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING))
	return
}
func getSetting8021xPhase2Auth(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_AUTH, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_AUTH))
	return
}
func getSetting8021xPhase2Autheap(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_AUTHEAP, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_AUTHEAP))
	return
}
func getSetting8021xPhase2CaCert(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CA_CERT, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CA_CERT))
	return
}
func getSetting8021xPhase2CaPath(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CA_PATH, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CA_PATH))
	return
}
func getSetting8021xPhase2SubjectMatch(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH))
	return
}
func getSetting8021xPhase2AltsubjectMatches(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES))
	return
}
func getSetting8021xPhase2ClientCert(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CLIENT_CERT, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CLIENT_CERT))
	return
}
func getSetting8021xPassword(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD))
	return
}
func getSetting8021xPasswordFlags(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_FLAGS, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_FLAGS))
	return
}
func getSetting8021xPasswordRaw(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_RAW, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_RAW))
	return
}
func getSetting8021xPasswordRawFlags(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_RAW_FLAGS, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_RAW_FLAGS))
	return
}
func getSetting8021xPrivateKey(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY))
	return
}
func getSetting8021xPrivateKeyPassword(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD))
	return
}
func getSetting8021xPrivateKeyPasswordFlags(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS))
	return
}
func getSetting8021xPhase2PrivateKey(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY))
	return
}
func getSetting8021xPhase2PrivateKeyPassword(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD))
	return
}
func getSetting8021xPhase2PrivateKeyPasswordFlags(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS))
	return
}
func getSetting8021xPin(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PIN, getSetting8021xKeyType(NM_SETTING_802_1X_PIN))
	return
}
func getSetting8021xPinFlags(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PIN_FLAGS, getSetting8021xKeyType(NM_SETTING_802_1X_PIN_FLAGS))
	return
}
func getSetting8021xSystemCaCerts(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_SYSTEM_CA_CERTS, getSetting8021xKeyType(NM_SETTING_802_1X_SYSTEM_CA_CERTS))
	return
}

// Setter
func setSetting8021xEap(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_EAP, value, getSetting8021xKeyType(NM_SETTING_802_1X_EAP))
	return
}
func setSetting8021xIdentity(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_IDENTITY, value, getSetting8021xKeyType(NM_SETTING_802_1X_IDENTITY))
	return
}
func setSetting8021xAnonymousIdentity(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_ANONYMOUS_IDENTITY, value, getSetting8021xKeyType(NM_SETTING_802_1X_ANONYMOUS_IDENTITY))
	return
}
func setSetting8021xPacFile(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PAC_FILE, value, getSetting8021xKeyType(NM_SETTING_802_1X_PAC_FILE))
	return
}
func setSetting8021xCaCert(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CA_CERT, value, getSetting8021xKeyType(NM_SETTING_802_1X_CA_CERT))
	return
}
func setSetting8021xCaPath(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CA_PATH, value, getSetting8021xKeyType(NM_SETTING_802_1X_CA_PATH))
	return
}
func setSetting8021xSubjectMatch(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_SUBJECT_MATCH, value, getSetting8021xKeyType(NM_SETTING_802_1X_SUBJECT_MATCH))
	return
}
func setSetting8021xAltsubjectMatches(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_ALTSUBJECT_MATCHES, value, getSetting8021xKeyType(NM_SETTING_802_1X_ALTSUBJECT_MATCHES))
	return
}
func setSetting8021xClientCert(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_CLIENT_CERT, value, getSetting8021xKeyType(NM_SETTING_802_1X_CLIENT_CERT))
	return
}
func setSetting8021xPhase1Peapver(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_PEAPVER, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_PEAPVER))
	return
}
func setSetting8021xPhase1Peaplabel(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_PEAPLABEL, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_PEAPLABEL))
	return
}
func setSetting8021xPhase1FastProvisioning(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE1_FAST_PROVISIONING))
	return
}
func setSetting8021xPhase2Auth(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_AUTH, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_AUTH))
	return
}
func setSetting8021xPhase2Autheap(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_AUTHEAP, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_AUTHEAP))
	return
}
func setSetting8021xPhase2CaCert(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CA_CERT, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CA_CERT))
	return
}
func setSetting8021xPhase2CaPath(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CA_PATH, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CA_PATH))
	return
}
func setSetting8021xPhase2SubjectMatch(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_SUBJECT_MATCH))
	return
}
func setSetting8021xPhase2AltsubjectMatches(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_ALTSUBJECT_MATCHES))
	return
}
func setSetting8021xPhase2ClientCert(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_CLIENT_CERT, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_CLIENT_CERT))
	return
}
func setSetting8021xPassword(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD, value, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD))
	return
}
func setSetting8021xPasswordFlags(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_FLAGS, value, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_FLAGS))
	return
}
func setSetting8021xPasswordRaw(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_RAW, value, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_RAW))
	return
}
func setSetting8021xPasswordRawFlags(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PASSWORD_RAW_FLAGS, value, getSetting8021xKeyType(NM_SETTING_802_1X_PASSWORD_RAW_FLAGS))
	return
}
func setSetting8021xPrivateKey(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY, value, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY))
	return
}
func setSetting8021xPrivateKeyPassword(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD, value, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD))
	return
}
func setSetting8021xPrivateKeyPasswordFlags(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS, value, getSetting8021xKeyType(NM_SETTING_802_1X_PRIVATE_KEY_PASSWORD_FLAGS))
	return
}
func setSetting8021xPhase2PrivateKey(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY))
	return
}
func setSetting8021xPhase2PrivateKeyPassword(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD))
	return
}
func setSetting8021xPhase2PrivateKeyPasswordFlags(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS, value, getSetting8021xKeyType(NM_SETTING_802_1X_PHASE2_PRIVATE_KEY_PASSWORD_FLAGS))
	return
}
func setSetting8021xPin(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PIN, value, getSetting8021xKeyType(NM_SETTING_802_1X_PIN))
	return
}
func setSetting8021xPinFlags(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_PIN_FLAGS, value, getSetting8021xKeyType(NM_SETTING_802_1X_PIN_FLAGS))
	return
}
func setSetting8021xSystemCaCerts(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_SYSTEM_CA_CERTS, value, getSetting8021xKeyType(NM_SETTING_802_1X_SYSTEM_CA_CERTS))
	return
}
