package desktop

import (
	"fmt"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/operations"
)

// Desktop is desktop itself.
type Desktop struct {
	app  *Application
	menu *Menu
}

// NewDesktop creates new desktop.
func NewDesktop(app *Application) *Desktop {
	return &Desktop{
		app:  app,
		menu: NewMenu(),
	}
}

func (desktop *Desktop) destroy() {
}

// GenMenuContent generates json format menu content used in DeepinMenu for Desktop itself.
func (desktop *Desktop) GenMenuContent() (*Menu, error) {
	sortSubMenu := NewMenu().SetIDGenerator(desktop.menu.genID)

	sortPolicies := desktop.app.settings.getSortPolicies()
	for _, sortPolicy := range sortPolicies {
		sortSubMenu.AppendItem(NewMenuItem(PoliciesName[sortPolicy], func(sortPolicy string) func() {
			return func() {
				desktop.app.emitRequestSort(sortPolicy)
			}
		}(sortPolicy), true))
	}

	sortMenuItem := NewMenuItem(Tr("_Sort by"), func() {}, true)
	sortMenuItem.subMenu = sortSubMenu

	menu := desktop.menu.AppendItem(sortMenuItem)

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
		for _, template := range templates {
			templateURI := template
			newSubMenu.AppendItem(NewMenuItem("", func() {
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
		menu.AddSeparator().AppendItem(NewMenuItem(Tr("_Display settings"), func() {
			// TODO
			fmt.Println("show display settings")
		}, true)).AppendItem(NewMenuItem(Tr("_Corner navigation"), func() {
			// TODO
			fmt.Println("show corner navigation")
		}, true)).AppendItem(NewMenuItem(Tr("Pe_rsonalize"), func() {
			// TODO
			fmt.Println("show personalize settings")
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
