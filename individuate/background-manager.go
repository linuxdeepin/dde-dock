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
	"dlib/dbus/property"
	"dlib/gio-2.0"
	"fmt"
	"math/rand"
	"time"
)

type BackgroundManager struct {
	AutoSwitch     *property.GSettingsBoolProperty   `access:"readwrite"`
	SwitchDuration *property.GSettingsIntProperty    `access:"readwrite"`
	CrossFadeMode  *property.GSettingsStringProperty `access:"readwrite"`
	CurrentPicture *property.GSettingsStringProperty
	PictureURIS    *property.GSettingsStrvProperty
	PictureIndex   *property.GSettingsIntProperty
}

const (
	INDIVIDUATE_ID = "com.deepin.dde.individuate"
)

var (
	indiviGSettings = gio.NewSettings(INDIVIDUATE_ID)
	_switchQuit     chan bool
)

func (bgManager *BackgroundManager) SetBackgroundPicture(uri string, replace bool) {
	var index int
	pictStrv := []string{}

	if replace {
		/* use 'uri' replace 'PictureURIS' */
		pictStrv = append(pictStrv, uri)
		index = 0
	} else {
		/* append 'uri' to 'PictureURIS' */
		pictStrv = bgManager.PictureURIS.GetValue().([]string)
		success, i := IsURIExist(uri, pictStrv)
		if success {
			indiviGSettings.SetInt("index", i)
			return
		}

		fmt.Println("add: strv ", pictStrv)
		index = len(pictStrv)
		fmt.Println("add: len ", index)
		pictStrv = append(pictStrv, uri)
	}

	indiviGSettings.SetStrv("picture-uris", pictStrv)
	indiviGSettings.SetInt("index", index)
}

func (bgManager *BackgroundManager) DeletePictureFromURIS(uri string) {
	if len(uri) <= 0 {
		return
	}

	tempURIS := []string{}
	uris := bgManager.PictureURIS.GetValue().([]string)
	index := int(bgManager.PictureIndex.GetValue().(int32))
	currentURI := bgManager.CurrentPicture.GetValue().(string)

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
		return
	}
	indiviGSettings.SetStrv("picture-uris", tempURIS)

	if uri == currentURI {
		index += 1
		if index > len(tempURIS) {
			index = 0
		}
	} else {
		if success, i := IsURIExist(currentURI, tempURIS); success {
			index = i
		}
	}
	indiviGSettings.SetInt("index", index)
}

func NewBackgroundManager() *BackgroundManager {
	bgManager := &BackgroundManager{}
	_switchQuit = make(chan bool)

	bgManager.AutoSwitch = property.NewGSettingsBoolProperty(
		bgManager, "AutoSwitch",
		indiviGSettings, "auto-switch")
	bgManager.SwitchDuration = property.NewGSettingsIntProperty(
		bgManager, "SwitchDuration",
		indiviGSettings, "background-duration")
	bgManager.CrossFadeMode = property.NewGSettingsStringProperty(
		bgManager, "CrossFadeMode",
		indiviGSettings, "cross-fade-auto-mode")
	bgManager.CurrentPicture = property.NewGSettingsStringProperty(
		bgManager, "CurrentPicture",
		indiviGSettings, "current-picture")
	bgManager.PictureURIS = property.NewGSettingsStrvProperty(
		bgManager, "PictureURIS",
		indiviGSettings, "picture-uris")
	bgManager.PictureIndex = property.NewGSettingsIntProperty(
		bgManager, "PictureIndex",
		indiviGSettings, "index")

	ListenGSetting(bgManager)

	return bgManager
}

func ListenGSetting(bgManager *BackgroundManager) {
	indiviGSettings.Connect("changed::picture-uris", func(s *gio.Settings, key string) {
		/* generate bg blur picture */
	})

	indiviGSettings.Connect("changed::index", func(s *gio.Settings, key string) {
		i := s.GetInt(key)
		uris := s.GetStrv("picture-uris")
		if len(uris) <= 0 {
			s.Reset("current-picture")
			return
		}
		if i > len(uris) {
			i = 0
		}
		fmt.Println("signal: index ", i)
		s.SetString("current-picture", uris[i])
	})

	indiviGSettings.Connect("changed::auto-switch", func(s *gio.Settings, key string) {
		v := s.GetBoolean(key)
		if v {
			go SwitchPictureThread(bgManager)
		} else {
			_switchQuit <- true
		}
	})
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

func SwitchPictureThread(bgManager *BackgroundManager) {
	for {
		secondNums := bgManager.SwitchDuration.GetValue().(time.Duration)
		timer := time.NewTimer(time.Second * secondNums)
		select {
		case <-timer.C:
			AutoSwitchPicture(bgManager)
		case <-_switchQuit:
			return
		}
	}
}

func AutoSwitchPicture(bgManager *BackgroundManager) {
	uris := bgManager.PictureURIS.GetValue().([]string)
	l := len(uris)
	if l <= 1 {
		return
	}
	index := int(bgManager.PictureIndex.GetValue().(int32))

	crossMode := bgManager.CrossFadeMode.GetValue().(string)
	if crossMode == "Sequential" {
		index += 1
		if index > l {
			index = 0
		}
	} else {
		rand.Seed(time.Now().UTC().UnixNano())
		index = rand.Intn(l - 1)
	}
	indiviGSettings.SetInt("index", index)
}
