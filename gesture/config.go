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

package gesture

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"pkg.deepin.io/lib/xdg/basedir"
)

const (
	ActionTypeShortcut    = "shortcut"
	ActionTypeCommandline = "commandline"
	ActionTypeBuiltin     = "built-in"
)

var (
	configSystemPath = "/usr/share/dde-daemon/gesture.json"
	configUserPath   = filepath.Join(basedir.GetUserConfigDir(), "deepin/dde-daemon/gesture.json")

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
		return fmt.Errorf("not found gesture info for: %s, %s, %d", name, direction, fingers)
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
		return nil, fmt.Errorf("file '%s' is empty", filename)
	}

	var infos gestureInfos
	err = json.Unmarshal(content, &infos)
	if err != nil {
		return nil, err
	}
	return infos, nil
}
