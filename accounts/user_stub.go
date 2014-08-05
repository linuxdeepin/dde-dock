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
	"path"
	"pkg.linuxdeepin.com/lib/dbus"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"strings"
)

func (obj *User) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		ACCOUNT_DEST,
		obj.objectPath,
		USER_MANAGER_IFC,
	}
}

func (obj *User) updatePropUserName(name string) {
	if len(name) < 1 {
		return
	}

	if obj.UserName != name {
		obj.UserName = name
		dbus.NotifyChange(obj, "UserName")
	}
}

func (obj *User) updatePropHomeDir(homeDir string) {
	if len(homeDir) < 1 {
		return
	}

	if obj.HomeDir != homeDir {
		obj.HomeDir = homeDir
		dbus.NotifyChange(obj, "HomeDir")
	}
}

func (obj *User) updatePropShell(shell string) {
	if len(shell) < 1 {
		return
	}

	if obj.Shell != shell {
		obj.Shell = shell
		dbus.NotifyChange(obj, "Shell")
	}
}

func (obj *User) updatePropIconFile(icon string) {
	if len(icon) < 1 {
		return
	}

	if obj.IconFile != icon {
		obj.IconFile = icon
		dbus.NotifyChange(obj, "IconFile")
	}
}

func (obj *User) updatePropBackgroundFile(bg string) {
	if len(bg) < 1 {
		return
	}

	if obj.BackgroundFile != bg {
		obj.BackgroundFile = bg
		dbus.NotifyChange(obj, "BackgroundFile")
	}
}

func (obj *User) updatePropAutomaticLogin(enable bool) {
	//if obj.AutomaticLogin != enable {
	obj.AutomaticLogin = enable
	dbus.NotifyChange(obj, "AutomaticLogin")
	//}
}

func (obj *User) updatePropAccountType(acctype int32) {
	//if obj.AccountType != acctype {
	obj.AccountType = acctype
	dbus.NotifyChange(obj, "AccountType")
	//}
}

func (obj *User) updatePropLocked(locked bool) {
	//if obj.Locked != locked {
	obj.Locked = locked
	dbus.NotifyChange(obj, "Locked")
	//}
}

func (obj *User) updatePropHistoryIcons(iconList []string) {
	if !isStrListEqual(obj.HistoryIcons, iconList) {
		obj.HistoryIcons = iconList
		dbus.NotifyChange(obj, "HistoryIcons")
	}
}

func (obj *User) updatePropIconList(iconList []string) {
	if !isStrListEqual(obj.IconList, iconList) {
		obj.IconList = iconList
		dbus.NotifyChange(obj, "IconList")
	}
}

func (obj *User) getPropIconFile() string {
	file := path.Join(USER_CONFIG_DIR, obj.UserName)
	icon := ""
	wFlag := false
	if !dutils.IsFileExist(file) {
		icon = getRandUserIcon()
		wFlag = true
	} else {
		if v, ok := dutils.ReadKeyFromKeyFile(file, "User",
			"Icon", ""); !ok {
			icon = getRandUserIcon()
			wFlag = true
		} else {
			icon = v.(string)
		}
	}

	if wFlag {
		dutils.WriteKeyToKeyFile(file, "User", "Icon", icon)
	}

	return icon
}

func (obj *User) getPropBackgroundFile() string {
	file := path.Join(USER_CONFIG_DIR, obj.UserName)
	bg := ""
	wFlag := false
	if !dutils.IsFileExist(file) {
		bg = USER_DEFAULT_BG
		wFlag = true
	} else {
		if v, ok := dutils.ReadKeyFromKeyFile(file, "User",
			"Background", ""); !ok {
			bg = USER_DEFAULT_BG
			wFlag = true
		} else {
			bg = v.(string)
		}
	}

	if wFlag {
		dutils.WriteKeyToKeyFile(file, "User", "Background", bg)
	}

	return bg
}

func (obj *User) getPropAutomaticLogin() bool {
	return isAutoLogin(obj.UserName)
}

func (obj *User) getPropAccountType() int32 {
	list := getAdministratorList()
	if strIsInList(obj.UserName, list) {
		return ACCOUNT_TYPE_ADMINISTACTOR
	}

	return ACCOUNT_TYPE_STANDARD
}

func (obj *User) getPropHistoryIcons() []string {
	list := []string{}
	file := path.Join(USER_CONFIG_DIR, obj.UserName)

	if !dutils.IsFileExist(file) {
		list = append(list, obj.IconFile)
	} else {
		if v, ok := dutils.ReadKeyFromKeyFile(file, "User",
			"HistoryIcons", []string{}); !ok {
			list = append(list, obj.IconFile)
		} else {
			list = append(list, v.([]string)...)
		}
	}

	tmp := []string{}
	for _, v := range list {
		if v == obj.IconFile {
			continue
		}
		tmp = append(tmp, v)
	}
	list = tmp

	return list
}

func (obj *User) getPropIconList() []string {
	list := []string{}

	sysList := getIconList(ICON_SYSTEM_DIR)
	list = append(list, sysList...)
	localList := getIconList(ICON_LOCAL_DIR)
	for _, l := range localList {
		if strings.Contains(l, obj.UserName+"-") {
			list = append(list, l)
		}
	}

	return list
}

func (obj *User) addHistoryIcon(iconPath string) {
	file := path.Join(USER_CONFIG_DIR, obj.UserName)
	if !dutils.IsFileExist(file) || !dutils.IsFileExist(iconPath) {
		return
	}

	list, _ := dutils.ReadKeyFromKeyFile(file, "User",
		"HistoryIcons", []string{})

	ret := []string{}
	ret = append(ret, iconPath)
	cnt := 1
	if list != nil {
		strs := list.([]string)
		as := deleteStrFromList(iconPath, strs)
		for _, l := range as {
			if cnt >= 10 {
				break
			}

			if ok := dutils.IsFileExist(l); !ok {
				continue
			}
			ret = append(ret, l)
			cnt++
		}
	}
	dutils.WriteKeyToKeyFile(file, "User", "HistoryIcons", ret)

	return
}

func (obj *User) deleteHistoryIcon(iconPath string) {
	file := path.Join(USER_CONFIG_DIR, obj.UserName)
	if !dutils.IsFileExist(file) {
		return
	}

	list, ok := dutils.ReadKeyFromKeyFile(file, "User",
		"HistoryIcons", []string{})
	if !ok || list == nil {
		return
	}

	tmp := deleteStrFromList(iconPath, list.([]string))
	dutils.WriteKeyToKeyFile(file, "User", "HistoryIcons", tmp)

	return
}

func (obj *User) setPropUserName(name string) {
	if len(name) < 1 {
		return
	}

	args := []string{}
	args = append(args, "-l")
	args = append(args, name)
	args = append(args, obj.UserName)
	execCommand(CMD_USERMOD, args)
}

func (obj *User) setPropHomeDir(homeDir string) {
	if len(homeDir) < 1 {
		return
	}

	args := []string{}
	args = append(args, "-m")
	args = append(args, "-d")
	args = append(args, homeDir)
	args = append(args, obj.UserName)
	execCommand(CMD_USERMOD, args)
}

func (obj *User) setPropShell(shell string) {
	if len(shell) < 1 {
		return
	}

	args := []string{}
	args = append(args, "-s")
	args = append(args, shell)
	args = append(args, obj.UserName)
	execCommand(CMD_USERMOD, args)
}

func (obj *User) setPropIconFile(icon string) {
	if len(icon) < 1 {
		return
	}

	file := path.Join(USER_CONFIG_DIR, obj.UserName)
	dutils.WriteKeyToKeyFile(file, "User", "Icon", icon)
	obj.addHistoryIcon(icon)
}

func (obj *User) setPropBackgroundFile(bg string) {
	if len(bg) < 1 {
		return
	}

	file := path.Join(USER_CONFIG_DIR, obj.UserName)
	dutils.WriteKeyToKeyFile(file, "User", "Background", bg)
}

func (obj *User) setPropAutomaticLogin(auto bool) {
	if auto {
		setAutomaticLogin(obj.UserName)
	} else {
		setAutomaticLogin("")
	}
}

func (obj *User) setPropAccountType(acctype int32) {
	t := obj.getPropAccountType()
	switch acctype {
	case ACCOUNT_TYPE_ADMINISTACTOR:
		if t != ACCOUNT_TYPE_ADMINISTACTOR {
			addUserToAdmList(obj.UserName)
		}
	case ACCOUNT_TYPE_STANDARD:
		if t == ACCOUNT_TYPE_ADMINISTACTOR {
			deleteUserFromAdmList(obj.UserName)
		}
	}
}

func (obj *User) setPropLocked(locked bool) {
	args := []string{}

	if locked {
		args = append(args, "-L")
		args = append(args, obj.UserName)
	} else {
		args = append(args, "-U")
		args = append(args, obj.UserName)
	}
	execCommand(CMD_USERMOD, args)
}
