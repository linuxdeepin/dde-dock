/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
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

package shortcuts

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"pkg.deepin.io/lib/gettext"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	// Under '/usr/share' or '/usr/local/share'
	systemActionsFile = "deepin/dde-daemon/keybinding/system_actions.json"
)

func ListSystemShortcut() Shortcuts {
	s := newSystemGSetting()
	defer s.Unref()
	return doListShortcut(s, systemIdNameMap(), KeyTypeSystem)
}

func resetSystemAccels() {
	s := newSystemGSetting()
	defer s.Unref()
	doResetAccels(s)
}

func disableSystemAccels(key string) {
	s := newSystemGSetting()
	defer s.Unref()
	doDisableAccles(s, key)
}

func addSystemAccel(key, accel string) {
	s := newSystemGSetting()
	defer s.Unref()
	doAddAccel(s, key, accel)
}

func delSystemAccel(key, accel string) {
	s := newSystemGSetting()
	defer s.Unref()
	doDelAccel(s, key, accel)
}

func systemIdNameMap() map[string]string {
	var idNameMap = map[string]string{
		"launcher":              gettext.Tr("Launcher"),
		"terminal":              gettext.Tr("Terminal"),
		"lock-screen":           gettext.Tr("Lock screen"),
		"show-dock":             gettext.Tr("Show/Hide the dock"),
		"logout":                gettext.Tr("Logout"),
		"terminal-quake":        gettext.Tr("Terminal Quake Window"),
		"screenshot":            gettext.Tr("Screenshot"),
		"screenshot-fullscreen": gettext.Tr("Full screenshot"),
		"screenshot-window":     gettext.Tr("Window screenshot"),
		"screenshot-delayed":    gettext.Tr("Delay screenshot"),
		"file-manager":          gettext.Tr("File manager"),
		"disable-touchpad":      gettext.Tr("Disable Touchpad"),
		"switch-layout":         gettext.Tr("Switch Layout"),
	}

	return idNameMap
}

func getSystemAction(id string) string {
	file := getSystemActionsFile()
	handler, err := getActionHandler(file)
	if err != nil {
		return findSysActionInTable(id)
	}

	for _, v := range handler.Actions {
		if v.Key == id {
			return v.Action
		}
	}

	return findSysActionInTable(id)
}

func findSysActionInTable(id string) string {
	switch id {
	case "launcher":
		return "dde-launcher"
	case "terminal":
		return "default-terminal"
	case "lock-screen":
		return "dde-lock"
	case "show-dock":
		return "dbus-send --type=method_call --dest=com.deepin.daemon.Dock /dde/dock/HideStateManager dde.dock.HideStateManager.ToggleShow"
	case "logout":
		return "dde-shutdown"
	case "terminal-quake":
		return "deepin-terminal --quake-mode"
	case "screenshot":
		return "deepin-screenshot"
	case "screenshot-fullscreen":
		return "deepin-screenshot -f"
	case "screenshot-window":
		return "deepin-screenshot -w"
	case "screenshot-delayed":
		return "deepin-screenshot -d 5"
	case "file-manager":
		return "nautilus"
	case "disable-touchpad":
		return "gsettings set com.deepin.dde.touchpad touchpad-enabled false"
	}

	return ""
}

type actionHandler struct {
	Actions []struct {
		Key    string `json:"Key"`
		Action string `json:"Action"`
	} `json:"Actions"`
}

func getActionHandler(file string) (*actionHandler, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var handler actionHandler
	err = json.Unmarshal(content, &handler)
	if err != nil {
		return nil, err
	}

	return &handler, nil
}

func getSystemActionsFile() string {
	var file = path.Join("/usr/share", systemActionsFile)
	if dutils.IsFileExist(file) {
		return file
	}

	file = path.Join("/usr/local/share", systemActionsFile)
	if dutils.IsFileExist(file) {
		return file
	}

	return ""
}
