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

func sumDSTTime(zone string, year int32) (int64, int64, bool) {
	czone := C.CString(zone)
	ret := C.get_dst_time(czone, C.int(year))
	tmp := uintptr(unsafe.Pointer(ret))
	l := unsafe.Sizeof(*ret)

	first := int64(*ret)
	second := int64(*(*C.longlong)(unsafe.Pointer(tmp + uintptr(l))))
	third := int64(*(*C.longlong)(unsafe.Pointer(tmp + uintptr(l)*2)))

	C.free(unsafe.Pointer(czone))
	C.free(unsafe.Pointer(ret))

	if third != 2 {
		return 0, 0, false
	}

	return first, second, true
}

func getOffsetByTimestamp(zone string, timestamp int64) int32 {
	czone := C.CString(zone)
	off := C.getoffset(czone, C.longlong(timestamp))
	C.free(unsafe.Pointer(czone))

	return int32(off)
}

func getYearBeginTime(zone string, year int32) int64 {
	czone := C.CString(zone)
	timestamp := C.get_year_begin_time(czone, C.int(year))
	C.free(unsafe.Pointer(czone))

	return int64(timestamp)
}

func newDSTInfo(zone string) *DSTInfo {
	dst, err := findDSTInfo(zone, dstDataFile)
	if err == errNoDST || err == nil {
		return dst
	}

	year := time.Now().Year()

	first, second, ok := sumDSTTime(zone, int32(year))
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

func newZoneSummary(zone string) *ZoneSummary {
	var info ZoneSummary
	year := time.Now().Year()

	info.Name = zone
	info.Desc = getZoneDesc(zone)
	//info.DST = newDSTInfo(zone, year)

	off := getOffsetByTimestamp(zone,
		getYearBeginTime(zone, int32(year)))
	info.RawOffset = off

	return &info
}
