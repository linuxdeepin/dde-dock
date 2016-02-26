/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
	LoadCategoryInfo() error
	FreeAppCategoryInfo()
	QueryID(app *gio.DesktopAppInfo) (CategoryID, error)
}
