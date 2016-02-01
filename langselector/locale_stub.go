/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package langselector

import (
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusPath      = "/com/deepin/daemon/LangSelector"
	dbusInterface = "com.deepin.daemon.LangSelector"
)

func (lang *LangSelector) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: dbusPath,
		Interface:  dbusInterface,
	}
}

func (lang *LangSelector) setPropCurrentLocale(locale string) {
	if lang.CurrentLocale != locale {
		lang.CurrentLocale = locale
		dbus.NotifyChange(lang, "CurrentLocale")
	}
}
