/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"github.com/BurntSushi/xgb/xproto"
)

type windowTitlesType map[xproto.Window]string

func newWindowTitles() windowTitlesType {
	return make(windowTitlesType)
}

func (a windowTitlesType) Equal(b windowTitlesType) bool {
	if len(a) != len(b) {
		return false
	}
	for keyA, valA := range a {
		valB, okB := b[keyA]
		if okB {
			if valA != valB {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
