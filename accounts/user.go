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
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
)

type User struct {
	Uid            string
	Gid            string
	UserName       string
	HomeDir        string
	Shell          string
	IconFile       string
	BackgroundFile string
	AutomaticLogin bool
	AccountType    int32
	Locked         bool
	LoginTime      uint64
	HistoryIcons   []string
	IconList       []string
	objectPath     string
	watcher        *fsnotify.Watcher
	watchQuit      chan bool
	endFlag        chan bool
	endWatchFlag   chan bool
}

func addUserToAdmList(name string) {
	tmps := []string{}
	tmps = append(tmps, "-a")
	tmps = append(tmps, name)
	tmps = append(tmps, "sudo")
	go execCommand(CMD_GPASSWD, tmps)
}

func deleteUserFromAdmList(name string) {
	tmps := []string{}
	tmps = append(tmps, "-d")
	tmps = append(tmps, name)
	tmps = append(tmps, "sudo")
	go execCommand(CMD_GPASSWD, tmps)
}

func getRandUserIcon() string {
	list := getIconList(ICON_SYSTEM_DIR)
	l := len(list)
	if l <= 0 {
		return ""
	}

	index := rand.Int31n(int32(l))
	return list[index]
}

func getIconList(dir string) []string {
	iconfd, err := os.Open(dir)
	if err != nil {
		logger.Errorf("Open '%s' failed: %v",
			dir, err)
		return []string{}
	}

	names, _ := iconfd.Readdirnames(0)
	list := []string{}
	for _, v := range names {
		if strings.Contains(v, "guest") {
			continue
		}

		tmp := strings.ToLower(v)
		//ok, _ := regexp.MatchString(`jpe?g$||png$||gif$`, tmp)
		ok1, _ := regexp.MatchString(`\.jpe?g$`, tmp)
		ok2, _ := regexp.MatchString(`\.png$`, tmp)
		ok3, _ := regexp.MatchString(`\.gif$`, tmp)
		if ok1 || ok2 || ok3 {
			list = append(list, path.Join(dir, v))
		}
	}

	return list
}

func getAdministratorList() []string {
	contents, err := ioutil.ReadFile(ETC_GROUP)
	if err != nil {
		logger.Errorf("ReadFile '%s' failed: %s", ETC_PASSWD, err)
		panic(err)
	}

	list := ""
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		strs := strings.Split(line, ":")
		if len(strs) != GROUP_SPLIT_LEN {
			continue
		}

		if strs[0] == "sudo" {
			list = strs[3]
			break
		}
	}

	return strings.Split(list, ",")
}

func setAutomaticLogin(name string) {
	dsp := getDefaultDisplayManager()
	logger.Infof("Set %s Auto For: %s", name, dsp)
	switch dsp {
	case "lightdm":
		if objUtil.IsFileExist(ETC_LIGHTDM_CONFIG) {
			objUtil.WriteKeyToKeyFile(ETC_LIGHTDM_CONFIG,
				LIGHTDM_AUTOLOGIN_GROUP,
				LIGHTDM_AUTOLOGIN_USER,
				name)
		}
	case "gdm":
		if objUtil.IsFileExist(ETC_GDM_CONFIG) {
			objUtil.WriteKeyToKeyFile(ETC_GDM_CONFIG,
				GDM_AUTOLOGIN_GROUP,
				GDM_AUTOLOGIN_USER,
				name)
		}
	case "kdm":
		if objUtil.IsFileExist(ETC_KDM_CONFIG) {
			objUtil.WriteKeyToKeyFile(ETC_KDM_CONFIG,
				KDM_AUTOLOGIN_GROUP,
				KDM_AUTOLOGIN_ENABLE,
				true)
			objUtil.WriteKeyToKeyFile(ETC_KDM_CONFIG,
				KDM_AUTOLOGIN_GROUP,
				KDM_AUTOLOGIN_USER,
				name)
		} else if objUtil.IsFileExist(USER_KDM_CONFIG) {
			objUtil.WriteKeyToKeyFile(ETC_KDM_CONFIG,
				KDM_AUTOLOGIN_GROUP,
				KDM_AUTOLOGIN_ENABLE,
				true)
			objUtil.WriteKeyToKeyFile(USER_KDM_CONFIG,
				KDM_AUTOLOGIN_GROUP,
				KDM_AUTOLOGIN_USER,
				name)
		}
	default:
		logger.Error("No support display manager")
	}
}

func isAutoLogin(username string) bool {
	dsp := getDefaultDisplayManager()
	//logger.Info("Display: ", dsp)

	switch dsp {
	case "lightdm":
		if objUtil.IsFileExist(ETC_LIGHTDM_CONFIG) {
			v, ok := objUtil.ReadKeyFromKeyFile(ETC_LIGHTDM_CONFIG,
				LIGHTDM_AUTOLOGIN_GROUP,
				LIGHTDM_AUTOLOGIN_USER,
				"")
			//logger.Info("AutoUser: ", v.(string))
			//logger.Info("UserName: ", username)
			if ok && v.(string) == username {
				return true
			}
		}
	case "gdm":
		if objUtil.IsFileExist(ETC_GDM_CONFIG) {
			v, ok := objUtil.ReadKeyFromKeyFile(ETC_GDM_CONFIG,
				GDM_AUTOLOGIN_GROUP,
				GDM_AUTOLOGIN_USER,
				"")
			if ok && v.(string) == username {
				return true
			}
		}
	case "kdm":
		if objUtil.IsFileExist(ETC_KDM_CONFIG) {
			v, ok := objUtil.ReadKeyFromKeyFile(ETC_KDM_CONFIG,
				KDM_AUTOLOGIN_GROUP,
				KDM_AUTOLOGIN_USER,
				"")
			if ok && v.(string) == username {
				return true
			}
		} else if objUtil.IsFileExist(USER_KDM_CONFIG) {
			v, ok := objUtil.ReadKeyFromKeyFile(USER_KDM_CONFIG,
				KDM_AUTOLOGIN_GROUP,
				KDM_AUTOLOGIN_USER,
				"")
			if ok && v.(string) == username {
				return true
			}
		}
	}

	return false
}

func getDefaultDisplayManager() string {
	contents, err := ioutil.ReadFile(ETC_DISPLAY_MANAGER)
	if err != nil {
		logger.Errorf("ReadFile '%s' failed: %s",
			ETC_DISPLAY_MANAGER, err)
		panic(err)
	}

	tmp := ""
	for _, b := range contents {
		if b == '\n' {
			tmp += ""
			continue
		}
		tmp += string(b)
	}

	return path.Base(tmp)
}

func (user *User) updateProps() {
	info, _ := getUserInfoByPath(user.objectPath)
	user.updatePropUserName(info.Name)
	user.updatePropHomeDir(info.Home)
	user.updatePropShell(info.Shell)
	user.updatePropLocked(info.Locked)
	user.updatePropAutomaticLogin(user.getPropAutomaticLogin())
	user.updatePropAccountType(user.getPropAccountType())
	user.updatePropIconList(user.getPropIconList())
	user.updatePropIconFile(user.getPropIconFile())
	user.updatePropBackgroundFile(user.getPropBackgroundFile())
	user.updatePropHistoryIcons(user.getPropHistoryIcons())
}

func newUser(path string) *User {
	info, ok := getUserInfoByPath(path)
	if !ok {
		return nil
	}

	obj := &User{}
	obj.objectPath = info.Path
	obj.Uid = info.Uid
	obj.Gid = info.Gid

	obj.watchQuit = make(chan bool)
	obj.endFlag = make(chan bool)
	obj.endWatchFlag = make(chan bool)

	var err error
	if obj.watcher, err = fsnotify.NewWatcher(); err != nil {
		logger.Error("New watcher in newUser failed:", err)
		panic(err)
	}
	obj.watchUserConfig()
	go obj.handUserConfigChanged()

	go obj.listenWatchQuit()

	return obj
}
