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
