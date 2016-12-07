/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package accounts

import (
	"errors"
	"fmt"
	"os"
	"path"
	"pkg.deepin.io/dde/api/lang_info"
	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/graphic"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
)

func (u *User) SetUserName(dbusMsg dbus.DMessage, name string) error {
	logger.Debug("[SetUserName] new name:", name)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, false); err != nil {
		logger.Debug("[SetUserName] access denied:", err)
		return err
	}

	if err := users.ModifyName(name, u.UserName); err != nil {
		logger.Warning("DoAction: modify username failed:", err)
		return err
	}

	u.setPropString(&u.UserName, "UserName", name)
	return nil
}

func (u *User) SetHomeDir(dbusMsg dbus.DMessage, home string) error {
	logger.Debug("[SetHomeDir] new home:", home)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, false); err != nil {
		logger.Debug("[SetHomeDir] access denied:", err)
		return err
	}

	if err := users.ModifyHome(home, u.UserName); err != nil {
		logger.Warning("DoAction: modify home failed:", err)
		return err
	}

	u.setPropString(&u.HomeDir, "HomeDir", home)
	return nil
}

func (u *User) SetShell(dbusMsg dbus.DMessage, shell string) error {
	logger.Debug("[SetShell] new shell:", shell)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
		logger.Debug("[SetShell] access denied:", err)
		return err
	}

	if err := users.ModifyShell(shell, u.UserName); err != nil {
		logger.Warning("DoAction: modify shell failed:", err)
		return err
	}

	u.setPropString(&u.Shell, "Shell", shell)
	return nil
}

func (u *User) SetPassword(dbusMsg dbus.DMessage, words string) error {
	logger.Debug("[SetPassword] start ...")
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, false); err != nil {
		logger.Debug("[SetPassword] access denied:", err)
		return err
	}

	if err := users.ModifyPasswd(words, u.UserName); err != nil {
		logger.Warning("DoAction: modify passwd failed:", err)
		return err
	}

	if err := users.LockedUser(false, u.UserName); err != nil {
		logger.Warning("DoAction: unlock user failed:", err)
		return err
	}
	u.setPropBool(&u.Locked, "Locked", false)
	return nil
}

func (u *User) SetAccountType(dbusMsg dbus.DMessage, ty int32) error {
	logger.Debug("[SetAccountType] type:", ty)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, false); err != nil {
		logger.Debug("[SetAccountType] access denied:", err)
		return err
	}

	if err := users.SetUserType(ty, u.UserName); err != nil {
		logger.Warning("DoAction: set user type failed:", err)
		return err
	}

	u.setPropInt32(&u.AccountType, "AccountType", ty)
	return nil
}

func (u *User) SetLocked(dbusMsg dbus.DMessage, locked bool) error {
	logger.Debug("[SetLocked] locaked:", locked)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, false); err != nil {
		logger.Debug("[SetLocked] access denied:", err)
		return err
	}

	if err := users.LockedUser(locked, u.UserName); err != nil {
		logger.Warning("DoAction: locked user failed:", err)
		return err
	}

	if locked && u.AutomaticLogin {
		users.SetAutoLoginUser("")
	}
	u.setPropBool(&u.Locked, "Locked", locked)
	return nil
}

func (u *User) SetAutomaticLogin(dbusMsg dbus.DMessage, auto bool) error {
	logger.Debug("[SetAutomaticLogin] auto", auto)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, false); err != nil {
		logger.Debug("[SetAutomaticLogin] access denied:", err)
		return err
	}

	if u.Locked {
		return fmt.Errorf("user %s has been locked", u.UserName)
	}

	var name = u.UserName
	if !auto {
		name = ""
	}

	if err := users.SetAutoLoginUser(name); err != nil {
		logger.Warning("DoAction: set auto login failed:", err)
		return err
	}

	u.setPropBool(&u.AutomaticLogin, "AutomaticLogin", auto)
	return nil
}

func (u *User) SetLocale(dbusMsg dbus.DMessage, locale string) error {
	logger.Debug("[SetLocale] locale:", locale)
	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
		logger.Debug("[SetLocale] access denied:", err)
		return err
	}

	if u.Locale == locale {
		return nil
	}

	if !lang_info.IsSupportedLocale(locale) {
		err := fmt.Errorf("Invalid locale %q", locale)
		logger.Debug("[SetLocale]", err)
		return err
	}

	if err := u.writeUserConfig(); err != nil {
		logger.Warning("[SetLocale]", err)
		return err
	}

	u.setPropString(&u.Locale, "Locale", locale)
	return nil
}

func (u *User) SetLayout(dbusMsg dbus.DMessage, layout string) error {
	logger.Debug("[SetLayout] new layout:", layout)
	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
		logger.Debug("[SetLayout] access denied:", err)
		return err
	}

	if u.Layout == layout {
		return nil
	}

	// TODO: check layout validity
	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
		return err
	}
	u.setPropString(&u.Layout, "Layout", layout)
	return nil
}

func (u *User) SetIconFile(dbusMsg dbus.DMessage, icon string) error {
	logger.Debug("[SetIconFile] new icon:", icon)
	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
		logger.Debug("[SetIconFile] access denied:", err)
		return err
	}

	srcIcon := dutils.DecodeURI(icon)
	icon = dutils.EncodeURI(icon, dutils.SCHEME_FILE)
	if u.IconFile == icon {
		return nil
	}

	if !graphic.IsSupportedImage(srcIcon) {
		err := fmt.Errorf("This icon '%s' not a image", icon)
		logger.Debug(err)
		return err
	}

	target, added, err := u.addIconFile(icon)
	if err != nil {
		logger.Warning("Set icon failed:", err)
		return err
	}

	u.addHistoryIcon(u.IconFile)
	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
		return err
	}
	u.setPropString(&u.IconFile, "IconFile", target)
	if added {
		u.setPropStrv(&u.IconList, "IconList", u.getAllIcons())
	}
	return nil
}

func (u *User) DeleteIconFile(dbusMsg dbus.DMessage, icon string) error {
	logger.Debug("[DeleteIconFile] icon:", icon)
	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
		logger.Debug("[DeleteIconFile] access denied:", err)
		return err
	}

	icon = dutils.EncodeURI(icon, dutils.SCHEME_FILE)
	if !u.IsIconDeletable(icon) {
		err := errors.New("This icon is not allowed to be deleted!")
		logger.Warning(err)
		return err
	}

	iconPath := dutils.DecodeURI(icon)
	if err := os.Remove(iconPath); err != nil {
		return err
	}

	u.DeleteHistoryIcon(dbusMsg, icon)
	u.setPropStrv(&u.IconList, "IconList", u.getAllIcons())
	return nil
}

func (u *User) SetBackgroundFile(dbusMsg dbus.DMessage, bg string) error {
	logger.Debug("[SetBackgroundFile] new background:", bg)
	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
		logger.Debug("[SetBackgroundFile] access denied:", err)
		return err
	}
	bg = dutils.EncodeURI(bg, dutils.SCHEME_FILE)
	if bg == u.BackgroundFile {
		genGaussianBlur(bg)
		return nil
	}

	if err := checkBackgroundValid(bg); err != nil {
		logger.Debug(err)
		return err
	}

	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
		return err
	}

	u.setPropString(&u.BackgroundFile, "BackgroundFile", bg)
	genGaussianBlur(bg)
	return nil
}

func (u *User) SetGreeterBackground(dbusMsg dbus.DMessage, bg string) error {
	logger.Debug("[SetGreeterBackground] new background:", bg)
	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
		logger.Debug("[SetGreeterBackground] access denied:", err)
		return err
	}
	bg = dutils.EncodeURI(bg, dutils.SCHEME_FILE)
	if u.GreeterBackground == bg {
		genGaussianBlur(bg)
		return nil
	}

	if err := checkBackgroundValid(bg); err != nil {
		logger.Debug(err)
		return err
	}

	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
		return err
	}

	u.setPropString(&u.GreeterBackground, "GreeterBackground", bg)
	genGaussianBlur(bg)
	return nil
}

func (u *User) SetHistoryLayout(dbusMsg dbus.DMessage, list []string) error {
	logger.Debug("[SetHistoryLayout] new history layout:", list)
	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
		logger.Debug("[SetHistoryLayout] access denied:", err)
		return err
	}

	if isStrvEqual(u.HistoryLayout, list) {
		return nil
	}

	// TODO: check layout list whether validity
	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
	}
	u.setPropStrv(&u.HistoryLayout, "HistoryLayout", list)
	return nil
}

func (u *User) DeleteHistoryIcon(dbusMsg dbus.DMessage, icon string) error {
	logger.Debug("[DeleteHistoryIcon] icon:", icon)
	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
		logger.Debug("[DeleteHistoryIcon] access denied:", err)
		return err
	}

	icon = dutils.EncodeURI(icon, dutils.SCHEME_FILE)
	u.deleteHistoryIcon(icon)
	return nil
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
	dir := path.Dir(dutils.DecodeURI(u.IconFile))

	filename := path.Join(dir, "bigger", baseName)
	if !dutils.IsFileExist(filename) {
		return ""
	}

	return dutils.EncodeURI(filename, dutils.SCHEME_FILE)
}

func checkBackgroundValid(bg string) error {
	bg = dutils.DecodeURI(bg)
	if !graphic.IsSupportedImage(bg) {
		return errors.New("unsupported image format")
	}
	// ok
	return nil
}
