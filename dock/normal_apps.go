package dock

import (
	"errors"
	"path/filepath"
	. "pkg.linuxdeepin.com/lib/gettext"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/utils"
	"strings"
)

type NormalApp struct {
	Id   string
	Icon string
	Name string
	Menu string

	changedCB func()

	path     string
	coreMenu *Menu
	dockItem *MenuItem
}

func NewNormalApp(id string) *NormalApp {
	basename := strings.ToLower(filepath.Base(id[:len(id)-8]))
	app := &NormalApp{Id: strings.Replace(basename, "_", "-", -1)}
	logger.Debug("NewNormalApp:", id)
	var core *DesktopAppInfo
	if strings.ContainsRune(id, filepath.Separator) {
		core = NewDesktopAppInfoFromFilename(id)
	} else {
		core = NewDesktopAppInfo(id)
		if core == nil {
			logger.Debug("guess desktop")
			if newId := guess_desktop_id(app.Id + ".desktop"); newId != "" {
				core = NewDesktopAppInfo(newId)
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
	core := NewDesktopAppInfo(app.Id)

	if core != nil {
		return core
	}

	if newId := guess_desktop_id(app.Id + ".desktop"); newId != "" {
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
		_, err := core.Launch(make([]*gio.File, 0), nil)
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
				core.LaunchAction(name, nil)
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
	_, err := core.Launch(nil, nil)
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
