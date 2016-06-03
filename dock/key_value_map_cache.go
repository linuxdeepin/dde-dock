package dock

type keyValueMapCache struct {
	Version string
	Data    map[string]string
}

func (c *keyValueMapCache) SetVersion(v string) {
	c.Version = v
}

func (c *keyValueMapCache) GetVersion() string {
	return c.Version
}

func (c *keyValueMapCache) ClearContent() {
	c.Data = make(map[string]string)
}

type keyValueMapCacheManager struct {
	*cacheManager
}

func newKeyValueMapCacheManager(version, file string) (*keyValueMapCacheManager, error) {
	m := &keyValueMapCacheManager{}
	cache := &keyValueMapCache{}
	m.cacheManager = newCacheManager(cache, version)
	err := m.OpenFile(file)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *keyValueMapCacheManager) GetData() map[string]string {
	return m.cache.(*keyValueMapCache).Data
}

func (m *keyValueMapCacheManager) GetValueByKey(key string) string {
	data := m.GetData()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return data[key]
}

func (m *keyValueMapCacheManager) SetKeyValue(key, value string) {
	data := m.GetData()
	m.mutex.Lock()
	if data[key] != value {
		logger.Debugf("SetKeyValue %v => %v", key, value)
		data[key] = value
		m.hasChanged = true
	}
	m.mutex.Unlock()
}

func (m *keyValueMapCacheManager) DeleteKey(key string) {
	data := m.GetData()
	logger.Debug("DeleteKey", key)
	m.mutex.Lock()
	delete(data, key)
	m.mutex.Unlock()
}
