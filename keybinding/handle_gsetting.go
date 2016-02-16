/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

import (
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"gir/gio-2.0"
)

func (m *Manager) listenGSettingChanged() {
	m.sysSetting.Connect("changed", func(s *gio.Settings, key string) {
		m.updateShortcutById(key, shortcuts.KeyTypeSystem)
	})
	m.sysSetting.GetStrv("launcher")

	m.mediaSetting.Connect("changed", func(s *gio.Settings, key string) {
		m.updateShortcutById(key, shortcuts.KeyTypeMedia)
	})
	m.mediaSetting.GetStrv("audio-forward")
}
