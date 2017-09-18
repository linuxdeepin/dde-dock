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

package audio

import (
	"encoding/json"
	dutils "pkg.deepin.io/lib/utils"
	"sync"
)

var (
	fileLocker    sync.Mutex
	configCache   *configInfo
	configHandler *dutils.Config
)

func init() {
	configHandler = new(dutils.Config)
	configHandler.SetConfigName("dde-daemon/audio")
}

type configInfo struct {
	Profiles   map[string]string // Profiles[cardName] = activeProfile
	Sink       string
	Source     string
	SinkPort   string
	SourcePort string

	SinkVolume   float64
	SourceVolume float64
}

func (info *configInfo) string() string {
	data, _ := json.Marshal(info)
	return string(data)
}

func (a *configInfo) equal(b *configInfo) bool {
	return (a.Sink == b.Sink &&
		a.Source == b.Source &&
		a.SinkPort == b.SinkPort &&
		a.SourcePort == b.SourcePort &&
		a.SinkVolume == b.SinkVolume &&
		a.SourceVolume == b.SourceVolume &&
		mapStrStrEqual(a.Profiles, b.Profiles))
}

func readConfigInfo() (*configInfo, error) {
	fileLocker.Lock()
	defer fileLocker.Unlock()

	if configCache != nil {
		return configCache, nil
	}

	var info configInfo
	err := configHandler.Load(&info)
	if err != nil {
		return nil, err
	}

	configCache = &info
	return configCache, nil
}

func saveConfigInfo(info *configInfo) error {
	fileLocker.Lock()
	defer fileLocker.Unlock()

	if configCache.equal(info) {
		logger.Debug("[saveConfigInfo] config info not changed")
		return nil
	} else {
		logger.Debug("[saveConfigInfo] will save:", info.string())
	}

	err := configHandler.Save(info)
	if err != nil {
		return err
	}

	configCache = info
	return nil
}

func mapStrStrEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}
	return true
}
