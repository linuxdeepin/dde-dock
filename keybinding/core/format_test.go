package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFormatAccel(t *testing.T) {
	Convey("Test format accel", t, func() {
		So(FormatAccel("control-%"),
			ShouldEqual, "<Control>%")
		So(FormatAccel("control-shift-%"),
			ShouldEqual, "<Control>%")
		So(FormatAccel("<Control><Alt>T"),
			ShouldEqual, "<Control><Alt>T")
	})
}

func TestXGBFormat(t *testing.T) {
	Convey("Test xgb format", t, func() {
		So(formatAccelToXGB("<Control>%"),
			ShouldEqual, "control-shift-percent")
		So(formatAccelToXGB("<Control><Alt>%"),
			ShouldEqual, "control-mod1-shift-percent")
	})
}

func TestIsAccelEqual(t *testing.T) {
	_, err := Initialize()
	if err != nil {
		return
	}

	Convey("Test accel equal", t, func() {
		So(IsAccelEqual("control-%", "control-shift-%"),
			ShouldEqual, true)
		So(IsAccelEqual("<Control>%", "control-shift-%"),
			ShouldEqual, true)
		So(IsAccelEqual("<Control>%", "control-t"),
			ShouldEqual, false)
	})
}

func TestIsKeyMatch(t *testing.T) {
	_, err := Initialize()
	if err != nil {
		return
	}

	Convey("Test key match", t, func() {
		So(IsKeyMatch("space", 0, 65), ShouldEqual, true)
		So(IsKeyMatch("control-space", 4, 65), ShouldEqual, true)
	})
}
