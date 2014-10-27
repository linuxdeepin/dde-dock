package launcher

import (
	C "launchpad.net/gocheck"
	"path/filepath"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	// "sort"
	"testing"
)

func TestCategory(t *testing.T) {
	C.TestingT(t)
}

type CategoryTestSuite struct {
	testDataDir string
}

var _ = C.Suite(NewCategoryTestSuite())

func NewCategoryTestSuite() *CategoryTestSuite {
	initCategory()
	s := &CategoryTestSuite{
		testDataDir: "testdata",
	}
	return s
}

func (s *CategoryTestSuite) TestGetCategoryInfos(c *C.C) {
	// infos := getCategoryInfos()
	// c.Check(infos[len(infos)-1].Id, C.Equals, -2)
	// c.Check(sort.IsSorted(CategoryInfosResult(infos[1:len(infos)-1])), C.Equals, true)
}

func (s *CategoryTestSuite) testGetDeepinCategory(c *C.C, name string, id CategoryId) {
	a := gio.NewDesktopAppInfo(name)
	if a == nil {
		c.Skip("create desktop app info failed")
		return
	}
	defer a.Unref()

	_id, err := getDeepinCategory(a)
	c.Check(err, C.IsNil)
	c.Check(_id, C.Equals, id)
}

func (s *CategoryTestSuite) TestGetDeepinCategory(c *C.C) {
	s.testGetDeepinCategory(c, filepath.Join(s.testDataDir, "deepin-music-player.desktop"), MultimediaID)
	s.testGetDeepinCategory(c, filepath.Join(s.testDataDir, "firefox.desktop"), NetworkID)
}
