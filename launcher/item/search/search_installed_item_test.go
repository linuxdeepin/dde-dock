/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package search

import (
	C "launchpad.net/gocheck"
	"os"
	"path"
	"pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/dde/daemon/launcher/item"
	"gir/gio-2.0"
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

	ch := make(chan Result)
	transaction, err = NewSearchInstalledItemTransaction(nil, nil, 4)
	c.Assert(transaction, C.IsNil)
	c.Assert(err, C.NotNil)
	c.Assert(err, C.Equals, ErrorSearchNullChannel)

	transaction, err = NewSearchInstalledItemTransaction(ch, nil, 4)
	c.Assert(transaction, C.NotNil)
	c.Assert(err, C.IsNil)
}

func (self *SearchInstalledItemTransactionTestSuite) testSearchInstalledItemTransaction(c *C.C, key string, fn func([]Result, *C.C), delay time.Duration, cancel bool) {
	old := os.Getenv("LANGUAGE")
	os.Setenv("LANGUAGE", "zh_CN.UTF-8")

	cancelChan := make(chan struct{})
	ch := make(chan Result)
	transaction, _ := NewSearchInstalledItemTransaction(ch, cancelChan, 4)
	itemInfo := item.New(gio.NewDesktopAppInfoFromFilename(path.Join(self.testDataDir, "firefox.desktop")))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(delay)
		transaction.Search(key, []interfaces.ItemInfo{itemInfo})
		wg.Done()
	}()
	if cancel {
		transaction.Cancel()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	res := []Result{}
	for data := range ch {
		res = append(res, data)
	}

	fn(res, c)

	os.Setenv("LANGUAGE", old)
}

func (self *SearchInstalledItemTransactionTestSuite) TestSearchInstalledItemTransaction(c *C.C) {
	self.testSearchInstalledItemTransaction(c, "f", func(res []Result, c *C.C) {
		c.Assert(len(res), C.Equals, 1)
	}, 0, false)
}

func (self *SearchInstalledItemTransactionTestSuite) TestSearchInstalledItemTransactionWithAItemNotExist(c *C.C) {
	self.testSearchInstalledItemTransaction(c, "IE", func(res []Result, c *C.C) {
		c.Assert(len(res), C.Equals, 0)
	}, 0, false)
}

func (self *SearchInstalledItemTransactionTestSuite) TestSearchInstalledItemTransactionCancel(c *C.C) {
	self.testSearchInstalledItemTransaction(c, "fire", func(res []Result, c *C.C) {
		c.Assert(len(res), C.Equals, 0)
	}, time.Second, true)
}
