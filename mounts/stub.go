/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mounts

import (
	"encoding/json"
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest = "com.deepin.daemon.DiskMount"
	dbusPath = "/com/deepin/daemon/DiskMount"
	dbusIFC  = "com.deepin.daemon.DiskMount"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) setPropDiskList(infos DiskInfos) {
	if toJSON(m.DiskList) == toJSON(infos) {
		return
	}

	m.DiskList = infos
	dbus.NotifyChange(m, "DiskList")
}

func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
