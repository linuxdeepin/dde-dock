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

package inputdevices

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSystemLayout(t *testing.T) {
	Convey("Get system layout", t, func() {
		layout, err := getSystemLayout("testdata/keyboard")
		So(err, ShouldBeNil)
		So(layout, ShouldEqual, "us;")
	})
}

func TestParseXKBFile(t *testing.T) {
	Convey("Parse xkb rule file", t, func() {
		handler, err := getLayoutListByFile("testdata/base.xml")
		So(err, ShouldBeNil)
		So(handler, ShouldNotBeNil)
	})
}

func TestStrList(t *testing.T) {
	var list = []string{"abc", "xyz", "123"}
	Convey("Add item to list", t, func() {
		ret, added := addItemToList("456", list)
		So(len(ret), ShouldEqual, 4)
		So(added, ShouldEqual, true)

		ret, added = addItemToList("123", list)
		So(len(ret), ShouldEqual, 3)
		So(added, ShouldEqual, false)
	})

	Convey("Delete item from list", t, func() {
		ret, deleted := delItemFromList("123", list)
		So(len(ret), ShouldEqual, 2)
		So(deleted, ShouldEqual, true)

		ret, deleted = delItemFromList("456", list)
		So(len(ret), ShouldEqual, 3)
		So(deleted, ShouldEqual, false)
	})

	Convey("Is item in list", t, func() {
		So(isItemInList("123", list), ShouldEqual, true)
		So(isItemInList("456", list), ShouldEqual, false)
	})
}

func TestSyndaemonExist(t *testing.T) {
	Convey("Test syndaemon exist", t, func() {
		So(isSyndaemonExist("testdata/syndaemon.pid"), ShouldEqual, false)
		So(isProcessExist("testdata/dde-desktop-cmdline", "dde-desktop"),
			ShouldEqual, true)
	})
}

func TestCurveControlPoints(t *testing.T) {
	// output svg path for debug
	for i := 1; i <= 7; i++ {
		p := getPressureCurveControlPoints(i)
		fmt.Printf(
			`<path d="M0,0 C%v,%v %v,%v 100,100" stroke="red" fill="none" style="stroke-width: 2px;"></path>`,
			p[0], p[1], p[2], p[3])
		fmt.Println("")
	}
}
