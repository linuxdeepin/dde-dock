/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package gesture

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	dutils "pkg.deepin.io/lib/utils"
	"sync"
)

const (
	ActionTypeShortcut    string = "shortcut"
	ActionTypeCommandline        = "commandline"
)

var (
	configSystemPath = "/usr/share/dde-daemon/gesture_config.json"
	configUserPath   = os.Getenv("HOME") + "/.config/deepin/dde-daemon/gesture_config.json"
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
	filename string
	Infos    gestureInfos
}

func newGestureManager() (*gestureManager, error) {
	var filename = configUserPath
	if !dutils.IsFileExist(configUserPath) {
		err := dutils.CopyFile(configSystemPath, configUserPath)
		if err != nil {
			filename = configSystemPath
		}
	}

	infos, err := newGestureInfosFromFile(filename)
	if err != nil {
		return nil, err
	}
	return &gestureManager{
		filename: filename,
		Infos:    infos,
	}, nil
}

func (m *gestureManager) Exec(name, direction string, fingers int32) error {
	m.locker.RLock()
	defer m.locker.RUnlock()
	info := m.Infos.Get(name, direction, fingers)
	if info == nil {
		return fmt.Errorf("Not found gesture info for: %s, %s, %d", name, direction, fingers)
	}

	var cmd = info.Action.Action
	if info.Action.Type == ActionTypeShortcut {
		cmd = fmt.Sprintf("xdotool key %s", cmd)
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
	data, err := json.Marshal(m.Infos)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.filename, data, 0644)
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

	var infos gestureInfos
	err = json.Unmarshal(content, &infos)
	if err != nil {
		return nil, err
	}
	return infos, nil
}
