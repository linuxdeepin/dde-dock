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

package accounts

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetLocaleFromFile(t *testing.T) {
	Convey("getLocaleFromFile", t, func(c C) {
		c.So(getLocaleFromFile("testdata/locale"), ShouldEqual, "zh_CN.UTF-8")
	})
}

func TestSystemLayout(t *testing.T) {
	Convey("Get system layout", t, func(c C) {
		layout, err := getSystemLayout("testdata/keyboard_us")
		c.So(err, ShouldBeNil)
		c.So(layout, ShouldEqual, "us;")
		layout, _ = getSystemLayout("testdata/keyboard_us_chr")
		c.So(layout, ShouldEqual, "us;chr")
	})
}

func TestAvailableShells(t *testing.T) {
	Convey("Get available shells", t, func(c C) {
		var ret = []string{"/bin/sh", "/bin/bash",
			"/bin/zsh", "/usr/bin/zsh",
			"/usr/bin/fish",
		}
		shells := getAvailableShells("testdata/shells")
		c.So(shells, ShouldResemble, ret)
	})
}
