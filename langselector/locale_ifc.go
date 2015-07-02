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
	"fmt"
	"pkg.linuxdeepin.com/dde-daemon/langselector/language_info"
	. "pkg.linuxdeepin.com/lib/gettext"
)

type LocaleInfo struct {
	// Locale name
	Locale string
	// Locale description
	Desc string
}

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

	if len(locale) == 0 || !language_info.IsLocaleValid(locale,
		language_info.LanguageListFile) {
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
			err := sendNotify("", "",
				Tr("System language is being changed, please wait..."))
			if err != nil {
				lang.logger.Warning("sendNotify failed:", err)
			}
		} else {
			err := sendNotify("", "",
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
	list, err := getLocaleInfoList(language_info.LanguageListFile)
	if err != nil {
		lang.logger.Warning(err)
		return nil
	}

	return list
}

func getLocaleInfoList(filename string) ([]LocaleInfo, error) {
	infoList, err := language_info.GetLanguageInfoList(filename)
	if err != nil {
		return nil, err
	}

	var list []LocaleInfo
	for _, info := range infoList {
		tmp := LocaleInfo{info.Locale, info.Description}
		list = append(list, tmp)
	}

	return list, nil
}
