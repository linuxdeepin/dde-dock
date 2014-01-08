/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package main

import (
	"dlib/dbus/property"
	"dlib/gio-2.0"
)

type ThemeInfo struct {
        Name string
        Type string // 'system' or 'custom'
}

type Manager struct {
	GtkTheme       *property.GSettingsStringProperty `access:"readwrite"`
	IconTheme      *property.GSettingsStringProperty `access:"readwrite"`
	FontTheme      *property.GSettingsStringProperty `access:"readwrite"`
	CursorTheme    *property.GSettingsStringProperty `access:"readwrite"`
	BackgroundFile *property.GSettingsStringProperty `access:"readwrite"`
	AutoSwitch     *property.GSettingsBoolProperty   `access:"readwrite"`
	SwitchDuration *property.GSettingsIntProperty    `access:"readwrite"`
	CrossFadeMode  *property.GSettingsStringProperty `access:"readwrite"`
	CrossInterval  *property.GSettingsIntProperty    `access:"readwrite"`

        AvailableGtkTheme []ThemeInfo
        AvailableIconTheme []ThemeInfo
        AvailableFontTheme []ThemeInfo
        AvailableCursorTheme []ThemeInfo
        AvailableBackground []ThemeInfo

	isAutoSwitch   bool
	quitAutoSwitch chan bool
}

const (
	MANAGER_DEST = "com.deepin.daemon.Individuation"
	MANAGER_PATH = "/com/deepin/daemon/Individuation"
	MANAGER_IFC  = "com.deepin.daemon.Individuation"

	GSD_SCHEMA_ID      = "org.gnome.desktop.interface"
	INDIVIDUATE_ID     = "com.deepin.dde.individuate"
	DEFAULT_BG_PICTURE = "/usr/share/backgrounds/default_background.jpg"

	SCHEMA_KEY_URIS           = "picture-uris"
	SCHEMA_KEY_INDEX          = "index"
	SCHEMA_KEY_CUR_PICT       = "current-picture"
	SCHEMA_KEY_CROSS_INTERVAL = "cross-fade-interval"
	SCHEMA_KEY_CROSS_MODE     = "cross-fade-auto-mode"
	SCHEMA_KEY_AUTO_SWITCH    = "auto-switch"
	SCHEMA_KEY_DURATION       = "background-duration"
	SCHEMA_KEY_GTK            = "gtk-theme"
	SCHEMA_KEY_ICON           = "icon-theme"
	SCHEMA_KEY_FONT           = "font-theme"
	SCHEMA_KEY_CURSOR         = "cursor-theme"

	DACCOUNTS_USER_PATH = "/com/deepin/daemon/Accounts/User"
)

var (
	gsdSettings     = gio.NewSettings(GSD_SCHEMA_ID)
	indiviGSettings = gio.NewSettings(INDIVIDUATE_ID)
)
