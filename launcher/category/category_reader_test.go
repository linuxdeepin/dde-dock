package category

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCategoryReader(t *testing.T) {
	Convey("test GetAllInfos", t, func() {
		infos := GetAllInfos("./testdata/categories.json")
		So(len(infos), ShouldEqual, 11)

		infos = GetAllInfos("")
		So(len(infos), ShouldEqual, 0)
	})
}
