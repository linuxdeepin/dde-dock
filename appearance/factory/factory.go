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

package factory

import (
	"pkg.deepin.io/dde-daemon/appearance/factory/background"
	"pkg.deepin.io/dde-daemon/appearance/factory/cursor_theme"
	"pkg.deepin.io/dde-daemon/appearance/factory/deepin_theme"
	"pkg.deepin.io/dde-daemon/appearance/factory/greeter_theme"
	"pkg.deepin.io/dde-daemon/appearance/factory/gtk_theme"
	"pkg.deepin.io/dde-daemon/appearance/factory/icon_theme"
	"pkg.deepin.io/dde-daemon/appearance/factory/sound_theme"
	. "pkg.deepin.io/dde-daemon/appearance/utils"
)

type FactoryInterface interface {
	Set(value string) error
	Delete(value string) error
	Destroy()
	GetFlag(value string) int32
	GetNameStrList() []string
	GetThumbnail(value string) string
	IsValueValid(value string) bool
	GetInfoByName(name string) (PathInfo, error)
}

const (
	ObjectTypeGtk         int32 = 0
	ObjectTypeIcon              = 1
	ObjectTypeCursor            = 2
	ObjectTypeSound             = 3
	ObjectTypeGreeter           = 4
	ObjectTypeBackground        = 5
	ObjectTypeDeepinTheme       = 6
)

func NewFactory(objType int32, handler func([]string)) FactoryInterface {
	switch objType {
	case ObjectTypeGtk:
		return gtk_theme.NewGtkTheme(handler)
	case ObjectTypeIcon:
		return icon_theme.NewIconTheme(handler)
	case ObjectTypeCursor:
		return cursor_theme.NewCursorTheme(handler)
	case ObjectTypeSound:
		return sound_theme.NewSoundTheme(handler)
	case ObjectTypeGreeter:
		return greeter_theme.NewGreeterTheme(handler)
	case ObjectTypeBackground:
		return background.NewBackground(handler)
	case ObjectTypeDeepinTheme:
		return deepin_theme.NewDeepinTheme(handler)
	}

	return nil
}
