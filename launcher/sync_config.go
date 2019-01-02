package launcher

import "encoding/json"

type syncConfig struct {
	m *Manager
}

func (*syncConfig) Name() string {
	return "launcher"
}

func (sc *syncConfig) Get() (interface{}, error) {
	var v syncData
	v.Version = "1.0"
	v.DisplayMode = sc.m.DisplayMode.GetString()
	v.Fullscreen = sc.m.Fullscreen.Get()
	return v, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var v syncData
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	sc.m.DisplayMode.SetString(v.DisplayMode)
	sc.m.Fullscreen.Set(v.Fullscreen)
	return nil
}

// version: 1.0
type syncData struct {
	Version     string `json:"version"`
	DisplayMode string `json:"display_mode"`
	Fullscreen  bool   `json:"fullscreen"`
}
