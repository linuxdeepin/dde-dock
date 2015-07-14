package desktop

import (
	"os/exec"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/operations"
	"pkg.deepin.io/lib/utils"
	"sort"
)

func getBaseName(uri string) string {
	return filepath.Base(utils.DecodeURI(uri))
}

type byName []string

func (s byName) Less(i, j int) bool {
	return getBaseName(s[i]) < getBaseName(s[j])
}

func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byName) Len() int {
	return len(s)
}

// Desktop is desktop itself.
type Desktop struct {
	app  *Application
	menu *Menu
}

// NewDesktop creates new desktop.
func NewDesktop(app *Application) *Desktop {
	return &Desktop{
		app: app,
	}
}

func (desktop *Desktop) destroy() {
}

// GenMenu generates json format menu content used in DeepinMenu for Desktop itself.
func (desktop *Desktop) GenMenu() (*Menu, error) {
	desktop.menu = NewMenu()
	menu := desktop.menu

	sortSubMenu := NewMenu().SetIDGenerator(desktop.menu.genID)

	sortPolicies := desktop.app.settings.getSortPolicies()
	for _, sortPolicy := range sortPolicies {
		sortSubMenu.AppendItem(NewMenuItem(sortPoliciesName[sortPolicy], func(sortPolicy string) func() {
			return func() {
				desktop.app.emitRequestSort(sortPolicy)
			}
		}(sortPolicy), true))
	}

	sortMenuItem := NewMenuItem(Tr("_Sort by"), func() {}, true)
	sortMenuItem.subMenu = sortSubMenu

	menu.AppendItem(sortMenuItem)

	newSubMenu := NewMenu().SetIDGenerator(menu.genID)
	newSubMenu.AppendItem(NewMenuItem(Tr("_Folder"), func() {
		desktop.app.emitRequestCreateDirectory()
	}, true)).AppendItem(NewMenuItem(Tr("_Text document"), func() {
		desktop.app.emitRequestCreateFile()
	}, true))

	templatePath := GetUserSpecialDir(glib.UserDirectoryDirectoryTemplates)
	job := operations.NewGetTemplateJob(templatePath)
	templates := job.Execute()
	if len(templates) != 0 {
		newSubMenu.AddSeparator()
		sort.Sort(byName(templates))
		for _, template := range templates {
			templateURI := template
			newSubMenu.AppendItem(NewMenuItem(getBaseName(templateURI), func() {
				desktop.app.emitRequestCreateFileFromTemplate(templateURI)
			}, true))
		}
	}

	newMenuItem := NewMenuItem(Tr("_New"), func() {}, true)
	newMenuItem.subMenu = newSubMenu

	menu.AppendItem(newMenuItem).AppendItem(NewMenuItem(Tr("Open in _terminal"), func() {
		runInTerminal(GetUserSpecialDir(glib.UserDirectoryDirectoryDesktop), "")
	}, true)).AppendItem(NewMenuItem(Tr("_Paste"), func() {
		desktop.app.emitRequestPaste(GetDesktopDir())
	}, operations.CanPaste(GetDesktopDir())))

	// TODO: plugin
	if true {
		ShowModule := func(module string) {
			go func() {
				conn, err := dbus.SessionBus()
				if err != nil {
					return
				}

				obj := conn.Object("com.deepin.dde.ControlCenter", "/com/deepin/dde/ControlCenter")
				if obj != nil {
					obj.Call("com.deepin.dde.ControlCenter.ShowModule", 0, module).Store()
				}
			}()
		}
		menu.AddSeparator().AppendItem(NewMenuItem(Tr("_Display settings"), func() {
			ShowModule("display")
		}, true)).AppendItem(NewMenuItem(Tr("_Corner navigation"), func() {
			exec.Command("/usr/lib/deepin-daemon/dde-zone").Start()
		}, true)).AppendItem(NewMenuItem(Tr("Pe_rsonalize"), func() {
			ShowModule("personalization")
		}, true))
	}

	return menu, nil

}

// DSS = "com.deepin.dde.ControlCenter"
// DSS_MODULE =
//     SYSTEM_INFO:"system_info"
//     DISPLAY:"display"
//     PERSON:"personalization"
//
// dss_dbus = null
// dss_ShowModule = (module) ->
//     try
//         dss_dbus = DCore.DBus.session(DSS) if dss_dbus is null or dss_dbus is undefined
//         dss_dbus?.ShowModule(module)
//     catch e
//         echo "dss_ShowModule #{module} error:#{e}"
