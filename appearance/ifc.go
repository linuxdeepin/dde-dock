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

package appearance

import (
	"fmt"
	"strconv"
	"strings"

	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/appearance/fonts"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
)

// Reset reset all themes and fonts settings to default values
func (m *Manager) Reset() {
	logger.Debug("Reset settings")

	var settingKeys = []string{
		gsKeyGtkTheme,
		gsKeyIconTheme,
		gsKeyCursorTheme,
		gsKeyFontSize,
	}
	for _, key := range settingKeys {
		userVal := m.setting.GetUserValue(key)
		if userVal != nil {
			logger.Debug("reset setting", key)
			m.setting.Reset(key)
		}
	}

	m.resetFonts()
}

// List list all available for the special type, return a json format list
func (m *Manager) List(ty string) (string, error) {
	logger.Debug("List for type:", ty)
	switch strings.ToLower(ty) {
	case TypeGtkTheme:
		return m.doShow(subthemes.ListGtkTheme())
	case TypeIconTheme:
		return m.doShow(subthemes.ListIconTheme())
	case TypeCursorTheme:
		return m.doShow(subthemes.ListCursorTheme())
	case TypeBackground:
		return m.doShow(background.ListBackground())
	case TypeStandardFont:
		return m.doShow(fonts.GetFamilyTable().ListStandard())
	case TypeMonospaceFont:
		return m.doShow(fonts.GetFamilyTable().ListMonospace())
	}
	return "", fmt.Errorf("Invalid type: %v", ty)
}

// Show show detail infos for the special type
// ret0: detail info, json format
func (m *Manager) Show(ty string, names []string) (string, error) {
	logger.Debugf("Show '%s' type '%s'", names, ty)
	switch strings.ToLower(ty) {
	case TypeGtkTheme:
		return m.doShow(subthemes.ListGtkTheme().ListGet(names))
	case TypeIconTheme:
		return m.doShow(subthemes.ListIconTheme().ListGet(names))
	case TypeCursorTheme:
		return m.doShow(subthemes.ListCursorTheme().ListGet(names))
	case TypeBackground:
		return m.doShow(background.ListBackground().ListGet(names))
	case TypeStandardFont, TypeMonospaceFont:
		return m.doShow(fonts.GetFamilyTable().GetFamilies(names))
	}
	return "", fmt.Errorf("Invalid type: %v", ty)
}

// Set set to the special 'value'
func (m *Manager) Set(ty, value string) error {
	logger.Debugf("Set '%s' for type '%s'", value, ty)
	var err error
	switch strings.ToLower(ty) {
	case TypeGtkTheme:
		if m.GtkTheme.Get() == value {
			return nil
		}
		err = m.doSetGtkTheme(value)
		if err == nil {
			m.GtkTheme.Set(value)
		}
	case TypeIconTheme:
		if m.IconTheme.Get() == value {
			return nil
		}
		err = m.doSetIconTheme(value)
		if err == nil {
			m.IconTheme.Set(value)
		}
	case TypeCursorTheme:
		if m.CursorTheme.Get() == value {
			return nil
		}
		err = m.doSetCursorTheme(value)
		if err == nil {
			m.CursorTheme.Set(value)
		}
	case TypeBackground:
		if m.Background.Get() == value {
			return nil
		}
		var uri string
		uri, err = m.doSetBackground(value)
		if err == nil && uri != m.Background.Get() {
			m.Background.Set(uri)
		}
	case TypeGreeterBackground:
		err = m.doSetGreeterBackground(value)
	case TypeStandardFont:
		if m.StandardFont.Get() == value {
			return nil
		}
		err = m.doSetStandardFont(value)
		if err == nil {
			m.StandardFont.Set(value)
		}
	case TypeMonospaceFont:
		if m.MonospaceFont.Get() == value {
			return nil
		}
		err = m.doSetMonnospaceFont(value)
		if err == nil {
			m.MonospaceFont.Set(value)
		}
	case TypeFontSize:
		size, e := strconv.ParseFloat(value, 64)
		if e != nil {
			return e
		}

		cur := m.FontSize.Get()
		if cur > size-0.01 && cur < size+0.01 {
			return nil
		}
		err = m.doSetFontSize(size)
		if err == nil {
			m.FontSize.Set(size)
		}
	default:
		return fmt.Errorf("Invalid type: %v", ty)
	}
	return err
}

// Delete delete the special 'name'
func (m *Manager) Delete(ty, name string) error {
	logger.Debugf("Delete '%s' type '%s'", name, ty)
	switch strings.ToLower(ty) {
	case TypeGtkTheme:
		return subthemes.ListGtkTheme().Delete(name)
	case TypeIconTheme:
		return subthemes.ListIconTheme().Delete(name)
	case TypeCursorTheme:
		return subthemes.ListCursorTheme().Delete(name)
	case TypeBackground:
		return background.ListBackground().Delete(name)
		//case TypeStandardFont:
		//case TypeMonospaceFont:
	}
	return fmt.Errorf("Invalid type: %v", ty)
}

// Thumbnail get thumbnail for the special 'name'
func (m *Manager) Thumbnail(ty, name string) (string, error) {
	logger.Debugf("Get thumbnail for '%s' type '%s'", name, ty)
	switch strings.ToLower(ty) {
	case TypeGtkTheme:
		return subthemes.GetGtkThumbnail(name)
	case TypeIconTheme:
		return subthemes.GetIconThumbnail(name)
	case TypeCursorTheme:
		return subthemes.GetCursorThumbnail(name)
	}
	return "", fmt.Errorf("Invalid type: %v", ty)
}

func (m *Manager) GetScaleFactor() float64 {
	return doGetScaleFactor()
}

func (m *Manager) SetScaleFactor(scale float64) error {
	if scale < 0 {
		return fmt.Errorf("Invalid scale value: %v", scale)
	}

	if getScaleStatus() {
		return fmt.Errorf("There is a scale job running.")
	}

	setScaleStatus(true)
	doSetScaleFactor(scale)
	return nil
}
