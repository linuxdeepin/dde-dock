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

	m.listenBgGsettings()
}

func (m *Manager) listenBgGsettings() {
	m.wrapBgSetting.Connect("changed::picture-uri", func(s *gio.Settings, key string) {
		uri := m.wrapBgSetting.GetString(gsKeyBackground)
		err := m.doSetBackground(uri)
		if err != nil {
			logger.Debugf("[Wrap background] set '%s' failed: %s", uri, err)
		}
	})
	m.wrapBgSetting.GetString(gsKeyBackground)

	if m.gnomeBgSetting != nil {
		m.gnomeBgSetting.Connect("changed::picture-uri", func(s *gio.Settings, key string) {
			uri := m.gnomeBgSetting.GetString(gsKeyBackground)
			err := m.doSetBackground(uri)
			if err != nil {
				logger.Debugf("[Gnome background] set '%s' failed: %s", uri, err)
			}
		})
		m.gnomeBgSetting.GetString(gsKeyBackground)
	}
}
