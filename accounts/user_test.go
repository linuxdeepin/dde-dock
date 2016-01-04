package accounts

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSystemLocale(t *testing.T) {
	Convey("Test system locale", t, func() {
		So(getSystemLanguage("testdata/locale"), ShouldEqual, "zh_CN")
	})
}
