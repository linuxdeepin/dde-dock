package dock

import (
	"encoding/gob"
	"os"
	"sync"
)

type cache interface {
	SetVersion(string)
	GetVersion() string
	ClearContent()
}

type cacheManager struct {
	cache           cache
	standardVersion string
	autoSaveEnabled bool
	hasChanged      bool
	file            *os.File
	mutex           sync.Mutex
}

func newCacheManager(c cache, version string) *cacheManager {
	return &cacheManager{
		cache:           c,
		standardVersion: version,
		autoSaveEnabled: true,
	}
}

func (m *cacheManager) Clear() {
	m.mutex.Lock()
	m.cache.SetVersion(m.standardVersion)
	m.cache.ClearContent()
	m.hasChanged = true
	m.mutex.Unlock()
}

func (m *cacheManager) OpenFile(file string) error {
	var err error
	m.file, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(m.file)
	m.mutex.Lock()
	err = dec.Decode(m.cache)
	m.mutex.Unlock()
	if err != nil {
		logger.Warning("Decode error:", err)
		m.Clear()
	}

	// check version match
	if m.cache.GetVersion() != m.standardVersion {
		logger.Warning("version not match")
		m.Clear()
	} else {
		logger.Info("load file ok")
		logger.Debugf("load content: %#v", m.cache)
	}
	return nil
}

func (m *cacheManager) Save() error {
	logger.Info("call save")
	enc := gob.NewEncoder(m.file)
	// seek head
	_, err := m.file.Seek(0, 0)
	if err != nil {
		logger.Warning("save error", err)
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	err = enc.Encode(m.cache)
	logger.Debugf("save content: %#v", m.cache)
	if err != nil {
		logger.Warning("save error:", err)
		return err
	}
	err = m.file.Sync()
	if err != nil {
		logger.Warning("save error:", err)
		return err
	}
	m.hasChanged = false
	return nil
}

func (m *cacheManager) SetAutoSaveEnabled(val bool) {
	m.autoSaveEnabled = val
}

func (m *cacheManager) AutoSave() error {
	if !m.autoSaveEnabled {
		logger.Debug("Skip save")
		return nil
	}
	if m.hasChanged {
		return m.Save()
	}
	return nil
}

func (m *cacheManager) Destroy() {
	m.file.Close()
}
