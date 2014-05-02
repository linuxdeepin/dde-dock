package main

import (
	"dlib"
)

// TODO doc

const NM_SETTING_GSM_SETTING_NAME = "gsm"

const (
	NM_SETTING_GSM_NUMBER         = "number"
	NM_SETTING_GSM_USERNAME       = "username"
	NM_SETTING_GSM_PASSWORD       = "password"
	NM_SETTING_GSM_PASSWORD_FLAGS = "password-flags"
	NM_SETTING_GSM_APN            = "apn"
	NM_SETTING_GSM_NETWORK_ID     = "network-id"
	NM_SETTING_GSM_NETWORK_TYPE   = "network-type"
	NM_SETTING_GSM_ALLOWED_BANDS  = "allowed-bands"
	NM_SETTING_GSM_PIN            = "pin"
	NM_SETTING_GSM_PIN_FLAGS      = "pin-flags"
	NM_SETTING_GSM_HOME_ONLY      = "home-only"
)

const (
	NM_SETTING_GSM_NETWORK_TYPE_ANY              = -1
	NM_SETTING_GSM_NETWORK_TYPE_UMTS_HSPA        = 0
	NM_SETTING_GSM_NETWORK_TYPE_GPRS_EDGE        = 1
	NM_SETTING_GSM_NETWORK_TYPE_PREFER_UMTS_HSPA = 2
	NM_SETTING_GSM_NETWORK_TYPE_PREFER_GPRS_EDGE = 3
	NM_SETTING_GSM_NETWORK_TYPE_PREFER_4G        = 4
	NM_SETTING_GSM_NETWORK_TYPE_4G               = 5
)

const (
	NM_SETTING_GSM_BAND_UNKNOWN = 0x00000000
	NM_SETTING_GSM_BAND_ANY     = 0x00000001
	NM_SETTING_GSM_BAND_EGSM    = 0x00000002 /*  900 MHz */
	NM_SETTING_GSM_BAND_DCS     = 0x00000004 /* 1800 MHz */
	NM_SETTING_GSM_BAND_PCS     = 0x00000008 /* 1900 MHz */
	NM_SETTING_GSM_BAND_G850    = 0x00000010 /*  850 MHz */
	NM_SETTING_GSM_BAND_U2100   = 0x00000020 /* WCDMA 3GPP UMTS 2100 MHz     (Class I) */
	NM_SETTING_GSM_BAND_U1800   = 0x00000040 /* WCDMA 3GPP UMTS 1800 MHz     (Class III) */
	NM_SETTING_GSM_BAND_U17IV   = 0x00000080 /* WCDMA 3GPP AWS 1700/2100 MHz (Class IV) */
	NM_SETTING_GSM_BAND_U800    = 0x00000100 /* WCDMA 3GPP UMTS 800 MHz      (Class VI) */
	NM_SETTING_GSM_BAND_U850    = 0x00000200 /* WCDMA 3GPP UMTS 850 MHz      (Class V) */
	NM_SETTING_GSM_BAND_U900    = 0x00000400 /* WCDMA 3GPP UMTS 900 MHz      (Class VIII) */
	NM_SETTING_GSM_BAND_U17IX   = 0x00000800 /* WCDMA 3GPP UMTS 1700 MHz     (Class IX) */
	NM_SETTING_GSM_BAND_U1900   = 0x00001000 /* WCDMA 3GPP UMTS 1900 MHz     (Class II) */
	NM_SETTING_GSM_BAND_U2600   = 0x00002000 /* WCDMA 3GPP UMTS 2600 MHz     (Class VII, internal) */
)

func newGsmConnectionData(id, uuid string) (data connectionData) {
	data = make(connectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, NM_SETTING_GSM_SETTING_NAME)

	addSettingField(data, fieldGsm)
	setSettingGsmPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	setSettingGsmPinFlags(data, NM_SETTING_SECRET_FLAG_NONE)

	addSettingField(data, fieldPpp)

	initSettingFieldIpv4(data)

	return
}

// Get available keys
func getSettingGsmAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, fieldGsm, NM_SETTING_GSM_NUMBER)
	keys = appendAvailableKeys(data, keys, fieldGsm, NM_SETTING_GSM_USERNAME)
	keys = appendAvailableKeys(data, keys, fieldGsm, NM_SETTING_GSM_PASSWORD)
	keys = appendAvailableKeys(data, keys, fieldGsm, NM_SETTING_GSM_PASSWORD_FLAGS)
	keys = appendAvailableKeys(data, keys, fieldGsm, NM_SETTING_GSM_APN)
	keys = appendAvailableKeys(data, keys, fieldGsm, NM_SETTING_GSM_NETWORK_ID)
	keys = appendAvailableKeys(data, keys, fieldGsm, NM_SETTING_GSM_NETWORK_TYPE)
	keys = appendAvailableKeys(data, keys, fieldGsm, NM_SETTING_GSM_HOME_ONLY)
	keys = appendAvailableKeys(data, keys, fieldGsm, NM_SETTING_GSM_PIN)
	return
}

// Get available values
func getSettingGsmAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_GSM_NUMBER:
		values = []kvalue{
			kvalue{"*99#", dlib.Tr("*99#")},
		}
	case NM_SETTING_GSM_PASSWORD_FLAGS:
		values = availableValuesNMSettingSecretFlag
	case NM_SETTING_GSM_APN:
		// TODO
		values = []kvalue{
			kvalue{"3gnet", dlib.Tr("China Unicom Internet")},
			kvalue{"3gwap", dlib.Tr("China Unicom MMS")},
			kvalue{"cmwap", dlib.Tr("China Mobile WAP")},
			kvalue{"cmnet", dlib.Tr("China Mobile Internet")},
			kvalue{"cmwap", dlib.Tr("China Mobile MMS")},
		}
	case NM_SETTING_GSM_NETWORK_TYPE:
		values = []kvalue{
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_ANY, dlib.Tr("Any")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_UMTS_HSPA, dlib.Tr("3G (UMTS/HSPA)")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_GPRS_EDGE, dlib.Tr("2G (GPRS/EDGE)")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_PREFER_UMTS_HSPA, dlib.Tr("Prefer 3G (UMTS/HSPA)")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_PREFER_GPRS_EDGE, dlib.Tr("Prefer 2G (GPRS/EDGE)")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_PREFER_4G, dlib.Tr("Prefer 4G (LTE)")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_4G, dlib.Tr("Use Only 4G (LTE)")},
		}
	}
	return
}

// Check whether the values are correct
func checkSettingGsmValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)
	// TODO
	return
}
