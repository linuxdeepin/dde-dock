package dock

import (
	"container/list"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_strListToSlice(t *testing.T) {
	l := list.New()
	l.PushBack("a")
	l.PushBack("b")
	l.PushBack("c")
	slice := strListToSlice(l)
	Convey("strListToSlice", t, func() {
		So(len(slice), ShouldEqual, 3)
		So(slice[0], ShouldEqual, "a")
		So(slice[1], ShouldEqual, "b")
		So(slice[2], ShouldEqual, "c")
	})
}

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
