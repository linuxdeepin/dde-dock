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

package language_info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type LanguageCodeInfo struct {
	LangCode    string
	CountryCode string
	Variant     string
}

type LanguageInfo struct {
	Locale      string `json:"Locale"`
	Description string `json:"Description"`
	LangCode    string `json:"LangCode"`
	CountryCode string `json:"CountryCode"`
}

type languageInfoGroup struct {
	LanguageList []LanguageInfo `json:"LanguageList"`
}

const LanguageListFile = "/usr/share/dde-daemon/lang/support_languages.json"

var (
	ErrInvalidLocale = fmt.Errorf("Invalid Locale")
)

func GetLanguageInfoList(config string) ([]LanguageInfo, error) {
	contents, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var lang languageInfoGroup
	err = json.Unmarshal(contents, &lang)
	if err != nil {
		return nil, err
	}

	return lang.LanguageList, nil
}

func GetCodeInfoByLocale(locale, config string) (LanguageCodeInfo, error) {
	langList, err := GetLanguageInfoList(config)
	if err != nil {
		return LanguageCodeInfo{}, err
	}

	var found bool
	var codeInfo LanguageCodeInfo
	for _, info := range langList {
		if locale == info.Locale {
			found = true

			var variant string
			strs := strings.Split(locale, ".")
			tmps := strings.Split(strs[0], "@")
			if len(tmps) > 1 {
				variant = tmps[1]
			}

			codeInfo.LangCode = info.LangCode
			codeInfo.CountryCode = info.CountryCode
			codeInfo.Variant = variant
			break
		}
	}

	if !found {
		return LanguageCodeInfo{}, ErrInvalidLocale
	}

	return codeInfo, nil
}

func IsLocaleValid(locale, config string) bool {
	langList, err := GetLanguageInfoList(config)
	if err != nil {
		return false
	}

	for _, info := range langList {
		if locale == info.Locale {
			return true
		}
	}

	return false
}
