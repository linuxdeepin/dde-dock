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
