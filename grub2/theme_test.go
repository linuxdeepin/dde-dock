/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package grub2

import (
	C "launchpad.net/gocheck"
)

func (*GrubTester) TestGetplJsonData(c *C.C) {
	testJSONContent := `{
    "BrightScheme":{"ItemColor":"#a6a6a6","SelectedItemColor":"#05abcf", "TerminalBox":"terminal_box_bright_*.png", "MenuPixmapStyle":"menu_box_bright_*.png", "ScrollbarThumb":"scrollbar_bright_*.png"},
    "DarkScheme":{"ItemColor":"#a6a6a6","SelectedItemColor":"#05abcf", "TerminalBox":"terminal_box_dark_*.png", "MenuPixmapStyle":"menu_box_dark_*.png", "ScrollbarThumb":"scrollbar_dark_*.png"},
    "CurrentScheme":{"ItemColor":"#a6a6a6","SelectedItemColor":"#05abcf", "TerminalBox":"terminal_box_dark_*.png", "MenuPixmapStyle":"menu_box_dark_*.png", "ScrollbarThumb":"scrollbar_dark_*.png"}
}`
	wantJSONData := &TplJSONData{
		ThemeScheme{"#a6a6a6", "#05abcf", "terminal_box_bright_*.png", "menu_box_bright_*.png", "scrollbar_bright_*.png"},
		ThemeScheme{"#a6a6a6", "#05abcf", "terminal_box_dark_*.png", "menu_box_dark_*.png", "scrollbar_dark_*.png"},
		ThemeScheme{"#a6a6a6", "#05abcf", "terminal_box_dark_*.png", "menu_box_dark_*.png", "scrollbar_dark_*.png"},
	}

	theme := &Theme{}
	jsonData, err := theme.getTplJSONData([]byte(testJSONContent))
	if err != nil {
		c.Error(err)
	}
	c.Check(*jsonData, C.Equals, *wantJSONData)
}

func (*GrubTester) TestGetCustomizedThemeContent(c *C.C) {
	testThemeTplContent := `# GRUB2 gfxmenu Linux Deepin theme
# Designed for 1024x768 resolution
# Global Property
title-text: ""
desktop-image: "background.png"
desktop-color: "#000000"
terminal-box: "{{.TerminalBox}}"
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
  menu_pixmap_style = "{{.MenuPixmapStyle}}"
  scrollbar = true
  scrollbar_width = 7
  scrollbar_thumb = "{{.ScrollbarThumb}}"
}
`
	wantThemeTxtContent := `# GRUB2 gfxmenu Linux Deepin theme
# Designed for 1024x768 resolution
# Global Property
title-text: ""
desktop-image: "background.png"
desktop-color: "#000000"
terminal-box: "terminal_box_bright_*.png"
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
  selected_item_color = "#05abcf"
  item_spacing = 0
  menu_pixmap_style = "menu_box_bright_*.png"
  scrollbar = true
  scrollbar_width = 7
  scrollbar_thumb = "scrollbar_bright_*.png"
}
`
	theme := &Theme{}
	tplValues := ThemeScheme{"#a6a6a6", "#05abcf", "terminal_box_bright_*.png", "menu_box_bright_*.png", "scrollbar_bright_*.png"}

	s, _ := theme.getCustomizedThemeContent([]byte(testThemeTplContent), tplValues)
	c.Check(string(s), C.Equals, wantThemeTxtContent)
}
