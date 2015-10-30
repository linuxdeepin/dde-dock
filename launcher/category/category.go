package category

import (
	"database/sql"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"

	lerrors "pkg.deepin.io/dde/daemon/launcher/errors"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
)

// category id and name.
const (
	UnknownID CategoryID = iota - 3
	OthersID
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
	_DataNewestIDFileName = "data_newest_id.ini"
	CategoryNameDBPath    = "/update/%s/desktop/desktop2014.db"
)

var (
	categoryNameTable = map[string]CategoryID{
		OtherCategoryName:        OthersID,
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

// Info for category.
type Info struct {
	id    CategoryID
	name  string
	items map[ItemID]struct{}
}

// ID returns category id.
func (c *Info) ID() CategoryID {
	return c.id
}

// Name returns category english name.
func (c *Info) Name() string {
	return c.name
}

// LocaleName returns category's locale name.
func (c *Info) LocaleName() string {
	return gettext.Tr(c.name)
}

// AddItem adds a new app.
func (c *Info) AddItem(itemID ItemID) {
	c.items[itemID] = struct{}{}
}

// RemoveItem removes a app.
func (c *Info) RemoveItem(itemID ItemID) {
	delete(c.items, itemID)
}

// Items returns all items belongs to this category.
func (c *Info) Items() []ItemID {
	items := []ItemID{}
	for itemID := range c.items {
		items = append(items, itemID)
	}
	return items
}

func getNewestDataID(dataDir string) (string, error) {
	file := glib.NewKeyFile()
	defer file.Free()

	ok, err := file.LoadFromFile(path.Join(dataDir, _DataNewestIDFileName), glib.KeyFileFlagsNone)
	if !ok {
		return "", err
	}

	id, err := file.GetString("newest", "data_id")
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetDBPath returns db path store category info.
func GetDBPath(dataDir string, template string) (string, error) {
	id, err := getNewestDataID(dataDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, fmt.Sprintf(template, id)), nil
}

// QueryID returns app's category.
func QueryID(app *gio.DesktopAppInfo, db *sql.DB) (CategoryID, error) {
	if app == nil {
		return OthersID, lerrors.NilArgument
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

func getDeepinCategory(basename string, db *sql.DB) (CategoryID, error) {
	if db == nil {
		return OthersID, errors.New("invalid db")
	}

	var categoryName string
	err := db.QueryRow(`
	select first_category_name
	from desktop
	where desktop_name = ?`,
		basename,
	).Scan(&categoryName)
	if err != nil {
		return OthersID, err
	}

	if categoryName == "" {
		return OthersID, errors.New("get empty category")
	}

	return getCategoryID(categoryName)
}

// IDList type alias for []CategoryID, used for sorting.
type IDList []CategoryID

func (list IDList) Less(i, j int) bool {
	return list[i] < list[j]
}

func (list IDList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list IDList) Len() int {
	return len(list)
}

func getXCategory(categories []string) CategoryID {
	candidateIDs := map[CategoryID]bool{OthersID: true}
	for _, category := range categories {
		if id, err := getCategoryID(category); err == nil {
			candidateIDs[id] = true
		}
	}

	if len(candidateIDs) > 1 && candidateIDs[OthersID] {
		delete(candidateIDs, OthersID)
	}

	var ids []CategoryID
	for id := range candidateIDs {
		ids = append(ids, id)
	}

	sort.Sort(IDList(ids))

	return ids[0]
}

func getCategoryID(name string) (CategoryID, error) {
	name = strings.ToLower(name)
	if id, ok := categoryNameTable[name]; ok {
		return id, nil
	}

	if id, ok := xCategoryNameIDMap[name]; ok {
		return id, nil
	}

	if id, ok := extraXCategoryNameIDMap[name]; ok {
		return id, nil
	}

	return OthersID, errors.New("unknown id")
}
