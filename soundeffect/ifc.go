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
	if !m.canPlayEvent(event) {
		return nil
	}

	err := m.player.PlaySystemSound(queryEvent(event))
	if err != nil {
		logger.Debugf("Play sound event '%s' failed: %v",
			queryEvent(event), err)
		return err
	}
	return nil
}

func (m *Manager) PlayThemeSound(theme, event string) error {
	if !m.canPlayEvent(event) {
		return nil
	}

	err := m.player.PlayThemeSound(theme, queryEvent(event))
	if err != nil {
		logger.Debugf("Play theme '%s' sound event '%s' failed: %v",
			theme, queryEvent(event), err)
		return err
	}
	return nil
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
