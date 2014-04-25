package main

import (
	nm "dbus/org/freedesktop/networkmanager"
	"dlib/dbus"
)

type apSecType uint32

const (
	apSecNone apSecType = iota
	apSecWep
	apSecPsk
	apSecEap
)

type accessPoint struct {
	Ssid         string
	Secured      bool
	SecuredInEap bool
	Strength     uint8
	Path         dbus.ObjectPath
}

func NewAccessPoint(apPath dbus.ObjectPath) (ap accessPoint, err error) {
	calcStrength := func(s uint8) uint8 {
		switch {
		case s <= 10:
			return 0
		case s <= 25:
			return 25
		case s <= 50:
			return 50
		case s <= 75:
			return 75
		case s <= 100:
			return 100
		}
		return 0
	}

	nmAp, err := nmNewAccessPoint(apPath)
	if err != nil {
		return
	}

	ap = accessPoint{
		Ssid:         string(nmAp.Ssid.Get()),
		Secured:      getApSecType(nmAp) != apSecNone,
		SecuredInEap: getApSecType(nmAp) == apSecEap,
		Strength:     calcStrength(nmAp.Strength.Get()),
		Path:         nmAp.Path,
	}
	return
}

func getApSecType(ap *nm.AccessPoint) apSecType {
	return doParseApSecType(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get())
}

func doParseApSecType(flags, wpaFlags, rsnFlags uint32) apSecType {
	r := apSecNone

	if (flags&NM_802_11_AP_FLAGS_PRIVACY != 0) && (wpaFlags == NM_802_11_AP_SEC_NONE) && (rsnFlags == NM_802_11_AP_SEC_NONE) {
		r = apSecWep
	}
	if wpaFlags != NM_802_11_AP_SEC_NONE {
		r = apSecPsk
	}
	if rsnFlags != NM_802_11_AP_SEC_NONE {
		r = apSecPsk
	}
	if (wpaFlags&NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0) || (rsnFlags&NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0) {
		r = apSecEap
	}
	return r
}
