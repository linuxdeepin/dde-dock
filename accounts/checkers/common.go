/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package checkers

func isStrInArray(str string, array []string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}

	return false
}
