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

package mime

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type userAppInfo struct {
	DesktopId     string   `json:"DesktopId"`
	SupportedMime []string `json:"SupportedMime"`
}

type userAppInfos []*userAppInfo

type userAppManager struct {
	appInfos userAppInfos
	filename string
	locker   sync.RWMutex
}

var (
	userAppFile = path.Join(os.Getenv("HOME"), ".config/deepin/dde-daemon/user_mime.json")
)

func (m *userAppManager) Get(mime string) userAppInfos {
	m.locker.RLock()
	defer m.locker.RUnlock()
	var ret userAppInfos
	for _, info := range m.appInfos {
		if info.HasMime(mime) {
			ret = append(ret, info)
		}
	}
	return ret
}

func (m *userAppManager) Add(mimes []string, desktopId string) bool {
	m.locker.Lock()
	defer m.locker.Unlock()
	var (
		found bool = false
		added bool = false
	)
	for _, info := range m.appInfos {
		if info.DesktopId == desktopId {
			found = true
			added = info.AddMimes(mimes)
			break
		}
	}
	if !found {
		added = true
		m.appInfos = append(m.appInfos, &userAppInfo{
			DesktopId:     desktopId,
			SupportedMime: mimes,
		})
	}
	return added
}

func (m *userAppManager) Delete(desktopId string) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	var (
		infos   userAppInfos
		deleted bool = false
	)
	for _, info := range m.appInfos {
		if info.DesktopId == desktopId {
			deleted = true
			continue
		}
		infos = append(infos, info)
	}
	if !deleted {
		return fmt.Errorf("Not found the application: %s", desktopId)
	}
	m.appInfos = infos
	return nil
}

func (m *userAppManager) Write() error {
	m.locker.RLock()
	defer m.locker.RUnlock()
	srcInfos, _ := readUserAppFile(m.filename)
	content := m.appInfos.String()
	if content == srcInfos.String() {
		logger.Debug("userAppManager.Write no need to write file")
		return nil
	}

	err := os.MkdirAll(path.Dir(m.filename), 0755)
	if err != nil {
		return err
	}
	logger.Debug("userAppManager.Write write file")
	return ioutil.WriteFile(m.filename, []byte(content), 0644)
}

func (infos userAppInfos) String() string {
	data, _ := json.Marshal(infos)
	return string(data)
}

func (info *userAppInfo) AddMimes(mimes []string) bool {
	var added bool
	for _, mime := range mimes {
		if info.HasMime(mime) {
			continue
		}
		added = true
		info.SupportedMime = append(info.SupportedMime, mime)
	}
	return added
}

func (info *userAppInfo) HasMime(mime string) bool {
	return isStrInList(mime, info.SupportedMime)
}

func newUserAppManager(filename string) (*userAppManager, error) {
	infos, err := readUserAppFile(filename)
	if err != nil {
		return nil, err
	}
	return &userAppManager{
		appInfos: infos,
		filename: filename,
	}, nil
}

func readUserAppFile(filename string) (userAppInfos, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var infos userAppInfos
	err = json.Unmarshal(content, &infos)
	if err != nil {
		return nil, err
	}
	return infos, nil
}
