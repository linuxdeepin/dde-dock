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

	XSETTINGS_FORMAT = 8
	XSETTINGS_ORDER  = 0
	XSETTINGS_SERIAL = 0

	XSETTINGS_INTERGER = 0
	XSETTINGS_STRING   = 1
	XSETTINGS_COLOR    = 2
)

var (
	X         *xgb.Conn
	sReply    *xproto.GetSelectionOwnerReply
	byteOrder binary.ByteOrder
)

func initXSettings() {
	var err error

	X, err = xgb.NewConn()
	if err != nil {
		logger.Println("Unable to connect X server:", err)
		panic(err)
	}

	if XSETTINGS_ORDER == 1 {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}

	newXWindow()

	sReply, err = xproto.GetSelectionOwner(X,
		getAtom(X, XSETTINGS_S0)).Reply()
	if err != nil {
		logger.Println("Unable to connect X server:", err)
		panic(err)
	}
	logger.Println("select owner wid:", sReply.Owner)
}

func main() {
	initXSettings()
        logger.Println("Deepin-Legacy")
	setXSettingsName("Net/ThemeName", "Deepin-Legacy")
	setXSettingsName("Net/IconThemeName", "Deepin-Legacy")
	time.Sleep(time.Second * 10)
	logger.Println("Adwaita, gnome")
	setXSettingsName("Net/ThemeName", "Adwaita")
	setXSettingsName("Net/IconThemeName", "gnome")
	time.Sleep(time.Second * 10)
	logger.Println("Deepin")
	setXSettingsName("Net/ThemeName", "Deepin")
	setXSettingsName("Net/IconThemeName", "Deepin")

        select{}
}
