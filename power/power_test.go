package power

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProccessExist(t *testing.T) {
	Convey("Test proccess whether exist", t, func() {
		pid, err := getPidFromFile("testdata/init_pid")
		So(err, ShouldBeNil)
		So(pid, ShouldEqual, 1)
	})
}
