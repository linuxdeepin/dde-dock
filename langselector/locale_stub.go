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
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusPath      = "/com/deepin/daemon/LangSelector"
	dbusInterface = "com.deepin.daemon.LangSelector"
)

func (lang *LangSelector) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: dbusPath,
		Interface:  dbusInterface,
	}
}

func (lang *LangSelector) setPropCurrentLocale(locale string) {
	if lang.CurrentLocale != locale {
		lang.CurrentLocale = locale
		dbus.NotifyChange(lang, "CurrentLocale")
	}
}
