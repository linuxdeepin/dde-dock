/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package fonts

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
	for _, v := range info.Lang {
		if lang == v {
			return true
		}
	}
	return false
}
