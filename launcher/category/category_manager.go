package category

import (
	"encoding/json"
	"errors"
	"gir/gio-2.0"
	"os"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

type CategoryJSONInfo struct {
	ID      string
	Locales map[string]map[string]string
	Name    string
}

func GetAllInfos(file string) []CategoryInfo {
	fallbackCategories := []CategoryInfo{
		NewInfo(OthersID, OthersName),
		NewInfo(InternetID, InternetName),
		NewInfo(OfficeID, OfficeName),
		NewInfo(DevelopmentID, DevelopmentName),
		NewInfo(ReadingID, ReadingName),
		NewInfo(GraphicsID, GraphicsName),
		NewInfo(GameID, GameName),
		NewInfo(MusicID, MusicName),
		NewInfo(SystemID, SystemName),
		NewInfo(VideoID, VideoName),
		NewInfo(ChatID, ChatName),
	}
	var categoryInfos []CategoryInfo
	f, err := os.Open(file)
	if err != nil {
		return fallbackCategories
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	var jsonInfo []CategoryJSONInfo
	if err := decoder.Decode(&jsonInfo); err != nil {
		return fallbackCategories
	}

	categoryInfos = make([]CategoryInfo, len(jsonInfo))
	for i, info := range jsonInfo {
		cid, _ := getCategoryID(info.ID)
		categoryInfos[i] = NewInfo(cid, info.Name)
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
