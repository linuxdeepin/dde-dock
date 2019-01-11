package audio

import (
	"encoding/json"

	"pkg.deepin.io/gir/gio-2.0"
)

type syncSoundEffect struct {
	Enabled bool `json:"enabled"`
}

type syncData struct {
	Version     string           `json:"version"`
	SoundEffect *syncSoundEffect `json:"soundeffect"`
}

type syncConfig struct {
	a *Audio
}

const (
	syncVersion = "1.0"
)

func (sc *syncConfig) Get() (interface{}, error) {
	s := gio.NewSettings(gsSchemaSoundEffect)
	defer s.Unref()
	return &syncData{
		Version: syncVersion,
		SoundEffect: &syncSoundEffect{
			Enabled: s.GetBoolean(gsKeyEnabled),
		},
	}, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var info syncData
	err := json.Unmarshal(data, &info)
	if err != nil {
		return err
	}
	if info.SoundEffect != nil {
		s := gio.NewSettings(gsSchemaSoundEffect)
		s.SetBoolean(gsKeyEnabled, info.SoundEffect.Enabled)
		s.Unref()
	}
	return nil
}
