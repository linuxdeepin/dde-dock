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

package power

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_checkTimeStabilized(t *testing.T) {
	Convey("checkTimeStabilized", t, func() {
		data := []uint64{
			9455,
			5467,
			3840,
			2962,
			2408,
			2408,
			1754,
			1698,
			1710,
			1675,
		}

		Convey("fewer than 3", func() {
			So(checkTimeStabilized(data[:2], data[2]), ShouldBeFalse)
		})
		Convey("not stablized", func() {
			So(checkTimeStabilized(data[:6], data[6]), ShouldBeFalse)
		})
		Convey("stablized", func() {
			So(checkTimeStabilized(data[:9], data[9]), ShouldBeTrue)
		})
	})
}
