/*
 * Copyright (C) 2015 ~ 2018 Deepin Technology Co., Ltd.
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

package dstore

import (
	"bufio"
	"encoding/json"
	"os"
)

type DQueryPkgNameTransaction struct {
	data map[string]string
}

// NewDQueryPkgNameTransaction returns package name of given desktop file.
func NewDQueryPkgNameTransaction(path string) (*DQueryPkgNameTransaction, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	t := &DQueryPkgNameTransaction{data: map[string]string{}}
	decoder := json.NewDecoder(bufio.NewReader(f))
	err = decoder.Decode(&t.data)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *DQueryPkgNameTransaction) Query(desktopID string) string {
	if t.data != nil {
		pkg := t.data[desktopID]
		return pkg
	}
	return ""
}
