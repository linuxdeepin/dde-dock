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
	"io"
	"os"
	"sync"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName = "com.deepin.daemon.Greeter"
	dbusPath        = "/com/deepin/daemon/Greeter"
	dbusInterface   = dbusServiceName
)

type Manager struct {
	service *dbusutil.Service
	mu      sync.Mutex

	methods *struct {
		UpdateGreeterQtTheme func() `in:"fd"`
	}
}

func (m *Manager) UpdateGreeterQtTheme(fd dbus.UnixFD) *dbus.Error {
	m.service.DelayAutoQuit()
	err := updateGreeterQtTheme(fd)
	if err != nil {
		logger.Warning(err)
	}
	return dbusutil.ToError(err)
}

func updateGreeterQtTheme(fd dbus.UnixFD) error {
	f := os.NewFile(uintptr(fd), "")
	defer f.Close()
	err := os.MkdirAll("/etc/lightdm/deepin", 0755)
	if err != nil {
		return err
	}
	const (
		themeFile     = "/etc/lightdm/deepin/qt-theme.ini"
		themeFileTemp = themeFile + ".tmp"
	)
	dest, err := os.Create(themeFileTemp)
	if err != nil {
		return err
	}
	// limit file size: 100KB
	src := io.LimitReader(f, 1024*100)
	_, err = io.Copy(dest, src)
	if err != nil {
		closeErr := dest.Close()
		if closeErr != nil {
			logger.Warning(closeErr)
		}
		return err
	}

	err = dest.Close()
	if err != nil {
		return err
	}

	err = os.Rename(themeFileTemp, themeFile)
	return err
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}
