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

package appearance

import (
	"fmt"
	"strconv"
	"strings"

	"pkg.deepin.io/dde/daemon/appearance/fonts"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"

	dutils "pkg.deepin.io/lib/utils"
)

// Reset reset all themes and fonts settings to default values
func (m *Manager) Reset() *dbus.Error {
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
	return nil
}

// List list all available for the special type, return a json format list
func (m *Manager) List(ty string) (string, *dbus.Error) {
	logger.Debug("List for type:", ty)
	jsonStr, err := m.list(ty)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	return jsonStr, nil
}

func (m *Manager) list(ty string) (string, error) {
	switch strings.ToLower(ty) {
	case TypeGtkTheme:
		return m.doShow(subthemes.ListGtkTheme())
	case TypeIconTheme:
		return m.doShow(subthemes.ListIconTheme())
	case TypeCursorTheme:
		return m.doShow(subthemes.ListCursorTheme())
	case TypeBackground:
		return m.doShow(m.listBackground())
	case TypeStandardFont:
		return m.doShow(fonts.GetFamilyTable().ListStandard())
	case TypeMonospaceFont:
		return m.doShow(fonts.GetFamilyTable().ListMonospace())
	}
	return "", fmt.Errorf("invalid type: %v", ty)

}

// Show show detail infos for the special type
// ret0: detail info, json format
func (m *Manager) Show(ty string, names []string) (string, *dbus.Error) {
	logger.Debugf("Show '%s' type '%s'", names, ty)
	jsonStr, err := m.show(ty, names)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	return jsonStr, nil
}

func (m *Manager) show(ty string, names []string) (string, error) {
	switch strings.ToLower(ty) {
	case TypeGtkTheme:
		return m.doShow(subthemes.ListGtkTheme().ListGet(names))
	case TypeIconTheme:
		return m.doShow(subthemes.ListIconTheme().ListGet(names))
	case TypeCursorTheme:
		return m.doShow(subthemes.ListCursorTheme().ListGet(names))
	case TypeBackground:
		return m.doShow(m.listBackground().ListGet(names))
	case TypeStandardFont, TypeMonospaceFont:
		return m.doShow(fonts.GetFamilyTable().GetFamilies(names))
	}
	return "", fmt.Errorf("invalid type: %v", ty)
}

// Set set to the special 'value'
func (m *Manager) Set(ty, value string) *dbus.Error {
	logger.Debugf("Set '%s' for type '%s'", value, ty)
	err := m.set(ty, value)
	return dbusutil.ToError(err)
}

func (m *Manager) set(ty, value string) error {
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
		file, err := m.doSetBackground(value)
		if err == nil {
			m.wsLoop.AddToShowed(file)
		}
	case TypeGreeterBackground:
		err = m.doSetGreeterBackground(value)
		m.currentGreeterBg = dutils.EncodeURI(value, dutils.SCHEME_FILE)
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
		err = m.doSetMonospaceFont(value)
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
		return fmt.Errorf("invalid type: %v", ty)
	}
	return err
}

// Delete delete the special 'name'
func (m *Manager) Delete(ty, name string) *dbus.Error {
	logger.Debugf("Delete '%s' type '%s'", name, ty)
	err := m.delete(ty, name)
	return dbusutil.ToError(err)
}

func (m *Manager) delete(ty, name string) error {
	switch strings.ToLower(ty) {
	case TypeGtkTheme:
		return subthemes.ListGtkTheme().Delete(name)
	case TypeIconTheme:
		return subthemes.ListIconTheme().Delete(name)
	case TypeCursorTheme:
		return subthemes.ListCursorTheme().Delete(name)
	case TypeBackground:
		return m.listBackground().Delete(name)
		//case TypeStandardFont:
		//case TypeMonospaceFont:
	}
	return fmt.Errorf("invalid type: %v", ty)
}

// Thumbnail get thumbnail for the special 'name'
func (m *Manager) Thumbnail(ty, name string) (string, *dbus.Error) {
	file, err := m.thumbnail(ty, name)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	return file, nil
}

func (m *Manager) thumbnail(ty, name string) (string, error) {
	logger.Debugf("Get thumbnail for '%s' type '%s'", name, ty)
	switch strings.ToLower(ty) {
	case TypeGtkTheme:
		return subthemes.GetGtkThumbnail(name)
	case TypeIconTheme:
		return subthemes.GetIconThumbnail(name)
	case TypeCursorTheme:
		return subthemes.GetCursorThumbnail(name)
	}
	return "", fmt.Errorf("invalid type: %v", ty)
}

func (m *Manager) GetScaleFactor() (float64, *dbus.Error) {
	return m.getScaleFactor(), nil
}

func (m *Manager) SetScaleFactor(scale float64) *dbus.Error {
	err := m.setScaleFactor(scale)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}
	return nil
}

func (m *Manager) SetScreenScaleFactors(v map[string]float64) *dbus.Error {
	err := m.setScreenScaleFactors(v)
	if err != nil {
		logger.Warning(err)
	}
	return dbusutil.ToError(err)
}

func (m *Manager) GetScreenScaleFactors() (map[string]float64, *dbus.Error) {
	v, err := m.getScreenScaleFactors()
	return v, dbusutil.ToError(err)
}
