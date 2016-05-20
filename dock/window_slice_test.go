package dock

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_diffSortedWindowSlice(t *testing.T) {
	a := windowSlice{1, 2, 3, 4}
	b := windowSlice{1, 3, 5, 6, 7}
	add, remove := diffSortedWindowSlice(a, b)

	Convey("diffSortedWindowSlice", t, func() {
		So(len(add), ShouldEqual, 3)
		So(add[0], ShouldEqual, 5)
		So(add[1], ShouldEqual, 6)
		So(add[2], ShouldEqual, 7)

		So(len(remove), ShouldEqual, 2)
		So(remove[0], ShouldEqual, 2)
		So(remove[1], ShouldEqual, 4)
	})
}
