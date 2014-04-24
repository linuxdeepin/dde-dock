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

package inputdevices

import (
        libsession "dbus/com/deepin/sessionmanager"
        Logger "dlib/logger"
        libutil "dlib/utils"
)

var (
        logObj  = Logger.NewLogger("input-device")
        utilObj = libutil.NewUtils()
        xsObj   *libsession.XSettings
)

func StartInputDevices() {
        var err error
        xsObj, err = libsession.NewXSettings("com.deepin.SessionManager",
                "/com/deepin/XSettings")
        if err != nil {
                logObj.Info("New XSettings Object Failed: ", err)
                return
        }

        initGdkEnv()
        listenDevsSettings()
        //go dlib.StartLoop()
        //select {}
}
