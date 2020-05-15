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
	"pkg.deepin.io/lib/dbus"
)

func Test_getValidName(t *testing.T) {
	Convey("getValidName", t, func(c C) {
		names := []string{"BAT0", "test.t", "test:t", "test-t", "test.1:2-3.4:5-6"}
		for _, name := range names {
			path := dbus.ObjectPath("/battery_" + getValidName(name))
			t.Log(path)
			c.So(path.IsValid(), ShouldBeTrue)
		}
	})
}
