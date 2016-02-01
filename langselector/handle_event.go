/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package langselector

import (
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
)

func (lang *LangSelector) onLocaleSuccess() {
	lang.lhelper.ConnectSuccess(func(ok bool, reason string) {
		err := lang.handleLocaleChanged(ok, reason)
		if err != nil {
			lang.logger.Warning(err)
			lang.setPropCurrentLocale(getLocale())
			e := sendNotify(localeIconFailed, "",
				Tr("System language failed to change, please try later."))
			lang.LocaleState = LocaleStateChanged
			if e != nil {
				lang.logger.Warning("sendNotify failed:", e)
			}
			return
		}
		e := sendNotify(localeIconFinished, "",
			Tr("System language has been changed, please log in again after logged out."))
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

	err = installI18nDependent(lang.CurrentLocale)
	if err != nil {
		lang.logger.Warning(err)
	}
	dbus.Emit(lang, "Changed", lang.CurrentLocale)

	return nil
}
