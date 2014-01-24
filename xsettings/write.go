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
	"bytes"
	"dlib/logger"
	"encoding/binary"
	"github.com/BurntSushi/xgb/xproto"
)

func writeInterger(value uint32, datas *[]byte) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("recover err:", err)
		}
	}()

	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, byteOrder, value)
	tmp := buf.Bytes()

	l := len(tmp)
	for i := 0; i < l; i++ {
		*datas = append(*datas, tmp[i])
	}
}

func writeInterger2(value uint16, datas *[]byte) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("recover err:", err)
		}
	}()

	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, byteOrder, value)
	tmp := buf.Bytes()

	l := len(tmp)
	for i := 0; i < l; i++ {
		*datas = append(*datas, tmp[i])
	}
}

func writeString(value string, datas *[]byte) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("recover err:", err)
		}
	}()

	tmp := []byte(value)
	l := len(tmp)

	for i := 0; i < l; i++ {
		*datas = append(*datas, tmp[i])
	}
}

func writeBytes(values []byte, datas *[]byte) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("recover err:", err)
		}
	}()

	l := len(values)
	for i := 0; i < l; i++ {
		*datas = append(*datas, values[i])
	}
}

func writeHeader(info *HeaderInfo, datas *[]byte) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("recover err:", err)
		}
	}()

	*datas = append(*datas, info.vType)
	*datas = append(*datas, 0)
	writeInterger2(info.nameLen, datas)
	writeString(info.name, datas)

	leftPad := 3 - (info.nameLen+3)%4
	for i := uint16(0); i < leftPad; i++ {
		*datas = append(*datas, 0)
	}

	writeInterger(info.lastSerial, datas)
	switch info.vType {
	case XSETTINGS_INTERGER:
		writeInterger(info.value.(uint32), datas)
	case XSETTINGS_STRING:
		name := info.value.(string)
		l := uint32(len(name))
		writeInterger(l, datas)
		writeString(name, datas)
		leftPad := 3 - (l+3)%4
		for i := uint32(0); i < leftPad; i++ {
			*datas = append(*datas, 0)
		}
	case XSETTINGS_COLOR:
		writeBytes(info.value.([]byte), datas)
	}
}

func changeXSettingsProperty(datas []byte) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("recover err:", err)
		}
	}()

	err := xproto.ChangePropertyChecked(X, xproto.PropModeReplace,
		sReply.Owner,
		getAtom(X, XSETTINGS_SETTINGS),
		getAtom(X, XSETTINGS_SETTINGS),
		bytesDataFormat, uint32(len(datas)), datas).Check()
	if err != nil {
		logger.Printf("Change Property '%s' Failed: %s\n",
			XSETTINGS_SETTINGS, err)
		panic(err)
	}
}

func setXSettingsName(name string, value interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("recover err:", err)
		}
	}()

	infos := readXSettings()

	isExist := false
	for _, v := range infos {
		if v.name == name {
			isExist = true
			v.value = value
		}
	}
	if !isExist {
		info := newHeaderInfo(name, value)
		infos = append(infos, info)
	}

	writeXSettingsData(infos)
}

func writeXSettingsData(infos []*HeaderInfo) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("recover err:", err)
		}
	}()

	datas := []byte{}

	datas = append(datas, xsettingsInfo.order)
	for i := 0; i < 3; i++ {
		datas = append(datas, 0)
	}
	writeInterger(xsettingsInfo.serial, &datas)
	l := uint32(len(infos))
	writeInterger(l, &datas)

	for _, v := range infos {
		writeHeader(v, &datas)
	}

	changeXSettingsProperty(datas)
}

func newHeaderInfo(name string, value interface{}) *HeaderInfo {
	info := &HeaderInfo{}
	switch value.(type) {
	case uint32:
		info.vType = XSETTINGS_INTERGER
	case string:
		info.vType = XSETTINGS_STRING
	case []byte:
		info.vType = XSETTINGS_COLOR
	default:
		panic("type invalid")
	}

	info.name = name
	info.lastSerial = xsettingsInfo.serial
	info.nameLen = uint16(len(name))
	info.value = value

	return info
}
