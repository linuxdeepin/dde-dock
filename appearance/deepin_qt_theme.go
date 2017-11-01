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

package appearance

import (
	"fmt"
	"path"
	"pkg.deepin.io/lib/keyfile"
	"pkg.deepin.io/lib/utils"
	"pkg.deepin.io/lib/xdg/basedir"
	"sync"
)

const (
	dQtSectionTheme = "Theme"
	dQtKeyIcon      = "IconThemeName"
	dQtKeyFont      = "Font"
	dQtKeyMonoFont  = "MonoFont"
	dQtKeyFontSize  = "FontSize"
)

var dQtFile = path.Join(basedir.GetUserConfigDir(), "deepin", "qt-theme.ini")

func setDQtTheme(file, section string, keys, values []string) error {
	_dQtLocker.Lock()
	defer _dQtLocker.Unlock()

	keyLen := len(keys)
	if keyLen != len(values) {
		return fmt.Errorf("keys - values not match: %d - %d", keyLen, len(values))
	}

	_, err := getDQtHandler(file)
	if err != nil {
		return err
	}

	for i := 0; i < keyLen; i++ {
		v, _ := _dQtHandler.GetString(section, keys[i])
		if v == values[i] {
			continue
		}
		_dQtHandler.SetString(section, keys[i], values[i])
		_needSave = true
	}
	return nil
}

func saveDQtTheme(file string) error {
	_dQtLocker.Lock()
	defer _dQtLocker.Unlock()

	if !_needSave {
		_dQtHandler = nil
		return nil
	}
	_needSave = false

	_, err := getDQtHandler(file)
	if err != nil {
		return err
	}

	err = _dQtHandler.SaveToFile(file)
	_dQtHandler = nil
	return err
}

var (
	_dQtHandler *keyfile.KeyFile
	_dQtLocker  sync.Mutex
	_needSave   bool = false
)

func getDQtHandler(file string) (*keyfile.KeyFile, error) {
	if _dQtHandler != nil {
		return _dQtHandler, nil
	}

	if !utils.IsFileExist(file) {
		err := utils.CreateFile(file)
		if err != nil {
			logger.Debug("Failed to create qt theme file:", file, err)
			return nil, err
		}
	}

	kf := keyfile.NewKeyFile()
	err := kf.LoadFromFile(file)
	if err != nil {
		logger.Debug("Failed to load qt theme file:", file, err)
		return nil, err
	}
	_dQtHandler = kf
	return _dQtHandler, nil
}
