/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package shortcuts

import (
	"pkg.deepin.io/lib/gettext"
)

func ListMetacityShortcut() Shortcuts {
	gala := newGalaGSetting()
	if gala != nil {
		defer gala.Unref()
		return doListShortcut(gala, metacityIdNameMap(), KeyTypeMetacity)
	}

	meta := newMetacityGSetting()
	if meta == nil {
		return nil
	}
	defer meta.Unref()
	return doListShortcut(meta, metacityIdNameMap(), KeyTypeMetacity)
}

func resetMetacityAccels() {
	meta := newMetacityGSetting()
	if meta != nil {
		defer meta.Unref()
		doResetAccels(meta)
	}

	gala := newGalaGSetting()
	if gala == nil {
		return
	}
	defer gala.Unref()
	doResetAccels(gala)
}

func disableMetacityAccels(key string) {
	meta := newMetacityGSetting()
	if meta != nil {
		defer meta.Unref()
		doDisableAccles(meta, key)
	}

	gala := newGalaGSetting()
	if gala == nil {
		return
	}
	defer gala.Unref()
	doDisableAccles(gala, key)
}

func addMetacityAccel(key, accel string) {
	meta := newMetacityGSetting()
	if meta != nil {
		defer meta.Unref()
		doAddAccel(meta, key, accel)
	}

	gala := newGalaGSetting()
	if gala == nil {
		return
	}
	defer gala.Unref()
	doAddAccel(gala, key, accel)
}

func delMetacityAccel(key, accel string) {
	meta := newMetacityGSetting()
	if meta != nil {
		defer meta.Unref()
		doDelAccel(meta, key, accel)
	}

	gala := newGalaGSetting()
	if gala == nil {
		return
	}
	defer gala.Unref()
	doDelAccel(gala, key, accel)
}

func metacityIdNameMap() map[string]string {
	var idNameMap = map[string]string{
		"expose-all-windows": gettext.Tr("Display windows of all workspaces"),
		"expose-windows":     gettext.Tr("Display windows of current workspace"),
		"preview-workspace":  gettext.Tr("Display workspace"),
	}

	return idNameMap
}
