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

package keybinding

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_shouldUseDDEKwin(t *testing.T) {
	Convey("parseKeystrokes", t, func(c C) {
		_, err := os.Stat("/usr/bin/kwin_no_scale")
		exist1 := err == nil
		exist2 := shouldUseDDEKwin()
		c.So(exist1, ShouldEqual, exist2)
	})
}
