package category

import (
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
)

type CategoryManager struct {
	categoryTable map[CategoryId]CategoryInfoInterface
}

func NewCategoryManager() *CategoryManager {
	m := &CategoryManager{
		categoryTable: map[CategoryId]CategoryInfoInterface{},
	}
	m.addCategory(
		&CategoryInfo{AllID, AllCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{OtherID, OtherCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{NetworkID, NetworkCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{MultimediaID, MultimediaCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{GamesID, GamesCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{GraphicsID, GraphicsCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{ProductivityID, ProductivityCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{IndustryID, IndustryCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{EducationID, EducationCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{DevelopmentID, DevelopmentCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{SystemID, SystemCategoryName, map[ItemId]struct{}{}},
		&CategoryInfo{UtilitiesID, UtilitiesCategoryName, map[ItemId]struct{}{}},
	)

	return m
}

func (m *CategoryManager) addCategory(c ...CategoryInfoInterface) {
	for _, info := range c {
		m.categoryTable[info.Id()] = info
	}
}

func (m *CategoryManager) GetCategory(id CategoryId) CategoryInfoInterface {
	category, ok := m.categoryTable[id]
	if ok {
		return category
	}

	return nil
}

func (m *CategoryManager) GetAllCategory() []CategoryId {
	ids := []CategoryId{}
	for id, _ := range m.categoryTable {
		ids = append(ids, id)
	}

	return ids
}

func (m *CategoryManager) AddItem(id ItemId, cid CategoryId) {
	m.categoryTable[cid].AddItem(id)
	m.categoryTable[AllID].AddItem(id)
}

func (m *CategoryManager) RemoveItem(id ItemId, cid CategoryId) {
	m.categoryTable[cid].RemoveItem(id)
	m.categoryTable[AllID].RemoveItem(id)
}
