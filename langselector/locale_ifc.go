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

type localeInfo struct {
	Locale string
	Desc   string
}

func (lang *LangSelector) SetLocale(locale string) error {
	if len(locale) == 0 || !language_info.IsLocaleValid(locale,
		language_info.LanguageListFile) {
		return fmt.Errorf("Invalid locale: %v", locale)
	}
	if lang.setDate == nil {
		return fmt.Errorf("SetDateTime object is nil")
	}

	if lang.CurrentLocale == locale {
		return nil
	}

	lang.LocaleState = LocaleStateChanging
	lang.setPropCurrentLocale(locale)
	go func() {
		if ok, _ := isNetworkEnable(); !ok {
			err := sendNotify("", "", Tr("System language is being changed, please wait..."))
			if err != nil {
				lang.logger.Warning("sendNotify failed:", err)
			}
		} else {
			err := sendNotify("", "", Tr("System language is being changed with an installation of lacked language packages, please wait..."))
			if err != nil {
				lang.logger.Warning("sendNotify failed:", err)
			}
		}
		lang.setDate.GenLocale(locale)
	}()

	return nil
}

func (lang *LangSelector) GetLocaleList() []localeInfo {
	list, err := getLocaleInfoList(language_info.LanguageListFile)
	if err != nil {
		lang.logger.Warning(err)
		return nil
	}

	return list
}

func getLocaleInfoList(filename string) ([]localeInfo, error) {
	infoList, err := language_info.GetLanguageInfoList(filename)
	if err != nil {
		return nil, err
	}

	var list []localeInfo
	for _, info := range infoList {
		tmp := localeInfo{info.Locale, info.Description}
		list = append(list, tmp)
	}

	return list, nil
}
