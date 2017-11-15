/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

package fonts

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var (
	familyHashCacheFile = path.Join(home, ".cache", "deepin", "dde-daemon", "fonts", "family_hash")
)

func (table FamilyHashTable) saveToFile() error {
	return doSaveObject(familyHashCacheFile, &table)
}

func loadCacheFromFile(file string, obj interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var r = bytes.NewBuffer(data)
	decoder := gob.NewDecoder(r)
	err = decoder.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}

func doSaveObject(file string, obj interface{}) error {
	var w bytes.Buffer
	encoder := gob.NewEncoder(&w)
	err := encoder.Encode(obj)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(file), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, w.Bytes(), 0644)
}

func writeFontConfig(content, file string) error {
	err := os.MkdirAll(path.Dir(file), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, []byte(content), 0644)
}

// If set pixelsize, wps-office-wps will not show some text.
//
//func configContent(standard, mono string, pixel float64) string {
func configContent(standard, mono string) string {
	return fmt.Sprintf(`<?xml version="1.0"?>
<!DOCTYPE fontconfig SYSTEM "fonts.dtd">
<fontconfig>
    <match target="pattern">
        <test qual="any" name="family">
            <string>serif</string>
        </test>
        <edit name="family" mode="assign" binding="strong">
            <string>%s</string>
            <string>%s</string>
        </edit>
    </match>

    <match target="pattern">
        <test qual="any" name="family">
            <string>sans-serif</string>
        </test>
        <edit name="family" mode="assign" binding="strong">
            <string>%s</string>
            <string>%s</string>
        </edit>
    </match>

    <match target="pattern">
        <test qual="any" name="family">
            <string>monospace</string>
        </test>
        <edit name="family" mode="assign" binding="strong">
            <string>%s</string>
            <string>%s</string>
        </edit>
    </match>

</fontconfig>`, standard, fallbackStandard,
		standard, fallbackStandard,
		mono, fallbackMonospace)
}
