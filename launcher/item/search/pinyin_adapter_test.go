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
	"gir/gio-2.0"
)

type PinYinTestSuite struct {
	testDataDir string
}

// var _ = C.Suite(&PinYinTestSuite{})

func (self *PinYinTestSuite) SetUpSuite(c *C.C) {
	self.testDataDir = "../../testdata"
}

func (self *PinYinTestSuite) TestPinYin(c *C.C) {
	names := []string{}
	oldLang := os.Getenv("LANGUAGE")
	os.Setenv("LANGUAGE", "zh_CN.UTF-8")
	addName := func(m *[]string, n string) {
		app := gio.NewDesktopAppInfoFromFilename(n)
		if app == nil {
			c.Skip("create desktop app info failed")
			return
		}
		defer app.Unref()
		name := app.GetDisplayName()
		c.Logf("add %q to names", name)
		*m = append(*m, name)
	}
	addName(&names, path.Join(self.testDataDir, "deepin-software-center.desktop"))
	addName(&names, path.Join(self.testDataDir, "firefox.desktop"))
	tree, err := NewPinYinSearchAdapter(names)
	if err != nil {
		c.Log(err)
		c.Fail()
	}
	search := func(key string, res []string) {
		keys, err := tree.Search(key)
		if err != nil {
			c.Log(err)
			c.Fail()
			return
		}
		c.Assert(keys, C.DeepEquals, res)
	}
	search("shang", []string{"深度商店"})
	search("sd", []string{"深度商店"})
	search("商店", []string{"深度商店"})
	search("firefox", []string{"Firefox 网络浏览器"})
	search("wang", []string{"Firefox 网络浏览器"})
	search("网络", []string{"Firefox 网络浏览器"})

	os.Setenv("LANGUAGE", oldLang)
}
