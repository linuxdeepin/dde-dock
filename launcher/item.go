package launcher

import (
	// "crypto/md5"
	"database/sql"
	"errors"
	// "fmt"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"pkg.linuxdeepin.com/lib/gio-2.0"
)

type ItemId string

type Xinfo struct {
	keywords    []string
	exec        string
	genericName string
	description string
	// #define FILENAME_WEIGHT 0.3
	// #define GENERIC_NAME_WEIGHT 0.01
	// #define KEYWORD_WEIGHT 0.1
	// #define CATEGORY_WEIGHT 0.01
	// #define NAME_WEIGHT 0.01
	// #define DISPLAY_NAME_WEIGHT 0.1
	// #define DESCRIPTION_WEIGHT 0.01
	// #define EXECUTABLE_WEIGHT 0.05
}

type ItemInfo struct {
	Path       string
	Name       string
	enName     string
	Id         ItemId
	Icon       string
	categoryId CategoryId
	xinfo      Xinfo
}

// TODO: add some method to ItemTable like remove/add
// type ItemTable map[ItemId]*ItemId

var itemTable = map[ItemId]*ItemInfo{}

func (i *ItemInfo) init(app *gio.DesktopAppInfo) {
	i.Id = getId(app)
	i.Path = app.GetFilename()
	i.Name = app.GetDisplayName()
	i.enName = app.GetString("Name")
	icon := app.GetIcon()
	if icon != nil {
		i.Icon = icon.ToString()
		if path.IsAbs(i.Icon) && !exist(i.Icon) {
			i.Icon = ""
		}
	}

	i.xinfo.keywords = make([]string, 0)
	keywords := app.GetKeywords()
	for _, keyword := range keywords {
		i.xinfo.keywords = append(i.xinfo.keywords, strings.ToLower(keyword))
	}
	i.xinfo.exec = app.GetExecutable()
	i.xinfo.genericName = app.GetGenericName()
	i.xinfo.description = app.GetDescription()
	i.categoryId = getCategory(app)
	categoryTable[i.categoryId].items[i.Id] = true
	categoryTable[AllID].items[i.Id] = true
	itemTable[i.Id] = i
}

func (i *ItemInfo) getCategoryId() CategoryId {
	return i.categoryId
}

func (i *ItemInfo) destroy() {
	// fmt.Printf("delete id from category#%d\n", cid)
	delete(categoryTable[i.getCategoryId()].items, i.Id)
	// logger.Info("delete id from category#-1")
	delete(categoryTable[OtherID].items, i.Id)
}

func getDeepinCategory(app *gio.DesktopAppInfo) (CategoryId, error) {
	filename := app.GetFilename()
	basename := path.Base(filename)
	dbPath, err := getDBPath(CategoryNameDBPath)
	if err != nil {
		return OtherID, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return OtherID, err
	}
	defer db.Close()

	var categoryName string
	err = db.QueryRow(
		`select first_category_name
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

	logger.Debug("get category name for", basename, "is", categoryName)
	id := findCategoryId(categoryName)
	logger.Debug(categoryName, id)
	return id, nil
}

func getXCategory(app *gio.DesktopAppInfo) CategoryId {
	candidateIds := map[CategoryId]bool{}
	categories := strings.Split(app.GetCategories(), ";")
	for _, category := range categories {
		if id, ok := nameIdMap[category]; ok {
			candidateIds[id] = true
		}
	}

	if len(candidateIds) == 0 {
		for _, category := range categories {
			if _, ok := nameIdMap[category]; !ok {
				candidateIds[findCategoryId(category)] = true
			}
		}
	}

	if len(candidateIds) > 1 && candidateIds[OtherID] {
		delete(candidateIds, OtherID)
	}

	ids := make([]CategoryId, 0)
	for id := range candidateIds {
		ids = append(ids, id)
	}

	return ids[0]
}

func getCategory(app *gio.DesktopAppInfo) CategoryId {
	id, err := getDeepinCategory(app)
	if err != nil {
		logger.Warningf("\"%s\" get category from database failed: %s", app.GetDisplayName(), err)
		return getXCategory(app)
	}
	logger.Debug("get category from database:", id)
	return id
}

func genId(filename string) ItemId {
	basename := path.Base(filename)
	// return ItemId(fmt.Sprintf("%x", md5.Sum([]byte(basename))))
	return ItemId(strings.Replace(basename[:len(basename)-8], "_", "-", -1)) // len(".desktop")
}

func getId(app *gio.DesktopAppInfo) ItemId {
	return genId(app.GetFilename())
}

func initItems() {
	allApps := gio.AppInfoGetAll()

	for _, app := range allApps {
		desktopApp := gio.ToDesktopAppInfo(app)
		// TODO: get keywords for pinyin searching.
		if app.ShouldShow() {
			itemInfo := &ItemInfo{}
			itemInfo.init(desktopApp)
		}
		app.Unref()
	}

	var err error
	names := make([]string, 0)
	for _, v := range itemTable {
		names = append(names, v.Name)
	}
	logger.Debug("Names:", names)
	pinyinSearchObj, err = NewPinYinSearch(names)
	if err != nil {
		logger.Warning("build pinyin search object failed:", err)
	}
}

func getItemInfos(id CategoryId) []ItemInfo {
	// logger.Info(id)
	infos := make([]ItemInfo, 0)
	if _, ok := categoryTable[id]; !ok {
		logger.Warning("category id:", id, "not exist")
		return infos
	}

	for k, _ := range categoryTable[id].items {
		// logger.Info("get item", k, "from category#", id)
		if _, ok := itemTable[k]; ok {
			infos = append(infos, *itemTable[k])
		}
	}

	return infos
}
