package main

import "crypto/md5"
import "encoding/hex"

import "dlib/dbus"
import "github.com/BurntSushi/xgb/xproto"
import "strconv"
import "strings"

func (e *AppEntry) setPropData(m map[xproto.Window]*WindowInfo) {
	for k, v := range m {
		title := v.Title
		key := strconv.FormatUint(uint64(k), 10)
		if value, ok := e.Data[key]; ok && value == title {
			continue
		}
		e.Data[key] = title
		dbus.NotifyChange(e, "Data")
	}
}

func (e *AppEntry) setPropTooltip(v string) {
	if e.Tooltip != v {
		e.Tooltip = v
		dbus.NotifyChange(e, "Tooltip")
	}
}
func (e *AppEntry) setPropIcon(v string) {
	if e.Icon != v {
		e.Icon = v
		dbus.NotifyChange(e, "Icon")
	}
}
func (e *AppEntry) setPropMenu(v string) {
	if e.Menu != v {
		e.Menu = v
		dbus.NotifyChange(e, "Menu")
	}
}
func (e *AppEntry) setPropStatus(v int32) {
	if e.Status != v {
		e.Status = v
		dbus.NotifyChange(e, "Status")
	}
}
func (e *AppEntry) setPropQuickWindowViewable(v bool) {
	if e.QuickWindowViewable != v {
		e.QuickWindowViewable = v
		dbus.NotifyChange(e, "QuickWindowViewable")
	}
}
func (e *AppEntry) setPropAllocation(v Rectangle) {
	if e.Allocation != v {
		e.Allocation = v
		dbus.NotifyChange(e, "Allocation")
	}
}

func (e *AppEntry) GetDBusInfo() dbus.DBusInfo {
	hasher := md5.New()
	hasher.Write([]byte(e.Id))
	// DBusObjectPath can't be start with digital number
	var id string
	if false {
		id = "d" + hex.EncodeToString(hasher.Sum(nil))
	} else {
		id = strings.Replace(e.Id, "-", "", -1)
		id = strings.Replace(id, ".", "", -1)
	}
	return dbus.DBusInfo{
		"dde.dock.entry." + id,
		"/dde/dock/entry/v1/" + id,
		"dde.dock.Entry",
	}
}
