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

package shortcuts

import (
	"testing"

	x "github.com/linuxdeepin/go-x11-client"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSplitStandardAccel(t *testing.T) {
	Convey("splitStandardAccel", t, func() {
		var keys []string
		var err error
		keys, err = splitStandardAccel("<Super>L")
		So(err, ShouldBeNil)
		So(keys, ShouldResemble, []string{"Super", "L"})

		// single key
		keys, err = splitStandardAccel("<Super>")
		So(err, ShouldBeNil)
		So(keys, ShouldResemble, []string{"Super"})

		keys, err = splitStandardAccel("Super_L")
		So(err, ShouldBeNil)
		So(keys, ShouldResemble, []string{"Super_L"})

		keys, err = splitStandardAccel("<Shift><Super>T")
		So(err, ShouldBeNil)
		So(keys, ShouldResemble, []string{"Shift", "Super", "T"})

		// abnormal situation:
		keys, err = splitStandardAccel("<Super>>")
		So(err, ShouldNotBeNil)

		keys, err = splitStandardAccel("<Super><")
		So(err, ShouldNotBeNil)

		keys, err = splitStandardAccel("Super<")
		So(err, ShouldNotBeNil)

		keys, err = splitStandardAccel("<Super><shiftT")
		So(err, ShouldNotBeNil)

		keys, err = splitStandardAccel("<Super><Shift><>T")
		So(err, ShouldNotBeNil)
	})
}

func TestParseStandardAccel(t *testing.T) {
	Convey("ParseStandardAccel", t, func() {
		var parsed ParsedAccel
		var err error

		parsed, err = ParseStandardAccel("Super_L")
		So(err, ShouldBeNil)
		So(parsed, ShouldResemble, ParsedAccel{Key: "Super_L"})

		parsed, err = ParseStandardAccel("Num_Lock")
		So(err, ShouldBeNil)
		So(parsed, ShouldResemble, ParsedAccel{Key: "Num_Lock"})

		parsed, err = ParseStandardAccel("<Control><Super>T")
		So(err, ShouldBeNil)
		So(parsed, ShouldResemble, ParsedAccel{
			Key:  "T",
			Mods: x.ModMask4 | x.ModMaskControl,
		})

		parsed, err = ParseStandardAccel("<Control><Alt><Shift><Super>T")
		So(err, ShouldBeNil)
		So(parsed, ShouldResemble, ParsedAccel{
			Key:  "T",
			Mods: x.ModMaskShift | x.ModMask4 | x.ModMask1 | x.ModMaskControl,
		})

		parsed, err = ParseStandardAccel("<Shift>XXXXX")
		So(err, ShouldBeNil)
		So(parsed, ShouldResemble, ParsedAccel{Key: "XXXXX", Mods: x.ModMaskShift})

		// abnormal situation:
		parsed, err = ParseStandardAccel("")
		So(err, ShouldNotBeNil)

		parsed, err = ParseStandardAccel("<lock><Shift>A")
		So(err, ShouldNotBeNil)
	})
}

func TestParsedAccelMethodString(t *testing.T) {
	Convey("ParsedAccel.String", t, func() {
		var parsed ParsedAccel
		parsed = ParsedAccel{
			Key:  "percent",
			Mods: x.ModMaskControl | x.ModMaskShift,
		}
		So(parsed.String(), ShouldEqual, "<Shift><Control>percent")

		parsed = ParsedAccel{
			Key:  "T",
			Mods: x.ModMaskShift | x.ModMask4 | x.ModMask1 | x.ModMaskControl,
		}
		So(parsed.String(), ShouldEqual, "<Shift><Control><Alt><Super>T")
	})
}
