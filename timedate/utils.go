/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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

package timedate

func isItemInList(item string, list []string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}

	return false
}

func addItemToList(item string, list []string) ([]string, bool) {
	if isItemInList(item, list) {
		return list, false
	}

	list = append(list, item)
	return list, true
}

func deleteItemFromList(item string, list []string) ([]string, bool) {
	var (
		ret   []string
		found bool = false
	)
	for _, v := range list {
		if v == item {
			found = true
			continue
		}

		ret = append(ret, v)
	}

	return ret, found
}

func filterNilString(list []string) ([]string, bool) {
	var (
		ret    []string
		hasNil bool = false
	)
	for _, v := range list {
		if len(v) == 0 {
			hasNil = true
			continue
		}
		ret = append(ret, v)
	}

	return ret, hasNil
}
