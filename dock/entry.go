/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"encoding/json"
	"strings"
	"sync"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
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
	nApp     *NormalApp
	rApp     *RuntimeApp
	nAppLock sync.RWMutex
	rAppLock sync.RWMutex

	Id   string
	Type string
	Data map[string]string

	DataChanged func(string, string)
}

func NewAppEntryWithRuntimeApp(rApp *RuntimeApp) *AppEntry {
	logger.Debugf("NewAppEntryWithRuntimeApp: %s, 0x%x", rApp.Id, rApp.CurrentInfo.Xid)
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
	logger.Debug("NewAppEntryWithNormalApp:", nApp.Id)
	e := &AppEntry{
		Id:   nApp.Id,
		Type: "App",
		Data: make(map[string]string),
	}
	e.setData(FieldStatus, NormalStatus)
	e.attachNormalApp(nApp)
	return e
}

func (e *AppEntry) HandleMenuItem(id string, timestamp uint32) {
	switch e.Data[FieldStatus] {
	case NormalStatus:
		e.nApp.HandleMenuItem(id, timestamp)
	case ActiveStatus:
		e.rApp.HandleMenuItem(id, timestamp)
	}
}

func (e *AppEntry) Activate(x, y int32, timestamp uint32) (bool, error) {
	switch e.Data[FieldStatus] {
	case NormalStatus:
		err := e.nApp.Activate(x, y, timestamp)
		return err == nil, err
	case ActiveStatus:
		err := e.rApp.Activate(x, y, timestamp)
		return err == nil, err
	}
	panic("should not reach")
}

func (e *AppEntry) ContextMenu(x, y int32)                                    {}
func (e *AppEntry) SecondaryActivate(x, y int32, timestamp uint32)            {}
func (e *AppEntry) HandleDragEnter(x, y int32, data string, timestamp uint32) {}
func (e *AppEntry) HandleDragLeave(x, y int32, data string, timestamp uint32) {}
func (e *AppEntry) HandleDragOver(x, y int32, data string, timestamp uint32)  {}
func (e *AppEntry) HandleDragDrop(x, y int32, data string, timestamp uint32) {
	paths := strings.Split(data, ",")
	logger.Debug("HandleDragDrop:", paths)
	if e.rApp != nil {
		logger.Debug("Launch from runtime app")
		core := e.rApp.createDesktopAppInfo()
		if core != nil {
			defer core.Destroy()
			_, err := core.LaunchUris(paths, gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
			if err != nil {
				logger.Warning("Launch Drop failed:", err)
			}
		} else {
			app, err :=
				gio.AppInfoCreateFromCommandline(e.rApp.exec,
					e.rApp.Id, gio.AppInfoCreateFlagsSupportsUris)
			if err != nil {
				logger.Warning("Create Launch app failed:", err)
				return
			}

			_, err = app.LaunchUris(paths, gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
			if err != nil {
				logger.Warning("Launch Drop failed:", err)
			}
		}
	} else if e.nApp != nil {
		logger.Debug("Launch from normal app")
		core := e.nApp.createDesktopAppInfo()
		if core != nil {
			defer core.Destroy()
			_, err := core.LaunchUris(paths, gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
			if err != nil {
				logger.Warning("Launch Drop failed:", err)
			}
		} else {
			// TODO:
			logger.Warning("TODO: AppEntry.nApp.core == nil")
		}
	}
}
func (e *AppEntry) HandleMouseWheel(x, y, delta int32, timestamp uint32) {}

func (e *AppEntry) setData(key, value string) {
	if e.Data[key] != value {
		e.Data[key] = value
		dbus.Emit(e, "DataChanged", key, value)
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
		if dockManager.hideStateManager.state == HideStateShown {
			dockManager.hideStateManager.updateStateWithDelay()
		}
	} else if e.nApp != nil {
		e.setData(FieldStatus, NormalStatus)
	} else {
		logger.Warning(e.Id + " goto an invalid status")
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
	e.nAppLock.Lock()
	defer e.nAppLock.Unlock()
	if e.nApp != nil {
		return
	}
	e.nApp = nApp
	logger.Debug("AttachNormalApp:", e.nApp.Id)
	e.nApp.setChangedCB(e.update)
	e.update()
}
func (e *AppEntry) detachNormalApp() {
	e.nAppLock.Lock()
	defer e.nAppLock.Unlock()
	if e.nApp != nil {
		logger.Debug("DetachNormalApp", e.nApp.Id)
		e.nApp.setChangedCB(nil)
		e.nApp = nil
		if e.rApp != nil {
			e.update()
		}
	}
}
func (e *AppEntry) attachRuntimeApp(rApp *RuntimeApp) {
	e.rAppLock.Lock()
	defer e.rAppLock.Unlock()
	if e.rApp != nil {
		return
	}
	e.rApp = rApp
	logger.Debug("AttachRuntimeApp:", e.rApp.Id)
	e.rApp.setChangedCB(e.update)
	e.update()
}
func (e *AppEntry) detachRuntimeApp() {
	e.rAppLock.Lock()
	defer e.rAppLock.Unlock()
	if e.rApp != nil {
		logger.Debug("DetachRuntimeApp:", e.rApp.Id)
		e.rApp.setChangedCB(nil)
		e.rApp = nil
		if e.nApp != nil {
			e.update()
		}
	}
}
