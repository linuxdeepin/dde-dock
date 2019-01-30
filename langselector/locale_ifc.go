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
	"fmt"

	"pkg.deepin.io/dde/api/language_support"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusPath      = "/com/deepin/daemon/LangSelector"
	dbusInterface = "com.deepin.daemon.LangSelector"

	localeIconStart    = "notification-change-language-start"
	localeIconFailed   = "notification-change-language-failed"
	localeIconFinished = "notification-change-language-finished"
)

func (*LangSelector) GetInterfaceName() string {
	return dbusInterface
}

// Set user desktop environment locale, the new locale will work after relogin.
// (Notice: this locale is only for the current user.)
//
// 设置用户会话的 locale，注销后生效，此改变只对当前用户生效。
//
// locale: see '/etc/locale.gen'
func (lang *LangSelector) SetLocale(locale string) *dbus.Error {
	lang.service.DelayAutoQuit()

	if !lang.isSupportedLocale(locale) {
		return dbusutil.ToError(fmt.Errorf("invalid locale: %v", locale))
	}

	lang.PropsMu.Lock()
	defer lang.PropsMu.Unlock()
	if lang.LocaleState == LocaleStateChanging || lang.CurrentLocale == locale {
		return nil
	}
	logger.Debugf("setLocale %q", locale)
	go lang.setLocale(locale)
	return nil
}

// Get locale info list that deepin supported
//
// 得到系统支持的 locale 信息列表
func (lang *LangSelector) GetLocaleList() ([]LocaleInfo, *dbus.Error) {
	lang.service.DelayAutoQuit()
	return lang.getCachedLocales(), nil
}

func (lang *LangSelector) GetLocaleDescription(locale string) (string, *dbus.Error) {
	lang.service.DelayAutoQuit()

	infos := lang.getCachedLocales()
	info, err := infos.Get(locale)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	return info.Desc, nil
}

// Reset set user desktop environment locale to system default locale
func (lang *LangSelector) Reset() *dbus.Error {
	lang.service.DelayAutoQuit()

	locale, err := getLocaleFromFile(systemLocaleFile)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return lang.SetLocale(locale)
}

func (lang *LangSelector) GetLanguageSupportPackages(locale string) ([]string, *dbus.Error) {
	ls, err := language_support.NewLanguageSupport()
	if err != nil {
		return nil, dbusutil.ToError(err)
	}

	pkgs := ls.ByLocale(locale, false)
	ls.Destroy()
	return pkgs, nil
}
