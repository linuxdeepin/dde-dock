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

func Test_windowInfosTypeEqual(t *testing.T) {
	Convey("windowInfosType Equal", t, func(c C) {
		wa := windowInfosType{
			0: {"a", false},
			1: {"b", false},
			2: {"c", true},
		}
		wb := windowInfosType{
			2: {"c", true},
			1: {"b", false},
			0: {"a", false},
		}
		c.So(wa.Equal(wb), ShouldBeTrue)

		wc := windowInfosType{
			1: {"b", false},
			2: {"c", false},
		}
		c.So(wc.Equal(wa), ShouldBeFalse)

		wd := windowInfosType{
			0: {"aa", false},
			1: {"b", false},
			2: {"c", false},
		}
		c.So(wd.Equal(wa), ShouldBeFalse)

		we := windowInfosType{
			0: {"a", false},
			1: {"b", false},
			3: {"c", false},
		}
		c.So(we.Equal(wa), ShouldBeFalse)

		wf := windowInfosType{
			0: {"a", false},
			1: {"b", false},
			2: {"c", false},
		}
		c.So(wf.Equal(wa), ShouldBeFalse)
	})
}
