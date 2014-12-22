package main

import (
	apidevice "dbus/com/deepin/api/device"
)

func requestUnblockAllDevice() {
	d, err := apidevice.NewDevice("com.deepin.api.Device", "/com/deepin/api/Device")
	if err != nil {
		logger.Error(err)
		return
	}
	err = d.UnblockDevice("all")
	if err != nil {
		logger.Error(err)
	}
}
