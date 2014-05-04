package main

import (
	"dlib/dbus"
)

const entryDestPrefix = "dde.dock.entry."
const entryPathPrefix = "/dde/dock/entry/v1/"

type EntryProxyer struct {
	entryId string
	core    *RemoteEntry

	Id          string
	Type        string
	Data        map[string]string
	DataChanged func(string, string)
}

func NewEntryProxyer(entryId string) (*EntryProxyer, error) {
	if core, err := NewRemoteEntry(entryDestPrefix+entryId, dbus.ObjectPath(entryPathPrefix+entryId)); err != nil {
		return nil, err
	} else {
		e := &EntryProxyer{
			core:    core,
			entryId: entryId,
			Id:      core.Id.Get(),
			Type:    core.Type.Get(),
			Data:    core.Data.Get(),
		}
		e.core.ConnectDataChanged(func(key, value string) {
			if e.DataChanged != nil {
				e.Data[key] = value
				e.DataChanged(key, value)
			}
		})
		return e, nil
	}
}

func (e *EntryProxyer) ContextMenu(x, y int32)              { e.core.ContextMenu(x, y) }
func (e *EntryProxyer) HandleMenuItem(id int32)             { e.core.HandleMenuItem(id) }
func (e *EntryProxyer) Activate(x, y int32)                 { e.core.Activate(x, y) }
func (e *EntryProxyer) SecondaryActivate(x, y int32)        { e.core.SecondaryActivate(x, y) }
func (e *EntryProxyer) OnDragEnter(x, y int32, data string) { e.core.OnDragEnter(x, y, data) }
func (e *EntryProxyer) OnDragLeave(x, y int32, data string) { e.core.OnDragLeave(x, y, data) }
func (e *EntryProxyer) OnDragOver(x, y int32, data string)  { e.core.OnDragOver(x, y, data) }
func (e *EntryProxyer) OnDragDrop(x, y int32, data string)  { e.core.OnDragDrop(x, y, data) }
func (e *EntryProxyer) OnMouseWheel(x, y, delta int32) {
	e.core.OnMouseWheel(x, y, delta)
}
func (e *EntryProxyer) ShowQuickWindow() { e.core.ShowQuickWindow() }

func (e *EntryProxyer) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		entryPathPrefix + e.entryId,
		"dde.dock.EntryProxyer",
	}
}
