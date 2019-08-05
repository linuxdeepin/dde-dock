package appearance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hexColorToXsColor(t *testing.T) {
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
		assert.Equal(t, test.xs, xsColor)
		if test.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func Test_xsColorToHexColor(t *testing.T) {
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
		assert.Equal(t, test.hex, hexColor)
		if test.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}
