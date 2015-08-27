package appearance

import (
	"pkg.deepin.io/lib/gio-2.0"
)

func (m *Manager) listenGSettingChanged() {
	m.setting.Connect("changed::theme", func(s *gio.Settings, key string) {
		m.doSetDTheme(m.setting.GetString(key))
	})
	m.setting.GetString(gsKeyTheme)

	m.setting.Connect("changed::font-size", func(s *gio.Settings, key string) {
		m.doSetFontSize(m.setting.GetInt(key))
	})
	m.setting.GetInt(gsKeyFontSize)
}
