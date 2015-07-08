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
	"pkg.deepin.io/dde/daemon/langselector/i18n_dependency"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
)

func (lang *LangSelector) onLocaleSuccess() {
	lang.lhelper.ConnectSuccess(func(ok bool, reason string) {
		err := lang.handleLocaleChanged(ok, reason)
		if err != nil {
			lang.logger.Warning(err)
			lang.setPropCurrentLocale(getLocale())
			e := sendNotify("", "", Tr("System language failed to change, please try later."))
			lang.LocaleState = LocaleStateChanged
			if e != nil {
				lang.logger.Warning("sendNotify failed:", e)
			}
			return
		}
		e := sendNotify("", "", Tr("System language has been changed, please log in again after logged out."))
		lang.LocaleState = LocaleStateChanged
		if e != nil {
			lang.logger.Warning("sendNotify failed:", e)
		}
	})
}

func (lang *LangSelector) handleLocaleChanged(ok bool, reason string) error {
	if !ok || lang.LocaleState != LocaleStateChanging {
		return ErrLocaleChangeFailed
	}

	err := writeUserLocale(lang.CurrentLocale)
	if err != nil {
		return err
	}

	err = i18n_dependency.InstallDependentPackages(lang.CurrentLocale)
	if err != nil {
		lang.logger.Warning(err)
	}
	dbus.Emit(lang, "Changed", lang.CurrentLocale)

	return nil
}
