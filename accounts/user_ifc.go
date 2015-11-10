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
	"fmt"
	"os"
	"path"
	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/graphic"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
)

func (u *User) SetUserName(dbusMsg dbus.DMessage, name string) (bool, error) {
	u.syncLocker.Lock()
	logger.Debug("[SetUserName] new name:", name)
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, false, "SetUserName")
	if err != nil {
		u.syncLocker.Unlock()
		logger.Debug("[SetUserName] access denied:", err)
		return false, err
	}

	go func() {
		defer u.syncLocker.Unlock()

		err := users.ModifyName(name, u.UserName)
		if err != nil {
			logger.Warning("DoAction: modify username failed:", err)
			triggerSigErr(pid, "SetUserName", err.Error())
			return
		}

		u.setPropString(&u.UserName, "UserName", name)
	}()

	return true, nil
}

func (u *User) SetHomeDir(dbusMsg dbus.DMessage, home string) (bool, error) {
	u.syncLocker.Lock()
	logger.Debug("[SetHomeDir] new home:", home)
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, false, "SetHomeDir")
	if err != nil {
		u.syncLocker.Unlock()
		logger.Debug("[SetHomeDir] access denied:", err)
		return false, err
	}

	go func() {
		defer u.syncLocker.Unlock()

		err := users.ModifyHome(home, u.UserName)
		if err != nil {
			logger.Warning("DoAction: modify home failed:", err)
			triggerSigErr(pid, "SetHomeDir", err.Error())
			return
		}

		u.setPropString(&u.HomeDir, "HomeDir", home)
	}()

	return true, nil
}

func (u *User) SetShell(dbusMsg dbus.DMessage, shell string) (bool, error) {
	u.syncLocker.Lock()
	logger.Debug("[SetShell] new shell:", shell)
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, true, "SetShell")
	if err != nil {
		u.syncLocker.Unlock()
		logger.Debug("[SetShell] access denied:", err)
		return false, err
	}

	go func() {
		defer u.syncLocker.Unlock()

		err := users.ModifyShell(shell, u.UserName)
		if err != nil {
			logger.Warning("DoAction: modify shell failed:", err)
			triggerSigErr(pid, "SetShell", err.Error())
			return
		}

		u.setPropString(&u.Shell, "Shell", shell)
	}()

	return true, nil
}

func (u *User) SetPassword(dbusMsg dbus.DMessage, words string) (bool, error) {
	u.syncLocker.Lock()
	logger.Debug("[SetPassword] start ...")
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, false, "SetPassword")
	if err != nil {
		u.syncLocker.Unlock()
		logger.Debug("[SetPassword] access denied:", err)
		return false, err
	}

	go func() {
		defer u.syncLocker.Unlock()

		err := users.ModifyPasswd(words, u.UserName)
		if err != nil {
			logger.Warning("DoAction: modify passwd failed:", err)
			triggerSigErr(pid, "SetPassword", err.Error())
			return
		}

		err = users.LockedUser(false, u.UserName)
		if err != nil {
			logger.Warning("DoAction: unlock user failed:", err)
		}
		u.setPropBool(&u.Locked, "Locked", false)
	}()

	return true, nil
}

func (u *User) SetAccountType(dbusMsg dbus.DMessage, ty int32) (bool, error) {
	u.syncLocker.Lock()
	logger.Debug("[SetAccountType] type:", ty)
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, false, "SetAccountType")
	if err != nil {
		u.syncLocker.Unlock()
		logger.Debug("[SetAccountType] access denied:", err)
		return false, err
	}

	go func() {
		defer u.syncLocker.Unlock()

		err := users.SetUserType(ty, u.UserName)
		if err != nil {
			logger.Warning("DoAction: set user type failed:", err)
			triggerSigErr(pid, "SetAccountType", err.Error())
			return
		}

		u.setPropInt32(&u.AccountType, "AccountType", ty)
	}()

	return true, nil
}

func (u *User) SetLocked(dbusMsg dbus.DMessage, locked bool) (bool, error) {
	u.syncLocker.Lock()
	logger.Debug("[SetLocked] locaked:", locked)
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, false, "SetLocked")
	if err != nil {
		u.syncLocker.Unlock()
		logger.Debug("[SetLocked] access denied:", err)
		return false, err
	}

	go func() {
		defer u.syncLocker.Unlock()

		err := users.LockedUser(locked, u.UserName)
		if err != nil {
			logger.Warning("DoAction: locked user failed:", err)
			triggerSigErr(pid, "SetLocked", err.Error())
			return
		}

		if locked && u.AutomaticLogin {
			users.SetAutoLoginUser("")
		}

		u.setPropBool(&u.Locked, "Locked", locked)
	}()

	return true, nil
}

func (u *User) SetAutomaticLogin(dbusMsg dbus.DMessage, auto bool) (bool, error) {
	u.syncLocker.Lock()
	logger.Debug("[SetAutomaticLogin] auto", auto)
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, false, "SetAutomaticLogin")
	if err != nil {
		u.syncLocker.Unlock()
		logger.Debug("[SetAutomaticLogin] access denied:", err)
		return false, err
	}

	if u.Locked {
		u.syncLocker.Unlock()
		return false, fmt.Errorf("%s has been locked", u.UserName)
	}

	var name = u.UserName
	if !auto {
		name = ""
	}

	go func() {
		defer u.syncLocker.Unlock()

		err := users.SetAutoLoginUser(name)
		if err != nil {
			logger.Warning("DoAction: set auto login failed:", err)
			triggerSigErr(pid, "SetAutomaticLogin", err.Error())
			return
		}

		u.setPropBool(&u.AutomaticLogin, "AutomaticLogin", auto)
	}()

	return true, nil
}

func (u *User) SetIconFile(dbusMsg dbus.DMessage, icon string) (bool, error) {
	logger.Debug("[SetIconFile] new icon:", icon)
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, true, "SetIconFile")
	if err != nil {
		logger.Debug("[SetIconFile] access denied:", err)
		return false, err
	}

	if u.IconFile == icon {
		return true, nil
	}

	if !graphic.IsSupportedImage(icon) {
		reason := fmt.Sprintf("This icon '%s' not a image", icon)
		logger.Debug(reason)
		triggerSigErr(pid, "SetIconFile", reason)
		return false, err
	}

	go func() {
		target, added, err := u.addIconFile(icon)
		if err != nil {
			logger.Warning("Set icon failed:", err)
			triggerSigErr(pid, "SetIconFile", err.Error())
			return
		}

		src := u.IconFile
		u.setPropString(&u.IconFile, "IconFile", target)
		u.addHistoryIcon(src)
		err = u.writeUserConfig()
		if err != nil {
			logger.Warning("Write user config failed:", err)
			triggerSigErr(pid, "SetIconFile", err.Error())
			u.setPropString(&u.IconFile, "IconFile", src)
			return
		}
		if added {
			u.setPropStrv(&u.IconList, "IconList", u.getAllIcons())
		}
	}()

	return true, nil
}

func (u *User) DeleteIconFile(dbusMsg dbus.DMessage, icon string) (bool, error) {
	logger.Debug("[DeleteIconFile] icon:", icon)
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, true, "DeleteIconFile")
	if err != nil {
		logger.Debug("[DeleteIconFile] access denied:", err)
		return false, err
	}

	if !u.IsIconDeletable(icon) {
		reason := "This icon is not allowed to be deleted!"
		logger.Warning(reason)
		triggerSigErr(pid, "DeleteHistoryIcon", reason)
		return false, fmt.Errorf(reason)
	}

	go func() {
		err := os.Remove(icon)
		if err != nil {
			triggerSigErr(pid, "DeleteIconFile", err.Error())
			return
		}

		u.DeleteHistoryIcon(dbusMsg, icon)
		u.setPropStrv(&u.IconList, "IconList", u.getAllIcons())
	}()

	return true, nil
}

func (u *User) SetBackgroundFile(dbusMsg dbus.DMessage, bg string) (bool, error) {
	logger.Debug("[SetBackgroundFile] new background:", bg)
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, true, "SetBackgroundFile")
	if err != nil {
		logger.Debug("[SetBackgroundFile] access denied:", err)
		return false, err
	}

	if bg == u.BackgroundFile {
		return true, nil
	}

	if !graphic.IsSupportedImage(bg) {
		reason := fmt.Sprintf("This background '%s' not a image", bg)
		logger.Debug(reason)
		triggerSigErr(pid, "SetBackgroundFile", reason)
		return false, err
	}

	go func() {
		src := u.BackgroundFile
		u.setPropString(&u.BackgroundFile, "BackgroundFile", bg)
		err = u.writeUserConfig()
		if err != nil {
			logger.Warning("Write user config failed:", err)
			triggerSigErr(pid, "SetBackgroundFile", err.Error())
			u.setPropString(&u.BackgroundFile, "BackgroundFile", src)
			return
		}
	}()

	return true, nil
}

func (u *User) DeleteHistoryIcon(dbusMsg dbus.DMessage, icon string) (bool, error) {
	logger.Debug("[DeleteHistoryIcon] icon:", icon)
	pid := dbusMsg.GetSenderPID()
	err := u.accessAuthentication(pid, true, "DeleteHistoryIcon")
	if err != nil {
		logger.Debug("[DeleteHistoryIcon] access denied:", err)
		return false, err
	}

	u.deleteHistoryIcon(icon)
	return true, nil
}

func (u *User) IsIconDeletable(icon string) bool {
	if u.IconFile == icon {
		return false
	}

	if !strings.Contains(icon, path.Join(userCustomIconsDir, u.UserName)) {
		return false
	}

	return true
}

// 获取当前头像的大图标
func (u *User) GetLargeIcon() string {
	baseName := path.Base(u.IconFile)
	dir := path.Dir(u.IconFile)

	filename := path.Join(dir, "bigger", baseName)
	if !dutils.IsFileExist(filename) {
		return ""
	}

	return filename
}
