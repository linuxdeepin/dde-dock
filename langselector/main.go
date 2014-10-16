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

package langselector

import (
	"pkg.linuxdeepin.com/lib"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
)

const (
	dbusSender = "com.deepin.daemon.LangSelector"
)

var _lang *LangSelector

func Start() *LangSelector {
	var logger = log.NewLogger(dbusSender)

	if !lib.UniqueOnSession(dbusSender) {
		logger.Warning("There is a LangSelector running...")
		return nil
	}

	logger.BeginTracing()

	_lang = newLangSelect(logger)
	if _lang == nil {
		logger.Fatal("Create LangSelector Failed")
	}

	err := dbus.InstallOnSession(_lang)
	if err != nil {
		logger.Fatal("Install Session DBus Failed:", err)
	}

	_lang.onGenLocaleStatus()

	return _lang
}

func Stop() {
	if _lang == nil {
		return
	}

	_lang.Destroy()
	_lang.logger.EndTracing()
	_lang = nil
}
