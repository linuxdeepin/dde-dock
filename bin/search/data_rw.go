/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
		logger.Warning("Open '%s' failed:", err)
		return
	}
	defer fp.Close()

	fp.WriteString(w.String())
	fp.Sync()
}

func readDatasFromFile(datas interface{}, filename string) bool {
	if !dutils.IsFileExist(filename) || datas == nil {
		logger.Warning("readDatasFromFile args error")
		return false
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warning("ReadFile '%s' failed:", err)
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
