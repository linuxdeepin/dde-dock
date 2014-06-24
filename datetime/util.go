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

package datetime

func timezoneIsValid(tz string) bool {
	for _, info := range zoneInfos {
		if tz == info.Zone {
			return true
		}
	}

	return false
}

func convertZoneToCity(tz string) string {
	city := ""
	for _, c := range tz {
		if c == '-' || c == '_' {
			city = city + " "
		} else {
			city = city + string(c)
		}
	}

	return city
}

func convertCityToZone(city string) string {
	tz := ""
	if isInUnderlineList(city) {
		for _, c := range city {
			if c == ' ' {
				tz = tz + "_"
			} else {
				tz = tz + string(c)
			}
		}
	} else {
		for _, c := range city {
			if c == ' ' {
				tz = tz + "-"
			} else {
				tz = tz + string(c)
			}
		}
	}

	return tz
}

func isInUnderlineList(key string) bool {
	ok := isElementExist(key, noUnderlineList)
	if ok {
		return false
	}

	return true
}

func isElementExist(element string, list []string) bool {
	for _, v := range list {
		if v == element {
			return true
		}
	}

	return false
}

func strArrayIsEqual(array1, array2 []string) bool {
	l1 := len(array1)
	l2 := len(array2)

	if l1 != l2 {
		return false
	}

	for i := 0; i < l1; i++ {
		if array1[i] != array2[i] {
			return false
		}
	}

	return true
}
