/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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
        "dbus/com/deepin/api/utils"
        "dlib/dbus"
        "dlib/logger"
        "os"
        "sync"
)

var (
        objManager  *Manager
        objUtil     *utils.Utils
        mutex       = new(sync.Mutex)
        logObject   = logger.NewLogger("daemon/themes")
        themeObjMap = make(map[string]*Theme)

        genId, destroyId = func() (func() int, func()) {
                count := 0
                return func() int {
                                mutex.Lock()
                                tmp := count
                                count++
                                mutex.Unlock()
                                return tmp
                        }, func() {
                                mutex.Lock()
                                count = 0
                                mutex.Unlock()
                        }
        }()
)

func destroyThemeObj(path string) {
        obj, ok := themeObjMap[path]
        if !ok {
                return
        }

        dbus.UnInstallObject(obj)
        delete(themeObjMap, path)
}

func destroyAllThemeObj() {
        for k, obj := range themeObjMap {
                dbus.UnInstallObject(obj)
                delete(themeObjMap, k)
        }
}

func updateThemeObj(pathNameMap map[string]PathInfo) {
        destroyAllThemeObj()
        destroyId()

        for path, info := range pathNameMap {
                obj := newTheme(path, info)
                err := dbus.InstallOnSession(obj)
                if err != nil {
                        logObject.Warning("Install Session Failed: %v", err)
                        panic(err)
                }
                themeObjMap[path] = obj
        }
}

func main() {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Fatal("Recove error in main: %v", err)
                }
        }()

        logObject.SetRestartCommand("/usr/lib/deepin-daemon/themes")
        var err error
        objUtil, err = utils.NewUtils("com.deepin.api.Utils",
                "/com/deepin/api/Utils")
        if err != nil {
                logObject.Info("New Utils Object Failed: %v\n", err)
                return
        }

        objManager = newManager()
        err = dbus.InstallOnSession(objManager)
        if err != nil {
                logObject.Warning("Install Session Failed: %v", err)
                panic(err)
        }

        //m.ThemeList = append(m.ThemeList, THEME_PATH+"Test")
        //m.ThemeList = append(m.ThemeList, THEME_PATH+"Deepin")
        updateThemeObj(objManager.pathNameMap)
        objManager.getThemeObject(objManager.CurrentTheme).setThemeViaXSettings()
        objThumb := &ThumbPath{}
        dbus.InstallOnSession(objThumb)
        objPre := &PreviewPath{}
        dbus.InstallOnSession(objPre)
        dbus.DealWithUnhandledMessage()

        if err = dbus.Wait(); err != nil {
                logObject.Warning("lost dbus session: %v", err)
                os.Exit(1)
        } else {
                os.Exit(0)
        }
}
