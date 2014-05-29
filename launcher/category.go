package launcher

import (
	"fmt"
	"sort"
	"strings"

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

	CategoryNameDBPath  = DataDir + "/update/%s/desktop/desktop2014.db"
	CategoryIndexDBPath = DataDir + "/update/%s/category/category.db"
)

type CategoryInfo struct {
	Id    CategoryId
	Name  string
	items map[ItemId]bool
}

var (
	nameIdMap = map[string]CategoryId{
		"all":          AllID,
		"other":        OtherID,
		"internet":     NetworkID,
		"multimedia":   MultimediaID,
		"games":        GamesID,
		"graphics":     GraphicsID,
		"productivity": ProductivityID,
		"industry":     IndustryID,
		"education":    EducationID,
		"development":  DevelopmentID,
		"system":       SystemID,
		"utilities":    UtilitiesID,
	}
	categoryTable = map[CategoryId]*CategoryInfo{
		AllID:          &CategoryInfo{AllID, "all", map[ItemId]bool{}},
		OtherID:        &CategoryInfo{OtherID, "other", map[ItemId]bool{}},
		NetworkID:      &CategoryInfo{NetworkID, "internet", map[ItemId]bool{}},
		MultimediaID:   &CategoryInfo{MultimediaID, "multimedia", map[ItemId]bool{}},
		GamesID:        &CategoryInfo{GamesID, "games", map[ItemId]bool{}},
		GraphicsID:     &CategoryInfo{GraphicsID, "graphics", map[ItemId]bool{}},
		ProductivityID: &CategoryInfo{ProductivityID, "productivity", map[ItemId]bool{}},
		IndustryID:     &CategoryInfo{IndustryID, "industry", map[ItemId]bool{}},
		EducationID:    &CategoryInfo{EducationID, "education", map[ItemId]bool{}},
		DevelopmentID:  &CategoryInfo{DevelopmentID, "development", map[ItemId]bool{}},
		SystemID:       &CategoryInfo{SystemID, "system", map[ItemId]bool{}},
		UtilitiesID:    &CategoryInfo{UtilitiesID, "utilities", map[ItemId]bool{}},
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
	for k, id := range XCategoryNameIdMap {
		// logger.Info(k, id)
		if _, ok := nameIdMap[k]; !ok {
			nameIdMap[k] = id
		}
	}
}

func findCategoryId(categoryName string) CategoryId {
	lowerCategoryName := strings.ToLower(categoryName)
	logger.Debug("categoryName:", lowerCategoryName)
	id, ok := nameIdMap[lowerCategoryName]
	logger.Debug("nameIdMap[\"%s\"]=%d\n", lowerCategoryName, id)
	if !ok {
		id, ok = extraXCategoryNameIdMap[lowerCategoryName]
		if !ok {
			return OtherID
		}
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
