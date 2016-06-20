/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_uniqStrSlice(t *testing.T) {
	slice := []string{"a", "b", "c", "c", "b", "a", "c"}
	slice = uniqStrSlice(slice)
	Convey("uniqStrSlice", t, func() {
		So(len(slice), ShouldEqual, 3)
		So(slice[0], ShouldEqual, "a")
		So(slice[1], ShouldEqual, "b")
		So(slice[2], ShouldEqual, "c")
	})
}

func Test_strSliceEqual(t *testing.T) {
	sa := []string{"a", "b", "c"}
	sb := []string{"a", "b", "c", "d"}
	sc := sa[:]
	Convey("strSliceEqual", t, func() {
		So(strSliceEqual(sa, sb), ShouldBeFalse)
		So(strSliceEqual(sa, sc), ShouldBeTrue)
	})
}

func Test_strSliceContains(t *testing.T) {
	Convey("strSliceContains", t, func() {
		slice := []string{"a", "b", "c"}
		So(strSliceContains(slice, "a"), ShouldBeTrue)
		So(strSliceContains(slice, "b"), ShouldBeTrue)
		So(strSliceContains(slice, "c"), ShouldBeTrue)
		So(strSliceContains(slice, "d"), ShouldBeFalse)
		So(strSliceContains(slice, "e"), ShouldBeFalse)
	})

}
