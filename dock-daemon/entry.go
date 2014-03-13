package main

import (
	"dlib/dbus"
)

const entryDestPrefix = "dde.dock.entry."
const entryPathPrefix = "/dde/dock/entry/v1/"

type Rectangle struct {
	X, Y          int16
	Width, Height uint16
}

type EntryProxyer struct {
	entryId string
	core    *RemoteEntry

	ID   string `dmusic`
	Type string `applet/other`

	Tooltip string
	Icon    string

	Status int32 `Actived/Normal/`

	QuickWindowVieable bool
	Allocation         Rectangle
}

func NewEntryProxyer(entryId string) (e *EntryProxyer, err error) {
	e = &EntryProxyer{}
	core, err := NewRemoteEntry(entryDestPrefix+entryId, dbus.ObjectPath(entryPathPrefix+entryId))
	if err != nil {
		return
	}

	// init properties
	e.entryId = entryId
	e.core = core
	e.ID = e.core.ID.Get()
	e.Type = e.core.Type.Get()
	e.Tooltip = e.core.Tooltip.Get()
	e.Icon = e.core.Icon.Get()
	e.Status = e.core.Status.Get()
	e.QuickWindowVieable = e.core.QuickWindowVieable.Get()
	r := e.core.Allocation.Get()
	e.Allocation = Rectangle{r[0].(int16), r[1].(int16), r[2].(uint16), r[3].(uint16)}

	// monitor properties changed
	e.core.ID.ConnectChanged(func() {
		dbus.NotifyChange(e, "ID")
	})
	e.core.Type.ConnectChanged(func() {
		dbus.NotifyChange(e, "Type")
	})
	e.core.Tooltip.ConnectChanged(func() {
		dbus.NotifyChange(e, "Tooltip")
	})
	e.core.Icon.ConnectChanged(func() {
		dbus.NotifyChange(e, "Icon")
	})
	e.core.Status.ConnectChanged(func() {
		dbus.NotifyChange(e, "Status")
	})
	e.core.QuickWindowVieable.ConnectChanged(func() {
		dbus.NotifyChange(e, "QuickWindowVieable")
	})
	e.core.Allocation.ConnectChanged(func() {
		dbus.NotifyChange(e, "Allocation")
	})

	return
}

func (e *EntryProxyer) QuickWindow(x, y int32)              { e.core.QuickWindow(x, y) }
func (e *EntryProxyer) ContextMenu(x, y int32)              { e.core.ContextMenu(x, y) }
func (e *EntryProxyer) Activate(x, y int32)                 { e.core.Activate(x, y) }
func (e *EntryProxyer) SecondaryActivate(x, y int32)        { e.core.SecondaryActivate(x, y) }
func (e *EntryProxyer) OnDragEnter(x, y int32, data string) { e.core.OnDragEnter(x, y, data) }
func (e *EntryProxyer) OnDragLeave(x, y int32, data string) { e.core.OnDragLeave(x, y, data) }
func (e *EntryProxyer) OnDragOver(x, y int32, data string)  { e.core.OnDragOver(x, y, data) }
func (e *EntryProxyer) OnDragDrop(x, y int32, data string)  { e.core.OnDragDrop(x, y, data) }
