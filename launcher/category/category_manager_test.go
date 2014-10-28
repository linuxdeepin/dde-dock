package category

import (
	C "launchpad.net/gocheck"
)

type CategoryManagerTestSuite struct {
	manager *CategoryManager
}

func (s *CategoryManagerTestSuite) SetUpTest(c *C.C) {
	s.manager = NewCategoryManager()
}

func (s *CategoryManagerTestSuite) TestGetCategory(c *C.C) {
	category := s.manager.GetCategory(AllID)
	c.Assert(category.Id(), C.Equals, AllID)

	category2 := s.manager.GetCategory(NetworkID)
	c.Assert(category2.Id(), C.Equals, NetworkID)

	category3 := s.manager.GetCategory(UnknownID)
	c.Assert(category3.Id(), C.IsNil)
}

func (s *CategoryManagerTestSuite) TestAddItem(c *C.C) {
	s.manager.AddItem("google-chrome", NetworkID)
	c.Assert(s.manager.categoryTable[NetworkID].Items(), C.HasLen, 1)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 1)

	s.manager.AddItem("firefox", NetworkID)
	c.Assert(s.manager.categoryTable[NetworkID].Items(), C.HasLen, 2)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 2)

	s.manager.AddItem("vim", DevelopmentID)
	c.Assert(s.manager.categoryTable[DevelopmentID].Items(), C.HasLen, 1)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 3)
}

func (s *CategoryManagerTestSuite) TestRemoveItem(c *C.C) {
	s.manager.AddItem("google-chrome", NetworkID)
	s.manager.AddItem("firefox", NetworkID)
	s.manager.AddItem("vim", DevelopmentID)
	c.Assert(s.manager.categoryTable[NetworkID].Items(), C.HasLen, 2)
	c.Assert(s.manager.categoryTable[DevelopmentID].Items(), C.HasLen, 1)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 3)

	s.manager.RemoveItem("vim", DevelopmentID)
	c.Assert(s.manager.categoryTable[NetworkID].Items(), C.HasLen, 2)
	c.Assert(s.manager.categoryTable[DevelopmentID].Items(), C.HasLen, 0)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 2)

	s.manager.RemoveItem("firefox", NetworkID)
	c.Assert(s.manager.categoryTable[NetworkID].Items(), C.HasLen, 1)
	c.Assert(s.manager.categoryTable[DevelopmentID].Items(), C.HasLen, 0)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 1)

	s.manager.RemoveItem("test", DevelopmentID)
	c.Assert(s.manager.categoryTable[NetworkID].Items(), C.HasLen, 1)
	c.Assert(s.manager.categoryTable[DevelopmentID].Items(), C.HasLen, 0)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 1)
}
