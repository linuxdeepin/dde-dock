package dstore

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetInstalledTime(t *testing.T) {
	Convey("DQueryTimeInstalledTransaction", t, func() {
		t, err := NewDQueryTimeInstalledTransaction("testdata/installTime.json")
		So(err, ShouldEqual, nil)
		So(t.Query("test"), ShouldEqual, 0)
		So(t.Query("chmsee"), ShouldNotEqual, 0)
	})
}
