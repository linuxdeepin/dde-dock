package category

import (
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

// Manager for categories.
type Manager struct {
	categoryTable map[CategoryID]CategoryInfo
}

// NewManager creates a new category manager.
func NewManager() *Manager {
	m := &Manager{
		categoryTable: map[CategoryID]CategoryInfo{},
	}
	m.addCategory(
		&Info{AllID, AllCategoryName, map[ItemID]struct{}{}},
		&Info{OthersID, OtherCategoryName, map[ItemID]struct{}{}},
		&Info{NetworkID, NetworkCategoryName, map[ItemID]struct{}{}},
		&Info{MultimediaID, MultimediaCategoryName, map[ItemID]struct{}{}},
		&Info{GamesID, GamesCategoryName, map[ItemID]struct{}{}},
		&Info{GraphicsID, GraphicsCategoryName, map[ItemID]struct{}{}},
		&Info{ProductivityID, ProductivityCategoryName, map[ItemID]struct{}{}},
		&Info{IndustryID, IndustryCategoryName, map[ItemID]struct{}{}},
		&Info{EducationID, EducationCategoryName, map[ItemID]struct{}{}},
		&Info{DevelopmentID, DevelopmentCategoryName, map[ItemID]struct{}{}},
		&Info{SystemID, SystemCategoryName, map[ItemID]struct{}{}},
		&Info{UtilitiesID, UtilitiesCategoryName, map[ItemID]struct{}{}},
	)

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
