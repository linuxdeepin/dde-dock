package soundeffect

// deepin sound theme 'key - event' map
var soundEventMap = map[string]string{
	keyLogin:         "sys-login",
	keyShutdown:      "sys-shutdown",
	keyLogout:        "sys-logout",
	keyWakeup:        "suspend-resume",
	keyNotification:  "message-out",
	keyUnableOperate: "app-error-critical",
	keyEmptyTrash:    "trash-empty",
	keyVolumeChange:  "audio-volume-change",
	keyBatteryLow:    "power-unplug-battery-low",
	keyPowerPlug:     "power-plug",
	keyPowerUnplug:   "power-unplug",
	keyDevicePlug:    "device-added",
	keyDeviceUnplug:  "device-removed",
	keyIconToDesktop: "send-to",
	keyScreenshot:    "screen-capture",
}

func (m *Manager) PlaySystemSound(event string) error {
	return m.doPlaySystemSound(event, false)
}

func (m *Manager) PlaySystemSoundSync(event string) error {
	return m.doPlaySystemSound(event, true)
}

func (m *Manager) PlayThemeSound(theme, event string) error {
	return m.doPlayThemeSound(theme, event, false)
}

func (m *Manager) PlayThemeSoundSync(theme, event string) error {
	return m.doPlayThemeSound(theme, event, true)
}

func (m *Manager) doPlaySystemSound(event string, sync bool) error {
	if !m.canPlayEvent(event) {
		return nil
	}

	event = queryEvent(event)
	if sync {
		return m.player.PlaySystemSoundSync(event)
	}
	return m.player.PlaySystemSound(event)
}

func (m *Manager) doPlayThemeSound(theme, event string, sync bool) error {
	if !m.canPlayEvent(event) {
		return nil
	}

	event = queryEvent(event)
	if sync {
		return m.player.PlayThemeSoundSync(theme, event)
	}
	return m.player.PlayThemeSound(theme, event)
}

func (m *Manager) canPlayEvent(event string) bool {
	if !isItemInList(event, m.setting.ListKeys()) {
		return true
	}

	return m.setting.GetBoolean(event)
}

func queryEvent(key string) string {
	value, ok := soundEventMap[key]
	if !ok {
		return key
	}
	return value
}

func isItemInList(item string, list []string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}
