/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package audio

import (
	"fmt"
	"pkg.deepin.io/lib/pulse"
	"sort"
	"sync"
)

var (
	cardLocker sync.Mutex
)

type CardInfo struct {
	Id            uint32
	Name          string
	ActiveProfile *ProfileInfo2
	Profiles      ProfileInfos2
	Ports         pulse.CardPortInfos
	core          *pulse.Card
}

func newCardInfo(card *pulse.Card) *CardInfo {
	var info = new(CardInfo)
	info.core = card
	info.update(card)
	return info
}

func (info *CardInfo) update(card *pulse.Card) {
	info.Id = card.Index
	propAlsaCardName := card.PropList["alsa.card_name"]
	if propAlsaCardName != "" {
		info.Name = propAlsaCardName
	} else {
		info.Name = card.Name
	}

	info.ActiveProfile = newProfileInfo2(card.ActiveProfile)
	sort.Sort(cProfileInfos2(card.Profiles))
	info.Profiles = newProfileInfos2(card.Profiles)
	info.filterProfile(card)
	info.Ports = card.Ports
}

func (info *CardInfo) tryGetProfileByPort(portName string) (string, error) {
	profile, _ := info.Ports.TrySelectProfile(portName)
	if len(profile) == 0 {
		return "", fmt.Errorf("Not found profile for port '%s'", portName)
	}
	return profile, nil
}

func (info *CardInfo) filterProfile(card *pulse.Card) {
	var profiles ProfileInfos2
	blacklist := profileBlacklist(card)
	for _, p := range info.Profiles {
		// skip unavailable and blacklisted profiles
		if p.Available == 0 || blacklist.Contains(p.Name) {
			continue
		}
		profiles = append(profiles, p)
	}
	info.Profiles = profiles
}

type CardInfos []*CardInfo

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

func (infos CardInfos) getAvailablePort(direction int) (uint32, pulse.CardPortInfo) {
	var (
		idY   uint32
		id    uint32
		portY pulse.CardPortInfo // Yes state
		port  pulse.CardPortInfo
	)
	for _, info := range infos {
		v := hasPortAvailable(info.Ports, direction, true)
		if v.Name != "" {
			if v.Priority > portY.Priority {
				portY = v
				idY = info.Id
			}
			continue
		}

		vv := hasPortAvailable(info.Ports, direction, false)
		if vv.Priority > port.Priority {
			port = vv
			id = info.Id
		}
	}

	if portY.Name != "" {
		return idY, portY
	}
	return id, port
}

func hasPortAvailable(infos pulse.CardPortInfos, direction int, onlyYes bool) pulse.CardPortInfo {
	var (
		portY pulse.CardPortInfo // Yes state
		portU pulse.CardPortInfo // Unknow state
	)
	for _, v := range infos {
		if v.Direction != direction {
			continue
		}

		if v.Available == pulse.AvailableTypeYes {
			if portY.Priority < v.Priority {
				portY = v
			}
		} else if v.Available == pulse.AvailableTypeUnknow {
			if portU.Priority < v.Priority {
				portU = v
			}
		}
	}

	if onlyYes {
		return portY
	}

	if portY.Name != "" {
		return portY
	}
	return portU
}

func hasPortChanged(oldInfos, newInfos pulse.CardPortInfos) (old, port pulse.CardPortInfo) {
	if len(oldInfos) != len(newInfos) {
		return
	}
	for _, v := range oldInfos {
		tmp, _ := newInfos.Get(v.Name, v.Direction)
		if v.Available == tmp.Available {
			continue
		}

		if port.Priority < tmp.Priority {
			old = v
			port = tmp
		}
	}
	return old, port
}
