package main

import (
	nm "dbus/org/freedesktop/networkmanager"
	"dlib/dbus"
)

type AccessPoint struct {
	Ssid     string
	NeedKey  bool
	Strength uint8
	Path     dbus.ObjectPath
}

const (
	ApKeyNone = iota
	ApKeyWep
	ApKeyPsk
	ApKeyEap
)

func parseFlags(ap *nm.AccessPoint) int {
	return doParseFlags(ap.Flags.Get(), ap.WpaFlags.Get(), ap.RsnFlags.Get())
}

func doParseFlags(flags, wpaFlags, rsnFlags uint32) int {
	r := ApKeyNone

	if (flags&NM_802_11_AP_FLAGS_PRIVACY != 0) && (wpaFlags == NM_802_11_AP_SEC_NONE) && (rsnFlags == NM_802_11_AP_SEC_NONE) {
		r = ApKeyWep
	}
	if wpaFlags != NM_802_11_AP_SEC_NONE {
		r = ApKeyPsk
	}
	if rsnFlags != NM_802_11_AP_SEC_NONE {
		r = ApKeyPsk
	}
	if (wpaFlags&NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0) || (rsnFlags&NM_802_11_AP_SEC_KEY_MGMT_802_1X != 0) {
		r = ApKeyEap
	}
	return r
}
