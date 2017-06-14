/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/iso"
	"pkg.deepin.io/lib/mobileprovider"
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
	addSetting(data, sectionCache)
	logicSetSettingVkMobileCountry(data, countryCode)
	logicSetSettingVkMobileProvider(data, plan.ProviderName)
	logicSetSettingVkMobilePlan(data, mobileprovider.MarshalPlan(plan))
	refileSectionCache(data)

	if active {
		cpath, _, err = nmAddAndActivateConnection(data, devPath, true)
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
	if (capabilities & nm.NM_DEVICE_MODEM_CAPABILITY_LTE) == capabilities {
		// all LTE modems treated as GSM/UMTS
		serviceType = connectionMobileGsm
	} else if (capabilities & nm.NM_DEVICE_MODEM_CAPABILITY_GSM_UMTS) == capabilities {
		serviceType = connectionMobileGsm
	} else if (capabilities & nm.NM_DEVICE_MODEM_CAPABILITY_CDMA_EVDO) == capabilities {
		serviceType = connectionMobileCdma
	} else {
		logger.Errorf("Unknown modem capabilities (0x%x)", capabilities)
	}
	return
}

func newMobileConnectionData(id, uuid, serviceType string) (data connectionData) {
	data = make(connectionData)

	addSetting(data, nm.NM_SETTING_CONNECTION_SETTING_NAME)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionAutoconnect(data, true)

	logicSetSettingVkMobileServiceType(data, serviceType)

	addSetting(data, nm.NM_SETTING_PPP_SETTING_NAME)
	logicSetSettingVkPppEnableLcpEcho(data, true)

	addSetting(data, nm.NM_SETTING_SERIAL_SETTING_NAME)
	setSettingSerialBaud(data, 115200)

	initSettingSectionIpv4(data)

	return
}

func initSettingSectionGsm(data connectionData) {
	setSettingConnectionType(data, nm.NM_SETTING_GSM_SETTING_NAME)
	addSetting(data, nm.NM_SETTING_GSM_SETTING_NAME)
	setSettingGsmNumber(data, "*99#")
	setSettingGsmPasswordFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
	setSettingGsmPinFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
}

// Get available keys
func getSettingGsmAvailableKeys(data connectionData) (keys []string) {
	if getSettingVkMobileProvider(data) == mobileProviderValueCustom {
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_GSM_SETTING_NAME, nm.NM_SETTING_GSM_NUMBER)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_GSM_SETTING_NAME, nm.NM_SETTING_GSM_USERNAME)
		if isSettingRequireSecret(getSettingGsmPasswordFlags(data)) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_GSM_SETTING_NAME, nm.NM_SETTING_GSM_PASSWORD)
		}
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_GSM_SETTING_NAME, nm.NM_SETTING_GSM_APN)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_GSM_SETTING_NAME, nm.NM_SETTING_GSM_NETWORK_ID)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_GSM_SETTING_NAME, nm.NM_SETTING_GSM_HOME_ONLY)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_GSM_SETTING_NAME, nm.NM_SETTING_GSM_PIN)
	}
	return
}

// Get available values
func getSettingGsmAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_GSM_PASSWORD_FLAGS:
		values = availableValuesSettingSecretFlags
	case nm.NM_SETTING_GSM_APN:
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
	if !isSettingExists(data, sectionCache) {
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
