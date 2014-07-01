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
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	Logger "pkg.linuxdeepin.com/lib/logger"
)

var (
	logObj     = Logger.NewLogger("input device")
	xsObj      *libsession.XSettings
	managerObj *Manager

	tpadSettings  = gio.NewSettings("com.deepin.dde.touchpad")
	mouseSettings = gio.NewSettings("com.deepin.dde.mouse")
	kbdSettings   = gio.NewSettings("com.deepin.dde.keyboard")
	layoutDescMap = make(map[string]string)
)

func Stop() {
	logObj.EndTracing()
}
func Start() {
	logObj.BeginTracing()
	logObj.SetRestartCommand("/usr/lib/deepin-daemon/inputdevices")

	var err error
	xsObj, err = libsession.NewXSettings("com.deepin.SessionManager",
		"/com/deepin/XSettings")
	if err != nil {
		logObj.Info("New XSettings Object Failed: ", err)
		return
	}

	listenDevsSettings()

	managerObj = NewManager()
	if err = dbus.InstallOnSession(managerObj); err != nil {
		logObj.Warning("Manager DBus Session Failed: ", err)
		panic(err)
	}

	datas := parseXML(_LAYOUT_XML_PATH)
	layoutDescMap = getLayoutList(datas)

	tpadFlag := false
	mouseFlag := false
	for _, info := range managerObj.Infos {
		if info.Id == "mouse" {
			mouse := NewMouse()
			if err := dbus.InstallOnSession(mouse); err != nil {
				logObj.Warning("Mouse DBus Session Failed: ", err)
				panic(err)
			}
			managerObj.mouseObj = mouse
			mouseFlag = true
		} else if info.Id == "touchpad" {
			tpad := NewTPad()
			if err := dbus.InstallOnSession(tpad); err != nil {
				logObj.Warning("TPad DBus Session Failed: ", err)
				panic(err)
			}
			tpadFlag = true
			managerObj.tpadObj = tpad
		} else if info.Id == "keyboard" {
			kbd := NewKeyboard()
			if err := dbus.InstallOnSession(kbd); err != nil {
				logObj.Warning("Kbd DBus Session Failed: ", err)
				panic(err)
			}
			managerObj.kbdObj = kbd
			setLayout(kbd.CurrentLayout.GetValue().(string))
		}
	}
	initGSettingsSet(tpadFlag)

	if mouseFlag {
		disableTPadWhenMouse()
	} else {
		info := deviceInfo{DEVICE_PATH + "Mouse", "mouse"}
		managerObj.Infos = append(managerObj.Infos, info)
		mouse := NewMouse()
		if err := dbus.InstallOnSession(mouse); err != nil {
			logObj.Warning("Mouse DBus Session Failed: ", err)
			panic(err)
		}
		managerObj.mouseObj = mouse
	}
}
