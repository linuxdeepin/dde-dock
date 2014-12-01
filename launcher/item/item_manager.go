package item

import (
	"errors"
	"fmt"
	"os"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/item/softwarecenter"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/utils"
	"pkg.linuxdeepin.com/lib/glib-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
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
		return errors.New("get package name failed")
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
		return errors.New(fmt.Sprintf("No such a item %q", id))
	}

	if err := sendToDesktop(m.GetItem(id).Path()); err != nil {
		return err
	}

	return nil
}

func (m *ItemManager) RemoveItemFromDesktop(id ItemId) error {
	if !m.HasItem(id) {
		return errors.New(fmt.Sprintf("No such a item %q", id))
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

type _Time struct {
	Time int64
	Id   ItemId
}

func (m *ItemManager) GetAllTimeInstalled() (infos map[ItemId]int64) {
	infos = map[ItemId]int64{}
	dataChan := make(chan ItemInfoInterface)
	go func() {
		for _, item := range m.itemTable {
			dataChan <- item
		}
		close(dataChan)
	}()

	const N = 20
	var wg sync.WaitGroup
	wg.Add(N)
	timeChan := make(chan _Time)
	for i := 0; i < N; i++ {
		go func() {
			for item := range dataChan {
				// NOTE:
				// the real installation time is hard to get.
				// using modification time as install time for now.

				// pkgName, err := GetPkgName(m.soft, item.Path())
				// if err != nil {
				// 	timeChan <- _Time{Id: item.Id(), Time: 0}
				// 	continue
				// }
				// t := GetTimeInstalled(pkgName)

				fi, err := os.Stat(item.Path())
				if err != nil {
					timeChan <- _Time{Id: item.Id(), Time: 0}
					continue
				}
				fi.ModTime()
				t := fi.ModTime().Unix()
				timeChan <- _Time{Id: item.Id(), Time: t}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(timeChan)
	}()

	for t := range timeChan {
		infos[t.Id] = t.Time
	}

	return
}

func (self *ItemManager) GetAllNewInstalledApps() ([]ItemId, error) {
	ids := []ItemId{}
	f := glib.NewKeyFile()
	defer f.Free()

	configFile := ConfigFilePath(_NewSoftwareRecordFile)
	_, err := f.LoadFromFile(configFile, glib.KeyFileFlagsNone)
	if err != nil {
		return ids, err
	}

	_, newApps, _ := f.GetStringList(_NewSoftwareGroupName, _NewSoftwareKeyName)
	for _, id := range newApps {
		ids = append(ids, ItemId(id))
	}

	return ids, nil
}

func (self *ItemManager) MarkNew(_id ItemId) error {
	id := string(_id)
	configFile := ConfigFilePath(_NewSoftwareRecordFile)
	if !dutils.IsFileExist(configFile) {
		f, err := os.Create(configFile)
		if err != nil {
			return err
		}
		f.Close()
	}

	f := glib.NewKeyFile()
	defer f.Free()
	_, err := f.LoadFromFile(configFile, glib.KeyFileFlagsNone)
	if err != nil {
		return err
	}

	_, newApps, _ := f.GetStringList(_NewSoftwareGroupName, _NewSoftwareKeyName)
	if dutils.IsElementInList(id, newApps) {
		return fmt.Errorf("%q is already the new installed application", id)
	}

	newApps = append(newApps, id)
	f.SetStringList(_NewSoftwareGroupName, _NewSoftwareKeyName, newApps)
	return SaveKeyFile(f, configFile)
}

func (self *ItemManager) MarkLaunched(_id ItemId) error {
	id := string(_id)
	configFile := ConfigFilePath(_NewSoftwareRecordFile)

	f := glib.NewKeyFile()
	defer f.Free()
	_, err := f.LoadFromFile(configFile, glib.KeyFileFlagsNone)
	if err != nil {
		return err
	}

	_, newApps, _ := f.GetStringList(_NewSoftwareGroupName, _NewSoftwareKeyName)

	if !dutils.IsElementInList(id, newApps) {
		return fmt.Errorf("%q is already not the new installed application", id)
	}

	newIds := []string{}
	for _, newAppId := range newApps {
		if newAppId == id {
			continue
		}
		newIds = append(newIds, newAppId)
	}

	f.SetStringList(_NewSoftwareGroupName, _NewSoftwareKeyName, newIds)
	return SaveKeyFile(f, configFile)
}
