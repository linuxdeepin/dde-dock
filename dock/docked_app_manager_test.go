package dock

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_uniqStrSlice(t *testing.T) {
	slice := []string{"a", "b", "c", "c", "b", "a", "c"}
	slice = uniqStrSlice(slice)
	Convey("uniqStrSlice", t, func() {
		So(len(slice), ShouldEqual, 3)
		So(slice[0], ShouldEqual, "a")
		So(slice[1], ShouldEqual, "b")
		So(slice[2], ShouldEqual, "c")
	})
}

func Test_strSliceEqual(t *testing.T) {
	sa := []string{"a", "b", "c"}
	sb := []string{"a", "b", "c", "d"}
	sc := sa[:]
	Convey("strSliceEqual", t, func() {
		So(strSliceEqual(sa, sb), ShouldBeFalse)
		So(strSliceEqual(sa, sc), ShouldBeTrue)
	})
}
