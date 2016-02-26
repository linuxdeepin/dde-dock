/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dstore

import (
	"bufio"
	"encoding/json"
	"os"
)

type DQueryTimeInstalledTransaction struct {
	data map[string]int64
}

func NewDQueryTimeInstalledTransaction(file string) (*DQueryTimeInstalledTransaction, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := &DQueryTimeInstalledTransaction{data: map[string]int64{}}
	decoder := json.NewDecoder(bufio.NewReader(f))
	err = decoder.Decode(&t.data)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *DQueryTimeInstalledTransaction) Query(pkgName string) int64 {
	timeInstalled := t.data[pkgName]
	return timeInstalled
}
