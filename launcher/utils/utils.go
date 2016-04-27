/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package utils

import (
	"dbus/org/freedesktop/notifications"
	"gir/gio-2.0"
	"os"
	"strings"
)

// dir default perm.
const (
	DirDefaultPerm os.FileMode = 0755

	dbusNotifyDest = "org.freedesktop.Notifications"
	dbusNotifyPath = "/org/freedesktop/Notifications"
)

// CreateDesktopAppInfo is a helper function for creating GDesktopAppInfo object.
// if name is a path, gio.NewDesktopAppInfoFromFilename is used.
// otherwise, name must be desktop id and gio.NewDesktopAppInfo is used.
func CreateDesktopAppInfo(name string) *gio.DesktopAppInfo {
	if strings.ContainsRune(name, os.PathSeparator) {
		return gio.NewDesktopAppInfoFromFilename(name)
	} else {
		return gio.NewDesktopAppInfo(name)
	}
}

func Notify(icon, summary, body string) {
	notifier, err := notifications.NewNotifier(dbusNotifyDest, dbusNotifyPath)
	if err != nil {
		return
	}

	go func() {
		notifier.Notify("Launcher", 0, icon, summary, body, nil, nil, 0)
		notifications.DestroyNotifier(notifier)
	}()

	return
}
