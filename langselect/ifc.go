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
	"errors"
	"pkg.linuxdeepin.com/lib/dbus"
	. "pkg.linuxdeepin.com/lib/gettext"
)

func (ls *LangSelect) SetLocale(locale string) error {
	if len(locale) < 1 {
		return errors.New("SetLocale args error")
	}

	if ls.CurrentLocale == locale {
		return nil
	}

	ls.sendNotify("", "", Tr("Language is changing, please wait"))
	setDate.GenLocale(locale)
	ls.changeLocaleFlag = true
	ls.CurrentLocale = locale
	dbus.NotifyChange(ls, "CurrentLocale")

	return nil
}

func (ls *LangSelect) GetLocaleList() (list []localeInfo) {
	return ls.getLocaleInfoList()
}
