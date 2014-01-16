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
	"os"
	"strings"
)

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func getPathFromURI(uri string) string {
	if !strings.Contains(uri, "file:///") {
		return uri
	}
	return "/" + strings.TrimLeft(uri, "file:///")
}

func getFontThemes() []ThemeType {
	fontTheme := []ThemeType{}

	fontTheme = append(fontTheme, ThemeType{Name: "Deepin", Type: "system"})
	fontTheme = append(fontTheme, ThemeType{Name: "Deepin1", Type: "system"})
	fontTheme = append(fontTheme, ThemeType{Name: "Deepin2", Type: "system"})
	fontTheme = append(fontTheme, ThemeType{Name: "Deepin3", Type: "system"})
	return fontTheme
}

func getBackgroundFiles() []ThemeType {
	bgTheme := []ThemeType{}

	bgTheme = append(bgTheme, ThemeType{Name: "Deepin", Type: "system"})
	bgTheme = append(bgTheme, ThemeType{Name: "Deepin1", Type: "system"})
	bgTheme = append(bgTheme, ThemeType{Name: "Deepin2", Type: "system"})
	bgTheme = append(bgTheme, ThemeType{Name: "Deepin3", Type: "system"})
	return bgTheme
}
