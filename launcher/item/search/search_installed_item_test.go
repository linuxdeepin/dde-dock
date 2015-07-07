package search

import (
	C "launchpad.net/gocheck"
	"os"
	"path"
	. "pkg.deepin.io/dde-daemon/launcher/interfaces"
	. "pkg.deepin.io/dde-daemon/launcher/item"
	"pkg.deepin.io/lib/gio-2.0"
	"sync"
	"time"
)

type SearchInstalledItemTransactionTestSuite struct {
	testDataDir string
}

var _ = C.Suite(&SearchInstalledItemTransactionTestSuite{})

func (self *SearchInstalledItemTransactionTestSuite) SetUpSuite(c *C.C) {
	self.testDataDir = "../../testdata"
}

func (self *SearchInstalledItemTransactionTestSuite) TestSearchInstalledItemTransactionConstructor(c *C.C) {
	var transaction *SearchInstalledItemTransaction
	var err error

	ch := make(chan SearchResult)
	transaction, err = NewSearchInstalledItemTransaction(nil, nil, 4)
	c.Assert(transaction, C.IsNil)
	c.Assert(err, C.NotNil)
	c.Assert(err, C.Equals, SearchErrorNullChannel)

	transaction, err = NewSearchInstalledItemTransaction(ch, nil, 4)
	c.Assert(transaction, C.NotNil)
	c.Assert(err, C.IsNil)
}

func (self *SearchInstalledItemTransactionTestSuite) testSearchInstalledItemTransaction(c *C.C, key string, fn func([]SearchResult, *C.C), delay time.Duration, cancel bool) {
	old := os.Getenv("LANGUAGE")
	os.Setenv("LANGUAGE", "zh_CN.UTF-8")

	cancelChan := make(chan struct{})
	ch := make(chan SearchResult)
	transaction, _ := NewSearchInstalledItemTransaction(ch, cancelChan, 4)
	itemInfo := NewItem(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "firefox.desktop")))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(delay)
		transaction.Search(key, []ItemInfoInterface{itemInfo})
		wg.Done()
	}()
	if cancel {
		transaction.Cancel()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	res := []SearchResult{}
	for data := range ch {
		res = append(res, data)
	}

	fn(res, c)

	os.Setenv("LANGUAGE", old)
}

func (self *SearchInstalledItemTransactionTestSuite) TestSearchInstalledItemTransaction(c *C.C) {
	self.testSearchInstalledItemTransaction(c, "f", func(res []SearchResult, c *C.C) {
		c.Assert(len(res), C.Equals, 1)
	}, 0, false)
}

func (self *SearchInstalledItemTransactionTestSuite) TestSearchInstalledItemTransactionWithAItemNotExist(c *C.C) {
	self.testSearchInstalledItemTransaction(c, "IE", func(res []SearchResult, c *C.C) {
		c.Assert(len(res), C.Equals, 0)
	}, 0, false)
}

func (self *SearchInstalledItemTransactionTestSuite) TestSearchInstalledItemTransactionCancel(c *C.C) {
	self.testSearchInstalledItemTransaction(c, "fire", func(res []SearchResult, c *C.C) {
		c.Assert(len(res), C.Equals, 0)
	}, time.Second, true)
}
