package inputdevices

import (
	"encoding/json"
)

type syncConfig struct {
	m *Manager
}

type syncMouseData struct {
	NaturalScroll bool `json:"natural_scroll"`
}

type syncTPadData struct {
	NaturalScroll bool `json:"natural_scroll"`
}

type syncData struct {
	Version  string         `json:"version"`
	Mouse    *syncMouseData `json:"mouse"`
	Touchpad *syncTPadData  `json:"touchpad"`
}

const (
	syncVersion = "1.0"
)

func (sc *syncConfig) Get() (interface{}, error) {
	return &syncData{
		Version: syncVersion,
		Mouse: &syncMouseData{
			NaturalScroll: sc.m.mouse.NaturalScroll.Get(),
		},
		Touchpad: &syncTPadData{
			NaturalScroll: sc.m.tpad.NaturalScroll.Get(),
		},
	}, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var info syncData
	err := json.Unmarshal(data, &info)
	if err != nil {
		return err
	}
	sc.m.mouse.NaturalScroll.Set(info.Mouse.NaturalScroll)
	sc.m.tpad.NaturalScroll.Set(info.Touchpad.NaturalScroll)
	return nil
}
