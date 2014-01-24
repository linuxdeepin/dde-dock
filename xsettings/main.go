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
	"encoding/binary"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"time"
)

type XSettingsInfo struct {
	order  byte
	serial uint32
}

type HeaderInfo struct {
	vType      byte
	nameLen    uint16
	name       string
	lastSerial uint32
	value      interface{}
}

const (
	XSETTINGS_S0       = "_XSETTINGS_S0"
	XSETTINGS_SETTINGS = "_XSETTINGS_SETTINGS"

	XSETTINGS_INTERGER = 0
	XSETTINGS_STRING   = 1
	XSETTINGS_COLOR    = 2
)

var (
	X               *xgb.Conn
	sReply          *xproto.GetSelectionOwnerReply
	byteOrder       binary.ByteOrder
	bytesDataFormat byte
	xsettingsInfo   *XSettingsInfo
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

func initXSettings() {
	var err error

	X, err = xgb.NewConn()
	if err != nil {
		logger.Println("Unable to connect X server:", err)
		panic(err)
	}

	sReply, err = xproto.GetSelectionOwner(X,
		getAtom(X, XSETTINGS_S0)).Reply()
	if err != nil {
		logger.Println("Unable to connect X server:", err)
		panic(err)
	}

	xsettingsInfo = &XSettingsInfo{}
}

func main() {
	initXSettings()
	logger.Println("Deepin-Legacy")
	setXSettingsName("Net/ThemeName", "Deepin-Legacy")
	setXSettingsName("Net/IconThemeName", "Deepin-Legacy")
	time.Sleep(time.Second * 5)
	logger.Println("Adwaita, gnome")
	setXSettingsName("Net/ThemeName", "Adwaita")
	setXSettingsName("Net/IconThemeName", "gnome")
	time.Sleep(time.Second * 5)
	logger.Println("Deepin")
	setXSettingsName("Net/ThemeName", "Deepin")
	setXSettingsName("Net/IconThemeName", "Deepin")
}
