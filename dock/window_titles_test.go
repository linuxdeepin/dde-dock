/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
