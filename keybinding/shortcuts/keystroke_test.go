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

package shortcuts

import (
	"testing"

	"github.com/linuxdeepin/go-x11-client/util/keysyms"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSplitKeystroke(t *testing.T) {
	Convey("splitKeystroke", t, func(c C) {
		var keys []string
		var err error
		keys, err = splitKeystroke("<Super>L")
		c.So(err, ShouldBeNil)
		c.So(keys, ShouldResemble, []string{"Super", "L"})

		// single key
		keys, err = splitKeystroke("<Super>")
		c.So(err, ShouldBeNil)
		c.So(keys, ShouldResemble, []string{"Super"})

		keys, err = splitKeystroke("Super_L")
		c.So(err, ShouldBeNil)
		c.So(keys, ShouldResemble, []string{"Super_L"})

		keys, err = splitKeystroke("<Shift><Super>T")
		c.So(err, ShouldBeNil)
		c.So(keys, ShouldResemble, []string{"Shift", "Super", "T"})

		// abnormal situation:
		_, err = splitKeystroke("<Super>>")
		c.So(err, ShouldNotBeNil)

		_, err = splitKeystroke("<Super><")
		c.So(err, ShouldNotBeNil)

		_, err = splitKeystroke("Super<")
		c.So(err, ShouldNotBeNil)

		_, err = splitKeystroke("<Super><shiftT")
		c.So(err, ShouldNotBeNil)

		_, err = splitKeystroke("<Super><Shift><>T")
		c.So(err, ShouldNotBeNil)
	})
}

func TestParseKeystroke(t *testing.T) {
	Convey("ParseKeystroke", t, func(c C) {
		var ks *Keystroke
		var err error

		ks, err = ParseKeystroke("Super_L")
		c.So(err, ShouldBeNil)
		c.So(ks, ShouldResemble, &Keystroke{
			Keystr: "Super_L",
			Keysym: keysyms.XK_Super_L,
		})

		ks, err = ParseKeystroke("Num_Lock")
		c.So(err, ShouldBeNil)
		c.So(ks, ShouldResemble, &Keystroke{
			Keystr: "Num_Lock",
			Keysym: keysyms.XK_Num_Lock,
		})

		ks, err = ParseKeystroke("<Control><Super>T")
		c.So(err, ShouldBeNil)
		c.So(ks, ShouldResemble, &Keystroke{
			Keystr: "T",
			Keysym: keysyms.XK_T,
			Mods:   keysyms.ModMaskSuper | keysyms.ModMaskControl,
		})

		ks, err = ParseKeystroke("<Control><Alt><Shift><Super>T")
		c.So(err, ShouldBeNil)
		c.So(ks, ShouldResemble, &Keystroke{
			Keystr: "T",
			Keysym: keysyms.XK_T,
			Mods:   keysyms.ModMaskShift | keysyms.ModMaskSuper | keysyms.ModMaskAlt | keysyms.ModMaskControl,
		})

		// abnormal situation:
		_, err = ParseKeystroke("<Shift>XXXXX")
		c.So(err, ShouldNotBeNil)

		_, err = ParseKeystroke("")
		c.So(err, ShouldNotBeNil)

		_, err = ParseKeystroke("<lock><Shift>A")
		c.So(err, ShouldNotBeNil)
	})
}

func TestKeystrokeMethodString(t *testing.T) {
	Convey("Keystroke.String", t, func(c C) {
		var ks Keystroke
		ks = Keystroke{
			Keystr: "percent",
			Mods:   keysyms.ModMaskControl | keysyms.ModMaskShift,
		}
		c.So(ks.String(), ShouldEqual, "<Shift><Control>percent")

		ks = Keystroke{
			Keystr: "T",
			Mods:   keysyms.ModMaskShift | keysyms.ModMaskSuper | keysyms.ModMaskAlt | keysyms.ModMaskControl,
		}
		c.So(ks.String(), ShouldEqual, "<Shift><Control><Alt><Super>T")
	})
}
