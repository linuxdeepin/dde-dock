package item

import (
	"fmt"
	C "launchpad.net/gocheck"
	"os"
	"path"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"gir/gio-2.0"
	"strings"
	"testing"
)

func TestItem(t *testing.T) {
	C.TestingT(t)
}

type ItemTestSuite struct {
	firefox         *gio.DesktopAppInfo
	notExistDesktop *gio.DesktopAppInfo
	testDataDir     string
}

var _ = C.Suite(&ItemTestSuite{})

func (s *ItemTestSuite) SetUpSuite(c *C.C) {
	s.testDataDir = "../testdata/"
	s.notExistDesktop = nil
}

func (s *ItemTestSuite) getItemInfoForDiffLang(name, lang string, c *C.C) ItemInfo {
	oldLangEnv := os.Getenv("LANGUAGE")
	defer os.Setenv("LANGUAGE", oldLangEnv)

	os.Setenv("LANGUAGE", lang)
	desktopPath := path.Join(s.testDataDir, name)
	desktop := gio.NewDesktopAppInfoFromFilename(desktopPath)
	desktop.GetCommandline()
	if desktop == nil {
		c.Skip(fmt.Sprintf("create desktop(%s) failed", desktopPath))
	}
	defer desktop.Unref()

	return New(desktop)
}

func (s *ItemTestSuite) TestNewItem(c *C.C) {
	c.Assert(New(s.notExistDesktop), C.IsNil)

	firefox := gio.NewDesktopAppInfoFromFilename(path.Join(s.testDataDir, "firefox.desktop"))
	c.Assert(firefox, C.NotNil)
}

func (s *ItemTestSuite) TestLocaleName(c *C.C) {
	var item ItemInfo
	item = s.getItemInfoForDiffLang("firefox.desktop", "en_US", c)
	c.Assert(item.LocaleName(), C.Equals, "Firefox Web Browser")

	item = s.getItemInfoForDiffLang("firefox.desktop", "zh_CN", c)
	c.Assert(item.LocaleName(), C.Equals, "Firefox 网络浏览器")
}

func (s *ItemTestSuite) TestID(c *C.C) {
	var item ItemInfo

	item = s.getItemInfoForDiffLang("firefox.desktop", "en_US", c)
	c.Assert(item.ID(), C.Equals, ItemID("firefox"))
	item = s.getItemInfoForDiffLang("firefox.desktop", "zh_CN", c)
	c.Assert(item.ID(), C.Equals, ItemID("firefox"))

	item = s.getItemInfoForDiffLang("deepin-music-player.desktop", "en_US", c)
	c.Assert(item.ID(), C.Equals, ItemID("deepin-music-player"))
	item = s.getItemInfoForDiffLang("deepin-music-player.desktop", "zh_CN", c)
	c.Assert(item.ID(), C.Equals, ItemID("deepin-music-player"))

	item = s.getItemInfoForDiffLang("qmmp_cue.desktop", "en_US", c)
	c.Assert(item.ID(), C.Equals, ItemID("qmmp-cue"))
	item = s.getItemInfoForDiffLang("qmmp_cue.desktop", "zh_CN", c)
	c.Assert(item.ID(), C.Equals, ItemID("qmmp-cue"))
}

func (s *ItemTestSuite) TestName(c *C.C) {
	var item ItemInfo
	item = s.getItemInfoForDiffLang("firefox.desktop", "en_US", c)
	c.Assert(item.Name(), C.Equals, "Firefox Web Browser")

	item = s.getItemInfoForDiffLang("firefox.desktop", "zh_CN", c)
	c.Assert(item.Name(), C.Equals, "Firefox Web Browser")
}

func (s *ItemTestSuite) TestKeywords(c *C.C) {
	var item ItemInfo
	item = s.getItemInfoForDiffLang("firefox.desktop", "en_US", c)
	expectedKeywords := strings.Split("Internet;WWW;Browser;Web;Explorer", ";")
	keywords := item.Keywords()

	c.Assert(keywords, C.HasLen, len(expectedKeywords))
	for i := 0; i < len(expectedKeywords); i++ {
		c.Assert(keywords[i], C.Equals, strings.ToLower(expectedKeywords[i]))
	}

	item = s.getItemInfoForDiffLang("firefox.desktop", "zh_CN", c)
	expectedKeywords = strings.Split("Internet;WWW;Browser;Web;Explorer;网页;浏览;上网;火狐;Firefox;ff;互联网;网站", ";")
	keywords = item.Keywords()

	c.Assert(keywords, C.HasLen, len(expectedKeywords))
	for i := 0; i < len(expectedKeywords); i++ {
		c.Assert(keywords[i], C.Equals, strings.ToLower(expectedKeywords[i]))
	}
}

func (s *ItemTestSuite) TestDescription(c *C.C) {
	var item ItemInfo
	item = s.getItemInfoForDiffLang("firefox.desktop", "en_US", c)
	c.Assert(item.Description(), C.Equals, "Browse the World Wide Web")

	item = s.getItemInfoForDiffLang("firefox.desktop", "zh_CN", c)
	c.Assert(item.Description(), C.Equals, "浏览互联网")
}

func (s *ItemTestSuite) TestGenericName(c *C.C) {
	var item ItemInfo
	item = s.getItemInfoForDiffLang("firefox.desktop", "en_US", c)
	c.Assert(item.GenericName(), C.Equals, "Web Browser")

	item = s.getItemInfoForDiffLang("firefox.desktop", "zh_CN", c)
	c.Assert(item.GenericName(), C.Equals, "网络浏览器")
}

func (s *ItemTestSuite) TestPath(c *C.C) {
	var item ItemInfo
	item = s.getItemInfoForDiffLang("firefox.desktop", "en_US", c)
	c.Assert(item.Path(), C.Equals, path.Join(s.testDataDir, "firefox.desktop"))
}

func (s *ItemTestSuite) TestIcon(c *C.C) {
	var item ItemInfo
	item = s.getItemInfoForDiffLang("firefox.desktop", "en_US", c)
	c.Assert(item.Icon(), C.Equals, "firefox")
}

func (s *ItemTestSuite) TestExecCmd(c *C.C) {
	var item ItemInfo
	item = s.getItemInfoForDiffLang("firefox.desktop", "en_US", c)
	c.Assert(item.ExecCmd(), C.Equals, "ls firefox %u")

	item = s.getItemInfoForDiffLang("firefox.desktop", "zh_CN", c)
	c.Assert(item.ExecCmd(), C.Equals, "ls firefox %u")
}
