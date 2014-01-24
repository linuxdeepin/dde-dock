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
	"dlib/logger"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func getAtom(X *xgb.Conn, name string) xproto.Atom {
	reply, err := xproto.InternAtom(X, false,
		uint16(len(name)), name).Reply()
	if err != nil {
		logger.Printf("'%s' Get Xproto Atom Failed: %s\n",
			name, err)
	}

	return reply.Atom
}

func newXWindow() {
	var err error
	m.wid, err = xproto.NewWindowId(X)
	if err != nil {
		logger.Println("New Window Id Failed:", err)
		panic(err)
	}

	setupInfo := xproto.Setup(X)
	/*
	   for _, screenInfo := setupInfo.Roots {
	   }
	*/
	screen := setupInfo.DefaultScreen(X)
	err = xproto.CreateWindowChecked(X,
		xproto.WindowClassCopyFromParent,
		m.wid, screen.Root, 0, 0,
		150, 150, 10, xproto.WindowClassInputOnly,
		screen.RootVisual, 0, nil).Check()
	if err != nil {
		logger.Println("Create Window Failed:", err)
		panic(err)
	}
}
