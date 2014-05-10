package main

func (b *Bluetooth) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	logger.Debug("OnPropertiesChanged()", name)
	switch name {
	// TODO
	}
}

// GetDevices return all devices object that marshaled by json.
func (b *Bluetooth) GetDevices() (devicesJSON string) {
	// TODO
	return
}
