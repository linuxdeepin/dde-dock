package dock

import "crypto/md5"
import "encoding/hex"

import "pkg.linuxdeepin.com/lib/dbus"

func (e *AppEntry) GetDBusInfo() dbus.DBusInfo {
	hasher := md5.New()
	hasher.Write([]byte(e.Id))
	// DBusObjectPath can't be start with digital number
	var id string
	id = "d" + hex.EncodeToString(hasher.Sum(nil))
	return dbus.DBusInfo{
		Dest:       "dde.dock.entry." + id,
		ObjectPath: "/dde/dock/entry/v1/" + id,
		Interface:  "dde.dock.Entry",
	}
}
