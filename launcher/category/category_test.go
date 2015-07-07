package category

import (
	C "launchpad.net/gocheck"
	. "pkg.deepin.io/dde-daemon/launcher/interfaces"
	"testing"
)

func TestCategory(t *testing.T) {
	C.TestingT(t)
}

type CategoryTestSuite struct {
	testDataDir string
}

var _ = C.Suite(&CategoryTestSuite{})

func (s *CategoryTestSuite) TestGetId(c *C.C) {
	cf := &CategoryInfo{AllID, "all", map[ItemId]struct{}{}}
	c.Assert(cf.Id(), C.Equals, AllID)
}
func (s *CategoryTestSuite) TestGetName(c *C.C) {
	cf := &CategoryInfo{AllID, "all", map[ItemId]struct{}{}}
	c.Assert(cf.Name(), C.Equals, "all")
}

func (s *CategoryTestSuite) TestGetAddItem(c *C.C) {
	cf := &CategoryInfo{AllID, "all", map[ItemId]struct{}{}}
	c.Assert(cf.items, C.DeepEquals, make(map[ItemId]struct{}, 0))
	cf.AddItem(ItemId("test"))
	c.Assert(cf.items, C.DeepEquals, map[ItemId]struct{}{ItemId("test"): struct{}{}})
}

func (s *CategoryTestSuite) TestRemoveItem(c *C.C) {
	cf := &CategoryInfo{AllID, "all", map[ItemId]struct{}{}}
	c.Assert(cf.items, C.DeepEquals, make(map[ItemId]struct{}, 0))
	cf.AddItem(ItemId("test"))
	c.Assert(cf.items, C.DeepEquals, map[ItemId]struct{}{ItemId("test"): struct{}{}})
	cf.RemoveItem(ItemId("test"))
	c.Assert(cf.items, C.DeepEquals, make(map[ItemId]struct{}, 0))
}

func (s *CategoryTestSuite) TestItems(c *C.C) {
	cf := &CategoryInfo{AllID, "all", map[ItemId]struct{}{}}
	c.Assert(cf.Items(), C.DeepEquals, []ItemId{})

	cf.AddItem(ItemId("test"))
	c.Assert(cf.Items(), C.HasLen, 1)

	cf.AddItem(ItemId("test2"))
	c.Assert(cf.Items(), C.HasLen, 2)

	cf.RemoveItem(ItemId("test"))
	c.Assert(cf.Items(), C.HasLen, 1)
}

// TODO: fake a db to test QueryCategoryId
