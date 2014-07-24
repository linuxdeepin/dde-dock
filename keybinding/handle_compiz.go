/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package keybinding

import (
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"strings"
)

var (
	coreSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_CORE,
		COMPIZ_SETTINGS_BASE_PATH+"core/")
	moveSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_MOVE,
		COMPIZ_SETTINGS_BASE_PATH+"move/")
	resizeSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_RESIZE,
		COMPIZ_SETTINGS_BASE_PATH+"resize/")
	vpswitchSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_VPSWITCH,
		COMPIZ_SETTINGS_BASE_PATH+"vpswitch/")
	putSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_PUT,
		COMPIZ_SETTINGS_BASE_PATH+"put/")
	wallSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_WALL,
		COMPIZ_SETTINGS_BASE_PATH+"wall/")
	shiftSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_SHIFT,
		COMPIZ_SETTINGS_BASE_PATH+"shift/")
	switcherSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_SWITCHER,
		COMPIZ_SETTINGS_BASE_PATH+"switcher/")
)

func formatCompizShortcut(shortcut string) string {
	Logger.Info("formatCompizShortcut:", shortcut)
	strs := strings.Split(shortcut, "-")
	l := len(strs)
	if l < 2 {
		return shortcut
	}

	tmp := ""
	for i := 0; i < l-1; i++ {
		tmp += "<" + strs[i] + ">"
	}
	tmp += strs[l-1]
	Logger.Info("formatCompizShortcut RET:", tmp)

	return tmp
}

func (m *Manager) setCompizSettings(id int32, key, value string) {
	k, ok := compizKeysMap[key]
	if !ok {
		Logger.Warningf("'%s' not in compizKeysMap", key)
		return
	}

	shortcut := formatCompizShortcut(value)
	if id >= 600 && id < 650 {
		coreSettings.SetString(k, shortcut)
	} else if id >= 650 && id < 700 {
		moveSettings.SetString(k, shortcut)
	} else if id >= 700 && id < 750 {
		resizeSettings.SetString(k, shortcut)
	} else if id >= 750 && id < 800 {
		vpswitchSettings.SetString(k, shortcut)
	} else if id >= 800 && id < 850 {
		putSettings.SetString(k, shortcut)
	} else if id >= 850 && id < 900 {
		wallSettings.SetString(k, shortcut)
	} else if id >= 900 && id < 950 {
		shiftSettings.SetString(k, shortcut)
	} else if id >= 950 && id < 1000 {
		switcherSettings.SetString(k, shortcut)
	}
}

func (m *Manager) listenCompizSettings() {
	coreSettings.Connect("changed", func(s *gio.Settings, key string) {
		shortcut := coreSettings.GetString(key)
		switch key {
		case "show-desktop-key":
			updateSystemSettings("show-desktop", shortcut)
		case "close-window-key":
			updateSystemSettings("close", shortcut)
		case "maximize-window-key":
			updateSystemSettings("maximize", shortcut)
		case "unmaximize-window-key":
			updateSystemSettings("unmaximize", shortcut)
		case "minimize-window-key":
			updateSystemSettings("minimize", shortcut)
		case "toggle-window-shaded-key":
			updateSystemSettings("toggle-shaded", shortcut)
		case "window-menu-key":
			updateSystemSettings("activate-window-menu", shortcut)
		}
	})

	moveSettings.Connect("changed", func(s *gio.Settings, key string) {
		shortcut := moveSettings.GetString(key)
		switch key {
		case "initiate-key":
			updateSystemSettings("begin-move", shortcut)
		}
	})

	resizeSettings.Connect("changed", func(s *gio.Settings, key string) {
		shortcut := resizeSettings.GetString(key)
		switch key {
		case "initiate-key":
			updateSystemSettings("begin-resize", shortcut)
		}
	})

	vpswitchSettings.Connect("changed", func(s *gio.Settings, key string) {
		shortcut := vpswitchSettings.GetString(key)
		switch key {
		case "switch-to-1-key":
			updateSystemSettings("switch-to-workspace-1", shortcut)
		case "switch-to-2-key":
			updateSystemSettings("switch-to-workspace-2", shortcut)
		case "switch-to-3-key":
			updateSystemSettings("switch-to-workspace-3", shortcut)
		case "switch-to-4-key":
			updateSystemSettings("switch-to-workspace-4", shortcut)
		}
	})

	putSettings.Connect("changed", func(s *gio.Settings, key string) {
		shortcut := putSettings.GetString(key)
		switch key {
		case "put-viewport-1-key":
			updateSystemSettings("put-viewport-1-key", shortcut)
		case "put-viewport-2-key":
			updateSystemSettings("put-viewport-2-key", shortcut)
		case "put-viewport-3-key":
			updateSystemSettings("put-viewport-3-key", shortcut)
		case "put-viewport-4-key":
			updateSystemSettings("put-viewport-4-key", shortcut)
		}
	})

	wallSettings.Connect("changed", func(s *gio.Settings, key string) {
		shortcut := wallSettings.GetString(key)
		switch key {
		case "left-key":
			updateSystemSettings("switch-to-workspace-left", shortcut)
		case "right-key":
			updateSystemSettings("switch-to-workspace-right", shortcut)
		case "up-key":
			updateSystemSettings("switch-to-workspace-up", shortcut)
		case "down-key":
			updateSystemSettings("switch-to-workspace-down", shortcut)
		case "left-window-key":
			updateSystemSettings("move-to-workspace-left", shortcut)
		case "right-window-key":
			updateSystemSettings("move-to-workspace-right", shortcut)
		case "up-window-key":
			updateSystemSettings("move-to-workspace-up", shortcut)
		case "down-window-key":
			updateSystemSettings("move-to-workspace-down", shortcut)
		}
	})

	shiftSettings.Connect("changed", func(s *gio.Settings, key string) {
		shortcut := shiftSettings.GetString(key)
		switch key {
		case "next-key":
			updateSystemSettings("next-key", shortcut)
		case "prev-key":
			updateSystemSettings("prev-key", shortcut)
		}
	})

	switcherSettings.Connect("changed", func(s *gio.Settings, key string) {
		shortcut := switcherSettings.GetString(key)
		switch key {
		case "next-key":
			updateSystemSettings("switch-applications", shortcut)
		case "prev-key":
			updateSystemSettings("switch-applications-backward", shortcut)
		}
	})
}
