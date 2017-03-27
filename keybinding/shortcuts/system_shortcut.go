/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
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
	dutils "pkg.deepin.io/lib/utils"
)

const (
	// Under '/usr/share' or '/usr/local/share'
	systemActionsFile = "dde-daemon/keybinding/system_actions.json"
)

type SystemShortcut struct {
	*GSettingsShortcut
	arg *ActionExecCmdArg
}

func (ss *SystemShortcut) SetName(name string) error {
	return ErrOpNotSupported
}

func (ss *SystemShortcut) GetAction() *Action {
	if ss.GetId() == "switch-layout" {
		return &Action{
			Type: ActionTypeSwitchKbdLayout,
		}
	}

	return &Action{
		Type: ActionTypeExecCmd,
		Arg:  ss.arg,
	}
}

func (ss *SystemShortcut) SetAction(newAction *Action) error {
	if newAction == nil {
		return ErrNilAction
	}
	if newAction.Type != ActionTypeExecCmd {
		return ErrInvalidActionType
	}

	arg, ok := newAction.Arg.(*ActionExecCmdArg)
	if !ok {
		return ErrTypeAssertionFail
	}
	ss.arg = arg
	return nil
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
		return "dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show"
	case "terminal-quake":
		return "deepin-terminal --quake-mode"
	case "screenshot":
		return "deepin-screenshot"
	case "screenshot-fullscreen":
		return "deepin-screenshot -f"
	case "deepin-screen-recorder":
		return "/usr/bin/deepin-screen-recorder"
	case "screenshot-window":
		return "deepin-screenshot -w"
	case "screenshot-delayed":
		return "deepin-screenshot -d 5"
	case "file-manager":
		return "/usr/lib/deepin-daemon/default-file-manager"
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
