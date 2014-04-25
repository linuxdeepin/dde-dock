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
        "os"
        "regexp"
        //"time"
)

func (op *AccountManager) listenUserListChanged() {
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
                logObject.Warningf("New Watcher Failed:%v", err)
                //panic(err)
                return
        }

        err = watcher.Watch(ETC_PASSWD)
        if err != nil {
                logObject.Warningf("Watch File '%s' Failed: %v", ETC_PASSWD, err)
                //panic(err)
                return
        }

        go func() {
                defer watcher.Close()
                for {
                        select {
                        case ev := <-watcher.Event:
                                if ev == nil {
                                        break
                                }

                                if ok, _ := regexp.MatchString(`\.swa?px?$`,
                                        ev.Name); ok {
                                        break
                                }
                                logObject.Info(ev)
                                if ev.IsDelete() {
                                        watcher.Watch(ETC_PASSWD)
                                } else if ev.IsCreate() {
                                        break
                                } else {
                                        op.emitUserListChanged()
                                }
                                //case err := <-watcher.Error:
                                //logObject.Warningf("Watch Error:%v", err)
                        }
                }
        }()
}

func (op *UserManager) listenUserInfoChanged(filename string) {
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
                logObject.Warningf("New Watcher Failed:%v", err)
                //panic(err)
                return
        }

        err = watcher.Watch(filename)
        if err != nil {
                logObject.Warningf("Watch '%s' failed: %s", filename, err)
                //panic(err)
                return
        }

        go func() {
                defer watcher.Close()
                for {
                        select {
                        case ev := <-watcher.Event:
                                if ev == nil {
                                        break
                                }

                                if ok, _ := regexp.MatchString(`\.swa?px?$`,
                                        ev.Name); ok {
                                        break
                                }
                                if ev.IsDelete() {
                                        watcher.Watch(filename)
                                } else if ev.IsCreate() {
                                        break
                                } else {
                                        op.updateUserInfo()
                                }
                                //case err := <-watcher.Error:
                                //logObject.Warningf("Watch Error:%v", err)
                        }
                }
        }()
}

func (op *UserManager) listenIconListChanged(filename string) {
        if ok := opUtils.IsFileExist(filename); !ok {
                if err := os.MkdirAll(filename, 0755); err != nil {
                        return
                }
        }
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
                logObject.Warningf("New Watcher Failed:%v", err)
                //panic(err)
                return
        }

        err = watcher.Watch(filename)
        if err != nil {
                logObject.Warningf("Watch '%s' failed: %s", filename, err)
                //panic(err)
                return
        }

        go func() {
                defer watcher.Close()
                for {
                        select {
                        case ev := <-watcher.Event:
                                if ev == nil {
                                        break
                                }

                                if ok, _ := regexp.MatchString(`\.swa?px?$`,
                                        ev.Name); ok {
                                        break
                                }
                                logObject.Info("Icon List Event:", ev)
                                op.setPropName("IconList")
                                //case err := <-watcher.Error:
                                //logObject.Warningf("Watch Error:%v", err)
                        }
                }
        }()
}

func (op *AccountManager) emitUserListChanged() {
        infos := getUserInfoList()
        destList := []string{}
        for _, info := range infos {
                path := USER_MANAGER_PATH + info.Uid
                destList = append(destList, path)
        }
        list, ret := compareStrList(op.UserList, destList)
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
                        op.UserDeleted(v)
                }
        }
}

func compareStrList(src, dest []string) ([]string, int) {
        sl := len(src)
        dl := len(dest)

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
