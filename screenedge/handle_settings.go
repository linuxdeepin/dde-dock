/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package screenedge

import (
	"gir/gio-2.0"
)

func handleSettingsChanged() {
	setting := zoneSettings()
	setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case "left-up":
			edgeActionMap[leftTopEdge] = setting.GetString(key)
		case "left-down":
			edgeActionMap[leftBottomEdge] = setting.GetString(key)
		case "right-up":
			edgeActionMap[rightTopEdge] = setting.GetString(key)
		case "right-down":
			edgeActionMap[rightBottomEdge] = setting.GetString(key)
		}
	})
	setting.GetString("left-up")
}
