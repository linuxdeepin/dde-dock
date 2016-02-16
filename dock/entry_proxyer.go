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
	"errors"

	"pkg.deepin.io/lib/dbus"
)

const entryDestPrefix = "dde.dock.entry."
const entryPathPrefix = "/dde/dock/entry/v1/"

// EntryProxyer为驻留程序以及打开程序的dbus接口。
type EntryProxyer struct {
	entryId string
	core    *RemoteEntry

	// Id属性为程序id。
	Id string
	// Type属性为程序类型，包括app，applet两种。
	Type string
	// Data包含其他程序相关属性，如icon，title，status，menu，app-xids。
	Data map[string]string
	// DataChanged会在Data属性有改变是被触发。
	DataChanged func(string, string)
}

func NewEntryProxyer(entryId string) (*EntryProxyer, error) {
	if core, err := NewRemoteEntry(entryDestPrefix+entryId, dbus.ObjectPath(entryPathPrefix+entryId)); err != nil {
		return nil, err
	} else {
		remoteEntryId := core.Id.Get()
		if "" == remoteEntryId {
			return nil, errors.New("Empty remote entry id")
		}
		e := &EntryProxyer{
			core:    core,
			entryId: entryId,
			Id:      remoteEntryId,
			Type:    core.Type.Get(),
			Data:    core.Data.Get(),
		}
		e.core.ConnectDataChanged(func(key, value string) {
			if e.Data != nil {
				e.Data[key] = value
				dbus.Emit(e, "DataChanged", key, value)
			}
		})
		return e, nil
	}
}

// ContextMenu接受鼠标事件的位置信息，然后生成
func (e *EntryProxyer) ContextMenu(x, y int32) { e.core.ContextMenu(x, y) }

func (e *EntryProxyer) HandleMenuItem(id string) {
	e.HandleMenuItemWithTimestamp(id, 0)
}

// HandleMenuItem对出入的id在右键菜单中对应的项做处理。
func (e *EntryProxyer) HandleMenuItemWithTimestamp(id string, timestamp uint32) {
	e.core.HandleMenuItem(id, timestamp)
}

func (e *EntryProxyer) Activate(x, y int32) (bool, error) {
	return e.ActivateWithTimestamp(x, y, 0)
}

// Activate在程序被点击时作出响应，接受鼠标事件的位置信息。
func (e *EntryProxyer) ActivateWithTimestamp(x, y int32, timestamp uint32) (bool, error) {
	return e.core.Activate(x, y, timestamp)
}

func (e *EntryProxyer) SecondaryActivate(x, y int32) {
	e.SecondaryActivateWithTimestamp(x, y, 0)
}

// SecondaryActivate与Activate作用相同，可用于其他鼠标点击事件，通常不会被使用。
func (e *EntryProxyer) SecondaryActivateWithTimestamp(x, y int32, timestamp uint32) {
	e.core.SecondaryActivate(x, y, timestamp)
}

// HandleDragEnter在前端触发DragEnter事件时被调用。
func (e *EntryProxyer) HandleDragEnter(x, y int32, data string, timestamp uint32) {
	e.core.HandleDragEnter(x, y, data, timestamp)
}

// HandleDragLeave在前端触发DragLeave事件时被调用。
func (e *EntryProxyer) HandleDragLeave(x, y int32, data string, timestamp uint32) {
	e.core.HandleDragLeave(x, y, data, timestamp)
}

// HandleDragOver在前端触发DragOver事件时被调用。
func (e *EntryProxyer) HandleDragOver(x, y int32, data string, timestamp uint32) {
	e.core.HandleDragOver(x, y, data, timestamp)
}

// HandleDragDrop在前端触发Dropp事件时被调用。
func (e *EntryProxyer) HandleDragDrop(x, y int32, data string, timestamp uint32) {
	e.core.HandleDragDrop(x, y, data, timestamp)
}

// HandleMouseWheel在前端触发鼠标滚轮事件时被调用。
func (e *EntryProxyer) HandleMouseWheel(x, y, delta int32, timestamp uint32) {
	e.core.HandleMouseWheel(x, y, delta, timestamp)
}

// ShowQuickWindow用于applet程序中，在需要显示时applet的窗口调用。
func (e *EntryProxyer) ShowQuickWindow() { e.core.ShowQuickWindow() }

func (e *EntryProxyer) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Dock",
		ObjectPath: entryPathPrefix + e.entryId,
		Interface:  "dde.dock.EntryProxyer",
	}
}
