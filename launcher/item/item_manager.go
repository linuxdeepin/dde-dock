package item

import (
	storeApi "dbus/com/deepin/store/api"
	"encoding/json"
	"fmt"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	. "pkg.deepin.io/dde/daemon/launcher/item/softwarecenter"
	. "pkg.deepin.io/dde/daemon/launcher/utils"
	"sync"
	"time"
)

const (
	_NewSoftwareRecordFile = "launcher/new_software.ini"
	_NewSoftwareGroupName  = "NewInstalledApps"
	_NewSoftwareKeyName    = "Ids"
)

type ItemManager struct {
	lock      sync.Mutex
	itemTable map[ItemId]ItemInfoInterface
	soft      SoftwareCenterInterface
}

func NewItemManager(soft SoftwareCenterInterface) *ItemManager {
	return &ItemManager{
		itemTable: map[ItemId]ItemInfoInterface{},
		soft:      soft,
	}
}

func (m *ItemManager) AddItem(item ItemInfoInterface) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.itemTable[item.Id()] = item
}

func (m *ItemManager) HasItem(id ItemId) bool {
	_, ok := m.itemTable[id]
	return ok
}

func (m *ItemManager) RemoveItem(id ItemId) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.itemTable, id)
}

func (m *ItemManager) GetItem(id ItemId) ItemInfoInterface {
	item, _ := m.itemTable[id]
	return item
}

func (m *ItemManager) GetAllItems() []ItemInfoInterface {
	infos := []ItemInfoInterface{}
	for _, item := range m.itemTable {
		infos = append(infos, item)
	}
	return infos
}

func (m *ItemManager) UninstallItem(id ItemId, purge bool, timeout time.Duration) error {
	item := m.GetItem(id)
	if item == nil {
		return fmt.Errorf("No such a item: %q", id)
	}

	pkgName, err := GetPkgName(m.soft, item.Path())
	if err != nil {
		return err
	}

	if pkgName == "" {
		return fmt.Errorf("get package name of %q failed", string(id))
	}

	transaction := NewUninstallTransaction(m.soft, pkgName, purge, timeout)
	return transaction.Exec()
}

func (m *ItemManager) IsItemOnDesktop(id ItemId) bool {
	item := m.GetItem(id)
	if item == nil {
		return false
	}
	return isOnDesktop(item.Path())
}

func (m *ItemManager) SendItemToDesktop(id ItemId) error {
	if !m.HasItem(id) {
		return fmt.Errorf("No such a item %q", id)
	}

	if err := sendToDesktop(m.GetItem(id).Path()); err != nil {
		return err
	}

	return nil
}

func (m *ItemManager) RemoveItemFromDesktop(id ItemId) error {
	if !m.HasItem(id) {
		return fmt.Errorf("No such a item %q", id)
	}

	if err := removeFromDesktop(m.GetItem(id).Path()); err != nil {
		return err
	}

	return nil
}

func (m *ItemManager) GetRate(id ItemId, f RateConfigFileInterface) uint64 {
	rate, _ := f.GetUint64(string(id), _RateRecordKey)
	return rate
}

func (m *ItemManager) SetRate(id ItemId, rate uint64, f RateConfigFileInterface) {
	f.SetUint64(string(id), _RateRecordKey, rate)
	SaveKeyFile(f, ConfigFilePath(_RateRecordFile))
}

func (m *ItemManager) GetAllFrequency(f RateConfigFileInterface) (infos map[ItemId]uint64) {
	infos = map[ItemId]uint64{}
	if f == nil {
		for id, _ := range m.itemTable {
			infos[id] = 0
		}
		return
	}

	for id, _ := range m.itemTable {
		infos[id] = m.GetRate(id, f)
	}

	return
}

func (m *ItemManager) GetAllTimeInstalled() (map[ItemId]int64, error) {
	infos := map[ItemId]int64{}
	var err error
	for id, _ := range m.itemTable {
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
		id := GenId(data[0].(string))
		t := int64(data[1].(float64))
		infos[id] = t
	}

	return infos, err
}

func (self *ItemManager) GetAllNewInstalledApps() ([]ItemId, error) {
	ids := []ItemId{}
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
		id := GenId(data[0].(string))
		ids = append(ids, id)
	}
	return ids, nil
}

func (self *ItemManager) MarkNew(_id ItemId) error {
	return nil
}

func (self *ItemManager) MarkLaunched(_id ItemId) error {
	store, err := storeApi.NewDStoreDesktop("com.deepin.store.Api", "/com/deepin/store/Api")
	if err != nil {
		return fmt.Errorf("create store api failed: %v", err)
	}
	defer storeApi.DestroyDStoreDesktop(store)

	_, ok := store.MarkLaunched(string(_id))
	return ok
}
