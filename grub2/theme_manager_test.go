package main

import (
	. "launchpad.net/gocheck"
)

var tm *ThemeManager

func init() {
	tm = NewThemeManager()
	Suite(tm)
}

func (tm *ThemeManager) TestGetterAndSetter(c *C) {
	f := "/dir/grub/themes/test-theme/theme.txt"
	want := ""
	tm.setEnabledThemeMainFile(f)
	c.Check(tm.getEnabledThemeMainFile(), Equals, want)
	
	f = "/boot/grub/themes/test-theme/theme.txt"
	want = f
	tm.setEnabledThemeMainFile(f)
	c.Check(tm.getEnabledThemeMainFile(), Equals, want)
}

func (tm *ThemeManager) TestGetThemeName(c *C) {
	var tests = []struct {
		s, want string
	}{
		{"/boot/grub/themes/name/theme.txt", "name"},
		{"/dir1/name/test.txt", "name"},
		{"", ""},
	}
	for _, t := range tests {
		c.Check(tm.getThemeName(t.s), Equals, t.want)
	}
}


















