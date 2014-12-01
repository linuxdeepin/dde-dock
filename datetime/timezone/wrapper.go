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

package timezone

// #cgo CFLAGS: -Wall -g
// #include <stdlib.h>
// #include "timestamp.h"
import "C"

import (
	"time"
	"unsafe"
)

func getDSTTime(zone string, year int32) (int64, int64, bool) {
	czone := C.CString(zone)
	defer C.free(unsafe.Pointer(czone))
	ret := C.get_dst_time(czone, C.int(year))

	t1 := int64(ret.enter)
	t2 := int64(ret.leave)

	if t1 == 0 || t2 == 0 {
		return 0, 0, false
	}

	return t1, t2, true
}

func getOffsetByTimestamp(zone string, timestamp int64) int32 {
	czone := C.CString(zone)
	defer C.free(unsafe.Pointer(czone))
	off := C.getoffset(czone, C.longlong(timestamp))

	return int32(off)
}

func getRawTimestamp(zone string, timestamp int64) int64 {
	czone := C.CString(zone)
	defer C.free(unsafe.Pointer(czone))
	ret := C.get_rawoffset_time(czone, C.longlong(timestamp))

	return int64(ret)
}

func newDSTInfo(zone string) *DSTInfo {
	dst, err := findDSTInfo(zone, dstDataFile)
	if err == errNoDST || err == nil {
		return dst
	}

	year := time.Now().Year()

	first, second, ok := getDSTTime(zone, int32(year))
	if !ok {
		return nil
	}
	off := getOffsetByTimestamp(zone, first)

	return &DSTInfo{
		Enter:     first,
		Leave:     second,
		DSTOffset: off,
	}
}

func newZoneInfo(zone string) *ZoneInfo {
	var info ZoneInfo

	info.Name = zone
	info.Desc = getZoneDesc(zone)
	dst := newDSTInfo(zone)
	if dst == nil {
		info.RawOffset = getOffsetByTimestamp(zone, 0)
	} else {
		info.RawOffset = getOffsetByTimestamp(zone,
			getRawTimestamp(zone, dst.Enter))
		info.DST = *dst
	}

	return &info
}
