/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package audio

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_isVolumeValid(t *testing.T) {
	Convey("isVolumeValid", t, func(c C) {
		c.So(isVolumeValid(0), ShouldBeTrue)
		c.So(isVolumeValid(-1), ShouldBeFalse)
	})
}
func Test_floatPrecision(t *testing.T) {
	Convey("floatPrecision", t, func(c C) {
		c.So(floatPrecision(3.1415926), ShouldEqual, 3.14)
		c.So(floatPrecision(2.718281828), ShouldEqual, 2.72)
	})
}
