package main

import "dlib/dbus"
import "dlib/gio-2.0"
import "crypto/md5"
import "encoding/hex"

//import "strings"

type Rectangle struct {
	X, Y          int16
	Width, Height uint16
}

type DesktopEntry struct {
	core *gio.DesktopAppInfo
	ID   string
	Type string

	Tooltip string
	Icon    string

	Status int32 `Actived/Normal/`

	QuickWindowVieable bool
	Allocation         Rectangle

	Data map[string]string
}

func NewDesktopEntry(id string) *DesktopEntry {
	e := &DesktopEntry{}
	e.core = gio.NewDesktopAppInfo(id)
	if e.core == nil {
		return nil
	}
	e.Tooltip = e.core.GetName()
	e.ID = e.core.GetId()
	e.Icon = e.core.GetIcon().ToString()
	return e
}

func (e *DesktopEntry) QuickWindow(x, y int32) {}
func (e *DesktopEntry) ContextMenu(x, y int32) {}
func (e *DesktopEntry) Activate(x, y int32) {
	e.core.Launch(nil, nil)
}
func (e *DesktopEntry) SecondaryActivate(x, y int32)        {}
func (e *DesktopEntry) OnDragEnter(x, y int32, data string) {}
func (e *DesktopEntry) OnDragLeave(x, y int32, data string) {}
func (e *DesktopEntry) OnDragOver(x, y int32, data string)  {}
func (e *DesktopEntry) OnDragDrop(x, y int32, data string)  {}

func (e *DesktopEntry) GetDBusInfo() dbus.DBusInfo {
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

func main() {
	for _, id := range loadAll() {
		if e := NewDesktopEntry(id + ".desktop"); e != nil {
			dbus.InstallOnSession(e)
		}
	}
	dbus.Wait()
}
