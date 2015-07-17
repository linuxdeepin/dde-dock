package desktop

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/operations"
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

// GenMenu generates json format menu content used in DeepinMenu for AppGroup.
func (item *AppGroup) GenMenu() (*Menu, error) {
	item.menu = NewMenu()
	return item.menu.AppendItem(NewMenuItem(Tr("_Open"), func() {
		ops := make([]int32, len(item.uris))
		for i := range item.uris {
			ops[i] = OpOpen
		}
		item.app.emitRequestOpen(item.uris, ops)
	}, true)).AddSeparator().AppendItem(NewMenuItem(Tr("_Rename"), func() {
		item.emitRequestRename()
	}, !item.multiple)).AddSeparator().AppendItem(NewMenuItem(Tr("_Ungroup"), func() {
		for _, uri := range item.uris {
			files := []string{}

			listJob := operations.NewListDirJob(uri, operations.ListJobFlagIncludeHidden)
			listJob.ListenProperty(func(p operations.ListProperty) {
				files = append(files, p.URI)
			})
			listJob.Execute()
			if err := listJob.GetError(); err != nil {
				fmt.Printf("list appgroup %s failed: %s\n", uri, err)
				continue
			}

			moveJob := operations.NewMoveJob(files, GetDesktopDir(), "", 0, nil)
			moveJob.Execute()
			if err := moveJob.GetError(); err != nil {
				fmt.Printf("dismiss appgroup %s failed: %s\n", uri, err)
			}
		}
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

func getGroupNameFromSoftwareCenter(files []string) (string, error) {
	db, err := sql.Open("sqlite3", "dbPath")
	if err != nil {
		return "", err
	}
	defer db.Close()

	category, err := getCategoryFromSoftwareCenter(db, files[0])
	if err != nil {
		return "", fmt.Errorf("get category failed from database: %s", err.Error())
	}
	if category == "" {
		return "", errors.New("empty category from database")
	}

	for _, file := range files[1:] {
		anotherCategory, err := getCategoryFromSoftwareCenter(db, file)
		if err != nil || anotherCategory != category {
			return "", errors.New("no same category from database")
		}
	}

	return category, nil
}

func containsWithCaseInsensitive(m map[string]struct{}, category string) bool {
	category = strings.ToLower(category)
	_, ok := m[category]
	return ok
}

func getCategoriesFromDesktop(file string, invalidCategories map[string]struct{}) (categories []string) {
	app := gio.NewDesktopAppInfoFromFilename(file)
	if app == nil {
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

// CategoryInfo includes simple information for software category.
type CategoryInfo struct {
	Name          string
	LowerCaseName string
	Count         int32
}

// CategoryInfos is an array of CategoryInfo, used by sort.Sort.
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

func getGroupNameFromDesktop(files []string) (string, error) {
	f := glib.NewKeyFile()
	defer f.Free()

	ok, err := f.LoadFromFile("/usr/share/dde/data/category_filter.ini", glib.KeyFileFlagsNone)
	if !ok {
		return "", fmt.Errorf("load category filter failed: %s", err.Error())
	}

	invalidCategories := getFilter(f, "Main", "filter")
	genericCategories := getFilter(f, "Aux", "filter")

	categoryCounts := countCategories(files, invalidCategories)
	if len(categoryCounts) == 0 {
		return "", errors.New("get no category from desktop")
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

	return candidateCategories[0].Name, nil
}

func getGroupName(files []string) string {
	name, err := getGroupNameFromSoftwareCenter(files)
	if name != "" {
		return name
	}

	if err != nil {
		fmt.Println(err)
	}

	name, err = getGroupNameFromDesktop(files)
	if name != "" {
		return name
	}

	if err != nil {
		fmt.Println(err)
	}

	return Tr("App Group")
}
