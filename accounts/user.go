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
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/graphic"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
)

const (
	UserTypeStandard int32 = iota
	UserTypeAdmin
)

const (
	layoutDelim           = ";"
	defaultLayout         = "us" + layoutDelim
	defaultLayoutFile     = "/etc/default/keyboard"
	defaultUserIcon       = "file:///var/lib/AccountsService/icons/default.png"
	defaultUserBackground = "file:///usr/share/backgrounds/default_background.jpg"

	maxWidth  = 200
	maxHeight = 200
)

const (
	confGroupUser            string = "User"
	confKeyXSession                 = "XSession"
	confKeySystemAccount            = "SystemAccount"
	confKeyIcon                     = "Icon"
	confKeyCustomIcon               = "CustomIcon"
	confKeyLocale                   = "Locale"
	confKeyLayout                   = "Layout"
	confKeyBackground               = "Background"
	confKeyGreeterBackground        = "GreeterBackground"
	confKeyHistoryIcons             = "HistoryIcons"
	confKeyHistoryLayout            = "HistoryLayout"
)

type User struct {
	UserName          string
	FullName          string
	Uid               string
	Gid               string
	HomeDir           string
	Shell             string
	Locale            string
	Layout            string
	IconFile          string
	customIcon        string
	BackgroundFile    string
	GreeterBackground string
	XSession          string

	// 用户是否被禁用
	Locked bool
	// 是否允许此用户自动登录
	AutomaticLogin bool
	SystemAccount  bool

	AccountType int32
	LoginTime   uint64

	IconList      []string
	HistoryLayout []string

	syncLocker   sync.Mutex
	configLocker sync.Mutex
}

func NewUser(userPath string) (*User, error) {
	info, err := users.GetUserInfoByUid(getUidFromUserPath(userPath))
	if err != nil {
		return nil, err
	}

	var u = &User{}
	u.setPropString(&u.UserName, "UserName", info.Name)

	comment := info.Comment()
	u.setPropString(&u.FullName, "FullName", comment.FullName())

	u.setPropString(&u.Uid, "Uid", info.Uid)
	u.setPropString(&u.Gid, "Gid", info.Gid)
	u.setPropString(&u.HomeDir, "HomeDir", info.Home)
	u.setPropString(&u.Shell, "Shell", info.Shell)
	u.setPropString(&u.IconFile, "IconFile", "")
	u.setPropString(&u.BackgroundFile, "BackgroundFile", "")

	u.setPropBool(&u.AutomaticLogin, "AutomaticLogin",
		users.IsAutoLoginUser(info.Name))

	u.updatePropLocked()
	u.updatePropAccountType()
	u.updateIconList()

	kFile, err := dutils.NewKeyFileFromFile(
		path.Join(userConfigDir, info.Name))
	if err != nil {
		xsession, _ := users.GetDefaultXSession()
		u.setPropString(&u.XSession, "XSession", xsession)
		u.setPropBool(&u.SystemAccount, "SystemAccount", false)
		u.setPropString(&u.Layout, "Layout", u.getDefaultLayout())
		u.setPropString(&u.Locale, "Locale", getLocaleFromFile(defaultLocaleFile))
		u.setPropString(&u.IconFile, "IconFile", defaultUserIcon)
		u.setPropString(&u.BackgroundFile, "BackgroundFile", defaultUserBackground)
		u.setPropString(&u.GreeterBackground, "GreeterBackground", defaultUserBackground)
		u.writeUserConfig()
		return u, nil
	}
	defer kFile.Free()

	var isSave bool = false
	xsession, _ := kFile.GetString(confGroupUser, confKeyXSession)
	u.setPropString(&u.XSession, "XSession", xsession)
	if u.XSession == "" {
		xsession, _ = users.GetDefaultXSession()
		u.setPropString(&u.XSession, "XSession", xsession)
		isSave = true
	}
	_, err = kFile.GetBoolean(confGroupUser, confKeySystemAccount)
	// only show non system account
	u.setPropBool(&u.SystemAccount, "SystemAccount", false)
	if err != nil {
		isSave = true
	}
	locale, _ := kFile.GetString(confGroupUser, confKeyLocale)
	u.setPropString(&u.Locale, "Locale", locale)
	if len(locale) == 0 {
		u.setPropString(&u.Locale, "Locale", getLocaleFromFile(defaultLocaleFile))
		isSave = true
	}
	layout, _ := kFile.GetString(confGroupUser, confKeyLayout)
	u.setPropString(&u.Layout, "Layout", layout)
	if len(layout) == 0 {
		u.setPropString(&u.Layout, "Layout", u.getDefaultLayout())
		isSave = true
	}
	icon, _ := kFile.GetString(confGroupUser, confKeyIcon)
	u.setPropString(&u.IconFile, "IconFile", icon)
	if len(u.IconFile) == 0 {
		u.setPropString(&u.IconFile, "IconFile", defaultUserIcon)
		isSave = true
	}

	u.customIcon, _ = kFile.GetString(confGroupUser, confKeyCustomIcon)

	// CustomInfo is the newly added field in the configuration file
	if u.customIcon == "" {
		if u.IconFile != defaultUserIcon && !isStrInArray(u.IconFile, u.IconList) {
			// u.IconFile is a custom icon, not a standard icon
			u.customIcon = u.IconFile
			isSave = true
		}
	}

	u.updateIconList()

	bg, _ := kFile.GetString(confGroupUser, confKeyBackground)
	u.setPropString(&u.BackgroundFile, "BackgroundFile", bg)
	if len(bg) == 0 {
		u.setPropString(&u.BackgroundFile, "BackgroundFile", defaultUserBackground)
		isSave = true
	}
	greeterBg, _ := kFile.GetString(confGroupUser, confKeyGreeterBackground)
	u.setPropString(&u.GreeterBackground, "GreeterBackground", greeterBg)
	if len(greeterBg) == 0 {
		u.setPropString(&u.GreeterBackground, "GreeterBackground", defaultUserBackground)
		isSave = true
	}

	_, hisLayout, _ := kFile.GetStringList(confGroupUser, confKeyHistoryLayout)
	u.setPropStrv(&u.HistoryLayout, "HistoryLayout", hisLayout)

	if isSave {
		u.writeUserConfig()
	}

	u.checkLeftSpace()
	return u, nil
}

func (u *User) destroy() {
	dbus.UnInstallObject(u)
}

func (u *User) updateIconList() {
	u.setPropStrv(&u.IconList, "IconList", u.getAllIcons())
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

func (u *User) writeUserConfig() error {
	u.configLocker.Lock()
	defer u.configLocker.Unlock()

	config := path.Join(userConfigDir, u.UserName)
	if !dutils.IsFileExist(config) {
		err := dutils.CreateFile(config)
		if err != nil {
			return err
		}
	}

	kFile, err := dutils.NewKeyFileFromFile(config)
	if err != nil {
		logger.Warningf("Load %s config file failed: %v", u.UserName, err)
		return err
	}
	defer kFile.Free()

	kFile.SetString(confGroupUser, confKeyXSession, u.XSession)
	kFile.SetBoolean(confGroupUser, confKeySystemAccount, u.SystemAccount)
	kFile.SetString(confGroupUser, confKeyLayout, u.Layout)
	kFile.SetString(confGroupUser, confKeyLocale, u.Locale)
	kFile.SetString(confGroupUser, confKeyIcon, u.IconFile)
	kFile.SetString(confGroupUser, confKeyCustomIcon, u.customIcon)
	kFile.SetString(confGroupUser, confKeyBackground, u.BackgroundFile)
	kFile.SetString(confGroupUser, confKeyGreeterBackground, u.GreeterBackground)
	kFile.SetStringList(confGroupUser, confKeyHistoryLayout, u.HistoryLayout)
	_, err = kFile.SaveToFile(config)
	if err != nil {
		logger.Warningf("Save %s config file failed: %v", u.UserName, err)
	}
	return err
}

func (u *User) updatePropLocked() {
	u.setPropBool(&u.Locked, "Locked", users.IsUserLocked(u.UserName))
}

func (u *User) updatePropAccountType() {
	if users.IsAdminUser(u.UserName) {
		u.setPropInt32(&u.AccountType, "AccountType", UserTypeAdmin)
	} else {
		u.setPropInt32(&u.AccountType, "AccountType", UserTypeStandard)
	}
}

func (u *User) accessAuthentication(pid uint32, check bool) error {
	var err error
	if check && u.isSelf(pid) {
		err = polkitAuthChangeOwnData(pid)
	} else {
		err = polkitAuthManagerUser(pid)
	}
	if err != nil {
		return err
	}

	return nil
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

func (u *User) getDefaultLayout() string {
	layout, err := getSystemLayout(defaultLayoutFile)
	if err != nil {
		logger.Warning("Failed to get system default layout:", err)
		return defaultLayout
	}
	return layout
}

// userPath must be composed with 'userDBusPath + uid'
func getUidFromUserPath(userPath string) string {
	items := strings.Split(userPath, userDBusPath)

	return items[1]
}

// ret0: output file
// ret1: scaled
// ret2: error
func scaleUserIcon(file string) (string, bool, error) {
	w, h, err := graphic.GetImageSize(file)
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

	defer debug.FreeOSMemory()
	return dest, true, graphic.ScaleImagePrefer(file, dest,
		maxWidth, maxHeight, graphic.FormatPng)
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

var (
	gaussianLocker sync.Mutex
	gaussianTasks  = make(map[string]bool)
)

func genGaussianBlur(file string) {
	gaussianLocker.Lock()
	file = dutils.DecodeURI(file)
	logger.Debug("[genGaussianBlur] task manager:", gaussianTasks)
	_, ok := gaussianTasks[file]
	if ok {
		logger.Debug("[genGaussianBlur] tash exists:", file)
		gaussianLocker.Unlock()
		return
	}
	gaussianTasks[file] = true
	gaussianLocker.Unlock()

	go func() {
		logger.Debug("[genGaussianBlur] will blur image:", file)
		exec.Command("/usr/lib/deepin-api/image-blur-helper",
			file).CombinedOutput()
		gaussianLocker.Lock()
		delete(gaussianTasks, file)
		gaussianLocker.Unlock()
	}()
}

func getSystemLayout(file string) (string, error) {
	fr, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	var (
		found   int
		layout  string
		variant string

		regLayout  = regexp.MustCompile(`^XKBLAYOUT=`)
		regVariant = regexp.MustCompile(`^XKBVARIANT=`)

		scanner = bufio.NewScanner(fr)
	)
	for scanner.Scan() {
		if found == 2 {
			break
		}

		var line = scanner.Text()
		if regLayout.MatchString(line) {
			layout = strings.Trim(getValueFromLine(line, "="), "\"")
			found += 1
			continue
		}

		if regVariant.MatchString(line) {
			variant = strings.Trim(getValueFromLine(line, "="), "\"")
			found += 1
		}
	}

	if len(layout) == 0 {
		return "", fmt.Errorf("Not found default layout")
	}

	return layout + layoutDelim + variant, nil
}

func getValueFromLine(line, delim string) string {
	array := strings.Split(line, delim)
	if len(array) != 2 {
		return ""
	}

	return strings.TrimSpace(array[1])
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
	finfos, err := ioutil.ReadDir("/usr/share/xsessions")
	if err != nil {
		return nil
	}

	var sessions []string
	for _, finfo := range finfos {
		if finfo.IsDir() || !strings.Contains(finfo.Name(), ".desktop") {
			continue
		}
		sessions = append(sessions, finfo.Name())
	}
	return sessions
}
