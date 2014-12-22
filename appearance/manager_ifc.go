/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package appearance

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	errInvalidType  = fmt.Errorf("Invalid Key Type")
	errInvalidValue = fmt.Errorf("Invalid Value")
)

func (m *Manager) GetFlag(t, name string) (int32, error) {
	t = strings.ToLower(t)
	switch t {
	case "theme":
		return m.dtheme.GetFlag(name), nil
	case "gtk":
		return m.gtk.GetFlag(name), nil
	case "icon":
		return m.icon.GetFlag(name), nil
	case "cursor":
		return m.cursor.GetFlag(name), nil
	case "sound":
		return m.sound.GetFlag(name), nil
	case "greeter":
		return m.greeter.GetFlag(name), nil
	case "background":
		return m.bg.GetFlag(name), nil
	}

	return -1, errInvalidType
}

// TODO: get i18n name
func (m *Manager) GetDisplayName(t, name string) string {
	return ""
}

func (m *Manager) Set(t, value string) error {
	t = strings.ToLower(t)
	switch t {
	case "theme":
		if !m.dtheme.IsValueValid(value) {
			logger.Debug("Invalid deepin theme:", value)
			return errInvalidValue
		}

		m.settings.SetString(deepinGSKeyTheme, value)
	case "gtk":
		if !m.gtk.IsValueValid(value) {
			logger.Debug("Invalid gtk theme:", value)
			return errInvalidValue
		}

		m.setTheme(t, value)
	case "icon":
		if !m.icon.IsValueValid(value) {
			logger.Debug("Invalid icon theme:", value)
			return errInvalidValue
		}

		m.setTheme(t, value)
	case "cursor":
		if !m.cursor.IsValueValid(value) {
			logger.Debug("Invalid cursor theme:", value)
			return errInvalidValue
		}

		m.setTheme(t, value)
	case "sound":
		if !m.sound.IsValueValid(value) {
			logger.Debug("Invalid sound theme:", value)
			return errInvalidValue
		}

		m.setTheme(t, value)
	case "greeter":
		if !m.greeter.IsValueValid(value) {
			logger.Debug("Invalid greeter theme:", value)
			return errInvalidValue
		}

		m.greeter.Set(value)
	case "background":
		if !m.bg.IsValueValid(value) {
			logger.Debug("Invalid background:", value)
			return errInvalidValue
		}

		m.setTheme(t, value)
	case "font-standard":
		if !m.font.IsStandardFontValid(value) {
			logger.Debug("Invalid font name:", value)
			return errInvalidValue
		}

		m.setTheme(t, value)
	case "font-mono":
		if !m.font.IsMonospacedFontValid(value) {
			logger.Debug("Invalid font monospace:", value)
			return errInvalidValue
		}

		m.setTheme(t, value)
	case "font-size":
		fontSize, _ := strconv.ParseInt(value, 10, 64)
		if !m.font.IsFontSizeValid(int32(fontSize)) {
			logger.Debug("Invalid font size:", value)
			return errInvalidValue
		}

		m.setTheme(t, int32(fontSize))
	default:
		return errInvalidType
	}

	return nil
}

func (m *Manager) GetThumbnail(t, name string) string {
	t = strings.ToLower(t)
	switch t {
	case "theme":
		if name == customDThemeId {
			t, ok := m.themeObjMap[name]
			if !ok {
				return ""
			}
			return m.bg.GetThumbnail(t.Background)
		}

		return m.dtheme.GetThumbnail(name)
	case "gtk":
		return m.gtk.GetThumbnail(name)
	case "icon":
		return m.icon.GetThumbnail(name)
	case "cursor":
		return m.cursor.GetThumbnail(name)
	case "greeter":
		return m.greeter.GetThumbnail(name)
	case "background":
		return m.bg.GetThumbnail(name)
	}

	return ""
}

func (m *Manager) GetFontList(t string) []string {
	return m.font.GetNameList(t)
}

func (m *Manager) GetFontStyles(name string) []string {
	return m.font.GetStyleListByName(name)
}

func (m *Manager) Delete(t, name string) {
	t = strings.ToLower(t)
	switch t {
	case "theme":
		m.dtheme.Delete(name)
	case "gtk":
		m.gtk.Delete(name)
	case "icon":
		m.icon.Delete(name)
	case "cursor":
		m.cursor.Delete(name)
	case "sound":
		m.sound.Delete(name)
	case "greeter":
		m.greeter.Delete(name)
	case "background":
		m.bg.Delete(name)
	}
}
