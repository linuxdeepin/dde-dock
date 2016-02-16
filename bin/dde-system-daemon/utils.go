/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
