package launcher

import (
	"fmt"
	"sort"
	"strings"

	"pkg.linuxdeepin.com/lib/glib-2.0"
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

	AllID    = -1
	OthersID = -2
	FavorID  = -3

	_AllCategoryName          = "all"
	_OthersCategoryName       = "others"
	_InternetCategoryName     = "internet"
	_MultimediaCategoryName   = "multimedia"
	_GamesCategoryName        = "games"
	_GraphicsCategoryName     = "graphics"
	_ProductivityCategoryName = "productivity"
	_IndustryCategoryName     = "industry"
	_EducationCategoryName    = "education"
	_DevelopmentCategoryName  = "development"
	_SystemCategoryName       = "system"
	_UtilitiesCategoryName    = "utilities"

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
		_AllCategoryName:          AllID,
		_OthersCategoryName:       OthersID,
		_InternetCategoryName:     NetworkID,
		_MultimediaCategoryName:   MultimediaID,
		_GamesCategoryName:        GamesID,
		_GraphicsCategoryName:     GraphicsID,
		_ProductivityCategoryName: ProductivityID,
		_IndustryCategoryName:     IndustryID,
		_EducationCategoryName:    EducationID,
		_DevelopmentCategoryName:  DevelopmentID,
		_SystemCategoryName:       SystemID,
		_UtilitiesCategoryName:    UtilitiesID,
	}
	categoryTable = map[CategoryId]*CategoryInfo{
		AllID:          &CategoryInfo{AllID, _AllCategoryName, map[ItemId]bool{}},
		OthersID:       &CategoryInfo{OthersID, _OthersCategoryName, map[ItemId]bool{}},
		NetworkID:      &CategoryInfo{NetworkID, _InternetCategoryName, map[ItemId]bool{}},
		MultimediaID:   &CategoryInfo{MultimediaID, _MultimediaCategoryName, map[ItemId]bool{}},
		GamesID:        &CategoryInfo{GamesID, _GamesCategoryName, map[ItemId]bool{}},
		GraphicsID:     &CategoryInfo{GraphicsID, _GraphicsCategoryName, map[ItemId]bool{}},
		ProductivityID: &CategoryInfo{ProductivityID, _ProductivityCategoryName, map[ItemId]bool{}},
		IndustryID:     &CategoryInfo{IndustryID, _IndustryCategoryName, map[ItemId]bool{}},
		EducationID:    &CategoryInfo{EducationID, _EducationCategoryName, map[ItemId]bool{}},
		DevelopmentID:  &CategoryInfo{DevelopmentID, _DevelopmentCategoryName, map[ItemId]bool{}},
		SystemID:       &CategoryInfo{SystemID, _SystemCategoryName, map[ItemId]bool{}},
		UtilitiesID:    &CategoryInfo{UtilitiesID, _UtilitiesCategoryName, map[ItemId]bool{}},
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
			return OthersID
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
