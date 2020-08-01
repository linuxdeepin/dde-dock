/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"

	dutils "pkg.deepin.io/lib/utils"
)

func writeDatasToFile(datas interface{}, filename string) {
	if datas == nil {
		logger.Warning("writeDatasToFile args error")
		return
	}

	var w bytes.Buffer
	enc := gob.NewEncoder(&w)
	if err := enc.Encode(datas); err != nil {
		logger.Warning("Gob Encode Datas Failed:", err)
		return
	}

	fp, err := os.Create(filename)
	if err != nil {
		logger.Warningf("failed to open %q: %v", filename, err)
		return
	}
	defer fp.Close()

	_, _ = fp.WriteString(w.String())
	_ = fp.Sync()
}

func readDatasFromFile(datas interface{}, filename string) bool {
	if !dutils.IsFileExist(filename) || datas == nil {
		logger.Warning("readDatasFromFile args error")
		return false
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warningf("failed to read file %q: %v", filename, err)
		return false
	}

	r := bytes.NewBuffer(contents)
	dec := gob.NewDecoder(r)
	if err = dec.Decode(datas); err != nil {
		logger.Warning("Decode Datas Failed:", err)
		return false
	}

	return true
}
