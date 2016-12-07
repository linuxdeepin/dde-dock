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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
	defaultLayout         = "us;"
	defaultUserIcon       = "file:///var/lib/AccountsService/icons/default.png"
	defaultUserBackground = "file:///usr/share/backgrounds/default_background.jpg"

	maxWidth  = 200
	maxHeight = 200
)

const (
	confGroupUser            string = "User"
	confKeyIcon                     = "Icon"
	confKeyLocale                   = "Locale"
	confKeyLayout                   = "Layout"
	confKeyBackground               = "Background"
	confKeyGreeterBackground        = "GreeterBackground"
	confKeyHistoryIcons             = "HistoryIcons"
	confKeyHistoryLayout            = "HistoryLayout"
)

type User struct {
	UserName          string
	Uid               string
	Gid               string
	HomeDir           string
	Shell             string
	Locale            string
	Layout            string
	IconFile          string
	BackgroundFile    string
	GreeterBackground string

	// 用户是否被禁用
	Locked bool
	// 是否允许此用户自动登录
	AutomaticLogin bool

	AccountType int32
	LoginTime   uint64

	IconList      []string
	HistoryLayout []string
	HistoryIcons  []string

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
		u.setPropString(&u.Layout, "Layout", defaultLayout)
		u.setPropString(&u.Locale, "Locale", getSystemLocale(defaultLocaleFile))
		u.setPropString(&u.IconFile, "IconFile", defaultUserIcon)
		u.setPropString(&u.BackgroundFile, "BackgroundFile", defaultUserBackground)
		u.setPropString(&u.GreeterBackground, "GreeterBackground", defaultUserBackground)
		u.writeUserConfig()
		return u, nil
	}
	defer kFile.Free()

	var isSave bool = false
	locale, _ := kFile.GetString(confGroupUser, confKeyLocale)
	u.setPropString(&u.Locale, "Locale", locale)
	if len(locale) == 0 {
		u.setPropString(&u.Locale, "Locale", getSystemLocale(defaultLocaleFile))
		isSave = true
	}
	layout, _ := kFile.GetString(confGroupUser, confKeyLayout)
	u.setPropString(&u.Layout, "Layout", layout)
	if len(layout) == 0 {
		u.setPropString(&u.Layout, "Layout", defaultLayout)
		isSave = true
	}
	icon, _ := kFile.GetString(confGroupUser, confKeyIcon)
	u.setPropString(&u.IconFile, "IconFile", icon)
	if len(u.IconFile) == 0 {
		u.setPropString(&u.IconFile, "IconFile", defaultUserIcon)
		isSave = true
	}

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
	_, hisIcons, _ := kFile.GetStringList(confGroupUser, confKeyHistoryIcons)
	u.setPropStrv(&u.HistoryIcons, "HistoryIcons", hisIcons)

	if isSave {
		u.writeUserConfig()
	}

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

	icon = dutils.DecodeURI(icon)
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

	return dutils.EncodeURI(dest, dutils.SCHEME_FILE), true, nil
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
		// for compatible reason, old config files may contain icon paths
		// that are not encoded.
		v = dutils.EncodeURI(v, dutils.SCHEME_FILE)
		if v == icon {
			continue
		}
		list = append(list, v)
	}
	u.setPropStrv(&u.HistoryIcons, "HistoryIcons", list)
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

	kFile.SetString(confGroupUser, confKeyLayout, u.Layout)
	kFile.SetString(confGroupUser, confKeyLocale, u.Locale)
	kFile.SetString(confGroupUser, confKeyIcon, u.IconFile)
	kFile.SetString(confGroupUser, confKeyBackground, u.BackgroundFile)
	kFile.SetString(confGroupUser, confKeyGreeterBackground, u.GreeterBackground)
	kFile.SetStringList(confGroupUser, confKeyHistoryIcons, u.HistoryIcons)
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
		return err
	}

	return nil
}

// userPath must be composed with 'userDBusPath + uid'
func getUidFromUserPath(userPath string) string {
	items := strings.Split(userPath, userDBusPath)

	return items[1]
}

func getSystemLocale(file string) string {
	// If file is big, please using bufio.Scanner
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
			return array[1]
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
