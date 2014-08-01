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

package langselect

import (
	"dbus/com/deepin/api/setdatetime"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
)

var (
	Logger  = log.NewLogger(DEST)
	setDate *setdatetime.SetDateTime
)

var _manager *LangSelect

func GetManager() *LangSelect {
	if _manager == nil {
		_manager = newLangSelect()
	}

	return _manager
}

func newLangSelect() *LangSelect {
	ls := &LangSelect{}

	ls.changeLocaleFlag = false
	ls.setPropCurrentLocale(ls.getCurrentLocale())

	return ls
}

func Start() {
	Logger.BeginTracing()

	var err error
	if setDate, err = setdatetime.NewSetDateTime("com.deepin.api.SetDateTime", "/com/deepin/api/SetDateTime"); err != nil {
		Logger.Fatal("New SetDateTime failed:", err)
		return
	}

	ls := GetManager()
	if err := dbus.InstallOnSession(ls); err != nil {
		Logger.Fatal("Install DBus Failed:", err)
		return
	}

	ls.listenLocaleChange()
}

func Stop() {
	Logger.EndTracing()
}
