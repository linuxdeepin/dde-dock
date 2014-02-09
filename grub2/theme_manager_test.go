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
	f := "/boot/grub/themes/test-theme/theme.txt"
	want := f
	tm.setEnabledThemeMainFile(f)
	c.Check(tm.getEnabledThemeMainFile(), Equals, want)

	f = "/dir/grub/themes/test-theme/theme.txt"
	want = ""
	tm.setEnabledThemeMainFile(f)
	c.Check(tm.getEnabledThemeMainFile(), Equals, want)
}

// TODO remove
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

// TODO remove
func (tm *ThemeManager) TestGetThemePath(c *C) {
	var tests = []struct {
		s, want string
	}{
		{"theme-name", "/boot/grub/themes/theme-name"},
		{"", "/boot/grub/themes"}, // TODO
	}
	for _, t := range tests {
		get, _ := tm.getThemePath(t.s)
		c.Check(get, Equals, t.want)
	}
}

func (tm *ThemeManager) TestGetThemeMainFile(c *C) {
	var tests = []struct {
		s, want string
	}{
		{"theme-name", "/boot/grub/themes/theme-name/theme.txt"},
		{"", "/boot/grub/themes/theme.txt"}, // TODO
	}
	for _, t := range tests {
		get, _ := tm.getThemeMainFile(t.s)
		c.Check(get, Equals, t.want)
	}
}

func (tm *ThemeManager) TestGetThemeTplFile(c *C) {
	var tests = []struct {
		s, want string
	}{
		{"theme-name", "/boot/grub/themes/theme-name/theme.tpl"},
		{"", "/boot/grub/themes/theme.tpl"}, // TODO
	}
	for _, t := range tests {
		get, _ := tm.getThemeTplFile(t.s)
		c.Check(get, Equals, t.want)
	}
}
