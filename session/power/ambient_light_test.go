package power

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_calcBrWithLightLevel(t *testing.T) {
	Convey("calcBrWithLightLevel", t, func() {
		var arr = []struct {
			lightLevel float64
			br         byte
		}{
			{-1, 0},
			{0, 0},
			{1, 2},
			{2, 3},
			{17, 29},
			{60, 48},
			{350, 62},
			{9999.9, 255},
			{10000, 255},
			{10000.1, 255},
		}

		for _, value := range arr {
			So(calcBrWithLightLevel(value.lightLevel), ShouldEqual, value.br)
		}
	})
}
