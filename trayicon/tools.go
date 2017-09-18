/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

package trayicon

import (
	x "github.com/linuxdeepin/go-x11-client"
)

func isValidWindow(win x.Window) bool {
	reply, err := x.GetWindowAttributes(XConn, win).Reply(XConn)
	return reply != nil && err == nil
}

func findRGBAVisualID() x.VisualID {
	screen := XConn.GetDefaultScreen()
	for _, dinfo := range screen.AllowedDepths {
		if dinfo.Depth == 32 {
			for _, vinfo := range dinfo.Visuals {
				return vinfo.VisualId
			}
		}
	}
	return screen.RootVisual
}
