/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dstore

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetInstalledTime(t *testing.T) {
	Convey("DQueryTimeInstalledTransaction", t, func() {
		t, err := NewDQueryTimeInstalledTransaction("testdata/installTime.json")
		So(err, ShouldEqual, nil)
		So(t.Query("test"), ShouldEqual, 0)
		So(t.Query("chmsee"), ShouldNotEqual, 0)
	})
}
