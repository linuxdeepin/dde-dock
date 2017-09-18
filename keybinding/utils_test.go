/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIBMHotkey(t *testing.T) {
	convey.Convey("Test ibm hotkey checker", t, func() {
		convey.So(checkIBMHotkey("testdata/hotkey"), convey.ShouldEqual, true)
		convey.So(checkIBMHotkey("testdata/hotkey_disable"), convey.ShouldEqual, false)
	})
}
