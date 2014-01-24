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
	"io"
)

func readInterger(buf io.Reader) uint32 {
	body := uint32(0)
	binary.Read(buf, m.order, &body)

	return body
}

func readColor(buf io.Reader) []uint16 {
	ret := []uint16{}
	var r uint16

	binary.Read(buf, m.order, &r)
	ret = append(ret, r)
	binary.Read(buf, m.order, &r)
	ret = append(ret, r)
	binary.Read(buf, m.order, &r)
	ret = append(ret, r)
	binary.Read(buf, m.order, &r)
	ret = append(ret, r)

	return ret
}

func readString(buf io.Reader) string {
	var nameLen uint32
	binary.Read(buf, m.order, &nameLen)
	if nameLen > 1000 {
		logger.Println("name len to long:", nameLen)
		panic("name len to long")
	}

	nameBuf := make([]byte, nameLen)
	binary.Read(buf, m.order, &nameBuf)

	leftPad := 3 - (nameLen+3)%4
	buf.Read(make([]byte, leftPad))

	return string(nameBuf)
}

func readString2(buf io.Reader) (string, uint16) {
	var nameLen uint16
	binary.Read(buf, m.order, &nameLen)

	nameBuf := make([]byte, nameLen)
	binary.Read(buf, m.order, &nameBuf)

	leftPad := 3 - (nameLen+3)%4
	buf.Read(make([]byte, leftPad))

	return string(nameBuf), nameLen
}

func readHeader(buf io.Reader) (byte, uint16, string, uint32) {
	var sType byte
	binary.Read(buf, m.order, &sType)
	buf.Read(make([]byte, 1))

	name, nameLen := readString2(buf)
	lastSerial := readInterger(buf)

	return sType, nameLen, name, lastSerial
}

func readXSettings() []*HeaderInfo {
	reply, err := xproto.GetProperty(X, false, sReply.Owner,
		getAtom(X, XSETTINGS_SETTINGS),
		getAtom(X, XSETTINGS_SETTINGS),
		0, 10240).Reply()
	if err != nil {
		logger.Println("Get Property Failed:", err)
		panic(err)
	}

	infos := []*HeaderInfo{}
	m.format = reply.Format
	data := reply.Value[:reply.ValueLen]
	xsettingsInfo.order = data[0]
	if data[0] == 1 {
		m.order = binary.BigEndian
	} else {
		m.order = binary.LittleEndian
	}

	buf := bytes.NewReader(data[4:])

	xsettingsInfo.serial = readInterger(buf)
	numSettings := readInterger(buf)

	//logger.Printf("serial: %d, numSettings: %d, suffix: %d\n",
	//serial, numSettings, suffix)

	for i := uint32(0); i < numSettings; i++ {
		sType, nameLen, name, lastSerial := readHeader(buf)
		info := &HeaderInfo{}
		info.vType = sType
		info.nameLen = nameLen
		info.name = name
		info.lastSerial = lastSerial
		switch sType {
		case XSETTINGS_INTERGER:
			v := readInterger(buf)
			//logger.Printf("%s = %d, start: %d, end: %d\n",
			//name, v, start, end)
			info.value = v
		case XSETTINGS_STRING:
			v := readString(buf)
			//logger.Printf("%s = %s, start: %d, end: %d\n",
			//name, v, start, end)
			info.value = v
		case XSETTINGS_COLOR:
			v := readColor(buf)
			//logger.Printf("%s = %d, %d, %d, %d, start: %d, end: %d\n",
			//name, v[0], v[1], v[2], v[3], start, end)
			info.value = v
		}
		infos = append(infos, info)
	}

	return infos
}
