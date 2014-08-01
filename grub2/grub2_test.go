package grub2

import (
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	// grub = &Grub2{}
	grub = NewGrub2()
	grub.config.Resolution = "1024x768"
	Suite(grub)
}

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
	wantDefaultEntry := "0"
	wantTimeout := "10"
	wantTheme := "/boot/grub/themes/demo/theme.txt"
	c.Check(len(grub.settings), Equals, wantSettingCount)
	c.Check(grub.settings["GRUB_DEFAULT"], Equals, wantDefaultEntry)
	c.Check(grub.settings["GRUB_TIMEOUT"], Equals, wantTimeout)
	c.Check(grub.settings["GRUB_THEME"], Equals, wantTheme)

	grub.fixSettings()
	// TODO
	// wantDistro := "`lsb_release -d -s 2> /dev/null || echo Debian`"
	// wantDefaultEntry = "LinuxDeepin GNU/Linux"
	wantTimeout = "10"
	wantTheme = "/boot/grub/themes/deepin/theme.txt"
	c.Check(len(grub.settings), Equals, wantSettingCount)
	// c.Check(grub.settings["GRUB_DISTRIBUTOR"], Equals, wantDistro)
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
	c.Check(grub.getSettingDefaultEntry(), Equals, wantDefaultEntry)
	grub.setSettingDefaultEntry(`Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic`)
	c.Check(grub.getSettingDefaultEntry(), Equals, wantDefaultEntry)

	// timeout
	wantTimeout := int32(10)
	c.Check(grub.getSettingTimeout(), Equals, wantTimeout)
	wantTimeout = int32(15)
	grub.setSettingTimeout(wantTimeout)
	c.Check(grub.getSettingTimeout(), Equals, wantTimeout)

	// gfxmode
	wantGfxmode := "1024x768"
	c.Check(grub.getSettingGfxmode(), Equals, wantGfxmode)
	wantGfxmode = "saved"
	grub.setSettingGfxmode(wantGfxmode)
	c.Check(grub.getSettingGfxmode(), Equals, wantGfxmode)

	// theme
	wantTheme := "/boot/grub/themes/demo/theme.txt"
	c.Check(grub.getSettingTheme(), Equals, wantTheme)
	wantTheme = "another_theme.txt"
	grub.setSettingTheme(wantTheme)
	c.Check(grub.getSettingTheme(), Equals, wantTheme)
}

func (grub *Grub2) TestSaveDefaultSettings(c *C) {
	testConfigContent := `GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
`
	wantConfigContent := `GRUB_BACKGROUND="<none>"
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_DEFAULT="0"
GRUB_GFXMODE="1024x768"
GRUB_THEME="/boot/grub/themes/deepin/theme.txt"
GRUB_TIMEOUT="10"
`
	// TODO
	// 	wantConfigContent := `GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
	// GRUB_DEFAULT="LinuxDeepin GNU/Linux"
	// GRUB_DISTRIBUTOR="` + "`" + `lsb_release -d -s 2> /dev/null || echo Debian` + "`" + `"
	// GRUB_THEME="/boot/grub/themes/deepin/theme.txt"
	// `
	grub.parseEntries(testMenuContent)
	grub.parseSettings(testConfigContent)
	grub.fixSettings()
	c.Check(grub.getSettingContentToSave(), Equals, wantConfigContent)
}

func (grub *Grub2) TestSaveSettings(c *C) {
	testConfigContent := `GRUB_DEFAULT="0"
GRUB_TIMEOUT="10"
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_GFXMODE="1024x768"
`
	wantConfigContent := `GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_DEFAULT="LinuxDeepin GNU/Linux"
GRUB_GFXMODE="auto"
GRUB_THEME="/boot/grub/themes/deepin/theme.txt"
GRUB_TIMEOUT="15"
`
	grub.parseEntries(testMenuContent)
	grub.parseSettings(testConfigContent)

	grub.setSettingDefaultEntry(`LinuxDeepin GNU/Linux`)
	grub.setSettingTimeout(15)
	grub.setSettingGfxmode("auto")
	grub.setSettingTheme("/boot/grub/themes/deepin/theme.txt")
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
