/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
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
	systemActionsFile = "dde-daemon/keybinding/system_actions.json"
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
		"wm-switcher":           gettext.Tr("Switch window effects"),
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
		return "dbus-send --print-reply --dest=com.deepin.dde.Launcher /com/deepin/dde/Launcher com.deepin.dde.Launcher.Toggle"
	case "terminal":
		return "/usr/lib/deepin-daemon/default-terminal"
	case "lock-screen":
		return "dbus-send --print-reply --dest=com.deepin.dde.lockFront /com/deepin/dde/lockFront com.deepin.dde.lockFront.Show"
	case "show-dock":
		return "dbus-send --type=method_call --dest=com.deepin.daemon.Dock /dde/dock/HideStateManager dde.dock.HideStateManager.ToggleShow"
	case "logout":
		return "dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Ping"
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
		return "gvfs-open ~"
	case "disable-touchpad":
		return "gsettings set com.deepin.dde.touchpad touchpad-enabled false"
	case "wm-switcher":
		return "dbus-send --type=method_call --dest=com.deepin.wm_switcher /com/deepin/wm_switcher com.deepin.wm_switcher.requestSwitchWM"
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
	var file = path.Join("/usr/local/share", systemActionsFile)
	if dutils.IsFileExist(file) {
		return file
	}

	file = path.Join("/usr/share", systemActionsFile)
	if dutils.IsFileExist(file) {
		return file
	}

	return ""
}
