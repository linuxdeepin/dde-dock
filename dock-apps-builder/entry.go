package main

import (
	"dlib/gio-2.0"
	"encoding/json"
	"strings"
)

const (
	FieldTitle   = "title"
	FieldIcon    = "icon"
	FieldMenu    = "menu"
	FieldAppXids = "app-xids"

	FieldStatus   = "app-status"
	ActiveStatus  = "active"
	NormalStatus  = "normal"
	InvalidStatus = "invalid"
)

type AppEntry struct {
	nApp *NormalApp
	rApp *RuntimeApp

	Id   string
	Type string
	Data map[string]string

	DataChanged func(string, string)
}

func NewAppEntryWithRuntimeApp(rApp *RuntimeApp) *AppEntry {
	LOGGER.Info("NewAppEntryWithRuntimeApp:", rApp.Id, rApp.CurrentInfo.Xid)
	e := &AppEntry{
		Id:   rApp.Id,
		Type: "App",
		Data: make(map[string]string),
	}
	e.setData(FieldStatus, ActiveStatus)
	e.attachRuntimeApp(rApp)
	return e
}
func NewAppEntryWithNormalApp(nApp *NormalApp) *AppEntry {
	LOGGER.Info("NewAppEntryWithNormalApp:", nApp.Id)
	e := &AppEntry{
		Id:   nApp.Id,
		Type: "App",
		Data: make(map[string]string),
	}
	e.setData(FieldStatus, NormalStatus)
	e.attachNormalApp(nApp)
	return e
}

func (e *AppEntry) HandleMenuItem(id string) {
	switch e.Data[FieldStatus] {
	case NormalStatus:
		e.nApp.HandleMenuItem(id)
	case ActiveStatus:
		e.rApp.HandleMenuItem(id)
	}
}

func (e *AppEntry) Activate(x, y int32) {
	switch e.Data[FieldStatus] {
	case NormalStatus:
		e.nApp.Activate(x, y)
	case ActiveStatus:
		e.rApp.Activate(x, y)
	}
}

func (e *AppEntry) ContextMenu(x, y int32)                  {}
func (e *AppEntry) SecondaryActivate(x, y int32)            {}
func (e *AppEntry) HandleDragEnter(x, y int32, data string) {}
func (e *AppEntry) HandleDragLeave(x, y int32, data string) {}
func (e *AppEntry) HandleDragOver(x, y int32, data string)  {}
func (e *AppEntry) HandleDragDrop(x, y int32, data string) {
	paths := strings.Split(data, ",")
	LOGGER.Debug("HandleDragDrop:", paths)
	if e.rApp != nil {
		LOGGER.Info("Launch from runtime app")
		if e.rApp.core != nil {
			_, err := e.rApp.core.LaunchUris(paths, nil)
			if err != nil {
				LOGGER.Error("Launch Drop failed:", err)
			}
		} else {
			app, err :=
				gio.AppInfoCreateFromCommandline(e.rApp.exec,
					e.rApp.Id, gio.AppInfoCreateFlagsSupportsUris)
			if err != nil {
				LOGGER.Error("Create Launch app failed:", err)
				return
			}

			_, err = app.LaunchUris(paths, nil)
			if err != nil {
				LOGGER.Error("Launch Drop failed:", err)
			}
		}
	} else if e.nApp != nil {
		LOGGER.Info("Launch from runtime app")
		if e.nApp.core != nil {
			_, err := e.nApp.core.LaunchUris(paths, nil)
			if err != nil {
				LOGGER.Error("Launch Drop failed:", err)
			}
		} else {
			// TODO:
			LOGGER.Error("TODO: AppEntry.nApp.core == nil")
		}
	}
}
func (e *AppEntry) HandleMouseWheel(x, y, delta int32) {}

func (e *AppEntry) setData(key, value string) {
	if e.Data[key] != value {
		e.Data[key] = value
		if e.DataChanged != nil {
			e.DataChanged(key, value)
		}
	}
}
func (e *AppEntry) getData(key string) string {
	return e.Data[key]
}

type XidInfo struct {
	Xid   uint32
	Title string
}

func (e *AppEntry) update() {
	if e.rApp != nil {
		e.setData(FieldStatus, ActiveStatus)
		xids := make([]XidInfo, 0)
		for k, v := range e.rApp.xids {
			xids = append(xids, XidInfo{uint32(k), v.Title})
		}
		b, _ := json.Marshal(xids)
		e.setData(FieldAppXids, string(b))
	} else if e.nApp != nil {
		e.setData(FieldStatus, NormalStatus)
	} else {
		LOGGER.Warning(e.Id + " goto an invalid status")
		return
	}
	//NOTE: sync this with NormalApp/RuntimeApp
	switch e.getData(FieldStatus) {
	case ActiveStatus:
		e.setData(FieldTitle, e.rApp.CurrentInfo.Title)
		e.setData(FieldIcon, e.rApp.CurrentInfo.Icon)
		e.setData(FieldMenu, e.rApp.Menu)
	case NormalStatus:
		e.setData(FieldTitle, e.nApp.Name)
		e.setData(FieldIcon, e.nApp.Icon)
		e.setData(FieldMenu, e.nApp.Menu)
	}
}
func (e *AppEntry) attachNormalApp(nApp *NormalApp) {
	if e.nApp != nil {
		return
	}
	e.nApp = nApp
	LOGGER.Info("AttachNormalApp:", e.nApp.Id)
	e.nApp.setChangedCB(e.update)
	e.update()
}
func (e *AppEntry) detachNormalApp() {
	if e.nApp != nil {
		LOGGER.Info("DetachNormalApp", e.nApp.Id)
		e.nApp.setChangedCB(nil)
		e.nApp = nil
		if e.rApp != nil {
			e.update()
		}
	}
}
func (e *AppEntry) attachRuntimeApp(rApp *RuntimeApp) {
	if e.rApp != nil {
		return
	}
	e.rApp = rApp
	LOGGER.Info("AttachRuntimeApp:", e.rApp.Id)
	e.rApp.setChangedCB(e.update)
	e.update()
}
func (e *AppEntry) detachRuntimeApp() {
	if e.rApp != nil {
		LOGGER.Info("DetachRuntimeApp:", e.rApp.Id)
		e.rApp.setChangedCB(nil)
		e.rApp = nil
		if e.nApp != nil {
			e.update()
		}
	}
}
