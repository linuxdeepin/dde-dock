package item

import (
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	. "pkg.deepin.io/dde-daemon/launcher/category"
	. "pkg.deepin.io/dde-daemon/launcher/interfaces"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/utils"
)

const (
	DesktopSuffixLen = len(".desktop")
)

// #define FILENAME_WEIGHT 0.3
// #define GENERIC_NAME_WEIGHT 0.01
// #define KEYWORD_WEIGHT 0.1
// #define CATEGORY_WEIGHT 0.01
// #define NAME_WEIGHT 0.01
// #define DISPLAY_NAME_WEIGHT 0.1
// #define DESCRIPTION_WEIGHT 0.01
// #define EXECUTABLE_WEIGHT 0.05
type Xinfo struct {
	keywords    []string
	exec        string
	genericName string
	description string
}

type ItemInfo struct {
	path          string
	name          string
	enName        string
	id            ItemId
	icon          string
	categoryId    CategoryId
	timeInstalled int64
	xinfo         Xinfo
}

func (i *ItemInfo) Path() string {
	return i.path
}

func (i *ItemInfo) Name() string {
	return i.name
}
func (i *ItemInfo) EnName() string {
	return i.enName
}

func (i *ItemInfo) Id() ItemId {
	return i.id
}

func (i *ItemInfo) Keywords() []string {
	return i.xinfo.keywords
}

func (i *ItemInfo) GenericName() string {
	return i.xinfo.genericName
}

func NewItem(app *gio.DesktopAppInfo) *ItemInfo {
	if app == nil {
		return nil
	}
	item := &ItemInfo{}
	item.init(app)
	return item
}

func (i *ItemInfo) init(app *gio.DesktopAppInfo) {
	i.id = getId(app)
	i.path = app.GetFilename()
	i.name = app.GetDisplayName()
	i.enName = app.GetString("Name")
	icon := app.GetIcon()
	if icon != nil {
		i.icon = icon.ToString()
		if path.IsAbs(i.icon) && !utils.IsFileExist(i.icon) {
			i.icon = ""
		}
	}

	i.xinfo.keywords = make([]string, 0)
	keywords := app.GetKeywords()
	for _, keyword := range keywords {
		i.xinfo.keywords = append(i.xinfo.keywords, strings.ToLower(keyword))
	}
	i.xinfo.exec = app.GetCommandline()
	i.xinfo.genericName = app.GetGenericName()
	i.xinfo.description = app.GetDescription()
	i.categoryId = OtherID
}

func (i *ItemInfo) Description() string {
	return i.xinfo.description
}

func (i *ItemInfo) ExecCmd() string {
	return i.xinfo.exec
}

func (i *ItemInfo) Icon() string {
	return i.icon
}

func (i *ItemInfo) GetCategoryId() CategoryId {
	return i.categoryId
}

func (i *ItemInfo) SetCategoryId(id CategoryId) {
	i.categoryId = id
}

func (i *ItemInfo) GetTimeInstalled() int64 {
	return i.timeInstalled
}

func (i *ItemInfo) SetTimeInstalled(timeInstalled int64) {
	i.timeInstalled = timeInstalled
}

func GenId(filename string) ItemId {
	if len(filename) <= DesktopSuffixLen {
		return ItemId("")
	}

	basename := path.Base(filename)
	return ItemId(strings.Replace(basename[:len(basename)-DesktopSuffixLen], "_", "-", -1))
}

func getId(app *gio.DesktopAppInfo) ItemId {
	return GenId(app.GetFilename())
}
