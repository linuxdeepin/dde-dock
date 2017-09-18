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
	"gir/gio-2.0"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/lib/dbus"
	"time"
)

func (m *Manager) listenGSettingChanged() {
	m.setting.Connect("changed", func(s *gio.Settings, key string) {
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
		default:
			return
		}
		if err != nil {
			logger.Warningf("Set %v failed: %v", key, err)
			return
		}
		dbus.Emit(m, "Changed", ty, value)
	})
	m.setting.GetDouble(gsKeyFontSize)

	m.listenBgGSettings()
}

func (m *Manager) listenBgGSettings() {
	m.wrapBgSetting.Connect("changed::picture-uri", func(s *gio.Settings, key string) {
		if m.wrapBgSetting == nil {
			return
		}

		logger.Debug("wrapBgSetting picture-uri changed")
		value := m.wrapBgSetting.GetString(key)
		uri, err := m.doSetBackground(value)
		if err != nil {
			// TODO: set bg to default bg
			logger.Warning("Ensure bg exists failed:", err, value)
			return
		}
		if uri != value {
			m.wrapBgSetting.SetString(key, uri)
		}
		dbus.Emit(m, "Changed", TypeBackground, uri)
	})
	m.wrapBgSetting.GetString(gsKeyBackground)

	if m.gnomeBgSetting == nil {
		return
	}
	m.gnomeBgSetting.Connect("changed::picture-uri", func(s *gio.Settings, key string) {
		if m.gnomeBgSetting == nil {
			return
		}

		// Wait for file copy finished
		time.Sleep(time.Millisecond * 500)
		uri := m.gnomeBgSetting.GetString(gsKeyBackground)
		old := m.wrapBgSetting.GetString(gsKeyBackground)
		logger.Debug("[Gnome background] changed:", key, uri, old)
		if uri == old {
			return
		}
		if !background.IsBackgroundFile(uri) {
			logger.Debugf("[Gnome background] Invalid background file '%v'", uri)
			return
		}

		v, err := m.doSetBackground(uri)
		if err != nil {
			logger.Debugf("[Gnome background] set '%s' failed: %s", uri, err)
			return
		}

		if v != old {
			m.wrapBgSetting.SetString(key, v)
		}
		logger.Debug("[Gnome background] sync wrap bg OVER ENDDDDDDDD:", uri)
	})
	m.gnomeBgSetting.GetString(gsKeyBackground)
}
