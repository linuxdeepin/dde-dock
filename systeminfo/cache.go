/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package systeminfo

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"sync"

	"pkg.deepin.io/lib/xdg/basedir"
)

var (
	cacheLocker sync.Mutex
	cacheFile   = filepath.Join(basedir.GetUserCacheDir(),
		"deepin/dde-daemon/systeminfo.cache")
)

func doReadCache(file string) (*SystemInfo, error) {
	cacheLocker.Lock()
	defer cacheLocker.Unlock()
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	var info SystemInfo
	decoder := gob.NewDecoder(fp)
	err = decoder.Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func doSaveCache(info *SystemInfo, file string) error {
	cacheLocker.Lock()
	defer cacheLocker.Unlock()
	err := os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		return err
	}

	fp, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fp.Close()

	encoder := gob.NewEncoder(fp)
	err = encoder.Encode(info)
	if err != nil {
		return err
	}
	return fp.Sync()
}
