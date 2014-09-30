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

package fonts

import (
	. "pkg.linuxdeepin.com/dde-daemon/appearance/utils"
	xsettings "pkg.linuxdeepin.com/dde-daemon/xsettings_wrapper"
)

type FontManager struct {
	standardList  []StyleInfo
	monospaceList []StyleInfo
}

func NewFontManager() *FontManager {
	font := &FontManager{}

	font.standardList, font.monospaceList = getStyleInfoList()
	xsettings.InitXSettings()
	InitWMSettings()

	return font
}

func getNameStrList(infos []StyleInfo) []string {
	var list []string
	for _, info := range infos {
		list = append(list, info.Id)
	}

	return list
}

func getStyleList(name string, infos []StyleInfo) []string {
	for _, info := range infos {
		if name == info.Id {
			return info.StyleList
		}
	}

	return nil
}
