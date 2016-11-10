/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
