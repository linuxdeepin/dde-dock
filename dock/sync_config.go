package dock

import (
	"encoding/json"
	"os"
)

type syncConfig struct {
	m *Manager
}

func (sc *syncConfig) Get() (interface{}, error) {
	var v syncData
	v.Version = syncConfigVersion
	v.WindowSize = sc.m.WindowSize.Get()
	v.DisplayMode = sc.m.DisplayMode.GetString()
	v.HideMode = sc.m.HideMode.GetString()
	v.Position = sc.m.Position.GetString()
	v.DockedApps = sc.m.DockedApps.Get()

	pluginSettingsJsonStr := sc.m.settings.GetString(settingKeyPluginSettings)
	err := json.Unmarshal([]byte(pluginSettingsJsonStr), &v.Plugins)
	if err != nil {
		logger.Warning(err)
	}

	return v, nil
}

func (sc *syncConfig) setPluginSettings(settings pluginSettings) {
	m := sc.m
	if m.pluginSettings.equal(settings) {
		return
	}
	m.pluginSettings.set(settings)
	// emit signal
	err := m.service.Emit(m, "PluginSettingsSynced")
	if err != nil {
		logger.Warning(err)
	}
}

func (sc *syncConfig) setDockedApps(dockedApps []string) {
	m := sc.m
	added, removed := diffStrSlice(m.DockedApps.Get(), dockedApps)

	for _, value := range added {
		desktopFile := unzipDesktopPath(value)
		_, err := os.Stat(desktopFile)
		if err == nil {
			_, err = m.requestDock(desktopFile, -1)
			if err != nil {
				logger.Warning(err)
			}
		}
	}

	for _, value := range removed {
		desktopFile := unzipDesktopPath(value)

		_, err := m.requestUndock(desktopFile)
		if err != nil {
			logger.Warning(err)
		}
	}

	m.DockedApps.Set(dockedApps)
}

func (sc *syncConfig) Set(data []byte) error {
	var v syncData
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	m := sc.m
	if v.WindowSize > 0 {
		m.WindowSize.Set(v.WindowSize)
	}
	m.DisplayMode.SetString(v.DisplayMode)
	m.HideMode.SetString(v.HideMode)
	m.Position.SetString(v.Position)
	sc.setDockedApps(v.DockedApps)
	sc.setPluginSettings(v.Plugins)
	return nil
}

func diffStrSlice(a, b []string) (added, removed []string) {
	// from a to b
	toMap := func(slice []string) map[string]struct{} {
		m := make(map[string]struct{})
		for _, value := range slice {
			m[value] = struct{}{}
		}
		return m
	}
	mapA := toMap(a)
	mapB := toMap(b)

	for keyB := range mapB {
		_, ok := mapA[keyB]
		if !ok {
			added = append(added, keyB)
		}
	}

	for keyA := range mapA {
		_, ok := mapB[keyA]
		if !ok {
			removed = append(removed, keyA)
		}
	}
	return
}

const (
	syncConfigVersion = "1.1"
)

type syncData struct {
	Version     string         `json:"version"`
	WindowSize  uint32         `json:"window_size"`
	DisplayMode string         `json:"display_mode"`
	HideMode    string         `json:"hide_mode"`
	Position    string         `json:"position"`
	DockedApps  []string       `json:"docked_apps"`
	Plugins     pluginSettings `json:"plugins"`
}
