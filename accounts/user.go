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
	"io/ioutil"
	"os"
	"path"
	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/graphic"
	dutils "pkg.deepin.io/lib/utils"
	"runtime/debug"
	"strings"
	"sync"
)

const (
	UserTypeStandard int32 = iota
	UserTypeAdmin
)

const (
	defaultUserIcon       = "/var/lib/AccountsService/icons/default.png"
	defaultUserBackground = "/usr/share/backgrounds/default_background.jpg"

	maxWidth  = 200
	maxHeight = 200
)

type User struct {
	UserName       string
	Uid            string
	Gid            string
	HomeDir        string
	Shell          string
	Language       string
	IconFile       string
	BackgroundFile string

	// 用户是否被禁用
	Locked bool
	// 是否允许此用户自动登录
	AutomaticLogin bool

	AccountType int32
	LoginTime   uint64

	IconList     []string
	HistoryIcons []string

	syncLocker sync.Mutex
}

func NewUser(userPath string) (*User, error) {
	info, err := users.GetUserInfoByUid(getUidFromUserPath(userPath))
	if err != nil {
		return nil, err
	}

	var u = &User{}
	u.setPropString(&u.UserName, "UserName", info.Name)
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

	u.setPropStrv(&u.IconList, "IconList", u.getAllIcons())

	kFile, err := dutils.NewKeyFileFromFile(
		path.Join(userConfigDir, info.Name))
	if err != nil {
		u.setPropString(&u.IconFile, "IconFile", defaultUserIcon)
		u.setPropString(&u.BackgroundFile, "BackgroundFile", defaultUserBackground)
		u.writeUserConfig()
		return u, nil
	}
	defer kFile.Free()

	var write bool = false
	lang, _ := kFile.GetString("User", "Language")
	u.setPropString(&u.Language, "Language", lang)
	if len(u.Language) == 0 {
		u.setPropString(&u.Language, "Language", getSystemLanguage(defaultLocaleFile))
		write = true
	}

	icon, _ := kFile.GetString("User", "Icon")
	u.setPropString(&u.IconFile, "IconFile", icon)
	if len(u.IconFile) == 0 {
		u.setPropString(&u.IconFile, "IconFile", defaultUserIcon)
		write = true
	}

	bg, _ := kFile.GetString("User", "Background")
	u.setPropString(&u.BackgroundFile, "BackgroundFile", bg)
	if len(u.BackgroundFile) == 0 {
		u.setPropString(&u.BackgroundFile, "BackgroundFile", defaultUserBackground)
		write = true
	}

	if write {
		u.writeUserConfig()
	}

	_, hisIcons, _ := kFile.GetStringList("User", "HistoryIcons")
	u.setPropStrv(&u.HistoryIcons, "HistoryIcons", hisIcons)

	return u, nil
}

func (u *User) destroy() {
	dbus.UnInstallObject(u)
}

func (u *User) getAllIcons() []string {
	icons := getUserStandardIcons()
	cusIcons := getUserCustomIcons(u.UserName)

	icons = append(icons, cusIcons...)
	return icons
}

func (u *User) addIconFile(icon string) (string, bool, error) {
	if isStrInArray(icon, u.IconList) {
		return icon, false, nil
	}

	md5, ok := dutils.SumFileMd5(icon)
	if !ok {
		return "", false, fmt.Errorf("Sum file '%s' md5 failed", icon)
	}

	tmp, scale, err := scaleUserIcon(icon, md5)
	if err != nil {
		return "", false, err
	}

	if scale {
		defer os.Remove(tmp)
	}

	dest := path.Join(userCustomIconsDir, u.UserName+"-"+md5)
	err = os.MkdirAll(path.Dir(dest), 0755)
	if err != nil {
		return "", false, err
	}
	err = dutils.CopyFile(tmp, dest)
	if err != nil {
		return "", false, err
	}

	return dest, true, nil
}

func (u *User) addHistoryIcon(icon string) {
	if len(icon) == 0 || icon == defaultUserIcon {
		return
	}

	icons := u.HistoryIcons
	if isStrInArray(icon, icons) {
		return
	}

	var list = []string{icon}
	for _, v := range icons {
		if len(list) >= 9 {
			break
		}

		list = append(list, v)
	}
	u.setPropStrv(&u.HistoryIcons, "HistoryIcons", list)
}

func (u *User) deleteHistoryIcon(icon string) {
	if len(icon) == 0 {
		return
	}

	icons := u.HistoryIcons
	var list []string
	for _, v := range icons {
		if v == icon {
			continue
		}
		list = append(list, v)
	}
	u.setPropStrv(&u.HistoryIcons, "HistoryIcons", list)
}

func (u *User) writeUserConfig() error {
	return doWriteUserConfig(path.Join(userConfigDir, u.UserName),
		u.IconFile, u.BackgroundFile, u.HistoryIcons)
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

func (u *User) accessAuthentication(pid uint32, check bool, action string) error {
	var self bool
	if check {
		uid, _ := getUidByPid(pid)
		if u.Uid == uid {
			self = true
		}
	}

	var err error
	if self {
		err = polkitAuthChangeOwnData(pid)
	} else {
		err = polkitAuthManagerUser(pid)
	}
	if err != nil {
		doEmitError(pid, action, err.Error())
		return err
	}

	return nil
}

func doWriteUserConfig(config, icon, bg string, hisIcons []string) error {
	if !dutils.IsFileExist(config) {
		err := dutils.CreateFile(config)
		if err != nil {
			return err
		}
	}

	kFile, err := dutils.NewKeyFileFromFile(config)
	if err != nil {
		return err
	}
	defer kFile.Free()

	kFile.SetString("User", "Icon", icon)
	kFile.SetString("User", "Background", bg)
	kFile.SetStringList("User", "HistoryIcons", hisIcons)
	_, content, err := kFile.ToData()

	return dutils.WriteStringToFile(config, content)
}

// userPath must be composed with 'userDBusPath + uid'
func getUidFromUserPath(userPath string) string {
	items := strings.Split(userPath, userDBusPath)

	return items[1]
}

func getSystemLanguage(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		line = strings.TrimSpace(line)
		array := strings.Split(line, "=")
		if len(array) != 2 {
			continue
		}

		if array[0] == "LANG" {
			tmp := strings.Split(array[1], ".")
			return tmp[0]
		}
	}
	return ""
}

func scaleUserIcon(file, md5 string) (string, bool, error) {
	w, h, err := graphic.GetImageSize(file)
	if err != nil {
		return "", false, err
	}

	if w < maxWidth && h < maxHeight {
		return file, false, nil
	}

	dest := path.Join("/tmp", md5)
	defer debug.FreeOSMemory()
	return dest, true, graphic.ScaleImagePrefer(file, dest,
		maxWidth, maxHeight, graphic.FormatPng)
}
