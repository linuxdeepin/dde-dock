/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import "crypto/md5"
import "encoding/hex"

import "pkg.deepin.io/lib/dbus"

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
