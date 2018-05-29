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

type Profile struct {
	Name        string
	Description string

	// The higher this value is, the more useful this profile is as a default.
	Priority uint32

	// 如果值是 0, 表示这个配置不可用，无法被激活
	// 如果值不为 0, 也不能保证此配置是可用的，它仅仅意味着不能肯定它是不可用的
	Available int
}

func newProfile(info pulse.ProfileInfo2) *Profile {
	return &Profile{
		Name:        info.Name,
		Description: info.Description,
		Priority:    info.Priority,
		Available:   info.Available,
	}
}

type ProfileList []*Profile

func newProfileList(src []pulse.ProfileInfo2) ProfileList {
	var result ProfileList
	for _, v := range src {
		result = append(result, newProfile(v))
	}
	return result
}

func (pl ProfileList) get(name string) (*Profile, error) {
	for _, info := range pl {
		if info.Name == name {
			return info, nil
		}
	}
	return nil, fmt.Errorf("invalid profile name: %v", name)
}

func getCommonProfiles(info1, info2 pulse.CardPortInfo) pulse.ProfileInfos2 {
	var commons pulse.ProfileInfos2
	if len(info1.Profiles) == 0 || len(info2.Profiles) == 0 {
		return commons
	}
	for _, profile := range info1.Profiles {
		if !info2.Profiles.Exists(profile.Name) {
			continue
		}
		commons = append(commons, profile)
	}
	if len(commons) != 0 {
		sort.Sort(commons)
	}
	return commons
}
