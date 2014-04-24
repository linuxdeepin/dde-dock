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

package main

import (
        "dlib"
        Logger "dlib/logger"
        libutils "dlib/utils"
)

//import "dde-daemon/deepin-daemon/inputdevices"

var (
        logObj   = Logger.NewLogger("deepin/daemon")
        utilsObj = libutils.NewUtils()
)

func main() {
        defer logObj.EndTracing()

        logObj.SetRestartCommand("/usr/lib/deepin-daemon/deepin-daemon")
        //enableTouchPad()
        //listenDevices()
        //inputdevices.StartInputDevices()

        startMprisDaemon()
        go dlib.StartLoop()

        select {}
}
