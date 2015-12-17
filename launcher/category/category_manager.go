package category

import (
	"errors"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"gir/gio-2.0"
)

type QueryIDTransition interface {
	Query(*gio.DesktopAppInfo) (CategoryID, error)
	Free()
}

// Manager for categories.
type Manager struct {
	categoryTable              map[CategoryID]CategoryInfo
	deepinQueryIDTransition    QueryIDTransition
	xCategoryQueryIDTransition QueryIDTransition
}

// NewManager creates a new category manager.
func NewManager(categories []CategoryInfo) *Manager {
	m := &Manager{
		categoryTable: map[CategoryID]CategoryInfo{},
	}
	m.AddCategory(categories...)
	return m
}

func (m *Manager) AddCategory(c ...CategoryInfo) {
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
	if m.deepinQueryIDTransition != nil {
		cid, e := m.deepinQueryIDTransition.Query(app)
		if e == nil {
			return cid, nil
		}
		err = e
	}

	if m.xCategoryQueryIDTransition != nil {
		return m.xCategoryQueryIDTransition.Query(app)
	}

	if err != nil {
		return OthersID, err
	}

	return OthersID, errors.New("No QueryIDTransition is created")
}

func (m *Manager) LoadAppCategoryInfo(deepin string, xcategory string) error {
	var err error
	m.deepinQueryIDTransition, err = NewDeepinQueryIDTransition(deepin)
	if err != nil {
		return err
	}

	m.xCategoryQueryIDTransition, err = NewXCategoryQueryIDTransition(xcategory)
	return err
}

func (m *Manager) FreeAppCategoryInfo() {
	m.deepinQueryIDTransition.Free()
	m.xCategoryQueryIDTransition.Free()
	m.deepinQueryIDTransition = nil
	m.xCategoryQueryIDTransition = nil
}
