/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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
	"fmt"

	"sync"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/composite"
	"github.com/linuxdeepin/go-x11-client/ext/damage"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
	"github.com/linuxdeepin/go-x11-client/util/wm/icccm"
)

type TrayIcon struct {
	win    x.Window
	notify bool
	data   []byte // window pixmap data
	damage damage.Damage
	mu     sync.Mutex
}

func NewTrayIcon(win x.Window) *TrayIcon {
	return &TrayIcon{
		win:    win,
		notify: true,
	}
}

func (icon *TrayIcon) getName() string {
	wmName, _ := ewmh.GetWMName(XConn, icon.win).Reply(XConn)
	if wmName != "" {
		return wmName
	}

	wmNameTextProp, err := icccm.GetWMName(XConn, icon.win).Reply(XConn)
	if err == nil {
		wmName, _ := wmNameTextProp.GetStr()
		if wmName != "" {
			return wmName
		}
	}

	wmClass, err := icccm.GetWMClass(XConn, icon.win).Reply(XConn)
	if err == nil {
		return fmt.Sprintf("[%s|%s]", wmClass.Class, wmClass.Instance)
	}

	return ""
}

func (icon *TrayIcon) getPixmapData() ([]byte, error) {
	pixmapId, err := XConn.GenerateID()
	if err != nil {
		return nil, err
	}
	pixmap := x.Pixmap(pixmapId)
	err = composite.NameWindowPixmapChecked(XConn, icon.win, pixmap).Check(XConn)
	if err != nil {
		return nil, err
	}
	defer x.FreePixmap(XConn, pixmap)

	geo, err := x.GetGeometry(XConn, x.Drawable(icon.win)).Reply(XConn)
	if err != nil {
		return nil, err
	}

	img, err := x.GetImage(XConn, x.ImageFormatZPixmap, x.Drawable(pixmap),
		0, 0, geo.Width, geo.Height, (1<<32)-1).Reply(XConn)
	if err != nil {
		return nil, err
	}
	return img.Data, nil
}
