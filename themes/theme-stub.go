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
 **/

package main

import (
	"dlib/dbus"
	//"dlib/gio-2.0"
	"math/rand"
	"time"
        "fmt"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		MANAGER_DEST,
		MANAGER_PATH,
		MANAGER_IFC,
	}
}

func (m *Manager) setPropThemeInfo(name string) {
	switch name {
	case "AvailableFontTheme":
		{
			m.AvailableBackground = getBackgroundFiles()
			dbus.NotifyChange(m, name)
		}
	case "AvailableBackground":
		{
			m.AvailableFontTheme = getFontThemes()
			dbus.NotifyChange(m, name)
		}
	case "AvailableIconTheme":
		{
			for _, v := range systemThemes {
				icon := ThemeType{Name: v.IconTheme, Type: "system"}
				m.AvailableIconTheme = append(m.AvailableIconTheme, icon)
			}
			dbus.NotifyChange(m, name)
		}
	case "AvailableGtkTheme":
		{
			for _, v := range systemThemes {
				gtk := ThemeType{Name: v.GtkTheme, Type: "system"}
				m.AvailableGtkTheme = append(m.AvailableGtkTheme, gtk)
			}
			dbus.NotifyChange(m, name)
		}
	case "AvailableCursorTheme":
		{
			for _, v := range systemThemes {
				cursor := ThemeType{Name: v.CursorTheme, Type: "system"}
				m.AvailableCursorTheme = append(m.AvailableCursorTheme, cursor)
			}
			dbus.NotifyChange(m, name)
		}
	case "AvailableWindowTheme":
		{
			for _, v := range systemThemes {
				window := ThemeType{Name: v.WindowTheme, Type: "system"}
				m.AvailableWindowTheme = append(m.AvailableWindowTheme, window)
			}
			dbus.NotifyChange(m, name)
		}
	default:
                fmt.Printf("'%s': invalid theme property\n", name)
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
	//gio.SettingsSync()
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

/*
 * get default picture when picture not exist
 */
func (m *Manager) parseFileNotExist() {
	tmp := []string{}
	uris := indiviGSettings.GetStrv(SCHEMA_KEY_URIS)
	uri := m.BackgroundFile.Get()
	if ok, i := isURIExist(uri, uris); ok {
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

func isURIExist(uri string, uris []string) (bool, int) {
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
