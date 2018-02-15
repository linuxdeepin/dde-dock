/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/gdkpixbuf"
	"pkg.deepin.io/lib/graphic"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	userDBusPath = "/com/deepin/daemon/Accounts/User"
	userDBusIFC  = "com.deepin.daemon.Accounts.User"
)

func (u *User) GetDBusExportInfo() dbusutil.ExportInfo {
	return dbusutil.ExportInfo{
		Path:      userDBusPath + u.Uid,
		Interface: userDBusIFC,
	}
}

func (u *User) SetFullName(sender dbus.Sender, name string) *dbus.Error {
	logger.Debugf("[SetFullName] new name: %q", name)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	if err := u.accessAuthentication(sender, true); err != nil {
		logger.Debug("[SetFullName] access denied:", err)
		return dbusutil.ToError(err)
	}

	if err := users.ModifyFullName(name, u.UserName); err != nil {
		logger.Warning("DoAction: modify full name failed:", err)
		return dbusutil.ToError(err)
	}
	u.setPropFullName(name)
	return nil
}

func (u *User) SetHomeDir(sender dbus.Sender, home string) *dbus.Error {
	logger.Debug("[SetHomeDir] new home:", home)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	if dutils.IsFileExist(home) {
		// if new home already exists, the `usermod -m -d` command will fail.
		return dbusutil.ToError(errors.New("new home already exists"))
	}

	if err := u.accessAuthentication(sender, true); err != nil {
		logger.Debug("[SetHomeDir] access denied:", err)
		return dbusutil.ToError(err)
	}

	if err := users.ModifyHome(home, u.UserName); err != nil {
		logger.Warning("DoAction: modify home failed:", err)
		return dbusutil.ToError(err)
	}

	u.setPropHomeDir(home)
	return nil
}

func (u *User) SetShell(sender dbus.Sender, shell string) *dbus.Error {
	logger.Debug("[SetShell] new shell:", shell)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	shells := getAvailableShells("/etc/shells")
	if len(shells) == 0 {
		err := fmt.Errorf("no available shell found")
		logger.Error("[SetShell] failed:", err)
		return dbusutil.ToError(err)
	}

	if !strv.Strv(shells).Contains(shell) {
		err := fmt.Errorf("not found the shell: %s", shell)
		logger.Warning("[SetShell] failed:", err)
		return dbusutil.ToError(err)
	}

	if err := u.accessAuthentication(sender, true); err != nil {
		logger.Debug("[SetShell] access denied:", err)
		return dbusutil.ToError(err)
	}

	if err := users.ModifyShell(shell, u.UserName); err != nil {
		logger.Warning("DoAction: modify shell failed:", err)
		return dbusutil.ToError(err)
	}

	u.setPropShell(shell)
	return nil
}

func (u *User) SetPassword(sender dbus.Sender, password string) *dbus.Error {
	logger.Debug("[SetPassword] start ...")
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	if err := u.accessAuthentication(sender, true); err != nil {
		logger.Debug("[SetPassword] access denied:", err)
		return dbusutil.ToError(err)
	}

	if err := users.ModifyPasswd(password, u.UserName); err != nil {
		logger.Warning("DoAction: modify password failed:", err)
		return dbusutil.ToError(err)
	}

	if err := users.LockedUser(false, u.UserName); err != nil {
		logger.Warning("DoAction: unlock user failed:", err)
		return dbusutil.ToError(err)
	}
	u.setPropLocked(false)
	return nil
}

func (u *User) SetAccountType(sender dbus.Sender, ty int32) *dbus.Error {
	logger.Debug("[SetAccountType] type:", ty)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	if err := u.accessAuthentication(sender, false); err != nil {
		logger.Debug("[SetAccountType] access denied:", err)
		return dbusutil.ToError(err)
	}

	if err := users.SetUserType(ty, u.UserName); err != nil {
		logger.Warning("DoAction: set user type failed:", err)
		return dbusutil.ToError(err)
	}

	u.setPropAccountType(ty)
	return nil
}

func (u *User) SetLocked(sender dbus.Sender, locked bool) *dbus.Error {
	logger.Debug("[SetLocked] locked:", locked)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	if err := polkitAuthChangeOwnData("", "", pid); err != nil {
		logger.Debug("[SetLocked] access denied:", err)
		return dbusutil.ToError(err)
	}

	if err := users.LockedUser(locked, u.UserName); err != nil {
		logger.Warning("DoAction: locked user failed:", err)
		return dbusutil.ToError(err)
	}

	if locked && u.AutomaticLogin {
		users.SetAutoLoginUser("", "")
	}
	u.setPropLocked(locked)
	return nil
}

func (u *User) SetAutomaticLogin(sender dbus.Sender, auto bool) *dbus.Error {
	logger.Debug("[SetAutomaticLogin] auto", auto)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	if err := polkitAuthAutoLogin(pid, auto); err != nil {
		logger.Debug("[SetAutomaticLogin] access denied:", err)
		return dbusutil.ToError(err)
	}

	if u.Locked {
		return dbusutil.ToError(fmt.Errorf("user %s has been locked", u.UserName))
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
		return dbusutil.ToError(err)
	}

	u.setPropAutomaticLogin(auto)
	return nil
}

func (u *User) EnableNoPasswdLogin(sender dbus.Sender, enabled bool) *dbus.Error {
	logger.Debug("[EnableNoPasswdLogin] enabled:", enabled)
	u.syncLocker.Lock()
	defer u.syncLocker.Unlock()

	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	if err := polkitAuthNoPasswdLogin(pid, enabled); err != nil {
		logger.Debug("[EnableNoPasswdLogin] access denied:", err)
		return dbusutil.ToError(err)
	}

	if u.Locked {
		return dbusutil.ToError(fmt.Errorf("user %s has been locked", u.UserName))
	}

	if u.NoPasswdLogin == enabled {
		return nil
	}

	if err := users.EnableNoPasswdLogin(u.UserName, enabled); err != nil {
		logger.Warning("DoAction: enable no password login failed:", err)
		return dbusutil.ToError(err)
	}

	u.setPropNoPasswdLogin(enabled)
	return nil
}

func (u *User) SetLocale(sender dbus.Sender, locale string) *dbus.Error {
	logger.Debug("[SetLocale] locale:", locale)
	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetLocale] access denied:", err)
			return dbusutil.ToError(err)
		}
	}

	if u.Locale == locale {
		return nil
	}

	if !lang_info.IsSupportedLocale(locale) {
		err := fmt.Errorf("invalid locale %q", locale)
		logger.Debug("[SetLocale]", err)
		return dbusutil.ToError(err)
	}

	err = u.writeUserConfigWithChange(confKeyLocale, locale)
	if err != nil {
		return dbusutil.ToError(err)
	}
	u.setPropLocale(locale)
	return nil
}

func (u *User) SetLayout(sender dbus.Sender, layout string) *dbus.Error {
	logger.Debug("[SetLayout] new layout:", layout)
	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetLayout] access denied:", err)
			return dbusutil.ToError(err)
		}
	}

	if u.Layout == layout {
		return nil
	}

	// TODO: check layout validity
	err = u.writeUserConfigWithChange(confKeyLayout, layout)
	if err != nil {
		return dbusutil.ToError(err)
	}
	u.setPropLayout(layout)
	return nil
}

func (u *User) SetIconFile(sender dbus.Sender, iconURI string) *dbus.Error {
	logger.Debug("[SetIconFile] new icon:", iconURI)
	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetIconFile] access denied:", err)
			return dbusutil.ToError(err)
		}
	}

	iconURI = dutils.EncodeURI(iconURI, dutils.SCHEME_FILE)
	iconFile := dutils.DecodeURI(iconURI)

	currentIconFile := u.getPropIconFile()
	if currentIconFile == iconURI {
		return nil
	}

	if !gdkpixbuf.IsSupportedImage(iconFile) {
		err := fmt.Errorf("%q is not a image file", iconFile)
		logger.Debug(err)
		return dbusutil.ToError(err)
	}

	newIconURI, added, err := u.setIconFile(iconURI)
	if err != nil {
		logger.Warning("Set icon failed:", err)
		return dbusutil.ToError(err)
	}

	if added {
		// newIconURI should be custom icon if added is true
		err = u.writeUserConfigWithChanges([]configChange{
			{confKeyCustomIcon, newIconURI},
			{confKeyIcon, newIconURI},
		})
		if err != nil {
			return dbusutil.ToError(err)
		}

		// remove old custom icon
		if u.customIcon != "" {
			logger.Debugf("remove old custom icon %q", u.customIcon)
			err := os.Remove(dutils.DecodeURI(u.customIcon))
			if err != nil {
				logger.Warning(err)
			}
		}
		u.customIcon = newIconURI
		u.updateIconList()
	} else {
		err = u.writeUserConfigWithChange(confKeyIcon, newIconURI)
		if err != nil {
			return dbusutil.ToError(err)
		}
	}

	u.setPropIconFile(newIconURI)
	return nil
}

// 只能删除不是用户当前图标的自定义图标
func (u *User) DeleteIconFile(sender dbus.Sender, icon string) *dbus.Error {
	logger.Debug("[DeleteIconFile] icon:", icon)
	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[DeleteIconFile] access denied:", err)
			return dbusutil.ToError(err)
		}
	}

	icon = dutils.EncodeURI(icon, dutils.SCHEME_FILE)
	if !u.IsIconDeletable(icon) {
		err := errors.New("this icon is not allowed to be deleted")
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	iconPath := dutils.DecodeURI(icon)
	if err := os.Remove(iconPath); err != nil {
		return dbusutil.ToError(err)
	}

	// set custom icon to empty
	err = u.writeUserConfigWithChange(confKeyCustomIcon, "")
	if err != nil {
		return dbusutil.ToError(err)
	}
	u.customIcon = ""

	u.updateIconList()
	return nil
}

func (u *User) SetDesktopBackgrounds(sender dbus.Sender, val []string) *dbus.Error {
	logger.Debugf("[SetDesktopBackgrounds] val: %#v", val)
	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetDesktopBackgrounds] access denied:", err)
			return dbusutil.ToError(err)
		}
	}

	if len(val) == 0 {
		return dbusutil.ToError(errors.New("val is empty"))
	}

	var newVal = make([]string, len(val))
	for idx, file := range val {
		newVal[idx] = dutils.EncodeURI(file, dutils.SCHEME_FILE)
	}

	if strv.Strv(u.DesktopBackgrounds).Equal(newVal) {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyDesktopBackgrounds, newVal)
	if err != nil {
		return dbusutil.ToError(err)
	}

	u.setPropDesktopBackgrounds(newVal)
	return nil
}

func (u *User) SetGreeterBackground(sender dbus.Sender, bg string) *dbus.Error {
	logger.Debug("[SetGreeterBackground] new background:", bg)
	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetGreeterBackground] access denied:", err)
			return dbusutil.ToError(err)
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
		return dbusutil.ToError(err)
	}

	err = u.writeUserConfigWithChange(confKeyGreeterBackground, bg)
	if err != nil {
		return dbusutil.ToError(err)
	}
	u.setPropGreeterBackground(bg)

	genGaussianBlur(bg)
	return nil
}

func (u *User) SetHistoryLayout(sender dbus.Sender, list []string) *dbus.Error {
	logger.Debug("[SetHistoryLayout] new history layout:", list)
	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	if !u.isSelf(pid) {
		err := polkitAuthManagerUser(pid)
		if err != nil {
			logger.Debug("[SetHistoryLayout] access denied:", err)
			return dbusutil.ToError(err)
		}
	}

	if isStrvEqual(u.HistoryLayout, list) {
		return nil
	}

	// TODO: check layout list whether validity
	err = u.writeUserConfigWithChange(confKeyHistoryLayout, list)
	if err != nil {
		return dbusutil.ToError(err)
	}
	u.setPropHistoryLayout(list)

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
