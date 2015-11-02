package category

import (
	"errors"
	"gir/gio-2.0"
	"path"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

type QueryIDTransaction interface {
	Query(string) (CategoryID, error)
	Free()
}

// Manager for categories.
type Manager struct {
	store                       DStore
	categoryTable               map[CategoryID]CategoryInfo
	deepinQueryIDTransaction    QueryIDTransaction
	xCategoryQueryIDTransaction QueryIDTransaction
	queryPkgNameTransaction     QueryPkgNameTransaction
}

// NewManager creates a new category manager.
func NewManager(store DStore, categories []CategoryInfo) *Manager {
	m := &Manager{
		store:                       store,
		categoryTable:               map[CategoryID]CategoryInfo{},
		deepinQueryIDTransaction:    nil,
		xCategoryQueryIDTransaction: nil,
		queryPkgNameTransaction:     nil,
	}
	m.addCategory(categories...)
	m.addCategory(NewInfo(AllID, AllName))
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
	m.categoryTable[cid].AddItem(id)
	m.categoryTable[AllID].AddItem(id)
}

// RemoveItem removes a app from category.
func (m *Manager) RemoveItem(id ItemID, cid CategoryID) {
	m.categoryTable[cid].RemoveItem(id)
	m.categoryTable[AllID].RemoveItem(id)
}

func (m *Manager) QueryID(app *gio.DesktopAppInfo) (CategoryID, error) {
	var err error
	if m.queryPkgNameTransaction != nil && m.deepinQueryIDTransaction != nil {
		desktopID := path.Base(app.GetFilename())
		pkgName := m.queryPkgNameTransaction.Query(desktopID)
		cid, e := m.deepinQueryIDTransaction.Query(pkgName)
		if e != nil {
			err = e
		}

		if cid != OthersID {
			return cid, nil
		}
	}

	if m.xCategoryQueryIDTransaction != nil {
		return m.xCategoryQueryIDTransaction.Query(app.GetCategories())
	}

	if err != nil {
		return OthersID, err
	}

	return OthersID, errors.New("No QueryIDTransaction is created or QueryIDTransaction failed")
}

func (m *Manager) LoadAppCategoryInfo(files ...string) error {
	m.FreeAppCategoryInfo()

	var err1 error
	m.queryPkgNameTransaction, err1 = m.store.NewQueryPkgNameTransaction(files[0])

	var err2 error
	m.deepinQueryIDTransaction, err2 = NewDeepinQueryIDTransaction(files[1])

	var err3 error
	m.xCategoryQueryIDTransaction, err3 = NewXCategoryQueryIDTransaction(files[2])

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}
	return nil
}

func (m *Manager) FreeAppCategoryInfo() {
	if m.deepinQueryIDTransaction != nil {
		m.deepinQueryIDTransaction.Free()
		m.deepinQueryIDTransaction = nil
	}
	if m.queryPkgNameTransaction != nil {
		m.queryPkgNameTransaction = nil
	}
	if m.xCategoryQueryIDTransaction != nil {
		m.xCategoryQueryIDTransaction.Free()
		m.xCategoryQueryIDTransaction = nil
	}
}
