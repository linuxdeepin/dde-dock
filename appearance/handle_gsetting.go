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

	"pkg.deepin.io/lib/gsettings"
)

func (m *Manager) listenGSettingChanged() {
	gsettings.ConnectChanged(appearanceSchema, "*", func(key string) {
		if m.setting == nil {
			return
		}

		var (
			ty    string
			value string
			err   error
		)
		switch key {
		case gsKeyGtkTheme:
			ty = TypeGtkTheme
			value = m.setting.GetString(key)
			err = m.doSetGtkTheme(value)
		case gsKeyIconTheme:
			ty = TypeIconTheme
			value = m.setting.GetString(key)
			err = m.doSetIconTheme(value)
		case gsKeyCursorTheme:
			ty = TypeCursorTheme
			value = m.setting.GetString(key)
			err = m.doSetCursorTheme(value)
		case gsKeyFontStandard:
			ty = TypeStandardFont
			value = m.setting.GetString(key)
			err = m.doSetStandardFont(value)
		case gsKeyFontMonospace:
			ty = TypeMonospaceFont
			value = m.setting.GetString(key)
			err = m.doSetMonnospaceFont(value)
		case gsKeyFontSize:
			ty = TypeFontSize
			size := m.setting.GetDouble(key)
			value = fmt.Sprint(size)
			err = m.doSetFontSize(size)
		case gsKeyBackgroundURIs:
			bgs := m.setting.GetStrv(key)
			m.currentDesktopBgs = bgs
			m.setDesktopBackgrounds(bgs)
		default:
			return
		}
		if err != nil {
			logger.Warningf("Set %v failed: %v", key, err)
			return
		}
		m.emitSignalChanged(ty, value)
	})

	m.listenBgGSettings()
}

func (m *Manager) emitSignalChanged(type0, value string) {
	m.service.Emit(m, "Changed", type0, value)
}

func (m *Manager) listenBgGSettings() {
	gsettings.ConnectChanged(wrapBgSchema, "picture-uri", func(key string) {
		if m.wrapBgSetting == nil {
			return
		}

		logger.Debug(wrapBgSchema, "changed")
		value := m.wrapBgSetting.GetString(key)
		err := m.doSetBackground(value)
		if err != nil {
			logger.Warning(err)
			return
		}
	})

	if m.gnomeBgSetting == nil {
		return
	}
	gsettings.ConnectChanged(gnomeBgSchema, "picture-uri", func(key string) {
		if m.gnomeBgSetting == nil {
			return
		}

		logger.Debug(gnomeBgSchema, "changed")
		value := m.gnomeBgSetting.GetString(gsKeyBackground)
		err := m.doSetBackground(value)
		if err != nil {
			logger.Warning(err)
			return
		}
	})
}
