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

package utils

import (
	"fmt"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	WMThemeName    = "theme"
	WMTitlebarFont = "titlebar-font"
)

const (
	wmGSettingsSchema = "org.gnome.desktop.wm.preferences"
)

var (
	errSchemaNotExist = fmt.Errorf("GSettings schema not exist")
)

var _wmSettings *gio.Settings

func InitWMSettings() error {
	if _wmSettings != nil {
		return nil
	}

	if !dutils.IsGSchemaExist(wmGSettingsSchema) {
		return errSchemaNotExist
	}

	_wmSettings = gio.NewSettings(wmGSettingsSchema)
	return nil
}

func WMSetString(key, value string) bool {
	if _wmSettings == nil {
		return false
	}

	return _wmSettings.SetString(key, value)
}
