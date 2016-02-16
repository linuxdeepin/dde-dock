/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher

import (
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

// ItemInfoExport is a wrapper struct used to export info to dbus .
type ItemInfoExport struct {
	Path          string
	Name          string
	ID            ItemID
	Icon          string
	CategoryID    CategoryID
	TimeInstalled int64
}

// NewItemInfoExport creates a new ItemInfoExport from ItemInfo.
func NewItemInfoExport(item ItemInfo) ItemInfoExport {
	if item == nil {
		return ItemInfoExport{}
	}
	return ItemInfoExport{
		Path:          item.Path(),
		Name:          item.LocaleName(),
		ID:            item.ID(),
		Icon:          item.Icon(),
		CategoryID:    item.CategoryID(),
		TimeInstalled: item.TimeInstalled(),
	}
}

// CategoryInfoExport is a wrapper struct used to export info to dbus.
type CategoryInfoExport struct {
	Name  string
	ID    CategoryID
	Items []ItemID
}

// NewCategoryInfoExport creates a new CategoryInfoExport from CategoryInfo.
func NewCategoryInfoExport(c CategoryInfo) CategoryInfoExport {
	if c == nil {
		return CategoryInfoExport{}
	}
	return CategoryInfoExport{
		Name:  c.Name(),
		ID:    c.ID(),
		Items: c.Items(),
	}
}

// FrequencyExport is a wrapper struct used to export info to dbus.
type FrequencyExport struct {
	ID        ItemID
	Frequency uint64
}

// TimeInstalledExport is a wrapper struct used to export info to dbus.
type TimeInstalledExport struct {
	ID   ItemID
	Time int64
}
