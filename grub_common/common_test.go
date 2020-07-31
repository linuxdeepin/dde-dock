package grub_common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGfxmodesSortDesc(t *testing.T) {
	gfxmodes := Gfxmodes{
		{1366, 768},
		{800, 600},
		{1360, 768},
		{1024, 768},
		{960, 720},
		{1024, 576},
	}
	gfxmodes.SortDesc()

	assert.Equal(t, Gfxmodes{
		{1366, 768},
		{1360, 768},
		{1024, 768},
		{960, 720},
		{1024, 576},
		{800, 600},
	}, gfxmodes)
}

func TestGfxmodesMax(t *testing.T) {
	gfxmodes := Gfxmodes{
		{800, 600},
		{1360, 768},
		{1024, 768},
		{1366, 768},
		{960, 720},
		{1024, 576},
	}
	assert.Equal(t, Gfxmode{1366, 768}, gfxmodes.Max())
	assert.Equal(t, Gfxmode{}, Gfxmodes(nil).Max())
}

func TestParseGfxmode(t *testing.T) {
	mode, err := ParseGfxmode("1024x768")
	assert.Nil(t, err)
	assert.Equal(t, Gfxmode{1024, 768}, mode)

	_, err = ParseGfxmode("auto")
	assert.NotNil(t, err)

	_, err = ParseGfxmode("1024x768x32")
	assert.NotNil(t, err)
}

func Test_parseBootArgDeepinGfxmode(t *testing.T) {
	cur, all, err := parseBootArgDeepinGfxmode("1,1280x1024,1366x768,1024x768")
	assert.Nil(t, err)
	assert.Equal(t, cur, Gfxmode{1366, 768})
	assert.Equal(t, all, Gfxmodes{
		{1280, 1024},
		{1366, 768},
		{1024, 768},
	})

	cur, all, err = parseBootArgDeepinGfxmode("0,1280x1024")
	assert.Nil(t, err)
	assert.Equal(t, cur, Gfxmode{1280, 1024})
	assert.Equal(t, all, Gfxmodes{
		{1280, 1024},
	})

	_, _, err = parseBootArgDeepinGfxmode("")
	assert.NotNil(t, err)

	_, _, err = parseBootArgDeepinGfxmode("3,1280x1024,1366x768,1024x768")
	assert.NotNil(t, err)

	_, _, err = parseBootArgDeepinGfxmode("-1,1280x1024,1366x768,1024x768")
	assert.NotNil(t, err)

	_, _, err = parseBootArgDeepinGfxmode("1,1280x1024,1366x768,1024x768,auto")
	assert.NotNil(t, err)
}
