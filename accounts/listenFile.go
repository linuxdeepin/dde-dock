/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package main

import (
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"regexp"
	"strings"
	//"time"
)

var (
	preGroupLine = 0
)

func watchAccountFiles() {
	mutex.Lock()
	defer mutex.Unlock()
	var err error

	println("BEGIN Watch Account All File")
	if watchAct == nil {
		watchAct, err = fsnotify.NewWatcher()
		if err != nil {
			logObject.Warning("New Watch Account Failed: ", err)
			logObject.Fatal(err)
			return
		}
	}

	// User List
	//watchAct.Watch(ETC_PASSWD)
	watchAct.Watch(ETC_SHADOW)
	watchAct.Watch(ETC_GROUP)
	println("END Watch All Account File")
}

func watchUserFiles() {
	mutex.Lock()
	defer mutex.Unlock()
	var err error

	println("BEGIN Watch User All File")
	if watchUser == nil {
		watchUser, err = fsnotify.NewWatcher()
		if err != nil {
			logObject.Warning("New Watch User Failed: ", err)
			logObject.Fatal(err)
			return
		}
	}

	// User Info
	//watchUser.Watch(ETC_GROUP)
	//watchUser.Watch(ETC_SHADOW)
	watchUser.Watch(ICON_SYSTEM_DIR)
	watchUser.Watch(ICON_LOCAL_DIR)
	if opUtils.IsFileExist(ETC_LIGHTDM_CONFIG) {
		watchUser.Watch(ETC_LIGHTDM_CONFIG)
	}
	if opUtils.IsFileExist(ETC_GDM_CONFIG) {
		watchUser.Watch(ETC_GDM_CONFIG)
	}
	if opUtils.IsFileExist(ETC_KDM_CONFIG) {
		watchUser.Watch(ETC_KDM_CONFIG)
	}
	println("END Watch All User File")
}

func removeWatchAccountFiles() {
	mutex.Lock()
	defer mutex.Unlock()
	if watchAct == nil {
		return
	}

	println("BEGIN Remove All Account File Watcher...")
	// User List
	//watchAct.RemoveWatch(ETC_PASSWD)
	watchAct.RemoveWatch(ETC_GROUP)
	watchAct.RemoveWatch(ETC_SHADOW)
	println("END Remove All Account File Watcher...")
}

func removeWatchUserFiles() {
	mutex.Lock()
	defer mutex.Unlock()
	if watchUser == nil {
		return
	}

	println("BEGIN Remove All User File Watcher...")
	// User Info
	//watchUser.RemoveWatch(ETC_GROUP)
	//watchUser.RemoveWatch(ETC_SHADOW)
	watchUser.RemoveWatch(ICON_SYSTEM_DIR)
	watchUser.RemoveWatch(ICON_LOCAL_DIR)
	if opUtils.IsFileExist(ETC_LIGHTDM_CONFIG) {
		watchUser.RemoveWatch(ETC_LIGHTDM_CONFIG)
	}
	if opUtils.IsFileExist(ETC_GDM_CONFIG) {
		watchUser.RemoveWatch(ETC_GDM_CONFIG)
	}
	if opUtils.IsFileExist(ETC_KDM_CONFIG) {
		watchUser.RemoveWatch(ETC_KDM_CONFIG)
	}
	println("END Remove All User File Watcher...")
}

func (op *AccountManager) listenUserListChanged() {
	for {
		select {
		case ev := <-watchAct.Event:
			if ev == nil {
				break
			}

			if ok, _ := regexp.MatchString(`\.swa?px?$`,
				ev.Name); ok {
				break
			}

			ok1, _ := regexp.MatchString(ETC_GROUP, ev.Name)
			//ok2, _ := regexp.MatchString(ETC_SHADOW, ev.Name)
			logObject.Info(ev)
			if ev.IsDelete() {
				removeWatchAccountFiles()
				watchAccountFiles()
			} else if ev.IsModify() {
				println("UPDATE User List")
				if ok1 {
					contents, _ := ioutil.ReadFile(ETC_GROUP)
					lines := strings.Split(string(contents), "\n")
					curGroupLine := len(lines)
					if preGroupLine != curGroupLine {
						preGroupLine = curGroupLine
						break
					}
				}
				op.emitUserListChanged()
				println("UPDATE User List END....")
			}
		case err := <-watchAct.Error:
			logObject.Warningf("Watch Error:%v", err)
		}
	}
}

func (op *UserManager) listenUserInfoChanged() {
	for {
		select {
		case ev := <-watchUser.Event:
			if ev == nil {
				break
			}

			if ok, _ := regexp.MatchString(`\.swa?px?$`,
				ev.Name); ok {
				break
			}
			logObject.Info(ev)

			//ok1, _ := regexp.MatchString(ETC_GROUP, ev.Name)
			//ok2, _ := regexp.MatchString(ETC_SHADOW, ev.Name)
			ok3, _ := regexp.MatchString(ICON_SYSTEM_DIR, ev.Name)
			ok4, _ := regexp.MatchString(ICON_LOCAL_DIR, ev.Name)
			ok5, _ := regexp.MatchString(ETC_LIGHTDM_CONFIG, ev.Name)
			ok6, _ := regexp.MatchString(ETC_GDM_CONFIG, ev.Name)
			ok7, _ := regexp.MatchString(ETC_KDM_CONFIG, ev.Name)

			if ok5 || ok6 || ok7 {
				if ev.IsDelete() {
					removeWatchUserFiles()
					watchUserFiles()
				} else if ev.IsModify() {
					println("UPDATE User Info")
					op.updateUserInfo()
					println("UPDATE User Info END...")
				}
			}

			if ok3 || ok4 {
				logObject.Info("Icon List Event:", ev)
				op.setPropName("IconList")
			}
		case err := <-watchUser.Error:
			logObject.Warningf("Watch Error:%v", err)
		}
	}
}

func (op *AccountManager) emitUserListChanged() {
	logObject.Info("EMIT User List Change Signal...")
	infos := getUserInfoList()
	destList := []string{}
	for _, info := range infos {
		path := USER_MANAGER_PATH + info.Uid
		destList = append(destList, path)
	}
	list, ret := compareStrList(op.UserList, destList)
	//logObject.Infof("***** compare ret: %v ------ %d", list, ret)
	switch ret {
	case 1:
		updateUserList()
		//go func() {
		//<-time.After(time.Millisecond * 500)
		op.setPropName("UserList")
		//}()
		for _, v := range list {
			op.UserAdded(v)
		}
	case -1:
		updateUserList()
		//go func() {
		//<-time.After(time.Millisecond * 500)
		op.setPropName("UserList")
		//}()
		for _, v := range list {
			//logObject.Info("======== User Deleted: ", v)
			op.UserDeleted(v)
		}
	}
	logObject.Info("EMIT User List Change Signal END...")
}

func compareStrList(src, dest []string) ([]string, int) {
	sl := len(src)
	dl := len(dest)

	//logObject.Info("--------Compare src: ", src)
	//logObject.Info("--------Compare dest: ", dest)
	tmp := []string{}
	if sl < dl {
		for i := 0; i < dl; i++ {
			j := 0
			for ; j < sl; j++ {
				if dest[i] == src[j] {
					break
				}
			}
			if j == sl {
				tmp = append(tmp, dest[i])
			}
		}
		return tmp, 1
	} else if sl > dl {
		for i := 0; i < sl; i++ {
			j := 0
			for ; j < dl; j++ {
				if src[i] == dest[j] {
					break
				}
			}
			if j == dl {
				tmp = append(tmp, src[i])
			}
		}
		return tmp, -1
	}

	return tmp, 0
}
