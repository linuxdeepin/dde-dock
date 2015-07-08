package item

import (
	"fmt"
	C "launchpad.net/gocheck"
	"math/rand"
	"os"
	"path"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/lib/gio-2.0"
	"sync"
	"time"
)

type ItemManagerTestSuite struct {
	softcenter  *MockSoftcenter
	m           ItemManagerInterface
	item        ItemInfoInterface
	timeout     time.Duration
	testDataDir string
	oldHome     string
	f           RateConfigFileInterface
}

var _ = C.Suite(&ItemManagerTestSuite{})

func createDesktopFailed(path string) string {
	return fmt.Sprintf("create desktop(%s) failed", path)
}

func (s *ItemManagerTestSuite) SetUpSuite(c *C.C) {
	rand.Seed(time.Now().UTC().UnixNano())
	s.testDataDir = "../testdata"
	firefoxDesktopPath := path.Join(s.testDataDir, "firefox.desktop")
	firefox := gio.NewDesktopAppInfoFromFilename(firefoxDesktopPath)
	if firefox == nil {
		c.Skip(createDesktopFailed(firefoxDesktopPath))
	}

	// according to the sources of glib.
	s.oldHome = os.Getenv("HOME")
	os.Setenv("HOME", s.testDataDir)

	var err error
	s.f, err = GetFrequencyRecordFile()
	if err != nil {
		c.Skip("get config file failed")
	}

	s.item = NewItem(firefox)
	firefox.Unref()
}

func (s *ItemManagerTestSuite) TearDownSuite(c *C.C) {
	os.Setenv("HOME", s.oldHome)
	s.f.Free()
}

func (s *ItemManagerTestSuite) SetUpTest(c *C.C) {
	s.softcenter = NewMockSoftcenter()
	s.m = NewItemManager(s.softcenter)
	s.timeout = time.Second * 10
}

func (s *ItemManagerTestSuite) TestItemManager(c *C.C) {
	c.Assert(s.m.GetItem(s.item.Id()), C.IsNil)
	c.Assert(s.m.HasItem(s.item.Id()), C.Equals, false)

	s.m.AddItem(s.item)
	c.Assert(s.m.GetItem(s.item.Id()).Id(), C.Equals, s.item.Id())
	c.Assert(s.m.HasItem(s.item.Id()), C.Equals, true)

	s.m.RemoveItem(s.item.Id())
	c.Assert(s.m.GetItem(s.item.Id()), C.IsNil)
	c.Assert(s.m.HasItem(s.item.Id()), C.Equals, false)
}

func (s *ItemManagerTestSuite) addTestItem(c *C.C, path string) ItemInfoInterface {
	desktop := gio.NewDesktopAppInfoFromFilename(path)
	if desktop == nil {
		c.Skip(createDesktopFailed(path))
	}
	item := NewItem(desktop)
	s.m.AddItem(item)
	desktop.Unref()
	return item
}

func (s *ItemManagerTestSuite) TestUnistallNotExistItem(c *C.C) {
	err := s.m.UninstallItem("test", false, s.timeout)
	c.Assert(err, C.NotNil)
	c.Assert(err.Error(), C.Matches, "No such a item:.*")
}

func (s *ItemManagerTestSuite) TestUnistallItemNoOnSoftCenter(c *C.C) {
	softcenter := s.addTestItem(c, path.Join(s.testDataDir, "deepin-software-center.desktop"))

	err := s.m.UninstallItem(softcenter.Id(), false, s.timeout)
	c.Assert(err, C.NotNil)
	c.Assert(err.Error(), C.Equals, fmt.Sprintf("get package name of %q failed", softcenter.Id()))

}

func (s *ItemManagerTestSuite) TestUnistallExistItem(c *C.C) {
	s.m.AddItem(s.item)

	err := s.m.UninstallItem(s.item.Id(), false, s.timeout)
	c.Assert(err, C.IsNil)
	c.Assert(s.softcenter.count, C.Equals, 1)
}

func (s *ItemManagerTestSuite) TestUnistallMultiItem(c *C.C) {
	var err error
	s.m.AddItem(s.item)
	player := s.addTestItem(c, path.Join(s.testDataDir, "deepin-music-player.desktop"))

	err = s.m.UninstallItem(s.item.Id(), false, s.timeout)
	c.Assert(err, C.IsNil)
	c.Assert(s.softcenter.count, C.Equals, 1)

	err = s.m.UninstallItem(player.Id(), false, s.timeout)
	c.Assert(err, C.IsNil)
	c.Assert(s.softcenter.count, C.Equals, 2)
}

func (s *ItemManagerTestSuite) TestUnistallMultiItemAsync(c *C.C) {
	// FIXME: is this test right.
	var err error
	s.m.AddItem(s.item)
	player := s.addTestItem(c, path.Join(s.testDataDir, "deepin-music-player.desktop"))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = s.m.UninstallItem(s.item.Id(), false, s.timeout)
		c.Assert(err, C.IsNil)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = s.m.UninstallItem(player.Id(), false, s.timeout)
		c.Assert(err, C.IsNil)
	}()

	wg.Wait()
	c.Assert(s.softcenter.count, C.Equals, 2)
}

func (s *ItemManagerTestSuite) TestGetRate(c *C.C) {
	c.Assert(s.m.GetRate(ItemId("firefox"), s.f), C.Equals, uint64(2))
	c.Assert(s.m.GetRate(ItemId("deepin-software-center"), s.f), C.Equals, uint64(0))
}

func (s *ItemManagerTestSuite) TestSetRate(c *C.C) {
	s.m.SetRate("firefox", uint64(3), s.f)
	c.Assert(s.m.GetRate("firefox", s.f), C.Equals, uint64(3))

	s.m.SetRate("firefox", uint64(2), s.f)
	c.Assert(s.m.GetRate("firefox", s.f), C.Equals, uint64(2))
}
