package mpris

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
