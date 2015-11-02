package dstore

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetPkgName(t *testing.T) {
	Convey("GetPkgName", t, func() {
		transition, err := NewDQueryPkgNameTransaction("testdata/package.json")
		So(err, ShouldBeNil)
		So(transition.Query("test.desktop"), ShouldEqual, "")
		So(transition.Query("Thunar.desktop"), ShouldEqual, "thunar")
	})
}
