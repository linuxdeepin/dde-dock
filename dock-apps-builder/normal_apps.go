package main

import "dlib/gio-2.0"

type NormalApp struct {
	Id   string
	Icon string
	Name string

	changedCB func()
	core      *gio.DesktopAppInfo
}

func NewNormalApp(id string) *NormalApp {
	app := &NormalApp{Id: id}
	app.core = gio.NewDesktopAppInfo(id)
	app.Icon = app.core.GetIcon().ToString()
	app.Name = app.core.GetDisplayName()
	return app
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
