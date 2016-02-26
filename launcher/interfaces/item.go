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
	"gir/glib-2.0"
	"time"
)

// ItemID is type for item's id.
type ItemID string

// ItemInfo is interface for item info.
type ItemInfo interface {
	Name() string
	Icon() string
	Path() string
	ID() ItemID
	ExecCmd() string
	Description() string
	LocaleName() string
	GenericName() string
	Keywords() []string
	CategoryID() CategoryID
	SetCategoryID(CategoryID)
	TimeInstalled() int64
	SetTimeInstalled(int64)
	LastModifiedTime() int64
	Refresh()
}

// ItemManager is interface for item manager.
type ItemManager interface {
	AddItem(ItemInfo)
	RemoveItem(ItemID)
	HasItem(ItemID) bool
	RefreshItem(ItemID)
	GetItem(ItemID) ItemInfo
	GetAllItems() []ItemInfo
	GetAllFrequency(*glib.KeyFile) map[ItemID]uint64
	GetAllTimeInstalled() (map[ItemID]int64, error)
	UninstallItem(ItemID, bool, time.Duration) error
	IsItemOnDesktop(ItemID) bool
	SendItemToDesktop(ItemID) error
	RemoveItemFromDesktop(ItemID) error
	GetFrequency(ItemID, *glib.KeyFile) uint64
	SetFrequency(ItemID, uint64, *glib.KeyFile)
	GetAllNewInstalledApps() ([]ItemID, error)
	MarkLaunched(ItemID) error
}
