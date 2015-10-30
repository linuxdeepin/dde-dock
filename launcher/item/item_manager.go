package item

import (
	storeApi "dbus/com/deepin/store/api"
	"encoding/json"
	"fmt"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/dde/daemon/launcher/item/dstore"
	. "pkg.deepin.io/dde/daemon/launcher/utils"
	"sync"
	"time"
)

const (
	_NewSoftwareRecordFile = "launcher/new_software.ini"
	_NewSoftwareGroupName  = "NewInstalledApps"
	_NewSoftwareKeyName    = "Ids"
)

// Manager controls all items.
type Manager struct {
	lock      sync.Mutex
	itemTable map[ItemID]ItemInfo
	soft      DStore
}

// NewManager creates a new item manager.
func NewManager(soft DStore) *Manager {
	return &Manager{
		itemTable: map[ItemID]ItemInfo{},
		soft:      soft,
	}
}

// AddItem adds a new app.
func (m *Manager) AddItem(item ItemInfo) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.itemTable[item.ID()] = item
}

// HasItem returns true if the app is existed.
func (m *Manager) HasItem(id ItemID) bool {
	_, ok := m.itemTable[id]
	return ok
}

// RemoveItem removes a app.
func (m *Manager) RemoveItem(id ItemID) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.itemTable, id)
}

// GetItem returns a Item struct object if app exist, otherwise return nil.
func (m *Manager) GetItem(id ItemID) ItemInfo {
	item, _ := m.itemTable[id]
	return item
}

// GetAllItems returns all apps.
func (m *Manager) GetAllItems() []ItemInfo {
	infos := []ItemInfo{}
	for _, item := range m.itemTable {
		infos = append(infos, item)
	}
	return infos
}

// UninstallItem will uninstall a app.
func (m *Manager) UninstallItem(id ItemID, purge bool, timeout time.Duration) error {
	item := m.GetItem(id)
	if item == nil {
		return fmt.Errorf("No such a item: %q", id)
	}

	pkgName, err := dstore.GetPkgName(m.soft, item.Path())
	if err != nil {
		return err
	}

	if pkgName == "" {
		return fmt.Errorf("get package name of %q failed", string(id))
	}

	transaction := dstore.NewUninstallTransaction(m.soft, pkgName, purge, timeout)
	return transaction.Exec()
}

// IsItemOnDesktop returns true if app exists on desktop.
func (m *Manager) IsItemOnDesktop(id ItemID) bool {
	item := m.GetItem(id)
	if item == nil {
		return false
	}
	return isOnDesktop(item.Path())
}

// SendItemToDesktop sends a app to desktop.
func (m *Manager) SendItemToDesktop(id ItemID) error {
	if !m.HasItem(id) {
		return fmt.Errorf("No such a item %q", id)
	}

	if err := sendToDesktop(m.GetItem(id).Path()); err != nil {
		return err
	}

	return nil
}

// RemoveItemFromDesktop removes app from desktop.
func (m *Manager) RemoveItemFromDesktop(id ItemID) error {
	if !m.HasItem(id) {
		return fmt.Errorf("No such a item %q", id)
	}

	if err := removeFromDesktop(m.GetItem(id).Path()); err != nil {
		return err
	}

	return nil
}

// GetFrequency returns a item's  use frequency.
func (m *Manager) GetFrequency(id ItemID, f RateConfigFile) uint64 {
	rate, _ := f.GetUint64(string(id), _RateRecordKey)
	return rate
}

// SetFrequency sets a item's  use frequency.
func (m *Manager) SetFrequency(id ItemID, rate uint64, f RateConfigFile) {
	f.SetUint64(string(id), _RateRecordKey, rate)
	SaveKeyFile(f, ConfigFilePath(_RateRecordFile))
}

// GetAllFrequency returns all items' use frequency
func (m *Manager) GetAllFrequency(f RateConfigFile) (infos map[ItemID]uint64) {
	infos = map[ItemID]uint64{}
	if f == nil {
		for id := range m.itemTable {
			infos[id] = 0
		}
		return
	}

	for id := range m.itemTable {
		infos[id] = m.GetFrequency(id, f)
	}

	return
}

// GetAllTimeInstalled returns all items installed time.
// TODO:
// 1. do it once.
// 2. update it when item changed.
func (m *Manager) GetAllTimeInstalled() (map[ItemID]int64, error) {
	infos := map[ItemID]int64{}
	var err error
	for id := range m.itemTable {
		infos[id] = 0
	}

	store, err := storeApi.NewDStoreDesktop("com.deepin.store.Api", "/com/deepin/store/Api")
	if err != nil {
		return infos, fmt.Errorf("create store api failed: %v", err)
	}
	defer storeApi.DestroyDStoreDesktop(store)

	datasStr, err := store.GetAllDesktops()
	if err != nil {
		return infos, fmt.Errorf("get all desktops' info failed: %v", err)
	}

	datas := [][]interface{}{}
	err = json.Unmarshal([]byte(datasStr), &datas)
	if err != nil {
		return infos, err
	}

	for _, data := range datas {
		id := GenID(data[0].(string))
		t := int64(data[1].(float64))
		infos[id] = t
	}

	return infos, err
}

// GetAllNewInstalledApps returns all apps newly installed.
func (self *Manager) GetAllNewInstalledApps() ([]ItemID, error) {
	ids := []ItemID{}
	store, err := storeApi.NewDStoreDesktop("com.deepin.store.Api", "/com/deepin/store/Api")
	if err != nil {
		return ids, fmt.Errorf("create store api failed: %v", err)
	}
	defer storeApi.DestroyDStoreDesktop(store)

	dataStr, err := store.GetNewDesktops()
	if err != nil {
		return ids, err
	}

	datas := [][]interface{}{}
	err = json.Unmarshal([]byte(dataStr), &datas)
	if err != nil {
		return ids, err
	}

	for _, data := range datas {
		id := GenID(data[0].(string))
		ids = append(ids, id)
	}
	return ids, nil
}

// MarkNew marks a item as newly installed.
func (self *Manager) MarkNew(_id ItemID) error {
	return nil
}

// MarkLaunched marks a item as launched, it won't be newly installed.
func (self *Manager) MarkLaunched(_id ItemID) error {
	store, err := storeApi.NewDStoreDesktop("com.deepin.store.Api", "/com/deepin/store/Api")
	if err != nil {
		return fmt.Errorf("create store api failed: %v", err)
	}
	defer storeApi.DestroyDStoreDesktop(store)

	_, ok := store.MarkLaunched(string(_id))
	return ok
}
