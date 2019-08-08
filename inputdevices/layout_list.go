/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package inputdevices

import (
	"encoding/xml"
	"io/ioutil"

	"pkg.deepin.io/dde/daemon/inputdevices/iso639"
	"pkg.deepin.io/lib/gettext"
	lib_locale "pkg.deepin.io/lib/locale"
	"pkg.deepin.io/lib/strv"
)

const (
	kbdLayoutsXml = "/usr/share/X11/xkb/rules/base.xml"
	kbdTextDomain = "xkeyboard-config"
)

type XKBConfigRegister struct {
	Layouts []XLayout `xml:"layoutList>layout"`
}

type XLayout struct {
	ConfigItem XConfigItem   `xml:"configItem"`
	Variants   []XConfigItem `xml:"variantList>variant>configItem"`
}

type XConfigItem struct {
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Languages   []string `xml:"languageList>iso639Id"`
}

func parseXML(filename string) (XKBConfigRegister, error) {
	var v XKBConfigRegister
	xmlByte, err := ioutil.ReadFile(filename)
	if err != nil {
		return v, err
	}

	err = xml.Unmarshal(xmlByte, &v)
	if err != nil {
		return v, err
	}

	return v, nil
}

type layoutMap map[string]layoutDetail

type layoutDetail struct {
	Languages   []string
	Description string
}

func getLayoutsFromFile(filename string) (layoutMap, error) {
	xmlData, err := parseXML(filename)
	if err != nil {
		return nil, err
	}

	result := make(layoutMap)
	for _, layout := range xmlData.Layouts {
		layoutName := layout.ConfigItem.Name
		desc := layout.ConfigItem.Description
		result[layoutName+layoutDelim] = layoutDetail{
			Languages:   layout.ConfigItem.Languages,
			Description: gettext.DGettext(kbdTextDomain, desc),
		}

		variants := layout.Variants
		for _, v := range variants {
			languages := v.Languages
			if len(v.Languages) == 0 {
				languages = layout.ConfigItem.Languages
			}
			result[layoutName+layoutDelim+v.Name] = layoutDetail{
				Languages:   languages,
				Description: gettext.DGettext(kbdTextDomain, v.Description),
			}
		}
	}

	return result, nil
}

func (layoutMap layoutMap) filterByLocales(locales []string) map[string]string {
	var localeLanguages []string
	for _, locale := range locales {
		components := lib_locale.ExplodeLocale(locale)
		lang := components.Language
		if lang != "" &&
			!strv.Strv(localeLanguages).Contains(lang) {
			localeLanguages = append(localeLanguages, lang)
		}
	}

	languages := make([]string, len(localeLanguages), 3*len(localeLanguages))
	copy(languages, localeLanguages)
	for _, t := range localeLanguages {
		a3Codes := iso639.ConvertA2ToA3(t)
		languages = append(languages, a3Codes...)
	}
	logger.Debug("languages:", languages)

	result := make(map[string]string)
	for layout, layoutDetail := range layoutMap {
		if layoutDetail.matchAnyLang(languages) {
			result[layout] = layoutDetail.Description
		}
	}
	return result
}

func (v *layoutDetail) matchAnyLang(languages []string) bool {
	for _, l := range languages {
		for _, ll := range v.Languages {
			if ll == l {
				return true
			}
		}
	}
	return false
}
