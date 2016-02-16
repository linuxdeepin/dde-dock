/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher

import (
	C "launchpad.net/gocheck"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/dde/daemon/launcher/mock"
)

type CategoryInfoExportTestSuite struct {
}

var _ = C.Suite(&CategoryInfoExportTestSuite{})

func (s *CategoryInfoExportTestSuite) TestContructor(c *C.C) {
	info := NewCategoryInfoExport(nil)
	c.Assert(info.Name, C.Equals, "")

	m := mock.NewCategoryInfo(CategoryID(1), "A")
	info = NewCategoryInfoExport(m)
	c.Assert(info.Name, C.Equals, "A")
	c.Assert(info.ID, C.Equals, CategoryID(1))
}
