/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package appearance

import (
	"fmt"
	"strconv"
	"strings"

	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/appearance/dtheme"
	"pkg.deepin.io/dde/daemon/appearance/fonts"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
	"pkg.deepin.io/lib/dbus"
)

// List list all available for the special type
func (m *Manager) List(ty string) (string, error) {
	logger.Debug("List for type:", ty)
	switch strings.ToLower(ty) {
	case TypeDTheme:
		return m.doShow(dtheme.ListDTheme())
	case TypeGtkTheme:
		return m.doShow(subthemes.ListGtkTheme())
	case TypeIconTheme:
		return m.doShow(subthemes.ListIconTheme())
	case TypeCursorTheme:
		return m.doShow(subthemes.ListCursorTheme())
	case TypeBackground:
		return m.doShow(background.ListBackground())
	case TypeStandardFont:
		return m.doShow(fonts.ListStandardFamily())
	case TypeMonospaceFont:
		return m.doShow(fonts.ListMonospaceFamily())
	}
	return "", fmt.Errorf("Invalid type: %v", ty)
}

// Show show detail info for the special type
// ret0: detail info, json format
func (m *Manager) Show(ty, name string) (string, error) {
	logger.Debugf("Show '%s' type '%s'", name, ty)
	switch strings.ToLower(ty) {
	case TypeDTheme:
		return m.doShow(dtheme.ListDTheme().Get(name))
	case TypeGtkTheme:
		return m.doShow(subthemes.ListGtkTheme().Get(name))
	case TypeIconTheme:
		return m.doShow(subthemes.ListIconTheme().Get(name))
	case TypeCursorTheme:
		return m.doShow(subthemes.ListCursorTheme().Get(name))
	case TypeBackground:
		return m.doShow(background.ListBackground().Get(name))
	case TypeStandardFont:
		return m.doShow(fonts.ListStandardFamily().Get(name))
	case TypeMonospaceFont:
		return m.doShow(fonts.ListMonospaceFamily().Get(name))
	}
	return "", fmt.Errorf("Invalid type: %v", ty)
}

// Set set to the special 'value'
func (m *Manager) Set(ty, value string) error {
	logger.Debugf("Set '%s' for type '%s'", value, ty)
	var err error
	switch strings.ToLower(ty) {
	case TypeDTheme:
		err = m.doSetDTheme(value)
	case TypeGtkTheme:
		err = m.doSetGtkTheme(value)
	case TypeIconTheme:
		err = m.doSetIconTheme(value)
	case TypeCursorTheme:
		err = m.doSetCursorTheme(value)
	case TypeBackground:
		err = m.doSetBackground(value)
	case TypeStandardFont:
		err = m.doSetStandardFont(value)
	case TypeMonospaceFont:
		err = m.doSetMonnospaceFont(value)
	case TypeFontSize:
		size, e := strconv.ParseInt(value, 10, 64)
		if e != nil {
			return e
		}
		err = m.doSetFontSize(int32(size))
	default:
		return fmt.Errorf("Invalid type: %v", ty)
	}

	if err != nil {
		return err
	}

	// Emit theme changed signal
	dbus.Emit(m, "Changed", ty, value)
	return nil
}

// Delete delete the special 'name'
func (m *Manager) Delete(ty, name string) error {
	logger.Debugf("Delete '%s' type '%s'", name, ty)
	switch strings.ToLower(ty) {
	case TypeDTheme:
		return dtheme.ListDTheme().Delete(name)
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
	case TypeDTheme:
		return dtheme.GetDThemeThumbnail(name)
	case TypeGtkTheme:
		return subthemes.GetGtkThumbnail(name)
	case TypeIconTheme:
		return subthemes.GetIconThumbnail(name)
	case TypeCursorTheme:
		return subthemes.GetCursorThumbnail(name)
	case TypeBackground:
		return background.ListBackground().Thumbnail(name)
		//case TypeStandardFont:
		//case TypeMonospaceFont:
	}
	return "", fmt.Errorf("Invalid type: %v", ty)
}
