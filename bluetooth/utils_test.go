package bluetooth

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var str = []string{"/bin/sh", "/bin/bash",
	"/bin/zsh", "/usr/bin/zsh",
	"/usr/bin/fish",
}

func TestIsStringInArray(t *testing.T) {
	Convey("IsStringInArray", t, func(c C) {
		ret := isStringInArray("testdata/shells", str)
		c.So(ret, ShouldEqual, false)
		ret = isStringInArray("/bin/sh", str)
		c.So(ret, ShouldEqual, true)
	})
}

func TestMarshalJSON(t *testing.T) {
	str1 := make(map[string]string)
	str1["name"] = "uniontech"
	str1["addr"] = "wuhan"
	Convey("MarshalJSON", t, func(c C) {
		ret := marshalJSON(str1)
		c.So(ret, ShouldEqual, `{"addr":"wuhan","name":"uniontech"}`)
	})
}
