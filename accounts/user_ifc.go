/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
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

package accounts

import (
	"errors"
	"fmt"
	"os"
	"path"

	"pkg.deepin.io/dde/api/lang_info"
	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gdkpixbuf"
	"pkg.deepin.io/lib/graphic"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
)

func (u *User) SetFullName(dbusMsg dbus.DMessage, name string) error {
	logger.Debugf("[SetFullName] new name: %q", name)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
		logger.Debug("[SetFullName] access denied:", err)
		return err
	}

	if err := users.ModifyFullName(name, u.UserName); err != nil {
		logger.Warning("DoAction: modify full name failed:", err)
		return err
	}
	u.setPropString(&u.FullName, "FullName", name)
	return nil
}

func (u *User) SetHomeDir(dbusMsg dbus.DMessage, home string) error {
	logger.Debug("[SetHomeDir] new home:", home)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	if !dutils.IsFileExist(home) {
		err := fmt.Errorf("Not found the home path: %s", home)
		return err
	}

	pid := dbusMsg.GetSenderPID()
	if err := u.accessAuthentication(pid, true); err != nil {
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

	shells := getAvailableShells("/etc/shells")
	if len(shells) == 0 {
		err := fmt.Errorf("No available shell found")
		logger.Error("[SetShell] failed:", err)
		return err
	}

	if !strv.Strv(shells).Contains(shell) {
		err := fmt.Errorf("Not found the shell: %s", shell)
		logger.Warning("[SetShell] failed:", err)
		return err
	}

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
	if err := u.accessAuthentication(pid, true); err != nil {
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
	if err := polkitAuthLogin(pid); err != nil {
		logger.Debug("[SetLocked] access denied:", err)
		return err
	}

	if err := users.LockedUser(locked, u.UserName); err != nil {
		logger.Warning("DoAction: locked user failed:", err)
		return err
	}

	if locked && u.AutomaticLogin {
		users.SetAutoLoginUser("", "")
	}
	u.setPropBool(&u.Locked, "Locked", locked)
	return nil
}

func (u *User) SetAutomaticLogin(dbusMsg dbus.DMessage, auto bool) error {
	logger.Debug("[SetAutomaticLogin] auto", auto)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid := dbusMsg.GetSenderPID()
	if err := polkitAuthLogin(pid); err != nil {
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

	session := u.XSession
	if session == "" {
		session = getUserSession(u.HomeDir)
	}
	if err := users.SetAutoLoginUser(name, session); err != nil {
		logger.Warning("DoAction: set auto login failed:", err)
		return err
	}

	u.setPropBool(&u.AutomaticLogin, "AutomaticLogin", auto)
	return nil
}

func (u *User) EnableNoPasswdLogin(dbusMsg dbus.DMessage, enabled bool) error {
	logger.Debug("[EnableNoPasswdLogin] enabled:", enabled)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid := dbusMsg.GetSenderPID()
	if err := polkitAuthLogin(pid); err != nil {
		logger.Debug("[EnableNoPasswdLogin] access denied:", err)
		return err
	}

	if u.Locked {
		return fmt.Errorf("user %s has been locked", u.UserName)
	}

	if u.NoPasswdLogin == enabled {
		return nil
	}

	if err := users.EnableNoPasswdLogin(u.UserName, enabled); err != nil {
		logger.Warning("DoAction: enable nopasswdlogin failed:", err)
		return err
	}

	u.setPropBool(&u.NoPasswdLogin, "NoPasswdLogin", users.CanNoPasswdLogin(u.UserName))
	return nil
}

func (u *User) SetLocale(dbusMsg dbus.DMessage, locale string) error {
	logger.Debug("[SetLocale] locale:", locale)
	pid := dbusMsg.GetSenderPID()
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetLocale] access denied:", err)
			return err
		}
	}

	if u.Locale == locale {
		return nil
	}

	if !lang_info.IsSupportedLocale(locale) {
		err := fmt.Errorf("Invalid locale %q", locale)
		logger.Debug("[SetLocale]", err)
		return err
	}

	oldLocale := u.Locale
	u.setPropString(&u.Locale, "Locale", locale)
	if err := u.writeUserConfig(); err != nil {
		logger.Warning("[SetLocale]", err)
		u.setPropString(&u.Locale, "Locale", oldLocale)
		return err
	}

	return nil
}

func (u *User) SetLayout(dbusMsg dbus.DMessage, layout string) error {
	logger.Debug("[SetLayout] new layout:", layout)
	pid := dbusMsg.GetSenderPID()
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetLayout] access denied:", err)
			return err
		}
	}

	if u.Layout == layout {
		return nil
	}

	// TODO: check layout validity
	oldLayout := u.Layout
	u.setPropString(&u.Layout, "Layout", layout)
	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
		u.setPropString(&u.Layout, "Layout", oldLayout)
		return err
	}
	return nil
}

func (u *User) SetIconFile(dbusMsg dbus.DMessage, iconURI string) error {
	logger.Debug("[SetIconFile] new icon:", iconURI)
	pid := dbusMsg.GetSenderPID()
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetIconFile] access denied:", err)
			return err
		}
	}

	iconURI = dutils.EncodeURI(iconURI, dutils.SCHEME_FILE)
	iconFile := dutils.DecodeURI(iconURI)
	if u.IconFile == iconURI {
		return nil
	}

	if !gdkpixbuf.IsSupportedImage(iconFile) {
		err := fmt.Errorf("%q is not a image file", iconFile)
		logger.Debug(err)
		return err
	}

	newIconURI, added, err := u.setIconFile(iconURI)
	if err != nil {
		logger.Warning("Set icon failed:", err)
		return err
	}

	oldIconURI := u.IconFile
	oldCustomIcon := u.customIcon
	u.setPropIconFile(newIconURI)

	// Whether we need to remove the old custom icon
	var removeOld bool

	if added {
		// newIconURI should be custom icon if added is true
		u.customIcon = newIconURI
		if isUserCustomIconURI(oldIconURI, u.UserName) {
			// old and new icons are custom icon, we need remove the old custom icon
			removeOld = true
		}
	}

	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
		// recover
		u.setPropIconFile(oldIconURI)
		u.customIcon = oldCustomIcon
		removeOld = false
		return err
	}

	if removeOld {
		logger.Debugf("remove old custom icon %q", oldIconURI)
		err := os.Remove(dutils.DecodeURI(oldIconURI))
		if err != nil {
			logger.Warning(err)
		}
	}

	if added {
		u.updateIconList()
	}
	return nil
}

func (u *User) DeleteIconFile(dbusMsg dbus.DMessage, icon string) error {
	logger.Debug("[DeleteIconFile] icon:", icon)
	pid := dbusMsg.GetSenderPID()
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[DeleteIconFile] access denied:", err)
			return err
		}
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

	oldCustomIcon := u.customIcon
	u.customIcon = ""
	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
		// recover
		u.customIcon = oldCustomIcon
		return err
	}
	u.updateIconList()
	return nil
}

func (u *User) SetDesktopBackgrounds(dbusMsg dbus.DMessage, val []string) error {
	logger.Debugf("[SetDesktopBackgrounds] val: %#v", val)
	pid := dbusMsg.GetSenderPID()
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetDesktopBackgrounds] access denied:", err)
			return err
		}
	}

	if len(val) == 0 {
		return errors.New("val is empty")
	}

	var newVal = make([]string, len(val))
	for idx, file := range val {
		newVal[idx] = dutils.EncodeURI(file, dutils.SCHEME_FILE)
	}

	if strv.Strv(u.DesktopBackgrounds).Equal(newVal) {
		return nil
	}

	oldVal := u.DesktopBackgrounds
	u.setPropStrv(&u.DesktopBackgrounds, confKeyDesktopBackgrounds, newVal)
	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
		u.setPropStrv(&u.DesktopBackgrounds, confKeyDesktopBackgrounds, oldVal)
		return err
	}
	return nil
}

func (u *User) SetGreeterBackground(dbusMsg dbus.DMessage, bg string) error {
	logger.Debug("[SetGreeterBackground] new background:", bg)
	pid := dbusMsg.GetSenderPID()
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetGreeterBackground] access denied:", err)
			return err
		}
	}
	bg = dutils.EncodeURI(bg, dutils.SCHEME_FILE)
	if u.GreeterBackground == bg {
		genGaussianBlur(bg)
		return nil
	}

	if !isBackgroundValid(bg) {
		err := ErrInvalidBackground{bg}
		logger.Warning(err)
		return err
	}

	oldGreeterBackground := u.GreeterBackground
	u.setPropString(&u.GreeterBackground, "GreeterBackground", bg)
	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
		u.setPropString(&u.GreeterBackground, "GreeterBackground", oldGreeterBackground)
		return err
	}

	genGaussianBlur(bg)
	return nil
}

func (u *User) SetHistoryLayout(dbusMsg dbus.DMessage, list []string) error {
	logger.Debug("[SetHistoryLayout] new history layout:", list)
	pid := dbusMsg.GetSenderPID()
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetHistoryLayout] access denied:", err)
			return err
		}
	}

	if isStrvEqual(u.HistoryLayout, list) {
		return nil
	}

	// TODO: check layout list whether validity
	oldHistoryLayout := u.HistoryLayout
	u.setPropStrv(&u.HistoryLayout, "HistoryLayout", list)
	if err := u.writeUserConfig(); err != nil {
		logger.Warning("Write user config failed:", err)
		u.setPropStrv(&u.HistoryLayout, "HistoryLayout", oldHistoryLayout)
	}
	return nil
}

func (u *User) IsIconDeletable(iconURI string) bool {
	if iconURI != u.IconFile && iconURI == u.customIcon {
		// iconURI is custom icon, and not current icon
		return true
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

var supportedFormats = strv.Strv([]string{"jpeg", "png", "bmp", "tiff"})

func isBackgroundValid(file string) bool {
	file = dutils.DecodeURI(file)
	format, err := graphic.SniffImageFormat(file)
	if err != nil {
		return false
	}

	if supportedFormats.Contains(format) {
		return true
	}
	return false
}

type ErrInvalidBackground struct {
	FileName string
}

func (err ErrInvalidBackground) Error() string {
	return fmt.Sprintf("%q is not a valid background file", err.FileName)
}
