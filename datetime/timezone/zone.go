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

import (
	"dbus/com/deepin/api/setdatetime"
	"fmt"
	"path"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

type DSTInfo struct {
	Enter     int64
	Leave     int64
	DSTOffset int32
}

type ZoneInfo struct {
	Name      string
	Desc      string
	RawOffset int32
	DST       DSTInfo
}

var (
	ErrInvalidZone = fmt.Errorf("Invalid Timezone")

	errNoDST = fmt.Errorf("The zone has not dst info")
)

const (
	zoneInfoDir = "/usr/share/zoneinfo"
	dstDataFile = "/usr/share/dde-daemon/dst_data"
)

func IsZoneValid(zone string) bool {
	if len(zone) == 0 {
		return false
	}

	filename := path.Join(zoneInfoDir, zone)
	if dutils.IsFileExist(filename) {
		return true
	}

	return false
}

func SetTimezone(zone string) error {
	if !IsZoneValid(zone) {
		return ErrInvalidZone
	}

	datetime, err := setdatetime.NewSetDateTime(
		"com.deepin.api.SetDateTime",
		"/com/deepin/api/SetDateTime")
	if err != nil {
		return err
	}
	defer setdatetime.DestroySetDateTime(datetime)

	_, err = datetime.SetTimezone(zone)
	if err != nil {
		return err
	}

	return nil
}

var _infos []ZoneInfo

func GetZoneInfoList() []ZoneInfo {
	if _infos != nil {
		return _infos
	}

	for _, tmp := range zoneWhiteList {
		info := newZoneInfo(tmp.zone)
		_infos = append(_infos, *info)
	}

	return _infos
}

func GetZoneInfo(zone string) (*ZoneInfo, error) {
	if !IsZoneValid(zone) {
		return nil, ErrInvalidZone
	}

	info := newZoneInfo(zone)

	return info, nil
}

func getZoneDesc(zone string) string {
	for _, tmp := range zoneWhiteList {
		if zone == tmp.zone {
			return tmp.desc
		}
	}

	return zone
}
