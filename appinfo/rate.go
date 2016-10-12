/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
