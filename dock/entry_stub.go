package dock

import "crypto/md5"
import "encoding/hex"

import "dlib/dbus"
import "strings"

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
