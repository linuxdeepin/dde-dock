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

func newProfileInfo2(info pulse.ProfileInfo2) *ProfileInfo2 {
	return &ProfileInfo2{
		Name:        info.Name,
		Description: info.Description,
		Priority:    info.Priority,
		Available:   info.Available,
	}
}

type ProfileInfos2 []*ProfileInfo2

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
