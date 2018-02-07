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
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_windowTitlesTypeEqual(t *testing.T) {
	Convey("windowTitlesType Equal", t, func() {
		a := windowTitlesType{
			0: "a",
			1: "b",
			2: "c",
		}
		b := windowTitlesType{
			2: "c",
			1: "b",
			0: "a",
		}
		So(a.Equal(b), ShouldBeTrue)

		c := windowTitlesType{
			1: "b",
			2: "c",
		}
		So(c.Equal(a), ShouldBeFalse)

		d := windowTitlesType{
			0: "aa",
			1: "b",
			2: "c",
		}
		So(d.Equal(a), ShouldBeFalse)

		e := windowTitlesType{
			0: "a",
			1: "b",
			3: "c",
		}
		So(e.Equal(a), ShouldBeFalse)
	})
}
