/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package accounts

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetLocaleFromFile(t *testing.T) {
	Convey("getLocaleFromFile", t, func() {
		So(getLocaleFromFile("testdata/locale"), ShouldEqual, "zh_CN.UTF-8")
	})
}

func TestSystemLayout(t *testing.T) {
	Convey("Get system layout", t, func() {
		layout, err := getSystemLayout("testdata/keyboard_us")
		So(err, ShouldBeNil)
		So(layout, ShouldEqual, "us;")
		layout, _ = getSystemLayout("testdata/keyboard_us_chr")
		So(layout, ShouldEqual, "us;chr")
	})
}

func TestAvailableShells(t *testing.T) {
	Convey("Get available shells", t, func() {
		var ret = []string{"/bin/sh", "/bin/bash",
			"/bin/zsh", "/usr/bin/zsh",
			"/usr/bin/fish",
		}
		shells := getAvailableShells("testdata/shells")
		So(shells, ShouldResemble, ret)
	})
}
