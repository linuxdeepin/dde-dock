package grub2

import (
	"io/ioutil"
	. "launchpad.net/gocheck"
	"pkg.linuxdeepin.com/lib/graphic"
	"pkg.linuxdeepin.com/lib/utils"
)

const (
	testGrubSettingsFile              = "testdata/grub_settings"
	testGrubThemeBackgroundSourceFile = "testdata/grub_theme_background_source"
)

var (
	tmpBaseDir    = "/tmp/dde_daemon_grub2_test"
	tmpConfigFile = tmpBaseDir + "/dde_daemon_grub2_test_settings"
	tmpThemeDir   = tmpBaseDir + "/dde_daemon_grub2_test_theme_dir"
	tmpGfxmode    = "1200x900"
)

func (*GrubTester) TestCustomArguments(c *C) {
	// prepare
	configFile = tmpBaseDir + "/grub2.json"
	SetDefaultGrubSettingFile(tmpConfigFile)
	SetDefaultThemeDir(tmpThemeDir)
	defer func() {
		configFile = ConfigFileDefault
		SetDefaultGrubSettingFile(DefaultGrubSettingFile)
		SetDefaultThemeDir(DefaultThemeDir)
	}()

	c.Check(grubSettingFile, Equals, tmpConfigFile)
	c.Check(themeDir, Equals, tmpThemeDir)
	c.Check(themeMainFile, Equals, tmpThemeDir+"/theme.txt")
	c.Check(themeTplFile, Equals, tmpThemeDir+"/theme.tpl")
	c.Check(themeJSONFile, Equals, tmpThemeDir+"/theme_tpl.json")
	c.Check(themeBgOrigSrcFile, Equals, tmpThemeDir+"/background_origin_source")
	c.Check(themeBgSrcFile, Equals, tmpThemeDir+"/background_source")
	c.Check(themeBgFile, Equals, tmpThemeDir+"/background.png")
	c.Check(themeBgThumbFile, Equals, tmpThemeDir+"/background_thumb.png")

	theme := NewTheme()
	c.Check(theme.themeDir, Equals, themeDir)
	c.Check(theme.mainFile, Equals, themeMainFile)
	c.Check(theme.tplFile, Equals, themeTplFile)
	c.Check(theme.jsonFile, Equals, themeJSONFile)
	c.Check(theme.bgSrcFile, Equals, themeBgSrcFile)
	c.Check(theme.bgFile, Equals, themeBgFile)
	c.Check(theme.bgThumbFile, Equals, themeBgThumbFile)
}

func (*GrubTester) TestSetup(c *C) {
	wantSettingsContent := `GRUB_BACKGROUND="<none>"
GRUB_CMDLINE_LINUX="locale=zh_CN.UTF-8 url=http://cdimage/nfsroot/deepin-2014/desktop/current/amd64/preseed/deepin.seed initrd=http://cdimage/nfsroot/deepin-2014/desktop/current/amd64/casper/initrd.lz"
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_DEFAULT="0"
GRUB_DISTRIBUTOR="` + "`lsb_release -d -s 2> /dev/null || echo Debian`" + `"
GRUB_GFXMODE="1200x900"
GRUB_THEME="/tmp/dde_daemon_grub2_test/dde_daemon_grub2_test_theme_dir/theme.txt"
GRUB_TIMEOUT="5"
`

	// prepare
	configFile = tmpBaseDir + "/grub2.json"
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
	c.Check(string(settingContentBuf), Equals, wantSettingsContent)
	g.readSettings()
	c.Check(g.getSettingDefaultEntry(), Equals, "0")
	c.Check(g.getSettingTimeout(), Equals, int32(5))
	c.Check(g.getSettingGfxmode(), Equals, tmpGfxmode)
	c.Check(g.getSettingTheme(), Equals, tmpThemeDir+"/theme.txt")
	w, h, _ = graphic.GetImageSize(tmpThemeDir + "/background.png")
	c.Check(w, Equals, 1200)
	c.Check(h, Equals, 900)
}

func (*GrubTester) TestSetupGfxmode(c *C) {
	// prepare
	configFile = tmpBaseDir + "/grub2.json"
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
	c.Check(g.getSettingDefaultEntry(), Equals, "0")
	c.Check(g.getSettingTimeout(), Equals, int32(5))
	c.Check(g.getSettingGfxmode(), Equals, getPrimaryScreenBestResolutionStr())
	c.Check(g.getSettingTheme(), Equals, tmpThemeDir+"/theme.txt")
	w, h, _ = graphic.GetImageSize(tmpThemeDir + "/background.png")
	c.Check(w, Equals, sw)
	c.Check(h, Equals, sh)

	// setup with none gfxmode
	utils.CopyFile(testGrubSettingsFile, tmpConfigFile)
	g.Setup("auto")
	g.readSettings()
	c.Check(g.getSettingDefaultEntry(), Equals, "0")
	c.Check(g.getSettingTimeout(), Equals, int32(5))
	c.Check(g.getSettingGfxmode(), Equals, "auto")
	c.Check(g.getSettingTheme(), Equals, tmpThemeDir+"/theme.txt")
	w, h, _ = graphic.GetImageSize(tmpThemeDir + "/background.png")
	c.Check(w, Equals, sw)
	c.Check(h, Equals, sh)

	// setup with wrong gfxmode format
	utils.CopyFile(testGrubSettingsFile, tmpConfigFile)
	g.Setup("1024x")
	g.readSettings()
	c.Check(g.getSettingDefaultEntry(), Equals, "0")
	c.Check(g.getSettingTimeout(), Equals, int32(5))
	c.Check(g.getSettingGfxmode(), Equals, "1024x")
	c.Check(g.getSettingTheme(), Equals, tmpThemeDir+"/theme.txt")
	w, h, _ = graphic.GetImageSize(tmpThemeDir + "/background.png")
	c.Check(w, Equals, sw)
	c.Check(h, Equals, sh)
}

func (*GrubTester) TestDoSetThemeBackgroundSourceFile(c *C) {
	// TODO
	// // prepare
	// configFile = tmpBaseDir + "/grub2.json"
	// utils.EnsureDirExist(tmpThemeDir)
	// utils.CopyFile(testGrubThemeBackgroundSourceFile, tmpThemeDir+"/background_source")
	// SetDefaultGrubSettingFile(tmpConfigFile)
	// SetDefaultThemeDir(tmpThemeDir)
	// 	defer func() {
	// 	configFile = ConfigFileDefault
	// 	SetDefaultGrubSettingFile(DefaultGrubSettingFile)
	// 	SetDefaultThemeDir(DefaultThemeDir)
	// }()

	// g := NewGrub2()
	// var w, h int
	// var sw, sh int
	// tmpsw, tmpsh := getPrimaryScreenBestResolution()
	// sw, sh = int(tmpsw), int(tmpsh)

	// setup := &SetupWrapper{}
	// setup.DoGenerateThemeBackground(w, h)
	// g.SetupTheme()
}
