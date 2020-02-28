package launcher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getXCategory(t *testing.T) {
	assert.Equal(t, CategoryOthers, getXCategory(nil))
	assert.Equal(t, CategoryMusic, getXCategory([]string{"audio"}))
	assert.Equal(t, CategoryVideo, getXCategory([]string{"video"}))
	assert.Equal(t, CategoryVideo, getXCategory([]string{"audiovideo", "player", "recorder"}))
}
