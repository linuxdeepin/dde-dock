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
	SetDefaultGrubConfigFile(tmpConfigFile)
	SetDefaultThemePath(tmpThemeDir)
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
	SetDefaultGrubConfigFile(tmpConfigFile)
	SetDefaultThemePath(tmpThemeDir)
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
