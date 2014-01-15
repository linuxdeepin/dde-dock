/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package main

import (
	"encoding/xml"
	"io/ioutil"
)

type XKBConfigRegister struct {
	LayoutList XLayoutList `xml:"layoutList"`
}

type XLayoutList struct {
	Layout []XLayout `xml:"layout"`
}

type XLayout struct {
	ConfigItem  XConfigItem  `xml:"configItem"`
	VariantList XVariantList `xml:"variantList"`
}

type XConfigItem struct {
	Name        string `xml:"name"`
	Description string `xml:"description"`
}

type XVariantList struct {
	Variant []XVariant `xml:"variant"`
}

type XVariant struct {
	ConfigItem XConfigItem `xml:"configItem"`
}

const (
	_LAYOUT_XML_PATH = "/usr/share/X11/xkb/rules/base.xml"
)

func ParseXML(filename string) XKBConfigRegister {
	xmlByte, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var v XKBConfigRegister
	err = xml.Unmarshal(xmlByte, &v)
	if err != nil {
		panic(err)
	}

	return v
}

func GetLayoutList(xmlData XKBConfigRegister) map[string]string {
	layouts := make(map[string]string)

	for _, layout := range xmlData.LayoutList.Layout {
		firstName := layout.ConfigItem.Name
		desc := layout.ConfigItem.Description
		layouts[firstName] = desc

		variants := layout.VariantList.Variant
		for _, v := range variants {
			lastName := v.ConfigItem.Name
			descTmp := v.ConfigItem.Description
			keyName := firstName + " " + lastName
			layouts[keyName] = descTmp
		}
	}

	return layouts
}
