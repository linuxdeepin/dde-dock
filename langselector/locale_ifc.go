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
	"fmt"
	. "pkg.deepin.io/lib/gettext"
)

const (
	localeIconStart    = "notification-change-language-start"
	localeIconFailed   = "notification-change-language-failed"
	localeIconFinished = "notification-change-language-finished"
)

// Set user desktop environment locale, the new locale will work after relogin.
// (Notice: this locale is only for the current user.)
//
// 设置用户会话的 locale，注销后生效，此改变只对当前用户生效。
//
// locale: see '/etc/locale.gen'
func (lang *LangSelector) SetLocale(locale string) error {
	if lang.LocaleState == LocaleStateChanging {
		return nil
	}

	if len(locale) == 0 || !lang.isSupportedLocale(locale) {
		return fmt.Errorf("Invalid locale: %v", locale)
	}
	if lang.lhelper == nil {
		return fmt.Errorf("LocaleHelper object is nil")
	}

	if lang.CurrentLocale == locale {
		return nil
	}

	go func() {
		lang.LocaleState = LocaleStateChanging
		lang.setPropCurrentLocale(locale)
		if ok, _ := isNetworkEnable(); !ok {
			err := sendNotify(localeIconStart, "",
				Tr("System language is being changed, please wait..."))
			if err != nil {
				lang.logger.Warning("sendNotify failed:", err)
			}
		} else {
			err := sendNotify(localeIconStart, "",
				Tr("System language is being changed with an installation of lacked language packages, please wait..."))
			if err != nil {
				lang.logger.Warning("sendNotify failed:", err)
			}
		}
		err := lang.lhelper.GenerateLocale(locale)
		if err != nil {
			lang.logger.Warning("GenerateLocale failed:", err)
			lang.LocaleState = LocaleStateChanged
		}
	}()

	return nil
}

// Get locale info list that deepin supported
//
// 得到系统支持的 locale 信息列表
func (lang *LangSelector) GetLocaleList() []LocaleInfo {
	return lang.getCachedLocales()
}

func (lang *LangSelector) GetLocaleDescription(locale string) (string, error) {
	infos := lang.getCachedLocales()
	info, err := infos.Get(locale)
	if err != nil {
		return "", err
	}
	return info.Desc, nil
}

// Reset set user desktop environment locale to system default locale
func (lang *LangSelector) Reset() error {
	locale, err := getLocaleFromFile(systemLocaleFile)
	if err != nil {
		return err
	}
	return lang.SetLocale(locale)
}
