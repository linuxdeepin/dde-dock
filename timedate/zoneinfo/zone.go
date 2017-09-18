/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
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

package zoneinfo

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	dutils "pkg.deepin.io/lib/utils"
)

const (
	zoneInfoDir = "/usr/share/zoneinfo"
)

type DSTInfo struct {
	// The timestamp of entering DST every year
	Enter int64
	// The timestamp of leaving DST every year
	Leave int64

	// The DST offset
	Offset int32
}

type ZoneInfo struct {
	// Timezone name, ex: "Asia/Shanghai"
	Name string
	// Timezone description, ex: "上海"
	Desc string

	// Timezone offset
	Offset int32

	DST DSTInfo
}

var (
	_zoneList []string

	// Error, invalid timezone
	ErrZoneInvalid = fmt.Errorf("Invalid time zone")
	errZoneNoDST   = fmt.Errorf("The time zone has no DST info")

	defaultZoneTab = "/usr/share/zoneinfo/zone1970.tab"
	defaultZoneDir = "/usr/share/zoneinfo"
)

/*
func init() {
	if _zoneList != nil {
		return
	}

	_zoneList = GetAllZones()
}
*/

// Check timezone validity
func IsZoneValid(zone string) bool {
	if len(zone) == 0 {
		return false
	}

	file := path.Join(defaultZoneDir, zone)
	return dutils.IsFileExist(file)
}

func GetAllZones() []string {
	if _zoneList != nil {
		return _zoneList
	}

	list, _ := getZoneListFromFile(defaultZoneTab)
	return list
}

// Query timezone detail info by timezone
func GetZoneInfo(zone string) (*ZoneInfo, error) {
	if !IsZoneValid(zone) {
		return nil, ErrZoneInvalid
	}

	info := newZoneInfo(zone)

	return info, nil
}

func getZoneListFromFile(file string) ([]string, error) {
	lines, err := getUncommentedZoneLines(file)
	if err != nil {
		return nil, err
	}

	var list []string
	for _, line := range lines {
		strv := strings.Split(line, "\t")
		list = append(list, strv[2])
	}

	return list, nil
}

func getUncommentedZoneLines(file string) ([]string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var (
		lines = strings.Split(string(content), "\n")
		match = regexp.MustCompile(`^#`)
		ret   []string
	)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		if match.MatchString(line) {
			continue
		}

		ret = append(ret, line)
	}

	return ret, nil
}
