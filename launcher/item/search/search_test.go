package search

import (
	"os"
	"path"
	"testing"
	"time"

	C "launchpad.net/gocheck"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/dde/daemon/launcher/item"
	"gir/gio-2.0"
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
	var transaction *Transaction
	var err error

	transaction, err = NewTransaction(nil, nil, nil, 4)
	c.Assert(transaction, C.IsNil)
	c.Assert(err, C.NotNil)
	c.Assert(err, C.Equals, ErrorSearchNullChannel)

	transaction, err = NewTransaction(nil, make(chan Result), nil, 4)
	c.Assert(transaction, C.NotNil)
	c.Assert(err, C.IsNil)
}

func (self *SearchTransactionTestSuite) testTransaction(c *C.C, pinyinObj PinYin, key string, fn func([]Result, *C.C), delay time.Duration, cancel bool) {
	old := os.Getenv("LANGUAGE")
	os.Setenv("LANGUAGE", "zh_CN.UTF-8")

	cancelChan := make(chan struct{})
	ch := make(chan Result)
	transaction, _ := NewTransaction(pinyinObj, ch, cancelChan, 4)
	firefoxItemInfo := item.New(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "firefox.desktop")))
	playerItemInfo := item.New(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "deepin-music-player.desktop")))
	chromeItemInfo := item.New(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "google-chrome.desktop")))
	go func() {
		time.Sleep(delay)
		transaction.Search(key, []ItemInfo{
			firefoxItemInfo,
			playerItemInfo,
			chromeItemInfo,
		})
		close(ch)
	}()
	if cancel {
		transaction.Cancel()
	}

	result := map[ItemID]Result{}
	for itemInfo := range ch {
		result[itemInfo.ID] = itemInfo
	}

	res := []Result{}
	for _, data := range result {
		res = append(res, data)
	}

	fn(res, c)

	os.Setenv("LANGUAGE", old)
}

func (self *SearchTransactionTestSuite) TestTransaction(c *C.C) {
	self.testTransaction(c, nil, "fire", func(res []Result, c *C.C) {
		c.Assert(len(res), C.Equals, 1)
	}, 0, false)
}

func (self *SearchTransactionTestSuite) TestSearchTransactionWithAItemNotExist(c *C.C) {
	self.testTransaction(c, nil, "IE", func(res []Result, c *C.C) {
		c.Assert(len(res), C.Equals, 0)
	}, 0, false)
}

func (self *SearchTransactionTestSuite) TestSearchTransactionWithPinYin(c *C.C) {
	old := os.Getenv("LANGUAGE")
	os.Setenv("LANGUAGE", "zh_CN.UTF-8")
	fireItemInfo := item.New(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "firefox.desktop")))
	chromeItemInfo := item.New(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "google-chrome.desktop")))
	defer os.Setenv("LANGUAGE", old)

	pinyinObj := NewMockPinYin(map[string][]string{
		"liu": []string{
			fireItemInfo.LocaleName(),
			chromeItemInfo.LocaleName(),
		},
	}, true)
	self.testTransaction(c, pinyinObj, "liu", func(res []Result, c *C.C) {
		c.Assert(len(res), C.Equals, 2)
	}, 0, false)

}

func (self *SearchTransactionTestSuite) TestSearchCancel(c *C.C) {
	self.testTransaction(c, nil, "fire", func(res []Result, c *C.C) {
		c.Assert(len(res), C.Equals, 0)
	}, time.Second, true)
}
