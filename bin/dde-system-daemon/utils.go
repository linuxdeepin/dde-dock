/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
