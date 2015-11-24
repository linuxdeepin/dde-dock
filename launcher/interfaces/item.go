package interfaces

import (
	"pkg.deepin.io/lib/glib-2.0"
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
}

// ItemManager is interface for item manager.
type ItemManager interface {
	AddItem(ItemInfo)
	RemoveItem(ItemID)
	HasItem(ItemID) bool
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
