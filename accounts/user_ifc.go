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

package accounts

import (
	"os"
	"path"
	"pkg.linuxdeepin.com/lib/dbus"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"strings"
)

func (obj *User) SetUserName(dbusMsg dbus.DMessage, username string) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In SetUserName:%v",
				err)
		}
	}()

	//if ok := dutils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	if obj.UserName != username {
		obj.setPropUserName(username)
	}

	return true
}

func (obj *User) SetHomeDir(dbusMsg dbus.DMessage, dir string) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In SetHomeDir:%v",
				err)
		}
	}()

	//if ok := dutils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	if obj.HomeDir != dir {
		obj.setPropHomeDir(dir)
	}

	return true
}

func (obj *User) SetShell(dbusMsg dbus.DMessage, shell string) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In SetShell:%v",
				err)
		}
	}()

	//if ok := dutils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	if obj.Shell != shell {
		obj.setPropShell(shell)
	}

	return true
}

func (obj *User) SetPassword(dbusMsg dbus.DMessage, words string) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In SetPassword:%v",
				err)
		}
	}()

	//if ok := dutils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	passwd := encodePasswd(words)
	changePasswd(obj.UserName, passwd)

	obj.setPropLocked(false)

	return true
}

func (obj *User) SetAutomaticLogin(dbusMsg dbus.DMessage, auto bool) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In SetAutomaticLogin:%v",
				err)
		}
	}()

	//if ok := dutils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	if obj.AutomaticLogin != auto {
		obj.setPropAutomaticLogin(auto)
	}

	return true
}

func (obj *User) SetAccountType(dbusMsg dbus.DMessage, t int32) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In SetAccountType:%v",
				err)
		}
	}()

	//if ok := dutils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	if obj.AccountType != t {
		obj.setPropAccountType(t)
	}

	return true
}

func (obj *User) SetLocked(dbusMsg dbus.DMessage, locked bool) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In SetLocked:%v",
				err)
		}
	}()

	//if ok := dutils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	if obj.Locked != locked {
		obj.setPropLocked(locked)
	}

	return true
}

func (obj *User) SetIconFile(dbusMsg dbus.DMessage, icon string) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In SetIconFile:%v",
				err)
		}
	}()

	//if ok := dutils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	if ok := dutils.IsFileExist(icon); !ok || obj.IconFile == icon {
		return false
	}

	if !strIsInList(icon, obj.IconList) {
		if ok := dutils.IsFileExist(ICON_LOCAL_DIR); !ok {
			if err := os.MkdirAll(ICON_LOCAL_DIR, 0755); err != nil {
				return false
			}
		}
		name := path.Base(icon)
		dest := path.Join(ICON_LOCAL_DIR, obj.UserName+"-"+name)
		if err := dutils.CopyFile(icon, dest); err != nil {
			return false
		}
		icon = dest
	}
	obj.setPropIconFile(icon)

	return true
}

func (obj *User) SetBackgroundFile(dbusMsg dbus.DMessage, bg string) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In SetBackgroundFile:%v",
				err)
		}
	}()

	//if ok := dutils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	if obj.BackgroundFile != bg {
		obj.setPropBackgroundFile(bg)
	}

	return true
}

func (obj *User) DeleteHistoryIcon(dbusMsg dbus.DMessage, icon string) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In DeleteHistoryIcon:%v",
				err)
		}
	}()

	//if ok := dutils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	obj.deleteHistoryIcon(icon)

	return true
}

func (u *User) DeleteIconFile(dbusMsg dbus.DMessage, icon string) (success bool) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In DeleteHistoryIcon:%v",
				err)
		}
	}()

	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	if err := rmAllFile(icon); err != nil {
		return false
	}

	u.deleteHistoryIcon(icon)

	return true
}

func (u *User) IsIconDeletable(icon string) bool {
	if icon == u.IconFile {
		return false
	}

	if strings.Contains(icon, path.Join(ICON_LOCAL_DIR, u.UserName)) {
		return true
	}

	return false
}
