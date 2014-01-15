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

package main

import (
	"fmt"
	"testing"
)

func TestBind(t *testing.T) {
	InitVariable()
	bm := NewBindManager()

	fmt.Println("CustomList: ", bm.CustomList)
	fmt.Println("SystemList: ", bm.SystemList)
	fmt.Println("MediaList", bm.MediaList)
	fmt.Println("WindowList", bm.WindowList)
	fmt.Println("WorkSpaceList", bm.WorkSpaceList)
	fmt.Println("ValidList: ", bm.ConflictValid)
	fmt.Println("InvalidList: ", bm.ConflictInvalid)

	check := bm.ChangeShortcut(10000, "<Super>E")
	if check == nil {
		t.Error("ChangeShortcut Error")
	}
	if check.IsConflict {
		fmt.Println("Conflict idList: ", check.IdList)
	}
	fmt.Println("ValidList: ", bm.ConflictValid)
	fmt.Println("InvalidList: ", bm.ConflictInvalid)

	check = bm.ChangeShortcut(10000, "<Alt>E")
	if check == nil {
		t.Error("ChangeShortcut Error")
	}
	if check.IsConflict {
		fmt.Println("Conflict idList: ", check.IdList)
	}
	fmt.Println("ValidList: ", bm.ConflictValid)
	fmt.Println("InvalidList: ", bm.ConflictInvalid)
}
