/**
 * Copyright (C) 2017 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package systeminfo

import (
	"encoding/gob"
	"os"
	"path"
	"sync"
)

var (
	cacheLocker sync.Mutex
	cacheFile   = path.Join(os.Getenv("HOME"),
		".cache/deepin/dde-daemon/systeminfo.cache")
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
	err := os.MkdirAll(path.Dir(file), 0755)
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
