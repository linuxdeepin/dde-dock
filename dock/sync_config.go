package dock

import "encoding/json"

type syncConfig struct {
	m *Manager
}

func (sc *syncConfig) Get() (interface{}, error) {
	var v syncData
	v.Version = "1.0"
	v.IconSize = sc.m.IconSize.Get()
	v.DisplayMode = sc.m.DisplayMode.GetString()
	v.HideMode = sc.m.HideMode.GetString()
	v.Position = sc.m.Position.GetString()
	// TODO: docked apps
	//v.DockedApps = sc.m.DockedApps.Get()
	return v, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var v syncData
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	m := sc.m
	m.IconSize.Set(v.IconSize)
	m.DisplayMode.SetString(v.DisplayMode)
	m.HideMode.SetString(v.HideMode)
	m.Position.SetString(v.Position)
	// TODO: docked apps
	return nil
}

// version: 1.0
type syncData struct {
	Version     string            `json:"version"` // such as "1.0.0"
	IconSize    uint32            `json:"icon_size"`
	DisplayMode string            `json:"display_mode"`
	HideMode    string            `json:"hide_mode"`
	Position    string            `json:"position"`
	DockedApps  DockedAppInfoList `json:"docked_apps"`
}

type DockedAppInfo struct {
	Key  string `json:"key"`
	File string `json:"file"`
}

type DockedAppInfoList []*DockedAppInfo
