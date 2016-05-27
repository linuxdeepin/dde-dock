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
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/utils"
	"testing"
)

func Test(t *testing.T) { C.TestingT(t) }

type GrubTester struct{}

func init() {
	// disable all output but error messages
	logger.SetLogLevel(log.LevelError)
	runWithoutDbus = true
	C.Suite(&GrubTester{})
}

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

func (*GrubTester) TestParseTitle(c *C.C) {
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
		c.Check(got, C.Equals, t.want)
	}
}

func (*GrubTester) TestParseEntries(c *C.C) {
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
	c.Check(len(grub.entries), C.Equals, len(wantEntyTitles))
	for i, entry := range grub.entries {
		c.Check(entry.getFullTitle(), C.Equals, wantEntyTitles[i])
	}
}

func (*GrubTester) TestParseSettings(c *C.C) {
	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContent)
	grub.parseSettings(testGrubSettingsContent)

	wantSettingCount := 7
	wantDefaultEntry := "0"
	wantTimeout := "10"
	wantTheme := "/boot/grub/themes/demo/theme.txt"
	c.Check(len(grub.settings), C.Equals, wantSettingCount)
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, wantDefaultEntry)
	c.Check(grub.doGetSettingTimeout(), C.Equals, wantTimeout)
	c.Check(grub.doGetSettingTheme(), C.Equals, wantTheme)
}

func (*GrubTester) TestParseInvalidSettings(c *C.C) {
	testGrubSettingsContent := `GRUB_DEFUALT=
GRUB_TIMEOUT
GRUB_THEME
`
	grub := NewGrub2()
	grub.parseSettings(testGrubSettingsContent)
	c.Check(len(grub.settings), C.Equals, 1)
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "")
	c.Check(grub.doGetSettingTimeout(), C.Equals, "")
	c.Check(grub.doGetSettingTheme(), C.Equals, "")
	c.Check(grub.getSettingContentToSave(), C.Equals, "")
}

func (*GrubTester) TestSettingDefaultEntry(c *C.C) {
	grub := NewGrub2()
	grub.doFixSettings()

	// default entry
	c.Check(grub.config.DefaultEntry, C.Equals, "0")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "0")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "0")

	// default entry if GRUB_DEFAULT not defined
	grub.doSetSettingDefaultEntry("")
	c.Check(grub.config.DefaultEntry, C.Equals, "")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "0")

	// custom entry index value
	grub.doSetSettingDefaultEntry("3")
	c.Check(grub.config.DefaultEntry, C.Equals, "3")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "3")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "3")

	// custom entry title value
	grub.doSetSettingDefaultEntry("LinuxDeepin GNU/Linux")
	c.Check(grub.config.DefaultEntry, C.Equals, "LinuxDeepin GNU/Linux")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "LinuxDeepin GNU/Linux")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "LinuxDeepin GNU/Linux")

	// custom sub entry value
	grub.doSetSettingDefaultEntry("3>1")
	c.Check(grub.config.DefaultEntry, C.Equals, "3>1")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "3>1")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "3")

	// load entry titles
	grub.parseEntries(testGrubMenuContent)

	// get default entry after titles loaded
	grub.doSetSettingDefaultEntry("0")
	c.Check(grub.config.DefaultEntry, C.Equals, "0")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "0")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("1")
	c.Check(grub.config.DefaultEntry, C.Equals, "1")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "1")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "Advanced options for LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("1>0") // sub entry
	c.Check(grub.config.DefaultEntry, C.Equals, "1>0")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "1>0")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "Advanced options for LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("1>3") // entry title not exists
	c.Check(grub.config.DefaultEntry, C.Equals, "1>3")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "1>3")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "Advanced options for LinuxDeepin GNU/Linux")

	// custom entry title value after titles loaded
	grub.doSetSettingDefaultEntry("LinuxDeepin GNU/Linux")
	c.Check(grub.config.DefaultEntry, C.Equals, "LinuxDeepin GNU/Linux")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "LinuxDeepin GNU/Linux")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("Advanced options for LinuxDeepin GNU/Linux")
	c.Check(grub.config.DefaultEntry, C.Equals, "Advanced options for LinuxDeepin GNU/Linux")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "Advanced options for LinuxDeepin GNU/Linux")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "Advanced options for LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic")
	c.Check(grub.config.DefaultEntry, C.Equals, "Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "Advanced options for LinuxDeepin GNU/Linux>LinuxDeepin GNU/Linux，Linux 3.11.0-15-generic")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "Advanced options for LinuxDeepin GNU/Linux")

	grub.doSetSettingDefaultEntry("Advanced options for LinuxDeepin GNU/Linux>Child Title That Not Exists") // sub entry title not exists
	c.Check(grub.config.DefaultEntry, C.Equals, "Advanced options for LinuxDeepin GNU/Linux>Child Title That Not Exists")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "Advanced options for LinuxDeepin GNU/Linux>Child Title That Not Exists")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "Advanced options for LinuxDeepin GNU/Linux")
}

func (*GrubTester) TestSettingTimeout(c *C.C) {
	grub := NewGrub2()
	grub.doFixSettings()

	// default timeout
	c.Check(grub.config.Timeout, C.Equals, "5")
	c.Check(grub.doGetSettingTimeout(), C.Equals, "5")
	c.Check(grub.getSettingTimeout(), C.Equals, int32(5))

	// default timeout if GRUB_TIMEOUT not defined
	grub.doSetSettingTimeout("")
	c.Check(grub.config.Timeout, C.Equals, "")
	c.Check(grub.doGetSettingTimeout(), C.Equals, "")
	c.Check(grub.getSettingTimeout(), C.Equals, int32(5))

	// custom timeout
	grub.doSetSettingTimeoutLogic(10)
	c.Check(grub.config.Timeout, C.Equals, "10")
	c.Check(grub.doGetSettingTimeout(), C.Equals, "10")
	c.Check(grub.getSettingTimeout(), C.Equals, int32(10))
}

func (*GrubTester) TestFixSettingDefaultEntry(c *C.C) {
	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContent)
	grub.parseSettings(testGrubSettingsContent)

	var needUpdate bool
	needUpdate = grub.doFixSettings()
	c.Check(needUpdate, C.Equals, true)

	c.Check(grub.config.DefaultEntry, C.Equals, "0")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "0")
}

func (*GrubTester) TestFixSettings(c *C.C) {
	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContent)
	grub.parseSettings(testGrubSettingsContent)

	var needUpdate bool
	needUpdate = grub.doFixSettings()
	c.Check(needUpdate, C.Equals, true)

	needUpdate = grub.doFixSettingDistro()
	c.Check(needUpdate, C.Equals, true)

	wantSettingCount := 8
	wantDefaultEntry := "0"
	wantDistro := "`/usr/bin/lsb_release -d -s 2>/dev/null || echo Debian`"
	wantTimeout := "5"
	wantTheme := "/boot/grub/themes/deepin/theme.txt"
	c.Check(len(grub.settings), C.Equals, wantSettingCount)
	c.Check(grub.doGetSettingDistributor(), C.Equals, wantDistro)
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, wantDefaultEntry)
	c.Check(grub.doGetSettingTimeout(), C.Equals, wantTimeout)
	c.Check(grub.doGetSettingTheme(), C.Equals, wantTheme)

	needUpdate = grub.doFixSettings()
	c.Check(needUpdate, C.Equals, false)

	needUpdate = grub.doFixSettingDistro()
	c.Check(needUpdate, C.Equals, false)
}

func (*GrubTester) TestSettingsGeneral(c *C.C) {
	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContent)
	grub.parseSettings(testGrubSettingsContent)

	entryTitles, _ := grub.GetSimpleEntryTitles()
	c.Check(len(entryTitles), C.Equals, 1)

	// gfxmode
	wantGfxmode := "1024x768"
	c.Check(grub.getSettingGfxmode(), C.Equals, wantGfxmode)
	wantGfxmode = "saved"
	grub.doSetSettingGfxmode(wantGfxmode)
	c.Check(grub.getSettingGfxmode(), C.Equals, wantGfxmode)

	// theme
	wantTheme := "/boot/grub/themes/demo/theme.txt"
	c.Check(grub.getSettingTheme(), C.Equals, wantTheme)
	wantTheme = "another_theme.txt"
	grub.doSetSettingTheme(wantTheme)
	c.Check(grub.getSettingTheme(), C.Equals, wantTheme)
}

func (*GrubTester) TestSaveDefaultSettings(c *C.C) {
	testGrubSettingsContent := `GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
`
	wantConfigContent := `GRUB_BACKGROUND="/boot/grub/themes/deepin/background.png"
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
	c.Check(grub.getSettingContentToSave(), C.Equals, wantConfigContent)
}

func (*GrubTester) TestSaveSettings(c *C.C) {
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
	c.Check(grub.getSettingContentToSave(), C.Equals, wantConfigContent)
}

func (*GrubTester) TestGetSimpleEntryTitles(c *C.C) {
	wantEntyTitles := []string{
		`LinuxDeepin GNU/Linux`,
		`Other OS`,
	}

	grub := NewGrub2()
	grub.parseEntries(testGrubMenuContentLong)
	entryTitles, _ := grub.GetSimpleEntryTitles()
	c.Check(len(entryTitles), C.Equals, len(wantEntyTitles))
	for i, title := range entryTitles {
		c.Check(title, C.Equals, wantEntyTitles[i])
	}
}

func (*GrubTester) TestRealSettings(c *C.C) {
	// prepare for real environment
	utils.CopyFile(testGrubSettingsFile, tmpConfigFile)
	utils.CopyFile(testGrubMenuFile, tmpGrubMenuFile)
	configFile = tmpBaseDir + "/tmp_grub2.json"
	grubMenuFile = tmpGrubMenuFile
	SetDefaultGrubSettingFile(tmpConfigFile)
	SetDefaultThemeDir(tmpThemeDir)
	defer func() {
		configFile = ConfigFileDefault
		grubSettingFile = DefaultGrubSettingFile
		SetDefaultGrubSettingFile(DefaultGrubSettingFile)
		SetDefaultThemeDir(DefaultThemeDir)
	}()

	c.Check(grubSettingFile, C.Equals, tmpConfigFile)
	c.Check(themeDir, C.Equals, tmpThemeDir)

	grub := NewGrub2()
	grub.readEntries()
	grub.readSettings()
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "0")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "Deepin 15.2 GNU/Linux")
	grub.setSettingDefaultEntry("Windows 10 (loader) (on /dev/sda1)")
	c.Check(grub.doGetSettingDefaultEntry(), C.Equals, "2")
	c.Check(grub.getSettingDefaultEntry(), C.Equals, "Windows 10 (loader) (on /dev/sda1)")
}
