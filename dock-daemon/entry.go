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
	entryId    string
	destPath   string
	objectPath dbus.ObjectPath
	core       *RemoteEntry

	Id   string `dmusic`
	Type string `applet/other`

	Tooltip string
	Icon    string
	Menu    string

	Status int32 `Actived/Normal/`

	QuickWindowViewable bool
	Allocation          Rectangle

	Data map[string]string
}
