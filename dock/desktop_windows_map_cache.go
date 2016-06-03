package dock

const desktopWindowsMapCacheVersion = "0.1"

type desktopWindowsMapCacheManager struct {
	*keyValuesMapCacheManager
}

func newDesktopWindowsMapCacheManager(file string) (*desktopWindowsMapCacheManager, error) {
	c := &desktopWindowsMapCacheManager{}
	var err error
	c.keyValuesMapCacheManager, err = newKeyValuesMapCacheManager(desktopWindowsMapCacheVersion, file)
	if err != nil {
		return nil, err
	}
	return c, nil
}
