package interfaces

import (
	"time"
)

type ItemManagerInterface interface {
	AddItem(ItemInfoInterface)
	RemoveItem(ItemId)
	HasItem(ItemId) bool
	GetItem(ItemId) ItemInfoInterface
	GetAllItems() []ItemInfoInterface
	GetAllFrequency(RateConfigFileInterface) map[ItemId]uint64
	GetAllTimeInstalled() (map[ItemId]int64, error)
	UninstallItem(ItemId, bool, time.Duration) error
	IsItemOnDesktop(ItemId) bool
	SendItemToDesktop(ItemId) error
	RemoveItemFromDesktop(ItemId) error
	GetRate(ItemId, RateConfigFileInterface) uint64
	SetRate(ItemId, uint64, RateConfigFileInterface)
	GetAllNewInstalledApps() ([]ItemId, error)
	MarkLaunched(ItemId) error
}
