/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package langselector

import (
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

const (
	dbusSender = "com.deepin.daemon.LangSelector"
)

var (
	_lang  *LangSelector
	logger = log.NewLogger("daemon/langselector")
)

func Start() *LangSelector {
	if !lib.UniqueOnSession(dbusSender) {
		logger.Warning("There is a LangSelector running...")
		return nil
	}

	logger.BeginTracing()

	_lang = newLangSelector(logger)
	if _lang == nil {
		logger.Error("Create LangSelector Failed")
		return nil
	}

	err := dbus.InstallOnSession(_lang)
	if err != nil {
		logger.Error("Install Session DBus Failed:", err)
		Stop()
		return nil
	}

	_lang.onLocaleSuccess()

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
