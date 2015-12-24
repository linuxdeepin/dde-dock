package dstore

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetPkgName(t *testing.T) {
	Convey("GetPkgName", t, func() {
		t, err := NewDQueryPkgNameTransaction("testdata/package.json")
		So(err, ShouldBeNil)
		So(t.Query("test.desktop"), ShouldEqual, "")
		So(t.Query("Thunar.desktop"), ShouldEqual, "thunar")
	})
}
