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
	. "github.com/smartystreets/goconvey/convey"
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"testing"
)

func Test_getAccelModKeys(t *testing.T) {
	Convey("Test getAccelModKeys", t, func() {
		pa, err := shortcuts.ParseStandardAccel("<Control>Alt_L")
		So(err, ShouldBeNil)
		So(getAccelModKeys(pa), ShouldResemble, []string{"Control", "Alt"})

		pa, err = shortcuts.ParseStandardAccel("<Control><Alt><Super>Shift_L")
		So(err, ShouldBeNil)
		So(getAccelModKeys(pa), ShouldResemble, []string{"Control", "Alt", "Super", "Shift"})

		pa, err = shortcuts.ParseStandardAccel("<Control><Alt><Super>X")
		So(err, ShouldBeNil)
		So(getAccelModKeys(pa), ShouldBeNil)
	})
}

func TestIBMHotkey(t *testing.T) {
	Convey("Test ibm hotkey checker", t, func() {
		So(checkIBMHotkey("testdata/hotkey"), ShouldBeTrue)
		So(checkIBMHotkey("testdata/hotkey_disable"), ShouldBeFalse)
	})
}
