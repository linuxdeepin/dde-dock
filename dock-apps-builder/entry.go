package main

import (
	"fmt"
)

//import "strings"

type Rectangle struct {
	X, Y          int16
	Width, Height uint16
}

const (
	InvalidStatus = iota
	ActiveStatus
	NormalStatus
)

type AppEntry struct {
	nApp *NormalApp
	rApp *RuntimeApp

	Data map[string]string
	Type string

	Id string

	Tooltip string
	Icon    string
	Menu    string

	Status int32

	QuickWindowViewable bool
	Allocation          Rectangle
}

func NewAppEntryWithRuntimeApp(rApp *RuntimeApp) *AppEntry {
	fmt.Println("NewAppEntryWithRuntimeApp:", rApp.Id, rApp.CurrentInfo.Xid)
	e := &AppEntry{
		Id:   rApp.Id,
		Type: "App",
		Data: make(map[string]string),
	}
	e.setPropStatus(ActiveStatus)
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
	e.setPropStatus(NormalStatus)
	e.attachNoramlApp(nApp)
	return e
}

func (e *AppEntry) QuickWindow(x, y int32) {}

func (e *AppEntry) HideQuickWindow() {}

func (e *AppEntry) ContextMenu(x, y int32) {}

func (e *AppEntry) HandleMenuItem(id int32) {
	switch e.Status {
	case NormalStatus:
		e.nApp.HandleMenuItem(id)
	case ActiveStatus:
		e.rApp.HandleMenuItem(id)
	}
}

func (e *AppEntry) Activate(x, y int32) {
	switch e.Status {
	case NormalStatus:
		e.nApp.Activate(x, y)
	case ActiveStatus:
		e.rApp.Activate(x, y)
	}
}

func (e *AppEntry) SecondaryActivate(x, y int32)        {}
func (e *AppEntry) OnDragEnter(x, y int32, data string) {}
func (e *AppEntry) OnDragLeave(x, y int32, data string) {}
func (e *AppEntry) OnDragOver(x, y int32, data string)  {}
func (e *AppEntry) OnDragDrop(x, y int32, data string)  {}

func (e *AppEntry) update() {
	if e.rApp != nil {
		e.setPropStatus(ActiveStatus)
	} else if e.nApp != nil {
		e.setPropStatus(NormalStatus)
	} else {
		LOGGER.Warning(e.Id + " goto an invalid status")
		return
	}
	//NOTE: sync this with NormalApp/RuntimeApp
	switch e.Status {
	case ActiveStatus:
		e.setPropTooltip(e.rApp.CurrentInfo.Title)
		e.setPropIcon(e.rApp.CurrentInfo.Icon)
		e.setPropMenu(e.rApp.Menu)
	case NormalStatus:
		e.setPropTooltip(e.nApp.Name)
		e.setPropIcon(e.nApp.Icon)
		e.setPropMenu(e.nApp.Menu)
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
	e.setPropData(rApp.xids)
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
