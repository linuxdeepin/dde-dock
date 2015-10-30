package launcher

import (
	C "launchpad.net/gocheck"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

type CategoryInfoExportTestSuite struct {
}

var _ = C.Suite(&CategoryInfoExportTestSuite{})

func (s *CategoryInfoExportTestSuite) TestContructor(c *C.C) {
	info := NewCategoryInfoExport(nil)
	c.Assert(info.Name, C.Equals, "")

	m := &MockCategoryInfo{CategoryID(1), "A", map[ItemID]bool{}}
	info = NewCategoryInfoExport(m)
	c.Assert(info.Name, C.Equals, "A")
	c.Assert(info.ID, C.Equals, CategoryID(1))
}
