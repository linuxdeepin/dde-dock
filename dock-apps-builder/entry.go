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

	Status int32

	QuickWindowVieable bool
	Allocation         Rectangle
}

func NewAppEntryWithRuntimeApp(rApp *RuntimeApp) *AppEntry {
	fmt.Println("NewAppEntryWithRuntimeApp:", rApp.Id, rApp.CurrentInfo.Xid)
	e := &AppEntry{
		Id:     rApp.Id,
		Type:   "App",
		Status: ActiveStatus,
		Data:   make(map[string]string),
	}
	e.attachRuntimeApp(rApp)
	return e
}
func NewAppEntryWithNormalApp(nApp *NormalApp) *AppEntry {
	fmt.Println("NewAppEntryWithNormalApp:", nApp.Id)
	e := &AppEntry{
		Id:     nApp.Id,
		Type:   "App",
		Status: NormalStatus,
		Data:   make(map[string]string),
	}
	e.attachNoramlApp(nApp)
	return e
}

func (e *AppEntry) QuickWindow(x, y int32) {}

func (e *AppEntry) ContextMenu(x, y int32) {}

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
		e.Status = ActiveStatus
	} else if e.nApp != nil {
		e.Status = NormalStatus
	} else {
		LOGGER.Warning(e.Id + " goto an invalid status")
		return
	}
	switch e.Status {
	case ActiveStatus:
		e.Tooltip = e.rApp.CurrentInfo.Title
		e.Icon = e.rApp.CurrentInfo.Icon
	case NormalStatus:
		e.Icon = e.nApp.Icon
		e.Tooltip = e.nApp.Name
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
