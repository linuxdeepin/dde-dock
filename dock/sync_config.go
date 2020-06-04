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
	v.WindowSizeEfficient = sc.m.WindowSizeEfficient.Get()
	v.WindowSizeFashion = sc.m.WindowSizeFashion.Get()
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
	added := dockedApps
	removed := m.DockedApps.Get()
	for _, value := range removed {
		desktopFile := unzipDesktopPath(value)
		_, err := m.requestUndock(desktopFile)
		if err != nil {
			logger.Warning(err)
		}
	}

	var index = 0
	for _, value := range added {
		desktopFile := unzipDesktopPath(value)
		_, err := os.Stat(desktopFile)
		if err == nil {
			_, err = m.requestDock(desktopFile, int32(index))
			if err != nil {
				logger.Warning(err)
			} else {
				index++
			}
		}
	}

	// emit signal
	err := m.service.Emit(m, "DockAppSettingsSynced")
	if err != nil {
		logger.Warning(err)
	}
}

func (sc *syncConfig) Set(data []byte) error {
	var v syncData
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	m := sc.m
	if v.WindowSizeEfficient > 0 {
		m.WindowSizeEfficient.Set(v.WindowSizeEfficient)
	}
	if v.WindowSizeFashion > 0 {
		m.WindowSizeFashion.Set(v.WindowSizeFashion)
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
	syncConfigVersion = "1.2"
)

type syncData struct {
	Version             string         `json:"version"`
	WindowSizeEfficient uint32         `json:"window_size_efficient"`
	WindowSizeFashion   uint32         `json:"window_size_fashion"`
	DisplayMode         string         `json:"display_mode"`
	HideMode            string         `json:"hide_mode"`
	Position            string         `json:"position"`
	DockedApps          []string       `json:"docked_apps"`
	Plugins             pluginSettings `json:"plugins"`
}
