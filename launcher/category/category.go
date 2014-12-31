package category

import (
	"database/sql"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"

	lerrors "pkg.linuxdeepin.com/dde-daemon/launcher/errors"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/glib-2.0"
)

const (
	UnknownID CategoryId = iota - 3
	OtherID
	AllID
	NetworkID
	MultimediaID
	GamesID
	GraphicsID
	ProductivityID
	IndustryID
	EducationID
	DevelopmentID
	SystemID
	UtilitiesID

	AllCategoryName          = "all"
	OtherCategoryName        = "others"
	NetworkCategoryName      = "internet"
	MultimediaCategoryName   = "multimedia"
	GamesCategoryName        = "games"
	GraphicsCategoryName     = "graphics"
	ProductivityCategoryName = "productivity"
	IndustryCategoryName     = "industry"
	EducationCategoryName    = "education"
	DevelopmentCategoryName  = "development"
	SystemCategoryName       = "system"
	UtilitiesCategoryName    = "utilities"

	SoftwareCenterDataDir = "/usr/share/deepin-software-center/data"
	_DataNewestIdFileName = "data_newest_id.ini"
	CategoryNameDBPath    = "/update/%s/desktop/desktop2014.db"
)

var (
	categoryNameTable = map[string]CategoryId{
		OtherCategoryName:        OtherID,
		AllCategoryName:          AllID,
		NetworkCategoryName:      NetworkID,
		MultimediaCategoryName:   MultimediaID,
		GamesCategoryName:        GamesID,
		GraphicsCategoryName:     GraphicsID,
		ProductivityCategoryName: ProductivityID,
		IndustryCategoryName:     IndustryID,
		EducationCategoryName:    EducationID,
		DevelopmentCategoryName:  DevelopmentID,
		SystemCategoryName:       SystemID,
		UtilitiesCategoryName:    UtilitiesID,
	}
)

type CategoryInfo struct {
	id    CategoryId
	name  string
	items map[ItemId]struct{}
}

func (c *CategoryInfo) Id() CategoryId {
	return c.id
}

func (c *CategoryInfo) Name() string {
	return c.name
}

func (c *CategoryInfo) AddItem(itemId ItemId) {
	c.items[itemId] = struct{}{}
}
func (c *CategoryInfo) RemoveItem(itemId ItemId) {
	delete(c.items, itemId)
}

func (c *CategoryInfo) Items() []ItemId {
	items := []ItemId{}
	for itemId, _ := range c.items {
		items = append(items, itemId)
	}
	return items
}

func getNewestDataId(dataDir string) (string, error) {
	file := glib.NewKeyFile()
	defer file.Free()

	ok, err := file.LoadFromFile(path.Join(dataDir, _DataNewestIdFileName), glib.KeyFileFlagsNone)
	if !ok {
		return "", err
	}

	id, err := file.GetString("newest", "data_id")
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetDBPath(dataDir string, template string) (string, error) {
	id, err := getNewestDataId(dataDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, fmt.Sprintf(template, id)), nil
}

func QueryCategoryId(app *gio.DesktopAppInfo, db *sql.DB) (CategoryId, error) {
	if app == nil {
		return UnknownID, lerrors.NilArgument
	}

	filename := app.GetFilename()
	basename := path.Base(filename)
	id, err := getDeepinCategory(basename, db)
	if err != nil {
		categories := strings.Split(strings.TrimRight(app.GetCategories(), ";"), ";")
		return getXCategory(categories), nil
	}
	return id, nil
}

func getDeepinCategory(basename string, db *sql.DB) (CategoryId, error) {
	if db == nil {
		return UnknownID, errors.New("invalid db")
	}

	var categoryName string
	err := db.QueryRow(`
	select first_category_name
	from desktop
	where desktop_name = ?`,
		basename,
	).Scan(&categoryName)
	if err != nil {
		return OtherID, err
	}

	if categoryName == "" {
		return OtherID, errors.New("get empty category")
	}

	return getCategoryId(categoryName)
}

type CategoryIdList []CategoryId

func (self CategoryIdList) Less(i, j int) bool {
	return self[i] < self[j]
}

func (self CategoryIdList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self CategoryIdList) Len() int {
	return len(self)
}

func getXCategory(categories []string) CategoryId {
	candidateIds := map[CategoryId]bool{OtherID: true}
	for _, category := range categories {
		if id, err := getCategoryId(category); err == nil {
			candidateIds[id] = true
		}
	}

	if len(candidateIds) > 1 && candidateIds[OtherID] {
		delete(candidateIds, OtherID)
	}

	ids := make([]CategoryId, 0)
	for id := range candidateIds {
		ids = append(ids, id)
	}

	sort.Sort(CategoryIdList(ids))

	return ids[0]
}

func getCategoryId(name string) (CategoryId, error) {
	name = strings.ToLower(name)
	if id, ok := categoryNameTable[name]; ok {
		return id, nil
	}

	if id, ok := xCategoryNameIdMap[name]; ok {
		return id, nil
	}

	if id, ok := extraXCategoryNameIdMap[name]; ok {
		return id, nil
	}

	return UnknownID, errors.New("unknown id")
}
