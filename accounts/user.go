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
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/gdkpixbuf"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	UserTypeStandard int32 = iota
	UserTypeAdmin
)

const (
	defaultUserIcon          = "file:///var/lib/AccountsService/icons/default.png"
	defaultUserBackgroundDir = "/usr/share/wallpapers/deepin/"

	maxWidth  = 200
	maxHeight = 200
)

const (
	confGroupUser             = "User"
	confKeyXSession           = "XSession"
	confKeySystemAccount      = "SystemAccount"
	confKeyIcon               = "Icon"
	confKeyCustomIcon         = "CustomIcon"
	confKeyLocale             = "Locale"
	confKeyLayout             = "Layout"
	confKeyDesktopBackgrounds = "DesktopBackgrounds"
	confKeyGreeterBackground  = "GreeterBackground"
	confKeyHistoryLayout      = "HistoryLayout"
)

var (
	defaultUserBackgroundOnce  sync.Once
	cacheDefaultUserBackground string
)

func getDefaultUserBackground() string {
	defaultUserBackgroundOnce.Do(func() {
		var baseName = "desktop.jpg"
		content, err := ioutil.ReadFile(filepath.Join(defaultUserBackgroundDir,
			"default.conf"))
		content = bytes.TrimSpace(content)
		if err == nil && len(content) != 0 {
			contentStr := string(content)
			_, err := os.Stat(filepath.Join(defaultUserBackgroundDir, contentStr))
			if err == nil {
				baseName = contentStr
			}
		}
		cacheDefaultUserBackground = "file://" + defaultUserBackgroundDir + baseName
		logger.Debug("default user background:", cacheDefaultUserBackground)
	})
	return cacheDefaultUserBackground
}

type User struct {
	service    *dbusutil.Service
	PropsMu    sync.RWMutex
	UserName   string
	FullName   string
	Uid        string
	Gid        string
	HomeDir    string
	Shell      string
	Locale     string
	Layout     string
	IconFile   string
	customIcon string
	// dbusutil-gen: equal=nil
	DesktopBackgrounds []string
	GreeterBackground  string
	XSession           string

	// 用户是否被禁用
	Locked bool
	// 是否允许此用户自动登录
	AutomaticLogin bool
	SystemAccount  bool

	NoPasswdLogin bool

	AccountType int32
	LoginTime   uint64

	// dbusutil-gen: equal=nil
	IconList []string
	// dbusutil-gen: equal=nil
	HistoryLayout []string

	syncLocker   sync.Mutex
	configLocker sync.Mutex

	methods *struct {
		SetFullName           func() `in:"name"`
		SetHomeDir            func() `in:"home"`
		SetShell              func() `in:"shell"`
		SetPassword           func() `in:"password"`
		SetAccountType        func() `in:"accountType"`
		SetLocked             func() `in:"locked"`
		SetAutomaticLogin     func() `in:"enabled"`
		EnableNoPasswdLogin   func() `in:"enabled"`
		SetLocale             func() `in:"locale"`
		SetLayout             func() `in:"layout"`
		SetIconFile           func() `in:"iconFile"`
		DeleteIconFile        func() `in:"iconFile"`
		SetDesktopBackgrounds func() `in:"backgrounds"`
		SetGreeterBackground  func() `in:"background"`
		SetHistoryLayout      func() `in:"layouts"`
		IsIconDeletable       func() `in:"icon"`
		GetLargeIcon          func() `out:"icon"`
	}
}

func NewUser(userPath string, service *dbusutil.Service) (*User, error) {
	userInfo, err := users.GetUserInfoByUid(getUidFromUserPath(userPath))
	if err != nil {
		return nil, err
	}

	var u = &User{
		service:        service,
		UserName:       userInfo.Name,
		FullName:       userInfo.Comment().FullName(),
		Uid:            userInfo.Uid,
		Gid:            userInfo.Gid,
		HomeDir:        userInfo.Home,
		Shell:          userInfo.Shell,
		AutomaticLogin: users.IsAutoLoginUser(userInfo.Name),
		NoPasswdLogin:  users.CanNoPasswdLogin(userInfo.Name),
		Locked:         users.IsUserLocked(userInfo.Name),
	}

	u.AccountType = u.getAccountType()
	u.IconList = u.getAllIcons()

	updateConfigPath(userInfo.Name)

	kf, err := dutils.NewKeyFileFromFile(
		path.Join(userConfigDir, userInfo.Name))
	if err != nil {
		xSession, _ := users.GetDefaultXSession()
		u.XSession = xSession
		u.SystemAccount = false
		u.Layout = getDefaultLayout()
		u.Locale = getDefaultLocale()
		u.IconFile = defaultUserIcon
		defaultUserBackground := getDefaultUserBackground()
		u.DesktopBackgrounds = []string{defaultUserBackground}
		u.GreeterBackground = defaultUserBackground
		u.writeUserConfig()
		return u, nil
	}
	defer kf.Free()

	var isSave = false
	xSession, _ := kf.GetString(confGroupUser, confKeyXSession)
	u.XSession = xSession
	if u.XSession == "" {
		xSession, _ = users.GetDefaultXSession()
		u.XSession = xSession
		isSave = true
	}
	_, err = kf.GetBoolean(confGroupUser, confKeySystemAccount)
	// only show non system account
	u.SystemAccount = false
	if err != nil {
		isSave = true
	}
	locale, _ := kf.GetString(confGroupUser, confKeyLocale)
	u.Locale = locale
	if locale == "" {
		u.Locale = getDefaultLocale()
		isSave = true
	}
	layout, _ := kf.GetString(confGroupUser, confKeyLayout)
	u.Layout = layout
	if layout == "" {
		u.Layout = getDefaultLayout()
		isSave = true
	}
	icon, _ := kf.GetString(confGroupUser, confKeyIcon)
	u.IconFile = icon
	if u.IconFile == "" {
		u.IconFile = defaultUserIcon
		isSave = true
	}

	u.customIcon, _ = kf.GetString(confGroupUser, confKeyCustomIcon)

	// CustomIcon is the newly added field in the configuration file
	if u.customIcon == "" {
		if u.IconFile != defaultUserIcon && !isStrInArray(u.IconFile, u.IconList) {
			// u.IconFile is a custom icon, not a standard icon
			u.customIcon = u.IconFile
			isSave = true
		}
	}

	u.IconList = u.getAllIcons()

	_, desktopBgs, _ := kf.GetStringList(confGroupUser, confKeyDesktopBackgrounds)
	u.DesktopBackgrounds = desktopBgs
	if len(desktopBgs) == 0 {
		u.DesktopBackgrounds = []string{getDefaultUserBackground()}
		isSave = true
	}

	greeterBg, _ := kf.GetString(confGroupUser, confKeyGreeterBackground)
	u.GreeterBackground = greeterBg
	if greeterBg == "" {
		u.GreeterBackground = getDefaultUserBackground()
		isSave = true
	}

	_, u.HistoryLayout, _ = kf.GetStringList(confGroupUser, confKeyHistoryLayout)

	if isSave {
		u.writeUserConfig()
	}

	u.checkLeftSpace()
	return u, nil
}

func (u *User) updateIconList() {
	u.IconList = u.getAllIcons()
	u.emitPropChangedIconList(u.IconList)
}

func (u *User) getAllIcons() []string {
	icons := getUserStandardIcons()
	if u.customIcon != "" {
		icons = append(icons, u.customIcon)
	}
	return icons
}

// ret0: new user icon uri
// ret1: added
// ret2: error
func (u *User) setIconFile(iconURI string) (string, bool, error) {
	if isStrInArray(iconURI, u.IconList) {
		return iconURI, false, nil
	}

	iconFile := dutils.DecodeURI(iconURI)
	tmp, scaled, err := scaleUserIcon(iconFile)
	if err != nil {
		return "", false, err
	}

	if scaled {
		logger.Debug("icon scaled", tmp)
		defer os.Remove(tmp)
	}

	dest := getNewUserCustomIconDest(u.UserName)
	err = os.MkdirAll(path.Dir(dest), 0755)
	if err != nil {
		return "", false, err
	}
	err = dutils.CopyFile(tmp, dest)
	if err != nil {
		return "", false, err
	}
	return dutils.EncodeURI(dest, dutils.SCHEME_FILE), true, nil
}

type configChange struct {
	key   string
	value interface{} // allowed type are bool, string, []string
}

func (u *User) writeUserConfigWithChanges(changes []configChange) error {
	u.configLocker.Lock()
	defer u.configLocker.Unlock()

	err := os.MkdirAll(userConfigDir, 0755)
	if err != nil {
		return err
	}

	config := path.Join(userConfigDir, u.UserName)
	if !dutils.IsFileExist(config) {
		err := dutils.CreateFile(config)
		if err != nil {
			return err
		}
	}

	kf, err := dutils.NewKeyFileFromFile(config)
	if err != nil {
		logger.Warningf("Load %s config file failed: %v", u.UserName, err)
		return err
	}
	defer kf.Free()

	kf.SetString(confGroupUser, confKeyXSession, u.XSession)
	kf.SetBoolean(confGroupUser, confKeySystemAccount, u.SystemAccount)
	kf.SetString(confGroupUser, confKeyLayout, u.Layout)
	kf.SetString(confGroupUser, confKeyLocale, u.Locale)
	kf.SetString(confGroupUser, confKeyIcon, u.IconFile)
	kf.SetString(confGroupUser, confKeyCustomIcon, u.customIcon)
	kf.SetStringList(confGroupUser, confKeyDesktopBackgrounds, u.DesktopBackgrounds)
	kf.SetString(confGroupUser, confKeyGreeterBackground, u.GreeterBackground)
	kf.SetStringList(confGroupUser, confKeyHistoryLayout, u.HistoryLayout)

	for _, change := range changes {
		switch val := change.value.(type) {
		case bool:
			kf.SetBoolean(confGroupUser, change.key, val)
		case string:
			kf.SetString(confGroupUser, change.key, val)
		case []string:
			kf.SetStringList(confGroupUser, change.key, val)
		default:
			return errors.New("unsupported value type")
		}
	}

	_, err = kf.SaveToFile(config)
	if err != nil {
		logger.Warningf("Save %s config file failed: %v", u.UserName, err)
	}
	return err
}

func (u *User) writeUserConfigWithChange(confKey string, value interface{}) error {
	return u.writeUserConfigWithChanges([]configChange{
		{confKey, value},
	})
}

func (u *User) writeUserConfig() error {
	return u.writeUserConfigWithChanges(nil)
}

func (u *User) updatePropLocked() {
	newVal := users.IsUserLocked(u.UserName)
	u.PropsMu.Lock()
	u.setPropLocked(newVal)
	u.PropsMu.Unlock()
}

func (u *User) updatePropAccountType() {
	newVal := u.getAccountType()
	u.PropsMu.Lock()
	u.setPropAccountType(newVal)
	u.PropsMu.Unlock()
}

func (u *User) updatePropCanNoPasswdLogin() {
	newVal := users.CanNoPasswdLogin(u.UserName)
	u.PropsMu.Lock()
	u.setPropNoPasswdLogin(newVal)
	u.PropsMu.Unlock()
}

func (u *User) updatePropAutomaticLogin() {
	newVal := users.IsAutoLoginUser(u.UserName)
	u.PropsMu.Lock()
	u.setPropAutomaticLogin(newVal)
	u.PropsMu.Unlock()
}

func (u *User) getAccountType() int32 {
	if users.IsAdminUser(u.UserName) {
		return UserTypeAdmin
	}
	return UserTypeStandard
}

func (u *User) checkAuth(sender dbus.Sender, selfPass bool, actionId string) error {
	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		return err
	}

	isSelf := u.isSelf(pid)
	if selfPass && isSelf {
		return nil
	}

	if actionId == "" {
		if isSelf {
			actionId = polkitActionChangeOwnData
		} else {
			actionId = polkitActionUserAdministration
		}
	}

	return checkAuth(actionId, pid)
}

func (u *User) checkAuthAutoLogin(sender dbus.Sender, enabled bool) error {
	var actionId string
	if enabled {
		actionId = polkitActionEnableAutoLogin
	} else {
		actionId = polkitActionDisableAutoLogin
	}

	return u.checkAuth(sender, false, actionId)
}

func (u *User) checkAuthNoPasswdLogin(sender dbus.Sender, enabled bool) error {
	var actionId string
	if enabled {
		actionId = polkitActionEnableNoPasswordLogin
	} else {
		actionId = polkitActionDisableNoPasswordLogin
	}
	return u.checkAuth(sender, false, actionId)
}

func (u *User) isSelf(pid uint32) bool {
	uid, _ := getUidByPid(pid)
	return u.Uid == uid
}

func (u *User) clearData() {
	// delete user config file
	configFile := path.Join(userConfigDir, u.UserName)
	err := os.Remove(configFile)
	if err != nil {
		logger.Warningf("remove user config failed:", err)
	}

	// delete user custom icon
	if u.customIcon != "" {
		customIconFile := dutils.DecodeURI(u.customIcon)
		err := os.Remove(customIconFile)
		if err != nil {
			logger.Warning("remove user custom icon failed:", err)
		}
	}
}

// userPath must be composed with 'userDBusPath + uid'
func getUidFromUserPath(userPath string) string {
	items := strings.Split(userPath, userDBusPathPrefix)

	return items[1]
}

// ret0: output file
// ret1: scaled
// ret2: error
func scaleUserIcon(file string) (string, bool, error) {
	w, h, err := gdkpixbuf.GetImageSize(file)
	if err != nil {
		return "", false, err
	}

	if w <= maxWidth && h <= maxHeight {
		return file, false, nil
	}

	dest, err := getTempFile()
	if err != nil {
		return "", false, err
	}

	err = gdkpixbuf.ScaleImagePrefer(file, dest, maxWidth, maxHeight, gdkpixbuf.GDK_INTERP_BILINEAR, gdkpixbuf.FormatPng)
	if err != nil {
		return "", false, err
	}

	return dest, true, nil
}

// return temp file path and error
func getTempFile() (string, error) {
	tmpfile, err := ioutil.TempFile("", "dde-daemon-accounts")
	if err != nil {
		return "", err
	}
	name := tmpfile.Name()
	tmpfile.Close()
	return name, nil
}

func getUserSession(homeDir string) string {
	session, ok := dutils.ReadKeyFromKeyFile(homeDir+"/.dmrc", "Desktop", "Session", "")
	if !ok {
		v := ""
		list := getSessionList()
		switch len(list) {
		case 0:
			v = ""
		case 1:
			v = list[0]
		default:
			if strv.Strv(list).Contains("deepin.desktop") {
				v = "deepin.desktop"
			} else {
				v = list[0]
			}
		}
		return v
	}
	return session.(string)
}

func getSessionList() []string {
	fileInfoList, err := ioutil.ReadDir("/usr/share/xsessions")
	if err != nil {
		return nil
	}

	var sessions []string
	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() || !strings.Contains(fileInfo.Name(), ".desktop") {
			continue
		}
		sessions = append(sessions, fileInfo.Name())
	}
	return sessions
}

// 迁移配置文件，复制文件从 $actConfigDir/users/$username 到 $userConfigDir/$username
func updateConfigPath(username string) {
	config := path.Join(userConfigDir, username)
	if dutils.IsFileExist(config) {
		return
	}

	err := os.MkdirAll(userConfigDir, 0755)
	if err != nil {
		logger.Warning("Failed to mkdir for user config:", err)
		return
	}

	oldConfig := path.Join(actConfigDir, "users", username)
	err = dutils.CopyFile(oldConfig, config)
	if err != nil {
		logger.Warning("Failed to update config:", username)
	}
}
