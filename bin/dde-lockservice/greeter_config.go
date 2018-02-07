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

package main

import (
	"os"
	"path"
	"pkg.deepin.io/lib/keyfile"
	"pkg.deepin.io/lib/utils"
)

const (
	greeterUserConfig   = "/var/lib/lightdm/lightdm-deepin-greeter/state_user"
	greeterGroupGeneral = "General"
	greeterKeyLastUser  = "last-user"
)

func setGreeterUser(file, username string) error {
	if !utils.IsFileExist(file) {
		err := os.MkdirAll(path.Dir(file), 0755)
		if err != nil {
			return err
		}
		err = utils.CreateFile(file)
		if err != nil {
			return err
		}
	}
	kf, err := loadKeyFile(file)
	if err != nil {
		kf = nil
		return err
	}

	v, err := kf.GetString(greeterGroupGeneral, greeterKeyLastUser)
	if v == username {
		return nil
	}

	kf.SetString(greeterGroupGeneral, greeterKeyLastUser, username)
	return kf.SaveToFile(file)
}

func getGreeterUser(file string) (string, error) {
	kf, err := loadKeyFile(file)
	if err != nil {
		kf = nil
		return "", err
	}

	return kf.GetString(greeterGroupGeneral, greeterKeyLastUser)
}

var _kf *keyfile.KeyFile

func loadKeyFile(file string) (*keyfile.KeyFile, error) {
	if _kf != nil {
		return _kf, nil
	}

	var kf = keyfile.NewKeyFile()
	return kf, kf.LoadFromFile(file)
}
