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
	dutils "dlib/utils"
	"github.com/howeyc/fsnotify"
	"regexp"
	"sync"
)

var _handleFlag = false

func (obj *Manager) watchUserListFile() {
	var err error

	if obj.listWatcher == nil {
		if obj.listWatcher, err = fsnotify.NewWatcher(); err != nil {
			logger.Error("New User List Watcher Failed:", err)
			panic(err)
		}
	}

	obj.listWatcher.Watch(ETC_SHADOW)
	obj.listWatcher.Watch(ACCOUNT_CONFIG_FILE)
}

func (obj *Manager) removeUserListFileWatch() {
	if obj.listWatcher == nil {
		return
	}

	obj.listWatcher.RemoveWatch(ETC_SHADOW)
	obj.listWatcher.RemoveWatch(ACCOUNT_CONFIG_FILE)
}

func (obj *Manager) watchUserInfoFile() {
	var err error

	if obj.infoWatcher == nil {
		if obj.infoWatcher, err = fsnotify.NewWatcher(); err != nil {
			logger.Error("New User Info Watcher Failed:", err)
			panic(err)
		}
	}

	obj.infoWatcher.Watch(ETC_GROUP)
	obj.infoWatcher.Watch(ETC_SHADOW)
	obj.infoWatcher.Watch(ICON_SYSTEM_DIR)
	obj.infoWatcher.Watch(ICON_LOCAL_DIR)
	if dutils.IsFileExist(ETC_LIGHTDM_CONFIG) {
		obj.infoWatcher.Watch(ETC_LIGHTDM_CONFIG)
	}
	if dutils.IsFileExist(ETC_GDM_CONFIG) {
		obj.infoWatcher.Watch(ETC_GDM_CONFIG)
	}
	if dutils.IsFileExist(ETC_KDM_CONFIG) {
		obj.infoWatcher.Watch(ETC_KDM_CONFIG)
	}
}

func (obj *Manager) removeUserInfoFileWatch() {
	if obj.infoWatcher == nil {
		return
	}

	obj.infoWatcher.RemoveWatch(ETC_GROUP)
	obj.infoWatcher.RemoveWatch(ETC_SHADOW)
	obj.infoWatcher.RemoveWatch(ICON_SYSTEM_DIR)
	obj.infoWatcher.RemoveWatch(ICON_LOCAL_DIR)
	if dutils.IsFileExist(ETC_LIGHTDM_CONFIG) {
		obj.infoWatcher.RemoveWatch(ETC_LIGHTDM_CONFIG)
	}
	if dutils.IsFileExist(ETC_GDM_CONFIG) {
		obj.infoWatcher.RemoveWatch(ETC_GDM_CONFIG)
	}
	if dutils.IsFileExist(ETC_KDM_CONFIG) {
		obj.infoWatcher.RemoveWatch(ETC_KDM_CONFIG)
	}
}

func (obj *Manager) handleUserListChanged() {
	for {
		select {
		case <-obj.listQuit:
			return
		case ev, ok := <-obj.listWatcher.Event:
			if !ok {
				if obj.listWatcher != nil {
					obj.removeUserListFileWatch()
				}
				obj.watchUserListFile()
				break
			}
			if ev == nil {
				break
			}

			if ok, _ := regexp.MatchString(`\.swa?px?$`, ev.Name); ok {
				break
			}

			logger.Info("User List Event:", ev)

			if ev.IsDelete() {
				obj.removeUserListFileWatch()
				obj.watchUserListFile()
			} else if ev.IsModify() {
				ok1, _ := regexp.MatchString(ACCOUNT_CONFIG_FILE, ev.Name)
				if ok1 {
					obj.updatePropAllowGuest(isAllowGuest())
					break
				}
				if !_handleFlag {
					var mutex sync.Mutex
					mutex.Lock()
					_handleFlag = true
					list, ret := compareStrList(obj.UserList, getUserList())
					switch ret {
					case 1:
						obj.updatePropUserList(getUserList())
						obj.updateAllUserInfo()
						for _, v := range list {
							obj.UserAdded(v)
						}
					case -1:
						obj.updatePropUserList(getUserList())
						obj.updateAllUserInfo()
						for _, v := range list {
							obj.UserDeleted(v)
						}
					}
					_handleFlag = false
					mutex.Unlock()
				}
			}
		case err, ok := <-obj.listWatcher.Error:
			if !ok || err != nil {
				if obj.listWatcher != nil {
					obj.removeUserListFileWatch()
				}
				obj.watchUserListFile()
				break
			}
		}
	}
}

func (obj *Manager) handleUserInfoChanged() {
	for {
		select {
		case <-obj.infoQuit:
			return
		case ev, ok := <-obj.infoWatcher.Event:
			if !ok {
				if obj.infoWatcher != nil {
					obj.removeUserInfoFileWatch()
				}
				obj.watchUserInfoFile()
				break
			}
			if ev == nil {
				break
			}

			if ok, _ := regexp.MatchString(`\.swa?px?$`, ev.Name); ok {
				break
			}

			logger.Info("User Info Event:", ev)
			ok3, _ := regexp.MatchString(ICON_SYSTEM_DIR, ev.Name)
			ok4, _ := regexp.MatchString(ICON_LOCAL_DIR, ev.Name)
			if ok3 || ok4 {
				if !_handleFlag {
					var mutex sync.Mutex
					mutex.Lock()
					_handleFlag = true
					obj.updateAllUserInfo()
					_handleFlag = false
					mutex.Unlock()
				}
				break
			}

			if ev.IsDelete() {
				obj.removeUserInfoFileWatch()
				obj.watchUserInfoFile()
				break
			}

			ok1, _ := regexp.MatchString(ETC_GROUP, ev.Name)
			if ok1 {
				for _, u := range obj.pathUserMap {
					u.updatePropAccountType(u.getPropAccountType())
				}
				break
			}

			ok2, _ := regexp.MatchString(ETC_SHADOW, ev.Name)
			if ok2 {
				infos := getUserInfoList()
				for _, info := range infos {
					u, ok := obj.pathUserMap[obj.FindUserByName(info.Name)]
					if !ok {
						continue
					}
					u.updatePropLocked(info.Locked)
				}
				break
			}

			ok5, _ := regexp.MatchString(ETC_LIGHTDM_CONFIG, ev.Name)
			ok6, _ := regexp.MatchString(ETC_GDM_CONFIG, ev.Name)
			ok7, _ := regexp.MatchString(ETC_KDM_CONFIG, ev.Name)
			if ok5 || ok6 || ok7 {
				for _, u := range obj.pathUserMap {
					u.updatePropAutomaticLogin(u.getPropAutomaticLogin())
				}
				break
			}
		case err, ok := <-obj.infoWatcher.Error:
			if !ok || err != nil {
				if obj.infoWatcher != nil {
					obj.removeUserInfoFileWatch()
				}
				obj.watchUserInfoFile()
				break
			}
		}
	}
}

func (obj *User) watchUserConfig() {
	if obj.watcher == nil {
		var err error
		if obj.watcher, err = fsnotify.NewWatcher(); err != nil {
			logger.Error("New watcher in newUser failed:", err)
			panic(err)
		}
	}

	obj.watcher.Watch(USER_CONFIG_FILE + obj.UserName)
}

func (obj *User) removeUserConfigWatch() {
	if obj.watcher == nil {
		return
	}

	obj.watcher.RemoveWatch(USER_CONFIG_FILE + obj.UserName)
}

func (obj *User) handUserConfigChanged() {
	for {
		select {
		case <-obj.quitFlag:
			return
		case ev, ok := <-obj.watcher.Event:
			if !ok {
				if obj.watcher != nil {
					obj.removeUserConfigWatch()
				}
				obj.watchUserConfig()
				break
			}
			if ev == nil {
				break
			}

			if ok, _ := regexp.MatchString(`\.swa?px?$`, ev.Name); ok {
				break
			}

			logger.Info("User Config Event:", ev)
			if ev.IsDelete() {
				obj.removeUserConfigWatch()
				obj.watchUserConfig()
			} else if ev.IsModify() {
				obj.updatePropIconList(obj.getPropIconList())
				obj.updatePropIconFile(obj.getPropIconFile())
				obj.updatePropBackgroundFile(obj.getPropBackgroundFile())
				obj.updatePropHistoryIcons(obj.getPropHistoryIcons())
			}
		case err, ok := <-obj.watcher.Error:
			if !ok || err != nil {
				if obj.watcher != nil {
					obj.removeUserConfigWatch()
				}
				obj.watchUserConfig()
				break
			}
		}
	}
}
