package dock

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"pkg.deepin.io/lib/dbusutil"
)

type pluginSettings map[string]map[string]interface{}

type pluginSettingsStorage struct {
	m      *Manager
	data   pluginSettings
	dataMu sync.Mutex

	timer       *time.Timer
	saving      bool
	saveStateMu sync.Mutex
}

func newPluginSettingsStorage(m *Manager) *pluginSettingsStorage {
	s := &pluginSettingsStorage{m: m}

	jsonStr := m.settings.GetString(settingKeyPluginSettings)
	var v pluginSettings
	err := json.Unmarshal([]byte(jsonStr), &v)
	if err == nil {
		s.data = v
	} else {
		logger.Warning("failed to load plugin settings:", err)
		s.data = make(pluginSettings)
	}

	s.timer = time.AfterFunc(3*time.Second, func() {
		s.save()
		s.saveStateMu.Lock()
		s.saving = false
		s.saveStateMu.Unlock()
	})
	return s
}

func (s *pluginSettingsStorage) requestSave() {
	s.saveStateMu.Lock()
	defer s.saveStateMu.Unlock()

	if s.saving {
		return
	} else {
		s.timer.Reset(1 * time.Second)
		s.saving = true
	}
}

func (s *pluginSettingsStorage) save() {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	jsonData, err := json.Marshal(s.data)
	if err != nil {
		logger.Warning(err)
		return
	}
	ok := s.m.settings.SetString(settingKeyPluginSettings, string(jsonData))
	if !ok {
		logger.Warning("failed to save plugin settings")
	}
}

func (s *pluginSettingsStorage) getJsonStr() (string, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	jsonData, err := json.Marshal(s.data)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	return string(jsonData), nil
}

func (s *pluginSettingsStorage) set(v pluginSettings) {
	s.dataMu.Lock()
	s.data = v
	s.dataMu.Unlock()
	s.requestSave()
}

func (s *pluginSettingsStorage) merge(v pluginSettings) {
	s.dataMu.Lock()

	for key1, value1 := range v {
		if s.data[key1] == nil && len(value1) > 0 {
			s.data[key1] = make(map[string]interface{})
		}

		for key2, value2 := range value1 {
			s.data[key1][key2] = value2
		}
	}

	s.dataMu.Unlock()
	s.requestSave()
}

func (s *pluginSettingsStorage) remove(key1 string, key2List []string) {
	s.dataMu.Lock()

	if len(key2List) == 0 {
		delete(s.data, key1)
	} else {
		if value1, ok := s.data[key1]; ok {
			for _, key2 := range key2List {
				delete(value1, key2)
			}
			if len(value1) == 0 {
				delete(s.data, key1)
			}
		}
	}

	s.dataMu.Unlock()
	s.requestSave()
}

func (s *pluginSettingsStorage) equal(v pluginSettings) bool {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()
	return reflect.DeepEqual(s.data, v)
}
