/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package appinfo

// #cgo pkg-config: glib-2.0
// #include <glib.h>
import "C"
import (
	"os"
	"path"

	"gir/glib-2.0"
	"pkg.deepin.io/lib/utils"
)

const (
	_DirDefaultPerm os.FileMode = 0755
)

// ConfigFilePath returns path in user's config dir.
func ConfigFilePath(name string) string {
	return path.Join(glib.GetUserConfigDir(), name)
}

// ConfigFile open the given keyfile, this file will be created if not existed.
func ConfigFile(name string) (*glib.KeyFile, error) {
	file := glib.NewKeyFile()
	conf := ConfigFilePath(name)
	if !utils.IsFileExist(conf) {
		os.MkdirAll(path.Dir(conf), _DirDefaultPerm)
		f, err := os.Create(conf)
		if err != nil {
			return nil, err
		}
		defer f.Close()
	}

	if ok, err := file.LoadFromFile(conf, glib.KeyFileFlagsNone); !ok {
		file.Free()
		return nil, err
	}
	return file, nil
}
