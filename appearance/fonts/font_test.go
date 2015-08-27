package fonts

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCompositeList(t *testing.T) {
	Convey("Test list composition", t, func() {
		So(compositeList([]string{"123", "234"}, []string{"234", "abc"}),
			ShouldResemble, []string{"123", "234", "abc"})
	})
}

func TestFontFamily(t *testing.T) {
	Convey("Test font family", t, func() {
		var families = Families{
			&Family{
				Id:     "Source Code Pro",
				Name:   "Source Code Pro",
				Styles: []string{"Regular", "Bold"},
			},
			&Family{
				Id:     "WenQuanYi Micro Hei",
				Name:   "文泉译微米黑",
				Styles: []string{"Normal", "Bold"},
			},
		}
		So(families.GetIds(), ShouldResemble,
			[]string{"Source Code Pro",
				"WenQuanYi Micro Hei"})
		So(families.Get("WenQuanYi Micro Hei").Id, ShouldEqual,
			"WenQuanYi Micro Hei")
	})
}
