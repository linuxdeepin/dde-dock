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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	authenticate "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.authenticate"
	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/gir/glib-2.0"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/gdkpixbuf"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
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
	confKeyUse24HourFormat    = "Use24HourFormat"
	confKeyUUID               = "UUID"
	confKeyWeekdayFormat      = "WeekdayFormat"
	confKeyShortDateFormat    = "ShortDateFormat"
	confKeyLongDateFormat     = "LongDateFormat"
	confKeyShortTimeFormat    = "ShortTimeFormat"
	confKeyLongTimeFormat     = "LongTimeFormat"
	confKeyWeekBegins         = "WeekBegins"

	defaultUse24HourFormat = true
	defaultWeekdayFormat   = 0
	defaultShortDateFormat = 0
	defaultLongDateFormat  = 0
	defaultShortTimeFormat = 0
	defaultLongTimeFormat  = 0
	defaultWeekBegins      = 0
)

func getDefaultUserBackground() string {
	filename := filepath.Join(defaultUserBackgroundDir, "desktop.bmp")
	_, err := os.Stat(filename)
	if err == nil {
		return "file://" + filename
	}

	return "file://" + filepath.Join(defaultUserBackgroundDir, "desktop.jpg")
}

type User struct {
	service         *dbusutil.Service
	PropsMu         sync.RWMutex
	UserName        string
	UUID            string
	FullName        string
	Uid             string
	Gid             string
	HomeDir         string
	Shell           string
	Locale          string
	Layout          string
	IconFile        string
	Use24HourFormat bool
	WeekdayFormat   int32
	ShortDateFormat int32
	LongDateFormat  int32
	ShortTimeFormat int32
	LongTimeFormat  int32
	WeekBegins      int32

	customIcon string
	// dbusutil-gen: equal=nil
	DesktopBackgrounds []string
	// dbusutil-gen: equal=isStrvEqual
	Groups            []string
	GreeterBackground string
	XSession          string

	PasswordStatus     string
	MaxPasswordAge     int32
	PasswordLastChange int32
	// 用户是否被禁用
	Locked bool
	// 是否允许此用户自动登录
	AutomaticLogin bool

	// deprecated property
	SystemAccount bool

	NoPasswdLogin bool

	AccountType int32
	LoginTime   uint64
	CreatedTime uint64

	// dbusutil-gen: equal=nil
	IconList []string
	// dbusutil-gen: equal=nil
	HistoryLayout []string

	configLocker sync.Mutex
	//nolint
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
		AddGroup              func() `in:"group"`
		DeleteGroup           func() `in:"group"`
		SetGroups             func() `in:"groups"`
		SetUse24HourFormat    func() `in:"value"`
		SetMaxPasswordAge     func() `in:"nDays"`
		IsPasswordExpired     func() `out:"expired"`
		SetWeekdayFormat      func() `in:"value"`
		SetShortDateFormat    func() `in:"value"`
		SetLongDateFormat     func() `in:"value"`
		SetShortTimeFormat    func() `in:"value"`
		SetLongTimeFormat     func() `in:"value"`
		SetWeekBegins         func() `in:"value"`
	}
}

func NewUser(userPath string, service *dbusutil.Service) (*User, error) {
	userInfo, err := users.GetUserInfoByUid(getUidFromUserPath(userPath))
	if err != nil {
		return nil, err
	}

	shadowInfo, err := users.GetShadowInfo(userInfo.Name)
	if err != nil {
		return nil, err
	}

	var u = &User{
		service:            service,
		UserName:           userInfo.Name,
		FullName:           userInfo.Comment().FullName(),
		Uid:                userInfo.Uid,
		Gid:                userInfo.Gid,
		HomeDir:            userInfo.Home,
		Shell:              userInfo.Shell,
		AutomaticLogin:     users.IsAutoLoginUser(userInfo.Name),
		NoPasswdLogin:      users.CanNoPasswdLogin(userInfo.Name),
		Locked:             shadowInfo.Status == users.PasswordStatusLocked,
		PasswordStatus:     shadowInfo.Status,
		MaxPasswordAge:     int32(shadowInfo.MaxDays),
		PasswordLastChange: int32(shadowInfo.LastChange),
	}

	u.AccountType = u.getAccountType()
	u.IconList = u.getAllIcons()
	u.Groups = u.getGroups()

	// NOTICE(jouyouyun): Got created time,  not accurate, can only be used as a reference
	u.CreatedTime, err = u.getCreatedTime()
	if err != nil {
		logger.Warning("Failed to get created time:", err)
	}

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
		u.Use24HourFormat = defaultUse24HourFormat
		u.UUID = dutils.GenUuid()
		u.WeekdayFormat = defaultWeekdayFormat
		u.ShortDateFormat = defaultShortDateFormat
		u.LongDateFormat = defaultLongDateFormat
		u.ShortTimeFormat = defaultShortTimeFormat
		u.LongTimeFormat = defaultLongTimeFormat
		u.WeekBegins = defaultWeekBegins

		err = u.writeUserConfig()
		if err != nil {
			logger.Warning(err)
		}
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

	greeterBg, ok := getUserGreeterBackground(kf)
	if ok {
		u.GreeterBackground = greeterBg
	} else {
		u.GreeterBackground = getDefaultUserBackground()
		isSave = true
	}

	_, u.HistoryLayout, _ = kf.GetStringList(confGroupUser, confKeyHistoryLayout)
	if !strv.Strv(u.HistoryLayout).Contains(u.Layout) {
		u.HistoryLayout = append(u.HistoryLayout, u.Layout)
		isSave = true
	}

	u.Use24HourFormat, err = kf.GetBoolean(confGroupUser, confKeyUse24HourFormat)
	if err != nil {
		u.Use24HourFormat = defaultUse24HourFormat
		isSave = true
	}

	u.WeekdayFormat, err = kf.GetInteger(confGroupUser, confKeyWeekdayFormat)
	if err != nil {
		u.WeekdayFormat = defaultWeekdayFormat
		isSave = true
	}

	u.ShortDateFormat, err = kf.GetInteger(confGroupUser, confKeyShortDateFormat)
	if err != nil {
		u.ShortDateFormat = defaultShortDateFormat
		isSave = true
	}

	u.LongDateFormat, err = kf.GetInteger(confGroupUser, confKeyLongDateFormat)
	if err != nil {
		u.LongDateFormat = defaultLongDateFormat
		isSave = true
	}

	u.ShortTimeFormat, err = kf.GetInteger(confGroupUser, confKeyShortTimeFormat)
	if err != nil {
		u.ShortTimeFormat = defaultShortTimeFormat
		isSave = true
	}

	u.LongTimeFormat, err = kf.GetInteger(confGroupUser, confKeyLongTimeFormat)
	if err != nil {
		u.LongTimeFormat = defaultLongTimeFormat
		isSave = true
	}

	u.WeekBegins, err = kf.GetInteger(confGroupUser, confKeyWeekBegins)
	if err != nil {
		u.WeekBegins = defaultWeekBegins
		isSave = true
	}

	u.UUID, err = kf.GetString(confGroupUser, confKeyUUID)
	if err != nil || u.UUID == "" {
		u.UUID = dutils.GenUuid()
		isSave = true
	}

	if isSave {
		err := u.writeUserConfig()
		if err != nil {
			logger.Warning(err)
		}
	}

	u.checkLeftSpace()
	return u, nil
}

func getUserGreeterBackground(kf *glib.KeyFile) (string, bool) {
	greeterBg, _ := kf.GetString(confGroupUser, confKeyGreeterBackground)
	if greeterBg == "" {
		return "", false
	}
	_, err := os.Stat(dutils.DecodeURI(greeterBg))
	if err != nil {
		logger.Warning(err)
		return "", false
	}
	return greeterBg, true
}

func (u *User) updateIconList() {
	u.IconList = u.getAllIcons()
	_ = u.emitPropChangedIconList(u.IconList)
}

func (u *User) getAllIcons() []string {
	icons := getUserStandardIcons()
	if u.customIcon != "" {
		icons = append(icons, u.customIcon)
	}
	return icons
}

func (u *User) getGroups() []string {
	groups, err := users.GetUserGroups(u.UserName)
	if err != nil {
		logger.Warning("failed to get user groups:", err)
		return nil
	}
	return groups
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
		defer func() {
			_ = os.Remove(tmp)
		}()
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
	value interface{} // allowed type are bool, string, []string , int32
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
	kf.SetString(confGroupUser, confKeyUUID, u.UUID)

	for _, change := range changes {
		switch val := change.value.(type) {
		case bool:
			kf.SetBoolean(confGroupUser, change.key, val)
		case string:
			kf.SetString(confGroupUser, change.key, val)
		case []string:
			kf.SetStringList(confGroupUser, change.key, val)
		case int32:
			kf.SetInteger(confGroupUser, change.key, val)
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

func (u *User) updatePropGroups() {
	newVal := u.getGroups()
	u.PropsMu.Lock()
	u.setPropGroups(newVal)
	u.PropsMu.Unlock()
}

func (u *User) updatePropAutomaticLogin() {
	newVal := users.IsAutoLoginUser(u.UserName)
	u.PropsMu.Lock()
	u.setPropAutomaticLogin(newVal)
	u.PropsMu.Unlock()
}

func (u *User) updatePropsPasswd(uInfo *users.UserInfo) {
	var userNameChanged bool
	var oldUserName string

	u.PropsMu.Lock()
	u.setPropGid(uInfo.Gid)

	if u.UserName != uInfo.Name {
		oldUserName = u.UserName
		userNameChanged = true
	}
	u.setPropUserName(uInfo.Name)

	u.setPropHomeDir(uInfo.Home)
	u.setPropShell(uInfo.Shell)
	fullName := uInfo.Comment().FullName()
	u.setPropFullName(fullName)
	u.PropsMu.Unlock()

	if userNameChanged {
		logger.Debugf("user name changed old: %q, new: %q", oldUserName, uInfo.Name)
		err := os.Rename(filepath.Join(userConfigDir, oldUserName),
			filepath.Join(userConfigDir, uInfo.Name))
		if err != nil {
			logger.Warning(err)
		}
	}
}

func (u *User) updatePropsShadow(shadowInfo *users.ShadowInfo) {
	u.PropsMu.Lock()

	u.setPropPasswordStatus(shadowInfo.Status)
	u.setPropLocked(shadowInfo.Status == users.PasswordStatusLocked)
	u.setPropMaxPasswordAge(int32(shadowInfo.MaxDays))
	u.setPropPasswordLastChange(int32(shadowInfo.LastChange))

	u.PropsMu.Unlock()
}

func (u *User) getAccountType() int32 {
	if users.IsAdminUser(u.UserName) {
		return users.UserTypeAdmin
	}
	return users.UserTypeStandard
}

func (u *User) checkAuth(sender dbus.Sender, selfPass bool, actionId string) error {
	uid, err := u.service.GetConnUID(string(sender))
	if err != nil {
		return err
	}

	isSelf := u.isSelf(uid)
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

	return checkAuth(actionId, string(sender))
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

func (u *User) isSelf(uid uint32) bool {
	return u.Uid == strconv.FormatInt(int64(uid), 10)
}

func (u *User) clearData() {
	// delete user config file
	configFile := path.Join(userConfigDir, u.UserName)
	err := os.Remove(configFile)
	if err != nil {
		logger.Warning("remove user config failed:", err)
	}

	// delete user custom icon
	if u.customIcon != "" {
		customIconFile := dutils.DecodeURI(u.customIcon)
		err := os.Remove(customIconFile)
		if err != nil {
			logger.Warning("remove user custom icon failed:", err)
		}
	}

	u.clearFingers()
}

func (u *User) clearFingers() {
	logger.Debug("clearFingers")

	sysBus, err := dbus.SystemBus()
	if err != nil {
		logger.Warning("connect to system bus failed:", err)
		return
	}

	fpObj := authenticate.NewFingerprint(sysBus)
	err = fpObj.DeleteAllFingers(0, u.UserName)
	if err != nil {
		logger.Warning("failed to delete enrolled fingers:", err)
		return
	}

	logger.Debug("clear fingers succesed")
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
	_ = tmpfile.Close()
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
