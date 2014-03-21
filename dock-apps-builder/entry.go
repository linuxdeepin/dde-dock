package main

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("NewAppEntryWithRuntimeApp:", rApp.Id, rApp.CurrentInfo.Xid)
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
	fmt.Println("NewAppEntryWithNormalApp:", nApp.Id)
	e := &AppEntry{
		Id:   nApp.Id,
		Type: "App",
		Data: make(map[string]string),
	}
	e.setData(FieldStatus, NormalStatus)
	e.attachNoramlApp(nApp)
	return e
}

func (e *AppEntry) HandleMenuItem(id int32) {
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

func (e *AppEntry) ContextMenu(x, y int32)              {}
func (e *AppEntry) SecondaryActivate(x, y int32)        {}
func (e *AppEntry) OnDragEnter(x, y int32, data string) {}
func (e *AppEntry) OnDragLeave(x, y int32, data string) {}
func (e *AppEntry) OnDragOver(x, y int32, data string)  {}
func (e *AppEntry) OnDragDrop(x, y int32, data string)  {}

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
func (e *AppEntry) attachNoramlApp(nApp *NormalApp) {
	if e.nApp != nil {
		return
	}
	e.nApp = nApp
	fmt.Println("AttachNormalApp:", e.nApp.Id)
	e.nApp.setChangedCB(e.update)
	e.update()
}
func (e *AppEntry) detachNormalApp() {
	if e.nApp != nil {
		fmt.Println("DetachNormalApp", e.nApp.Id)
		e.nApp = nil
		e.nApp.setChangedCB(nil)
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
	fmt.Println("AttachRuntimeApp:", e.rApp.Id)
	e.rApp.setChangedCB(e.update)
	e.update()
}
func (e *AppEntry) detachRuntimeApp() {
	if e.rApp != nil {
		fmt.Println("DetachRuntimeApp:", e.rApp.Id)
		e.rApp.setChangedCB(nil)
		e.rApp = nil
		if e.nApp != nil {
			e.update()
		}
	}
}
