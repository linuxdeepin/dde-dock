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
	"pkg.linuxdeepin.com/lib/dbus"
)

func (ls *LangSelect) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DEST,
		PATH,
		IFC,
	}
}

func (ls *LangSelect) getCurrentLocale() (locale string) {
	tmp, ok := ls.getUserLocale()
	if !ok {
		tmp, ok = ls.getDefaultLocale()
		if !ok {
			locale = DEFAULT_LOCALE
			return
		}
		locale = tmp
	} else {
		locale = tmp
	}

	return
}

func (ls *LangSelect) setPropCurrentLocale(locale string) {
	if ls.CurrentLocale != locale {
		ls.CurrentLocale = locale
		dbus.NotifyChange(ls, "CurrentLocale")
	}
}
