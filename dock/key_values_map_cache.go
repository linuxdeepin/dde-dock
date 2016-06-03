package dock

type keyValuesMapCache struct {
	Version string
	Data    map[string][]string
}

func (c *keyValuesMapCache) SetVersion(v string) {
	c.Version = v
}

func (c *keyValuesMapCache) GetVersion() string {
	return c.Version
}

func (c *keyValuesMapCache) ClearContent() {
	c.Data = make(map[string][]string)
}

type keyValuesMapCacheManager struct {
	*cacheManager
}

func newKeyValuesMapCacheManager(version, file string) (*keyValuesMapCacheManager, error) {
	m := &keyValuesMapCacheManager{}
	cache := &keyValuesMapCache{}
	m.cacheManager = newCacheManager(cache, version)
	err := m.OpenFile(file)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *keyValuesMapCacheManager) GetData() map[string][]string {
	return m.cache.(*keyValuesMapCache).Data
}

func (m *keyValuesMapCacheManager) GetKeyByValue(value string) string {
	data := m.GetData()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for key, slice := range data {
		if isStrInSlice(value, slice) {
			return key
		}
	}
	return ""
}

func (c *keyValuesMapCacheManager) AddKeyValue(key, value string) {
	const valuesLimit = 50

	logger.Debugf("AddKeyValue %v => %v", key, value)
	if key == "" || value == "" {
		logger.Warning("AddKeyValue failed: key or value empty")
		return
	}

	key0 := c.GetKeyByValue(value)
	if key0 != "" {
		logger.Warningf("AddKeyValue failed: %v => %v exist", key0, value)
		return
	}

	data := c.GetData()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	slice, ok := data[key]
	c.hasChanged = true
	if ok {
		// 如果尺寸超过限制值，说明缓存不起作用，清除之前的缓存
		if len(slice) > valuesLimit {
			data[key] = []string{value}
			return
		}

		data[key] = append(slice, value)
	} else {
		// new key
		data[key] = []string{value}
	}
}

func (c *keyValuesMapCacheManager) DeleteKeyValue(key, value string) {
	data := c.GetData()
	c.mutex.Lock()
	slice, ok := data[key]
	if ok {
		var index int = -1
		for i, v := range slice {
			if v == value {
				index = i
			}
		}
		if index != -1 {
			data[key] = append(slice[:index], slice[index+1:]...)
			logger.Debugf("DeleteKeyValue %v => %v", key, value)
			c.hasChanged = true
		}
	}
	c.mutex.Unlock()
}
