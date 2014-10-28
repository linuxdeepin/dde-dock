package launcher

import (
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
)

type MockCategoryInfo struct {
	id    CategoryId
	name  string
	items map[ItemId]bool
}

func (c *MockCategoryInfo) Id() CategoryId {
	return c.id
}

func (c *MockCategoryInfo) Name() string {
	return c.name
}

func (c *MockCategoryInfo) AddItem(itemId ItemId) {
	c.items[itemId] = true
}
func (c *MockCategoryInfo) RemoveItem(itemId ItemId) {
	delete(c.items, itemId)
}

func (c *MockCategoryInfo) Items() []ItemId {
	items := []ItemId{}
	for itemId, _ := range c.items {
		items = append(items, itemId)
	}
	return items
}
