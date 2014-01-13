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
	"dlib/dbus/property"
	"dlib/gio-2.0"
	"fmt"
)

func NewManager() *Manager {
	m := &Manager{}

	m.GtkTheme = property.NewGSettingsStringProperty(
		m, "GtkTheme",
		indiviGSettings, SCHEMA_KEY_GTK)
	m.IconTheme = property.NewGSettingsStringProperty(
		m, "IconTheme",
		indiviGSettings, SCHEMA_KEY_ICON)
	m.FontTheme = property.NewGSettingsStringProperty(
		m, "FontTheme",
		indiviGSettings, SCHEMA_KEY_FONT)
	m.CursorTheme = property.NewGSettingsStringProperty(
		m, "CursorTheme",
		indiviGSettings, SCHEMA_KEY_CURSOR)
	m.BackgroundFile = property.NewGSettingsStringProperty(
		m, "BackgroundFile",
		indiviGSettings, SCHEMA_KEY_CUR_PICT)
	m.AutoSwitch = property.NewGSettingsBoolProperty(
		m, "AutoSwitch",
		indiviGSettings, SCHEMA_KEY_AUTO_SWITCH)
	m.SwitchDuration = property.NewGSettingsIntProperty(
		m, "SwitchDuration",
		indiviGSettings, SCHEMA_KEY_DURATION)
	m.CrossFadeMode = property.NewGSettingsStringProperty(
		m, "CrossFadeMode",
		indiviGSettings, SCHEMA_KEY_CROSS_MODE)
	m.CrossInterval = property.NewGSettingsIntProperty(
		m, "CrossInterval",
		indiviGSettings, SCHEMA_KEY_CROSS_INTERVAL)

	m.isAutoSwitch = false
	m.quitAutoSwitch = make(chan bool)

	InitThemeInfo(m)
	ListenSettings(m)

	return m
}

func InitThemeInfo(m *Manager) {
	m.AvailableFontTheme = GetFontThemes()
	m.AvailableBackground = GetBackgroundFiles()

	for _, v := range systemThemes {
		gtk := ThemeType{Name: v.GtkTheme, Type: "system"}
		m.AvailableGtkTheme = append(m.AvailableGtkTheme, gtk)

		icon := ThemeType{Name: v.IconTheme, Type: "system"}
		m.AvailableIconTheme = append(m.AvailableIconTheme, icon)

		cursor := ThemeType{Name: v.CursorTheme, Type: "system"}
		m.AvailableCursorTheme = append(m.AvailableCursorTheme, cursor)

		window := ThemeType{Name: v.WindowTheme, Type: "system"}
		m.AvailableWindowTheme = append(m.AvailableWindowTheme, window)
	}
}

func ListenSettings(m *Manager) {
	indiviGSettings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case SCHEMA_KEY_CUR_PICT:
			{
				if m.isAutoSwitch {
					m.quitAutoSwitch <- true
				}
				uri := s.GetString(SCHEMA_KEY_CUR_PICT)
				filename := GetPathFromURI(uri)
				fmt.Println("\tlisten uri: ", uri)
				fmt.Println("\tlisten path: ", filename)
				isExist := IsFileExist(filename)
				if !isExist {
					ParseFileNotExist(m)
					return
				}
				userManager.BackgroundFile.Set(filename)
				tmp := []string{}
				if m.AutoSwitch.Get() {
					defer func() {
						go SwitchPictureThread(m)
					}()
					uris := s.GetStrv(SCHEMA_KEY_URIS)
					ok, i := IsURIExist(uri, uris)
					if ok {
						s.SetInt(SCHEMA_KEY_INDEX, i)
						return
					}
					tmp = append(tmp, uris...)
				}
				tmp = append(tmp, uri)
				l := len(tmp)
				s.SetStrv(SCHEMA_KEY_URIS, tmp)
				s.SetInt(SCHEMA_KEY_INDEX, l-1)
				break
			}
		case SCHEMA_KEY_AUTO_SWITCH:
			{
				fmt.Println("\tisAutoSwitch: ", m.isAutoSwitch)
				if m.isAutoSwitch {
					m.quitAutoSwitch <- true
				}
				autoSwitch := s.GetBoolean(SCHEMA_KEY_AUTO_SWITCH)
				fmt.Println("\tautoSwitch: ", autoSwitch)
				if autoSwitch {
					go SwitchPictureThread(m)
				}
				break
			}
		case SCHEMA_KEY_URIS:
			{
				/* generate bg blur picture */
				uris := indiviGSettings.GetStrv(SCHEMA_KEY_URIS)
				for _, v := range uris {
					go accountsExtends.BackgroundBlurPictPath(currentUid, GetPathFromURI(v))
				}
				break
			}
		case SCHEMA_KEY_DURATION:
			{
				if m.isAutoSwitch {
					m.quitAutoSwitch <- true
				}
				if m.AutoSwitch.Get() {
					go SwitchPictureThread(m)
				}
				break
			}
		case SCHEMA_KEY_GTK:
			{
				gtk := infaceSettings.GetString(SCHEMA_KEY_GTK)
				if gtk == m.GtkTheme.Get() {
					break
				}
				infaceSettings.SetString(SCHEMA_KEY_GTK,
					m.GtkTheme.Get())
				break
			}
		case SCHEMA_KEY_ICON:
			{
				icon := infaceSettings.GetString(SCHEMA_KEY_ICON)
				if icon == m.IconTheme.Get() {
					break
				}
				infaceSettings.SetString(SCHEMA_KEY_ICON,
					m.IconTheme.Get())
				break
			}
		case SCHEMA_KEY_FONT:
			{
                                font := infaceSettings.GetString(SCHEMA_KEY_FONT)
                                if font == m.FontTheme.Get() {
                                        break
                                }
				break
			}
		case SCHEMA_KEY_CURSOR:
			{
				cursor := infaceSettings.GetString(SCHEMA_KEY_CURSOR)
				if cursor == m.CursorTheme.Get() {
					break
				}
				infaceSettings.SetString(SCHEMA_KEY_CURSOR,
					m.CursorTheme.Get())
				break
			}
		default:
			break
		}
		gio.SettingsSync()
	})

	infaceSettings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case SCHEMA_KEY_GTK:
			{
				gtk := infaceSettings.GetString(SCHEMA_KEY_GTK)
				if gtk == m.GtkTheme.Get() {
					break
				}
				m.GtkTheme.Set(gtk)
				break
			}
		case SCHEMA_KEY_ICON:
			{
				icon := infaceSettings.GetString(SCHEMA_KEY_ICON)
				if icon == m.IconTheme.Get() {
					break
				}
				m.IconTheme.Set(icon)
				break
			}
		case SCHEMA_KEY_FONT:
			{
                                font := infaceSettings.GetString(SCHEMA_KEY_FONT)
                                if font == m.FontTheme.Get() {
                                        break
                                }
                                m.FontTheme.Set(font)
				break
			}
		case SCHEMA_KEY_CURSOR:
			{
				cursor := infaceSettings.GetString(SCHEMA_KEY_CURSOR)
				if cursor == m.CursorTheme.Get() {
					break
				}
				m.CursorTheme.Set(cursor)
				break
			}
		default:
			break
		}
		gio.SettingsSync()
	})
}
