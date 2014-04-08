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
        libkeybind "dbus/com/deepin/daemon/keybinding"
        libdbus "dbus/org/freedesktop/dbus"
        liblogger "dlib/logger"
)

var (
        logObj      = liblogger.NewLogger("daemon/mpris-daemon")
        dbusObj     *libdbus.DBusDaemon
        mediaKeyObj *libkeybind.MediaKey
)

func init() {
        var err error

        dbusObj, err = libdbus.NewDBusDaemon("org.freedesktop.DBus",
                "/")
        if err != nil {
                logObj.Info("New DBusDaemon Failed: ", err)
                panic(err)
        }

        mediaKeyObj, err = libkeybind.NewMediaKey(
                "com.deepin.daemon.KeyBinding",
                "/com/deepin/daemon/MediaKey")
        if err != nil {
                logObj.Info("New MediaKey Object Failed: ", err)
                panic(err)
        }

        logObj.SetRestartCommand("/usr/lib/deepin-daemon/mpris-daemon")
}

func main() {
        listenAudioSignal()

        select {}
}
