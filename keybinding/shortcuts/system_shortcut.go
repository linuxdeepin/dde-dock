/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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
	return defaultSysActionCmdMap[id]
}

// key is id, value is commandline.
var defaultSysActionCmdMap = map[string]string{
	"launcher":               "dbus-send --print-reply --dest=com.deepin.dde.Launcher /com/deepin/dde/Launcher com.deepin.dde.Launcher.Toggle",
	"terminal":               "/usr/lib/deepin-daemon/default-terminal",
	"terminal-quake":         "deepin-terminal --quake-mode",
	"lock-screen":            "dbus-send --print-reply --dest=com.deepin.dde.lockFront /com/deepin/dde/lockFront com.deepin.dde.lockFront.Show",
	"logout":                 "dbus-send --print-reply --dest=com.deepin.dde.shutdownFront /com/deepin/dde/shutdownFront com.deepin.dde.shutdownFront.Show",
	"deepin-screen-recorder": "dbus-send --print-reply --dest=com.deepin.ScreenRecorder /com/deepin/ScreenRecorder com.deepin.ScreenRecorder.stopRecord",
	"system-monitor":         "/usr/bin/deepin-system-monitor",
	"color-picker":           "dbus-send --print-reply --dest=com.deepin.Picker /com/deepin/Picker com.deepin.Picker.Show",
	// screenshot actions:
	"screenshot":            screenshotCmdPrefix + "StartScreenshot",
	"screenshot-fullscreen": screenshotCmdPrefix + "FullscreenScreenshot",
	"screenshot-window":     screenshotCmdPrefix + "TopWindowScreenshot",
	"screenshot-delayed":    screenshotCmdPrefix + "DelayScreenshot int64:5",
	"file-manager":          "/usr/lib/deepin-daemon/default-file-manager",
	"disable-touchpad":      "gsettings set com.deepin.dde.touchpad touchpad-enabled false",
	"wm-switcher":           "dbus-send --type=method_call --dest=com.deepin.WMSwitcher /com/deepin/WMSwitcher com.deepin.WMSwitcher.RequestSwitchWM",
	"turn-off-screen":       "sleep 0.5; xset dpms force off",
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
