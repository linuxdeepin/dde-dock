package dock

import (
	"errors"
	"fmt"
	"path/filepath"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/utils"
	"strings"
)

type NormalApp struct {
	Id        string
	DesktopID string
	Icon      string
	Name      string
	Menu      string

	changedCB func()

	path     string
	coreMenu *Menu
	dockItem *MenuItem
}

func NewNormalApp(desktopID string) *NormalApp {
	app := &NormalApp{Id: normalizeAppID(trimDesktop(desktopID)), DesktopID: desktopID}
	logger.Info("NewNormalApp:", app.Id, "for desktop", desktopID)
	var core *DesktopAppInfo
	if strings.ContainsRune(desktopID, filepath.Separator) {
		core = NewDesktopAppInfoFromFilename(desktopID)
	} else {
		core = NewDesktopAppInfo(desktopID)
		if core == nil {
			newId := guess_desktop_id(app.Id)
			logger.Info(fmt.Sprintf("guess desktop: %q", newId))
			if newId != "" {
				core = NewDesktopAppInfo(newId)
				app.DesktopID = newId
			}
		}
	}
	if core == nil {
		return nil
	}
	defer core.Unref()
	app.path = core.GetFilename()
	app.Icon = getAppIcon(core.DesktopAppInfo)
	logger.Debug(app.Id, "::app icon:", app.Icon)
	app.Name = core.GetDisplayName()
	logger.Debug("Name", app.Name)
	app.buildMenu(core)
	return app
}

func (app *NormalApp) createDesktopAppInfo() *DesktopAppInfo {
	core := NewDesktopAppInfo(app.DesktopID)

	if core != nil {
		return core
	}

	if newId := guess_desktop_id(app.Id); newId != "" {
		core = NewDesktopAppInfo(newId)
		if core != nil {
			return core
		}
	}

	return NewDesktopAppInfoFromFilename(app.path)
}

func (app *NormalApp) buildMenu(core *DesktopAppInfo) {
	app.coreMenu = NewMenu()
	app.coreMenu.AppendItem(NewMenuItem(Tr("_Run"), func() {
		core := app.createDesktopAppInfo()
		if core == nil {
			logger.Warning("Run app failed")
			return
		}
		defer core.Unref()
		_, err := core.Launch(make([]*gio.File, 0), gio.GetGdkAppLaunchContext())
		if err != nil {
			logger.Warning("Launch App Failed: ", err)
		}
	}, true))
	app.coreMenu.AddSeparator()
	for _, actionName := range core.ListActions() {
		name := actionName //NOTE: don't directly use 'actionName' with closure in an forloop
		app.coreMenu.AppendItem(NewMenuItem(
			core.GetActionName(actionName),
			func() {
				core := app.createDesktopAppInfo()
				if core == nil {
					logger.Warning("start action", name,
						"failed")
					return
				}
				defer core.Unref()
				core.LaunchAction(name, gio.GetGdkAppLaunchContext())
			},
			true,
		))
	}
	app.coreMenu.AddSeparator()
	dockItem := NewMenuItem(
		Tr("_Undock"),
		func() {
			DOCKED_APP_MANAGER.Undock(app.Id)
		},
		true,
	)
	app.coreMenu.AppendItem(dockItem)

	app.Menu = app.coreMenu.GenerateJSON()
}

func (app *NormalApp) HandleMenuItem(id string) {
	if app.coreMenu != nil {
		app.coreMenu.HandleAction(id)
	}
}

func NewNormalAppFromFilename(name string) *NormalApp {
	app := &NormalApp{}
	core := NewDesktopAppInfoFromFilename(name)
	if core == nil {
		return nil
	}
	defer core.Unref()
	app.path = core.GetFilename()
	app.Icon = core.GetIcon().ToString()
	app.Name = core.GetDisplayName()
	app.buildMenu(core)
	return app
}

func (app *NormalApp) Activate(x, y int32) error {
	// FIXME:
	// the launch will be successful even if the desktop file is not
	// existed.
	if !utils.IsFileExist(app.path) {
		return errors.New("invalid")
	}

	core := app.createDesktopAppInfo()
	if core == nil {
		return errors.New("create desktop app info failed")
	}
	defer core.Unref()
	_, err := core.Launch(nil, gio.GetGdkAppLaunchContext())
	if err != nil {
		logger.Warning("launch", app.Id, "failed:", err)
	}
	return err
}

func (app *NormalApp) setChangedCB(cb func()) {
	app.changedCB = cb
}

func (app *NormalApp) notifyChanged() {
	if app.changedCB != nil {
		app.changedCB()
	}
}
