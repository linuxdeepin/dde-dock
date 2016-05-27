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
	"io/ioutil"
	C "launchpad.net/gocheck"
	"pkg.deepin.io/lib/graphic"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/utils"
)

const (
	testGrubSettingsFile              = "testdata/grub_settings"
	testGrubMenuFile                  = "testdata/grub.cfg"
	testGrubThemeBackgroundSourceFile = "testdata/grub_theme_background_source"
)

var (
	tmpBaseDir      = "./testdata"
	tmpConfigFile   = tmpBaseDir + "/tmp_settings"
	tmpGrubMenuFile = tmpBaseDir + "/tmp_grub.cfg"
	tmpThemeDir     = tmpBaseDir + "/tmp_theme_dir"
	tmpGfxmode      = "1200x900"
)

func (*GrubTester) TestCustomArguments(c *C.C) {
	// prepare
	configFile = tmpBaseDir + "/tmp_grub2.json"
	SetDefaultGrubSettingFile(tmpConfigFile)
	SetDefaultThemeDir(tmpThemeDir)
	defer func() {
		configFile = ConfigFileDefault
		SetDefaultGrubSettingFile(DefaultGrubSettingFile)
		SetDefaultThemeDir(DefaultThemeDir)
	}()

	c.Check(grubSettingFile, C.Equals, tmpConfigFile)
	c.Check(themeDir, C.Equals, tmpThemeDir)
	c.Check(themeMainFile, C.Equals, tmpThemeDir+"/theme.txt")
	c.Check(themeTplFile, C.Equals, tmpThemeDir+"/theme.tpl")
	c.Check(themeJSONFile, C.Equals, tmpThemeDir+"/theme_tpl.json")
	c.Check(themeBgOrigSrcFile, C.Equals, tmpThemeDir+"/background_origin_source")
	c.Check(themeBgSrcFile, C.Equals, tmpThemeDir+"/background_source")
	c.Check(themeBgFile, C.Equals, tmpThemeDir+"/background.png")
	c.Check(themeBgThumbFile, C.Equals, tmpThemeDir+"/background_thumb.png")

	theme := NewTheme()
	c.Check(theme.themeDir, C.Equals, themeDir)
	c.Check(theme.mainFile, C.Equals, themeMainFile)
	c.Check(theme.tplFile, C.Equals, themeTplFile)
	c.Check(theme.jsonFile, C.Equals, themeJSONFile)
	c.Check(theme.bgSrcFile, C.Equals, themeBgSrcFile)
	c.Check(theme.bgFile, C.Equals, themeBgFile)
	c.Check(theme.bgThumbFile, C.Equals, themeBgThumbFile)
}

func (*GrubTester) TestSetup(c *C.C) {
	wantSettingsContent := `GRUB_BACKGROUND="./testdata/tmp_theme_dir/background.png"
GRUB_CMDLINE_LINUX="locale=zh_CN.UTF-8 url=http://cdimage/nfsroot/deepin-2014/desktop/current/amd64/preseed/deepin.seed initrd=http://cdimage/nfsroot/deepin-2014/desktop/current/amd64/casper/initrd.lz"
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_DEFAULT="0"
GRUB_DISTRIBUTOR="` + "`/usr/bin/lsb_release -d -s 2>/dev/null || echo Debian`" + `"
GRUB_GFXMODE="1200x900"
GRUB_THEME="./testdata/tmp_theme_dir/theme.txt"
GRUB_TIMEOUT="5"
`

	// prepare
	configFile = tmpBaseDir + "/tmp_grub2.json"
	utils.EnsureDirExist(tmpThemeDir)
	utils.CopyFile(testGrubThemeBackgroundSourceFile, tmpThemeDir+"/background_source")
	SetDefaultGrubSettingFile(tmpConfigFile)
	SetDefaultThemeDir(tmpThemeDir)
	defer func() {
		configFile = ConfigFileDefault
		SetDefaultGrubSettingFile(DefaultGrubSettingFile)
		SetDefaultThemeDir(DefaultThemeDir)
	}()

	g := NewGrub2()
	var w, h int

	// setup with target gfxmode
	utils.CopyFile(testGrubSettingsFile, tmpConfigFile)
	g.Setup(tmpGfxmode)
	settingContentBuf, _ := ioutil.ReadFile(tmpConfigFile)
	c.Check(string(settingContentBuf), C.Equals, wantSettingsContent)
	g.readSettings()
	c.Check(g.getSettingDefaultEntry(), C.Equals, "0")
	c.Check(g.getSettingTimeout(), C.Equals, int32(5))
	c.Check(g.getSettingGfxmode(), C.Equals, tmpGfxmode)
	c.Check(g.getSettingTheme(), C.Equals, tmpThemeDir+"/theme.txt")
	w, h, _ = graphic.GetImageSize(tmpThemeDir + "/background.png")
	c.Check(w, C.Equals, 1200)
	c.Check(h, C.Equals, 900)
}

func (*GrubTester) TestSetupGfxmode(c *C.C) {
	// prepare
	configFile = tmpBaseDir + "/tmp_grub2.json"
	utils.EnsureDirExist(tmpThemeDir)
	utils.CopyFile(testGrubThemeBackgroundSourceFile, tmpThemeDir+"/background_source")
	SetDefaultGrubSettingFile(tmpConfigFile)
	SetDefaultThemeDir(tmpThemeDir)
	defer func() {
		configFile = ConfigFileDefault
		SetDefaultGrubSettingFile(DefaultGrubSettingFile)
		SetDefaultThemeDir(DefaultThemeDir)
	}()

	g := NewGrub2()
	var w, h int
	var sw, sh int
	tmpsw, tmpsh := getPrimaryScreenBestResolution()
	sw, sh = int(tmpsw), int(tmpsh)

	// setup with none gfxmode
	utils.CopyFile(testGrubSettingsFile, tmpConfigFile)
	g.Setup("")
	g.readSettings()
	c.Check(g.getSettingDefaultEntry(), C.Equals, "0")
	c.Check(g.getSettingTimeout(), C.Equals, int32(5))
	c.Check(g.getSettingGfxmode(), C.Equals, getDefaultGfxmode())
	c.Check(g.getSettingTheme(), C.Equals, tmpThemeDir+"/theme.txt")
	w, h, _ = graphic.GetImageSize(tmpThemeDir + "/background.png")
	c.Check(w, C.Equals, sw)
	c.Check(h, C.Equals, sh)

	// setup with none gfxmode
	utils.CopyFile(testGrubSettingsFile, tmpConfigFile)
	g.Setup("auto")
	g.readSettings()
	c.Check(g.getSettingDefaultEntry(), C.Equals, "0")
	c.Check(g.getSettingTimeout(), C.Equals, int32(5))
	c.Check(g.getSettingGfxmode(), C.Equals, "auto")
	c.Check(g.getSettingTheme(), C.Equals, tmpThemeDir+"/theme.txt")
	w, h, _ = graphic.GetImageSize(tmpThemeDir + "/background.png")
	c.Check(w, C.Equals, sw)
	c.Check(h, C.Equals, sh)

	// setup with wrong gfxmode format

	// will following error message in this case, so we just disable
	// all output and recovery it when done
	oldLogLevel := logger.GetLogLevel()
	logger.SetLogLevel(log.LevelDisable)
	defer logger.SetLogLevel(oldLogLevel)

	utils.CopyFile(testGrubSettingsFile, tmpConfigFile)
	g.Setup("1024x")
	g.readSettings()
	c.Check(g.getSettingDefaultEntry(), C.Equals, "0")
	c.Check(g.getSettingTimeout(), C.Equals, int32(5))
	c.Check(g.getSettingGfxmode(), C.Equals, "1024x")
	c.Check(g.getSettingTheme(), C.Equals, tmpThemeDir+"/theme.txt")
	w, h, _ = graphic.GetImageSize(tmpThemeDir + "/background.png")
	c.Check(w, C.Equals, sw)
	c.Check(h, C.Equals, sh)
}
