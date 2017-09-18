/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package appinfo

import (
	"gir/glib-2.0"
	"io/ioutil"
	"os"
)

const (
	_RateRecordFile = "launcher/rate.ini"
	_RateRecordKey  = "rate"
)

// GetFrequencyRecordFile returns the file which records items' use frequency.
func GetFrequencyRecordFile() (*glib.KeyFile, error) {
	return ConfigFile(_RateRecordFile)
}

func GetFrequency(id string, f *glib.KeyFile) uint64 {
	rate, _ := f.GetUint64(id, _RateRecordKey)
	return rate
}

func SetFrequency(id string, freq uint64, f *glib.KeyFile) {
	f.SetUint64(id, _RateRecordKey, freq)
	saveKeyFile(f, ConfigFilePath(_RateRecordFile))
}

// saveKeyFile saves key file.
func saveKeyFile(file *glib.KeyFile, path string) error {
	_, content, err := file.ToData()
	if err != nil {
		return err
	}

	stat, err := os.Lstat(path)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, []byte(content), stat.Mode())
	if err != nil {
		return err
	}
	return nil
}
