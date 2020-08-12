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
	"pkg.deepin.io/lib/imgutil"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	userDBusPathPrefix = "/com/deepin/daemon/Accounts/User"
	userDBusInterface  = "com.deepin.daemon.Accounts.User"
)

func (*User) GetInterfaceName() string {
	return userDBusInterface
}

func (u *User) SetFullName(sender dbus.Sender, name string) *dbus.Error {
	logger.Debugf("[SetFullName] new name: %q", name)

	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetFullName] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if u.FullName != name {
		if err := users.ModifyFullName(name, u.UserName); err != nil {
			logger.Warning("DoAction: modify full name failed:", err)
			return dbusutil.ToError(err)
		}

		u.FullName = name
		_ = u.emitPropChangedFullName(name)
	}

	return nil
}

func (u *User) SetHomeDir(sender dbus.Sender, home string) *dbus.Error {
	logger.Debug("[SetHomeDir] new home:", home)

	err := u.checkAuth(sender, false, "")
	if err != nil {
		logger.Debug("[SetHomeDir] access denied:", err)
		return dbusutil.ToError(err)
	}

	if dutils.IsFileExist(home) {
		// if new home already exists, the `usermod -m -d` command will fail.
		return dbusutil.ToError(errors.New("new home already exists"))
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if u.HomeDir != home {
		if err := users.ModifyHome(home, u.UserName); err != nil {
			logger.Warning("DoAction: modify home failed:", err)
			return dbusutil.ToError(err)
		}
		u.HomeDir = home
		_ = u.emitPropChangedHomeDir(home)
	}

	return nil
}

func (u *User) SetShell(sender dbus.Sender, shell string) *dbus.Error {
	logger.Debug("[SetShell] new shell:", shell)

	err := u.checkAuth(sender, false, "")
	if err != nil {
		logger.Debug("[SetShell] access denied:", err)
		return dbusutil.ToError(err)
	}

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

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if u.Shell != shell {
		if err := users.ModifyShell(shell, u.UserName); err != nil {
			logger.Warning("DoAction: modify shell failed:", err)
			return dbusutil.ToError(err)
		}
		u.Shell = shell
		_ = u.emitPropChangedShell(shell)
	}

	return nil
}

func (u *User) SetPassword(sender dbus.Sender, password string) *dbus.Error {
	logger.Debug("[SetPassword] start ...")

	err := u.checkAuth(sender, false, "")
	if err != nil {
		logger.Debug("[SetPassword] access denied:", err)
		return dbusutil.ToError(err)
	}

	if err := users.ModifyPasswd(password, u.UserName); err != nil {
		logger.Warning("DoAction: modify password failed:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if u.Locked {
		if err := users.LockedUser(false, u.UserName); err != nil {
			logger.Warning("DoAction: unlock user failed:", err)
			return dbusutil.ToError(err)
		}
		u.Locked = false
		_ = u.emitPropChangedLocked(false)
	}
	return nil
}

func (u *User) SetMaxPasswordAge(sender dbus.Sender, nDays int32) *dbus.Error {
	err := u.checkAuth(sender, false, "")
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	err = users.ModifyMaxPasswordAge(u.UserName, int(nDays))
	if err != nil {
		logger.Warning("failed to set max password age:", err)
		return dbusutil.ToError(err)
	}

	return nil
}

func (u *User) IsPasswordExpired() (bool, *dbus.Error) {
	v, err := users.IsPasswordExpired(u.UserName)
	return v, dbusutil.ToError(err)
}

func (u *User) SetLocked(sender dbus.Sender, locked bool) *dbus.Error {
	logger.Debug("[SetLocked] locked:", locked)

	err := u.checkAuth(sender, false, polkitActionUserAdministration)
	if err != nil {
		logger.Debug("[SetLocked] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if u.Locked != locked {
		if err := users.LockedUser(locked, u.UserName); err != nil {
			logger.Warning("DoAction: locked user failed:", err)
			return dbusutil.ToError(err)
		}

		u.Locked = locked
		_ = u.emitPropChangedLocked(locked)

		if locked && u.AutomaticLogin {
			if err := users.SetAutoLoginUser("", ""); err != nil {
				logger.Warning("failed to clear auto login user:", err)
				return dbusutil.ToError(err)
			}
			u.AutomaticLogin = false
			_ = u.emitPropChangedAutomaticLogin(false)
		}
	}
	return nil
}

func (u *User) AddGroup(sender dbus.Sender, group string) *dbus.Error {
	logger.Debugf("add group %s for %s", group, u.UserName)
	err := u.checkAuth(sender, false, polkitActionUserAdministration)
	if err != nil {
		logger.Debug("[AddGroup] access denied:", err)
		return dbusutil.ToError(err)
	}

	err = users.AddGroupForUser(group, u.UserName)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}

func (u *User) DeleteGroup(sender dbus.Sender, group string) *dbus.Error {
	logger.Debugf("delete group %s for %s", group, u.UserName)
	err := u.checkAuth(sender, false, polkitActionUserAdministration)
	if err != nil {
		logger.Debug("[DeleteGroup] access denied:", err)
		return dbusutil.ToError(err)
	}

	err = users.DeleteGroupForUser(group, u.UserName)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}

func (u *User) SetGroups(sender dbus.Sender, groups []string) *dbus.Error {
	logger.Debugf("set groups %v for %s", groups, u.UserName)
	err := u.checkAuth(sender, false, polkitActionUserAdministration)
	if err != nil {
		logger.Debug("[SetGroups] access denied:", err)
		return dbusutil.ToError(err)
	}

	err = users.SetGroupsForUser(groups, u.UserName)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}

func (u *User) SetAutomaticLogin(sender dbus.Sender, enabled bool) *dbus.Error {
	logger.Debug("[SetAutomaticLogin] enable:", enabled)

	err := u.checkAuthAutoLogin(sender, enabled)
	if err != nil {
		logger.Debug("[SetAutomaticLogin] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if u.Locked {
		return dbusutil.ToError(fmt.Errorf("user %s has been locked", u.UserName))
	}

	if u.AutomaticLogin == enabled {
		return nil
	}

	var name = u.UserName
	if !enabled {
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

	u.AutomaticLogin = enabled
	_ = u.emitPropChangedAutomaticLogin(enabled)
	return nil
}

func (u *User) EnableNoPasswdLogin(sender dbus.Sender, enabled bool) *dbus.Error {
	logger.Debug("[EnableNoPasswdLogin] enabled:", enabled)

	err := u.checkAuthNoPasswdLogin(sender, enabled)
	if err != nil {
		logger.Debug("[EnableNoPasswdLogin] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

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

	u.NoPasswdLogin = enabled
	_ = u.emitPropChangedNoPasswdLogin(enabled)
	return nil
}

func (u *User) SetLocale(sender dbus.Sender, locale string) *dbus.Error {
	logger.Debug("[SetLocale] locale:", locale)

	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetLocale] access denied:", err)
		return dbusutil.ToError(err)
	}

	if !lang_info.IsSupportedLocale(locale) {
		err := fmt.Errorf("invalid locale %q", locale)
		logger.Debug("[SetLocale]", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if u.Locale == locale {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyLocale, locale)
	if err != nil {
		return dbusutil.ToError(err)
	}
	u.Locale = locale
	_ = u.emitPropChangedLocale(locale)
	return nil
}

func (u *User) SetLayout(sender dbus.Sender, layout string) *dbus.Error {
	logger.Debug("[SetLayout] new layout:", layout)

	err := u.checkAuth(sender, true, polkitActionSetKeyboardLayout)
	if err != nil {
		logger.Debug("[SetLayout] access denied:", err)
		return dbusutil.ToError(err)
	}

	// TODO: check layout validity

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if u.Layout == layout {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyLayout, layout)
	if err != nil {
		return dbusutil.ToError(err)
	}
	u.Layout = layout
	_ = u.emitPropChangedLayout(layout)
	return nil
}

func (u *User) SetIconFile(sender dbus.Sender, iconURI string) *dbus.Error {
	logger.Debug("[SetIconFile] new icon:", iconURI)

	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetIconFile] access denied:", err)
		return dbusutil.ToError(err)
	}

	iconURI = dutils.EncodeURI(iconURI, dutils.SCHEME_FILE)
	iconFile := dutils.DecodeURI(iconURI)

	if !gdkpixbuf.IsSupportedImage(iconFile) {
		err := fmt.Errorf("%q is not a image file", iconFile)
		logger.Debug(err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if u.IconFile == iconURI {
		return nil
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

	u.IconFile = newIconURI
	_ = u.emitPropChangedIconFile(newIconURI)
	return nil
}

// 只能删除不是用户当前图标的自定义图标
func (u *User) DeleteIconFile(sender dbus.Sender, icon string) *dbus.Error {
	logger.Debug("[DeleteIconFile] icon:", icon)

	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[DeleteIconFile] access denied:", err)
		return dbusutil.ToError(err)
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

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

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

	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetDesktopBackgrounds] access denied:", err)
		return dbusutil.ToError(err)
	}

	if len(val) == 0 {
		return dbusutil.ToError(errors.New("val is empty"))
	}

	var newVal = make([]string, len(val))
	for idx, file := range val {
		newVal[idx] = dutils.EncodeURI(file, dutils.SCHEME_FILE)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if strv.Strv(u.DesktopBackgrounds).Equal(newVal) {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyDesktopBackgrounds, newVal)
	if err != nil {
		return dbusutil.ToError(err)
	}

	u.DesktopBackgrounds = newVal
	_ = u.emitPropChangedDesktopBackgrounds(newVal)
	return nil
}

func (u *User) SetGreeterBackground(sender dbus.Sender, bg string) *dbus.Error {
	logger.Debug("[SetGreeterBackground] new background:", bg)

	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetGreeterBackground] access denied:", err)
		return dbusutil.ToError(err)
	}

	bg = dutils.EncodeURI(bg, dutils.SCHEME_FILE)

	if !isBackgroundValid(bg) {
		err := ErrInvalidBackground{bg}
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer func() {
		u.PropsMu.Unlock()
		genGaussianBlur(bg)
	}()

	if u.GreeterBackground == bg {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyGreeterBackground, bg)
	if err != nil {
		return dbusutil.ToError(err)
	}
	u.GreeterBackground = bg
	_ = u.emitPropChangedGreeterBackground(bg)
	return nil
}

func (u *User) SetHistoryLayout(sender dbus.Sender, list []string) *dbus.Error {
	logger.Debug("[SetHistoryLayout] new history layout:", list)

	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetHistoryLayout] access denied:", err)
		return dbusutil.ToError(err)
	}

	// TODO: check layout list whether validity

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if isStrvEqual(u.HistoryLayout, list) {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyHistoryLayout, list)
	if err != nil {
		return dbusutil.ToError(err)
	}
	u.HistoryLayout = list
	_ = u.emitPropChangedHistoryLayout(list)
	return nil
}

func (u *User) SetUse24HourFormat(sender dbus.Sender, value bool) *dbus.Error {
	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetUse24HourFormat] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if value == u.Use24HourFormat {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyUse24HourFormat, value)
	if err != nil {
		return dbusutil.ToError(err)
	}

	u.Use24HourFormat = value
	err = u.emitPropChangedUse24HourFormat(value)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}

func (u *User) IsIconDeletable(iconURI string) bool {
	u.PropsMu.RLock()
	defer u.PropsMu.RUnlock()

	if iconURI != u.IconFile && iconURI == u.customIcon {
		// iconURI is custom icon, and not current icon
		return true
	}
	return true
}

// 获取当前头像的大图标
func (u *User) GetLargeIcon() string {
	u.PropsMu.RLock()
	baseName := path.Base(u.IconFile)
	dir := path.Dir(dutils.DecodeURI(u.IconFile))
	u.PropsMu.RUnlock()

	filename := path.Join(dir, "bigger", baseName)
	if !dutils.IsFileExist(filename) {
		return ""
	}

	return dutils.EncodeURI(filename, dutils.SCHEME_FILE)
}

var supportedFormats = strv.Strv([]string{"gif", "jpeg", "png", "bmp", "tiff"})

func isBackgroundValid(file string) bool {
	file = dutils.DecodeURI(file)
	format, err := imgutil.SniffFormat(file)
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

func (u *User) SetWeekdayFormat(sender dbus.Sender, value int32) *dbus.Error {
	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetWeekdayFormat] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if value == u.WeekdayFormat {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyWeekdayFormat, value)
	if err != nil {
		return dbusutil.ToError(err)
	}

	u.WeekdayFormat = value
	err = u.emitPropChangedWeekdayFormat(value)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}

func (u *User) SetShortDateFormat(sender dbus.Sender, value int32) *dbus.Error {
	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetShortDateFormat] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if value == u.ShortDateFormat {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyShortDateFormat, value)
	if err != nil {
		return dbusutil.ToError(err)
	}

	u.ShortDateFormat = value
	err = u.emitPropChangedShortDateFormat(value)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}

func (u *User) SetLongDateFormat(sender dbus.Sender, value int32) *dbus.Error {
	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetLongDateFormat] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if value == u.LongDateFormat {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyLongDateFormat, value)
	if err != nil {
		return dbusutil.ToError(err)
	}

	u.LongDateFormat = value
	err = u.emitPropChangedLongDateFormat(value)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}

func (u *User) SetShortTimeFormat(sender dbus.Sender, value int32) *dbus.Error {
	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetShortTimeFormat] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if value == u.ShortTimeFormat {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyShortTimeFormat, value)
	if err != nil {
		return dbusutil.ToError(err)
	}

	u.ShortTimeFormat = value
	err = u.emitPropChangedShortTimeFormat(value)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}

func (u *User) SetLongTimeFormat(sender dbus.Sender, value int32) *dbus.Error {
	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetLongTimeFormat] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if value == u.LongTimeFormat {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyLongTimeFormat, value)
	if err != nil {
		return dbusutil.ToError(err)
	}

	u.LongTimeFormat = value
	err = u.emitPropChangedLongTimeFormat(value)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}

func (u *User) SetWeekBegins(sender dbus.Sender, value int32) *dbus.Error {
	err := u.checkAuth(sender, true, "")
	if err != nil {
		logger.Debug("[SetWeekBegins] access denied:", err)
		return dbusutil.ToError(err)
	}

	u.PropsMu.Lock()
	defer u.PropsMu.Unlock()

	if value == u.WeekBegins {
		return nil
	}

	err = u.writeUserConfigWithChange(confKeyWeekBegins, value)
	if err != nil {
		return dbusutil.ToError(err)
	}

	u.WeekBegins = value
	err = u.emitPropChangedWeekBegins(value)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}
