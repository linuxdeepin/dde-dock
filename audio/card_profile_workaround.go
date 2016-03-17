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
	"encoding/json"
	"fmt"
	"pkg.deepin.io/lib/pulse"
	"sort"
	"sync"
)

const (
	CardBuildin   = 0
	CardBluethooh = 1
	CardUnknow    = 2
)

type ProfileInfo2 struct {
	Name        string
	Description string

	// The higher this value is, the more useful this profile is as a default.
	Priority uint32

	// 如果值是 0, 表示这个配置不可用，无法被激活
	// 如果值不为 0, 也不能保证此配置是可用的，它仅仅意味着不能肯定它是不可用的
	Available int
}
type ProfileInfos2 []*ProfileInfo2

type CardInfo struct {
	Id uint32

	Name string

	ActiveProfile *ProfileInfo2
	Profiles      ProfileInfos2
}
type CardInfos []*CardInfo

var (
	cardLocker sync.Mutex
)

func newProfileInfos2(infos []pulse.ProfileInfo2) ProfileInfos2 {
	var pinfos ProfileInfos2
	for _, v := range infos {
		pinfos = append(pinfos, newProfileInfo2(v))
	}
	return pinfos
}

func (infos ProfileInfos2) get(name string) (*ProfileInfo2, error) {
	for _, info := range infos {
		if info.Name == name {
			return info, nil
		}
	}
	return nil, fmt.Errorf("Invalid profile name: %v", name)
}

func newProfileInfo2(info pulse.ProfileInfo2) *ProfileInfo2 {
	return &ProfileInfo2{
		Name:        info.Name,
		Description: info.Description,
		Priority:    info.Priority,
		Available:   info.Available,
	}
}

func newCardInfo(card *pulse.Card) *CardInfo {
	var info = new(CardInfo)
	info.update(card)
	return info
}

func (info *CardInfo) update(card *pulse.Card) {
	info.Id = card.Index
	info.Name = card.Name
	info.ActiveProfile = newProfileInfo2(card.ActiveProfile)
	sort.Sort(cProfileInfos2(card.Profiles))
	info.Profiles = newProfileInfos2(card.Profiles)
	info.filterProfile(card)
}

func (info *CardInfo) filterProfile(card *pulse.Card) {
	var profiles ProfileInfos2
	blacklist := profileBlacklist(card)
	for _, p := range info.Profiles {
		_, ok := blacklist[p.Name]
		if ok || p.Available == 0 {
			continue
		}
		profiles = append(profiles, p)
	}
	info.Profiles = profiles
}

func newCardInfos(cards []*pulse.Card) CardInfos {
	var infos CardInfos
	cardLocker.Lock()
	defer cardLocker.Unlock()
	for _, v := range cards {
		infos = append(infos, newCardInfo(v))
	}
	return infos
}

func (infos CardInfos) string() string {
	return toJSON(infos)
}

func (infos CardInfos) get(id uint32) (*CardInfo, error) {
	cardLocker.Lock()
	defer cardLocker.Unlock()
	for _, info := range infos {
		if info.Id == id {
			return info, nil
		}
	}
	return nil, fmt.Errorf("Invalid card id: %v", id)
}

func (infos CardInfos) add(info *CardInfo) (CardInfos, bool) {
	tmp, _ := infos.get(info.Id)
	if tmp != nil {
		return infos, false
	}

	cardLocker.Lock()
	defer cardLocker.Unlock()
	infos = append(infos, info)
	return infos, true
}

func (infos CardInfos) delete(id uint32) (CardInfos, bool) {
	var (
		ret     CardInfos
		deleted bool
	)
	for _, info := range infos {
		if info.Id == id {
			deleted = true
			continue
		}
		ret = append(ret, info)
	}
	return ret, deleted
}

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

type cProfileInfos2 []pulse.ProfileInfo2

func (infos cProfileInfos2) exist(name string) bool {
	for _, info := range infos {
		if info.Name == name {
			return true
		}
	}
	return false
}

func (infos cProfileInfos2) Len() int {
	return len(infos)
}

func (infos cProfileInfos2) Less(i, j int) bool {
	return infos[i].Priority > infos[j].Priority
}

func (infos cProfileInfos2) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}

//select New Card Profile By priority, protocl.
func selectNewCardProfile(c *pulse.Card) {
	blacklist := profileBlacklist(c)
	if blacklist[c.ActiveProfile.Name] != "true" {
		logger.Debug("use profile:", c.ActiveProfile)
		return
	}

	var profiles cProfileInfos2
	for _, p := range c.Profiles {
		if "true" == blacklist[p.Name] {
			continue
		}
		profiles = append(profiles, p)
	}

	// sort profiles by priority
	logger.Debug("[selectNewCardProfile] before sort:", profiles)
	sort.Sort(profiles)
	logger.Debug("[selectNewCardProfile] after sort:", profiles)

	if len(profiles) > 0 {
		logger.Debug("re-select card profile:", profiles[0])
		if c.ActiveProfile.Name != profiles[0].Name {
			c.SetProfile(profiles[0].Name)
		}
		return
	}
}

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}
