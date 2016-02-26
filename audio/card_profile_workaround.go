/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import (
	"pkg.deepin.io/lib/pulse"
)

const (
	CardBuildin   = 0
	CardBluethooh = 1
	CardUnknow    = 2
)

func cardType(c *pulse.Card) int {
	if c.PropList[PropDeviceFromFactor] == "internal" {
		return CardBuildin
	}
	if c.PropList[PropDeviceBus] == "bluetooth" {
		return CardBluethooh
	}
	return CardUnknow
}

func profileBlacklist(c *pulse.Card) map[string]string {
	switch cardType(c) {
	case CardBluethooh:
		// TODO: bluez not full support headset_head_unit, please skip
		return map[string]string{"off": "true", "headset_head_unit": "true"}
	case CardBuildin, CardUnknow:
		fallthrough
	default:
		return map[string]string{"off": "true"}
	}
}

//select New Card Profile By priority, protocl.
func selectNewCardProfile(c *pulse.Card) {
	profiles := []pulse.ProfileInfo2{}
	blacklist := profileBlacklist(c)
	if blacklist[c.ActiveProfile.Name] != "true" {
		logger.Info("use profile:", c.ActiveProfile)
		return
	}

	for _, p := range c.Profiles {
		if "true" == blacklist[p.Name] {
			continue
		}
		profiles = append(profiles, p)
	}

	//TODO: sort profiles by priority

	if len(profiles) > 0 {
		logger.Info("re-select card profile:", profiles[0])
		c.SetProfile(profiles[0].Name)
		return
	}
}
