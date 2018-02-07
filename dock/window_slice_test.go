/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

func Test_diffSortedWindowSlice(t *testing.T) {
	a := windowSlice{1, 2, 3, 4}
	b := windowSlice{1, 3, 5, 6, 7}
	add, remove := diffSortedWindowSlice(a, b)

	Convey("diffSortedWindowSlice", t, func() {
		So(len(add), ShouldEqual, 3)
		So(add[0], ShouldEqual, 5)
		So(add[1], ShouldEqual, 6)
		So(add[2], ShouldEqual, 7)

		So(len(remove), ShouldEqual, 2)
		So(remove[0], ShouldEqual, 2)
		So(remove[1], ShouldEqual, 4)
	})
}
