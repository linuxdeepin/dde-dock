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
	"io/ioutil"
	"strconv"
	"strings"
)

type dstData struct {
	zone string
	dst  DSTInfo
}

func parseDSTDataFile(filename string) ([]dstData, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(contents), "\n")
	var infos []dstData
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		strs := strings.Split(line, ";")
		if len(strs) != 4 {
			continue
		}

		enter, _ := strconv.ParseInt(strs[1], 10, 64)
		leave, _ := strconv.ParseInt(strs[2], 10, 64)
		offset, _ := strconv.ParseInt(strs[3], 10, 64)
		info := dstData{
			zone: strs[0],
			dst: DSTInfo{
				Enter:     enter,
				Leave:     leave,
				DSTOffset: int32(offset),
			},
		}
		infos = append(infos, info)
	}

	return infos, nil
}

func findDSTInfo(zone, filename string) (*DSTInfo, error) {
	infos, err := parseDSTDataFile(filename)
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		if info.zone == zone {
			return &info.dst, nil
		}
	}

	return nil, errNoDST
}
