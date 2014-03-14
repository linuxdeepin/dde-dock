package main

import "dlib/dbus"
import "crypto/md5"
import "encoding/hex"

//import "strings"

type Rectangle struct {
	X, Y          int16
	Width, Height uint16
}

type AppEntry struct {
	nApp *NormalApp
	rApp *RuntimeApp
	ID   string
	Type string

	Tooltip string
	Icon    string

	Status int32 `Actived/Normal/`

	QuickWindowVieable bool
	Allocation         Rectangle

	Data map[string]string
}

func NewAppEntry(id string) *AppEntry {
	e := &AppEntry{}
	return e
}

func (e *AppEntry) QuickWindow(x, y int32) {}
func (e *AppEntry) ContextMenu(x, y int32) {}
func (e *AppEntry) Activate(x, y int32) {
}
func (e *AppEntry) SecondaryActivate(x, y int32)        {}
func (e *AppEntry) OnDragEnter(x, y int32, data string) {}
func (e *AppEntry) OnDragLeave(x, y int32, data string) {}
func (e *AppEntry) OnDragOver(x, y int32, data string)  {}
func (e *AppEntry) OnDragDrop(x, y int32, data string)  {}

func (e *AppEntry) GetDBusInfo() dbus.DBusInfo {
	hasher := md5.New()
	hasher.Write([]byte(e.ID))
	// DBusObjectPath can't be start with digital number
	id := "d" + hex.EncodeToString(hasher.Sum(nil))
	return dbus.DBusInfo{
		"dde.dock.entry." + id,
		"/dde/dock/entry/v1/" + id,
		"dde.dock.Entry",
	}
}
