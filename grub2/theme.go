/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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

package grub2

import (
	"sync"

	"pkg.deepin.io/lib/dbusutil"
)

const DefaultThemeDir = "/boot/grub/themes/deepin"

var (
	themeDir    = DefaultThemeDir
	themeBgFile = themeDir + "/background.png"
)

// Theme is a dbus object which provide properties and methods to
// setup deepin grub2 theme.
type Theme struct {
	g        *Grub2
	service  *dbusutil.Service
	themeDir string
	bgFile   string

	PropsMu sync.RWMutex
	methods *struct {
		SetBackgroundSourceFile func() `in:"filename"`
		GetBackground           func() `out:"background"`
	}

	signals *struct {
		BackgroundChanged struct{}
	}
}

// NewTheme create Theme object.
func NewTheme(g *Grub2) *Theme {
	theme := &Theme{}
	theme.g = g
	theme.service = g.service
	theme.themeDir = themeDir
	theme.bgFile = themeBgFile

	return theme
}
