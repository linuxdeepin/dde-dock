/*
 * Copyright (C) 2020 ~ 2021 Deepin Technology Co., Ltd.
 *
 * Author:     weizhixiang <1138871845@qq.com>
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

package gesture

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_isInWindowBlacklist(t *testing.T) {
	Convey("isInWindowBlacklist", t, func(c C) {
		slice := []string{"window1", "window2", "window3"}
		c.So(isInWindowBlacklist("window1", slice), ShouldBeTrue)
		c.So(isInWindowBlacklist("window2", slice), ShouldBeTrue)
		c.So(isInWindowBlacklist("window3", slice), ShouldBeTrue)
		c.So(isInWindowBlacklist("window4", slice), ShouldBeFalse)
	})
}
