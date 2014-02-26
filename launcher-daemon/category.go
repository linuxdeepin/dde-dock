package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"dlib/glib-2.0"
)

type CategoryId int32

const (
	NetworkID CategoryId = iota
	MultimediaID
	GamesID
	GraphicsID
	ProductivityID
	IndustryID
	EducationID
	DevelopmentID
	SystemID
	UtilitiesID

	AllID   = -1
	OtherID = -2
	FavorID = -3

	DataDir      = "/usr/share/deepin-software-center/data"
	DataNewestId = DataDir + "/data_newest_id.ini"

	CategoryNameDBPath  = DataDir + "/update/%s/desktop/desktop.db"
	CategoryIndexDBPath = DataDir + "/update/%s/category/category.db"
)

type CategoryInfo struct {
	Id    CategoryId
	Name  string
	items map[ItemId]bool
}

var (
	nameIdMap     = map[string]CategoryId{}
	categoryTable = map[CategoryId]*CategoryInfo{
		AllID:   &CategoryInfo{AllID, "all", map[ItemId]bool{}},
		OtherID: &CategoryInfo{OtherID, "other", map[ItemId]bool{}},
	}
)

func getDBPath(template string) (string, error) {
	file := glib.NewKeyFile()
	defer file.Free()

	ok, err := file.LoadFromFile(DataNewestId, glib.KeyFileFlagsNone)
	if !ok {
		return "", err
	}

	id, err := file.GetString("newest", "data_id")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(template, id), nil
}

func initCategory() {
	dbPath, err := getDBPath(CategoryIndexDBPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := `select distinct first_category_name, first_category_index
	from  category_name;`

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var id CategoryId
		rows.Scan(&name, &id)
		// fmt.Println(name, id)
		lowerName := strings.ToLower(name)
		nameIdMap[lowerName] = id
		// fmt.Println("category id:", id)
		categoryTable[id] = &CategoryInfo{
			id,
			lowerName,
			map[ItemId]bool{},
		}
	}
	rows.Close()

	for _, v := range XCategoryNameIdMap {
		// fmt.Println(v.Name(), v.Id())
		nameIdMap[strings.ToLower(v.Name())] = v.Id()
	}
}

func findCategoryId(categoryName string) CategoryId {
	lowerCategoryName := strings.ToLower(categoryName)
	// fmt.Println("categoryName:", lowerCategoryName)
	id, ok := nameIdMap[lowerCategoryName]
	// fmt.Printf("nameIdMap[\"%s\"]=%d\n", lowerCategoryName, id)
	if !ok {
		return OtherID
	}
	return id
}

type CategoryInfoExport struct {
	Id    CategoryId
	Name  string
	Items []string
}

type CategoryInfosResult []CategoryInfoExport

func (res CategoryInfosResult) Len() int {
	return len(res)
}

func (res CategoryInfosResult) Swap(i, j int) {
	res[i], res[j] = res[j], res[i]
}

func (res CategoryInfosResult) Less(i, j int) bool {
	if res[i].Id == -1 || res[j].Id == -2 {
		return true
	} else if res[i].Id == -2 || res[j].Id == -1 {
		return false
	} else {
		return res[i].Id < res[j].Id
	}
}

type ItemIdList []string

func (l ItemIdList) Len() int {
	return len(l)
}

func (l ItemIdList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l ItemIdList) Less(i, j int) bool {
	return itemTable[ItemId(l[i])].Name < itemTable[ItemId(l[j])].Name
}

func getCategoryInfos() CategoryInfosResult {
	infos := make(CategoryInfosResult, 0)
	for _, v := range categoryTable {
		if v.Id == AllID {
			continue
		}
		info := CategoryInfoExport{v.Id, v.Name, make([]string, 0)}
		for k, _ := range v.items {
			info.Items = append(info.Items, string(k))
		}
		sort.Sort(ItemIdList(info.Items))
		infos = append(infos, info)
	}
	sort.Sort(infos)
	return infos
}
