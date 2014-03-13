package main

import "dlib/dbus"
import "fmt"
import "os"

type Rectangle struct {
	X, Y          int16
	Width, Height uint16
}

type Entry struct {
	ID   string `dmusic`
	Type string `applet/other`

	Tooltip string
	Icon    string

	Status int32 `Actived/Normal/`

	QuickWindowVieable bool
	Allocation         Rectangle
}

func (*Entry) QuickWindow(x, y int32)              {}
func (*Entry) ContextMenu(x, y int32)              {}
func (*Entry) Activate(x, y int32)                 {}
func (*Entry) SecondaryActivate(x, y int32)        {}
func (*Entry) OnDragEnter(x, y int32, data string) {}
func (*Entry) OnDragLeave(x, y int32, data string) {}
func (*Entry) OnDragOver(x, y int32, data string)  {}
func (*Entry) OnDragDrop(x, y int32, data string)  {}

func (*Entry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		fmt.Sprintf("dde.dock.entry.app.chrome-%d", os.Getpid()),
		"/dde/dock/entry/v1",
		"dde.dock.Entry",
	}
}

func main() {
	dbus.InstallOnSession(&Entry{Allocation: Rectangle{0, 0, 300, 400}})
	dbus.Wait()
}
