/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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
	"os"
	"os/user"
	"path/filepath"

	ddbus "pkg.deepin.io/dde/daemon/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/xdg/basedir"
)

var (
	// for locale-helper
	_ = Tr("Authentication is required to switch language")
)

func (lang *LangSelector) listenHelperSignal() {
	lang.helper.ConnectSuccess(func(ok bool, reason string) {
		err := lang.handleLocaleChanged(ok, reason)
		if err != nil {
			lang.logger.Warning(err)
			err := sendNotify(localeIconFailed, "",
				Tr("System language failed to change, please try later"))
			if err != nil {
				lang.logger.Warning("sendNotify failed:", err)
			}
			err = syncUserLocale(lang.CurrentLocale)
			if err != nil {
				lang.logger.Warning("Sync user object locale failed:", err)
			}

			// recover lang.CurrentLocale
			lang.PropsMu.Lock()
			lang.setPropCurrentLocale(getCurrentUserLocale())
			lang.setPropLocaleState(LocaleStateChanged)
			lang.PropsMu.Unlock()
			return
		}
		err = sendNotify(localeIconFinished, "",
			Tr("System language has been changed, please log in again after logged out"))
		if err != nil {
			lang.logger.Warning("sendNotify failed:", err)
		}

		lang.PropsMu.Lock()
		defer lang.PropsMu.Unlock()

		err = syncUserLocale(lang.CurrentLocale)
		if err != nil {
			lang.logger.Warning("Sync user object locale failed:", err)
		}

		lang.setPropLocaleState(LocaleStateChanged)
	})
}

func (lang *LangSelector) handleLocaleChanged(ok bool, reason string) error {
	if !ok {
		return ErrLocaleChangeFailed
	}

	lang.PropsMu.RLock()
	if lang.LocaleState != LocaleStateChanging {
		lang.PropsMu.RUnlock()
		return ErrLocaleChangeFailed
	}

	currentLocale := lang.CurrentLocale
	lang.PropsMu.RUnlock()

	err := writeUserLocale(currentLocale)
	if err != nil {
		return err
	}

	fontCfgFile := filepath.Join(basedir.GetUserConfigDir(), "fontconfig/conf.d/99-deepin.conf")
	if err := os.Remove(fontCfgFile); err != nil {
		lang.logger.Warningf("remove font config file %q failed: %v", fontCfgFile, err)
	}

	err = installLangSupportPackages(currentLocale)
	if err != nil {
		lang.logger.Warning(err)
	}
	lang.service.Emit(lang, "Changed", currentLocale)

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
