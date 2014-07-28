// Package main provides ...
package dock

import (
	"testing"
)

func TestGenerateMenuJson(t *testing.T) {
	f := NewNormalApp("firefox.desktop")
	f.buildMenu()
	t.Log(f.Menu)
}
