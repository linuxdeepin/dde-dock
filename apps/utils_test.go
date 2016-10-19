/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package apps

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func Test_getHomeByUid(t *testing.T) {
	Convey("getHomeByUid", t, func() {
		uid := os.Getuid()
		home, err := getHomeByUid(uid)
		So(err, ShouldBeNil)
		So(home, ShouldEqual, os.Getenv("HOME"))
	})
}
