package category

import (
	"errors"
	"gir/gio-2.0"
	"pkg.deepin.io/dde/daemon/dstore"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

func GetAllInfos(file string) []CategoryInfo {
	dstoreCategoryInfos := dstore.GetAllInfos(file)
	categoryInfos := make([]CategoryInfo, len(dstoreCategoryInfos))
	for i, categoryInfo := range dstoreCategoryInfos {
		cid, _ := getCategoryID(categoryInfo.ID)
		categoryInfos[i] = NewInfo(cid, categoryInfo.Name)
	}
	return categoryInfos
}

// Manager for categories.
type Manager struct {
	store              DStore
	categoryTable      map[CategoryID]CategoryInfo
	queryIDTransaction QueryCategoryTransaction
}

// NewManager creates a new category manager.
func NewManager(store DStore, categories []CategoryInfo) *Manager {
	m := &Manager{
		store:              store,
		categoryTable:      map[CategoryID]CategoryInfo{},
		queryIDTransaction: nil,
	}
	m.addCategory(categories...)
	m.addCategory(NewInfo(AllID, dstore.AllName))
	return m
}

func (m *Manager) addCategory(c ...CategoryInfo) {
	for _, info := range c {
		m.categoryTable[info.ID()] = info
	}
}

// GetCategory returns category info according to id.
func (m *Manager) GetCategory(id CategoryID) CategoryInfo {
	category, ok := m.categoryTable[id]
	if ok {
		return category
	}

	return nil
}

// GetAllCategory returns all categories.
func (m *Manager) GetAllCategory() []CategoryID {
	ids := []CategoryID{}
	for id := range m.categoryTable {
		ids = append(ids, id)
	}

	return ids
}

// AddItem adds a app to category.
func (m *Manager) AddItem(id ItemID, cid CategoryID) {
	if category, ok := m.categoryTable[cid]; ok {
		category.AddItem(id)
	}
	m.categoryTable[AllID].AddItem(id)
}

// RemoveItem removes a app from category.
func (m *Manager) RemoveItem(id ItemID, cid CategoryID) {
	if category, ok := m.categoryTable[cid]; ok {
		category.RemoveItem(id)
	}
	m.categoryTable[AllID].RemoveItem(id)
}

func (m *Manager) QueryID(app *gio.DesktopAppInfo) (CategoryID, error) {
	if m.queryIDTransaction != nil {
		c, e := m.queryIDTransaction.Query(app)
		if e != nil {
			return OthersID, e
		}
		cid, e := getCategoryID(c)
		return cid, e
	}

	return OthersID, errors.New("No QueryIDTransaction is created or QueryIDTransaction failed")
}

func (m *Manager) LoadCategoryInfo() error {
	m.FreeAppCategoryInfo()

	var err error
	m.queryIDTransaction, err = m.store.NewQueryCategoryTransaction()

	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) FreeAppCategoryInfo() {
	if m.queryIDTransaction != nil {
		m.queryIDTransaction = nil
	}
}
