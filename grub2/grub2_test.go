package main

import (
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

var (
	testMenuContent = `
menuentry 'LinuxDeepin GNU/Linux' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-simple' {
recordfail
}
submenu 'Advanced options for LinuxDeepin GNU/Linux' $menuentry_id_option 'gnulinux-advanced' {
	menuentry 'LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-3.11.0-15-generic-advanced' {
	recordfail
		echo	'载入 Linux 3.11.0-15-generic ...'
	}
`
	testMenuContentLong = `
menuentry 'LinuxDeepin GNU/Linux' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-simple' {
recordfail
	load_video
}
submenu 'Advanced options for LinuxDeepin GNU/Linux' $menuentry_id_option 'gnulinux-advanced' {
	menuentry 'LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-3.11.0-15-generic-advanced' {
	recordfail
		echo	'载入 Linux 3.11.0-15-generic ...'
	}
    submenu 'Inner submenu for test' {
    	menuentry 'Menuentry in Level 3' {
        }
    }
}
menuentry 'Other OS' {
}
`
	testConfigContent = `
# comment line
GRUB_DEFAULT="0"
GRUB_HIDDEN_TIMEOUT="0"
GRUB_HIDDEN_TIMEOUT_QUIET="true"
# comment line
GRUB_TIMEOUT="10"
GRUB_GFXMODE="1024x768"
  GRUB_BACKGROUND=/boot/grub/background.png
  GRUB_THEME="/boot/grub/themes/demo/theme.txt"
`
)

var grub *Grub2

func init() {
	grub = &Grub2{}
	Suite(grub)
}

func (grub *Grub2) TestParseTitle(c *C) {
	var tests = []struct {
		s, want string
	}{
		{`menuentry 'LinuxDeepin GNU/Linux' --class linux $menuentry_id_option 'gnulinux-simple'`, `LinuxDeepin GNU/Linux`},
		{`  menuentry 'LinuxDeepin GNU/Linux' --class linux`, `LinuxDeepin GNU/Linux`},
		{`submenu 'Advanced options for LinuxDeepin GNU/Linux'`, `Advanced options for LinuxDeepin GNU/Linux`},
		{``, ``},
	}
	for _, t := range tests {
		got, _ := grub.parseTitle(t.s)
		c.Check(got, Equals, t.want)
	}
}

func (grub *Grub2) TestParseEntries(c *C) {
	wantEntyTitles := []string{
		`LinuxDeepin GNU/Linux`,
		`Advanced options for LinuxDeepin GNU/Linux`,
		`Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic`,
		`Advanced options for LinuxDeepin GNU/Linux>Inner submenu for test`,
		`Advanced options for LinuxDeepin GNU/Linux>Inner submenu for test>Menuentry in Level 3`,
		`Other OS`,
	}

	grub.parseEntries(testMenuContentLong)
	c.Check(len(grub.entries), Equals, len(wantEntyTitles))
	for i, entry := range grub.entries {
		c.Check(entry.getFullTitle(), Equals, wantEntyTitles[i])
	}
}

func (grub *Grub2) TestParseSettings(c *C) {
	grub.parseEntries(testMenuContent)
	grub.parseSettings(testConfigContent)

	wantSettingCount := 7
	wantDefaultEntry := "LinuxDeepin GNU/Linux"
	wantTimeout := "10"
	wantTheme := "/boot/grub/themes/demo/theme.txt"

	c.Check(len(grub.settings), Equals, wantSettingCount)
	c.Check(grub.settings["GRUB_DEFAULT"], Equals, wantDefaultEntry)
	c.Check(grub.settings["GRUB_TIMEOUT"], Equals, wantTimeout)
	c.Check(grub.settings["GRUB_THEME"], Equals, wantTheme)
}

func (grub *Grub2) TestSetterAndGetter(c *C) {
	grub.parseEntries(testMenuContent)
	grub.parseSettings(testConfigContent)

	entryTitles, _ := grub.GetSimpleEntryTitles()
	c.Check(len(entryTitles), Equals, 1)

	// default entry
	wantDefaultEntry := `LinuxDeepin GNU/Linux`
	c.Check(grub.getDefaultEntry(), Equals, wantDefaultEntry)
	grub.setDefaultEntry(`Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic`)
	c.Check(grub.getDefaultEntry(), Equals, wantDefaultEntry)

	// timeout
	wantTimeout := int32(10)
	c.Check(grub.getTimeout(), Equals, wantTimeout)
	wantTimeout = int32(15)
	grub.setTimeout(wantTimeout)
	c.Check(grub.getTimeout(), Equals, wantTimeout)

	// gfxmode
	wantGfxmode := "1024x768"
	c.Check(grub.getGfxmode(), Equals, wantGfxmode)
	wantGfxmode = "saved"
	grub.setGfxmode(wantGfxmode)
	c.Check(grub.getGfxmode(), Equals, wantGfxmode)

	// theme
	wantTheme := "/boot/grub/themes/demo/theme.txt"
	c.Check(grub.getTheme(), Equals, wantTheme)
	wantTheme = "another_theme.txt"
	grub.setTheme(wantTheme)
	c.Check(grub.getTheme(), Equals, wantTheme)
}

func (grub *Grub2) TestSaveDefaultSettings(c *C) {
	testConfigContent := `GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
`
	wantConfigContent := `GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_DEFAULT="LinuxDeepin GNU/Linux"
`
	grub.parseEntries(testMenuContent)
	grub.parseSettings(testConfigContent)
	c.Check(grub.getSettingContentToSave(), Equals, wantConfigContent)
}

func (grub *Grub2) TestSaveSettings(c *C) {
	testConfigContent := `GRUB_DEFAULT="0"
GRUB_TIMEOUT="10"
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_GFXMODE="1024x768"
`
	wantConfigContent := `GRUB_DEFAULT="LinuxDeepin GNU/Linux"
GRUB_TIMEOUT="15"
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_GFXMODE="auto"
GRUB_THEME="/boot/grub/themes/demo/theme.txt"
`

	grub.parseEntries(testMenuContent)
	grub.parseSettings(testConfigContent)

	grub.setDefaultEntry(`LinuxDeepin GNU/Linux`)
	grub.setTimeout(15)
	grub.setGfxmode("auto")
	grub.setTheme("/boot/grub/themes/demo/theme.txt")
	c.Check(grub.getSettingContentToSave(), Equals, wantConfigContent)
}

func (grub *Grub2) TestGetEntryTitles(c *C) {
	wantEntyTitles := []string{
		`LinuxDeepin GNU/Linux`,
		`Other OS`,
	}

	grub.parseEntries(testMenuContentLong)
	entryTitles, _ := grub.GetSimpleEntryTitles()
	c.Check(len(entryTitles), Equals, len(wantEntyTitles))
	for i, title := range entryTitles {
		c.Check(title, Equals, wantEntyTitles[i])
	}
}
