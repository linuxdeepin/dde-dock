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

package utils

func IsStrInList(value string, list []string) bool {
	for _, s := range list {
		if s == value {
			return true
		}
	}

	return false
}

func IsYearValid(value int32) bool {
	if value >= 1970 {
		return true
	}

	return false
}

func IsMonthValid(value int32) bool {
	if value > 0 && value < 13 {
		return true
	}

	return false
}

func IsDayValid(year, month, value int32) bool {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		if value > 0 && value < 32 {
			return true
		}
	case 4, 6, 9, 11:
		if value > 0 && value < 31 {
			return true
		}
	case 2:
		if IsLeapYear(year) {
			if value > 0 && value < 30 {
				return true
			}
		} else {
			if value > 0 && value < 29 {
				return true
			}
		}
	}

	return false
}

func IsHourValid(value int32) bool {
	if value >= 0 && value < 24 {
		return true
	}

	return false
}

func IsMinuteValid(value int32) bool {
	if value >= 0 && value < 60 {
		return true
	}

	return false
}

func IsSecondValid(value int32) bool {
	if value >= 0 && value < 61 {
		return true
	}

	return false
}

func IsLeapYear(value int32) bool {
	if value%400 == 0 ||
		(value%4 == 0 && value%100 != 0) {
		return true
	}
	return false
}
