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
        "dlib/dbus"
        "dlib/gio-2.0"
        "dlib/logger"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{_DATE_TIME_DEST, _DATE_TIME_PATH, _DATA_TIME_IFC}
}

func (op *Manager) setAutoSetTime(auto bool) (bool, error) {
        var (
                ret     bool
                err     error
        )

        if auto {
                ret, err = setDate.SetNtpUsing(true)
        } else {
                ret, err = setDate.SetNtpUsing(false)
        }

        if err != nil {
                logger.Printf("Set NTP - %d Failed: %s\n",
                        auto, err)
                return false, err
        }
        return ret, nil
}

func (op *Manager) setPropName(name string) {
        switch name {
        case "CurrentTimezone":
                tz, _, err := setDate.GetTimezone()
                if err != nil {
                        logger.Printf("Get Time Zone Failed: %s\n", err)
                        return
                }
                op.CurrentTimezone = tz
                dbus.NotifyChange(op, name)
        case "UserTimezoneList":
                list := dateSettings.GetStrv("user-timezone-list")
                if !strArrayIsEqual(list, op.UserTimezoneList) {
                        op.UserTimezoneList = list
                        dbus.NotifyChange(op, "UserTimezoneList")
                }
        }
}

func (op *Manager) listenSettings() {
        dateSettings.Connect("changed::is-auto-set", func(s *gio.Settings, name string) {
                op.setAutoSetTime(s.GetBoolean("is-auto-set"))
        })
        dateSettings.Connect("changed::user-timezone-list", func(s *gio.Settings, name string) {
                op.setPropName("UserTimezoneList")
        })
}

func (op *Manager) listenZone() {
        err := zoneWatcher.Watch(_TIME_ZONE_FILE)
        if err != nil {
                logger.Printf("Watch '%s' Failed: %s\n", _TIME_ZONE_FILE, err)
                return
        }

        go func() {
                for {
                        select {
                        case ev := <-zoneWatcher.Event:
                                logger.Println("Watcher Event: ", ev)
                                if ev.IsDelete() {
                                        zoneWatcher.Watch(_TIME_ZONE_FILE)
                                } else {
                                        //if ev.IsModify() {
                                        op.setPropName("CurrentTimezone")
                                        //}
                                }
                        case err := <-zoneWatcher.Error:
                                logger.Println("Watcher Event: ", err)
                        }
                }
        }()
}
