package interfaces

import (
	"gir/gio-2.0"
)

// CategoryID is the type for category id.
type CategoryID int64

// CategoryInfo is interface for category info.
type CategoryInfo interface {
	ID() CategoryID
	Name() string
	LocaleName() string
	Items() []ItemID
	AddItem(ItemID)
	RemoveItem(ItemID)
}

// CategoryManager is interface for category manager.
type CategoryManager interface {
	AddItem(ItemID, CategoryID)
	RemoveItem(ItemID, CategoryID)
	GetAllCategory() []CategoryID
	GetCategory(id CategoryID) CategoryInfo
	LoadAppCategoryInfo(deepin string, xcategory string) error
	FreeAppCategoryInfo()
	QueryID(app *gio.DesktopAppInfo) (CategoryID, error)
}
