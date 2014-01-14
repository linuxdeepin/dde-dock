/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 *
 * Function: Manager Background switch/add/delete etc...
 **/

package main

import (
	"dlib/gio-2.0"
	"fmt"
	"math/rand"
	"time"
)

func (m *Manager) DeletePictureFromURIS(uri string) {
	if len(uri) <= 0 {
		return
	}

	tempURIS := []string{}
	uris := indiviGSettings.GetStrv(SCHEMA_KEY_URIS)
	index := indiviGSettings.GetInt(SCHEMA_KEY_INDEX)
	currentURI := m.BackgroundFile.Get()

	fmt.Println("del: uris ", uris)
	for _, v := range uris {
		if v != uri {
			tempURIS = append(tempURIS, v)
		}
	}

	fmt.Println("del: tmp ", tempURIS)
	if len(tempURIS) <= 0 {
		indiviGSettings.Reset("picture-uris")
		indiviGSettings.SetInt("index", 0)
		m.BackgroundFile.Set(tempURIS[0])
		return
	}
	indiviGSettings.SetStrv("picture-uris", tempURIS)

	if uri == currentURI {
		index += 1
		if index > len(tempURIS) {
			index = 0
		}
		m.BackgroundFile.Set(tempURIS[index])
	} else {
		if success, i := IsURIExist(currentURI, tempURIS); success {
			index = i
		}
	}
	indiviGSettings.SetInt("index", index)
}

func IsURIExist(uri string, uris []string) (bool, int) {
	if len(uris) <= 0 {
		return false, -1
	}

	for i, v := range uris {
		if v == uri {
			return true, i
		}
	}

	return false, -1
}

func (m *Manager) switchPictureThread() {
	m.isAutoSwitch = true
	for {
		secondNums := m.SwitchDuration.Get()
		timer := time.NewTimer(time.Second * time.Duration(secondNums))
		select {
		case <-timer.C:
			m.autoSwitchPicture()
		case <-m.quitAutoSwitch:
			m.isAutoSwitch = false
			return
		}
	}
}

func (m *Manager) autoSwitchPicture() {
	uris := indiviGSettings.GetStrv(SCHEMA_KEY_URIS)
	l := len(uris)
	if l <= 1 {
		return
	}
	index := int(indiviGSettings.GetInt(SCHEMA_KEY_INDEX))

	/*fmt.Println("\nAutoSwitchPicture...")*/
	//fmt.Println("\turis: ", uris)
	//fmt.Println("\tlen: ", l)
	/*fmt.Println("\tindex: ", index)*/

	crossMode := m.CrossFadeMode.Get()
	//fmt.Println("\tmode: ", crossMode)
	if crossMode == "Sequential" {
		index += 1
		if index >= l {
			index = 0
		}
		fmt.Println("\tSequential index: ", index)
	} else {
		rand.Seed(time.Now().UTC().UnixNano())
		index = rand.Intn(l - 1)
		fmt.Println("\tOther index: ", index)
	}
	m.BackgroundFile.Set(uris[index])
	//fmt.Println("\turi: ", uris[index])
	indiviGSettings.SetInt(SCHEMA_KEY_INDEX, index)
	gio.SettingsSync()
}

/*
 * get default picture when picture not exist
 */
func (m *Manager) parseFileNotExist() {
	tmp := []string{}
	uris := indiviGSettings.GetStrv(SCHEMA_KEY_URIS)
	uri := m.BackgroundFile.Get()
	if ok, i := IsURIExist(uri, uris); ok {
		for j, v := range uris {
			if j == i {
				continue
			}
			tmp = append(tmp, v)
		}
	}
	l := len(tmp)
	if l <= 0 {
		tmp = []string{DEFAULT_BG_PICTURE}
	}
	indiviGSettings.SetStrv(SCHEMA_KEY_URIS, tmp)
	m.BackgroundFile.Set(tmp[0])
	indiviGSettings.SetInt(SCHEMA_KEY_INDEX, 0)
}
