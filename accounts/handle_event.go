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

package main

import (
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

	obj.listWatcher.Watch(ETC_PASSWD)
}

func (obj *Manager) removeUserListFileWatch() {
	if obj.listWatcher == nil {
		return
	}

	obj.listWatcher.RemoveWatch(ETC_PASSWD)
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
	if objUtil.IsFileExist(ETC_LIGHTDM_CONFIG) {
		obj.infoWatcher.Watch(ETC_LIGHTDM_CONFIG)
	}
	if objUtil.IsFileExist(ETC_GDM_CONFIG) {
		obj.infoWatcher.Watch(ETC_GDM_CONFIG)
	}
	if objUtil.IsFileExist(ETC_KDM_CONFIG) {
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
	if objUtil.IsFileExist(ETC_LIGHTDM_CONFIG) {
		obj.infoWatcher.RemoveWatch(ETC_LIGHTDM_CONFIG)
	}
	if objUtil.IsFileExist(ETC_GDM_CONFIG) {
		obj.infoWatcher.RemoveWatch(ETC_GDM_CONFIG)
	}
	if objUtil.IsFileExist(ETC_KDM_CONFIG) {
		obj.infoWatcher.RemoveWatch(ETC_KDM_CONFIG)
	}
}

func (obj *Manager) handleUserListChanged() {
	for {
		select {
		case ev := <-obj.listWatcher.Event:
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
				if !_handleFlag {
					var mutex sync.Mutex
					mutex.Lock()
					_handleFlag = true
					list, ret := compareStrList(obj.UserList, getUserList())
					switch ret {
					case 1:
						obj.setPropUserList(getUserList())
						obj.updateAllUserInfo()
						for _, v := range list {
							obj.UserAdded(v)
						}
					case -1:
						obj.setPropUserList(getUserList())
						obj.updateAllUserInfo()
						for _, v := range list {
							obj.UserDeleted(v)
						}
					}
					_handleFlag = false
					mutex.Unlock()
				}
			}
		}
	}
}

func (obj *Manager) handleUserInfoChanged() {
	for {
		select {
		case ev := <-obj.infoWatcher.Event:
			if ev == nil {
				break
			}

			if ok, _ := regexp.MatchString(`\.swa?px?$`, ev.Name); ok {
				break
			}

			logger.Info("User Info Event:", ev)
			ok1, _ := regexp.MatchString(ETC_GROUP, ev.Name)
			ok2, _ := regexp.MatchString(ETC_SHADOW, ev.Name)
			//ok3, _ := regexp.MatchString(ICON_SYSTEM_DIR, ev.Name)
			//ok4, _ := regexp.MatchString(ICON_LOCAL_DIR, ev.Name)
			ok5, _ := regexp.MatchString(ETC_LIGHTDM_CONFIG, ev.Name)
			ok6, _ := regexp.MatchString(ETC_GDM_CONFIG, ev.Name)
			ok7, _ := regexp.MatchString(ETC_KDM_CONFIG, ev.Name)

			if ok1 || ok2 || ok5 || ok6 || ok7 {
				if ev.IsDelete() {
					obj.removeUserInfoFileWatch()
					obj.watchUserInfoFile()
				}
			}

			if !_handleFlag {
				var mutex sync.Mutex
				mutex.Lock()
				_handleFlag = true
				obj.updateAllUserInfo()
				_handleFlag = false
				mutex.Unlock()
			}
		}
	}
}
