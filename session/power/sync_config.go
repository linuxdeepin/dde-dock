package power

import "encoding/json"

type syncConfig struct {
	m *Manager
}

func (sc *syncConfig) Get() (interface{}, error) {
	return &syncData{
		Version:         "1.0",
		ScreenBlackLock: sc.m.ScreenBlackLock.Get(),
		SleepLock:       sc.m.SleepLock.Get(),
	}, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var v syncData
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	sc.m.ScreenBlackLock.Set(v.ScreenBlackLock)
	sc.m.SleepLock.Set(v.SleepLock)
	return nil
}

// version: 1.0
type syncData struct {
	Version         string `json:"version"`
	ScreenBlackLock bool   `json:"screen_black_lock"`
	SleepLock       bool   `json:"sleep_lock"`
}
