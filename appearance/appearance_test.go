package appearance

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_hexColorToXsColor(t *testing.T) {

	Convey("hexColorToXsColor", t, func(c C) {
		var tests = []struct {
			hex string
			xs  string
			err bool
		}{
			{"#ffffff", "65535,65535,65535,65535", false},
			{"#ffffffff", "65535,65535,65535,65535", false},
			{"#ff69b4", "65535,26985,46260,65535", false},
			{"#00000000", "0,0,0,0", false},
			{"abc", "", true},
			{"#FfFff", "", true},
			{"", "", true},
		}

		for _, test := range tests {
			xsColor, err := hexColorToXsColor(test.hex)
			c.So(err != nil, ShouldEqual, test.err)
			c.So(xsColor, ShouldEqual, test.xs)
		}
	})
}

func Test_xsColorToHexColor(t *testing.T) {

	Convey("xsColorToHexColor", t, func(c C) {
		var tests = []struct {
			hex string
			xs  string
			err bool
		}{
			{"#FFFFFF", "65535,65535,65535,65535", false},
			{"#FF69B4", "65535,26985,46260,65535", false},
			{"", "", true},
			{"", "123,456,678", true},
			{"", "-1,-2,03,45", true},
			{"", "65535,26985,46260,65535,", true},
			{"", "65535,26985,46260,165535", true},
		}

		for _, test := range tests {
			hexColor, err := xsColorToHexColor(test.xs)
			c.So(err != nil, ShouldEqual, test.err)
			c.So(hexColor, ShouldEqual, test.hex)
		}
	})
}

func Test_byteArrayToHexColor(t *testing.T) {
	Convey("byteArrayToHexColor", t, func(c C) {
		var tests = []struct {
			byteArray [4]byte
			hex       string
		}{
			{[4]byte{0xff, 0xff, 0xff}, "#FFFFFF00"},
			{[4]byte{0xff, 0xff, 0xff, 0xff}, "#FFFFFF"},
			{[4]byte{0xff, 0x69, 0xb4}, "#FF69B400"},
			{[4]byte{0xff, 0x69, 0xb4, 0xff}, "#FF69B4"},
		}

		for _, test := range tests {
			hexColor := byteArrayToHexColor(test.byteArray)
			c.So(hexColor, ShouldEqual, test.hex)
		}
	})
}

func Test_parseHexColor(t *testing.T) {
	Convey("parseHexColor", t, func(c C) {
		var tests = []struct {
			byteArray [4]byte
			hex       string
		}{
			{[4]byte{0xff, 0xff, 0xff}, "#FFFFFF00"},
			{[4]byte{0xff, 0xff, 0xff, 0xff}, "#FFFFFF"},
			{[4]byte{0xff, 0x69, 0xb4}, "#FF69B400"},
			{[4]byte{0xff, 0x69, 0xb4, 0xff}, "#FF69B4"},
		}

		for _, test := range tests {
			byteArray, err := parseHexColor(test.hex)
			c.So(err, ShouldBeNil)
			c.So(byteArray, ShouldEqual, test.byteArray)
		}
	})
}
