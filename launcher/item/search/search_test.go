package search

import (
	C "launchpad.net/gocheck"
	"os"
	"path"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	. "pkg.deepin.io/dde/daemon/launcher/item"
	"pkg.deepin.io/lib/gio-2.0"
	"testing"
	"time"
)

func TestSearch(t *testing.T) {
	C.TestingT(t)
}

type SearchTransactionTestSuite struct {
	testDataDir string
}

var _ = C.Suite(&SearchTransactionTestSuite{})

func (self *SearchTransactionTestSuite) SetUpSuite(c *C.C) {
	self.testDataDir = "../../testdata"
}

func (self *SearchTransactionTestSuite) TestSearchTransactionConstructor(c *C.C) {
	var transaction *SearchTransaction
	var err error

	transaction, err = NewSearchTransaction(nil, nil, nil, 4)
	c.Assert(transaction, C.IsNil)
	c.Assert(err, C.NotNil)
	c.Assert(err, C.Equals, SearchErrorNullChannel)

	transaction, err = NewSearchTransaction(nil, make(chan SearchResult), nil, 4)
	c.Assert(transaction, C.NotNil)
	c.Assert(err, C.IsNil)
}

func (self *SearchTransactionTestSuite) testSearchTransaction(c *C.C, pinyinObj PinYinInterface, key string, fn func([]SearchResult, *C.C), delay time.Duration, cancel bool) {
	old := os.Getenv("LANGUAGE")
	os.Setenv("LANGUAGE", "zh_CN.UTF-8")

	cancelChan := make(chan struct{})
	ch := make(chan SearchResult)
	transaction, _ := NewSearchTransaction(pinyinObj, ch, cancelChan, 4)
	firefoxItemInfo := NewItem(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "firefox.desktop")))
	playerItemInfo := NewItem(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "deepin-music-player.desktop")))
	chromeItemInfo := NewItem(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "google-chrome.desktop")))
	go func() {
		time.Sleep(delay)
		transaction.Search(key, []ItemInfoInterface{
			firefoxItemInfo,
			playerItemInfo,
			chromeItemInfo,
		})
		close(ch)
	}()
	if cancel {
		transaction.Cancel()
	}

	result := map[ItemId]SearchResult{}
	for itemInfo := range ch {
		result[itemInfo.Id] = itemInfo
	}

	res := []SearchResult{}
	for _, data := range result {
		res = append(res, data)
	}

	fn(res, c)

	os.Setenv("LANGUAGE", old)
}

func (self *SearchTransactionTestSuite) TestSearchTransaction(c *C.C) {
	self.testSearchTransaction(c, nil, "fire", func(res []SearchResult, c *C.C) {
		c.Assert(len(res), C.Equals, 1)
	}, 0, false)
}

func (self *SearchTransactionTestSuite) TestSearchTransactionWithAItemNotExist(c *C.C) {
	self.testSearchTransaction(c, nil, "IE", func(res []SearchResult, c *C.C) {
		c.Assert(len(res), C.Equals, 0)
	}, 0, false)
}

func (self *SearchTransactionTestSuite) TestSearchTransactionWithPinYin(c *C.C) {
	old := os.Getenv("LANGUAGE")
	os.Setenv("LANGUAGE", "zh_CN.UTF-8")
	fireItemInfo := NewItem(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "firefox.desktop")))
	chromeItemInfo := NewItem(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "google-chrome.desktop")))
	os.Setenv("LANGUAGE", old)

	pinyinObj := NewMockPinYin(map[string][]string{
		// both GenericName contains 浏
		"liu": []string{
			fireItemInfo.GenericName(),
			chromeItemInfo.GenericName(),
		},
	}, true)
	self.testSearchTransaction(c, pinyinObj, "liu", func(res []SearchResult, c *C.C) {
		c.Assert(len(res), C.Equals, 2)
	}, 0, false)

}

func (self *SearchTransactionTestSuite) TestSearchCancel(c *C.C) {
	self.testSearchTransaction(c, nil, "fire", func(res []SearchResult, c *C.C) {
		c.Assert(len(res), C.Equals, 0)
	}, time.Second, true)
}
