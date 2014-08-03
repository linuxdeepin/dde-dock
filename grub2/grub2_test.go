package grub2

import (
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type GrubTester struct{}

var _ = Suite(&GrubTester{})

const (
	testGrubMenuContent = `
menuentry 'LinuxDeepin GNU/Linux' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-simple' {
recordfail
}
submenu 'Advanced options for LinuxDeepin GNU/Linux' $menuentry_id_option 'gnulinux-advanced' {
	menuentry 'LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic' --class linuxdeepin --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-3.11.0-15-generic-advanced' {
	recordfail
		echo	'载入 Linux 3.11.0-15-generic ...'
	}
`
	testGrubMenuContentLong = `
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
	testGrubSettingsContent = `
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

func (*GrubTester) TestParseTitle(c *C) {
	var tests = []struct {
		s, want string
	}{
		{`menuentry 'LinuxDeepin GNU/Linux' --class linux $menuentry_id_option 'gnulinux-simple'`, `LinuxDeepin GNU/Linux`},
		{`  menuentry 'LinuxDeepin GNU/Linux' --class linux`, `LinuxDeepin GNU/Linux`},
		{`submenu 'Advanced options for LinuxDeepin GNU/Linux'`, `Advanced options for LinuxDeepin GNU/Linux`},
		{``, ``},
	}
	grub := NewGrub2()
	for _, t := range tests {
		got, _ := grub.parseTitle(t.s)
		c.Check(got, Equals, t.want)
	}
}

func (*GrubTester) TestParseEntries(c *C) {
	wantEntyTitles := []string{
		`LinuxDeepin GNU/Linux`,
		`Advanced options for LinuxDeepin GNU/Linux`,
		`Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic`,
		`Advanced options for LinuxDeepin GNU/Linux>Inner submenu for test`,
		`Advanced options for LinuxDeepin GNU/Linux>Inner submenu for test>Menuentry in Level 3`,
		`Other OS`,
	}

	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContentLong)
	c.Check(len(grub.entries), Equals, len(wantEntyTitles))
	for i, entry := range grub.entries {
		c.Check(entry.getFullTitle(), Equals, wantEntyTitles[i])
	}
}

func (*GrubTester) TestParseSettings(c *C) {
	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContent)
	grub.parseSettings(testGrubSettingsContent)

	wantSettingCount := 7
	wantDefaultEntry := "0"
	wantTimeout := "10"
	wantTheme := "/boot/grub/themes/demo/theme.txt"
	c.Check(len(grub.settings), Equals, wantSettingCount)
	c.Check(grub.doGetSettingDefaultEntry(), Equals, wantDefaultEntry)
	c.Check(grub.doGetSettingTimeout(), Equals, wantTimeout)
	c.Check(grub.doGetSettingTheme(), Equals, wantTheme)
}

func (*GrubTester) TestParseInvalidSettings(c *C) {
	testGrubSettingsContent := `GRUB_DEFUALT=
GRUB_TIMEOUT
GRUB_THEME
`
	grub := NewGrub2()
	grub.parseSettings(testGrubSettingsContent)
	c.Check(len(grub.settings), Equals, 1)
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "")
	c.Check(grub.doGetSettingTimeout(), Equals, "")
	c.Check(grub.doGetSettingTheme(), Equals, "")
	c.Check(grub.getSettingContentToSave(), Equals, "")
}

func (*GrubTester) TestSettingDefaultEntry(c *C) {
	grub := NewGrub2()
	grub.doFixSettings()

	// default entry
	c.Check(grub.config.DefaultEntry, Equals, "0")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "0")
	c.Check(grub.getSettingDefaultEntry(), Equals, "0")

	// default entry if GRUB_DEFAULT not defined
	grub.doSetSettingDefaultEntry("")
	c.Check(grub.config.DefaultEntry, Equals, "")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "")
	c.Check(grub.getSettingDefaultEntry(), Equals, "0")

	// custom entry index value
	grub.doSetSettingDefaultEntry("3")
	c.Check(grub.config.DefaultEntry, Equals, "3")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "3")
	c.Check(grub.getSettingDefaultEntry(), Equals, "3")

	// custom entry title value
	grub.doSetSettingDefaultEntry("LinuxDeepin GNU/Linux")
	c.Check(grub.config.DefaultEntry, Equals, "LinuxDeepin GNU/Linux")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "LinuxDeepin GNU/Linux")
	c.Check(grub.getSettingDefaultEntry(), Equals, "LinuxDeepin GNU/Linux")

	// custom sub entry value
	grub.doSetSettingDefaultEntry("3>1")
	c.Check(grub.config.DefaultEntry, Equals, "3>1")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "3>1")
	c.Check(grub.getSettingDefaultEntry(), Equals, "3")

	// load entry titles
	grub.parseEntries(testGrubMenuContent)

	// get default entry after titles loaded
	grub.doSetSettingDefaultEntry("0")
	c.Check(grub.config.DefaultEntry, Equals, "0")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "0")
	c.Check(grub.getSettingDefaultEntry(), Equals, "LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("1")
	c.Check(grub.config.DefaultEntry, Equals, "1")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "1")
	c.Check(grub.getSettingDefaultEntry(), Equals, "Advanced options for LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("1>0") // sub entry
	c.Check(grub.config.DefaultEntry, Equals, "1>0")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "1>0")
	c.Check(grub.getSettingDefaultEntry(), Equals, "Advanced options for LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("1>3") // entry title not exists
	c.Check(grub.config.DefaultEntry, Equals, "1>3")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "1>3")
	c.Check(grub.getSettingDefaultEntry(), Equals, "Advanced options for LinuxDeepin GNU/Linux")

	// custom entry title value after titles loaded
	grub.doSetSettingDefaultEntry("LinuxDeepin GNU/Linux")
	c.Check(grub.config.DefaultEntry, Equals, "LinuxDeepin GNU/Linux")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "LinuxDeepin GNU/Linux")
	c.Check(grub.getSettingDefaultEntry(), Equals, "LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("Advanced options for LinuxDeepin GNU/Linux")
	c.Check(grub.config.DefaultEntry, Equals, "Advanced options for LinuxDeepin GNU/Linux")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "Advanced options for LinuxDeepin GNU/Linux")
	c.Check(grub.getSettingDefaultEntry(), Equals, "Advanced options for LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic")
	c.Check(grub.config.DefaultEntry, Equals, "Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic")
	c.Check(grub.getSettingDefaultEntry(), Equals, "Advanced options for LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("Advanced options for LinuxDeepin GNU/Linux>Child Title That Not Exists") // sub entry title not exists
	c.Check(grub.config.DefaultEntry, Equals, "Advanced options for LinuxDeepin GNU/Linux>Child Title That Not Exists")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "Advanced options for LinuxDeepin GNU/Linux>Child Title That Not Exists")
	c.Check(grub.getSettingDefaultEntry(), Equals, "Advanced options for LinuxDeepin GNU/Linux")
}

func (*GrubTester) TestSettingTimeout(c *C) {
	grub := NewGrub2()
	grub.doFixSettings()

	// default timeout
	c.Check(grub.config.Timeout, Equals, "5")
	c.Check(grub.doGetSettingTimeout(), Equals, "5")
	c.Check(grub.getSettingTimeout(), Equals, int32(5))

	// default timeout if GRUB_TIMEOUT not defined
	grub.doSetSettingTimeout("")
	c.Check(grub.config.Timeout, Equals, "")
	c.Check(grub.doGetSettingTimeout(), Equals, "")
	c.Check(grub.getSettingTimeout(), Equals, int32(5))

	// custom timeout
	grub.doSetSettingTimeoutLogic(10)
	c.Check(grub.config.Timeout, Equals, "10")
	c.Check(grub.doGetSettingTimeout(), Equals, "10")
	c.Check(grub.getSettingTimeout(), Equals, int32(10))
}

func (*GrubTester) TestFixSettingDefaultEntry(c *C) {
	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContent)
	grub.parseSettings(testGrubSettingsContent)

	var needUpdate bool
	needUpdate = grub.doFixSettings()
	c.Check(needUpdate, Equals, true)

	c.Check(grub.config.DefaultEntry, Equals, "0")
	c.Check(grub.doGetSettingDefaultEntry(), Equals, "0")
}

func (*GrubTester) TestFixSettings(c *C) {
	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContent)
	grub.parseSettings(testGrubSettingsContent)

	var needUpdate bool
	needUpdate = grub.doFixSettings()
	c.Check(needUpdate, Equals, true)

	needUpdate = grub.doFixSettingDistro()
	c.Check(needUpdate, Equals, true)

	wantSettingCount := 8
	wantDefaultEntry := "0"
	wantDistro := "`lsb_release -d -s 2> /dev/null || echo Debian`"
	wantTimeout := "5"
	wantTheme := "/boot/grub/themes/deepin/theme.txt"
	c.Check(len(grub.settings), Equals, wantSettingCount)
	c.Check(grub.doGetSettingDistributor(), Equals, wantDistro)
	c.Check(grub.doGetSettingDefaultEntry(), Equals, wantDefaultEntry)
	c.Check(grub.doGetSettingTimeout(), Equals, wantTimeout)
	c.Check(grub.doGetSettingTheme(), Equals, wantTheme)

	needUpdate = grub.doFixSettings()
	c.Check(needUpdate, Equals, false)

	needUpdate = grub.doFixSettingDistro()
	c.Check(needUpdate, Equals, false)
}

func (*GrubTester) TestSettingsGeneral(c *C) {
	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContent)
	grub.parseSettings(testGrubSettingsContent)

	entryTitles, _ := grub.GetSimpleEntryTitles()
	c.Check(len(entryTitles), Equals, 1)

	// gfxmode
	wantGfxmode := "1024x768"
	c.Check(grub.getSettingGfxmode(), Equals, wantGfxmode)
	wantGfxmode = "saved"
	grub.doSetSettingGfxmode(wantGfxmode)
	c.Check(grub.getSettingGfxmode(), Equals, wantGfxmode)

	// theme
	wantTheme := "/boot/grub/themes/demo/theme.txt"
	c.Check(grub.getSettingTheme(), Equals, wantTheme)
	wantTheme = "another_theme.txt"
	grub.doSetSettingTheme(wantTheme)
	c.Check(grub.getSettingTheme(), Equals, wantTheme)
}

func (*GrubTester) TestSaveDefaultSettings(c *C) {
	testGrubSettingsContent := `GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
`
	wantConfigContent := `GRUB_BACKGROUND="<none>"
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_DEFAULT="0"
GRUB_GFXMODE="1024x768"
GRUB_THEME="/boot/grub/themes/deepin/theme.txt"
GRUB_TIMEOUT="5"
`
	grub := NewGrub2()
	grub.config.Resolution = "1024x768"
	grub.parseSettings(testGrubSettingsContent)
	grub.doFixSettings()
	c.Check(grub.getSettingContentToSave(), Equals, wantConfigContent)
}

func (*GrubTester) TestSaveSettings(c *C) {
	testGrubSettingsContent := `GRUB_DEFAULT="0"
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
	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContent)
	grub.parseSettings(testGrubSettingsContent)

	grub.doSetSettingDefaultEntry(`LinuxDeepin GNU/Linux`)
	grub.doSetSettingTimeoutLogic(15)
	grub.doSetSettingGfxmode("auto")
	grub.doSetSettingTheme("/boot/grub/themes/deepin/theme.txt")
	c.Check(grub.getSettingContentToSave(), Equals, wantConfigContent)
}

func (*GrubTester) TestGetEntryTitles(c *C) {
	wantEntyTitles := []string{
		`LinuxDeepin GNU/Linux`,
		`Other OS`,
	}

	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContentLong)
	entryTitles, _ := grub.GetSimpleEntryTitles()
	c.Check(len(entryTitles), Equals, len(wantEntyTitles))
	for i, title := range entryTitles {
		c.Check(title, Equals, wantEntyTitles[i])
	}
}
