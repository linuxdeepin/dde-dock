package dock

import (
	. "dlib/gettext"
	"dlib/gio-2.0"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type NormalApp struct {
	Id   string
	Icon string
	Name string
	Menu string

	changedCB func()

	core     *gio.DesktopAppInfo
	coreMenu *Menu
	dockItem *MenuItem
}

func NewNormalApp(id string) *NormalApp {
	app := &NormalApp{Id: strings.ToLower(filepath.Base(id[:len(id)-8]))}
	logger.Info("NewNormalApp:", id)
	if filepath.IsAbs(id) {
		app.core = gio.NewDesktopAppInfoFromFilename(id)
	} else {
		app.core = gio.NewDesktopAppInfo(id)
		if app.core == nil {
			logger.Info("guess desktop")
			if newId := guess_desktop_id(app.Id + ".desktop"); newId != "" {
				app.core = gio.NewDesktopAppInfo(newId)
			}
		}
	}
	if app.core == nil {
		return nil
	}
	app.Icon = getAppIcon(app.core)
	logger.Info("app icon:", app.Icon)
	app.Name = app.core.GetDisplayName()
	logger.Info("Name", app.Name)
	app.buildMenu()
	return app
}

func (app *NormalApp) buildMenu() {
	app.coreMenu = NewMenu()
	app.coreMenu.AppendItem(NewMenuItem(Tr("_Run"), func() {
		_, err := app.core.Launch(make([]*gio.File, 0), nil)
		logger.Warning("Launch App Failed: ", err)
	}, true))
	app.coreMenu.AddSeparator()
	for _, actionName := range app.core.ListActions() {
		name := actionName //NOTE: don't directly use 'actionName' with closure in an forloop
		app.coreMenu.AppendItem(NewMenuItem(
			app.core.GetActionName(actionName),
			func() { app.core.LaunchAction(name, nil) },
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
	app.core = gio.NewDesktopAppInfoFromFilename(name)
	app.Icon = app.core.GetIcon().ToString()
	app.Name = app.core.GetDisplayName()
	app.buildMenu()
	return app
}

func (app *NormalApp) Activate(x, y int32) error {
	// FIXME:
	// the launch will be successful even if the desktop file is not
	// existed.
	f, err := os.Open(app.core.GetFilename())
	if err != nil {
		return errors.New("invalid")
	}
	f.Close()

	b, err := app.core.Launch(nil, nil)
	logger.Warning(b)
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
