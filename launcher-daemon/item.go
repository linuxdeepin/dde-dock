package main

import (
	// "crypto/md5"
	"database/sql"
	// "fmt"
	"log"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	pinyin "dbus/com/deepin/api/search"
	"dlib/gio-2.0"
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
	Path        string
	Name        string
	enName      string
	Id          ItemId
	Icon        string
	categoryIds map[CategoryId]bool
	xinfo       Xinfo
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

	i.categoryIds = map[CategoryId]bool{}
	i.xinfo.keywords = make([]string, 0)
	keywords := app.GetKeywords()
	for _, keyword := range keywords {
		i.xinfo.keywords = append(i.xinfo.keywords, strings.ToLower(keyword))
	}
	i.xinfo.exec = app.GetExecutable()
	i.xinfo.genericName = app.GetGenericName()
	i.xinfo.description = app.GetDescription()
	for _, id := range getCategories(app) {
		i.categoryIds[id] = true
		categoryTable[id].items[i.Id] = true
	}
	categoryTable[AllID].items[i.Id] = true
	itemTable[i.Id] = i
}

func (i *ItemInfo) getCategoryIds() []CategoryId {
	ids := make([]CategoryId, 0)
	for k, _ := range i.categoryIds {
		ids = append(ids, k)
	}
	return ids
}

func (i *ItemInfo) destroy() {
	for _, cid := range i.getCategoryIds() {
		// fmt.Printf("delete id from category#%d\n", cid)
		delete(categoryTable[cid].items, i.Id)
	}
	// logger.Info("delete id from category#-1")
	delete(categoryTable[OtherID].items, i.Id)
}

func getDeepinCategory(app *gio.DesktopAppInfo) (CategoryId, error) {
	filename := app.GetFilename()
	basename := path.Base(filename)
	dbPath, err := getDBPath(CategoryNameDBPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
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
	id := findCategoryId(categoryName)
	// logger.Info(categoryName, id)
	return id, nil
}

func getXCategories(app *gio.DesktopAppInfo) []CategoryId {
	candidateIds := map[CategoryId]bool{}
	categories := strings.Split(app.GetCategories(), ";")
	for _, category := range categories {
		candidateIds[findCategoryId(category)] = true
	}

	if len(candidateIds) > 1 && candidateIds[OtherID] {
		delete(candidateIds, OtherID)
	}

	ids := make([]CategoryId, 0)
	for id := range candidateIds {
		ids = append(ids, id)
	}

	return ids
}

func getCategories(app *gio.DesktopAppInfo) []CategoryId {
	var categories = make([]CategoryId, 0)
	id, err := getDeepinCategory(app)
	if err != nil {
		return getXCategories(app)
	} else {
		return append(categories, id)
	}
}

func genId(filename string) ItemId {
	basename := path.Base(filename)
	// return ItemId(fmt.Sprintf("%x", md5.Sum([]byte(basename))))
	return ItemId(basename[:len(basename)-8]) // len(".desktop")
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
	tree, err = pinyin.NewSearch("com.deepin.api.Search", "/com/deepin/api/Search")
	if err != nil {
		return
	}
	names := make(map[string]string, 0)
	for _, v := range itemTable {
		names[v.Name] = v.Name
	}
	treeId, _ = tree.NewTrieWithString(names, "DDELauncherDaemon")
}

func getItemInfos(id CategoryId) []ItemInfo {
	// logger.Info(id)
	infos := make([]ItemInfo, 0)
	if _, ok := categoryTable[id]; !ok {
		logger.Info("category id:", id, "not exist")
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
