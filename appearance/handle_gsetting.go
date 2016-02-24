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
	"gir/gio-2.0"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"time"
)

func (m *Manager) listenGSettingChanged() {
	m.setting.Connect("changed::theme", func(s *gio.Settings, key string) {
		value := m.setting.GetString(key)
		if m.Theme == value {
			return
		}
		m.doSetDTheme(value)
	})
	m.setting.GetString(gsKeyTheme)

	m.setting.Connect("changed::font-size", func(s *gio.Settings, key string) {
		m.doSetFontSize(m.setting.GetInt(key))
	})
	m.setting.GetInt(gsKeyFontSize)

	m.listenBgGSettings()
}

func (m *Manager) listenBgGSettings() {
	if m.gnomeBgSetting == nil {
		return
	}
	m.gnomeBgSetting.Connect("changed::picture-uri", func(s *gio.Settings, key string) {
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

		err := m.doSetBackground(uri)
		if err != nil {
			logger.Debugf("[Gnome background] set '%s' failed: %s", uri, err)
			return
		}
		logger.Debug("[Gnome background] sync wrap bg OVER ENDDDDDDDD:", uri)
	})
	m.gnomeBgSetting.GetString(gsKeyBackground)
}
