/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package dock

import (
	"github.com/BurntSushi/xgb/xproto"
)

type ExportWindowInfo struct {
	Title string
	Flash bool
}

type windowInfosType map[xproto.Window]ExportWindowInfo

func newWindowInfos() windowInfosType {
	return make(windowInfosType)
}

func (a windowInfosType) Equal(b windowInfosType) bool {
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
