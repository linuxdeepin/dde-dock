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

package dock

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_uniqStrSlice(t *testing.T) {
	slice := []string{"a", "b", "c", "c", "b", "a", "c"}
	slice = uniqStrSlice(slice)
	Convey("uniqStrSlice", t, func(c C) {
		c.So(len(slice), ShouldEqual, 3)
		c.So(slice[0], ShouldEqual, "a")
		c.So(slice[1], ShouldEqual, "b")
		c.So(slice[2], ShouldEqual, "c")
	})
}

func Test_strSliceEqual(t *testing.T) {
	sa := []string{"a", "b", "c"}
	sb := []string{"a", "b", "c", "d"}
	sc := sa[:]
	Convey("strSliceEqual", t, func(c C) {
		c.So(strSliceEqual(sa, sb), ShouldBeFalse)
		c.So(strSliceEqual(sa, sc), ShouldBeTrue)
	})
}

func Test_strSliceContains(t *testing.T) {
	Convey("strSliceContains", t, func(c C) {
		slice := []string{"a", "b", "c"}
		c.So(strSliceContains(slice, "a"), ShouldBeTrue)
		c.So(strSliceContains(slice, "b"), ShouldBeTrue)
		c.So(strSliceContains(slice, "c"), ShouldBeTrue)
		c.So(strSliceContains(slice, "d"), ShouldBeFalse)
		c.So(strSliceContains(slice, "e"), ShouldBeFalse)
	})

}
