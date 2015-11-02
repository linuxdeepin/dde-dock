package dstore

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetInstalledTime(t *testing.T) {
	Convey("DQueryTimeInstalledTransaction", t, func() {
		transition, err := NewDQueryTimeInstalledTransaction("testdata/installTime.json")
		So(err, ShouldEqual, nil)
		So(transition.Query("test"), ShouldEqual, 0)
		So(transition.Query("chmsee"), ShouldNotEqual, 0)
	})
}
