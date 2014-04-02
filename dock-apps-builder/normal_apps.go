package main

import "dlib/gio-2.0"
import "fmt"
import "path/filepath"

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
	app := &NormalApp{Id: filepath.Base(id[:len(id)-8])}
	LOGGER.Info(id)
	if filepath.IsAbs(id) {
		app.core = gio.NewDesktopAppInfoFromFilename(id)
	} else {
		app.core = gio.NewDesktopAppInfo(id)
	}
	if app.core == nil {
		return nil
	}
	app.Icon = app.core.GetIcon().ToString()
	app.Name = app.core.GetDisplayName()
	app.buildMenu()
	return app
}

func (app *NormalApp) buildMenu() {
	app.coreMenu = NewMenu()
	app.coreMenu.AppendItem(NewMenuItem("_Run", func() {
		_, err := app.core.Launch(make([]*gio.File, 0), nil)
		LOGGER.Warning("Launch App Failed: ", err)
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
		"_Undock",
		func() { /*TODO: do the real work*/
			fmt.Println("Undock")
		},
		true,
	)
	app.coreMenu.AppendItem(dockItem)

	app.Menu = app.coreMenu.GenerateJSON()
}

func (app *NormalApp) HandleMenuItem(id int32) {
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

func (app *NormalApp) Activate(x, y int32) {
	app.core.Launch(nil, nil)
}

func (app *NormalApp) setChangedCB(cb func()) {
	app.changedCB = cb
}

func (app *NormalApp) notifyChanged() {
	if app.changedCB != nil {
		app.changedCB()
	}
}
