package category

import (
	C "launchpad.net/gocheck"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"testing"
)

func TestCategory(t *testing.T) {
	C.TestingT(t)
}

type CategoryTestSuite struct {
	testDataDir string
}

var _ = C.Suite(&CategoryTestSuite{})

func (s *CategoryTestSuite) TestGetID(c *C.C) {
	cf := &Info{AllID, "all", map[ItemID]struct{}{}}
	c.Assert(cf.ID(), C.Equals, AllID)
}
func (s *CategoryTestSuite) TestGetName(c *C.C) {
	cf := &Info{AllID, "all", map[ItemID]struct{}{}}
	c.Assert(cf.Name(), C.Equals, "all")
}

func (s *CategoryTestSuite) TestGetAddItem(c *C.C) {
	cf := &Info{AllID, "all", map[ItemID]struct{}{}}
	c.Assert(cf.items, C.DeepEquals, make(map[ItemID]struct{}, 0))
	cf.AddItem(ItemID("test"))
	c.Assert(cf.items, C.DeepEquals, map[ItemID]struct{}{ItemID("test"): struct{}{}})
}

func (s *CategoryTestSuite) TestRemoveItem(c *C.C) {
	cf := &Info{AllID, "all", map[ItemID]struct{}{}}
	c.Assert(cf.items, C.DeepEquals, make(map[ItemID]struct{}, 0))
	cf.AddItem(ItemID("test"))
	c.Assert(cf.items, C.DeepEquals, map[ItemID]struct{}{ItemID("test"): struct{}{}})
	cf.RemoveItem(ItemID("test"))
	c.Assert(cf.items, C.DeepEquals, make(map[ItemID]struct{}, 0))
}

func (s *CategoryTestSuite) TestItems(c *C.C) {
	cf := &Info{AllID, "all", map[ItemID]struct{}{}}
	c.Assert(cf.Items(), C.DeepEquals, []ItemID{})

	cf.AddItem(ItemID("test"))
	c.Assert(cf.Items(), C.HasLen, 1)

	cf.AddItem(ItemID("test2"))
	c.Assert(cf.Items(), C.HasLen, 2)

	cf.RemoveItem(ItemID("test"))
	c.Assert(cf.Items(), C.HasLen, 1)
}

// TODO: fake a db to test QueryID
