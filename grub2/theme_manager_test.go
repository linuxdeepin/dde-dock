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

func (tm *ThemeManager) TestGetCustomizedThemeContent(c *C) {
	testThemeTplContent := `# GRUB2 gfxmenu Linux Deepin theme
# Designed for 1024x768 resolution
# Global Property
title-text: ""
desktop-image: "{{.Background}}"
desktop-color: "#000000"
terminal-box: "terminal_*.png"
terminal-font: "Fixed Regular 13"

# Show the boot menu
+ boot_menu {
  left = 15%
  top = 20%
  width = 70%
  height = 60%
  item_font = "Courier 10 Pitch Bold 16"
  selected_item_font = "Courier 10 Pitch Bold 24"
  item_color = "{{.ItemColor}}"
  selected_item_color = "{{.SelectedItemColor}}"
  item_spacing = 0
  menu_pixmap_style = "empty_*.png"
  scrollbar = true
  scrollbar_width = 7
  scrollbar_thumb = "sb_th_*.png"
}
`
	wantThemeTxtContent := `# GRUB2 gfxmenu Linux Deepin theme
# Designed for 1024x768 resolution
# Global Property
title-text: ""
desktop-image: "background.jpg"
desktop-color: "#000000"
terminal-box: "terminal_*.png"
terminal-font: "Fixed Regular 13"

# Show the boot menu
+ boot_menu {
  left = 15%
  top = 20%
  width = 70%
  height = 60%
  item_font = "Courier 10 Pitch Bold 16"
  selected_item_font = "Courier 10 Pitch Bold 24"
  item_color = "#a6a6a6"
  selected_item_color = "#fefefe"
  item_spacing = 0
  menu_pixmap_style = "empty_*.png"
  scrollbar = true
  scrollbar_width = 7
  scrollbar_thumb = "sb_th_*.png"
}
`
	tplData := ThemeTpl{"background.jpg", "#a6a6a6", "#fefefe"}
	s, _ := tm.getCustomizedThemeContent(testThemeTplContent, tplData)
	c.Check(s, Equals, wantThemeTxtContent)
}


















