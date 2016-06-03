package dock

const desktopHashFileMapCacheVersion = "0.1"

type desktopHashFileMapCacheManager struct {
	*keyValueMapCacheManager
}

func newDesktopHashFileMapCacheManager(file string) (*desktopHashFileMapCacheManager, error) {
	c := &desktopHashFileMapCacheManager{}
	var err error
	c.keyValueMapCacheManager, err = newKeyValueMapCacheManager(desktopHashFileMapCacheVersion, file)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *desktopHashFileMapCacheManager) GetAppInfo(desktopHash string) *AppInfo {
	desktopFile := c.GetValueByKey(desktopHash)
	if desktopFile == "" {
		return nil
	}
	appInfo := NewAppInfoFromFile(desktopFile)
	return appInfo
}
