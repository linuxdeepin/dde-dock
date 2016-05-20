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
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"pkg.deepin.io/lib/dbus"
	"sync"
)

const (
	entryDBusObjPathPrefix = "/dde/dock/entry/v1/"
	entryDBusDestPrefix    = "dde.dock.entry."
	entryDBusInterface     = "dde.dock.Entry"

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

	hashId string
	Id     string
	Type   string
	Data   map[string]string

	DataChanged func(string, string)
}

func NewAppEntry(id string) *AppEntry {
	e := &AppEntry{
		Id:   id,
		Type: "App",
		Data: make(map[string]string),
	}

	// Set hash id
	hasher := md5.New()
	hasher.Write([]byte(e.Id))
	// DBusObjectPath can't be start with digital number
	e.hashId = "d" + hex.EncodeToString(hasher.Sum(nil))
	return e
}

func NewAppEntryWithRuntimeApp(rApp *RuntimeApp) *AppEntry {
	logger.Debugf("NewAppEntryWithRuntimeApp: app id %s, win %v", rApp.Id, rApp.CurrentInfo.window)
	e := NewAppEntry(rApp.Id)
	e.setData(FieldStatus, ActiveStatus)
	e.attachRuntimeApp(rApp)
	return e
}

func NewAppEntryWithNormalApp(nApp *NormalApp) *AppEntry {
	logger.Debug("NewAppEntryWithNormalApp:", nApp.Id)
	e := NewAppEntry(nApp.Id)
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
	// x, y useless
	// only handle file://
	logger.Debug("HandleDragDrop:", data)
	if e.rApp != nil {
		e.rApp.HandleDragDrop(data, timestamp)
	} else if e.nApp != nil {
		e.nApp.HandleDragDrop(data, timestamp)
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
		for k, v := range e.rApp.windowInfoTable {
			xids = append(xids, XidInfo{uint32(k), v.Title})
		}
		b, _ := json.Marshal(xids)
		e.setData(FieldAppXids, string(b))
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
