// +build ignore

package item

import (
	"fmt"
	C "launchpad.net/gocheck"
	"math/rand"
	"os"
	"path"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
	"sync"
	"time"
)

type ItemManagerTestSuite struct {
	softcenter  *MockSoftcenter
	m           ItemManager
	item        ItemInfo
	timeout     time.Duration
	testDataDir string
	oldHome     string
	f           *glib.KeyFile
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

	s.item = New(firefox)
	firefox.Unref()
}

func (s *ItemManagerTestSuite) TearDownSuite(c *C.C) {
	os.Setenv("HOME", s.oldHome)
	s.f.Free()
}

func (s *ItemManagerTestSuite) SetUpTest(c *C.C) {
	s.softcenter = NewMockSoftcenter()
	s.m = NewManager(s.softcenter)
	s.timeout = time.Second * 10
}

func (s *ItemManagerTestSuite) TestItemManager(c *C.C) {
	c.Assert(s.m.GetItem(s.item.ID()), C.IsNil)
	c.Assert(s.m.HasItem(s.item.ID()), C.Equals, false)

	s.m.AddItem(s.item)
	c.Assert(s.m.GetItem(s.item.ID()).ID(), C.Equals, s.item.ID())
	c.Assert(s.m.HasItem(s.item.ID()), C.Equals, true)

	s.m.RemoveItem(s.item.ID())
	c.Assert(s.m.GetItem(s.item.ID()), C.IsNil)
	c.Assert(s.m.HasItem(s.item.ID()), C.Equals, false)
}

func (s *ItemManagerTestSuite) addTestItem(c *C.C, path string) ItemInfo {
	desktop := gio.NewDesktopAppInfoFromFilename(path)
	if desktop == nil {
		c.Skip(createDesktopFailed(path))
	}
	item := New(desktop)
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

	err := s.m.UninstallItem(softcenter.ID(), false, s.timeout)
	c.Assert(err, C.NotNil)
	c.Assert(err.Error(), C.Equals, fmt.Sprintf("get package name of %q failed", softcenter.ID()))

}

func (s *ItemManagerTestSuite) TestUnistallExistItem(c *C.C) {
	s.m.AddItem(s.item)

	err := s.m.UninstallItem(s.item.ID(), false, s.timeout)
	c.Assert(err, C.IsNil)
	c.Assert(s.softcenter.count, C.Equals, 1)
}

func (s *ItemManagerTestSuite) TestUnistallMultiItem(c *C.C) {
	var err error
	s.m.AddItem(s.item)
	player := s.addTestItem(c, path.Join(s.testDataDir, "deepin-music-player.desktop"))

	err = s.m.UninstallItem(s.item.ID(), false, s.timeout)
	c.Assert(err, C.IsNil)
	c.Assert(s.softcenter.count, C.Equals, 1)

	err = s.m.UninstallItem(player.ID(), false, s.timeout)
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
		err = s.m.UninstallItem(s.item.ID(), false, s.timeout)
		c.Assert(err, C.IsNil)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = s.m.UninstallItem(player.ID(), false, s.timeout)
		c.Assert(err, C.IsNil)
	}()

	wg.Wait()
	c.Assert(s.softcenter.count, C.Equals, 2)
}

func (s *ItemManagerTestSuite) TestSetFrequency(c *C.C) {
	s.m.SetFrequency("firefox", uint64(3), s.f)
	c.Assert(s.m.GetFrequency("firefox", s.f), C.Equals, uint64(3))

	s.m.SetFrequency("firefox", uint64(2), s.f)
	c.Assert(s.m.GetFrequency("firefox", s.f), C.Equals, uint64(2))
}
