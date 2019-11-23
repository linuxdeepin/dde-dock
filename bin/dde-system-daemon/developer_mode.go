/*
 * Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
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

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

type developerModeInfo struct {
	Enabled bool `json:"enabled"`

	changed bool
}

const (
	developerModeFile = "/var/lib/deepin/developer_mode.json"
)

var (
	_developerModeLocker sync.RWMutex
)

func (info *developerModeInfo) Enable(v bool) {
	if info.Enabled == v {
		return
	}
	info.Enabled = v
	info.changed = true
}

func (info *developerModeInfo) Save(filename string) error {
	if !info.changed {
		return nil
	}

	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	_developerModeLocker.Lock()
	defer _developerModeLocker.Unlock()
	err = os.MkdirAll(filepath.Dir(developerModeFile), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func (*Daemon) IsDeveloperMode() (bool, *dbus.Error) {
	info, _ := newDeveloperModeInfo(developerModeFile)
	if info == nil {
		info = &developerModeInfo{Enabled: false}
	}
	return info.Enabled, nil
}

func (*Daemon) EnableDeveloperMode() *dbus.Error {
	info, _ := newDeveloperModeInfo(developerModeFile)
	if info == nil {
		info = &developerModeInfo{Enabled: false}
	}
	if info.Enabled {
		return nil
	}

	info.Enable(true)
	err := info.Save(developerModeFile)
	if err != nil {
		return dbusutil.ToError(err)
	}
	// TODO(jouyouyun): do not delete file, implement in kernel
	modifyFileAttr(developerModeFile, "+i")
	return nil
}

func newDeveloperModeInfo(filename string) (*developerModeInfo, error) {
	_developerModeLocker.RLock()
	defer _developerModeLocker.RUnlock()
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var info developerModeInfo
	err = json.Unmarshal(contents, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func modifyFileAttr(filename, attr string) {
	cmdFiles := []string{
		"/usr/bin/chattr",
		"/bin/chattr",
	}
	for _, cmd := range cmdFiles {
		outs, err := exec.Command(cmd, attr, filename).CombinedOutput()
		if err != nil {
			logger.Warning("Failed to modify file attr:", attr, string(outs), err)
			continue
		}
		return
	}
}
