/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
