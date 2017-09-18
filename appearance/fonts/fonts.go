/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package fonts

import (
	"strings"
)

type Font struct {
	Id         string
	Name       string
	Family     string
	FamilyName string
	File       string

	Styles []string
	Lang   []string

	Monospace bool
	Deletable bool
}
type Fonts []*Font

// Some fonts do not follow the standard,
// so add a whitelist to handle these fonts.
var idWhiteList = []string{
	"NSimSun-18030",
}

func ListFont() Fonts {
	return fcInfosToFonts()
}

func (infos Fonts) ListStandard() Fonts {
	var ret Fonts
	for _, info := range infos {
		if !info.supportedCurLang() {
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func (infos Fonts) ListMonospace() Fonts {
	var ret Fonts
	for _, info := range infos {
		if !info.Monospace {
			continue
		}

		ret = append(ret, info)
	}
	return ret
}

func (infos Fonts) Get(id string) *Font {
	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return nil
}

func (infos Fonts) convertToFamilies() Families {
	var ret Families
	for _, info := range infos {
		if isItemInList(info.Id, idWhiteList) {
			ret = ret.add(&Family{
				Id:     info.Id,
				Name:   info.Name,
				Styles: info.Styles,
				//Files:  []string{info.File},
			})
			continue
		}

		ret = ret.add(&Family{
			Id:     info.Family,
			Name:   info.FamilyName,
			Styles: info.Styles,
			//Files:  []string{info.File},
		})
	}
	return ret
}

func (info *Font) supportedCurLang() bool {
	lang := getCurLang()
	// 由于 FcFontList 返回的结果中 lang 字段与利用 fc-query 方法获取的不同，比如有个字体的 lang 字段就丢失了 zh-cn 。
	// 这是有可能 FontConfig 的bug，只能暂时这样解决。
	if lang == "zh-cn" {
		lang = "zh"
	}
	for _, v := range info.Lang {
		if strings.HasPrefix(v, lang) {
			return true
		}
	}
	return false
}
