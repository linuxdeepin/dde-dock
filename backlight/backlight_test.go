package backlight

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"sort"
	"testing"
)

func TestSyspathInfo(t *testing.T) {
	Convey("Test syspath info", t, func() {
		info, _ := NewSyspathInfo("testdata/backlight/intel_backlight")
		So(info, ShouldResemble, &SyspathInfo{
			Path:          "testdata/backlight/intel_backlight",
			Type:          BacklightRaw,
			MaxBrightness: 100,
		})
	})
}

func TestSyspathInfos(t *testing.T) {
	Convey("Test syspath infos", t, func() {
		var infos = SyspathInfos{
			&SyspathInfo{
				Path:          "testdata/backlight/acpi_backlight",
				Type:          BacklightPlatform,
				MaxBrightness: 100,
			},
			&SyspathInfo{
				Path:          "testdata/backlight/intel_backlight",
				Type:          BacklightRaw,
				MaxBrightness: 100,
			},
		}
		infos = infos.sortLCD()
		So(infos[0].Path, ShouldEqual, "testdata/backlight/intel_backlight")
		So(infos[1].Path, ShouldEqual, "testdata/backlight/acpi_backlight")

		info, _ := infos.Get("testdata/backlight/intel_backlight")
		So(info.Path, ShouldEqual, "testdata/backlight/intel_backlight")
	})
}

func TestBrightness(t *testing.T) {
	Convey("Test brightness", t, func() {
		v, _ := doGetBrightness("testdata/backlight/intel_backlight/brightness")
		So(v, ShouldEqual, 20)
		v, _ = doGetBrightness("testdata/backlight/intel_backlight/max_brightness")
		So(v, ShouldEqual, 100)

		var testFile = "testdata/brightness"
		doSetBrightness("testdata", 50)
		v, _ = doGetBrightness(testFile)
		So(v, ShouldEqual, 50)
		os.Remove(testFile)
	})
}

func TestType(t *testing.T) {
	Convey("Test type", t, func() {
		ty, _ := getType("testdata/backlight/intel_backlight")
		So(ty, ShouldEqual, BacklightRaw)
		ty, _ = getType("testdata/backlight/acpi_backlight")
		So(ty, ShouldEqual, BacklightPlatform)
	})
}

func TestSyspathList(t *testing.T) {
	Convey("Test syspath list", t, func() {
		list := doListSyspath("testdata/backlight")
		sort.Strings(list)
		So(list, ShouldResemble, []string{
			"testdata/backlight/acpi_backlight",
			"testdata/backlight/intel_backlight",
		})
	})
}
