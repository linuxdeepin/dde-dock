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

func TestGetPkgName(t *testing.T) {
	Convey("GetPkgName", t, func() {
		t, err := NewDQueryPkgNameTransaction("testdata/package.json")
		So(err, ShouldBeNil)
		So(t.Query("test.desktop"), ShouldEqual, "")
		So(t.Query("Thunar.desktop"), ShouldEqual, "thunar")
	})
}
