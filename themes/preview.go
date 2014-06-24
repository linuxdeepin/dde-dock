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

package themes

import (
	dutils "pkg.linuxdeepin.com/lib/utils"
	"path"
)

func getGtkPreview(name string) string {
	list := getGtkThemeList()

	for _, l := range list {
		if name == l.Name {
			dest := ""
			if l.T == THEME_TYPE_SYSTEM {
				dest = path.Join(PERSON_SYS_THUMB_PATH,
					"WindowThemes", name+"-preview.png")
			} else {
				dest = path.Join(PERSON_LOCAL_THUMB_PATH,
					"WindowThemes", name+"-preview.png")
			}
			if dutils.IsFileExist(dest) {
				return dest
			}
			break
		}
	}

	return ""
}

func getIconPreview(name string) string {
	list := getIconThemeList()

	for _, l := range list {
		if name == l.Name {
			dest := ""
			if l.T == THEME_TYPE_SYSTEM {
				dest = path.Join(PERSON_SYS_THUMB_PATH,
					"IconThemes", name+"-preview.png")
			} else {
				dest = path.Join(PERSON_LOCAL_THUMB_PATH,
					"IconThemes", name+"-preview.png")
			}
			if dutils.IsFileExist(dest) {
				return dest
			}
			break
		}
	}

	return ""
}

func getCursorPreview(name string) string {
	list := getCursorThemeList()

	for _, l := range list {
		if name == l.Name {
			dest := ""
			if l.T == THEME_TYPE_SYSTEM {
				dest = path.Join(PERSON_SYS_THUMB_PATH,
					"CursorThemes", name+"-preview.png")
			} else {
				dest = path.Join(PERSON_LOCAL_THUMB_PATH,
					"CursorThemes", name+"-preview.png")
			}
			if dutils.IsFileExist(dest) {
				return dest
			}
			break
		}
	}

	return ""
}
