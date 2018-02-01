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

package gesture

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sync"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/gsettings"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	ActionTypeShortcut    string = "shortcut"
	ActionTypeCommandline        = "commandline"
	ActionTypeBuiltin            = "built-in"
)

var (
	configSystemPath = "/usr/share/dde-daemon/gesture.json"
	configUserPath   = os.Getenv("HOME") + "/.config/deepin/dde-daemon/gesture.json"

	gestureSchemaId = "com.deepin.dde.gesture"
	gsKeyEnabled    = "enabled"
)

type ActionInfo struct {
	Type   string
	Action string
}

type gestureInfo struct {
	Name      string
	Direction string
	Fingers   int32
	Action    ActionInfo
}
type gestureInfos []*gestureInfo

type gestureManager struct {
	locker   sync.RWMutex
	userFile string
	Infos    gestureInfos

	setting *gio.Settings
	enabled bool
}

func newGestureManager() (*gestureManager, error) {
	var filename = configUserPath
	if !dutils.IsFileExist(configUserPath) {
		filename = configSystemPath
	}

	infos, err := newGestureInfosFromFile(filename)
	if err != nil {
		return nil, err
	}

	setting, err := dutils.CheckAndNewGSettings(gestureSchemaId)
	if err != nil {
		return nil, err
	}

	return &gestureManager{
		userFile: configUserPath,
		Infos:    infos,
		setting:  setting,
		enabled:  setting.GetBoolean(gsKeyEnabled),
	}, nil
}

func (m *gestureManager) Exec(name, direction string, fingers int32) error {
	m.locker.RLock()
	defer m.locker.RUnlock()

	if !m.enabled {
		logger.Debug("Gesture had been disabled")
		return nil
	}

	info := m.Infos.Get(name, direction, fingers)
	if info == nil {
		return fmt.Errorf("Not found gesture info for: %s, %s, %d", name, direction, fingers)
	}

	logger.Debug("[Exec] action info:", info.Name, info.Direction, info.Fingers,
		info.Action.Type, info.Action.Action)
	if isKbdAlreadyGrabbed() {
		return fmt.Errorf("There has some proccess grabed keyboard, not exec action")
	}
	var cmd = info.Action.Action
	switch info.Action.Type {
	case ActionTypeCommandline:
		break
	case ActionTypeShortcut:
		cmd = fmt.Sprintf("xdotool key %s", cmd)
		break
	case ActionTypeBuiltin:
		f, ok := builtinSets[cmd]
		if !ok {
			return fmt.Errorf("Invalid built-in action: %s", cmd)
		}
		return f()
	default:
		return fmt.Errorf("Invalid action type: %s", info.Action.Type)
	}

	out, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(out))
	}
	return nil
}

func (m *gestureManager) Write() error {
	m.locker.Lock()
	defer m.locker.Unlock()
	err := os.MkdirAll(path.Dir(m.userFile), 0755)
	if err != nil {
		return err
	}
	data, err := json.Marshal(m.Infos)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.userFile, data, 0644)
}

func (m *gestureManager) handleGSettingsChanged() {
	gsettings.ConnectChanged(gestureSchemaId, gsKeyEnabled, func(key string) {
		m.enabled = m.setting.GetBoolean(key)
	})
}

func (infos gestureInfos) Get(name, direction string, fingers int32) *gestureInfo {
	for _, info := range infos {
		if info.Name == name && info.Direction == direction &&
			info.Fingers == fingers {
			return info
		}
	}
	return nil
}

func (infos gestureInfos) Set(name, direction string, fingers int32, action ActionInfo) error {
	info := infos.Get(name, direction, fingers)
	if info == nil {
		return fmt.Errorf("Not found gesture info for: %s, %s, %d", name, direction, fingers)
	}
	info.Action = action
	return nil
}

func newGestureInfosFromFile(filename string) (gestureInfos, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("File '%s' is empty", filename)
	}

	var infos gestureInfos
	err = json.Unmarshal(content, &infos)
	if err != nil {
		return nil, err
	}
	return infos, nil
}
