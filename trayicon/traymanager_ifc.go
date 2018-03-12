/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package trayicon

import (
	"errors"

	x "github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName = "com.deepin.dde.TrayManager"
	dbusInterface   = dbusServiceName
	dbusPath        = "/com/deepin/dde/TrayManager"
)

func (*TrayManager) GetInterfaceName() string {
	return dbusInterface
}

// Manage方法获取系统托盘图标的管理权。
func (m *TrayManager) Manage() (bool, *dbus.Error) {
	logger.Debug("call Manage by dbus")

	err := m.sendClientMsgMANAGER()
	if err != nil {
		logger.Warning(err)
		return false, dbusutil.ToError(err)
	}
	return true, nil
}

// GetName返回传入的系统图标的窗口id的窗口名。
func (m *TrayManager) GetName(win uint32) (string, *dbus.Error) {
	m.mutex.Lock()
	icon, ok := m.icons[x.Window(win)]
	m.mutex.Unlock()
	if !ok {
		return "", dbusutil.ToError(errors.New("icon not found"))
	}
	return icon.getName(), nil
}

// EnableNotification设置对应id的窗口是否可以通知。
func (m *TrayManager) EnableNotification(win uint32, enable bool) *dbus.Error {
	m.mutex.Lock()
	icon, ok := m.icons[x.Window(win)]
	m.mutex.Unlock()
	if !ok {
		return dbusutil.ToError(errors.New("icon not found"))
	}

	icon.mu.Lock()
	icon.notify = enable
	icon.mu.Unlock()
	return nil
}
