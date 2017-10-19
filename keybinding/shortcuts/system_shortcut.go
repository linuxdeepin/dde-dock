/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package shortcuts

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"sync"

	dutils "pkg.deepin.io/lib/utils"
)

const (
	// Under '/usr/share' or '/usr/local/share'
	systemActionsFile   = "dde-daemon/keybinding/system_actions.json"
	screenshotCmdPrefix = "dbus-send --print-reply --dest=com.deepin.Screenshot /com/deepin/Screenshot com.deepin.Screenshot."
)

type SystemShortcut struct {
	*GSettingsShortcut
	arg *ActionExecCmdArg
}

func (ss *SystemShortcut) SetName(name string) error {
	return ErrOpNotSupported
}

func (ss *SystemShortcut) GetAction() *Action {
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

var loadSysActionsFileOnce sync.Once
var actionsCache *actionHandler

func getSystemActionCmd(id string) string {
	loadSysActionsFileOnce.Do(func() {
		file := getSystemActionsFile()
		actions, err := loadSystemActionsFile(file)
		if err != nil {
			logger.Warning("failed to load system actions file:", err)
			return
		}
		actionsCache = actions
	})

	if actionsCache != nil {
		if cmd, ok := actionsCache.getCmd(id); ok {
			return cmd
		}
	}
	return getDefaultSysActionCmd(id)
}

func getDefaultSysActionCmd(id string) string {
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
	case "deepin-screen-recorder":
		return "/usr/bin/deepin-screen-recorder"
		// screenshot actions:
	case "screenshot":
		return screenshotCmdPrefix + "StartScreenshot"
	case "screenshot-fullscreen":
		return screenshotCmdPrefix + "FullscreenScreenshot"
	case "screenshot-window":
		return screenshotCmdPrefix + "TopWindowScreenshot"
	case "screenshot-delayed":
		return screenshotCmdPrefix + "DelayScreenshot int64:5"

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

func (a *actionHandler) getCmd(id string) (cmd string, ok bool) {
	for _, v := range a.Actions {
		if v.Key == id {
			return v.Action, true
		}
	}
	return "", false
}

func loadSystemActionsFile(file string) (*actionHandler, error) {
	logger.Debug("load system action file:", file)

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
