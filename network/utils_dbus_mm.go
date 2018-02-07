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
	mmdbus "dbus/org/freedesktop/modemmanager1"
	"pkg.deepin.io/lib/dbus"
)

// modem capabilities
const (
	MM_MODEM_CAPABILITY_NONE         = 0
	MM_MODEM_CAPABILITY_POTS         = 1 << 0
	MM_MODEM_CAPABILITY_CDMA_EVDO    = 1 << 1
	MM_MODEM_CAPABILITY_GSM_UMTS     = 1 << 2
	MM_MODEM_CAPABILITY_LTE          = 1 << 3
	MM_MODEM_CAPABILITY_LTE_ADVANCED = 1 << 4
	MM_MODEM_CAPABILITY_IRIDIUM      = 1 << 5
	MM_MODEM_CAPABILITY_ANY          = 0xFFFFFFF
)

// modem network access technologies
const (
	MM_MODEM_ACCESS_TECHNOLOGY_UNKNOWN     = 0
	MM_MODEM_ACCESS_TECHNOLOGY_POTS        = 1 << 0
	MM_MODEM_ACCESS_TECHNOLOGY_GSM         = 1 << 1
	MM_MODEM_ACCESS_TECHNOLOGY_GSM_COMPACT = 1 << 2
	MM_MODEM_ACCESS_TECHNOLOGY_GPRS        = 1 << 3
	MM_MODEM_ACCESS_TECHNOLOGY_EDGE        = 1 << 4
	MM_MODEM_ACCESS_TECHNOLOGY_UMTS        = 1 << 5
	MM_MODEM_ACCESS_TECHNOLOGY_HSDPA       = 1 << 6
	MM_MODEM_ACCESS_TECHNOLOGY_HSUPA       = 1 << 7
	MM_MODEM_ACCESS_TECHNOLOGY_HSPA        = 1 << 8
	MM_MODEM_ACCESS_TECHNOLOGY_HSPA_PLUS   = 1 << 9
	MM_MODEM_ACCESS_TECHNOLOGY_1XRTT       = 1 << 10
	MM_MODEM_ACCESS_TECHNOLOGY_EVDO0       = 1 << 11
	MM_MODEM_ACCESS_TECHNOLOGY_EVDOA       = 1 << 12
	MM_MODEM_ACCESS_TECHNOLOGY_EVDOB       = 1 << 13
	MM_MODEM_ACCESS_TECHNOLOGY_LTE         = 1 << 14
	MM_MODEM_ACCESS_TECHNOLOGY_ANY         = 0xFFFFFFFF
)

// modem modes
const (
	MM_MODEM_MODE_NONE = 0
	MM_MODEM_MODE_CS   = 1 << 0
	MM_MODEM_MODE_2G   = 1 << 1
	MM_MODEM_MODE_3G   = 1 << 2
	MM_MODEM_MODE_4G   = 1 << 3
	MM_MODEM_MODE_ANY  = 0xFFFFFFF
)

const (
	moblieNetworkType2G      = "2G"
	moblieNetworkType3G      = "3G"
	moblieNetworkType4G      = "4G"
	moblieNetworkTypeUnknown = "Unknown"
)

func mmNewModem(modemPath dbus.ObjectPath) (modem *mmdbus.Modem, err error) {
	modem, err = mmdbus.NewModem(dbusMmDest, modemPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func mmDestroyModem(modem *mmdbus.Modem) {
	if modem == nil {
		logger.Error("Modem to destroy is nil")
		return
	}
	mmdbus.DestroyModem(modem)
}

func mmGetModemDeviceIdentifier(modemPath dbus.ObjectPath) (devId string, err error) {
	modem, err := mmNewModem(modemPath)
	if err != nil {
		return
	}
	defer mmdbus.DestroyModem(modem)

	devId = modem.DeviceIdentifier.Get()
	return
}

func mmGetModemDeviceSysPath(modemPath dbus.ObjectPath) (sysPath string, err error) {
	modem, err := mmNewModem(modemPath)
	if err != nil {
		return
	}
	defer mmdbus.DestroyModem(modem)

	sysPath = modem.Device.Get()
	return
}

func mmGetModemDeviceSignalQuality(modemPath dbus.ObjectPath) (signalQuality uint32, err error) {
	modem, err := mmNewModem(modemPath)
	if err != nil {
		return
	}
	defer mmdbus.DestroyModem(modem)

	signalQuality = mmDoGetModemDeviceSignalQuality(modem.SignalQuality.Get())
	return
}
func mmDoGetModemDeviceSignalQuality(signalQualityData []interface{}) (signalQuality uint32) {
	if len(signalQualityData) > 0 {
		signalQuality = signalQualityData[0].(uint32)
	}
	return
}

func mmGetModemDeviceAccessTechnologies(modemPath dbus.ObjectPath) (accessTechnologies uint32, err error) {
	modem, err := mmNewModem(modemPath)
	if err != nil {
		return
	}
	defer mmdbus.DestroyModem(modem)

	accessTechnologies = modem.AccessTechnologies.Get()
	return
}

// determine 2g/3g/4g network type through access technologies
func mmGetModemMobileNetworkType(modemPath dbus.ObjectPath) (networkType string) {
	technologies, _ := mmGetModemDeviceAccessTechnologies(modemPath)
	return mmDoGetModemMobileNetworkType(technologies)
}
func mmDoGetModemMobileNetworkType(technologies uint32) (networkType string) {
	switch {
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_ANY) == MM_MODEM_ACCESS_TECHNOLOGY_ANY:
		return moblieNetworkType4G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_ANY) == MM_MODEM_ACCESS_TECHNOLOGY_ANY:
		return moblieNetworkType4G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_LTE) == MM_MODEM_ACCESS_TECHNOLOGY_LTE:
		return moblieNetworkType4G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_EVDOB) == MM_MODEM_ACCESS_TECHNOLOGY_EVDOB:
		return moblieNetworkType3G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_EVDOA) == MM_MODEM_ACCESS_TECHNOLOGY_EVDOA:
		return moblieNetworkType3G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_EVDO0) == MM_MODEM_ACCESS_TECHNOLOGY_EVDO0:
		return moblieNetworkType3G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_1XRTT) == MM_MODEM_ACCESS_TECHNOLOGY_1XRTT:
		return moblieNetworkType3G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_HSPA_PLUS) == MM_MODEM_ACCESS_TECHNOLOGY_HSPA_PLUS:
		return moblieNetworkType3G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_HSPA) == MM_MODEM_ACCESS_TECHNOLOGY_HSPA:
		return moblieNetworkType3G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_HSUPA) == MM_MODEM_ACCESS_TECHNOLOGY_HSUPA:
		return moblieNetworkType3G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_HSDPA) == MM_MODEM_ACCESS_TECHNOLOGY_HSDPA:
		return moblieNetworkType3G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_UMTS) == MM_MODEM_ACCESS_TECHNOLOGY_UMTS:
		return moblieNetworkType3G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_EDGE) == MM_MODEM_ACCESS_TECHNOLOGY_EDGE:
		return moblieNetworkType2G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_GPRS) == MM_MODEM_ACCESS_TECHNOLOGY_GPRS:
		return moblieNetworkType2G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_GSM_COMPACT) == MM_MODEM_ACCESS_TECHNOLOGY_GSM_COMPACT:
		return moblieNetworkType2G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_GSM) == MM_MODEM_ACCESS_TECHNOLOGY_GSM:
		return moblieNetworkType2G
	case (technologies & MM_MODEM_ACCESS_TECHNOLOGY_POTS) == MM_MODEM_ACCESS_TECHNOLOGY_POTS:
		return moblieNetworkType2G
	}
	return moblieNetworkTypeUnknown
}
