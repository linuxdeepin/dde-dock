/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package accounts

import (
	"dlib/dbus"
	Logger "dlib/logger"
	Utils "dlib/utils"
)

var (
	logger  = Logger.NewLogger(ACCOUNT_DEST)
	objUtil = Utils.NewUtils()
)

func Start() {
	logger.BeginTracing()
	defer logger.EndTracing()

	obj := GetManager()
	if err := dbus.InstallOnSystem(obj); err != nil {
		logger.Error("Install DBus Failed:", err)
		panic(err)
	}

	obj.updateAllUserInfo()

	dbus.DealWithUnhandledMessage()
}

func Stop() {
	obj := GetManager()

	obj.endFlag <- true
	obj.listEndFlag <- true
	obj.infoEndFlag <- true
	obj.infoWatcher.Close()
	obj.listWatcher.Close()
	obj.destroyAllUser()
	dbus.UnInstallObject(obj)
}
