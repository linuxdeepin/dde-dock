package category

import (
	C "launchpad.net/gocheck"
)

type ManagerTestSuite struct {
	manager *Manager
}

func (s *ManagerTestSuite) SetUpTest(c *C.C) {
	s.manager = NewManager(nil, GetAllInfos(""))
}

func (s *ManagerTestSuite) TestGetCategory(c *C.C) {
	category := s.manager.GetCategory(AllID)
	c.Assert(category.ID(), C.Equals, AllID)

	category2 := s.manager.GetCategory(InternetID)
	c.Assert(category2.ID(), C.Equals, InternetID)

	category3 := s.manager.GetCategory(OthersID)
	c.Assert(category3.ID(), C.IsNil)
}

func (s *ManagerTestSuite) TestAddItem(c *C.C) {
	s.manager.AddItem("google-chrome", InternetID)
	c.Assert(s.manager.categoryTable[InternetID].Items(), C.HasLen, 1)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 1)

	s.manager.AddItem("firefox", InternetID)
	c.Assert(s.manager.categoryTable[InternetID].Items(), C.HasLen, 2)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 2)

	s.manager.AddItem("vim", DevelopmentID)
	c.Assert(s.manager.categoryTable[DevelopmentID].Items(), C.HasLen, 1)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 3)
}

func (s *ManagerTestSuite) TestRemoveItem(c *C.C) {
	s.manager.AddItem("google-chrome", InternetID)
	s.manager.AddItem("firefox", InternetID)
	s.manager.AddItem("vim", DevelopmentID)
	c.Assert(s.manager.categoryTable[InternetID].Items(), C.HasLen, 2)
	c.Assert(s.manager.categoryTable[DevelopmentID].Items(), C.HasLen, 1)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 3)

	s.manager.RemoveItem("vim", DevelopmentID)
	c.Assert(s.manager.categoryTable[InternetID].Items(), C.HasLen, 2)
	c.Assert(s.manager.categoryTable[DevelopmentID].Items(), C.HasLen, 0)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 2)

	s.manager.RemoveItem("firefox", InternetID)
	c.Assert(s.manager.categoryTable[InternetID].Items(), C.HasLen, 1)
	c.Assert(s.manager.categoryTable[DevelopmentID].Items(), C.HasLen, 0)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 1)

	s.manager.RemoveItem("test", DevelopmentID)
	c.Assert(s.manager.categoryTable[InternetID].Items(), C.HasLen, 1)
	c.Assert(s.manager.categoryTable[DevelopmentID].Items(), C.HasLen, 0)
	c.Assert(s.manager.categoryTable[AllID].Items(), C.HasLen, 1)
}
