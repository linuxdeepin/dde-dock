/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProccessExist(t *testing.T) {
	Convey("Test proccess whether exist", t, func() {
		pid, err := getPidFromFile("testdata/init_pid")
		So(err, ShouldBeNil)
		So(pid, ShouldEqual, 1)
	})
}
