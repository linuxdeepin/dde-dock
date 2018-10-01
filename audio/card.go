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

package audio

import (
	"fmt"
	"sort"

	"pkg.deepin.io/lib/pulse"
)

type Card struct {
	Id            uint32
	Name          string
	ActiveProfile *Profile
	Profiles      ProfileList
	Ports         pulse.CardPortInfos
	core          *pulse.Card
}

type CardExport struct {
	Id    uint32
	Name  string
	Ports []CardPortExport
}

type CardPortExport struct {
	Name        string
	Description string
	Direction   int
}

func newCard(card *pulse.Card) *Card {
	var info = new(Card)
	info.core = card
	info.update(card)
	return info
}

func (c *Card) update(card *pulse.Card) {
	c.Id = card.Index
	propDeviceProductName := card.PropList["device.product.name"]
	propAlsaCardName := card.PropList["alsa.card_name"]
	propDeviceDescription := card.PropList["device.description"]
	if propDeviceProductName != "" {
		c.Name = propDeviceProductName
	} else if propAlsaCardName != "" {
		c.Name = propAlsaCardName
	} else if propDeviceDescription != "" {
		c.Name = propDeviceDescription
	} else {
		c.Name = card.Name
	}

	c.ActiveProfile = newProfile(card.ActiveProfile)
	sort.Sort(card.Profiles)
	c.Profiles = newProfileList(card.Profiles)
	c.filterProfile(card)
	c.Ports = card.Ports
}

func (c *Card) tryGetProfileByPort(portName string) (string, error) {
	profile, _ := c.Ports.TrySelectProfile(portName)
	if len(profile) == 0 {
		return "", fmt.Errorf("not found profile for port '%s'", portName)
	}
	return profile, nil
}

func (c *Card) filterProfile(card *pulse.Card) {
	var profiles ProfileList
	blacklist := profileBlacklist(card)
	for _, p := range c.Profiles {
		// skip unavailable and blacklisted profiles
		if p.Available == 0 || blacklist.Contains(p.Name) {
			// TODO : p.Available == 0 ?
			continue
		}
		profiles = append(profiles, p)
	}
	c.Profiles = profiles
}

type CardList []*Card

func newCardList(cards []*pulse.Card) CardList {
	var result CardList
	for _, v := range cards {
		result = append(result, newCard(v))
	}
	return result
}

func (cl CardList) string() string {
	var list []CardExport
	for _, cardInfo := range cl {
		var ports []CardPortExport
		for _, portInfo := range cardInfo.Ports {
			ports = append(ports, CardPortExport{
				Name:        portInfo.Name,
				Description: portInfo.Description,
				Direction:   portInfo.Direction,
			})
		}

		list = append(list, CardExport{
			Id:    cardInfo.Id,
			Name:  cardInfo.Name,
			Ports: ports,
		})
	}
	return toJSON(list)
}

func (cl CardList) get(id uint32) (*Card, error) {
	for _, info := range cl {
		if info.Id == id {
			return info, nil
		}
	}
	return nil, fmt.Errorf("invalid card id: %v", id)
}

func (cl CardList) add(info *Card) (CardList, bool) {
	card, _ := cl.get(info.Id)
	if card != nil {
		return cl, false
	}

	return append(cl, info), true
}

func (cl CardList) delete(id uint32) (CardList, bool) {
	var (
		ret     CardList
		deleted bool
	)
	for _, info := range cl {
		if info.Id == id {
			deleted = true
			continue
		}
		ret = append(ret, info)
	}
	return ret, deleted
}

func (cl CardList) getAvailablePort(direction int) (uint32, pulse.CardPortInfo) {
	var (
		idY   uint32
		id    uint32
		portY pulse.CardPortInfo // Yes state
		port  pulse.CardPortInfo
	)
	for _, info := range cl {
		v := hasPortAvailable(info.Ports, direction, true)
		if v.Name != "" {
			if portY.Priority < v.Priority || portY.Name == "" {
				portY = v
				idY = info.Id
			}
			continue
		}

		vv := hasPortAvailable(info.Ports, direction, false)
		if port.Priority < vv.Priority || port.Name == "" {
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
		portU pulse.CardPortInfo // Unknown state
	)
	for _, v := range infos {
		if v.Direction != direction {
			continue
		}

		if v.Available == pulse.AvailableTypeYes {
			if portY.Priority < v.Priority || portY.Name == "" {
				portY = v
			}
		} else if v.Available == pulse.AvailableTypeUnknow {
			if portU.Priority < v.Priority || portU.Name == "" {
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
