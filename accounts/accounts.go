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
	dutils "pkg.linuxdeepin.com/lib/utils"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"strconv"
	"strings"
)

type Manager struct {
	UserList    []string
	AllowGuest  bool
	GuestIcon   string
	pathUserMap map[string]*User

	listWatcher *fsnotify.Watcher
	infoWatcher *fsnotify.Watcher
	listQuit    chan bool
	infoQuit    chan bool

	UserAdded   func(string)
	UserDeleted func(string)
}

var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = newManager()
	}

	return _manager
}

func newManager() *Manager {
	obj := &Manager{}

	var err error
	if obj.listWatcher, err = fsnotify.NewWatcher(); err != nil {
		logger.Error("New User List Watcher Failed:", err)
		panic(err)
	}
	if obj.infoWatcher, err = fsnotify.NewWatcher(); err != nil {
		logger.Error("New User Info Watcher Failed:", err)
		panic(err)
	}

	obj.listQuit = make(chan bool)
	obj.infoQuit = make(chan bool)
	obj.pathUserMap = make(map[string]*User)
	obj.updatePropUserList(getUserList())
	obj.updatePropAllowGuest(isAllowGuest())

	obj.watchUserListFile()
	obj.watchUserInfoFile()
	go obj.handleUserListChanged()
	go obj.handleUserInfoChanged()

	return obj
}

func getUserInfoList() []UserInfo {
	contents, err := ioutil.ReadFile(ETC_PASSWD)
	if err != nil {
		logger.Errorf("ReadFile '%s' failed: %s", ETC_PASSWD, err)
		panic(err)
	}

	infos := []UserInfo{}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		strs := strings.Split(line, ":")

		/* len of each line in /etc/passwd by spliting ':' is 7 */
		if len(strs) != PASSWD_SPLIT_LEN {
			continue
		}

		info := newUserInfo(strs[0], strs[2], strs[3],
			strs[5], strs[6])
		if checkUserIsHuman(&info) {
			infos = append(infos, info)
		}
	}

	return infos
}

func newUserInfo(name, uid, gid, home, shell string) UserInfo {
	info := UserInfo{}

	info.Name = name
	info.Uid = uid
	info.Gid = gid
	info.Home = home
	info.Shell = shell
	info.Path = USER_MANAGER_PATH + uid

	return info
}

func checkUserIsHuman(info *UserInfo) bool {
	if info.Name == "root" {
		return false
	}

	shells := strings.Split(info.Shell, "/")
	tmpShell := shells[len(shells)-1]
	if SHELL_END_FALSE == tmpShell ||
		SHELL_END_NOLOGIN == tmpShell {
		return false
	}

	if !detectedViaShadowFile(info) {
		id, _ := strconv.ParseInt(info.Uid, 10, 64)
		if id < 1000 {
			return false
		}
	}

	return true
}

func detectedViaShadowFile(info *UserInfo) bool {
	contents, err := ioutil.ReadFile(ETC_SHADOW)
	if err != nil {
		logger.Errorf("ReadFile '%s' failed: %s", ETC_SHADOW, err)
		panic(err)
	}

	isHuman := false
	info.Locked = false
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		strs := strings.Split(line, ":")
		if len(strs) != SHADOW_SPLIT_LEN {
			continue
		}

		if strs[0] != info.Name {
			continue
		}
		pw := strs[1]
		//加盐密码最短为13
		if len(pw) < 13 {
			break
		}

		if pw[0] == '!' {
			info.Locked = true
		}

		isHuman = true
	}

	return isHuman
}

func getUserList() []string {
	infos := getUserInfoList()
	list := []string{}

	for _, info := range infos {
		list = append(list, info.Path)
	}

	return list
}

func isAllowGuest() bool {
	if v, ok := dutils.ReadKeyFromKeyFile(ACCOUNT_CONFIG_FILE,
		ACCOUNT_GROUP_KEY, ACCOUNT_KEY_GUEST, true); !ok {
		dutils.WriteKeyToKeyFile(ACCOUNT_CONFIG_FILE,
			ACCOUNT_GROUP_KEY, ACCOUNT_KEY_GUEST, false)

		return false
	} else {
		if ret, ok := v.(bool); ok {
			return ret
		}
	}

	return false
}

func getUserInfoByPath(path string) (UserInfo, bool) {
	infos := getUserInfoList()

	for _, info := range infos {
		if path == info.Path {
			return info, true
		}
	}

	return UserInfo{}, false
}

func getUserInfoByName(name string) (UserInfo, bool) {
	infos := getUserInfoList()

	for _, info := range infos {
		if name == info.Name {
			return info, true
		}
	}

	return UserInfo{}, false
}
