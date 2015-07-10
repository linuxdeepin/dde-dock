package desktop

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
	"sort"
	"strings"
)

var _ = sort.Sort

const (
	// AppGroupPrefix is the prefix of AppGroup's name.
	AppGroupPrefix string = ".deepin_rich_dir_"
)

// AppGroup represents appgroup on desktop.
type AppGroup struct {
	*Item
}

// NewAppGroup creates a app group.
func NewAppGroup(app *Application, uris []string) *AppGroup {
	return &AppGroup{NewItem(app, uris)}
}

// GenMenuContent generates json format menu content used in DeepinMenu for AppGroup.
func (item *AppGroup) GenMenuContent() (*Menu, error) {
	item.menu = NewMenu()
	return item.menu.AppendItem(NewMenuItem(Tr("_Open"), func() {
		item.app.emitRequestOpen(item.uris)
	}, true)).AddSeparator().AppendItem(NewMenuItem(Tr("_Rename"), func() {
		item.emitRequestRename()
	}, !item.multiple)).AddSeparator().AppendItem(NewMenuItem(Tr("_Ungroup"), func() {
		// TODO
		// item.app.emitRequestDismissAppGroup(item.uri)
	}, true)).AddSeparator().AppendItem(NewMenuItem(Tr("_Delete"), func() {
		item.emitRequestDelete()
	}, true)), nil
}

func getCategoryFromSoftwareCenter(db *sql.DB, file string) (string, error) {
	var category string
	err := db.QueryRow(`select first_category_name from desktop where desktop_name = ?`, filepath.Base(file)).Scan(&category)
	if err != nil {
		return "", err
	}

	return category, nil
}

func getGroupNameFromSoftwareCenter(files []string) string {
	db, err := sql.Open("sqlite3", "dbPath")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer db.Close()

	category, err := getCategoryFromSoftwareCenter(db, files[0])
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if category == "" {
		return ""
	}

	for _, file := range files[1:] {
		anotherCategory, err := getCategoryFromSoftwareCenter(db, file)
		if err != nil || anotherCategory != category {
			return ""
		}
	}

	return category
}

func containsWithCaseInsensitive(m map[string]struct{}, category string) bool {
	category = strings.ToLower(category)
	_, ok := m[category]
	return ok
}

func getCategoriesFromDesktop(file string, invalidCategories map[string]struct{}) (categories []string) {
	app := gio.NewDesktopAppInfoFromFilename(file)
	if app != nil {
		return
	}
	defer app.Unref()

	for _, category := range strings.Split(app.GetCategories(), ";") {
		if category != "" && !containsWithCaseInsensitive(invalidCategories, category) {
			categories = append(categories, category)
		}
	}

	return
}

type CategoryInfo struct {
	Name          string
	LowerCaseName string
	Count         int32
}

type CategoryInfos []CategoryInfo

func (info CategoryInfos) Less(i, j int) bool {
	if info[i].Count < info[j].Count {
		return true
	} else if info[i].Count > info[j].Count {
		return false
	} else {
		return info[i].Name < info[j].Name
	}
}

func (info CategoryInfos) Swap(i, j int) {
	info[i], info[j] = info[j], info[i]
}

func (info CategoryInfos) Len() int {
	return len(info)
}

func countCategories(files []string, invalidCategories map[string]struct{}) []CategoryInfo {
	count := []CategoryInfo{}
	m := map[string]*CategoryInfo{}
	for _, file := range files {
		categories := getCategoriesFromDesktop(file, invalidCategories)
		for _, category := range categories {
			lowerCaseName := strings.ToLower(category)
			if info, ok := m[strings.ToLower(category)]; ok {
				info.Count = info.Count + 1
			} else {
				m[lowerCaseName] = &CategoryInfo{
					Name:          category,
					LowerCaseName: lowerCaseName,
					Count:         1,
				}
				count = append(count, *m[lowerCaseName])
			}
		}
	}

	return count
}

func getFilter(f *glib.KeyFile, group string, keyname string) map[string]struct{} {
	_, list, _ := f.GetStringList(group, keyname)
	keyMap := map[string]struct{}{}
	for _, el := range list {
		keyMap[el] = struct{}{}
	}

	return keyMap
}

func getGroupNameFromDesktop(files []string) string {
	f := glib.NewKeyFile()
	defer f.Free()

	ok, err := f.LoadFromFile("/usr/share/dde/data/category_filter.ini", glib.KeyFileFlagsNone)
	if !ok {
		fmt.Println(err)
		return ""
	}

	invalidCategories := getFilter(f, "Main", "filter")
	genericCategories := getFilter(f, "Aux", "filter")

	categoryCounts := countCategories(files, invalidCategories)
	if len(categoryCounts) == 0 {
		return ""
	}

	// remove generic categories.
	candidateCategories := []CategoryInfo{}
	for _, category := range categoryCounts {
		if _, ok := genericCategories[category.LowerCaseName]; !ok {
			candidateCategories = append(candidateCategories, category)
		}
	}

	if len(candidateCategories) == 0 {
		// make sure there is one category at least.
		candidateCategories = categoryCounts
	}

	sort.Sort(CategoryInfos(candidateCategories))

	return candidateCategories[0].Name
}

func getGroupName(files []string) string {
	name := getGroupNameFromSoftwareCenter(files)
	if name != "" {
		return name
	}

	name = getGroupNameFromDesktop(files)
	if name != "" {
		return name
	}

	return Tr("App Group")
}
