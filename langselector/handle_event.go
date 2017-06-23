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
	"os"
	"os/user"
	"path/filepath"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/xdg/basedir"
)

var (
	// for locale-helper
	_ = Tr("Authentication is required to switch language")
)

func (lang *LangSelector) onLocaleSuccess() {
	lang.lhelper.ConnectSuccess(func(ok bool, reason string) {
		err := lang.handleLocaleChanged(ok, reason)
		if err != nil {
			lang.logger.Warning(err)
			lang.setPropCurrentLocale(getCurrentUserLocale())
			e := sendNotify(localeIconFailed, "",
				Tr("System language failed to change, please try later"))
			if e != nil {
				lang.logger.Warning("sendNotify failed:", e)
			}
			e = syncUserLocale(lang.CurrentLocale)
			if e != nil {
				lang.logger.Warning("Sync user object locale failed:", e)
			}
			lang.LocaleState = LocaleStateChanged
			return
		}
		e := sendNotify(localeIconFinished, "",
			Tr("System language has been changed, please log in again after logged out"))
		if e != nil {
			lang.logger.Warning("sendNotify failed:", e)
		}
		e = syncUserLocale(lang.CurrentLocale)
		if e != nil {
			lang.logger.Warning("Sync user object locale failed:", e)
		}
		lang.LocaleState = LocaleStateChanged
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

	fontCfgFile := filepath.Join(basedir.GetUserConfigDir(), "fontconfig/conf.d/99-deepin.conf")
	if err := os.Remove(fontCfgFile); err != nil {
		lang.logger.Warningf("remove font config file %q failed: %v", fontCfgFile, err)
	}

	err = installI18nDependent(lang.CurrentLocale)
	if err != nil {
		lang.logger.Warning(err)
	}
	dbus.Emit(lang, "Changed", lang.CurrentLocale)

	return nil
}

func syncUserLocale(locale string) error {
	cur, err := user.Current()
	if err != nil {
		return err
	}

	u, err := ddbus.NewUserByUid(cur.Uid)
	if err != nil {
		return err
	}
	err = u.SetLocale(locale)
	ddbus.DestroyUser(u)
	return err
}
