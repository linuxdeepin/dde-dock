package main

import "dlib/gio-2.0"

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
	app := &NormalApp{Id: id}
	app.core = gio.NewDesktopAppInfo(id)
	app.Icon = app.core.GetIcon().ToString()
	app.Name = app.core.GetDisplayName()
	app.buildMenu()
	return app
}

func (app *NormalApp) buildMenu() {
	app.coreMenu = NewMenu()
	for _, actionName := range app.core.ListActions() {
		name := actionName //NOTE: don't directly use 'actionName' with closure in an forloop
		app.coreMenu.AppendItem(&MenuItem{
			name,
			func() { app.core.LaunchAction(name, nil) },
			true,
		})
	}
	dockItem := &MenuItem{
		Name:   "_Undock",
		Action: func() { /*TODO: do the real work*/
		},
		Enabled: true,
	}
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
