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

	m := &MockCategoryInfo{CategoryId(1), "A", map[ItemId]bool{}}
	info = NewCategoryInfoExport(m)
	c.Assert(info.Name, C.Equals, "A")
	c.Assert(info.Id, C.Equals, CategoryId(1))
}
