package main

import (
	"dlib/dbus"
)

func (e *EntryProxyer) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"dde.dock.EntryManager",
		entryPathPrefix + e.entryId,
		"dde.dock.EntryProxyer",
	}
}

func NewEntryProxyer(entryId string) (e *EntryProxyer, err error) {
	e = &EntryProxyer{}
	e.destPath = entryDestPrefix + entryId
	e.objectPath = dbus.ObjectPath(entryPathPrefix + entryId)
	core, err := NewRemoteEntry(e.destPath, e.objectPath)
	if err != nil {
		return
	}

	// init properties
	e.entryId = entryId
	e.core = core

	e.setProperty("Id")
	e.setProperty("Type")
	e.setProperty("Tooltip")
	e.setProperty("Icon")
	e.setProperty("Status")
	e.setProperty("QuickWindowViewable")
	e.setProperty("Allocation")
	e.setProperty("Data")
	e.setProperty("Menu")

	// monitor properties changed
	e.core.Id.ConnectChanged(func() {
		e.setProperty("Id")
	})
	e.core.Type.ConnectChanged(func() {
		e.setProperty("Type")
	})
	e.core.Tooltip.ConnectChanged(func() {
		e.setProperty("Tooltip")
	})
	e.core.Icon.ConnectChanged(func() {
		e.setProperty("Icon")
	})
	e.core.Status.ConnectChanged(func() {
		e.setProperty("Status")
	})
	e.core.QuickWindowViewable.ConnectChanged(func() {
		e.setProperty("QuickWindowViewable")
	})
	e.core.Allocation.ConnectChanged(func() {
		e.setProperty("Allocation")
	})
	e.core.Data.ConnectChanged(func() {
		e.setProperty("Data")
	})
	e.core.Data.ConnectChanged(func() {
		e.setProperty("Menu")
	})

	return
}

func (e *EntryProxyer) QuickWindow(x, y int32)              { e.core.QuickWindow(x, y) }
func (e *EntryProxyer) HideQuickWindow()                    { e.core.HideQuickWindow() }
func (e *EntryProxyer) ContextMenu(x, y int32)              { e.core.ContextMenu(x, y) }
func (e *EntryProxyer) HandleMenuItem(id int32)             { e.core.HandleMenuItem(id) }
func (e *EntryProxyer) Activate(x, y int32)                 { e.core.Activate(x, y) }
func (e *EntryProxyer) SecondaryActivate(x, y int32)        { e.core.SecondaryActivate(x, y) }
func (e *EntryProxyer) OnDragEnter(x, y int32, data string) { e.core.OnDragEnter(x, y, data) }
func (e *EntryProxyer) OnDragLeave(x, y int32, data string) { e.core.OnDragLeave(x, y, data) }
func (e *EntryProxyer) OnDragOver(x, y int32, data string)  { e.core.OnDragOver(x, y, data) }
func (e *EntryProxyer) OnDragDrop(x, y int32, data string)  { e.core.OnDragDrop(x, y, data) }

func (e *EntryProxyer) setProperty(prop string) {
	switch prop {
	case "Id":
		e.Id = e.core.Id.Get()
	case "Type":
		e.Type = e.core.Type.Get()
	case "Tooltip":
		e.Tooltip = e.core.Tooltip.Get()
	case "Icon":
		e.Icon = e.core.Icon.Get()
	case "Status":
		e.Status = e.core.Status.Get()
	case "QuickWindowViewable":
		e.QuickWindowViewable = e.core.QuickWindowViewable.Get()
	case "Allocation":
		r := e.core.Allocation.Get()
		e.Allocation = Rectangle{r[0].(int16), r[1].(int16), r[2].(uint16), r[3].(uint16)}
	case "Data":
		e.Data = e.core.Data.Get()
	case "Menu":
		e.Menu = e.core.Menu.Get()
	}
	dbus.NotifyChange(e, prop)
}
