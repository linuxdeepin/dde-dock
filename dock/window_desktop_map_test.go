package dock

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_windowDesktopMapLoad(t *testing.T) {
	Convey("windowDesktopMapLoad", t, func() {
		m, err := newWindowDesktopMapFromFile("./testdata/window-desktop-map.gob")
		So(err, ShouldBeNil)
		if m != nil {
			t.Logf("m.content: %#v", m.content)

			m.NewRel("winA1", "desktopA")
			m.NewRel("winA2", "desktopA")
			m.NewRel("winA3", "desktopA")
			m.NewRel("winB1", "desktopB")
			m.NewRel("winB2", "desktopB")
			m.NewRel("winC1", "desktopC")

			m.DelRel("winA1", "desktopA")

			t.Logf("m.content: %#v", m.content)
			m.Save()
		}
	})
}
