/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

type windowSlice []xproto.Window

func (a windowSlice) Len() int           { return len(a) }
func (a windowSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a windowSlice) Less(i, j int) bool { return uint32(a[i]) < uint32(a[j]) }

func (winSlice windowSlice) Contains(win xproto.Window) bool {
	for _, window := range winSlice {
		if window == win {
			return true
		}
	}
	return false
}

// from a to b
// return [add, remove]
func diffSortedWindowSlice(a, b windowSlice) (add, remove windowSlice) {
	ia := 0
	ib := 0
	lenA := len(a)
	lenB := len(b)

	for ia < lenA && ib < lenB {
		va := uint32(a[ia])
		vb := uint32(b[ib])
		if va == vb {
			ia++
			ib++
		} else if va < vb {
			// remove
			remove = append(remove, a[ia])
			ia++
		} else {
			// va > vb
			// add
			add = append(add, b[ib])
			ib++
		}
	}

	for ia < lenA {
		remove = append(remove, a[ia])
		ia++
	}

	for ib < lenB {
		add = append(add, b[ib])
		ib++
	}
	return
}
