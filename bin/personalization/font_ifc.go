/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"errors"
	"os/exec"
	"path"
	dutils "pkg.deepin.io/lib/utils"
	"strconv"
)

const (
	_QT_CONFIG_FILE    = ".config/Trolltech.conf"
	_DEFAULT_FONT_SIZE = 11

	_QT_KEY_GROUP = "Qt"
	_QT_KEY_FONT  = "font"
	_QT_FONT_ARGS = ",-1,5,50,0,0,0,0,0"
)

func (fs *FontSettings) setQtFont(name, sizeStr string) error {
	if len(name) < 1 || len(sizeStr) < 1 {
		return errors.New("Set Qt Font args error")
	}

	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		return errors.New("Get home dir failed")
	}

	filename := path.Join(homeDir, _QT_CONFIG_FILE)
	value := "\"" + name + "," + sizeStr + _QT_FONT_ARGS + "\""
	if !dutils.WriteKeyToKeyFile(filename, _QT_KEY_GROUP,
		_QT_KEY_FONT, value) {
		return errors.New("Set Qt font failed")
	}

	return nil
}

func (fs *FontSettings) SetDocumentFont(name string, size int32) error {
	if len(name) < 1 {
		return errors.New("Set document font args error")
	}

	if size < 9 && size > 26 {
		size = _DEFAULT_FONT_SIZE
	}

	sizeStr := strconv.FormatInt(int64(size), 10)
	if err := fs.setQtFont(name, sizeStr); err != nil {
		return err
	}

	if fs.xs == nil {
		fs.initSettings()
	}
	fs.xs.SetString("Gtk/FontName", name+" "+sizeStr)

	return nil
}

func (fs *FontSettings) SetTitleFont(name string, size int32) error {
	if len(name) < 1 {
		return errors.New("Set document font args error")
	}

	if size < 9 && size > 26 {
		size = _DEFAULT_FONT_SIZE
	}

	sizeStr := strconv.FormatInt(int64(size), 10)
	if err := fs.setQtFont(name, sizeStr); err != nil {
		return err
	}

	if fs.wmSettings == nil {
		fs.initSettings()
	}
	fs.wmSettings.SetString("titlebar-font", name+" "+sizeStr)

	return nil
}

func (fs *FontSettings) SetMonoFont(name string, size int32) error {
	if len(name) < 1 {
		return errors.New("Set document font args error")
	}

	if size < 9 && size > 26 {
		size = _DEFAULT_FONT_SIZE
	}

	sizeStr := strconv.FormatInt(int64(size), 10)
	if err := fs.setQtFont(name, sizeStr); err != nil {
		return err
	}

	if out, err := exec.Command("/usr/bin/gconftool",
		"-t", "string",
		"-s", "/desktop/gnome/interface/monospace_font_name",
		name+" "+sizeStr).CombinedOutput(); err != nil {
		return errors.New(string(out))
	}

	return nil
}
