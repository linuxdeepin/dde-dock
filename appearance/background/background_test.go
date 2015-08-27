package background

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestScanner(t *testing.T) {
	Convey("Scanner bg", t, func() {
		So(scanner("testdata/Theme1/wallpapers"), ShouldResemble,
			[]string{
				"testdata/Theme1/wallpapers/desktop.jpg",
			})
		So(scanner("testdata/Theme2/wallpapers"), ShouldBeNil)
	})
}

func TestGetDirsFromDTheme(t *testing.T) {
	Convey("Get bg dirs from dtheme", t, func() {
		So(getDirsFromDTheme("testdata"), ShouldResemble,
			[]string{
				"testdata/Theme1/wallpapers",
				"testdata/Theme2/wallpapers"})
	})
}
