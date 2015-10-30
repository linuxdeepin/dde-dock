package launcher

import (
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

type MockCategoryInfo struct {
	id    CategoryID
	name  string
	items map[ItemID]bool
}

func (c *MockCategoryInfo) ID() CategoryID {
	return c.id
}

func (c *MockCategoryInfo) Name() string {
	return c.name
}

func (c *MockCategoryInfo) LocaleName() string {
	// TODO: locale name
	return c.name
}

func (c *MockCategoryInfo) AddItem(itemID ItemID) {
	c.items[itemID] = true
}
func (c *MockCategoryInfo) RemoveItem(itemID ItemID) {
	delete(c.items, itemID)
}

func (c *MockCategoryInfo) Items() []ItemID {
	items := []ItemID{}
	for itemID := range c.items {
		items = append(items, itemID)
	}
	return items
}
