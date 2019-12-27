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

package audio

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"pkg.deepin.io/lib/xdg/basedir"
)

var (
	fileLocker  sync.Mutex
	configCache *config
	configFile  = filepath.Join(basedir.GetUserConfigDir(), "deepin/dde-daemon/audio.json")
)

type config struct {
	Profiles   map[string]string // Profiles[cardName] = activeProfile
	Sink       string
	Source     string
	SinkPort   string
	SourcePort string

	SinkVolume   float64
	SourceVolume float64
}

func (c *config) string() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func (c *config) equal(b *config) bool {
	if c == nil && b == nil {
		return true
	}
	if c == nil || b == nil {
		return false
	}
	return c.Sink == b.Sink &&
		c.Source == b.Source &&
		c.SinkPort == b.SinkPort &&
		c.SourcePort == b.SourcePort &&
		c.SinkVolume == b.SinkVolume &&
		c.SourceVolume == b.SourceVolume &&
		mapStrStrEqual(c.Profiles, b.Profiles)
}

func readConfig() (*config, error) {
	fileLocker.Lock()
	defer fileLocker.Unlock()

	if configCache != nil {
		return configCache, nil
	}

	var info config
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &info)
	if err != nil {
		return nil, err
	}

	configCache = &info
	return configCache, nil
}

func saveConfig(info *config) error {
	fileLocker.Lock()
	defer fileLocker.Unlock()

	if configCache.equal(info) {
		logger.Debug("[saveConfigInfo] config info not changed")
		return nil
	} else {
		logger.Debug("[saveConfigInfo] will save:", info.string())
	}

	content, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(configFile), 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFile, content, 0644)
	if err != nil {
		return err
	}

	configCache = info
	return nil
}

func removeConfig() error {
	fileLocker.Lock()
	defer fileLocker.Unlock()
	return os.Remove(configFile)
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
