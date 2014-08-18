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
	"dbus/com/deepin/api/greeterutils"
	libsession "dbus/com/deepin/sessionmanager"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/log"
)

var (
	logger     = log.NewLogger(DEVICE_DEST)
	xsObj      *libsession.XSettings
	greeterObj *greeterutils.GreeterUtils
	managerObj *Manager

	tpadSettings  = gio.NewSettings("com.deepin.dde.touchpad")
	mouseSettings = gio.NewSettings("com.deepin.dde.mouse")
	kbdSettings   = gio.NewSettings("com.deepin.dde.keyboard")
	layoutDescMap = make(map[string]string)
)

func Stop() {
	logger.EndTracing()
}
func Start() {
	logger.BeginTracing()
	logger.SetRestartCommand("/usr/lib/deepin-daemon/dde-session-daemon")

	var err error
	xsObj, err = libsession.NewXSettings("com.deepin.SessionManager",
		"/com/deepin/XSettings")
	if err != nil {
		logger.Warning("New XSettings Object Failed: ", err)
		return
	}

	if greeterObj, err = greeterutils.NewGreeterUtils("com.deepin.api.GreeterUtils", "/com/deepin/api/GreeterUtils"); err != nil {
		logger.Warning("New GreeterUtils failed:", err)
		return
	}

	listenDevsSettings()

	managerObj = NewManager()
	if err = dbus.InstallOnSession(managerObj); err != nil {
		logger.Fatal("Manager DBus Session Failed: ", err)
	}

	datas := parseXML(_LAYOUT_XML_PATH)
	layoutDescMap = getLayoutList(datas)

	mouse := NewMouse()
	if err := dbus.InstallOnSession(mouse); err != nil {
		logger.Fatal("Mouse DBus Session Failed: ", err)
	}
	managerObj.mouseObj = mouse

	kbd := NewKeyboard()
	if err := dbus.InstallOnSession(kbd); err != nil {
		logger.Fatal("Kbd DBus Session Failed: ", err)
	}
	managerObj.kbdObj = kbd
	//setLayoutOptions()
	setLayout(kbd.CurrentLayout.GetValue().(string))

	tpadFlag := false
	for _, info := range managerObj.Infos {
		if info.Id == "touchpad" {
			tpad := NewTPad()
			if err := dbus.InstallOnSession(tpad); err != nil {
				logger.Fatal("TPad DBus Session Failed: ", err)
			}
			tpadFlag = true
			managerObj.tpadObj = tpad
			break
		}
	}
	initGSettingsSet(tpadFlag)

	if managerObj.mouseObj.Exist {
		disableTPadWhenMouse()
	} else {
		tpadSettings.SetBoolean(TPAD_KEY_ENABLE, true)
	}
}
