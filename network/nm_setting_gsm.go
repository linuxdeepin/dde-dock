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

import (
	"pkg.linuxdeepin.com/lib/dbus"
	. "pkg.linuxdeepin.com/lib/gettext"
	"pkg.linuxdeepin.com/lib/iso"
	"pkg.linuxdeepin.com/lib/mobileprovider"
)

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

const mobileProviderValueCustom = "<custom>"

func newMobileConnectionForDevice(id, uuid string, devPath dbus.ObjectPath, active bool) (cpath dbus.ObjectPath, err error) {
	logger.Infof("new mobile connection, id=%s, uuid=%s, devPath=%s", id, uuid, devPath)

	// guess default plan for mobile device
	countryCode, _ := iso.GetLocaleCountryCode()
	serviceType := getMobileDeviceServicType(devPath)
	plan, err := getDefaultPlanForMobileDevice(countryCode, serviceType)
	if err != nil {
		return
	}

	data := newMobileConnectionData("mobile", uuid, serviceType)
	addSettingSection(data, sectionCache)
	logicSetSettingVkMobileCountry(data, countryCode)
	logicSetSettingVkMobileProvider(data, plan.ProviderName)
	logicSetSettingVkMobilePlan(data, mobileprovider.MarshalPlan(plan))
	refileSectionCache(data)

	if active {
		cpath, _, err = nmAddAndActivateConnection(data, devPath)
	} else {
		cpath, err = nmAddConnection(data)
	}
	return
}
func getDefaultPlanForMobileDevice(countryCode, serviceType string) (plan mobileprovider.Plan, err error) {
	if serviceType == connectionMobileGsm {
		plan, err = mobileprovider.GetDefaultGSMPlanForCountry(countryCode)
	} else {
		plan, err = mobileprovider.GetDefaultCDMAPlanForCountry(countryCode)
	}
	if err != nil {
		logger.Error(err)
	}
	return
}
func getMobileDeviceServicType(devPath dbus.ObjectPath) (serviceType string) {
	capabilities := nmGetDeviceModemCapabilities(devPath)
	if (capabilities & NM_DEVICE_MODEM_CAPABILITY_LTE) == capabilities {
		// all LTE modems treated as GSM/UMTS
		serviceType = connectionMobileGsm
	} else if (capabilities & NM_DEVICE_MODEM_CAPABILITY_GSM_UMTS) == capabilities {
		serviceType = connectionMobileGsm
	} else if (capabilities & NM_DEVICE_MODEM_CAPABILITY_CDMA_EVDO) == capabilities {
		serviceType = connectionMobileCdma
	} else {
		logger.Errorf("Unknown modem capabilities (0x%x)", capabilities)
	}
	return
}

func newMobileConnectionData(id, uuid, serviceType string) (data connectionData) {
	data = make(connectionData)

	addSettingSection(data, sectionConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionAutoconnect(data, true)

	logicSetSettingVkMobileServiceType(data, serviceType)

	addSettingSection(data, sectionPpp)
	logicSetSettingVkPppEnableLcpEcho(data, true)

	addSettingSection(data, sectionSerial)
	setSettingSerialBaud(data, 115200)

	initSettingSectionIpv4(data)

	return
}

func initSettingSectionGsm(data connectionData) {
	setSettingConnectionType(data, NM_SETTING_GSM_SETTING_NAME)
	addSettingSection(data, sectionGsm)
	setSettingGsmNumber(data, "*99#")
	setSettingGsmPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	setSettingGsmPinFlags(data, NM_SETTING_SECRET_FLAG_NONE)
}

// Get available keys
func getSettingGsmAvailableKeys(data connectionData) (keys []string) {
	if getSettingVkMobileProvider(data) == mobileProviderValueCustom {
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_NUMBER)
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_USERNAME)
		if isSettingRequireSecret(getSettingGsmPasswordFlags(data)) {
			keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_PASSWORD)
		}
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_APN)
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_NETWORK_ID)
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_NETWORK_TYPE)
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_HOME_ONLY)
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_PIN)
	}
	return
}

// Get available values
func getSettingGsmAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_GSM_PASSWORD_FLAGS:
		values = availableValuesSettingSecretFlags
	case NM_SETTING_GSM_APN:
	case NM_SETTING_GSM_NETWORK_TYPE:
		values = []kvalue{
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_ANY, Tr("Any")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_UMTS_HSPA, Tr("3G (UMTS/HSPA)")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_GPRS_EDGE, Tr("2G (GPRS/EDGE)")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_PREFER_UMTS_HSPA, Tr("Prefer 3G (UMTS/HSPA)")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_PREFER_GPRS_EDGE, Tr("Prefer 2G (GPRS/EDGE)")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_PREFER_4G, Tr("Prefer 4G (LTE)")},
			kvalue{NM_SETTING_GSM_NETWORK_TYPE_4G, Tr("Use Only 4G (LTE)")},
		}
	}
	return
}

// Check whether the values are correct
func checkSettingGsmValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	ensureSettingGsmApnNoEmpty(data, errs)
	ensureSettingGsmNumberNoEmpty(data, errs)
	return
}

func syncMoibleConnectionId(data connectionData) {
	// sync connection name
	if !isSettingSectionExists(data, sectionCache) {
		return
	}
	providerName := getSettingVkMobileProvider(data)
	if providerName == mobileProviderValueCustom {
		switch getSettingVkMobileServiceType(data) {
		case connectionMobileGsm:
			setSettingConnectionId(data, Tr("Custom")+" GSM")
		case connectionMobileCdma:
			setSettingConnectionId(data, Tr("Custom")+" CDMA")
		}
	} else {
		if plan, err := mobileprovider.UnmarshalPlan(getSettingVkMobilePlan(data)); err == nil {
			if plan.IsGSM {
				if len(plan.Name) > 0 {
					setSettingConnectionId(data, providerName+" "+plan.Name)
				} else {
					setSettingConnectionId(data, providerName+" "+Tr("Default"))
				}
			} else {
				setSettingConnectionId(data, providerName)
			}
		}
	}
}
